package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ygo-skc/skc-go/common/v2/model"
	"github.com/ygo-skc/skc-go/common/v2/util"
	"github.com/ygo-skc/skc-go/common/v2/ygo"
	"google.golang.org/grpc/status"
)

const (
	cardScoreByFormatAndDateQuery = `
SELECT
	ci.card_number,
	card_color,
	card_name,
	card_attribute,
	card_effect,
	monster_type,
	monster_attack,
	monster_defense,
	score
FROM
	card_scores AS cs FORCE INDEX (FORMAT)
	JOIN card_info AS ci ON ci.card_number = cs.card_number
WHERE
	cs.format = ?
	AND cs.effective_date = ?
ORDER BY
	%s`

	cardScoreQuery = `
SELECT
	score_versions.format,
	score_versions.effective_date,
	COALESCE(scores.score, 0) AS score,
	? AS card_number
FROM
	(
		SELECT DISTINCT
			FORMAT,
			effective_date
		FROM
			card_scores
	) AS score_versions
	LEFT JOIN card_scores AS scores ON scores.format = score_versions.format
	AND scores.effective_date = score_versions.effective_date
	AND scores.card_number = ?
ORDER BY
	score_versions.effective_date DESC`

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
	AND scores.card_number = cards.card_number
	AND scores.effective_date = score_versions.effective_date
WHERE
	cards.card_number IN (%s)
ORDER BY
	score_versions.effective_date DESC`
)

type ScoreRepository interface {
	GetScoresByFormatAndDate(context.Context, string, string, ygo.CardRestrictionSortOrder) ([]*ygo.CardScoreEntry, uint32, *status.Status)

	GetCardScoreByID(context.Context, string, time.Time, func(*ygo.CardScore, *ygo.ScoreEntry, time.Time)) (*ygo.CardScore, *status.Status)
	GetCardScoresByIDs(context.Context, []string, time.Time, func(*ygo.CardScore, *ygo.ScoreEntry, time.Time)) (map[string]*ygo.CardScore, *status.Status)
}
type YGOScoreRepository struct{}

func (imp YGOScoreRepository) GetScoresByFormatAndDate(
	ctx context.Context, format string, effectiveDate string, sortOrder ygo.CardRestrictionSortOrder) ([]*ygo.CardScoreEntry, uint32, *status.Status) {

	logger := util.RetrieveLogger(ctx)
	logger.Info(fmt.Sprintf("Retrieving scores using format %s and date %s", format, effectiveDate))

	var sortingSubQuery string
	switch sortOrder {
	case ygo.CardRestrictionSortOrder_CARD_NAME:
		sortingSubQuery = "card_name"
	case ygo.CardRestrictionSortOrder_SCORE_THEN_COLOR:
		sortingSubQuery = "score DESC, card_color, card_name"
	}

	query := fmt.Sprintf(cardScoreByFormatAndDateQuery, sortingSubQuery)
	if rows, err := skcDBConn.Query(query, format, effectiveDate); err != nil {
		return make([]*ygo.CardScoreEntry, 0), 0, handleQueryError(logger, err)
	} else {
		var (
			id, color, name, attribute, effect string
			monsterType                        *string
			atk, def                           *uint32
			score                              uint32
		)
		entries := make([]*ygo.CardScoreEntry, 0, 600)
		var numEntries uint32

		for rows.Next() {
			if err := rows.Scan(&id, &color, &name, &attribute, &effect, &monsterType, &atk, &def, &score); err != nil {
				return make([]*ygo.CardScoreEntry, 0), 0, handleRowParsingError(util.RetrieveLogger(ctx), err)
			}
			entries = append(entries, &ygo.CardScoreEntry{
				Card: model.NewYGOCardProtoBuilder(id, name).
					WithColor(color).
					WithAttribute(attribute).
					WithEffect(effect).
					WithMonsterType(monsterType).
					WithAttack(atk).
					WithDefense(def).
					Build(),
				Score: score,
			})
			numEntries++
		}
		return entries, numEntries, nil
	}
}

func (imp YGOScoreRepository) GetCardScoreByID(ctx context.Context, cardID string, todaysDate time.Time,
	parser func(*ygo.CardScore, *ygo.ScoreEntry, time.Time)) (*ygo.CardScore, *status.Status) {

	logger := util.RetrieveLogger(ctx)
	logger.Info("Retrieving card score data")

	if rows, err := skcDBConn.Query(cardScoreQuery, cardID, cardID); err != nil {
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
