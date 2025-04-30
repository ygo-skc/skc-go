package api

import (
	"context"
	"net/http"

	"github.com/ygo-skc/skc-go/skc-db-service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *Server) GetYGOCard(ctx context.Context, req *pb.YGOCardRequest) (*pb.YGOCardResponse, error) {
	if c, err := skcDBInterface.GetDesiredCardInDBUsingID(ctx, req.ID); err != nil && err.StatusCode == http.StatusNotFound {
		return nil, status.Errorf(codes.NotFound, "%s", err.Message)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Message)
	} else {
		res := &pb.YGOCardResponse{
			ID:        c.ID,
			Color:     c.Color,
			Name:      c.Name,
			Attribute: c.Attribute,
			Effect:    c.Effect,
		}

		if c.MonsterType != nil {
			res.MonsterType = wrapperspb.String(*c.MonsterType)
		}

		if c.Attack != nil {
			res.Attack = wrapperspb.UInt32(*c.Attack)
		}

		if c.Defense != nil {
			res.Defense = wrapperspb.UInt32(*c.Defense)
		}

		return res, nil
	}
}
