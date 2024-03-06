package handlers

import (
	"github.com/labstack/echo/v4"

	"space-traders/service/views/components/fleet"
	"space-traders/service/views/components/shared"
	"space-traders/service/views/components/ship"
)

func (vh *ViewHandler) MountFleetRoutes(e *echo.Echo) {
	e.GET("/fleet", vh.GetFleet)
	e.GET("/fleet/list", vh.GetFleetList)
	e.GET("/fleet/ship/:symbol", vh.GetFleetShip)
}

func (vh *ViewHandler) GetFleet(c echo.Context) error {
	return fleet.Page().Render(c.Request().Context(), c.Response())
}

func (vh *ViewHandler) GetFleetList(c echo.Context) error {
	resp, _, err := vh.Client.FleetAPI.GetMyShips(c.Request().Context()).Execute()
	if err != nil {
		c.Logger().Error(err.Error())
		return shared.Error(err).Render(c.Request().Context(), c.Response())
	}

	return fleet.Fleet(*resp).Render(c.Request().Context(), c.Response())
}

func (vh *ViewHandler) GetFleetShip(c echo.Context) error {
	shipSymbol := c.Param("symbol")
	resp, _, err := vh.Client.FleetAPI.GetMyShip(c.Request().Context(), shipSymbol).Execute()
	if err != nil {
		c.Logger().Error(err.Error())
		return shared.Error(err).Render(c.Request().Context(), c.Response())
	}

	return ship.Page(*resp).Render(c.Request().Context(), c.Response())
}
