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
	GetCardColorsProto(context.Context) (*ygo.CardColors, *model.APIError)

	GetCardByIDProto(context.Context, string) (*ygo.Card, *model.APIError)
	GetCardByID(context.Context, string) (*model.YGOCard, *model.APIError)

	GetCardsByIDProto(context.Context, model.CardIDs) (*ygo.Cards, *model.APIError)
	GetCardsByID(context.Context, model.CardIDs) (*model.BatchCardData[model.CardIDs], *model.APIError)

	GetCardsByNameProto(context.Context, model.CardNames) (*ygo.Cards, *model.APIError)
	GetCardsByName(context.Context, model.CardNames) (*model.BatchCardData[model.CardNames], *model.APIError)

	GetRandomCardProto(context.Context, []string) (*ygo.Card, *model.APIError)
	GetRandomCard(context.Context, []string) (*model.YGOCard, *model.APIError)
}

type YGOServiceV1 struct {
	client ygo.CardServiceClient
}

func NewYGOServiceV1(client ygo.CardServiceClient) YGOServiceV1 {
	return YGOServiceV1{
		client: client,
	}
}

func (svc YGOServiceV1) GetCardColorsProto(ctx context.Context) (*ygo.CardColors, *model.APIError) {
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

func (svc YGOServiceV1) GetCardByIDProto(ctx context.Context, cardID string) (*ygo.Card, *model.APIError) {
	return getCardByID(ctx, svc.client, cardID)
}

func (svc YGOServiceV1) GetCardByID(ctx context.Context, cardID string) (*model.YGOCard, *model.APIError) {
	c, err := getCardByID(ctx, svc.client, cardID)
	if err == nil {
		card := model.YGOCardRESTFromProto(c)
		return &card, nil
	}
	return nil, err
}

func getCardByID(ctx context.Context, cardServiceClient ygo.CardServiceClient, cardID string) (*ygo.Card, *model.APIError) {
	logger := util.LoggerFromContext(ctx)
	logger.Info(fmt.Sprintf("Fetching card info using ID: %v", cardID))

	if cards, err := cardServiceClient.GetCardByID(ctx, &ygo.ResourceID{ID: cardID}); err != nil {
		logger.Error(
			fmt.Sprintf("There was an issue calling YGO Service. Operation: %s. Code %s. Error: %s",
				"Get Card By ID", status.Code(err), err))
		return nil, &model.APIError{Message: "There was an error fetching card info", StatusCode: http.StatusInternalServerError}
	} else {
		return cards, nil
	}
}

func (svc YGOServiceV1) GetCardsByIDProto(ctx context.Context, cardIDs model.CardIDs) (*ygo.Cards, *model.APIError) {
	return getCardsByID(ctx, svc.client, cardIDs)
}

func (svc YGOServiceV1) GetCardsByID(ctx context.Context, cardIDs model.CardIDs) (*model.BatchCardData[model.CardIDs], *model.APIError) {
	c, err := getCardsByID(ctx, svc.client, cardIDs)
	if err == nil {
		return model.BatchCardDataFromProto[model.CardIDs](c), nil
	}
	return nil, err
}

func getCardsByID(ctx context.Context, cardServiceClient ygo.CardServiceClient, cardIDs model.CardIDs) (*ygo.Cards, *model.APIError) {
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
		return cards, nil
	}
}

func (svc YGOServiceV1) GetCardsByNameProto(ctx context.Context, cardNames model.CardNames) (*ygo.Cards, *model.APIError) {
	return getCardsByName(ctx, svc.client, cardNames)
}

func (svc YGOServiceV1) GetCardsByName(ctx context.Context, cardNames model.CardNames) (*model.BatchCardData[model.CardNames], *model.APIError) {
	c, err := getCardsByName(ctx, svc.client, cardNames)
	if err == nil {
		return model.BatchCardDataFromProto[model.CardNames](c), nil
	}
	return nil, err
}

func getCardsByName(ctx context.Context, cardServiceClient ygo.CardServiceClient, cardNames model.CardNames) (*ygo.Cards, *model.APIError) {
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
		return cards, nil
	}
}

func (svc YGOServiceV1) GetRandomCardProto(ctx context.Context, blackListedIDs []string) (*ygo.Card, *model.APIError) {
	return getRandomCard(ctx, svc.client, blackListedIDs)
}

func (svc YGOServiceV1) GetRandomCard(ctx context.Context, blackListedIDs []string) (*model.YGOCard, *model.APIError) {
	c, err := getRandomCard(ctx, svc.client, blackListedIDs)
	if err == nil {
		card := model.YGOCardRESTFromProto(c)
		return &card, nil
	}
	return nil, err
}

func getRandomCard(ctx context.Context,
	cardServiceClient ygo.CardServiceClient, blackListedIDs []string) (*ygo.Card, *model.APIError) {
	logger := util.LoggerFromContext(ctx)
	logger.Info("Getting random card")

	if card, err := cardServiceClient.GetRandomCard(ctx, &ygo.BlackListed{BlackListedRefs: blackListedIDs}); err != nil {
		logger.Error(
			fmt.Sprintf("There was an issue calling YGO Service. Operation: %s. Code %s. Error: %s",
				"Random Card", status.Code(err), err))
		return nil, &model.APIError{Message: "There was an error fetching random card", StatusCode: http.StatusInternalServerError}
	} else {
		return card, nil
	}
}
