package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ygo-skc/skc-go/common/v2/model"
	"github.com/ygo-skc/skc-go/common/v2/util"
	cUtil "github.com/ygo-skc/skc-go/common/v2/util"
	"github.com/ygo-skc/skc-go/common/v2/ygo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	productDetailsQuery = `
SELECT
	product_id,
	product_locale,
	product_name,
	product_type,
	product_sub_type,
	product_release_date
FROM
	products
WHERE
	product_id = ?`

	cardsByProductIDQuery = `
SELECT
	%s,
	product_position,
	card_rarity
FROM
	product_contents
WHERE
	product_id = ?
ORDER BY
	product_position`

	productInfoByIDs = `
SELECT
	product_id,
	product_locale,
	product_name,
	product_type,
	product_sub_type,
	product_release_date,
	product_content_total
FROM
	product_info
WHERE
	product_id IN (%s)`
)

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

type ProductRepository interface {
	GetCardsByProductID(context.Context, string) (*ygo.Product, *status.Status)

	GetProductSummaryByID(context.Context, string) (*ygo.ProductSummary, *status.Status)
	GetProductsSummaryByID(context.Context, model.ProductIDs) (*ygo.Products, *status.Status)
}
type YGOProductRepository struct{}

func (imp YGOProductRepository) GetCardsByProductID(ctx context.Context, productID string) (*ygo.Product, *status.Status) {
	logger := cUtil.RetrieveLogger(ctx)
	logger.Info(fmt.Sprintf("Retrieving product data using ID %s", productID))

	if product, err := queryProductInfo(logger, productID); err != nil {
		return nil, err
	} else {
		query := fmt.Sprintf(cardsByProductIDQuery, cardAttributes)
		if rows, err := skcDBConn.Query(query, productID); err != nil {
			return nil, handleQueryError(logger, err)
		} else {
			if items, rarityDistribution, err := parseRowsForProductItems(ctx, rows); err != nil {
				return nil, err
			} else {
				product.Items = items
				product.TotalItems = uint32(len(items))
				product.RarityDistribution = rarityDistribution
				return product, nil
			}
		}
	}
}

func (imp YGOProductRepository) GetProductSummaryByID(ctx context.Context, productID string) (*ygo.ProductSummary, *status.Status) {
	if results, err := imp.GetProductsSummaryByID(ctx, []string{productID}); err != nil {
		return nil, err
	} else {
		if product, exists := results.Products[productID]; !exists {
			return nil, status.New(codes.NotFound, "No results found")
		} else {
			return product, nil
		}
	}
}

func (imp YGOProductRepository) GetProductsSummaryByID(ctx context.Context, products model.ProductIDs) (*ygo.Products, *status.Status) {
	logger := cUtil.RetrieveLogger(ctx)
	logger.Info(fmt.Sprintf("Retrieving summary of the following products: %v", products))

	args, numProducts := buildVariableQuerySubjects(products)
	productData := make(map[string]*ygo.ProductSummary, numProducts)

	query := fmt.Sprintf(productInfoByIDs, variablePlaceholders(numProducts))

	if rows, err := skcDBConn.Query(query, args...); err != nil {
		return nil, handleQueryError(logger, err)
	} else {
		for rows.Next() {
			var id, locale, name, t, subType, releaseDate string
			var totalItems uint32

			if err := rows.Scan(&id, &locale, &name, &t, &subType, &releaseDate, &totalItems); err != nil {
				return nil, handleRowParsingError(logger, err)
			}

			productData[id] = &ygo.ProductSummary{ID: id, Locale: locale, Name: name, Type: t, SubType: subType, ReleaseDate: releaseDate, TotalItems: totalItems}
		}
	}

	return &ygo.Products{
		Products:         productData,
		UnknownResources: model.FindMissingKeys(productData, products),
	}, nil
}
