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
	logger.Info(fmt.Sprintf("Fetching card details using cardID %v", req.ID))

	if c, err := skcDBInterface.GetCardByID(ctx, req.ID); err != nil && err.StatusCode == http.StatusNotFound {
		logger.Info(fmt.Sprintf("%s Not found in DB", req.ID))
		return nil, status.Errorf(codes.NotFound, "%s", err.Message)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Message)
	} else {

		return c.ToPB(), nil
	}
}

func (s *ygoServiceServer) QueryCards(ctx context.Context, req *ygo.Resources) (*ygo.Cards, error) {
	logger, ctx := util.NewRequestSetup(ctx, "Query Cards")
	logger.Info(fmt.Sprintf("Fetching card details using cardIDs: %v", req.IDs))

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

		logger.Info(fmt.Sprintf("Valid card IDs: %v. Invalid IDs: %v", validIDs, cards.UnknownResources))
		return &ygo.Cards{CardInfo: pbCards, UnknownResources: cards.UnknownResources}, nil
	}
}
