package handlers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"space-traders/repository/postgres"
	"space-traders/service/views/components/login"
	"space-traders/service/views/components/register"
)

func (vh *ViewHandler) MountLoginRoutes(e *echo.Echo) {
	e.GET("/login", vh.GetLogin)
	e.GET("/register", vh.GetRegister)
	e.POST("/login", vh.HandleLogin)
	e.POST("/register", vh.HandleRegister)
	e.GET("/logout", vh.Logout)
}

func (vh *ViewHandler) GetLogin(c echo.Context) error {
	return login.Page().Render(c.Request().Context(), c.Response())
}

func (vh *ViewHandler) HandleLogin(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := vh.userDB.GetUserWithAPIKeyByUsername(c.Request().Context(), pgtype.Text{String: username, Valid: true})
	if err != nil {
		c.Logger().Error(err.Error())
		return login.LoginFailure().Render(c.Request().Context(), c.Response())
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(password))
	if err != nil {
		c.Logger().Error(err.Error())
		return login.LoginFailure().Render(c.Request().Context(), c.Response())
	}

	token, err := vh.generateUserJWT(username)
	if err != nil {
		c.Logger().Error("Failed to generate JWT")
		return login.LoginFailure().Render(c.Request().Context(), c.Response())
	}

	cookie := &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	}

	c.SetCookie(cookie)
	return login.LoginSuccess().Render(c.Request().Context(), c.Response())
}

func (vh *ViewHandler) Logout(c echo.Context) error {
	cookie := &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Expires:  time.Now().Add(-1 * time.Hour),
	}
	c.SetCookie(cookie)

	vh.Client.GetConfig().AddDefaultHeader("Authorization", "")

	return c.Redirect(http.StatusFound, "/")
}

func (vh *ViewHandler) GetRegister(c echo.Context) error {
	return register.Page().Render(c.Request().Context(), c.Response())
}

func (vh *ViewHandler) HandleRegister(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	apiKey := c.FormValue("api-key")

	if username == "" || password == "" || apiKey == "" {
		c.Logger().Error("Missing required fields")
		return register.RegisterFailure().Render(c.Request().Context(), c.Response())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.Logger().Error(err.Error())
		return register.RegisterFailure().Render(c.Request().Context(), c.Response())
	}

	u := postgres.CreateUserParams{
		UserUid:  pgtype.Text{String: uuid.NewString(), Valid: true},
		Username: pgtype.Text{String: username, Valid: true},
		Password: pgtype.Text{String: string(hashedPassword), Valid: true},
		Email:    pgtype.Text{String: "", Valid: false},
	}
	_, err = vh.userDB.CreateUser(c.Request().Context(), u)
	if err != nil {
		c.Logger().Error(err.Error())
		return register.RegisterFailure().Render(c.Request().Context(), c.Response())
	}

	err = vh.userDB.CreateAPIKey(c.Request().Context(), postgres.CreateAPIKeyParams{
		Key:      pgtype.Text{String: apiKey, Valid: true},
		Username: pgtype.Text{String: username, Valid: true},
	})
	if err != nil {
		c.Logger().Error(err.Error())
		return register.RegisterFailure().Render(c.Request().Context(), c.Response())
	}

	return register.RegisterSuccess().Render(c.Request().Context(), c.Response())
}

func (vh ViewHandler) generateUserJWT(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	token.Claims = claims

	secret := []byte(vh.cfg.JWT_SECRET)

	t, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return t, nil
}
