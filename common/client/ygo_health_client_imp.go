package client

import (
	context "context"
	"log/slog"
	"net/http"

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

	if h, err := imp.client.APIStatus(ctx, &emptypb.Empty{}); err != nil {
		logger.Error("Issue retrieving YGO Service status", slog.Any("err", err))
		return nil, &model.APIError{Message: "API is down", StatusCode: http.StatusInternalServerError}
	} else {
		return h, nil
	}
}
