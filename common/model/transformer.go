package model

import "github.com/ygo-skc/skc-go/common/ygo"

func YGOCardRESTFromPB(c *ygo.Card) YGOCard {
	ygoCardGRPC := YGOCardGRPC{Card: c}
	return YGOCardREST{
		ID:          ygoCardGRPC.GetID(),
		Color:       ygoCardGRPC.GetColor(),
		Name:        ygoCardGRPC.GetName(),
		Attribute:   ygoCardGRPC.GetAttribute(),
		Effect:      ygoCardGRPC.GetEffect(),
		MonsterType: ygoCardGRPC.GetMonsterType(),
		Attack:      ygoCardGRPC.GetAttack(),
		Defense:     ygoCardGRPC.GetDefense(),
	}
}

type YGOCardsTransformer interface {
	ToSelf(*ygo.Cards) *ygo.Cards
	ToBatchCardDataUsingID(*ygo.Cards) *BatchCardData[CardIDs]
	ToBatchCardDataUsingName(*ygo.Cards) *BatchCardData[CardNames]
}
type YGOCardsTransformerV1 struct{}

func (t YGOCardsTransformerV1) ToSelf(c *ygo.Cards) *ygo.Cards {
	return c
}

func (t YGOCardsTransformerV1) ToBatchCardDataUsingID(c *ygo.Cards) *BatchCardData[CardIDs] {
	return BatchCardDataFromPB[CardIDs](c)
}

func (t YGOCardsTransformerV1) ToBatchCardDataUsingName(c *ygo.Cards) *BatchCardData[CardNames] {
	return BatchCardDataFromPB[CardNames](c)
}

func BatchCardDataFromPB[T CardIDs | CardNames](c *ygo.Cards) *BatchCardData[T] {
	batchCardData := make(CardDataMap, len(c.CardInfo))
	for k, v := range c.CardInfo {
		batchCardData[k] = YGOCardRESTFromPB(v)
	}
	return &BatchCardData[T]{CardInfo: batchCardData, UnknownResources: c.UnknownResources}
}
