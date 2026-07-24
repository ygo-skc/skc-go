package client

import (
	context "context"
	"log/slog"

	"github.com/ygo-skc/skc-go/common/v3/model"
	"github.com/ygo-skc/skc-go/common/v3/util"
	"github.com/ygo-skc/skc-go/common/v3/ygo"
)

type YGOCardClientImp interface {
	GetCardColorsProto(context.Context) (*ygo.GetCardColorsResponse, *model.APIError)
	GetCardByIDProto(context.Context, string) (*ygo.Card, *model.APIError)
	GetCardsByIDProto(context.Context, model.CardIDs) (*ygo.Cards, *model.APIError)
	GetCardsByNameProto(context.Context, model.CardNames) (*ygo.Cards, *model.APIError)
	GetCardsReferencingNameInEffectProto(context.Context, []string) (*ygo.CardList, *model.APIError)
	GetArchetypalCardsUsingCardNameProto(context.Context, string) (*ygo.CardList, *model.APIError)
	GetExplicitArchetypalInclusionsProto(context.Context, string) (*ygo.CardList, *model.APIError)
	GetExplicitArchetypalExclusionsProto(context.Context, string) (*ygo.CardList, *model.APIError)
	GetRandomCardProto(context.Context, []string) (*ygo.Card, *model.APIError)
}
type YGOCardClientImpV1 struct {
	client ygo.CardServiceClient
}

func (imp YGOCardClientImpV1) GetCardColorsProto(ctx context.Context) (*ygo.GetCardColorsResponse, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Retrieving card colors")
	cColors, err := imp.client.GetCardColors(ctx, &ygo.GetCardColorsRequest{})
	if err != nil {
		logger.Error("Issue calling YGO Card Service", slog.Any("err", err))
		return nil, rpcErrorToAPIError(err)
	}
	return cColors, nil
}

func (imp YGOCardClientImpV1) GetCardByIDProto(ctx context.Context, cardID string) (*ygo.Card, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Fetching card info using ID", slog.String("ygo_service.resource", cardID))
	res, err := imp.client.GetCardByID(ctx, &ygo.GetCardByIDRequest{Subject: &ygo.ResourceID{Id: cardID}})
	if err != nil {
		logger.Error("Issue calling YGO Card Service", slog.Any("err", err))
		return nil, rpcErrorToAPIError(err)
	}
	return res.Card, nil
}

func (imp YGOCardClientImpV1) GetCardsByIDProto(ctx context.Context, cardIDs model.CardIDs) (*ygo.Cards, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Fetching multi-card info using ID", slog.Any("card_ids", cardIDs))
	res, err := imp.client.GetCardsByID(ctx, &ygo.GetCardsByIDRequest{Subjects: &ygo.ResourceIDs{Ids: cardIDs}})
	if err != nil {
		logger.Error("Issue calling YGO Card Service", slog.Any("err", err))
		return nil, rpcErrorToAPIError(err)
	}
	if res.Cards.UnknownResources == nil {
		res.Cards.UnknownResources = make([]string, 0)
	}
	return res.Cards, nil
}

func (imp YGOCardClientImpV1) GetCardsByNameProto(ctx context.Context, cardNames model.CardNames) (*ygo.Cards, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Fetching multi-card info using card name", slog.Int("card_name_count", len(cardNames)))
	res, err := imp.client.GetCardsByName(ctx, &ygo.GetCardsByNameRequest{Subjects: &ygo.ResourceNames{Names: cardNames}})
	if err != nil {
		logger.Error("Issue calling YGO Card Service", slog.Any("err", err))
		return nil, rpcErrorToAPIError(err)
	}
	if res.Cards.UnknownResources == nil {
		res.Cards.UnknownResources = make([]string, 0)
	}
	return res.Cards, nil
}

func (imp YGOCardClientImpV1) GetCardsReferencingNameInEffectProto(ctx context.Context, namesOfCards []string) (*ygo.CardList, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Fetching cards referencing names in card text", slog.Any("names", namesOfCards))
	res, err := imp.client.GetCardsReferencingNameInEffect(ctx, &ygo.GetCardsReferencingNameInEffectRequest{Subjects: &ygo.ResourceNames{Names: namesOfCards}})
	if err != nil {
		logger.Error("Issue calling YGO Card Service", slog.Any("err", err))
		return nil, rpcErrorToAPIError(err)
	}
	return res.Cards, nil
}

/*
Archetype functionality
*/
func (imp YGOCardClientImpV1) GetArchetypalCardsUsingCardNameProto(ctx context.Context, archetype string) (*ygo.CardList, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Fetching archetypal cards", slog.String("ygo_service.resource", archetype))
	res, err := imp.client.GetArchetypalCardsUsingCardName(ctx, &ygo.GetArchetypalCardsUsingCardNameRequest{Subject: &ygo.Archetype{Name: archetype}})
	if err != nil {
		logger.Error("Issue calling YGO Card Service", slog.Any("err", err))
		return nil, rpcErrorToAPIError(err)
	}
	return res.Cards, nil
}

func (imp YGOCardClientImpV1) GetExplicitArchetypalInclusionsProto(ctx context.Context, archetype string) (*ygo.CardList, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Fetching explicit archetype inclusions", slog.String("ygo_service.resource", archetype))
	res, err := imp.client.GetExplicitArchetypalInclusions(ctx, &ygo.GetExplicitArchetypalInclusionsRequest{Subject: &ygo.Archetype{Name: archetype}})
	if err != nil {
		logger.Error("Issue calling YGO Card Service", slog.Any("err", err))
		return nil, rpcErrorToAPIError(err)
	}
	return res.Cards, nil
}

func (imp YGOCardClientImpV1) GetExplicitArchetypalExclusionsProto(ctx context.Context, archetype string) (*ygo.CardList, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Fetching explicit archetype exclusions", slog.String("ygo_service.resource", archetype))
	res, err := imp.client.GetExplicitArchetypalExclusions(ctx, &ygo.GetExplicitArchetypalExclusionsRequest{Subject: &ygo.Archetype{Name: archetype}})
	if err != nil {
		logger.Error("Issue calling YGO Card Service", slog.Any("err", err))
		return nil, rpcErrorToAPIError(err)
	}
	return res.Cards, nil
}

/*
Random card functionality
*/
func (imp YGOCardClientImpV1) GetRandomCardProto(ctx context.Context, blackListedIDs []string) (*ygo.Card, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Fetching random card")
	res, err := imp.client.GetRandomCard(ctx, &ygo.GetRandomCardRequest{Blacklist: &ygo.BlackListed{BlackListedRefs: blackListedIDs}})
	if err != nil {
		logger.Error("Issue calling YGO Card Service", slog.Any("err", err))
		return nil, rpcErrorToAPIError(err)
	}
	return res.Card, nil
}
