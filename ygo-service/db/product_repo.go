package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/ygo-skc/skc-go/common/v3/model"
	"github.com/ygo-skc/skc-go/common/v3/util"
	"github.com/ygo-skc/skc-go/common/v3/ygo"
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
	product_id, product_locale, product_name, product_type, product_sub_type, product_release_date, product_content_total
FROM
	product_info
WHERE
	product_id IN (%s)`

	productReleasedOnDate = `
SELECT 
	product_id, product_locale, product_name, product_type, product_sub_type, product_release_date, product_content_total
FROM
	product_info
WHERE
	DATE_FORMAT(product_release_date, '%m-%d') = ?`
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
		}

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
		rarityDistribution[rarity]++
	}

	if err := rows.Err(); err != nil {
		return nil, nil, handleQueryError(util.RetrieveLogger(ctx), err)
	}
	return items, rarityDistribution, nil
}

type ProductRepository interface {
	GetCardsByProductID(context.Context, string) (*ygo.Product, *status.Status)

	GetProductSummaryByID(context.Context, string) (*ygo.ProductSummary, *status.Status)
	GetProductsSummaryByID(context.Context, model.ProductIDs) (*ygo.Products, *status.Status)

	GetProductsReleasedSameDay(context.Context, time.Time) ([]*ygo.ProductSummary, *status.Status)
}
type YGOProductRepository struct{}

func (imp YGOProductRepository) GetCardsByProductID(ctx context.Context, productID string) (*ygo.Product, *status.Status) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Retrieving product data", slog.String("product_id", productID))

	product, err := queryProductInfo(ctx, logger, productID)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(cardsByProductIDQuery, cardAttributes)
	rows, dbErr := skcDBConn.QueryContext(ctx, query, productID)
	if dbErr != nil {
		return nil, handleQueryError(logger, dbErr)
	}
	defer rows.Close()

	items, rarityDistribution, err := parseRowsForProductItems(ctx, rows)
	if err != nil {
		return nil, err
	}

	product.Items = items
	product.TotalItems = uint32(len(items))
	product.RarityDistribution = rarityDistribution
	return product, nil
}

func (imp YGOProductRepository) GetProductSummaryByID(ctx context.Context, productID string) (*ygo.ProductSummary, *status.Status) {
	results, err := imp.GetProductsSummaryByID(ctx, []string{productID})
	if err != nil {
		return nil, err
	}

	product, exists := results.Products[productID]
	if !exists {
		return nil, status.New(codes.NotFound, "No results found")
	}
	return product, nil
}

func (imp YGOProductRepository) GetProductsSummaryByID(ctx context.Context, products model.ProductIDs) (*ygo.Products, *status.Status) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Retrieving product summaries", slog.Any("product_ids", products))

	if len(products) == 0 {
		return &ygo.Products{Products: make(map[string]*ygo.ProductSummary)}, nil
	}

	args := buildVariableQuerySubjects(products)
	productData := make(map[string]*ygo.ProductSummary, len(products))

	query := fmt.Sprintf(productInfoByIDs, variablePlaceholders(len(products)))

	rows, err := skcDBConn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, handleQueryError(logger, err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, locale, name, t, subType, releaseDate string
		var totalItems uint32

		if err := rows.Scan(&id, &locale, &name, &t, &subType, &releaseDate, &totalItems); err != nil {
			return nil, handleRowParsingError(logger, err)
		}

		productData[id] = &ygo.ProductSummary{Id: id, Locale: locale, Name: name, Type: t, SubType: subType, ReleaseDate: releaseDate, TotalItems: totalItems}
	}

	if err := rows.Err(); err != nil {
		return nil, handleQueryError(logger, err)
	}

	return &ygo.Products{
		Products:         productData,
		UnknownResources: model.FindMissingKeys(productData, products),
	}, nil
}

func (imp YGOProductRepository) GetProductsReleasedSameDay(ctx context.Context, date time.Time) ([]*ygo.ProductSummary, *status.Status) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Retrieving products released on the same month/day provided",
		slog.String("month", date.Month().String()),
		slog.Int("day", date.Day()))

	rows, err := skcDBConn.QueryContext(ctx, productReleasedOnDate, date.Format("01-02"))
	if err != nil {
		return nil, handleQueryError(logger, err)
	}
	defer rows.Close()

	products := make([]*ygo.ProductSummary, 0, 10)
	for rows.Next() {
		var id, locale, name, t, subType, releaseDate string
		var totalItems uint32

		if err := rows.Scan(&id, &locale, &name, &t, &subType, &releaseDate, &totalItems); err != nil {
			return nil, handleRowParsingError(logger, err)
		}

		products = append(products, &ygo.ProductSummary{
			Id:          id,
			Locale:      locale,
			Name:        name,
			Type:        t,
			SubType:     subType,
			ReleaseDate: releaseDate,
			TotalItems:  totalItems})
	}

	if err := rows.Err(); err != nil {
		return nil, handleQueryError(logger, err)
	}

	logger.Info("Succesfully retrieved products", slog.Int("num_products", len(products)))

	return products, nil
}
