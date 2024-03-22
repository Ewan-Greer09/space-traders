package handlers

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"space-traders/repository/mysql"
	"space-traders/service/views/register"
)

func (h *ViewHandler) MountRegisterRoutes(e *echo.Echo) {
	e.GET("/register", h.RegisterPage)
	e.GET("/register/submit", h.RegisterSubmit)
}

func (h *ViewHandler) RegisterPage(c echo.Context) error {
	return register.Page().Render(c.Request().Context(), c.Response())
}

func (h *ViewHandler) RegisterSubmit(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirm-password")
	email := c.FormValue("email")
	apiKey := c.FormValue("api-key")

	if username == "" || password == "" || email == "" || apiKey == "" {
		return register.RegisterResponseError("All fields are required").Render(c.Request().Context(), c.Response())
	}

	if password != confirmPassword {
		return register.RegisterResponseError("Passwords do not match").Render(c.Request().Context(), c.Response())
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return register.RegisterResponseError("There was an issue. Please wait a moment and try again").Render(c.Request().Context(), c.Response())
	}

	err = h.userDB.CreateUser(c.Request().Context(), mysql.CreateUserParams{
		UserUid:  sql.NullString{String: uuid.NewString(), Valid: true},
		Username: sql.NullString{String: username, Valid: true},
		Password: sql.NullString{String: string(encryptedPassword), Valid: true},
		Email:    sql.NullString{String: email, Valid: true},
	})
	if err != nil {
		return register.RegisterResponseError("Failed to create user").Render(c.Request().Context(), c.Response())
	}

	err = h.userDB.CreateAPIKey(c.Request().Context(), mysql.CreateAPIKeyParams{
		ApiKey:   sql.NullString{String: apiKey, Valid: true},
		Username: sql.NullString{String: username, Valid: true},
	})
	if err != nil {
		return register.RegisterResponseError("Failed to create API key").Render(c.Request().Context(), c.Response())
	}

	c.Response().Header().Add("HX-Redirect", "http://localhost:3000/login")
	return register.RegisterResponseSuccess("User created successfully").Render(c.Request().Context(), c.Response())
}
