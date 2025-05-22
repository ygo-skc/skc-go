package service

import (
	context "context"
	"fmt"
	"net/http"

	"github.com/ygo-skc/skc-go/common/model"
	"github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/common/ygo"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type YGOService interface {
	CardColors(context.Context) (*ygo.CardColors, *model.APIError)

	QueryCard(context.Context, string) (*ygo.Card, *model.APIError)
	QueryCardREST(context.Context, string) (*model.YGOCardREST, *model.APIError)

	QueryCards(context.Context, []string) (*ygo.Cards, *model.APIError)
	QueryCardsREST(context.Context, []string) (*model.BatchCardData[model.CardIDs], *model.APIError)

	RandomCard(context.Context, []string) (*ygo.Card, *model.APIError)
	RandomCardREST(context.Context, []string) (*model.YGOCardREST, *model.APIError)
}

type YGOServiceV1 struct {
	client           ygo.CardServiceClient
	cardTransformer  model.YGOCardTransformer
	cardsTransformer model.YGOCardsTransformer
}

func NewYGOServiceV1(client ygo.CardServiceClient) YGOServiceV1 {
	return YGOServiceV1{
		client:           client,
		cardTransformer:  model.YGOCardTransformerV1{},
		cardsTransformer: model.YGOCardsTransformerV1{},
	}
}

func (svc YGOServiceV1) CardColors(ctx context.Context) (*ygo.CardColors, *model.APIError) {
	logger := util.LoggerFromContext(ctx)
	logger.Info("Retrieving card colors")

	if cColors, err := svc.client.Colors(ctx, &emptypb.Empty{}); err != nil {
		logger.Error(
			fmt.Sprintf("There was an issue calling YGO Service. Operation: %s. Code %s. Error: %s",
				"Query Card", status.Code(err), err))
		return nil, &model.APIError{Message: "There was an error fetching card info", StatusCode: http.StatusInternalServerError}
	} else {
		return cColors, nil
	}
}

func (svc YGOServiceV1) QueryCard(ctx context.Context, cardID string) (*ygo.Card, *model.APIError) {
	return queryCard(ctx, svc.client, cardID, svc.cardTransformer.ToSelf)
}

func (svc YGOServiceV1) QueryCardREST(ctx context.Context, cardID string) (*model.YGOCardREST, *model.APIError) {
	return queryCard(ctx, svc.client, cardID, svc.cardTransformer.ToREST)
}

func queryCard[T *ygo.Card | *model.YGOCardREST | *model.YGOCardGRPC](ctx context.Context,
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

func (svc YGOServiceV1) QueryCards(ctx context.Context, cardIDs []string) (*ygo.Cards, *model.APIError) {
	return queryCards(ctx, svc.client, cardIDs, svc.cardsTransformer.ToSelf)
}

func (svc YGOServiceV1) QueryCardsREST(ctx context.Context, cardIDs []string) (*model.BatchCardData[model.CardIDs], *model.APIError) {
	return queryCards(ctx, svc.client, cardIDs, svc.cardsTransformer.ToREST)
}

func queryCards[T *ygo.Cards | *model.BatchCardData[model.CardIDs]](ctx context.Context,
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

func (svc YGOServiceV1) RandomCard(ctx context.Context, blackListedIDs []string) (*ygo.Card, *model.APIError) {
	return randomCard(ctx, svc.client, blackListedIDs, svc.cardTransformer.ToSelf)
}

func (svc YGOServiceV1) RandomCardREST(ctx context.Context, blackListedIDs []string) (*model.YGOCardREST, *model.APIError) {
	return randomCard(ctx, svc.client, blackListedIDs, svc.cardTransformer.ToREST)
}

func randomCard[T *ygo.Card | *model.YGOCardREST | *model.YGOCardGRPC](ctx context.Context,
	cardServiceClient ygo.CardServiceClient, blackListedIDs []string, transformer func(*ygo.Card) T) (T, *model.APIError) {
	logger := util.LoggerFromContext(ctx)
	logger.Info("Getting random card")

	if cards, err := cardServiceClient.RandomCard(ctx, &ygo.BlackListedResources{BlackListedRefs: blackListedIDs}); err != nil {
		logger.Error(
			fmt.Sprintf("There was an issue calling YGO Service. Operation: %s. Code %s. Error: %s",
				"Random Card", status.Code(err), err))
		return nil, &model.APIError{Message: "There was an error fetching random card", StatusCode: http.StatusInternalServerError}
	} else {
		return transformer(cards), nil
	}
}
