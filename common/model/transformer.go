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
		MonsterType: util.ProtoStringValue(monsterType),
		Attack:      util.ProtoUInt32Value(atk),
		Defense:     util.ProtoUInt32Value(def),
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

func YGOCardListRESTFromProto(c *ygo.CardList) []YGOCard {
	cards := make([]YGOCard, len(c.Cards))
	for i, c := range c.Cards {
		ygoCardGRPC := YGOCardGRPC{Card: c}
		cards[i] = YGOCardREST{
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
	return cards
}

func BatchCardDataFromProto[T CardIDs | CardNames](c *ygo.Cards) *BatchCardData[T] {
	batchCardData := make(CardDataMap, len(c.CardInfo))
	for k, v := range c.CardInfo {
		batchCardData[k] = YGOCardRESTFromProto(v)
	}
	return &BatchCardData[T]{CardInfo: batchCardData, UnknownResources: c.UnknownResources}
}
