package client

import (
	context "context"
	"log/slog"

	"github.com/ygo-skc/skc-go/common/v2/model"
	"github.com/ygo-skc/skc-go/common/v2/util"
	"github.com/ygo-skc/skc-go/common/v2/ygo"
	"google.golang.org/protobuf/types/known/emptypb"
)

type YGOCardClientImp interface {
	GetCardColorsProto(context.Context) (*ygo.CardColors, *model.APIError)
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

func (imp YGOCardClientImpV1) GetCardColorsProto(ctx context.Context) (*ygo.CardColors, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Retrieving card colors")
	cColors, err := imp.client.GetCardColors(ctx, &emptypb.Empty{})
	if err != nil {
		logger.Error("Issue calling YGO Card Service", slog.Any("err", err))
		return nil, rpcErrorToAPIError(err)
	}
	return cColors, nil
}

func (imp YGOCardClientImpV1) GetCardByIDProto(ctx context.Context, cardID string) (*ygo.Card, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Fetching card info", slog.String("card_id", cardID))
	card, err := imp.client.GetCardByID(ctx, &ygo.ResourceID{ID: cardID})
	if err != nil {
		logger.Error("Issue calling YGO Card Service", slog.Any("err", err))
		return nil, rpcErrorToAPIError(err)
	}
	return card, nil
}

func (imp YGOCardClientImpV1) GetCardsByIDProto(ctx context.Context, cardIDs model.CardIDs) (*ygo.Cards, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Fetching card info", slog.Any("card_ids", cardIDs))
	cards, err := imp.client.GetCardsByID(ctx, &ygo.ResourceIDs{IDs: cardIDs})
	if err != nil {
		logger.Error("Issue calling YGO Card Service", slog.Any("err", err))
		return nil, rpcErrorToAPIError(err)
	}
	if cards.UnknownResources == nil {
		cards.UnknownResources = make([]string, 0)
	}
	return cards, nil
}

func (imp YGOCardClientImpV1) GetCardsByNameProto(ctx context.Context, cardNames model.CardNames) (*ygo.Cards, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Fetching card info", slog.Int("card_name_count", len(cardNames)))
	cards, err := imp.client.GetCardsByName(ctx, &ygo.ResourceNames{Names: cardNames})
	if err != nil {
		logger.Error("Issue calling YGO Card Service", slog.Any("err", err))
		return nil, rpcErrorToAPIError(err)
	}
	if cards.UnknownResources == nil {
		cards.UnknownResources = make([]string, 0)
	}
	return cards, nil
}

func (imp YGOCardClientImpV1) GetCardsReferencingNameInEffectProto(ctx context.Context, namesOfCards []string) (*ygo.CardList, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Fetching cards referencing names in card text", slog.Any("names", namesOfCards))
	cards, err := imp.client.GetCardsReferencingNameInEffect(ctx, &ygo.ResourceNames{Names: namesOfCards})
	if err != nil {
		logger.Error("Issue calling YGO Card Service", slog.Any("err", err))
		return nil, rpcErrorToAPIError(err)
	}
	return cards, nil
}

/*
Archetype functionality
*/
func (imp YGOCardClientImpV1) GetArchetypalCardsUsingCardNameProto(ctx context.Context, archetype string) (*ygo.CardList, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Fetching archetypal cards", slog.String("archetype", archetype))
	cards, err := imp.client.GetArchetypalCardsUsingCardName(ctx, &ygo.Archetype{Archetype: archetype})
	if err != nil {
		logger.Error("Issue calling YGO Card Service", slog.Any("err", err))
		return nil, rpcErrorToAPIError(err)
	}
	return cards, nil
}

func (imp YGOCardClientImpV1) GetExplicitArchetypalInclusionsProto(ctx context.Context, archetype string) (*ygo.CardList, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Fetching explicit archetype inclusions", slog.String("archetype", archetype))
	cards, err := imp.client.GetExplicitArchetypalInclusions(ctx, &ygo.Archetype{Archetype: archetype})
	if err != nil {
		logger.Error("Issue calling YGO Card Service", slog.Any("err", err))
		return nil, rpcErrorToAPIError(err)
	}
	return cards, nil
}

func (imp YGOCardClientImpV1) GetExplicitArchetypalExclusionsProto(ctx context.Context, archetype string) (*ygo.CardList, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Fetching explicit archetype exclusions", slog.String("archetype", archetype))
	cards, err := imp.client.GetExplicitArchetypalExclusions(ctx, &ygo.Archetype{Archetype: archetype})
	if err != nil {
		logger.Error("Issue calling YGO Card Service", slog.Any("err", err))
		return nil, rpcErrorToAPIError(err)
	}
	return cards, nil
}

/*
Random card functionality
*/
func (imp YGOCardClientImpV1) GetRandomCardProto(ctx context.Context, blackListedIDs []string) (*ygo.Card, *model.APIError) {
	logger := util.RetrieveLogger(ctx)
	logger.Info("Getting random card")
	card, err := imp.client.GetRandomCard(ctx, &ygo.BlackListed{BlackListedRefs: blackListedIDs})
	if err != nil {
		logger.Error("Issue calling YGO Card Service", slog.Any("err", err))
		return nil, rpcErrorToAPIError(err)
	}
	return card, nil
}
