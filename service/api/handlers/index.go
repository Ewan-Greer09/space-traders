package handlers

import (
	"github.com/labstack/echo/v4"

	"space-traders/service/views/components/index"
	"space-traders/service/views/components/widgets"
)

func (vh *ViewHandler) MountIndexRoutes(e *echo.Echo) {
	e.GET("/", vh.GetIndex)
	e.GET("/move-ship", vh.MoveShip)
	e.GET("/submit-form", vh.SubmitForm)
	e.GET("set-destination", vh.HandleNavigation)
}

func (vh *ViewHandler) GetIndex(c echo.Context) error {
	return index.Page().Render(c.Request().Context(), c.Response())
}

func (vh *ViewHandler) MoveShip(c echo.Context) error {

	resp, _, err := vh.Client.FleetAPI.GetMyShips(c.Request().Context()).Execute()
	if err != nil {
		return err
	}

	return widgets.NavWidget(resp.Data).Render(c.Request().Context(), c.Response())
}

func (vh *ViewHandler) SubmitForm(c echo.Context) error {
	shipSymbol := c.FormValue("ship-select")
	resp, _, err := vh.Client.FleetAPI.GetMyShip(c.Request().Context(), shipSymbol).Execute()
	if err != nil {
		return err
	}

	sysResp, _, err := vh.Client.SystemsAPI.GetSystemWaypoints(c.Request().Context(), resp.Data.Nav.SystemSymbol).Execute()
	if err != nil {
		return err
	}

	return widgets.SelectedShip(resp.Data, sysResp.Data).Render(c.Request().Context(), c.Response())
}

func (vh *ViewHandler) HandleNavigation(c echo.Context) error {
	waypointSymbol := c.FormValue("waypoint-symbol")
	shipSymbol := c.FormValue("ship-symbol")

	resp, err := vh.myClient.GetMyShips(shipSymbol)
	if err != nil {
		return err
	}

	if resp.Data.Nav.Status == "DOCKED" {
		err := vh.myClient.SendToOrbit(shipSymbol)
		if err != nil {
			return err
		}
	}

	navResp, err := vh.myClient.NavigateShip(shipSymbol, waypointSymbol)
	if err != nil {
		return err
	}

	return widgets.SelectedDest(*navResp).Render(c.Request().Context(), c.Response())
}
