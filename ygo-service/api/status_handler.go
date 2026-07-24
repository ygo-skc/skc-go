package api

import (
	"context"

	"github.com/ygo-skc/skc-go/common/v3/health"
	"github.com/ygo-skc/skc-go/common/v3/util"
)

func (s *healthServiceServer) APIStatus(ctx context.Context, req *health.APIStatusRequest) (*health.APIStatusResponse, error) {
	logger, _ := util.NewLogger(context.Background(), "Status")
	logger.Info("Retrieving status of gRPC service")

	return &health.APIStatusResponse{Version: "3.0.0"}, nil
}
