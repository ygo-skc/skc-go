package model

import (
	"github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/common/ygo"
)

func NewYgoCardProto(id string, color string, name string, attribute string, effect string, monsterType *string,
	atk *uint32, def *uint32) *ygo.Card {
	return &ygo.Card{
		ID:          id,
		Color:       color,
		Name:        name,
		Attribute:   attribute,
		Effect:      effect,
		MonsterType: util.PBStringValue(monsterType),
		Attack:      util.PBUInt32Value(atk),
		Defense:     util.PBUInt32Value(def),
	}
}

func YGOCardRESTFromProto(c *ygo.Card) YGOCard {
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

func BatchCardDataFromProto[T CardIDs | CardNames](c *ygo.Cards) *BatchCardData[T] {
	batchCardData := make(CardDataMap, len(c.CardInfo))
	for k, v := range c.CardInfo {
		batchCardData[k] = YGOCardRESTFromProto(v)
	}
	return &BatchCardData[T]{CardInfo: batchCardData, UnknownResources: c.UnknownResources}
}
