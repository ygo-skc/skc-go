package client

import (
	context "context"
	"fmt"
	"net/http"

	"github.com/ygo-skc/skc-go/common/v2/model"
	"github.com/ygo-skc/skc-go/common/v2/util"
	"github.com/ygo-skc/skc-go/common/v2/ygo"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type YGOCardClientImp interface {
	GetCardColorsProto(context.Context) (*ygo.CardColors, *model.APIError)

	GetCardByIDProto(context.Context, string) (*ygo.Card, *model.APIError)
	GetCardByID(context.Context, string) (*model.YGOCard, *model.APIError)

	GetCardsByIDProto(context.Context, model.CardIDs) (*ygo.Cards, *model.APIError)
	GetCardsByID(context.Context, model.CardIDs) (*model.BatchCardData[model.CardIDs], *model.APIError)

	GetCardsByNameProto(context.Context, model.CardNames) (*ygo.Cards, *model.APIError)
	GetCardsByName(context.Context, model.CardNames) (*model.BatchCardData[model.CardNames], *model.APIError)

	GetCardsReferencingNameInEffectProto(context.Context, []string) (*ygo.CardList, *model.APIError)
	GetCardsReferencingNameInEffect(context.Context, []string) ([]model.YGOCard, *model.APIError)

	GetArchetypalCardsUsingCardNameProto(context.Context, string) (*ygo.CardList, *model.APIError)
	GetArchetypalCardsUsingCardName(context.Context, string) ([]model.YGOCard, *model.APIError)

	GetExplicitArchetypalInclusionsProto(context.Context, string) (*ygo.CardList, *model.APIError)
	GetExplicitArchetypalInclusions(context.Context, string) ([]model.YGOCard, *model.APIError)

	GetExplicitArchetypalExclusionsProto(context.Context, string) (*ygo.CardList, *model.APIError)
	GetExplicitArchetypalExclusions(context.Context, string) ([]model.YGOCard, *model.APIError)

	GetRandomCardProto(context.Context, []string) (*ygo.Card, *model.APIError)
	GetRandomCard(context.Context, []string) (*model.YGOCard, *model.APIError)
}
type YGOCardClientImpV1 struct {
	client ygo.CardServiceClient
}

const (
	ygoCardClientErr = "There was an issue calling YGO Card Service. Operation: %s. Code %s. Error: %s"
)

func (imp YGOCardClientImpV1) GetCardColorsProto(ctx context.Context) (*ygo.CardColors, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Retrieving card colors")

	if cColors, err := imp.client.GetCardColors(ctx, &emptypb.Empty{}); err != nil {
		logger.Error(fmt.Sprintf(ygoCardClientErr, "Get Card Colors", status.Code(err), err))
		return nil, &model.APIError{Message: "Error fetching card color data", StatusCode: http.StatusInternalServerError}
	} else {
		return cColors, nil
	}
}

func (imp YGOCardClientImpV1) GetCardByIDProto(ctx context.Context, cardID string) (*ygo.Card, *model.APIError) {
	return getCardByID(ctx, imp.client, cardID)
}

func (imp YGOCardClientImpV1) GetCardByID(ctx context.Context, cardID string) (*model.YGOCard, *model.APIError) {
	c, err := getCardByID(ctx, imp.client, cardID)
	if err == nil {
		card := model.YGOCardRESTFromProto(c)
		return &card, nil
	}
	return nil, err
}

func getCardByID(ctx context.Context, client ygo.CardServiceClient, cardID string) (*ygo.Card, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info(fmt.Sprintf("Fetching card info using ID: %v", cardID))

	if cards, err := client.GetCardByID(ctx, &ygo.ResourceID{ID: cardID}); err != nil {
		logger.Error(fmt.Sprintf(ygoCardClientErr, "Get Card By ID", status.Code(err), err))
		return nil, &model.APIError{Message: "Error fetching card info", StatusCode: http.StatusInternalServerError}
	} else {
		return cards, nil
	}
}

func (imp YGOCardClientImpV1) GetCardsByIDProto(ctx context.Context, cardIDs model.CardIDs) (*ygo.Cards, *model.APIError) {
	return getCardsByID(ctx, imp.client, cardIDs)
}

func (imp YGOCardClientImpV1) GetCardsByID(ctx context.Context, cardIDs model.CardIDs) (*model.BatchCardData[model.CardIDs], *model.APIError) {
	c, err := getCardsByID(ctx, imp.client, cardIDs)
	if err == nil {
		return model.BatchCardDataFromProto[model.CardIDs](c, model.CardIDAsKey), nil
	}
	return nil, err
}

func getCardsByID(ctx context.Context, client ygo.CardServiceClient, cardIDs model.CardIDs) (*ygo.Cards, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info(fmt.Sprintf("Fetching card info for the following IDs: %v", cardIDs))

	if cards, err := client.GetCardsByID(ctx, &ygo.ResourceIDs{IDs: cardIDs}); err != nil {
		logger.Error(fmt.Sprintf(ygoCardClientErr, "Get Cards By ID", status.Code(err), err))
		return nil, &model.APIError{Message: "Error fetching batch card info", StatusCode: http.StatusInternalServerError}
	} else {
		if cards.UnknownResources == nil {
			cards.UnknownResources = make([]string, 0)
		}
		return cards, nil
	}
}

func (imp YGOCardClientImpV1) GetCardsByNameProto(ctx context.Context, cardNames model.CardNames) (*ygo.Cards, *model.APIError) {
	return getCardsByName(ctx, imp.client, cardNames)
}

func (imp YGOCardClientImpV1) GetCardsByName(ctx context.Context, cardNames model.CardNames) (*model.BatchCardData[model.CardNames], *model.APIError) {
	c, err := getCardsByName(ctx, imp.client, cardNames)
	if err == nil {
		return model.BatchCardDataFromProto[model.CardNames](c, model.CardNameAsKey), nil
	}
	return nil, err
}

func getCardsByName(ctx context.Context, client ygo.CardServiceClient, cardNames model.CardNames) (*ygo.Cards, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info(fmt.Sprintf("Fetching card info using %d card name(s)", len(cardNames)))

	if cards, err := client.GetCardsByName(ctx, &ygo.ResourceNames{Names: cardNames}); err != nil {
		logger.Error(fmt.Sprintf(ygoCardClientErr, "Get Cards By Name", status.Code(err), err))
		return nil, &model.APIError{Message: "Error fetching batch card info", StatusCode: http.StatusInternalServerError}
	} else {
		if cards.UnknownResources == nil {
			cards.UnknownResources = make([]string, 0)
		}
		return cards, nil
	}
}

func (imp YGOCardClientImpV1) GetCardsReferencingNameInEffectProto(ctx context.Context, namesOfCards []string) (*ygo.CardList, *model.APIError) {
	return getCardsReferencingNameInEffect(ctx, imp.client, namesOfCards)
}

func (imp YGOCardClientImpV1) GetCardsReferencingNameInEffect(ctx context.Context, namesOfCards []string) ([]model.YGOCard, *model.APIError) {
	c, err := getCardsReferencingNameInEffect(ctx, imp.client, namesOfCards)
	if err == nil {
		return model.YGOCardListRESTFromProto(c), nil
	}
	return nil, err
}

func getCardsReferencingNameInEffect(ctx context.Context, client ygo.CardServiceClient, namesOfCards []string) (*ygo.CardList, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info(fmt.Sprintf("Fetching cards that reference the following names in their text %v", namesOfCards))

	if cards, err := client.GetCardsReferencingNameInEffect(ctx, &ygo.ResourceNames{Names: namesOfCards}); err != nil {
		logger.Error(fmt.Sprintf(ygoCardClientErr, "Search Card References Using Text", status.Code(err), err))
		return nil, &model.APIError{Message: "Error searching card text for references", StatusCode: http.StatusInternalServerError}
	} else {
		return cards, nil
	}
}

/*
Archetype functionality
*/
func (imp YGOCardClientImpV1) GetArchetypalCardsUsingCardNameProto(ctx context.Context, archetype string) (*ygo.CardList, *model.APIError) {
	return getArchetypalCardsUsingCardName(ctx, imp.client, archetype)
}

func (imp YGOCardClientImpV1) GetArchetypalCardsUsingCardName(ctx context.Context, archetype string) ([]model.YGOCard, *model.APIError) {
	c, err := getArchetypalCardsUsingCardName(ctx, imp.client, archetype)
	if err == nil {
		return model.YGOCardListRESTFromProto(c), nil
	}
	return nil, err
}

func getArchetypalCardsUsingCardName(ctx context.Context, client ygo.CardServiceClient,
	archetype string) (*ygo.CardList, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info(fmt.Sprintf("Fetching cards with %s in their name", archetype))

	if cards, err := client.GetArchetypalCardsUsingCardName(ctx, &ygo.Archetype{Archetype: archetype}); err != nil {
		logger.Error(fmt.Sprintf(ygoCardClientErr, "Get Archetypal Cards Using Name", status.Code(err), err))
		return nil, &model.APIError{Message: "Error fetching archetypal data", StatusCode: http.StatusInternalServerError}
	} else {
		return cards, nil
	}
}

func (imp YGOCardClientImpV1) GetExplicitArchetypalInclusionsProto(ctx context.Context, archetype string) (*ygo.CardList, *model.APIError) {
	return getArchetypalCardsUsingCardName(ctx, imp.client, archetype)
}

func (imp YGOCardClientImpV1) GetExplicitArchetypalInclusions(ctx context.Context, archetype string) ([]model.YGOCard, *model.APIError) {
	c, err := getExplicitArchetypalInclusions(ctx, imp.client, archetype)
	if err == nil {
		return model.YGOCardListRESTFromProto(c), nil
	}
	return nil, err
}

func getExplicitArchetypalInclusions(ctx context.Context, client ygo.CardServiceClient, archetype string) (*ygo.CardList, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info(fmt.Sprintf("Fetching cards that are explicitly included from archetype %s", archetype))

	if cards, err := client.GetExplicitArchetypalInclusions(ctx, &ygo.Archetype{Archetype: archetype}); err != nil {
		logger.Error(fmt.Sprintf(ygoCardClientErr, "Get Explicit Archetype Inclusions", status.Code(err), err))
		return nil, &model.APIError{Message: "Error fetching explicit archetype inclusions", StatusCode: http.StatusInternalServerError}
	} else {
		return cards, nil
	}
}

func (imp YGOCardClientImpV1) GetExplicitArchetypalExclusionsProto(ctx context.Context, archetype string) (*ygo.CardList, *model.APIError) {
	return getExplicitArchetypalExclusions(ctx, imp.client, archetype)
}

func (imp YGOCardClientImpV1) GetExplicitArchetypalExclusions(ctx context.Context, archetype string) ([]model.YGOCard, *model.APIError) {
	c, err := getExplicitArchetypalExclusions(ctx, imp.client, archetype)
	if err == nil {
		return model.YGOCardListRESTFromProto(c), nil
	}
	return nil, err
}

func getExplicitArchetypalExclusions(ctx context.Context, client ygo.CardServiceClient, archetype string) (*ygo.CardList, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info(fmt.Sprintf("Fetching cards that are explicitly excluded from archetype %s", archetype))

	if cards, err := client.GetExplicitArchetypalExclusions(ctx, &ygo.Archetype{Archetype: archetype}); err != nil {
		logger.Error(fmt.Sprintf(ygoCardClientErr, "Get Explicit Archetype Exclusions", status.Code(err), err))
		return nil, &model.APIError{Message: "Error fetching explicit archetype exclusions", StatusCode: http.StatusInternalServerError}
	} else {
		return cards, nil
	}
}

/*
Random card functionality
*/
func (imp YGOCardClientImpV1) GetRandomCardProto(ctx context.Context, blackListedIDs []string) (*ygo.Card, *model.APIError) {
	return getRandomCard(ctx, imp.client, blackListedIDs)
}

func (imp YGOCardClientImpV1) GetRandomCard(ctx context.Context, blackListedIDs []string) (*model.YGOCard, *model.APIError) {
	c, err := getRandomCard(ctx, imp.client, blackListedIDs)
	if err == nil {
		card := model.YGOCardRESTFromProto(c)
		return &card, nil
	}
	return nil, err
}

func getRandomCard(ctx context.Context,
	client ygo.CardServiceClient, blackListedIDs []string) (*ygo.Card, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Getting random card")

	if card, err := client.GetRandomCard(ctx, &ygo.BlackListed{BlackListedRefs: blackListedIDs}); err != nil {
		logger.Error(fmt.Sprintf(ygoCardClientErr, "Random Card", status.Code(err), err))
		return nil, &model.APIError{Message: "Error fetching random card", StatusCode: http.StatusInternalServerError}
	} else {
		return card, nil
	}
}
