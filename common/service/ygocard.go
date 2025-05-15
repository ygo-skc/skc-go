package service

import (
	context "context"
	"fmt"
	"net/http"

	"github.com/ygo-skc/skc-go/common/model"
	"github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/common/ygo"
	"google.golang.org/grpc/status"
)

const (
	BATCH_CARD_INFO_OPERATION = "Batch Card Info"
	BATCH_CARD_INFO_ERROR     = "There was an error fetching card info"
)

func QueryCards[T *ygo.Cards | *model.BatchCardData[model.CardIDs]](ctx context.Context,
	cardServiceClient ygo.CardServiceClient, cardIDs []string,
	transformer func(*ygo.Cards) T) (T, *model.APIError) {
	logger := util.LoggerFromContext(ctx)
	logger.Info(fmt.Sprintf("Fetching card info for the following IDs: %v", cardIDs))

	if cards, err := cardServiceClient.QueryCards(ctx, &ygo.Resources{IDs: cardIDs}); err != nil {
		logger.Error(
			fmt.Sprintf("There was an issue calling YGO Service. Operation: %s. Code %s. Error: %s",
				BATCH_CARD_INFO_OPERATION,
				status.Code(err),
				err))
		return nil, &model.APIError{Message: BATCH_CARD_INFO_ERROR, StatusCode: http.StatusInternalServerError}
	} else {
		return transformer(cards), nil
	}
}
