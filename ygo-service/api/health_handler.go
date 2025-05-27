package api

import (
	"context"

	"github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/ygo-service/health"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *healthServiceServer) APIStatus(ctx context.Context, req *emptypb.Empty) (*health.APIStatusDetails, error) {
	logger, _ := util.NewRequestSetup(context.Background(), "Status")
	logger.Info("Retrieving status of gRPC service")

	return &health.APIStatusDetails{Version: "1.2.0"}, nil
}
