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
	"github.com/ygo-skc/skc-go/common/ygo"
)

var (
	skcDBConn  *sql.DB
	spaceRegex = regexp.MustCompile(`[ ]+`)
)

const (
	// errors
	genericError = "Error occurred while querying DB"

	cardAttributes = "card_number, card_color, card_name, card_attribute, card_effect, monster_type, monster_attack, monster_defense"

	// queries
	dbVersionQuery    = "SELECT VERSION()"
	cardColorIDsQuery = "SELECT color_id, card_color from card_colors ORDER BY color_id"

	cardByCardIDQuery   = "SELECT %s FROM card_info WHERE card_number = ?"
	cardsByCardIDsQuery = "SELECT %s FROM card_info WHERE card_number IN (%s)"

	cardsByCardNamesQuery      = "SELECT %s FROM card_info WHERE card_name IN (%s)"
	searchCardUsingEffectQuery = "SELECT %s FROM card_info WHERE MATCH(card_effect) AGAINST(? IN BOOLEAN MODE) AND card_number != ? ORDER BY color_id, card_name"

	archetypeInclusionSubQuery = `SELECT %s FROM card_info WHERE MATCH (card_effect) AGAINST ('+"This card is always treated as" +"%s"' IN BOOLEAN MODE)`
	archetypeExclusionSubQuery = `SELECT %s FROM card_info WHERE MATCH (card_effect) AGAINST ('+"This card is not treated as" +"%s"'  IN BOOLEAN MODE)`

	archetypalCardsUsingCardNameQuery    = "SELECT %s FROM card_info WHERE card_name LIKE BINARY ? ORDER BY card_name"
	archetypalCardsUsingCardTextQuery    = `SELECT a.* FROM (%s) a WHERE a.card_effect REGEXP 'always treated as a.*"%s".* card' ORDER BY card_name`
	nonArchetypalCardsUsingCardTextQuery = `SELECT a.* FROM (%s) a WHERE a.card_effect REGEXP 'not treated as.*"%s".* card' ORDER BY card_name`

	randomCardQuery              = "SELECT %s FROM card_info WHERE card_color != 'Token' ORDER BY RAND() LIMIT 1"
	randomCardWithBlackListQuery = "SELECT %s FROM card_info WHERE card_number NOT IN (%s) AND card_color != 'Token' ORDER BY RAND() LIMIT 1"
)

type CardRepository interface {
	GetDBVersion(context.Context) (string, error)
	GetCardColorIDs(context.Context) (*ygo.CardColors, *model.APIError)

	GetCardByID(context.Context, string) (*ygo.Card, *model.APIError)
	GetCardsByIDs(context.Context, model.CardIDs) (*ygo.Cards, *model.APIError)

	GetCardsByNames(context.Context, model.CardNames) (*ygo.Cards, *model.APIError)
	SearchForCardRefUsingEffect(context.Context, string, string) (*ygo.CardList, *model.APIError)

	GetArchetypalCardsUsingCardName(context.Context, string) (*ygo.CardList, *model.APIError)
	GetExplicitArchetypalInclusions(context.Context, string) (*ygo.CardList, *model.APIError)
	GetExplicitArchetypalExclusions(context.Context, string) (*ygo.CardList, *model.APIError)

	GetRandomCard(context.Context, []string) (*ygo.Card, *model.APIError)
}

const (
	productDetailsQuery   = "SELECT product_id, product_locale, product_name, product_type, product_sub_type, product_release_date FROM products where product_id = ?"
	cardsByProductIDQuery = "SELECT %s, product_position, card_rarity FROM product_contents WHERE product_id= ? ORDER BY product_position"
)

type ProductRepository interface {
	GetCardsByProduct(context.Context, string) (*ygo.Product, *model.APIError)
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
		MonsterType: util.ProtoStringValue(monsterType),
		Attack:      util.ProtoUInt32Value(atk),
		Defense:     util.ProtoUInt32Value(def),
	}, nil
}

func parseRowsForCards(ctx context.Context, rows *sql.Rows, keyFn func(*ygo.Card) string) (map[string]*ygo.Card, *model.APIError) {
	cards := make(map[string]*ygo.Card)

	for rows.Next() {
		var id, color, name, attribute, effect string
		var monsterType *string
		var atk, def *uint32

		if err := rows.Scan(&id, &color, &name, &attribute, &effect, &monsterType, &atk, &def); err != nil {
			return nil, handleRowParsingError(util.LoggerFromContext(ctx), err)
		} else {
			card := model.NewYgoCardProto(id, color, name, attribute, effect, monsterType, atk, def)
			cards[keyFn(card)] = card
		}
	}

	return cards, nil
}

func parseRowsForCardList(ctx context.Context, rows *sql.Rows) ([]*ygo.Card, *model.APIError) {
	cardList := make([]*ygo.Card, 0)

	for rows.Next() {
		var id, color, name, attribute, effect string
		var monsterType *string
		var atk, def *uint32

		if err := rows.Scan(&id, &color, &name, &attribute, &effect, &monsterType, &atk, &def); err != nil {
			return nil, handleRowParsingError(util.LoggerFromContext(ctx), err)
		} else {
			cardList = append(cardList, model.NewYgoCardProto(id, color, name, attribute, effect, monsterType, atk, def))
		}
	}

	return cardList, nil
}

func queryProductInfo(logger *slog.Logger, productID string) (*ygo.Product, *model.APIError) {
	var id, locale, name, releaseDate, t, subType string

	if err := skcDBConn.QueryRow(productDetailsQuery, productID).Scan(&id, &locale, &name, &releaseDate, &t, &subType); err != nil {
		return nil, handleQueryError(logger, err)
	}
	return &ygo.Product{ID: id, Locale: locale, Name: name, ReleaseDate: releaseDate, Type: t, SubType: subType}, nil
}

func parseRowsForProductItems(ctx context.Context, rows *sql.Rows) ([]*ygo.ProductItem, map[string]uint32, *model.APIError) {
	items := make([]*ygo.ProductItem, 0)
	itemByCardIDxPosition := make(map[string]*ygo.ProductItem)
	rarityDistribution := make(map[string]uint32)

	for rows.Next() {
		var id, color, name, attribute, effect string
		var monsterType *string
		var atk, def *uint32
		var productPosition, rarity string

		if err := rows.Scan(&id, &color, &name, &attribute, &effect, &monsterType, &atk, &def, &productPosition, &rarity); err != nil {
			return nil, nil, handleRowParsingError(util.LoggerFromContext(ctx), err)
		} else {
			// either create a new ProductItem or use reference to existing Item and update the rarities
			key := fmt.Sprintf("%s-%s", id, productPosition)
			if _, exists := itemByCardIDxPosition[key]; exists {
				itemByCardIDxPosition[key].Rarities = append(itemByCardIDxPosition[key].Rarities, rarity)
			} else {
				item := &ygo.ProductItem{
					Card:     model.NewYgoCardProto(id, color, name, attribute, effect, monsterType, atk, def),
					Position: productPosition,
					Rarities: []string{rarity},
				}
				items = append(items, item)
				itemByCardIDxPosition[key] = item
			}

			// running total of all rarities
			if num, exists := rarityDistribution[rarity]; exists {
				rarityDistribution[rarity] = num + 1
			} else {
				rarityDistribution[rarity] = 1
			}
		}
	}

	return items, rarityDistribution, nil
}
