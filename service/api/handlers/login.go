package handlers

import (
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"

	"space-traders/service/views/components/login"
	"space-traders/service/views/components/register"
)

func (vh *ViewHandler) MountLoginRoutes(e *echo.Echo) {
	e.GET("/login", vh.GetLogin)
	e.POST("/login", vh.HandleLogin)
	e.POST("/login/username", vh.HandleLoginWithUsername)
	e.GET("/register", vh.GetRegister)

}

func (vh *ViewHandler) GetLogin(c echo.Context) error {
	return login.Page().Render(c.Request().Context(), c.Response())
}

func (vh *ViewHandler) HandleLogin(c echo.Context) error {
	key := c.FormValue("api-key")

	vh.Client.GetConfig().AddDefaultHeader("Authorization", "Bearer "+key)

	// validate that the key is valid
	_, _, err := vh.Client.AgentsAPI.GetMyAgent(c.Request().Context()).Execute()
	if err != nil {
		return login.LoginFailure().Render(c.Request().Context(), c.Response())
	}

	cookie := &http.Cookie{
		Name:     "session",
		Value:    key,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	}
	c.SetCookie(cookie)

	return login.LoginSuccess().Render(c.Request().Context(), c.Response())
}

func (vh *ViewHandler) HandleLoginWithUsername(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := vh.userDB.GetOneByUsername(c.Request().Context(), pgtype.Text{String: username})
	if err != nil {
		return login.LoginFailure().Render(c.Request().Context(), c.Response())
	}

	if user.Password.String != password {
		return login.LoginFailure().Render(c.Request().Context(), c.Response())
	}

	cookie := &http.Cookie{
		Name:     "session",
		Value:    user.ApiKey.String,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	}
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/")
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
