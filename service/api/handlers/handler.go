package handlers

import (
	"context"

	openAPI "github.com/UnseenBook/spacetraders-go-sdk"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	"space-traders/repository/postgres"
	"space-traders/service/config"
	"space-traders/service/views/components/shared"
)

type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}
}

type ViewHandler struct {
	Client *openAPI.APIClient
	userDB *postgres.Queries
	cfg    *config.Config
}

func (vh *ViewHandler) MountSharedRoutes(e *echo.Echo) {
	e.GET("/com/header", vh.GetHeader)
	e.GET("/com/footer", vh.GetFooter)
}

func NewViewHandler(config *config.Config) *ViewHandler {
	cfg := openAPI.NewConfiguration()
	cfg.AddDefaultHeader("Content-Type", "application/json")
	cfg.AddDefaultHeader("Accept", "application/json")

	// create a connection pool
	pool, err := pgxpool.New(context.Background(), config.DATABASE_URL)
	if err != nil {
		panic(err)
	}

	return &ViewHandler{
		Client: openAPI.NewAPIClient(cfg),
		userDB: postgres.New(pool),
		cfg:    config,
	}
}

func (vh *ViewHandler) GetHeader(c echo.Context) error {
	vh.Client.GetConfig().AddDefaultHeader("Authorization", "")

	_, r, err := vh.Client.DefaultAPI.GetStatus(c.Request().Context()).Execute()
	if err != nil {
		c.Logger().Error(err.Error())
		return shared.Error(err).Render(c.Request().Context(), c.Response())
	}

	if r.StatusCode != 200 {
		return shared.Header(false).Render(c.Request().Context(), c.Response())
	}

	return shared.Header(true).Render(c.Request().Context(), c.Response())
}

func (vh *ViewHandler) GetFooter(c echo.Context) error {
	return shared.Footer().Render(c.Request().Context(), c.Response())
}

func (vh *ViewHandler) AddKeyToReq() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/login" || c.Path() == "/register" || c.Path() == "/com/header" || c.Path() == "/com/footer" {
				c.Logger().Info("Skipping auth for path: ", c.Path())
				return next(c)
			}

			cookie, err := c.Cookie("session")
			if err != nil || cookie.Value == "" {
				c.Logger().Error(err.Error())
				return next(c)
			}

			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					c.Logger().Error("Invalid token method")
					return nil, echo.NewHTTPError(401, "Invalid token")
				}

				return []byte(vh.cfg.JWT_SECRET), nil

			})
			if err != nil {
				c.Logger().Error(err.Error())
				next(c)
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				c.Logger().Error("Invalid token")
				return next(c)
			}

			user, err := vh.userDB.GetUserWithAPIKeyByUsername(c.Request().Context(), pgtype.Text{String: claims["username"].(string), Valid: true})
			if err != nil {
				c.Logger().Error(err.Error())
				return next(c)
			}

			c.Set("apiKey", user.Key.String) // save this in the context to facilitate user lookups using the apikey as pkey
			c.Set("username", claims["username"])

			vh.Client.GetConfig().AddDefaultHeader("Authorization", "Bearer "+user.Key.String)

			return next(c)
		}
	}
}

func (vh *ViewHandler) Favicon(c echo.Context) error {
	return c.File("service/views/favicon.ico")
}
