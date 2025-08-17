package api

import (
	"context"

	"github.com/ygo-skc/skc-go/common/health"
	"github.com/ygo-skc/skc-go/common/util"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *healthServiceServer) APIStatus(ctx context.Context, req *emptypb.Empty) (*health.APIStatusDetails, error) {
	logger, _ := util.NewLogger(context.Background(), "Status")
	logger.Info("Retrieving status of gRPC service")

	return &health.APIStatusDetails{Version: "1.4.0"}, nil
}
