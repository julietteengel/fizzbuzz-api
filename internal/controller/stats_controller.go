package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/julietteengel/fizzbuzz-api/common/errors"
	"github.com/julietteengel/fizzbuzz-api/internal/service"
)

type StatsController struct {
	service service.IStatsService
}

func NewStatsController(service service.IStatsService) *StatsController {
	return &StatsController{
		service: service,
	}
}

// GetStats returns statistics about the most frequently requested parameters.
// @Summary Get FizzBuzz statistics
// @Description Returns statistics about the most frequently requested FizzBuzz parameters
// @Tags stats
// @Produce json
// @Success 200 {object} model.StatsResponse
// @Success 204 "No statistics available yet"
// @Failure 500 {string} string "Service error message (translated)"
// @Router /api/v1/stats [get]
func (c *StatsController) GetStats(ctx echo.Context) error {
	stats, err := c.service.GetMostFrequent(ctx.Request().Context())
	if err != nil {
		return errors.WrapErrorHTTP(ctx, err, errors.ServiceError)
	}

	if stats == nil {
		return ctx.NoContent(http.StatusNoContent)
	}

	return ctx.JSON(http.StatusOK, stats)
}