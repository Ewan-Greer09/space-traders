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
		ViewHandler: handlers.NewViewHandler(),
		Cfg:         config.MustLoadConfig(),
	}
	a.Routes()

	return a
}

func (a *API) Start() error {
	if err := a.e.Start(fmt.Sprintf("%s:%s", a.Cfg.Host, a.Cfg.Host)); err != nil {
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
	a.e.Use(getSessionCookie())
	a.e.Use(a.ViewHandler.AddKeyToReq())

	a.e.Static("/static", "service/views/pages/css")

	a.ViewHandler.MountIndexRoutes(a.e)
	a.ViewHandler.MountFleetRoutes(a.e)
	a.ViewHandler.MountLoginRoutes(a.e)
	a.ViewHandler.MountSharedRoutes(a.e)
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
