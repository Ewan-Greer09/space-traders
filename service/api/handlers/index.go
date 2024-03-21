package handlers

import (
	"github.com/labstack/echo/v4"

	"space-traders/service/views/components/index"
	"space-traders/service/views/components/widgets"
)

func (vh *ViewHandler) MountIndexRoutes(e *echo.Echo) {
	e.GET("/", vh.GetIndex)
	e.GET("/widgets/navigation", vh.GetNavigationWidget)
}

func (vh *ViewHandler) GetIndex(c echo.Context) error {
	return index.Page().Render(c.Request().Context(), c.Response())
}

func (vh *ViewHandler) GetNavigationWidget(c echo.Context) error {
	return widgets.NavWidget().Render(c.Request().Context(), c.Response())
}
