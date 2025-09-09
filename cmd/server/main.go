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
	"github.com/julietteengel/fizzbuzz-api/internal/database"
	"github.com/julietteengel/fizzbuzz-api/internal/repository"
	"github.com/julietteengel/fizzbuzz-api/internal/service"
)

// @title FizzBuzz REST API
// @version 1.0
// @description A production-ready FizzBuzz REST API server built with Go and Echo framework following clean architecture principles.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email juliette.engel@skema.edu
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

func main() {
	fx.New( //creates and initializes an App
		fx.Provide( // registers any number of constructor functions, teaching the application how to instantiate various types.
			config.Load,
			database.NewGormDB,
			repository.NewStatsRepository,
			service.NewFizzBuzzService,
			service.NewStatsService,
			controller.NewFizzBuzzController,
			controller.NewStatsController,
			newEcho,
		),
		fx.Invoke(setupRoutes), //Invoke registers functions that are executed eagerly on application start.
	).Run()
}

func newEcho() *echo.Echo {
	e := echo.New()

	//1. Middleware Stack:
	e.Use(middleware.Logger())    // Logs every request
	e.Use(middleware.Recover())   // Prevents crashes
	e.Use(middleware.CORS())      // Browser compatibility
	e.Use(middleware.RequestID()) // Request tracing

	return e
}

func setupRoutes(
	lc fx.Lifecycle,
	e *echo.Echo,
	cfg *config.Config,
	fizzBuzzController *controller.FizzBuzzController,
	statsController *controller.StatsController,
) {
	// API Routes
	//2. Route Grouping:
	api := e.Group("/api/v1") // Prefix all API routes
	api.POST("/fizzbuzz", fizzBuzzController.GenerateFizzBuzz)
	api.GET("/stats", statsController.GetStats)

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
			return e.Shutdown(ctx) // Properly close connections
		},
	})
}
