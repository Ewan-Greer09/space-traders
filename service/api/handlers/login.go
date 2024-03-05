package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"space-traders/service/views/components/login"
)

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
