package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/ygo-skc/skc-go/common/model"
	"github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/common/ygo"
	"google.golang.org/grpc/status"
)

func queryCard(logger *slog.Logger, query string, args []interface{}) (*ygo.Card, *status.Status) {
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

func parseRowsForCards(ctx context.Context, rows *sql.Rows, keyFn func(*ygo.Card) string) (map[string]*ygo.Card, *status.Status) {
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

func parseRowsForCardList(ctx context.Context, rows *sql.Rows) ([]*ygo.Card, *status.Status) {
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

func queryProductInfo(logger *slog.Logger, productID string) (*ygo.Product, *status.Status) {
	var id, locale, name, releaseDate, t, subType string

	if err := skcDBConn.QueryRow(productDetailsQuery, productID).Scan(&id, &locale, &name, &releaseDate, &t, &subType); err != nil {
		return nil, handleQueryError(logger, err)
	}
	return &ygo.Product{ID: id, Locale: locale, Name: name, ReleaseDate: releaseDate, Type: t, SubType: subType}, nil
}

func parseRowsForProductItems(ctx context.Context, rows *sql.Rows) ([]*ygo.ProductItem, map[string]uint32, *status.Status) {
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
