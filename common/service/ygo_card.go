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
	GetCardColors(context.Context) (*ygo.CardColors, *model.APIError)

	GetCardByID(context.Context, string) (*ygo.Card, *model.APIError)
	GetCardByIDREST(context.Context, string) (*model.YGOCardREST, *model.APIError)

	GetCardsByID(context.Context, model.CardIDs) (*ygo.Cards, *model.APIError)
	GetCardsByIDREST(context.Context, model.CardIDs) (*model.BatchCardData[model.CardIDs], *model.APIError)

	GetCardsByName(context.Context, model.CardNames) (*ygo.Cards, *model.APIError)
	GetCardsByNameREST(context.Context, model.CardNames) (*model.BatchCardData[model.CardNames], *model.APIError)

	GetRandomCard(context.Context, []string) (*ygo.Card, *model.APIError)
	GetRandomCardREST(context.Context, []string) (*model.YGOCardREST, *model.APIError)
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

func (svc YGOServiceV1) GetCardColors(ctx context.Context) (*ygo.CardColors, *model.APIError) {
	logger := util.LoggerFromContext(ctx)
	logger.Info("Retrieving card colors")

	if cColors, err := svc.client.GetCardColors(ctx, &emptypb.Empty{}); err != nil {
		logger.Error(
			fmt.Sprintf("There was an issue calling YGO Service. Operation: %s. Code %s. Error: %s",
				"Get Card Colors", status.Code(err), err))
		return nil, &model.APIError{Message: "There was an error fetching card info", StatusCode: http.StatusInternalServerError}
	} else {
		return cColors, nil
	}
}

func (svc YGOServiceV1) GetCardByID(ctx context.Context, cardID string) (*ygo.Card, *model.APIError) {
	return getCardByID(ctx, svc.client, cardID, svc.cardTransformer.ToSelf)
}

func (svc YGOServiceV1) GetCardByIDREST(ctx context.Context, cardID string) (*model.YGOCardREST, *model.APIError) {
	return getCardByID(ctx, svc.client, cardID, svc.cardTransformer.ToREST)
}

func getCardByID[T *ygo.Card | *model.YGOCardREST | *model.YGOCardGRPC](ctx context.Context,
	cardServiceClient ygo.CardServiceClient, cardID string, transformer func(*ygo.Card) T) (T, *model.APIError) {
	logger := util.LoggerFromContext(ctx)
	logger.Info(fmt.Sprintf("Fetching card info using ID: %v", cardID))

	if cards, err := cardServiceClient.GetCardByID(ctx, &ygo.ResourceID{ID: cardID}); err != nil {
		logger.Error(
			fmt.Sprintf("There was an issue calling YGO Service. Operation: %s. Code %s. Error: %s",
				"Get Card By ID", status.Code(err), err))
		return nil, &model.APIError{Message: "There was an error fetching card info", StatusCode: http.StatusInternalServerError}
	} else {
		return transformer(cards), nil
	}
}

func (svc YGOServiceV1) GetCardsByID(ctx context.Context, cardIDs model.CardIDs) (*ygo.Cards, *model.APIError) {
	return getCardsByID(ctx, svc.client, cardIDs, svc.cardsTransformer.ToSelf)
}

func (svc YGOServiceV1) GetCardsByIDREST(ctx context.Context, cardIDs model.CardIDs) (*model.BatchCardData[model.CardIDs], *model.APIError) {
	return getCardsByID(ctx, svc.client, cardIDs, svc.cardsTransformer.ToBatchCardDataUsingID)
}

func getCardsByID[T *ygo.Cards | *model.BatchCardData[model.CardIDs]](ctx context.Context,
	cardServiceClient ygo.CardServiceClient, cardIDs model.CardIDs,
	transformer func(*ygo.Cards) T) (T, *model.APIError) {
	logger := util.LoggerFromContext(ctx)
	logger.Info(fmt.Sprintf("Fetching card info for the following IDs: %v", cardIDs))

	if cards, err := cardServiceClient.GetCardsByID(ctx, &ygo.ResourceIDs{IDs: cardIDs}); err != nil {
		logger.Error(
			fmt.Sprintf("There was an issue calling YGO Service. Operation: %s. Code %s. Error: %s",
				"Get Cards By ID", status.Code(err), err))
		return nil, &model.APIError{Message: "There was an error fetching batch card info", StatusCode: http.StatusInternalServerError}
	} else {
		if cards.UnknownResources == nil {
			cards.UnknownResources = make([]string, 0)
		}
		return transformer(cards), nil
	}
}

func (svc YGOServiceV1) GetCardsByName(ctx context.Context, cardNames model.CardNames) (*ygo.Cards, *model.APIError) {
	return getCardsByName(ctx, svc.client, cardNames, svc.cardsTransformer.ToSelf)
}

func (svc YGOServiceV1) GetCardsByNameREST(ctx context.Context, cardNames model.CardNames) (*model.BatchCardData[model.CardNames], *model.APIError) {
	return getCardsByName(ctx, svc.client, cardNames, svc.cardsTransformer.ToBatchCardDataUsingName)
}

func getCardsByName[T *ygo.Cards | *model.BatchCardData[model.CardNames]](ctx context.Context,
	cardServiceClient ygo.CardServiceClient, cardNames model.CardNames,
	transformer func(*ygo.Cards) T) (T, *model.APIError) {
	logger := util.LoggerFromContext(ctx)
	logger.Info(fmt.Sprintf("Fetching card info using %d card name(s)", len(cardNames)))

	if cards, err := cardServiceClient.GetCardsByName(ctx, &ygo.ResourceNames{Names: cardNames}); err != nil {
		logger.Error(
			fmt.Sprintf("There was an issue calling YGO Service. Operation: %s. Code %s. Error: %s",
				"Get Cards By Name", status.Code(err), err))
		return nil, &model.APIError{Message: "There was an error fetching batch card info", StatusCode: http.StatusInternalServerError}
	} else {
		if cards.UnknownResources == nil {
			cards.UnknownResources = make([]string, 0)
		}
		return transformer(cards), nil
	}
}

func (svc YGOServiceV1) GetRandomCard(ctx context.Context, blackListedIDs []string) (*ygo.Card, *model.APIError) {
	return getRandomCard(ctx, svc.client, blackListedIDs, svc.cardTransformer.ToSelf)
}

func (svc YGOServiceV1) GetRandomCardREST(ctx context.Context, blackListedIDs []string) (*model.YGOCardREST, *model.APIError) {
	return getRandomCard(ctx, svc.client, blackListedIDs, svc.cardTransformer.ToREST)
}

func getRandomCard[T *ygo.Card | *model.YGOCardREST | *model.YGOCardGRPC](ctx context.Context,
	cardServiceClient ygo.CardServiceClient, blackListedIDs []string, transformer func(*ygo.Card) T) (T, *model.APIError) {
	logger := util.LoggerFromContext(ctx)
	logger.Info("Getting random card")

	if cards, err := cardServiceClient.GetRandomCard(ctx, &ygo.BlackListed{BlackListedRefs: blackListedIDs}); err != nil {
		logger.Error(
			fmt.Sprintf("There was an issue calling YGO Service. Operation: %s. Code %s. Error: %s",
				"Random Card", status.Code(err), err))
		return nil, &model.APIError{Message: "There was an error fetching random card", StatusCode: http.StatusInternalServerError}
	} else {
		return transformer(cards), nil
	}
}
