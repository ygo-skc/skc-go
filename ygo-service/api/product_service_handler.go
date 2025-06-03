package api

import (
	"context"

	"github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/common/ygo"
)

func (s *ygoProductServiceServer) GetCardsByProductID(ctx context.Context, req *ygo.ResourceID) (*ygo.Product, error) {
	_, newCtx := util.NewRequestSetup(ctx, "Product Details")

	p, err := productRepo.GetCardsByProductID(newCtx, req.ID)
	return p, err.Err()
}

func (s *ygoProductServiceServer) GetProductSummaryByID(ctx context.Context, req *ygo.ResourceID) (*ygo.ProductSummary, error) {
	_, newCtx := util.NewRequestSetup(ctx, "Product Summary")

	p, err := productRepo.GetProductSummaryByID(newCtx, req.ID)
	return p, err.Err()
}

func (s *ygoProductServiceServer) GetProductsSummaryByID(ctx context.Context, req *ygo.ResourceIDs) (*ygo.Products, error) {
	_, newCtx := util.NewRequestSetup(ctx, "Products Summary")

	products, err := productRepo.GetProductsSummaryByID(newCtx, req.IDs)
	return products, err.Err()
}
