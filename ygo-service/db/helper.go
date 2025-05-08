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
	cUtil "github.com/ygo-skc/skc-go/common/util"
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
	queryCardColorIDs = "SELECT color_id, card_color from card_colors ORDER BY color_id"

	cardByCardIDQuery              = "SELECT card_number, card_color, card_name, card_attribute, card_effect, monster_type, monster_attack, monster_defense FROM card_info WHERE card_number = ?"
	cardsByCardIDsQuery            = "SELECT card_number, card_color, card_name, card_attribute, card_effect, monster_type, monster_attack, monster_defense FROM card_info WHERE card_number IN (%s)"
	queryCardUsingCardNames        = "SELECT card_number, card_color, card_name, card_attribute, card_effect, monster_type, monster_attack, monster_defense FROM card_info WHERE card_name IN (%s)"
	queryCardsUsingProductID       = "SELECT DISTINCT(card_number), card_color,card_name,card_attribute,card_effect,monster_type,monster_attack,monster_defense FROM product_contents WHERE product_id= ? ORDER BY card_name"
	queryRandomCardID              = "SELECT card_number FROM card_info WHERE card_color != 'Token' ORDER BY RAND() LIMIT 1"
	queryRandomCardIDWithBlackList = "SELECT card_number FROM card_info WHERE card_number NOT IN (%s) AND card_color != 'Token' ORDER BY RAND() LIMIT 1"

	queryCardsInArchetypeUsingName  = "SELECT card_number, card_color, card_name, card_attribute, card_effect, monster_type, monster_attack, monster_defense FROM card_info WHERE card_name LIKE BINARY ? ORDER BY card_name"
	queryCardsTreatedAsArchetype    = "SELECT card_number, card_color, card_name, card_attribute, card_effect, monster_type, monster_attack, monster_defense FROM card_info WHERE MATCH(card_effect) AGAINST(? IN BOOLEAN MODE) ORDER BY card_name"
	queryCardsNotTreatedAsArchetype = "SELECT card_number, card_color, card_name, card_attribute, card_effect, monster_type, monster_attack, monster_defense FROM card_info WHERE MATCH(card_effect) AGAINST(? IN BOOLEAN MODE) ORDER BY card_name"

	findRelatedCardsUsingCardEffect = "SELECT card_number, card_color, card_name, card_attribute, card_effect, monster_type, monster_attack, monster_defense FROM card_info WHERE MATCH(card_effect) AGAINST(? IN BOOLEAN MODE) AND card_number != ? ORDER BY color_id, card_name"
)

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
	return &model.APIError{Message: genericError, StatusCode: http.StatusInternalServerError}
}

func handleRowParsingError(logger *slog.Logger, err error) *model.APIError {
	logger.Error(fmt.Sprintf("Error parsing data from DB - %v", err))
	return &model.APIError{Message: genericError, StatusCode: http.StatusInternalServerError}
}

func parseRowsForCard(ctx context.Context, rows *sql.Rows) ([]model.YGOCardREST, *model.APIError) {
	logger := cUtil.LoggerFromContext(ctx)
	cards := []model.YGOCardREST{}

	for rows.Next() {
		var card model.YGOCardREST
		if err := rows.Scan(&card.ID, &card.Color, &card.Name, &card.Attribute, &card.Effect, &card.MonsterType, &card.Attack, &card.Defense); err != nil {
			return nil, handleRowParsingError(logger, err)
		} else {
			cards = append(cards, card)
		}
	}

	return cards, nil // no parsing error
}
