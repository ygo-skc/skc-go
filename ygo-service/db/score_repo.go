package db

import (
	"context"
	"fmt"

	cUtil "github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/common/ygo"
	"google.golang.org/grpc/codes"
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

	cardScoreQuery = `
SELECT
	score_versions.format,
	score_versions.effective_date,
	COALESCE(scores.score, 0) AS score
FROM
	(
		SELECT DISTINCT
			FORMAT,
			effective_date
		FROM
			card_scores
	) score_versions
	CROSS JOIN cards
	LEFT JOIN card_scores scores ON scores.format = score_versions.format
	AND scores.effective_date = score_versions.effective_date
	AND scores.card_number = cards.card_number
WHERE
	cards.card_number = ?
ORDER BY
	score_versions.effective_date DESC;`
)

type ScoreRepository interface {
	GetDatesForFormat(context.Context, string) ([]string, *status.Status)

	GetCardScoreByID(context.Context, string) ([]*ygo.ScoreInstance, *status.Status)
	GetCardScoresByIDs(context.Context, string) (*ygo.Product, *status.Status)
}
type YGOScoreRepository struct{}

func (imp YGOScoreRepository) GetDatesForFormat(ctx context.Context, format string) ([]string, *status.Status) {
	logger := cUtil.RetrieveLogger(ctx)
	logger.Info(fmt.Sprintf("Retrieving dates for ygo format %s", format))

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

func (imp YGOScoreRepository) GetCardScoreByID(ctx context.Context, cardID string) ([]*ygo.ScoreInstance, *status.Status) {
	logger := cUtil.RetrieveLogger(ctx)
	logger.Info(fmt.Sprintf("Retrieving card score data using ID %s", cardID))

	if rows, err := skcDBConn.Query(cardScoreQuery, cardID); err != nil {
		return nil, handleQueryError(logger, err)
	} else {
		scores := make([]*ygo.ScoreInstance, 0, 5)
		var (
			format        string
			effectiveDate string
			score         uint32
		)

		for rows.Next() {
			if err := rows.Scan(&format, &effectiveDate, &score); err != nil {
				return nil, handleRowParsingError(logger, err)
			} else {
				scores = append(scores, &ygo.ScoreInstance{Format: format, EffectiveDate: effectiveDate, Score: score})
			}
		}

		if len(scores) == 0 {
			logger.Error(fmt.Sprintf("Scores not retrieved since card ID %s DNE", cardID))
			return nil, status.New(codes.NotFound, "Resource not found")
		}
		return scores, nil
	}
}
