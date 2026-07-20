package client

import (
	context "context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ygo-skc/skc-go/common/v2/model"
	"github.com/ygo-skc/skc-go/common/v2/util"
	"github.com/ygo-skc/skc-go/common/v2/ygo"
	"google.golang.org/grpc/status"
)

type YGOProductClientImp interface {
	GetCardsByProductIDProto(context.Context, string) (*ygo.Product, *model.APIError)

	GetProductSummaryByIDProto(context.Context, string) (*ygo.ProductSummary, *model.APIError)

	GetProductsSummaryByIDProto(context.Context, model.ProductIDs) (*ygo.Products, *model.APIError)
	GetProductsSummaryByID(context.Context, model.ProductIDs) (*model.BatchProductSummaryData[model.ProductIDs], *model.APIError)
}
type YGOProductClientImpV1 struct {
	client ygo.ProductServiceClient
}

func (imp YGOProductClientImpV1) GetCardsByProductIDProto(ctx context.Context, productID string) (*ygo.Product, *model.APIError) {
	return getCardsByProductID(ctx, imp.client, productID)
}

func getCardsByProductID(ctx context.Context, productServiceClient ygo.ProductServiceClient, productID string) (*ygo.Product, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Retrieving cards for product", slog.String("product_id", productID))

	if p, err := productServiceClient.GetCardsByProductID(ctx, &ygo.ResourceID{ID: productID}); err != nil {
		logger.Error("Issue calling YGO Product Service",
			slog.String("operation", "get_cards_by_product"), slog.Any("grpc_code", status.Code(err)), slog.Any("err", err))
		return nil, &model.APIError{Message: fmt.Sprintf("Error fetching cards for product %s", productID), StatusCode: http.StatusInternalServerError}
	} else {
		return p, nil
	}
}

func (imp YGOProductClientImpV1) GetProductSummaryByIDProto(ctx context.Context, productID string) (*ygo.ProductSummary, *model.APIError) {
	return getProductSummaryByID(ctx, imp.client, productID)
}

func getProductSummaryByID(ctx context.Context, productServiceClient ygo.ProductServiceClient, productID string) (*ygo.ProductSummary, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Retrieving product summary", slog.String("product_id", productID))

	if ps, err := productServiceClient.GetProductSummaryByID(ctx, &ygo.ResourceID{ID: productID}); err != nil {
		logger.Error("Issue calling YGO Product Service",
			slog.String("operation", "get_product_summary"), slog.Any("grpc_code", status.Code(err)), slog.Any("err", err))
		return nil, &model.APIError{Message: fmt.Sprintf("Error fetching product summary for product %s", productID), StatusCode: http.StatusInternalServerError}
	} else {
		return ps, nil
	}
}

func (imp YGOProductClientImpV1) GetProductsSummaryByIDProto(ctx context.Context, productID model.ProductIDs) (*ygo.Products, *model.APIError) {
	return getProductsSummaryByID(ctx, imp.client, productID)
}

func (imp YGOProductClientImpV1) GetProductsSummaryByID(ctx context.Context,
	productID model.ProductIDs) (*model.BatchProductSummaryData[model.ProductIDs], *model.APIError) {
	p, err := getProductsSummaryByID(ctx, imp.client, productID)
	if err == nil {
		return model.BatchProductSummaryFromProductsProto(p, model.ProductIDAsKey), nil
	}
	return nil, err
}

func getProductsSummaryByID(ctx context.Context, productServiceClient ygo.ProductServiceClient, productIDs model.ProductIDs) (*ygo.Products, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Retrieving product summaries", slog.Any("product_ids", productIDs))

	if ps, err := productServiceClient.GetProductsSummaryByID(ctx, &ygo.ResourceIDs{IDs: productIDs}); err != nil {
		logger.Error("Issue calling YGO Product Service",
			slog.String("operation", "get_products_summary"), slog.Any("grpc_code", status.Code(err)), slog.Any("err", err))
		return nil, &model.APIError{Message: fmt.Sprintf("Error fetching product summary for product(s) %v", productIDs), StatusCode: http.StatusInternalServerError}
	} else {
		return ps, nil
	}
}
