package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/ygo-skc/skc-go/common/model"
	"github.com/ygo-skc/skc-go/common/util"
	cUtil "github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/common/ygo"
)

var (
	skcDBConn  *sql.DB
	spaceRegex = regexp.MustCompile(`[ ]+`)
)

const (
	// errors
	genericError = "Error occurred while querying DB"

	// queries
	dbVersionQuery    = "SELECT VERSION()"
	cardColorIDsQuery = "SELECT color_id, card_color from card_colors ORDER BY color_id"

	cardByCardIDQuery   = "SELECT card_number, card_color, card_name, card_attribute, card_effect, monster_type, monster_attack, monster_defense FROM card_info WHERE card_number = ?"
	cardsByCardIDsQuery = "SELECT card_number, card_color, card_name, card_attribute, card_effect, monster_type, monster_attack, monster_defense FROM card_info WHERE card_number IN (%s)"

	cardsByCardNamesQuery = "SELECT card_number, card_color, card_name, card_attribute, card_effect, monster_type, monster_attack, monster_defense FROM card_info WHERE card_name IN (%s)"

	randomCardQuery              = "SELECT card_number, card_color, card_name, card_attribute, card_effect, monster_type, monster_attack, monster_defense FROM card_info WHERE card_color != 'Token' ORDER BY RAND() LIMIT 1"
	randomCardWithBlackListQuery = "SELECT card_number, card_color, card_name, card_attribute, card_effect, monster_type, monster_attack, monster_defense FROM card_info WHERE card_number NOT IN (%s) AND card_color != 'Token' ORDER BY RAND() LIMIT 1"
)

type CardRepository interface {
	GetDBVersion(context.Context) (string, error)
	GetCardColorIDs(context.Context) (*ygo.CardColors, *model.APIError)

	GetCardByID(context.Context, string) (*ygo.Card, *model.APIError)
	GetCardsByIDs(context.Context, model.CardIDs) (*ygo.Cards, *model.APIError)

	GetCardsByNames(context.Context, model.CardNames) (*ygo.Cards, *model.APIError)

	GetRandomCard(context.Context, []string) (*ygo.Card, *model.APIError)
}

func convertToFullText(subject string) string {
	fullTextSubject := spaceRegex.ReplaceAllString(strings.ReplaceAll(subject, "-", " "), " +")
	return fmt.Sprintf(`"+%s"`, fullTextSubject) // match phrase, not all words in text will match only consecutive matches of words in phrase
}

func buildVariableQuerySubjects(subjects []string) ([]interface{}, int) {
	numSubjects := len(subjects)
	args := make([]interface{}, numSubjects)

	for index, cardId := range subjects {
		args[index] = cardId
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

func queryCard(logger *slog.Logger, query string, args []interface{}) (*ygo.Card, *model.APIError) {
	var id, color, name, attribute, effect string
	var monsterType *string
	var atk, def *uint32

	if err := skcDBConn.QueryRow(query, args...).Scan(&id, &color, &name, &attribute, &effect, &monsterType, &atk, &def); err != nil {
		return nil, handleQueryError(logger, err)
	}

	return &ygo.Card{
		ID:          id,
		Color:       color,
		Name:        name,
		Attribute:   attribute,
		Effect:      effect,
		MonsterType: util.PBStringValue(monsterType),
		Attack:      util.PBUInt32Value(atk),
		Defense:     util.PBUInt32Value(def),
	}, nil
}

func parseRowsForCards(ctx context.Context, rows *sql.Rows, keyFn func(*ygo.Card) string) (*map[string]*ygo.Card, *model.APIError) {
	cards := make(map[string]*ygo.Card)

	for rows.Next() {
		var id, color, name, attribute, effect string
		var monsterType *string
		var atk, def *uint32

		if err := rows.Scan(&id, &color, &name, &attribute, &effect, &monsterType, &atk, &def); err != nil {
			return nil, handleRowParsingError(cUtil.LoggerFromContext(ctx), err)
		} else {
			card := &ygo.Card{
				ID:          id,
				Color:       color,
				Name:        name,
				Attribute:   attribute,
				Effect:      effect,
				MonsterType: util.PBStringValue(monsterType),
				Attack:      util.PBUInt32Value(atk),
				Defense:     util.PBUInt32Value(def),
			}
			cards[keyFn(card)] = card
		}
	}

	return &cards, nil // no parsing error
}
