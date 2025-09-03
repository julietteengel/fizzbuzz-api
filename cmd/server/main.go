package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"

	"github.com/julietteengel/fizzbuzz-api/internal/config"
	"github.com/julietteengel/fizzbuzz-api/internal/controller"
	"github.com/julietteengel/fizzbuzz-api/internal/service"
)

func main() {
	fx.New(
		fx.Provide(
			config.Load,
			service.NewFizzBuzzService,
			controller.NewFizzBuzzController,
			newEcho,
		),
		fx.Invoke(setupRoutes),
	).Run()
}

func newEcho() *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())

	return e
}

func setupRoutes(
	lc fx.Lifecycle,
	e *echo.Echo,
	cfg *config.Config,
	fizzBuzzController *controller.FizzBuzzController,
) {
	// API Routes
	api := e.Group("/api/v1")
	api.POST("/fizzbuzz", fizzBuzzController.GenerateFizzBuzz)
	
	// Health check (outside API group)
	e.GET("/health", fizzBuzzController.HealthCheck)

	// Server lifecycle
	port := ":" + cfg.Server.Port

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := e.Start(port); err != nil && err != http.ErrServerClosed {
					fmt.Printf("Failed to start server: %v\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Shutting down server...")
			return e.Shutdown(ctx)
		},
	})
}