package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/julietteengel/fizzbuzz-api/internal/config"
	"github.com/julietteengel/fizzbuzz-api/internal/controller"
	"github.com/julietteengel/fizzbuzz-api/internal/service"
)

func main() {
	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize services
	fizzBuzzService := service.NewFizzBuzzService()

	// Initialize controllers
	fizzBuzzController := controller.NewFizzBuzzController(fizzBuzzService)

	// Routes
	e.POST("/fizzbuzz", fizzBuzzController.GenerateFizzBuzz)
	e.GET("/health", fizzBuzzController.HealthCheck)

	// Start server
	port := ":8080"
	log.Printf("Starting server on port %s", port)
	if err := e.Start(port); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}