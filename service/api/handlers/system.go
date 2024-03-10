package handlers

import (
	"github.com/labstack/echo/v4"

	"space-traders/service/views/components/system"
)

func (vh *ViewHandler) MountSystemRoutes(e *echo.Echo) {
	e.GET("/system", vh.SystemPage)
	e.GET("/system/info", vh.SystemInfo)
	e.GET("/system/waypoints", vh.SystemWaypoints)
	e.GET("/system/locations", vh.SystemLocations)
}

func (vh *ViewHandler) SystemPage(c echo.Context) error {
	return system.Page().Render(c.Request().Context(), c.Response())
}

func (vh *ViewHandler) SystemInfo(c echo.Context) error {
	agentResp, _, err := vh.Client.AgentsAPI.GetMyAgent(c.Request().Context()).Execute()
	if err != nil {
		return err
	}

	resp, _, err := vh.Client.SystemsAPI.GetSystem(c.Request().Context(), agentResp.Data.Headquarters[:len(agentResp.Data.Headquarters)-3]).Execute()
	if err != nil {
		return err
	}

	return system.SystemInfo(resp.Data).Render(c.Request().Context(), c.Response())
}

func (vh *ViewHandler) SystemWaypoints(c echo.Context) error {
	agentResp, _, err := vh.Client.AgentsAPI.GetMyAgent(c.Request().Context()).Execute()
	if err != nil {
		return err
	}

	resp, _, err := vh.Client.SystemsAPI.GetSystem(c.Request().Context(), agentResp.Data.Headquarters[:len(agentResp.Data.Headquarters)-3]).Execute()
	if err != nil {
		return err
	}

	return system.WaypointList(resp.GetData().Waypoints).Render(c.Request().Context(), c.Response())
}

//TODO: create a db to store the system data, periodically refresh in a go routine?

func (vh *ViewHandler) SystemLocations(c echo.Context) error {
	// agentResp, _, err := vh.Client.AgentsAPI.GetMyAgent(c.Request().Context()).Execute()
	// if err != nil {
	// 	return err
	// }

	// resp, _, err := vh.Client.SystemsAPI.GetSystem(c.Request().Context(), agentResp.Data.Headquarters[:len(agentResp.Data.Headquarters)-3]).Execute()
	// if err != nil {
	// 	return err
	// }

	return system.SystemLocations().Render(c.Request().Context(), c.Response())
}
