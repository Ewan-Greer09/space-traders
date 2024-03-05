package handlers

import (
	"context"
	"os"

	openAPI "github.com/UnseenBook/spacetraders-go-sdk"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"

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
	userDB *pgx.Conn
}

func NewViewHandler() *ViewHandler {
	cfg := openAPI.NewConfiguration()
	cfg.AddDefaultHeader("Content-Type", "application/json")
	cfg.AddDefaultHeader("Accept", "application/json")
	// cfg.AddDefaultHeader("Authorization", "Bearer "+os.Getenv("API_ACCESS_TOKEN"))

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	return &ViewHandler{
		Client: openAPI.NewAPIClient(cfg),
		userDB: conn,
	}
}

func (vh *ViewHandler) GetHeader(c echo.Context) error {
	resp, _, err := vh.Client.DefaultAPI.GetStatus(c.Request().Context()).Execute()
	if err != nil {
		c.Logger().Error(err.Error())
		return shared.Error(err).Render(c.Request().Context(), c.Response())
	}
	_, status := resp.GetStatusOk()

	return shared.Header(status).Render(c.Request().Context(), c.Response())
}

func (vh *ViewHandler) GetFooter(c echo.Context) error {
	return shared.Footer().Render(c.Request().Context(), c.Response())
}

func (vh *ViewHandler) AddKeyToReq() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session := c.Get("session")
			if session != nil {
				vh.Client.GetConfig().AddDefaultHeader("Authorization", "Bearer "+session.(string))
			}

			return next(c)
		}
	}
}
