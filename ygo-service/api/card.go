package api

import (
	"context"
	"net/http"

	"github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/ygo-service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetYGOCard(ctx context.Context, req *pb.YGOCardRequest) (*pb.YGOCardResponse, error) {
	if c, err := skcDBInterface.GetDesiredCardInDBUsingID(ctx, req.ID); err != nil && err.StatusCode == http.StatusNotFound {
		return nil, status.Errorf(codes.NotFound, "%s", err.Message)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Message)
	} else {
		res := &pb.YGOCardResponse{
			ID:          c.ID,
			Color:       c.Color,
			Name:        c.Name,
			Attribute:   c.Attribute,
			Effect:      c.Effect,
			MonsterType: util.PBStringValue(c.MonsterType),
			Attack:      util.PBUInt32Value(c.Attack),
			Defense:     util.PBUInt32Value(c.Defense),
		}

		return res, nil
	}
}
