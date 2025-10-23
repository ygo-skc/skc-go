package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/ygo-skc/skc-go/common/v2/model"
	"github.com/ygo-skc/skc-go/common/v2/util"
	"github.com/ygo-skc/skc-go/common/v2/ygo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	genericError = "Error occurred while querying DB"
)

func handleQueryError(logger *slog.Logger, err error) *status.Status {
	logger.Error(fmt.Sprintf("Error fetching data from DB - %v", err))

	if err == sql.ErrNoRows {
		return status.New(codes.NotFound, "No results found")
	}
	return status.New(codes.Internal, genericError)
}

func handleRowParsingError(logger *slog.Logger, err error) *status.Status {
	logger.Error(fmt.Sprintf("Error parsing data from DB - %v", err))
	return status.New(codes.Internal, genericError)
}
func queryCard(logger *slog.Logger, query string, args []interface{}) (*ygo.Card, *status.Status) {
	var (
		id, color, name, attribute, effect string
		monsterType                        *string
		atk, def                           *uint32
	)

	if err := skcDBConn.QueryRow(query, args...).Scan(&id, &color, &name, &attribute, &effect, &monsterType, &atk, &def); err != nil {
		return nil, handleQueryError(logger, err)
	}

	card := model.NewYGOCardProtoBuilder(id, name).WithColor(color).
		WithAttribute(attribute).WithEffect(effect).WithMonsterType(monsterType).WithAttack(atk).WithDefense(def).Build()
	return card, nil
}

func parseRowsForCards(ctx context.Context, rows *sql.Rows, keyFn func(*ygo.Card) string) (map[string]*ygo.Card, *status.Status) {
	cards := make(map[string]*ygo.Card)

	var (
		id, color, name, attribute, effect string
		monsterType                        *string
		atk, def                           *uint32
	)
	for rows.Next() {
		if err := rows.Scan(&id, &color, &name, &attribute, &effect, &monsterType, &atk, &def); err != nil {
			return nil, handleRowParsingError(util.RetrieveLogger(ctx), err)
		} else {
			card := model.NewYGOCardProtoBuilder(id, name).WithColor(color).
				WithAttribute(attribute).WithEffect(effect).WithMonsterType(monsterType).WithAttack(atk).WithDefense(def).Build()
			cards[keyFn(card)] = card
		}
	}

	return cards, nil
}

func parseRowsForCardList(ctx context.Context, rows *sql.Rows) ([]*ygo.Card, *status.Status) {
	cardList := make([]*ygo.Card, 0)

	var (
		id, color, name, attribute, effect string
		monsterType                        *string
		atk, def                           *uint32
	)
	for rows.Next() {
		if err := rows.Scan(&id, &color, &name, &attribute, &effect, &monsterType, &atk, &def); err != nil {
			return nil, handleRowParsingError(util.RetrieveLogger(ctx), err)
		} else {
			cardList = append(cardList,
				model.NewYGOCardProtoBuilder(id, name).WithColor(color).
					WithAttribute(attribute).WithEffect(effect).WithMonsterType(monsterType).WithAttack(atk).WithDefense(def).Build())
		}
	}

	return cardList, nil
}

func queryProductInfo(logger *slog.Logger, productID string) (*ygo.Product, *status.Status) {
	var id, locale, name, t, subType, releaseDate string

	if err := skcDBConn.QueryRow(productDetailsQuery, productID).Scan(&id, &locale, &name, &t, &subType, &releaseDate); err != nil {
		return nil, handleQueryError(logger, err)
	}
	return &ygo.Product{ID: id, Locale: locale, Name: name, ReleaseDate: releaseDate, Type: t, SubType: subType}, nil
}

func parseRowsForProductItems(ctx context.Context, rows *sql.Rows) ([]*ygo.ProductItem, map[string]uint32, *status.Status) {
	items := make([]*ygo.ProductItem, 0)
	itemByCardIDxPosition := make(map[string]*ygo.ProductItem)
	rarityDistribution := make(map[string]uint32)
	var (
		id, color, name, attribute, effect string
		monsterType                        *string
		atk, def                           *uint32
		productPosition, rarity            string
	)
	for rows.Next() {
		if err := rows.Scan(&id, &color, &name, &attribute, &effect, &monsterType, &atk, &def, &productPosition, &rarity); err != nil {
			return nil, nil, handleRowParsingError(util.RetrieveLogger(ctx), err)
		} else {
			// either create a new ProductItem or use reference to existing Item and update the rarities
			key := fmt.Sprintf("%s-%s", id, productPosition)
			if _, exists := itemByCardIDxPosition[key]; exists {
				itemByCardIDxPosition[key].Rarities = append(itemByCardIDxPosition[key].Rarities, rarity)
			} else {
				item := &ygo.ProductItem{
					Card: model.NewYGOCardProtoBuilder(id, name).WithColor(color).
						WithAttribute(attribute).WithEffect(effect).WithMonsterType(monsterType).WithAttack(atk).WithDefense(def).Build(),
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

func convertToFullText(subject string) string {
	fullTextSubject := spaceRegex.ReplaceAllString(strings.ReplaceAll(subject, "-", " "), " +")
	return fmt.Sprintf(`"+%s"`, fullTextSubject) // match phrase, not all words in text will match only consecutive matches of words in phrase
}

func buildVariableQuerySubjects(subjects []string) ([]interface{}, int) {
	numSubjects := len(subjects)
	args := make([]interface{}, numSubjects)

	for index, subject := range subjects {
		args[index] = subject
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
