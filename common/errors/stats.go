package errors

import "net/http"

// Stats specific errors
var (
	StatsNotFoundError = ControllerError{
		Name:          "StatsNotFoundError",
		HttpErrorCode: http.StatusNoContent,
		Translation: Translation{
			Fr: "Aucune statistique disponible.",
			En: "No statistics available.",
		},
	}

	StatsRetrievalError = ControllerError{
		Name:          "StatsRetrievalError",
		HttpErrorCode: http.StatusInternalServerError,
		Translation: Translation{
			Fr: "Erreur lors de la récupération des statistiques.",
			En: "Failed to retrieve statistics.",
		},
	}
)
