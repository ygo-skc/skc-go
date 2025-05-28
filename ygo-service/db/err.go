package db

import (
	"database/sql"
	"fmt"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	genericError = "Error occurred while querying DB"
)

func handleQueryError(logger *slog.Logger, err error) *status.Status {
	logger.Error(fmt.Sprintf("Error fetching data from DB - %v", err))

	if err == sql.ErrNoRows {
		return status.New(codes.NotFound, "No results found")
	}
	return status.New(codes.Internal, genericError)
}

func handleRowParsingError(logger *slog.Logger, err error) *status.Status {
	logger.Error(fmt.Sprintf("Error parsing data from DB - %v", err))
	return status.New(codes.Internal, genericError)
}
