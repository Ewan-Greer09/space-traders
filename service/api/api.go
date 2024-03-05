package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"space-traders/service/api/handlers"
)

type API struct {
	e           *echo.Echo
	ViewHandler *handlers.ViewHandler
}

func NewAPI() *API {
	e := echo.New()
	e.HideBanner = true

	a := &API{e: e, ViewHandler: handlers.NewViewHandler()}
	a.Routes()

	return a
}

func (a *API) Start() error {
	if err := a.e.Start("127.0.0.1:3000"); err != nil {
		return err
	}

	return nil
}

func (a *API) Routes() {
	a.e.Use(middleware.Logger())
	a.e.Use(middleware.Recover())
	a.e.Use(middleware.CORS())
	a.e.Use(middleware.Gzip())
	a.e.Use(middleware.Secure())
	a.e.Use(getSessionCookie())
	a.e.Use(a.ViewHandler.AddKeyToReq())

	a.e.Static("/static", "service/views/pages/css")
	com := a.e.Group("/com")

	// group for displaying pages starting with a blank path
	page := a.e.Group("")

	/*
		Page Includes
	*/

	com.GET("/header", a.ViewHandler.GetHeader)
	com.GET("/footer", a.ViewHandler.GetFooter)

	// not a page, but is at the root of the site path
	page.GET("/logout", a.ViewHandler.Logout)

	/*
		Index
	*/

	index := com.Group("/index")
	_ = index

	page.GET("/", a.ViewHandler.GetIndex) // returns the index page

	/*
		Agents
	*/

	agents := com.Group("/agents")
	_ = agents

	// agents.GET("/", a.ViewHandler.GetAgents) // returns the agents pages main content

	/*
		Contracts
	*/

	contracts := com.Group("/contracts")
	_ = contracts

	// contracts.GET("/", a.ViewHandler.GetContracts) // returns the contracts pages main content

	/*
		Factions
	*/

	factions := com.Group("/factions")
	_ = factions

	// factions.GET("/", a.ViewHandler.GetFactions) // returns the factions pages main content

	/*
		Fleet
	*/

	fleet := com.Group("/fleet")

	page.GET("/fleet", a.ViewHandler.GetFleet)             // returns the fleet page
	page.GET("/fleet/:symbol", a.ViewHandler.GetFleetShip) // returns the fleet ship page (with symbol as a parameter)

	fleet.GET("/list", a.ViewHandler.GetFleetList)

	/*
		Systems
	*/

	systems := com.Group("/systems")
	_ = systems

	// systems.GET("/", a.ViewHandler.GetSystems) // returns the systems pages main content

	/*
		Login
	*/

	login := com.Group("/login")

	page.GET("/login", a.ViewHandler.GetLogin) // returns the login page
	login.POST("/login", a.ViewHandler.HandleLogin)
}

func getSessionCookie() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("session")
			if err != nil {
				return next(c)
			}

			c.Set("session", cookie.Value)

			return next(c)
		}
	}
}
