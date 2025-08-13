package api

import (
	"context"

	"github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/common/ygo"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ygoCardServiceServer) GetCardColors(ctx context.Context, req *emptypb.Empty) (*ygo.CardColors, error) {
	_, newCtx := util.NewLogger(ctx, "Card Colors")

	c, err := cardRepo.GetCardColorIDs(newCtx)
	return c, err.Err()
}

func (s *ygoCardServiceServer) GetCardByID(ctx context.Context, req *ygo.ResourceID) (*ygo.Card, error) {
	_, newCtx := util.NewLogger(ctx, "Query Card By ID")

	c, err := cardRepo.GetCardByID(newCtx, req.ID)
	return c, err.Err()
}

func (s *ygoCardServiceServer) GetCardsByID(ctx context.Context, req *ygo.ResourceIDs) (*ygo.Cards, error) {
	_, newCtx := util.NewLogger(ctx, "Query Cards By ID")

	c, err := cardRepo.GetCardsByIDs(newCtx, req.IDs)
	return c, err.Err()
}

func (s *ygoCardServiceServer) GetCardsByName(ctx context.Context, req *ygo.ResourceNames) (*ygo.Cards, error) {
	_, newCtx := util.NewLogger(ctx, "Query Cards By Name")

	c, err := cardRepo.GetCardsByNames(newCtx, req.Names)
	return c, err.Err()
}

func (s *ygoCardServiceServer) GetCardsReferencingNameInEffect(ctx context.Context, req *ygo.ResourceNames) (*ygo.CardList, error) {
	_, newCtx := util.NewLogger(ctx, "Find Refs Using Card Effect")

	c, err := cardRepo.GetCardsReferencingNameInEffect(newCtx, req.Names)
	return c, err.Err()
}

func (s *ygoCardServiceServer) GetArchetypalCardsUsingCardName(ctx context.Context, req *ygo.Archetype) (*ygo.CardList, error) {
	_, newCtx := util.NewLogger(ctx, "Query Archetypal Cards Using Card Name")

	c, err := cardRepo.GetArchetypalCardsUsingCardName(newCtx, req.Archetype)
	return c, err.Err()
}

func (s *ygoCardServiceServer) GetExplicitArchetypalInclusions(ctx context.Context, req *ygo.Archetype) (*ygo.CardList, error) {
	_, newCtx := util.NewLogger(ctx, "Query Archetypal Inclusions")

	c, err := cardRepo.GetExplicitArchetypalInclusions(newCtx, req.Archetype)
	return c, err.Err()
}

func (s *ygoCardServiceServer) GetExplicitArchetypalExclusions(ctx context.Context, req *ygo.Archetype) (*ygo.CardList, error) {
	_, newCtx := util.NewLogger(ctx, "Query Archetypal Exclusions")

	c, err := cardRepo.GetExplicitArchetypalExclusions(newCtx, req.Archetype)
	return c, err.Err()
}

func (s *ygoCardServiceServer) GetRandomCard(ctx context.Context, req *ygo.BlackListed) (*ygo.Card, error) {
	_, newCtx := util.NewLogger(ctx, "Random Card")

	c, err := cardRepo.GetRandomCard(newCtx, req.BlackListedRefs)
	return c, err.Err()
}
