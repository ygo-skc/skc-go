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

func BatchCardDataFromPB[T CardIDs | CardNames](c *ygo.Cards) *BatchCardData[T] {
	batchCardData := make(CardDataMap, len(c.CardInfo))
	for k, v := range c.CardInfo {
		batchCardData[k] = YGOCardRESTFromPB(v)
	}
	return &BatchCardData[T]{CardInfo: batchCardData, UnknownResources: c.UnknownResources}
}
