package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ygo-skc/skc-go/common/model"
	"github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/ygo-service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) QueryCard(ctx context.Context, req *pb.YGOResource) (*pb.YGOCard, error) {
	logger, ctx := util.NewRequestSetup(ctx, "Query Card")
	logger.Info(fmt.Sprintf("Fetching card details using %v", req))

	if c, err := skcDBInterface.GetDesiredCardInDBUsingID(ctx, req.ID); err != nil && err.StatusCode == http.StatusNotFound {
		return nil, status.Errorf(codes.NotFound, "%s", err.Message)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Message)
	} else {

		return cardToPB(c), nil
	}
}

func (s *Server) QueryCards(ctx context.Context, req *pb.YGOResources) (*pb.YGOCards, error) {
	logger, ctx := util.NewRequestSetup(ctx, "Query Card")
	logger.Info(fmt.Sprintf("Fetching card details using %v", req))

	if cards, err := skcDBInterface.GetDesiredCardInDBUsingMultipleCardIDs(ctx, req.ID); err != nil && err.StatusCode == http.StatusNotFound {
		return nil, status.Errorf(codes.NotFound, "%s", err.Message)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Message)
	} else {
		pbCards := make([]*pb.YGOCard, len(cards.CardInfo))
		i := 0
		for _, c := range cards.CardInfo {
			pbCards[i] = cardToPB(c)
			i++
		}
		return &pb.YGOCards{Cards: pbCards, UnknownResources: cards.UnknownResources}, nil
	}
}

func cardToPB(c model.Card) *pb.YGOCard {
	return &pb.YGOCard{
		ID:          c.ID,
		Color:       c.Color,
		Name:        c.Name,
		Attribute:   c.Attribute,
		Effect:      c.Effect,
		MonsterType: util.PBStringValue(c.MonsterType),
		Attack:      util.PBUInt32Value(c.Attack),
		Defense:     util.PBUInt32Value(c.Defense),
	}
}
