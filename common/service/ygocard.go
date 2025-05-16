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

func QueryCard[T *ygo.Card | *model.YGOCardREST | *model.YGOCardGRPC](ctx context.Context,
	cardServiceClient ygo.CardServiceClient, cardID string, transformer func(*ygo.Card) T) (T, *model.APIError) {
	logger := util.LoggerFromContext(ctx)
	logger.Info(fmt.Sprintf("Fetching card info using ID: %v", cardID))

	if cards, err := cardServiceClient.QueryCard(ctx, &ygo.Resource{ID: cardID}); err != nil {
		logger.Error(
			fmt.Sprintf("There was an issue calling YGO Service. Operation: %s. Code %s. Error: %s",
				"Query Card", status.Code(err), err))
		return nil, &model.APIError{Message: "There was an error fetching card info", StatusCode: http.StatusInternalServerError}
	} else {
		return transformer(cards), nil
	}
}

func QueryCards[T *ygo.Cards | *model.BatchCardData[model.CardIDs]](ctx context.Context,
	cardServiceClient ygo.CardServiceClient, cardIDs []string,
	transformer func(*ygo.Cards) T) (T, *model.APIError) {
	logger := util.LoggerFromContext(ctx)
	logger.Info(fmt.Sprintf("Fetching card info for the following IDs: %v", cardIDs))

	if cards, err := cardServiceClient.QueryCards(ctx, &ygo.Resources{IDs: cardIDs}); err != nil {
		logger.Error(
			fmt.Sprintf("There was an issue calling YGO Service. Operation: %s. Code %s. Error: %s",
				"Query Cards", status.Code(err), err))
		return nil, &model.APIError{Message: "There was an error fetching batch card info", StatusCode: http.StatusInternalServerError}
	} else {
		if cards.UnknownResources == nil {
			cards.UnknownResources = make([]string, 0)
		}
		return transformer(cards), nil
	}
}
