package client

import (
	context "context"
	"fmt"
	"net/http"

	"github.com/ygo-skc/skc-go/common/model"
	"github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/common/ygo"
	"google.golang.org/grpc/status"
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
	return getCardsByProductID(ctx, imp.client, productID)
}

func getCardsByProductID(ctx context.Context, productServiceClient ygo.ProductServiceClient, productID string) (*ygo.Product, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info(fmt.Sprintf("Retrieving cards for product w/ ID %s", productID))

	if p, err := productServiceClient.GetCardsByProductID(ctx, &ygo.ResourceID{ID: productID}); err != nil {
		logger.Error(
			fmt.Sprintf("There was an issue calling YGO Service. Operation: %s. Code %s. Error: %s",
				"Get Cards By Product", status.Code(err), err))
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
	logger.Info(fmt.Sprintf("Retrieving summary of product w/ ID %s", productID))

	if ps, err := productServiceClient.GetProductSummaryByID(ctx, &ygo.ResourceID{ID: productID}); err != nil {
		logger.Error(
			fmt.Sprintf("There was an issue calling YGO Service. Operation: %s. Code %s. Error: %s",
				"Get Product Summary", status.Code(err), err))
		return nil, &model.APIError{Message: fmt.Sprintf("Error fetching product summary for product %s", productID), StatusCode: http.StatusInternalServerError}
	} else {
		return ps, nil
	}
}

func (imp YGOProductClientImpV1) GetProductsSummaryByIDProto(ctx context.Context, productID model.ProductIDs) (*ygo.Products, *model.APIError) {
	return getProductsSummaryByID(ctx, imp.client, productID)
}

func getProductsSummaryByID(ctx context.Context, productServiceClient ygo.ProductServiceClient, productIDs model.ProductIDs) (*ygo.Products, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info(fmt.Sprintf("Retrieving summary of product w/ ID %s", productIDs))

	if ps, err := productServiceClient.GetProductsSummaryByID(ctx, &ygo.ResourceIDs{IDs: productIDs}); err != nil {
		logger.Error(
			fmt.Sprintf("There was an issue calling YGO Service. Operation: %s. Code %s. Error: %s",
				"Get Products Summary", status.Code(err), err))
		return nil, &model.APIError{Message: fmt.Sprintf("Error fetching product summary for product(s) %v", productIDs), StatusCode: http.StatusInternalServerError}
	} else {
		return ps, nil
	}
}
