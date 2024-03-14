package handlers

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"

	"space-traders/repository/postgres"
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

	err = vh.addSystems(c)
	if err != nil {
		c.Logger().Error(err.Error())
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
		_, err = vh.userDB.CreateSystem(context.Background(), postgres.CreateSystemParams{
			Symbol:       pgtype.Text{String: resp.Data.Symbol, Valid: true},
			SectorSymbol: pgtype.Text{String: resp.Data.SectorSymbol, Valid: true},
			Type:         pgtype.Text{String: string(resp.Data.Type), Valid: true},
			X:            pgtype.Int4{Int32: int32(resp.Data.X), Valid: true},
			Y:            pgtype.Int4{Int32: int32(resp.Data.Y), Valid: true},
		})
		if err != nil {
			return err
		}
	}

	return nil
}
