package api

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"space-traders/service/api/handlers"
	"space-traders/service/config"
	"space-traders/service/views/index"
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
	a.reigsterRoutes()

	return a
}

func (a *API) Start() error {
	if err := a.e.Start(fmt.Sprintf("%s:%s", a.Cfg.Host, a.Cfg.Port)); err != nil {
		return err
	}
	return nil
}

func (api *API) reigsterRoutes() {
	api.e.Use(
		middleware.RequestID(),
		middleware.Logger(),
		middleware.Recover(),
		middleware.CORS(),
		middleware.Gzip(),
		middleware.Secure(),
		api.ViewHandler.AddKeyToReq(),
		api.ViewHandler.LoginRedirect(),
	)

	api.e.Static("/static", "service/static")

	MountRoutes(api.e, *api.ViewHandler)

	api.e.GET("/", api.ViewHandler.HandlerHeader)

	api.e.GET("/ship/data", func(c echo.Context) error {
		return index.ShipData("This is some demo ship data").Render(c.Request().Context(), c.Response())
	})

	api.e.GET("/system/data", func(c echo.Context) error {
		return index.SystemData("This is some demo system data").Render(c.Request().Context(), c.Response())
	})
}
