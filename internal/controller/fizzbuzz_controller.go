package controller

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/julietteengel/fizzbuzz-api/internal/model"
	"github.com/julietteengel/fizzbuzz-api/internal/service"
	"github.com/julietteengel/fizzbuzz-api/pkg/errors"
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
	lang := ctx.Request().Header.Get("Accept-Language")

	if err := ctx.Bind(&request); err != nil {
		return c.handleError(ctx, errors.InvalidRequestError, lang)
	}

	// Validation
	if request.Int1 <= 0 {
		return c.handleError(ctx, errors.ValidationInt1Error, lang)
	}

	if request.Int2 <= 0 {
		return c.handleError(ctx, errors.ValidationInt2Error, lang)
	}

	if request.Limit <= 0 || request.Limit > 10000 {
		return c.handleError(ctx, errors.ValidationLimitError, lang)
	}

	if len(request.Str1) == 0 || len(request.Str1) > 100 {
		return c.handleError(ctx, errors.ValidationStr1Error, lang)
	}

	if len(request.Str2) == 0 || len(request.Str2) > 100 {
		return c.handleError(ctx, errors.ValidationStr2Error, lang)
	}

	response, err := c.service.GenerateFizzBuzz(request)
	if err != nil {
		return c.handleError(ctx, errors.ServiceError, lang)
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *FizzBuzzController) handleError(ctx echo.Context, err errors.ControllerError, lang string) error {
	message := err.Translation.En
	if lang == "fr" || lang == "fr-FR" {
		message = err.Translation.Fr
	}

	return ctx.JSON(err.HttpErrorCode, model.ErrorResponse{
		Error:   err.Name,
		Message: message,
	})
}

func (c *FizzBuzzController) HealthCheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, model.HealthCheckResponse{
		Status:    "ok",
		Timestamp: time.Now(),
	})
}