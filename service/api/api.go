package api

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"space-traders/service/api/handlers"
	"space-traders/service/config"
)

type API struct {
	e           *echo.Echo
	ViewHandler *handlers.ViewHandler
	Cfg         *config.Config
}

func NewAPI() *API {
	e := echo.New()
	e.HideBanner = true

	a := &API{
		e:           e,
		ViewHandler: handlers.NewViewHandler(config.MustLoadConfig()),
		Cfg:         config.MustLoadConfig(),
	}
	a.Routes()

	return a
}

func (a *API) Start() error {
	if err := a.e.Start(fmt.Sprintf("%s:%s", a.Cfg.Host, a.Cfg.Port)); err != nil {
		return err
	}
	return nil
}

func (a *API) Routes() {
	a.e.Use(middleware.RequestID())
	a.e.Use(middleware.Logger())
	a.e.Use(middleware.Recover())
	a.e.Use(middleware.CORS())
	a.e.Use(middleware.Gzip())
	a.e.Use(middleware.Secure())
	a.e.Use(a.ViewHandler.AddKeyToReq())

	a.e.Static("/static", "service/views")
	a.e.GET("/favicon.ico", a.ViewHandler.Favicon)

	a.ViewHandler.MountIndexRoutes(a.e)
	a.ViewHandler.MountFleetRoutes(a.e)
	a.ViewHandler.MountLoginRoutes(a.e)
	a.ViewHandler.MountSharedRoutes(a.e)
	a.ViewHandler.MountSystemRoutes(a.e)
}
