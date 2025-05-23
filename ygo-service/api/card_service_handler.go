package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/common/ygo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ygoServiceServer) GetCardColors(ctx context.Context, req *emptypb.Empty) (*ygo.CardColors, error) {
	logger, ctx := util.NewRequestSetup(ctx, "Card Colors")
	logger.Info("Retrieving card colors")

	if c, err := skcDBInterface.GetCardColorIDs(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Message)
	} else {
		return c, nil
	}
}

func (s *ygoServiceServer) GetCardByID(ctx context.Context, req *ygo.ResourceID) (*ygo.Card, error) {
	logger, ctx := util.NewRequestSetup(ctx, "Query Card By ID")
	logger.Info(fmt.Sprintf("Getting card details using card ID: %v", req.ID))

	if c, err := skcDBInterface.GetCardByID(ctx, req.ID); err != nil && err.StatusCode == http.StatusNotFound {
		return nil, status.Errorf(codes.NotFound, "%s", err.Message)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Message)
	} else {
		return c, nil
	}
}

func (s *ygoServiceServer) GetCardsByID(ctx context.Context, req *ygo.ResourceIDs) (*ygo.Cards, error) {
	logger, ctx := util.NewRequestSetup(ctx, "Query Cards By ID")
	logger.Info(fmt.Sprintf("Getting card details using card ID's: %v", req.IDs))

	if cards, err := skcDBInterface.GetCardsByIDs(ctx, req.IDs); err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Message)
	} else {
		logger.Info(fmt.Sprintf("The following Card ID's were invalid: %v", cards.UnknownResources))
		return cards, nil
	}
}

func (s *ygoServiceServer) GetCardsByName(ctx context.Context, req *ygo.ResourceNames) (*ygo.Cards, error) {
	logger, ctx := util.NewRequestSetup(ctx, "Query Cards By Name")
	logger.Info(fmt.Sprintf("Getting card details using %d card name(s)", len(req.Names)))

	if cards, err := skcDBInterface.GetCardsByNames(ctx, req.Names); err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Message)
	} else {
		logger.Info(fmt.Sprintf("The following cards were not found: %v", cards.UnknownResources))
		return cards, nil
	}
}

func (s *ygoServiceServer) GetArchetypalCardsUsingCardName(ctx context.Context, req *ygo.Archetype) (*ygo.CardList, error) {
	logger, ctx := util.NewRequestSetup(ctx, "Query Archetypal Cards Using Card Name")

	if cards, err := skcDBInterface.GetArchetypalCardsUsingCardName(ctx, req.Archetype); err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Message)
	} else {
		logger.Info(fmt.Sprintf("Found %d cards in archetype %s", len(cards.Cards), req.Archetype))
		return cards, nil
	}
}

func (s *ygoServiceServer) GetRandomCard(ctx context.Context, req *ygo.BlackListed) (*ygo.Card, error) {
	logger, ctx := util.NewRequestSetup(ctx, "Random Card")
	logger.Info(fmt.Sprintf("Getting random card from DB. Client has provided %d blacklisted IDs", len(req.BlackListedRefs)))

	if c, err := skcDBInterface.GetRandomCard(ctx, req.BlackListedRefs); err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Message)
	} else {
		return c, nil
	}
}
