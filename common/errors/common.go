package errors

import "net/http"

// Common errors used across multiple controllers
var (
	InvalidRequestError = ControllerError{
		Name:          "InvalidRequestError",
		HttpErrorCode: http.StatusBadRequest,
		Translation: Translation{
			Fr: "Impossible de parser le corps de la requête.",
			En: "Failed to parse request body.",
		},
	}

	ServiceError = ControllerError{
		Name:          "ServiceError",
		HttpErrorCode: http.StatusInternalServerError,
		Translation: Translation{
			Fr: "Erreur lors du traitement de la requête.",
			En: "Failed to process request.",
		},
	}
)
