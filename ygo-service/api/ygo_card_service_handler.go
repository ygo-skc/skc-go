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
	logger.Info(fmt.Sprintf("Fetching card details using %v", req))

	if c, err := skcDBInterface.GetDesiredCardInDBUsingID(ctx, req.ID); err != nil && err.StatusCode == http.StatusNotFound {
		return nil, status.Errorf(codes.NotFound, "%s", err.Message)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Message)
	} else {

		return c.ToPB(), nil
	}
}

func (s *ygoServiceServer) QueryCards(ctx context.Context, req *ygo.Resources) (*ygo.Cards, error) {
	logger, ctx := util.NewRequestSetup(ctx, "Query Card")
	logger.Info(fmt.Sprintf("Fetching card details using %v", req))

	if cards, err := skcDBInterface.GetDesiredCardInDBUsingMultipleCardIDs(ctx, req.IDs); err != nil && err.StatusCode == http.StatusNotFound {
		return nil, status.Errorf(codes.NotFound, "%s", err.Message)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Message)
	} else {
		pbCards := make(map[string]*ygo.Card, len(cards.CardInfo))
		for _, c := range cards.CardInfo {
			pbCards[c.GetID()] = c.(model.Card).ToPB()
		}
		return &ygo.Cards{CardInfo: pbCards, UnknownResources: cards.UnknownResources}, nil
	}
}
