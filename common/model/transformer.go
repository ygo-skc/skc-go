package model

import "github.com/ygo-skc/skc-go/common/ygo"

type YGOCardTransformer interface {
	ToSelf(*ygo.Card) *ygo.Card
	ToREST(*ygo.Card) *YGOCardREST
}

type YGOCardTransformerV1 struct{}

func (t YGOCardTransformerV1) ToSelf(c *ygo.Card) *ygo.Card {
	return c
}

func (t YGOCardTransformerV1) ToREST(c *ygo.Card) *YGOCardREST {
	return YGOCardRESTFromPB(c)
}

func YGOCardRESTFromPB(c *ygo.Card) *YGOCardREST {
	ygoCardGRPC := YGOCardGRPC{Card: c}
	return &YGOCardREST{ID: ygoCardGRPC.GetID(), Color: ygoCardGRPC.GetColor(), Name: ygoCardGRPC.GetName(),
		Attribute: ygoCardGRPC.GetAttribute(), Effect: ygoCardGRPC.GetEffect(), MonsterType: ygoCardGRPC.GetMonsterType(),
		Attack: ygoCardGRPC.GetAttack(), Defense: ygoCardGRPC.GetDefense()}
}

type YGOCardsTransformer interface {
	ToSelf(*ygo.Cards) *ygo.Cards
	ToREST(*ygo.Cards) *BatchCardData[CardIDs]
}
type YGOCardsTransformerV1 struct{}

func (t YGOCardsTransformerV1) ToSelf(c *ygo.Cards) *ygo.Cards {
	return c
}

func (t YGOCardsTransformerV1) ToREST(c *ygo.Cards) *BatchCardData[CardIDs] {
	return BatchCardDataFromPB(c)
}

func BatchCardDataFromPB(c *ygo.Cards) *BatchCardData[CardIDs] {
	batchCardData := make(CardDataMap, len(c.CardInfo))
	for k, v := range c.CardInfo {
		batchCardData[k] = YGOCardRESTFromPB(v)
	}
	return &BatchCardData[CardIDs]{CardInfo: batchCardData, UnknownResources: c.UnknownResources}
}
