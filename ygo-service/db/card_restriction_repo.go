package db

import (
	"context"

	"github.com/ygo-skc/skc-go/common/v2/util"
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
	logger.Info("Retrieving effective dates")

	if rows, err := skcDBConn.Query(datesForFormatQuery, format); err != nil {
		return nil, handleQueryError(logger, err)
	} else {
		scores := make([]string, 0, 5)
		var date string

		for rows.Next() {
			if err := rows.Scan(&date); err != nil {
				return nil, handleRowParsingError(logger, err)
			} else {
				scores = append(scores, date)
			}
		}
		return scores, nil
	}
}
