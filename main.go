package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jicol-95/arran/config"
	"github.com/jicol-95/arran/consumer"
	"github.com/jicol-95/arran/dal"
	"github.com/jicol-95/arran/domain"
	"github.com/jicol-95/arran/handler"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	addMiddleware(e)
	logger := e.Logger

	cfg := config.NewArranConfig()

	db, err := dal.InitDB(cfg.PostgresConfig)

	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}

	err = dal.RunDatabaseMigrations(cfg.PostgresConfig)

	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}

	tm := dal.NewTransactionManager(db)
	exampleResourceRepo := dal.NewExampleResourceRepository(db)
	exampleResourceService := domain.NewExampleResourceService(logger, tm, exampleResourceRepo)

	_, err = consumer.ProcessExampleResourceTopic(cfg.Kafka, exampleResourceService, logger)

	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}

	e.GET("/metrics", echoprometheus.NewHandler())
	e.GET("/rest/health", handler.HealthHandler)

	e.POST("/rest/v1/example-resources", handler.ExampleResourcePostHandler(tm, exampleResourceService))
	e.PUT("/rest/v1/example-resources/:id", handler.ExampleResourceUpdateByIdHandler(tm, exampleResourceService))
	e.GET("/rest/v1/example-resources/:id", handler.ExampleResourceGetByIdHandler(tm, exampleResourceService))
	e.DELETE("/rest/v1/example-resources/:id", handler.ExampleResourceDeleteByIdHandler(tm, exampleResourceService))

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
	e.Use(echoprometheus.NewMiddleware("arran"))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
}

func waitForInterruptSignal(e *echo.Echo) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
	e.Logger.Info("Interrupt signal received, gracefully shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
