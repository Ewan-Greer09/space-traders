package handlers

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"space-traders/repository/mysql"
	"space-traders/service/views/register"
)

func (h *ViewHandler) RegisterPage(c echo.Context) error {
	return register.Page().Render(c.Request().Context(), c.Response())
}

type RegisterForm struct {
	Username        string `json:"username" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm-password" validate:"required"`
	Email           string `json:"email" validate:"required"`
	ApiKey          string `json:"api-key" validate:"required"`
}

func (r RegisterForm) Validate() error {
	v := validator.New()
	err := v.Struct(r)
	if err != nil {
		return err
	}

	if strings.Compare(r.Password, r.ConfirmPassword) != 0 {
		return errors.New("password and new password do not match")
	}

	return nil
}

func (h *ViewHandler) RegisterSubmit(c echo.Context) error {
	var form RegisterForm
	if err := c.Bind(&form); err != nil {
		return register.RegisterResponseError("Invalid form data").Render(c.Request().Context(), c.Response())
	}

	if form.Password != form.ConfirmPassword {
		return register.RegisterResponseError("Passwords do not match").Render(c.Request().Context(), c.Response())
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		return register.RegisterResponseError("There was an issue. Please wait a moment and try again").Render(c.Request().Context(), c.Response())
	}

	err = h.userDB.CreateUser(c.Request().Context(), mysql.CreateUserParams{
		UserUid:  sql.NullString{String: uuid.NewString(), Valid: true},
		Username: sql.NullString{String: form.Username, Valid: true},
		Password: sql.NullString{String: string(encryptedPassword), Valid: true},
		Email:    sql.NullString{String: form.Email, Valid: true},
	})
	if err != nil {
		return register.RegisterResponseError("Failed to create user").Render(c.Request().Context(), c.Response())
	}

	err = h.userDB.CreateAPIKey(c.Request().Context(), mysql.CreateAPIKeyParams{
		ApiKey:   sql.NullString{String: form.ApiKey, Valid: true},
		Username: sql.NullString{String: form.Username, Valid: true},
	})
	if err != nil {
		return register.RegisterResponseError("Failed to create API key").Render(c.Request().Context(), c.Response())
	}

	c.Response().Header().Add("HX-Redirect", "/login")
	return register.RegisterResponseSuccess("User created successfully").Render(c.Request().Context(), c.Response())
}
