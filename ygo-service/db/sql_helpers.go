package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/ygo-skc/skc-go/common/v2/ygo"
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

func queryProductInfo(logger *slog.Logger, productID string) (*ygo.Product, *status.Status) {
	var id, locale, name, t, subType, releaseDate string

	if err := skcDBConn.QueryRow(productDetailsQuery, productID).Scan(&id, &locale, &name, &t, &subType, &releaseDate); err != nil {
		return nil, handleQueryError(logger, err)
	}
	return &ygo.Product{ID: id, Locale: locale, Name: name, ReleaseDate: releaseDate, Type: t, SubType: subType}, nil
}

func convertToFullText(subject string) string {
	fullTextSubject := spaceRegex.ReplaceAllString(strings.ReplaceAll(subject, "-", " "), " +")
	return fmt.Sprintf(`"+%s"`, fullTextSubject) // match phrase, not all words in text will match only consecutive matches of words in phrase
}

func buildVariableQuerySubjects(subjects []string) ([]interface{}, int) {
	numSubjects := len(subjects)
	args := make([]interface{}, numSubjects)

	for index, subject := range subjects {
		args[index] = subject
	}

	return args, numSubjects
}

func variablePlaceholders(totalFields int) string {
	switch totalFields {
	case 0:
		return ""
	case 1:
		return "?"
	default:
		return fmt.Sprintf("?%s", strings.Repeat(", ?", totalFields-1))
	}
}
