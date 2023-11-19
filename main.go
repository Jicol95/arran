package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jicol-95/arran/dao"
	"github.com/jicol-95/arran/handler"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	if err := dao.RunDatabaseMigrations(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	e := echo.New()
	addMiddleware(e)

	e.GET("/metrics", echoprometheus.NewHandler())
	e.GET("rest/health", handler.HealthHandler)
	e.GET("rest/v1/error", handler.ErrorExampleHandler)

	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	waitForInterruptSignal(e)
}

func addMiddleware(e *echo.Echo) {
	e.Logger.SetLevel(log.INFO)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(echoprometheus.NewMiddleware("arran")) // adds middleware to gather metrics
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.PATCH},
	}))
}

func waitForInterruptSignal(e *echo.Echo) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	e.Logger.Info("Interrupt signal received, gracefully shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
