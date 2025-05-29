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
}
type YGOProductClientImpV1 struct {
	client *ygo.ProductServiceClient
}

func (imp YGOProductClientImpV1) GetCardsByProductIDProto(ctx context.Context, productID string) (*ygo.Product, *model.APIError) {
	return getCardsByProductID(ctx, imp.client, productID)
}

func getCardsByProductID(ctx context.Context, productServiceClient *ygo.ProductServiceClient, productID string) (*ygo.Product, *model.APIError) {
	logger := util.LoggerFromContext(ctx)
	logger.Info("Retrieving card colors")

	if cColors, err := (*productServiceClient).GetCardsByProductID(ctx, &ygo.ResourceID{ID: productID}); err != nil {
		logger.Error(
			fmt.Sprintf("There was an issue calling YGO Service. Operation: %s. Code %s. Error: %s",
				"Get Card Colors", status.Code(err), err))
		return nil, &model.APIError{Message: "Error fetching card color data", StatusCode: http.StatusInternalServerError}
	} else {
		return cColors, nil
	}
}
