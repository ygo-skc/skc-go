package db

import (
	"context"

	"github.com/ygo-skc/skc-go/common/model"
	cUtil "github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/common/ygo"
)

type YGOProductRepository struct{}

func (imp YGOProductRepository) GetCardsByProduct(ctx context.Context, productID string) (*ygo.Product, *model.APIError) {
	logger := cUtil.LoggerFromContext(ctx)

	if product, err := queryProductInfo(logger, productID); err != nil {
		return nil, err
	} else {
		if rows, err := skcDBConn.Query(cardsByProductIDQuery, productID); err != nil {
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
