package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/a-h/templ"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/html"

	"space-traders/service/views/components/system"
)

func (vh *ViewHandler) MountSystemRoutes(e *echo.Echo) {
	e.GET("/system", vh.SystemPage)
	e.GET("/system/info", vh.SystemInfo)
	e.GET("/system/waypoints", vh.SystemWaypoints)
	e.GET("/system/locations", vh.SystemLocations)
	e.GET("/system/waypoint/:id", vh.SystemWaypointInfo)
	e.GET("/system/waypoint/:id/info", vh.SystemWaypointInfoContent)
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

type Coordinate struct {
	X int32
	Y int32
}

// SystemWaypoint represents a waypoint with a name and coordinates
type SystemWaypoint struct {
	Name   string     // Name of the waypoint
	Coords Coordinate // Coordinates of the waypoint
}

func (vh *ViewHandler) SystemLocations(c echo.Context) error {
	agentResp, _, err := vh.Client.AgentsAPI.GetMyAgent(c.Request().Context()).Execute()
	if err != nil {
		return err
	}

	resp, _, err := vh.Client.SystemsAPI.GetSystem(c.Request().Context(), agentResp.Data.Headquarters[:len(agentResp.Data.Headquarters)-3]).Execute()
	if err != nil {
		return err
	}

	var points []opts.ScatterData
	for _, item := range resp.Data.Waypoints {
		points = append(points, opts.ScatterData{Value: []interface{}{item.X, item.Y}, Name: fmt.Sprint(item.Symbol + " - " + string(item.Type))})
	}

	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Solar System Representation",
			Subtitle: "Scatter plot of items in the solar system",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "X",
			SplitLine: &opts.SplitLine{
				Show: false,
			},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Y",
			SplitLine: &opts.SplitLine{
				Show: false,
			},
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      true,
			Enterable: true,
			Formatter: "{b}: {c}",
		}),
	)

	scatter.AddSeries("Waypoint", points, charts.WithItemStyleOpts(opts.ItemStyle{
		Color: "#e74c3c",
	})).SetSeriesOptions(charts.WithLabelOpts(opts.Label{
		Show: false,
	}))

	buff := bytes.Buffer{}
	scatter.Render(&buff)

	str, err := removeUnwantedElements(buff.String())
	if err != nil {
		return err
	}

	return system.SystemLocations(Unsafe(str)).Render(c.Request().Context(), c.Response())
}

func Unsafe(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}

// removeUnwantedElements takes an HTML string and returns a new HTML string
// with only the script tag in the header and the content in the body.
func removeUnwantedElements(htmlStr string) (string, error) {
	doc, err := html.Parse(bytes.NewBufferString(htmlStr))
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "head" || n.Data == "body") {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if n.Data == "head" && c.Type == html.ElementNode && c.Data == "script" {
					html.Render(&b, c) // Render script tag in head
				} else if n.Data == "body" {
					html.Render(&b, c) // Render everything in body
				}
			}
		} else {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}

	f(doc)

	return b.String(), nil
}

func (vh *ViewHandler) SystemWaypointInfo(c echo.Context) error {
	waypointID := c.Param("id")
	return system.WaypointPage(waypointID).Render(c.Request().Context(), c.Response())
}

func (vh *ViewHandler) SystemWaypointInfoContent(c echo.Context) error {
	waypointID := c.Param("id")
	waypointIDParts := strings.Split(waypointID, "-")

	for i, part := range waypointIDParts {
		waypointIDParts[i] = strings.ToUpper(part)
	}

	resp, _, err := vh.Client.SystemsAPI.GetWaypoint(c.Request().Context(), fmt.Sprintf("%s-%s", waypointIDParts[0], waypointIDParts[1]), strings.ToUpper(waypointID)).Execute()
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	return system.WaypointInfo(resp.Data).Render(c.Request().Context(), c.Response())
}
