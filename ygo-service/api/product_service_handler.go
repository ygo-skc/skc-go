package api

import (
	"context"
	"fmt"

	"github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/common/ygo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ygoProductServiceServer) GetCardsByProduct(ctx context.Context, req *ygo.ResourceID) (*ygo.Product, error) {
	logger, ctx := util.NewRequestSetup(ctx, "Product Details")
	logger.Info(fmt.Sprintf("Retrieving product information using ID %s", req.ID))

	if p, err := productRepo.GetCardsByProduct(ctx, req.ID); err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Message)
	} else {
		logger.Info(fmt.Sprintf("Found %d items for product %s. Name of product: %s", p.TotalItems, p.ID, p.Name))
		return p, nil
	}
}
