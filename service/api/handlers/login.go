package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"space-traders/service/views/login"
)

func (h *ViewHandler) LoginPage(c echo.Context) error {
	return login.Page().Render(c.Request().Context(), c.Response())
}

func (h *ViewHandler) LoginSubmit(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		return login.LoginResponseError("Username and password are required").Render(c.Request().Context(), c.Response())
	}

	user, err := h.userDB.GetUserWithAPIKeyByUsername(c.Request().Context(), sql.NullString{String: username, Valid: true})
	if err != nil {
		return login.LoginResponseError("Invalid username").Render(c.Request().Context(), c.Response())
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(password))
	if err != nil {
		return login.LoginResponseError("Invalid password").Render(c.Request().Context(), c.Response())
	}

	token, err := h.generateUserJWT(username)
	if err != nil {
		return login.LoginResponseError("Failed to generate token").Render(c.Request().Context(), c.Response())
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
	c.Response().Header().Set("HX-Redirect", "http://localhost:3000/")
	return login.LoginResponseSuccess("Login Successful").Render(c.Request().Context(), c.Response())
}

func (vh ViewHandler) generateUserJWT(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	token.Claims = claims

	secret := []byte(vh.cfg.JwtSecret)

	t, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return t, nil
}

func (h *ViewHandler) Logout(e echo.Context) error {
	cookie := &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Expires:  time.Now().Add(-1 * time.Hour),
	}

	e.SetCookie(cookie)
	return e.Redirect(http.StatusFound, "/")
}
