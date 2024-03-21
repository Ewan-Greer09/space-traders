package handlers

import (
	"context"
	"database/sql"

	"github.com/labstack/echo/v4"

	"space-traders/repository/mysql"
	"space-traders/service/views/components/fleet"
	"space-traders/service/views/components/shared"
	"space-traders/service/views/components/ship"
)

func (vh *ViewHandler) MountFleetRoutes(e *echo.Echo) {
	e.GET("/fleet", vh.GetFleet)
	e.GET("/fleet/list", vh.GetFleetList)
	e.GET("/fleet/ship/:symbol", vh.GetFleetShip)
	e.GET("/fleet/ship/:symbol/nav", vh.GetShipsNav)
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

	err = vh.addSystems(c)
	if err != nil {
		c.Logger().Info(err.Error())
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

func (vh *ViewHandler) addSystems(c echo.Context) error {
	resp, _, err := vh.Client.FleetAPI.GetMyShips(c.Request().Context()).Execute()
	if err != nil {
		return err
	}

	for _, ship := range resp.Data {
		// get the system details
		resp, _, err := vh.Client.SystemsAPI.GetSystem(context.Background(), ship.Nav.SystemSymbol).Execute()
		if err != nil {
			return err
		}

		// store the system details in the database
		err = vh.userDB.CreateSystem(c.Request().Context(), mysql.CreateSystemParams{
			Symbol:       sql.NullString{String: resp.Data.Symbol, Valid: true},
			SectorSymbol: sql.NullString{String: resp.Data.SectorSymbol, Valid: true},
			Type:         sql.NullString{String: string(resp.Data.GetType()), Valid: true},
			X:            sql.NullInt32{Int32: int32(resp.Data.X), Valid: true},
			Y:            sql.NullInt32{Int32: int32(resp.Data.Y), Valid: true},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (vh *ViewHandler) GetShipsNav(c echo.Context) error {
	return nil
}
