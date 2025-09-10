package errors

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// Translation represents error messages in multiple languages
type Translation struct {
	Fr string `json:"fr"`
	En string `json:"en"`
}

// ControllerError represents a standardized error structure
type ControllerError struct {
	Name          string      `json:"name"`
	HttpErrorCode int         `json:"code"`
	Translation   Translation `json:"translation"`
}

func (e ControllerError) Error() string {
	return e.Translation.En
}

// WrapErrorHTTP logs the original error and returns a properly formatted HTTP error
func WrapErrorHTTP(c echo.Context, originalErr error, controllerError ControllerError) error {
	if originalErr != nil {
		log.Errorf("Error %s: %v", controllerError.Name, originalErr)
	}

	lang := c.Request().Header.Get("Accept-Language")
	message := controllerError.Translation.En
	if lang == "fr" || lang == "fr-FR" {
		message = controllerError.Translation.Fr
	}

	return echo.NewHTTPError(controllerError.HttpErrorCode, message)
}