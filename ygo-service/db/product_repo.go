package db

import (
	"context"
	"fmt"

	"github.com/ygo-skc/skc-go/common/model"
	cUtil "github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/common/ygo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
