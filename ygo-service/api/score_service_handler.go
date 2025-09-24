package api

import (
	"context"

	"github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/common/ygo"
)

func (s *ygoScoreServiceServer) GetCardScoreByID(ctx context.Context, req *ygo.ResourceID) (*ygo.CardScore, error) {
	_, newCtx := util.NewLogger(ctx, "Card Score")

	if scoreHistory, err := scoreRepo.GetCardScoreByID(newCtx, req.ID); err != nil {
		return nil, err.Err()
	} else {
		return &ygo.CardScore{ScoreHistory: scoreHistory}, nil
	}
}

func (s *ygoScoreServiceServer) GetCardScoresByIDs(ctx context.Context, req *ygo.ResourceIDs) (*ygo.CardScores, error) {
	// _, newCtx := util.NewLogger(ctx, "Card Score")

	return nil, nil
}
