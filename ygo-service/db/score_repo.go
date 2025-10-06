package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ygo-skc/skc-go/common/util"
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
	COALESCE(scores.score, 0) AS score,
	cards.card_number
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

	multiCardScoreQuery = `
SELECT
	score_versions.format,
	score_versions.effective_date,
	COALESCE(scores.score, 0) AS score,
	cards.card_number
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
	cards.card_number IN (%s)
ORDER BY
	score_versions.effective_date DESC;`
)

type ScoreRepository interface {
	GetDatesForFormat(context.Context, string) ([]string, *status.Status)

	GetCardScoreByID(context.Context, string) ([]*ygo.ScoreEntry, *status.Status)
	GetCardScoresByIDs(context.Context, []string) (map[string][]*ygo.ScoreEntry, *status.Status)
}
type YGOScoreRepository struct{}

func (imp YGOScoreRepository) GetDatesForFormat(ctx context.Context, format string) ([]string, *status.Status) {
	logger := cUtil.RetrieveLogger(ctx)
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

func (imp YGOScoreRepository) GetCardScoreByID(ctx context.Context, cardID string) ([]*ygo.ScoreEntry, *status.Status) {
	logger := cUtil.RetrieveLogger(ctx)
	logger.Info("Retrieving card score data")

	if rows, err := skcDBConn.Query(cardScoreQuery, cardID); err != nil {
		return nil, handleQueryError(logger, err)
	} else {
		scores := make([]*ygo.ScoreEntry, 0, 5)

		for rows.Next() {
			if score, _, err := parseRowsForScoreEntry(ctx, rows); err != nil {
				return nil, err
			} else {
				scores = append(scores, score)
			}
		}

		if len(scores) == 0 {
			logger.Error("Scores not retrieved since card ID DNE")
			return nil, status.New(codes.NotFound, "Resource not found")
		}
		return scores, nil
	}
}

func (imp YGOScoreRepository) GetCardScoresByIDs(ctx context.Context, cardIDs []string) (map[string][]*ygo.ScoreEntry, *status.Status) {
	logger := cUtil.RetrieveLogger(ctx)
	logger.Info(fmt.Sprintf("Retrieving card score data using ID's: %v", cardIDs))

	args, numCards := buildVariableQuerySubjects(cardIDs)
	query := fmt.Sprintf(multiCardScoreQuery, variablePlaceholders(numCards))

	if rows, err := skcDBConn.Query(query, args...); err != nil {
		return nil, handleQueryError(logger, err)
	} else {
		scoresByID := make(map[string][]*ygo.ScoreEntry)

		for rows.Next() {
			if score, cardID, err := parseRowsForScoreEntry(ctx, rows); err != nil {
				return nil, err
			} else {
				if _, exists := scoresByID[cardID]; !exists {
					scoresByID[cardID] = make([]*ygo.ScoreEntry, 0, 5)
				}
				scoresByID[cardID] = append(scoresByID[cardID], score)
			}
		}
		return scoresByID, nil
	}
}

func parseRowsForScoreEntry(ctx context.Context, rows *sql.Rows) (*ygo.ScoreEntry, string, *status.Status) {
	var (
		format        string
		effectiveDate string
		score         uint32
		cardID        string
	)

	if err := rows.Scan(&format, &effectiveDate, &score, &cardID); err != nil {
		return nil, "", handleRowParsingError(util.RetrieveLogger(ctx), err)
	} else {
		return &ygo.ScoreEntry{Format: format, EffectiveDate: effectiveDate, Score: score}, cardID, nil
	}
}
