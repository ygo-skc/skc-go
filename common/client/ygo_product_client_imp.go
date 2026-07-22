package client

import (
	context "context"
	"log/slog"

	"github.com/ygo-skc/skc-go/common/v3/model"
	"github.com/ygo-skc/skc-go/common/v3/util"
	"github.com/ygo-skc/skc-go/common/v3/ygo"
)

type YGOProductClientImp interface {
	GetCardsByProductIDProto(context.Context, string) (*ygo.Product, *model.APIError)
	GetProductSummaryByIDProto(context.Context, string) (*ygo.ProductSummary, *model.APIError)
	GetProductsSummaryByIDProto(context.Context, model.ProductIDs) (*ygo.Products, *model.APIError)
}
type YGOProductClientImpV1 struct {
	client ygo.ProductServiceClient
}

func (imp YGOProductClientImpV1) GetCardsByProductIDProto(ctx context.Context, productID string) (*ygo.Product, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Retrieving cards associated to product", slog.String("ygo_service.resource", productID))
	p, err := imp.client.GetCardsByProductID(ctx, &ygo.ResourceID{ID: productID})
	if err != nil {
		logger.Error("Issue calling YGO Product Service", slog.Any("err", err))
		return nil, rpcErrorToAPIError(err)
	}
	return p, nil
}

func (imp YGOProductClientImpV1) GetProductSummaryByIDProto(ctx context.Context, productID string) (*ygo.ProductSummary, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Retrieving product summary", slog.String("ygo_service.resource", productID))
	ps, err := imp.client.GetProductSummaryByID(ctx, &ygo.ResourceID{ID: productID})
	if err != nil {
		logger.Error("Issue calling YGO Product Service", slog.Any("err", err))
		return nil, rpcErrorToAPIError(err)
	}
	return ps, nil
}

func (imp YGOProductClientImpV1) GetProductsSummaryByIDProto(ctx context.Context, productIDs model.ProductIDs) (*ygo.Products, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Retrieving product summaries", slog.Any("product_ids", productIDs))
	ps, err := imp.client.GetProductsSummaryByID(ctx, &ygo.ResourceIDs{IDs: productIDs})
	if err != nil {
		logger.Error("Issue calling YGO Product Service", slog.Any("err", err))
		return nil, rpcErrorToAPIError(err)
	}
	return ps, nil
}
