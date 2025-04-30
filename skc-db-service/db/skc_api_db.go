package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	cModel "github.com/ygo-skc/skc-go/common/model"
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
	queryDBVersion    = "SELECT VERSION()"
	queryCardColorIDs = "SELECT color_id, card_color from card_colors ORDER BY color_id"

	queryCardUsingCardID           = "SELECT card_number, card_color, card_name, card_attribute, card_effect, monster_type, monster_attack, monster_defense FROM card_info WHERE card_number = ?"
	queryCardUsingCardIDs          = "SELECT card_number, card_color, card_name, card_attribute, card_effect, monster_type, monster_attack, monster_defense FROM card_info WHERE card_number IN (%s)"
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

func handleQueryError(logger *slog.Logger, err error) *cModel.APIError {
	logger.Error(fmt.Sprintf("Error fetching data from DB - %v", err))
	return &cModel.APIError{Message: genericError, StatusCode: http.StatusInternalServerError}
}

func handleRowParsingError(logger *slog.Logger, err error) *cModel.APIError {
	logger.Error(fmt.Sprintf("Error parsing data from DB - %v", err))
	return &cModel.APIError{Message: genericError, StatusCode: http.StatusInternalServerError}
}

// interface
type SKCDatabaseAccessObject interface {
	GetSKCDBVersion(context.Context) (string, error)

	GetDesiredCardInDBUsingID(context.Context, string) (cModel.Card, *cModel.APIError)
	GetDesiredCardInDBUsingMultipleCardIDs(context.Context, []string) (cModel.BatchCardData[cModel.CardIDs], *cModel.APIError)
}

// impl
type SKCDAOImplementation struct{}

// Get version of MYSQL being used by SKC DB.
func (imp SKCDAOImplementation) GetSKCDBVersion(ctx context.Context) (string, error) {
	var version string
	if err := skcDBConn.QueryRow(queryDBVersion).Scan(&version); err != nil {
		cUtil.LoggerFromContext(ctx).Info(fmt.Sprintf("Error getting SKC DB version - %v", err))
		return version, &cModel.APIError{Message: genericError, StatusCode: http.StatusInternalServerError}
	}

	return version, nil
}

// Leverages GetDesiredCardInDBUsingMultipleCardIDs to get information on a specific card using its identifier
func (imp SKCDAOImplementation) GetDesiredCardInDBUsingID(ctx context.Context, cardID string) (cModel.Card, *cModel.APIError) {
	if results, err := imp.GetDesiredCardInDBUsingMultipleCardIDs(ctx, []string{cardID}); err != nil {
		return cModel.Card{}, err
	} else {
		if card, exists := results.CardInfo[cardID]; !exists {
			return cModel.Card{}, &cModel.APIError{Message: fmt.Sprintf("No results found when querying by card ID %s", cardID), StatusCode: http.StatusNotFound}
		} else {
			return card, nil
		}
	}
}

func (imp SKCDAOImplementation) GetDesiredCardInDBUsingMultipleCardIDs(ctx context.Context, cardIDs []string) (cModel.BatchCardData[cModel.CardIDs], *cModel.APIError) {
	logger := cUtil.LoggerFromContext(ctx)
	logger.Info("Retrieving card data from DB")

	args, numCards := buildVariableQuerySubjects(cardIDs)
	cardData := make(cModel.CardDataMap, numCards) // used to store results

	query := fmt.Sprintf(queryCardUsingCardIDs, variablePlaceholders(numCards))

	if rows, err := skcDBConn.Query(query, args...); err != nil {
		return cModel.BatchCardData[cModel.CardIDs]{}, handleQueryError(logger, err)
	} else {
		if cards, err := parseRowsForCard(ctx, rows); err != nil {
			return cModel.BatchCardData[cModel.CardIDs]{}, err
		} else {
			for _, card := range cards {
				cardData[card.ID] = card
			}
		}
	}

	return cModel.BatchCardData[cModel.CardIDs]{CardInfo: cardData, UnknownResources: cardData.FindMissingIDs(cardIDs)}, nil
}

func parseRowsForCard(ctx context.Context, rows *sql.Rows) ([]cModel.Card, *cModel.APIError) {
	logger := cUtil.LoggerFromContext(ctx)
	cards := []cModel.Card{}

	for rows.Next() {
		var card cModel.Card
		if err := rows.Scan(&card.ID, &card.Color, &card.Name, &card.Attribute, &card.Effect, &card.MonsterType, &card.Attack, &card.Defense); err != nil {
			return nil, handleRowParsingError(logger, err)
		} else {
			cards = append(cards, card)
		}
	}

	return cards, nil // no parsing error
}
