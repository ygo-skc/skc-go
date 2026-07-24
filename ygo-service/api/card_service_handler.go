package api

import (
	"context"

	"github.com/ygo-skc/skc-go/common/v3/util"
	"github.com/ygo-skc/skc-go/common/v3/ygo"
)

func (s *ygoCardServiceServer) GetCardColors(ctx context.Context, req *ygo.GetCardColorsRequest) (*ygo.GetCardColorsResponse, error) {
	_, newCtx := util.NewLogger(ctx, "Card Colors")

	c, err := cardRepo.GetCardColorIDs(newCtx)
	return &ygo.GetCardColorsResponse{Values: c}, err.Err()
}

func (s *ygoCardServiceServer) GetCardByID(ctx context.Context, req *ygo.GetCardByIDRequest) (*ygo.GetCardByIDResponse, error) {
	_, newCtx := util.NewLogger(ctx, "Query Card By ID")

	c, err := cardRepo.GetCardByID(newCtx, req.Subject.Id)
	return &ygo.GetCardByIDResponse{Card: c}, err.Err()
}

func (s *ygoCardServiceServer) GetCardsByID(ctx context.Context, req *ygo.GetCardsByIDRequest) (*ygo.GetCardsByIDResponse, error) {
	_, newCtx := util.NewLogger(ctx, "Query Cards By ID")

	c, err := cardRepo.GetCardsByIDs(newCtx, req.Subjects.Ids)
	return &ygo.GetCardsByIDResponse{Cards: c}, err.Err()
}

func (s *ygoCardServiceServer) GetCardsByName(ctx context.Context, req *ygo.GetCardsByNameRequest) (*ygo.GetCardsByNameResponse, error) {
	_, newCtx := util.NewLogger(ctx, "Query Cards By Name")

	c, err := cardRepo.GetCardsByNames(newCtx, req.Subjects.Names)
	return &ygo.GetCardsByNameResponse{Cards: c}, err.Err()
}

func (s *ygoCardServiceServer) GetCardsReferencingNameInEffect(ctx context.Context, req *ygo.GetCardsReferencingNameInEffectRequest) (*ygo.GetCardsReferencingNameInEffectResponse, error) {
	_, newCtx := util.NewLogger(ctx, "Find Refs Using Card Effect")

	c, err := cardRepo.GetCardsReferencingNameInEffect(newCtx, req.Subjects.Names)
	return &ygo.GetCardsReferencingNameInEffectResponse{Cards: c}, err.Err()
}

func (s *ygoCardServiceServer) GetArchetypalCardsUsingCardName(ctx context.Context, req *ygo.GetArchetypalCardsUsingCardNameRequest) (*ygo.GetArchetypalCardsUsingCardNameResponse, error) {
	_, newCtx := util.NewLogger(ctx, "Query Archetypal Cards Using Card Name")

	c, err := cardRepo.GetArchetypalCardsUsingCardName(newCtx, req.Subject.Name)
	return &ygo.GetArchetypalCardsUsingCardNameResponse{Cards: c}, err.Err()
}

func (s *ygoCardServiceServer) GetExplicitArchetypalInclusions(ctx context.Context, req *ygo.GetExplicitArchetypalInclusionsRequest) (*ygo.GetExplicitArchetypalInclusionsResponse, error) {
	_, newCtx := util.NewLogger(ctx, "Query Archetypal Inclusions")

	c, err := cardRepo.GetExplicitArchetypalInclusions(newCtx, req.Subject.Name)
	return &ygo.GetExplicitArchetypalInclusionsResponse{Cards: c}, err.Err()
}

func (s *ygoCardServiceServer) GetExplicitArchetypalExclusions(ctx context.Context, req *ygo.GetExplicitArchetypalExclusionsRequest) (*ygo.GetExplicitArchetypalExclusionsResponse, error) {
	_, newCtx := util.NewLogger(ctx, "Query Archetypal Exclusions")

	c, err := cardRepo.GetExplicitArchetypalExclusions(newCtx, req.Subject.Name)
	return &ygo.GetExplicitArchetypalExclusionsResponse{Cards: c}, err.Err()
}

func (s *ygoCardServiceServer) GetRandomCard(ctx context.Context, req *ygo.GetRandomCardRequest) (*ygo.GetRandomCardResponse, error) {
	_, newCtx := util.NewLogger(ctx, "Random Card")

	c, err := cardRepo.GetRandomCard(newCtx, req.Blacklist.BlackListedRefs)
	return &ygo.GetRandomCardResponse{Card: c}, err.Err()
}
