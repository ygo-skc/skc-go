package api

import (
	"context"

	"github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/common/ygo"
)

func (s *ygoProductServiceServer) GetCardsByProduct(ctx context.Context, req *ygo.ResourceID) (*ygo.Product, error) {
	_, newCtx := util.NewRequestSetup(ctx, "Product Details")

	p, err := productRepo.GetCardsByProduct(newCtx, req.ID)
	return p, err.Err()
}
