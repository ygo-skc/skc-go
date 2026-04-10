package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/ygo-skc/skc-go/common/v2/ygo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	genericError = "Error occurred while querying DB"
)

var (
	spaceRegex = regexp.MustCompile(`[ ]+`)
	quoteRegex = regexp.MustCompile(`['"]`)
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

// removes quotes from text as full text search does not handle them well
func convertToFullText(subject string) string {
	return fmt.Sprintf(`"%s"`, quoteRegex.ReplaceAllString(subject, "")) // match phrase
}

func buildVariableQuerySubjects(subjects []string) ([]any, int) {
	numSubjects := len(subjects)
	args := make([]any, numSubjects)

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
