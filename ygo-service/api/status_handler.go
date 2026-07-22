package api

import (
	"context"

	"github.com/ygo-skc/skc-go/common/v3/health"
	"github.com/ygo-skc/skc-go/common/v3/util"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *healthServiceServer) APIStatus(ctx context.Context, req *emptypb.Empty) (*health.APIStatusDetails, error) {
	logger, _ := util.NewLogger(context.Background(), "Status")
	logger.Info("Retrieving status of gRPC service")

	return &health.APIStatusDetails{Version: "2.1.10"}, nil
}
