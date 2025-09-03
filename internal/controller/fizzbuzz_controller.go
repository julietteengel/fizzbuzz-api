package controller

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/julietteengel/fizzbuzz-api/internal/model"
	"github.com/julietteengel/fizzbuzz-api/internal/service"
	"github.com/julietteengel/fizzbuzz-api/common/errors"
)

type FizzBuzzController struct {
	service service.IFizzBuzzService
}

func NewFizzBuzzController(service service.IFizzBuzzService) *FizzBuzzController {
	return &FizzBuzzController{
		service: service,
	}
}

func (c *FizzBuzzController) GenerateFizzBuzz(ctx echo.Context) error {
	var request model.FizzBuzzRequest

	if err := ctx.Bind(&request); err != nil {
		return errors.WrapErrorHTTP(ctx, err, errors.InvalidRequestError)
	}

	// Validation
	if request.Int1 <= 0 {
		return errors.WrapErrorHTTP(ctx, nil, errors.ValidationInt1Error)
	}

	if request.Int2 <= 0 {
		return errors.WrapErrorHTTP(ctx, nil, errors.ValidationInt2Error)
	}

	if request.Limit <= 0 || request.Limit > 10000 {
		return errors.WrapErrorHTTP(ctx, nil, errors.ValidationLimitError)
	}

	if len(request.Str1) == 0 || len(request.Str1) > 100 {
		return errors.WrapErrorHTTP(ctx, nil, errors.ValidationStr1Error)
	}

	if len(request.Str2) == 0 || len(request.Str2) > 100 {
		return errors.WrapErrorHTTP(ctx, nil, errors.ValidationStr2Error)
	}

	response, err := c.service.GenerateFizzBuzz(ctx.Request().Context(), request)
	if err != nil {
		return errors.WrapErrorHTTP(ctx, err, errors.ServiceError)
	}

	return ctx.JSON(http.StatusOK, response)
}


func (c *FizzBuzzController) HealthCheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, model.HealthCheckResponse{
		Status:    "ok",
		Timestamp: time.Now(),
	})
}