package api

import (
	"github.com/labstack/echo/v4"

	"space-traders/service/api/handlers"
)

func MountRoutes(e *echo.Echo, h handlers.ViewHandler) {
	// Login
	e.GET("/login", h.LoginPage)
	e.GET("/login/submit", h.LoginSubmit)
	e.GET("/logout", h.Logout)

	// Registration
	e.GET("/register", h.RegisterPage)
	e.GET("/register/submit", h.RegisterSubmit)

	// Utilities
	e.GET("/favicon.ico", h.Favicon)
}
