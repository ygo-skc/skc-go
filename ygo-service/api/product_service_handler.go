package api

import (
	"context"

	"github.com/ygo-skc/skc-go/common/v3/util"
	"github.com/ygo-skc/skc-go/common/v3/ygo"
)

func (s *ygoProductServiceServer) GetCardsByProductID(ctx context.Context, req *ygo.GetCardsByProductIDRequest) (*ygo.GetCardsByProductIDResponse, error) {
	_, newCtx := util.NewLogger(ctx, "Product Details")

	p, err := productRepo.GetCardsByProductID(newCtx, req.Subject.Id)
	return &ygo.GetCardsByProductIDResponse{Product: p}, err.Err()
}

func (s *ygoProductServiceServer) GetProductSummaryByID(ctx context.Context, req *ygo.GetProductSummaryByIDRequest) (*ygo.GetProductSummaryByIDResponse, error) {
	_, newCtx := util.NewLogger(ctx, "Product Summary")

	p, err := productRepo.GetProductSummaryByID(newCtx, req.Subject.Id)
	return &ygo.GetProductSummaryByIDResponse{ProductSummary: p}, err.Err()
}

func (s *ygoProductServiceServer) GetProductsSummaryByID(ctx context.Context, req *ygo.GetProductsSummaryByIDRequest) (*ygo.GetProductsSummaryByIDResponse, error) {
	_, newCtx := util.NewLogger(ctx, "Products Summary")

	products, err := productRepo.GetProductsSummaryByID(newCtx, req.Subjects.Ids)
	return &ygo.GetProductsSummaryByIDResponse{Products: products}, err.Err()
}
