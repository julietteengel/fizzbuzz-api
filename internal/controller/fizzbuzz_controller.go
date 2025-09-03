package controller

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/julietteengel/fizzbuzz-api/internal/model"
	"github.com/julietteengel/fizzbuzz-api/internal/service"
)

type FizzBuzzController struct {
	service service.FizzBuzzService
}

func NewFizzBuzzController(service service.FizzBuzzService) *FizzBuzzController {
	return &FizzBuzzController{
		service: service,
	}
}

func (c *FizzBuzzController) GenerateFizzBuzz(ctx echo.Context) error {
	var request model.FizzBuzzRequest

	if err := ctx.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Failed to parse request body",
		})
	}

	response, err := c.service.GenerateFizzBuzz(request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrorResponse{
			Error:   "service_error",
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *FizzBuzzController) HealthCheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, model.HealthCheckResponse{
		Status:    "ok",
		Timestamp: time.Now(),
	})
}