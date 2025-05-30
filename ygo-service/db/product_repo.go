package db

import (
	"context"
	"fmt"

	cUtil "github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/common/ygo"
	"google.golang.org/grpc/status"
)

type ProductRepository interface {
	GetCardsByProductID(context.Context, string) (*ygo.Product, *status.Status)
}
type YGOProductRepository struct{}

func (imp YGOProductRepository) GetCardsByProductID(ctx context.Context, productID string) (*ygo.Product, *status.Status) {
	logger := cUtil.LoggerFromContext(ctx)
	logger.Info(fmt.Sprintf("Retrieving product data using ID %s", productID))

	if product, err := queryProductInfo(logger, productID); err != nil {
		return nil, err
	} else {
		query := fmt.Sprintf(cardsByProductIDQuery, cardAttributes)
		if rows, err := skcDBConn.Query(query, productID); err != nil {
			return nil, handleQueryError(cUtil.LoggerFromContext(ctx), err)
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
