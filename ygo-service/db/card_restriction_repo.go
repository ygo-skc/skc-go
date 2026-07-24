package db

import (
	"context"
	"log/slog"

	"github.com/ygo-skc/skc-go/common/v3/util"
	"google.golang.org/grpc/status"
)

const (
	datesForFormatQuery = `
SELECT
	UNIQUE effective_date
FROM
	card_scores
WHERE
	format = ?
ORDER BY
	effective_date DESC;`
)

type CardRestrictionRepository interface {
	GetDatesForFormat(context.Context, string) ([]string, *status.Status)
}
type YGOCardRestrictionRepository struct{}

func (imp YGOCardRestrictionRepository) GetDatesForFormat(ctx context.Context, format string) ([]string, *status.Status) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Retrieving effective dates", slog.String("format", format))

	rows, err := skcDBConn.QueryContext(ctx, datesForFormatQuery, format)
	if err != nil {
		return nil, handleQueryError(logger, err)
	}
	defer rows.Close()

	scores := make([]string, 0, 5)
	var date string

	for rows.Next() {
		if err := rows.Scan(&date); err != nil {
			return nil, handleRowParsingError(logger, err)
		}
		scores = append(scores, date)
	}

	if err := rows.Err(); err != nil {
		return nil, handleQueryError(logger, err)
	}
	return scores, nil
}
