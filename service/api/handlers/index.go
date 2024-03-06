package handlers

import (
	"github.com/labstack/echo/v4"

	"space-traders/service/views/components/index"
)

func (vh *ViewHandler) MountIndexRoutes(e *echo.Echo) {
	e.GET("/", vh.GetIndex)
}

func (vh *ViewHandler) GetIndex(c echo.Context) error {
	return index.Page().Render(c.Request().Context(), c.Response())
}
