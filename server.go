package main

import (
	"go-echo/config"
	"go-echo/routes"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/unrolled/secure"
)

func main() {
	config.ConnectDB()

	secureMiddleware := secure.New(secure.Options{
		AllowedHosts:            []string{"localhost:9000", "www.google.com"},
		FrameDeny:               true,
		CustomFrameOptionsValue: "SAMEORIGIN",
		ContentTypeNosniff:      true,
		BrowserXssFilter:        true,
	})

	r := echo.New()

	r.Use(config.MiddlewareLogging)
	r.Use(middleware.CORS())
	r.Use(echo.WrapMiddleware(secureMiddleware.Handler))

	r.HTTPErrorHandler = config.ErrorHandler

	routes.ApiRoute(r)
	routes.PublicRoute(r)

	lock := make(chan error)
	go func(lock chan error) { lock <- r.Start(":9000") }(lock)

	time.Sleep(1 * time.Millisecond)
	config.MakeLogEntry(nil).Warning("application started without ssl/tls enabled")

	err := <-lock
	if err != nil {
		config.MakeLogEntry(nil).Panic("failed to start application")
	}
}
