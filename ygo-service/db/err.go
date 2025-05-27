package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ygo-skc/skc-go/common/model"
)

const (
	genericError = "Error occurred while querying DB"
)

func handleQueryError(logger *slog.Logger, err error) *model.APIError {
	logger.Error(fmt.Sprintf("Error fetching data from DB - %v", err))

	if err == sql.ErrNoRows {
		return &model.APIError{
			Message:    "No results found",
			StatusCode: http.StatusNotFound,
		}
	}
	return &model.APIError{Message: genericError, StatusCode: http.StatusInternalServerError}
}

func handleRowParsingError(logger *slog.Logger, err error) *model.APIError {
	logger.Error(fmt.Sprintf("Error parsing data from DB - %v", err))
	return &model.APIError{Message: genericError, StatusCode: http.StatusInternalServerError}
}
