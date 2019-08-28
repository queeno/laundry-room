package laundry

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type laundryApp struct {}

func NewLaundryApp() laundryApp {
	return laundryApp{}
}

func (c laundryApp) Run(ctx context.Context) {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.File("/","laundry.html")
	e.Static("/static", "static")

	go func() {
		if err := e.Start(":8080"); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	select {
	case <-ctx.Done():
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}
	}
}
