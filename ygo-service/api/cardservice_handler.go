package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ygo-skc/skc-go/common/model"
	"github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/common/ygo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ygoServiceServer) QueryCard(ctx context.Context, req *ygo.Resource) (*ygo.Card, error) {
	logger, ctx := util.NewRequestSetup(ctx, "Query Card")
	logger.Info(fmt.Sprintf("Fetching card details using card ID: %v", req.ID))

	if c, err := skcDBInterface.GetCardByID(ctx, req.ID); err != nil && err.StatusCode == http.StatusNotFound {
		return nil, status.Errorf(codes.NotFound, "%s", err.Message)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Message)
	} else {

		return c.ToPB(), nil
	}
}

func (s *ygoServiceServer) QueryCards(ctx context.Context, req *ygo.Resources) (*ygo.Cards, error) {
	logger, ctx := util.NewRequestSetup(ctx, "Query Cards")
	logger.Info(fmt.Sprintf("Fetching card details using card ID's: %v", req.IDs))

	if cards, err := skcDBInterface.GetCardsByIDs(ctx, req.IDs); err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Message)
	} else {
		pbCards := make(map[string]*ygo.Card, len(cards.CardInfo))
		validIDs := make(model.CardIDs, len(cards.CardInfo))
		ind := 0
		for id, c := range cards.CardInfo {
			pbCards[id] = c.(model.YGOCardREST).ToPB()
			validIDs[ind] = id
			ind++
		}

		logger.Info(fmt.Sprintf("The following Card ID's were invalid: %v", cards.UnknownResources))
		return &ygo.Cards{CardInfo: pbCards, UnknownResources: cards.UnknownResources}, nil
	}
}

func (s *ygoServiceServer) RandomCard(ctx context.Context, req *ygo.BlackListedResources) (*ygo.Card, error) {
	logger, ctx := util.NewRequestSetup(ctx, "Random Card")
	logger.Info(fmt.Sprintf("Getting random card from DB. Client has provided %d blacklisted IDs", len(req.BlackListedRefs)))

	if c, err := skcDBInterface.GetRandomCard(ctx, req.BlackListedRefs); err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Message)
	} else {
		return c.ToPB(), nil
	}
}
