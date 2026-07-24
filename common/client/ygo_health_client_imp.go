package client

import (
	context "context"
	"log/slog"

	"github.com/ygo-skc/skc-go/common/v3/health"
	"github.com/ygo-skc/skc-go/common/v3/model"
	"github.com/ygo-skc/skc-go/common/v3/util"
)

type YGOHealthClientImp interface {
	GetAPIStatus(context.Context) (*health.APIStatusResponse, *model.APIError)
}

type YGOHealthClientImpV1 struct {
	client health.HealthServiceClient
}

func (imp YGOHealthClientImpV1) GetAPIStatus(ctx context.Context) (*health.APIStatusResponse, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	h, err := imp.client.APIStatus(ctx, &health.APIStatusRequest{})
	if err != nil {
		logger.Error("Issue retrieving YGO Service status", slog.Any("err", err))
		return nil, rpcErrorToAPIError(err)
	}
	return h, nil
}
