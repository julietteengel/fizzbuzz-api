package errors

import "net/http"

// FizzBuzz specific errors
var (
	ValidationInt1Error = ControllerError{
		Name:          "ValidationInt1Error",
		HttpErrorCode: http.StatusBadRequest,
		Translation: Translation{
			Fr: "Le paramètre int1 doit être supérieur à 0.",
			En: "Parameter int1 must be greater than 0.",
		},
	}

	ValidationInt2Error = ControllerError{
		Name:          "ValidationInt2Error",
		HttpErrorCode: http.StatusBadRequest,
		Translation: Translation{
			Fr: "Le paramètre int2 doit être supérieur à 0.",
			En: "Parameter int2 must be greater than 0.",
		},
	}

	ValidationLimitError = ControllerError{
		Name:          "ValidationLimitError",
		HttpErrorCode: http.StatusBadRequest,
		Translation: Translation{
			Fr: "Le paramètre limit doit être entre 1 et 10000.",
			En: "Parameter limit must be between 1 and 10000.",
		},
	}

	ValidationStr1Error = ControllerError{
		Name:          "ValidationStr1Error",
		HttpErrorCode: http.StatusBadRequest,
		Translation: Translation{
			Fr: "Le paramètre str1 doit contenir entre 1 et 100 caractères.",
			En: "Parameter str1 must be between 1 and 100 characters.",
		},
	}

	ValidationStr2Error = ControllerError{
		Name:          "ValidationStr2Error",
		HttpErrorCode: http.StatusBadRequest,
		Translation: Translation{
			Fr: "Le paramètre str2 doit contenir entre 1 et 100 caractères.",
			En: "Parameter str2 must be between 1 and 100 characters.",
		},
	}

	FizzBuzzGenerationError = ControllerError{
		Name:          "FizzBuzzGenerationError",
		HttpErrorCode: http.StatusInternalServerError,
		Translation: Translation{
			Fr: "Erreur lors de la génération de la séquence FizzBuzz.",
			En: "Failed to generate FizzBuzz sequence.",
		},
	}
)