package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ygo-skc/skc-go/common/v2/util"
	"github.com/ygo-skc/skc-go/common/v2/ygo"
	"google.golang.org/grpc/status"
)

const (
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
	GetScoresByFormatAndDate(context.Context, string, string) ([]*ygo.CardScoreEntry, *status.Status)

	GetCardScoreByID(context.Context, string, time.Time, func(*ygo.CardScore, *ygo.ScoreEntry, time.Time)) (*ygo.CardScore, *status.Status)
	GetCardScoresByIDs(context.Context, []string, time.Time, func(*ygo.CardScore, *ygo.ScoreEntry, time.Time)) (map[string]*ygo.CardScore, *status.Status)
}
type YGOScoreRepository struct{}

func (imp YGOScoreRepository) GetDatesForFormat(ctx context.Context, format string) ([]string, *status.Status) {
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

func (imp YGOScoreRepository) GetCardScoreByID(ctx context.Context, cardID string, todaysDate time.Time,
	parser func(*ygo.CardScore, *ygo.ScoreEntry, time.Time)) (*ygo.CardScore, *status.Status) {

	logger := util.RetrieveLogger(ctx)
	logger.Info("Retrieving card score data")

	if rows, err := skcDBConn.Query(cardScoreQuery, cardID); err != nil {
		return nil, handleQueryError(logger, err)
	} else {
		score := &ygo.CardScore{
			CurrentScoreByFormat: make(map[string]uint32, 3),
			UniqueFormats:        make([]string, 0, 3),
			ScheduledChanges:     make([]string, 0, 3),
			ScoreHistory:         make([]*ygo.ScoreEntry, 0, 5),
		}

		for rows.Next() {
			if entry, _, err := parseRowsForScoreEntry(ctx, rows); err != nil {
				return nil, err
			} else {
				parser(score, entry, todaysDate)
			}
		}
		return score, nil
	}
}

func (imp YGOScoreRepository) GetCardScoresByIDs(ctx context.Context, cardIDs []string, todaysDate time.Time,
	parser func(*ygo.CardScore, *ygo.ScoreEntry, time.Time)) (map[string]*ygo.CardScore, *status.Status) {

	logger := util.RetrieveLogger(ctx)
	logger.Info(fmt.Sprintf("Retrieving card score data using ID's: %v", cardIDs))

	args, numCards := buildVariableQuerySubjects(cardIDs)
	query := fmt.Sprintf(multiCardScoreQuery, variablePlaceholders(numCards))

	if rows, err := skcDBConn.Query(query, args...); err != nil {
		return nil, handleQueryError(logger, err)
	} else {
		scoresByID := make(map[string]*ygo.CardScore)

		for rows.Next() {
			if score, cardID, err := parseRowsForScoreEntry(ctx, rows); err != nil {
				return nil, err
			} else {
				if _, exists := scoresByID[cardID]; !exists {
					scoresByID[cardID] = &ygo.CardScore{
						CurrentScoreByFormat: make(map[string]uint32, 3),
						UniqueFormats:        make([]string, 0, 3),
						ScheduledChanges:     make([]string, 0, 3),
						ScoreHistory:         make([]*ygo.ScoreEntry, 0, 5),
					}
				}
				parser(scoresByID[cardID], score, todaysDate)
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
