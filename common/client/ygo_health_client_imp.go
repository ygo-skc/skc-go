package client

import (
	context "context"
	"log/slog"

	"github.com/ygo-skc/skc-go/common/v2/health"
	"github.com/ygo-skc/skc-go/common/v2/model"
	"github.com/ygo-skc/skc-go/common/v2/util"
	"google.golang.org/protobuf/types/known/emptypb"
)

type YGOHealthClientImp interface {
	GetAPIStatus(context.Context) (*health.APIStatusDetails, *model.APIError)
}

type YGOHealthClientImpV1 struct {
	client health.HealthServiceClient
}

func (imp YGOHealthClientImpV1) GetAPIStatus(ctx context.Context) (*health.APIStatusDetails, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	h, err := imp.client.APIStatus(ctx, &emptypb.Empty{})
	if err != nil {
		logger.Error("Issue retrieving YGO Service status", slog.Any("err", err))
		return nil, rpcErrorToAPIError(err)
	}
	return h, nil
}
