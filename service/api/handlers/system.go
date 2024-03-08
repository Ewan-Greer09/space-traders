package handlers

import (
	"github.com/labstack/echo/v4"

	"space-traders/service/views/components/system"
)

func (vh *ViewHandler) MountSystemRoutes(e *echo.Echo) {
	e.GET("/system", vh.SystemPage)
}

func (vh *ViewHandler) SystemPage(c echo.Context) error {
	return system.Page().Render(c.Request().Context(), c.Response())
}
