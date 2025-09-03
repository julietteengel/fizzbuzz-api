package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/julietteengel/fizzbuzz-api/internal/config"
	"github.com/julietteengel/fizzbuzz-api/internal/controller"
	"github.com/julietteengel/fizzbuzz-api/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Echo
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())

	// Initialize services
	fizzBuzzService := service.NewFizzBuzzService()

	// Initialize controllers
	fizzBuzzController := controller.NewFizzBuzzController(fizzBuzzService)

	// API Routes
	api := e.Group("/api/v1")
	api.POST("/fizzbuzz", fizzBuzzController.GenerateFizzBuzz)
	
	// Health check (outside API group)
	e.GET("/health", fizzBuzzController.HealthCheck)

	// Start server
	port := ":" + cfg.Server.Port
	log.Printf("Starting %s server on port %s", cfg.App.Environment, port)
	if err := e.Start(port); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}