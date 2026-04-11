package model

import (
	"github.com/ygo-skc/skc-go/common/v2/util"
	"github.com/ygo-skc/skc-go/common/v2/ygo"
)

type YGOCardProtoBuilder struct {
	c *ygo.Card
}

func NewYGOCardProtoBuilder(id, name string) *YGOCardProtoBuilder {
	return &YGOCardProtoBuilder{
		c: &ygo.Card{
			ID:   id,
			Name: name,
		},
	}
}
func (b *YGOCardProtoBuilder) WithColor(color string) *YGOCardProtoBuilder {
	b.c.Color = color
	return b
}
func (b *YGOCardProtoBuilder) WithAttribute(attribute string) *YGOCardProtoBuilder {
	b.c.Attribute = attribute
	return b
}
func (b *YGOCardProtoBuilder) WithEffect(effect string) *YGOCardProtoBuilder {
	b.c.Effect = effect
	return b
}
func (b *YGOCardProtoBuilder) WithMonsterType(monsterType *string) *YGOCardProtoBuilder {
	b.c.MonsterType = util.ProtoStringValue(monsterType)
	return b
}
func (b *YGOCardProtoBuilder) WithAttack(atk *uint32) *YGOCardProtoBuilder {
	b.c.Attack = util.ProtoUInt32Value(atk)
	return b
}
func (b *YGOCardProtoBuilder) WithDefense(def *uint32) *YGOCardProtoBuilder {
	b.c.Defense = util.ProtoUInt32Value(def)
	return b
}
func (b *YGOCardProtoBuilder) Build() *ygo.Card {
	return b.c
}
