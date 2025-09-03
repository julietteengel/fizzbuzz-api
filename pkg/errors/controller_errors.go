package errors

import "net/http"

type Translation struct {
	Fr string `json:"fr"`
	En string `json:"en"`
}

type ControllerError struct {
	Name          string      `json:"name"`
	HttpErrorCode int         `json:"code"`
	Translation   Translation `json:"translation"`
}

func (e ControllerError) Error() string {
	return e.Translation.En
}

var (
	InvalidRequestError = ControllerError{
		Name:          "InvalidRequestError",
		HttpErrorCode: http.StatusBadRequest,
		Translation: Translation{
			Fr: "Impossible de parser le corps de la requête.",
			En: "Failed to parse request body.",
		},
	}

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

	ServiceError = ControllerError{
		Name:          "ServiceError",
		HttpErrorCode: http.StatusInternalServerError,
		Translation: Translation{
			Fr: "Erreur lors de la génération de la séquence FizzBuzz.",
			En: "Failed to generate FizzBuzz sequence.",
		},
	}

	GenericParamsError = ControllerError{
		Name:          "GenericParamsError",
		HttpErrorCode: http.StatusBadRequest,
		Translation: Translation{
			Fr: "Impossible de récupérer les paramètres.",
			En: "Cannot get params.",
		},
	}
)