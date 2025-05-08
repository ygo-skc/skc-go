package model

import (
	"sort"
	"strings"

	"github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/common/ygo"
)

type YGOCard interface {
	YGOResource
	GetColor() string
	GetAttribute() string
	GetEffect() string
	GetMonsterType() *string
	GetAttack() *uint32
	GetDefense() *uint32
}
type YGOCards []YGOCard

/*
YGOCardREST Struct and CardDescriptor conformance
*/
type YGOCardREST struct {
	ID          string  `db:"card_number" json:"cardID"`
	Color       string  `db:"card_color" json:"cardColor"`
	Name        string  `db:"card_name" json:"cardName"`
	Attribute   string  `db:"card_attribute" json:"cardAttribute"`
	Effect      string  `db:"card_effect" json:"cardEffect"`
	MonsterType *string `db:"monster_type" json:"monsterType,omitempty"`
	Attack      *uint32 `db:"monster_attack" json:"monsterAttack,omitempty"`
	Defense     *uint32 `db:"monster_defense" json:"monsterDefense,omitempty"`
}

func (c YGOCardREST) GetID() string           { return c.ID }
func (c YGOCardREST) GetColor() string        { return c.Color }
func (c YGOCardREST) GetName() string         { return c.Name }
func (c YGOCardREST) GetAttribute() string    { return c.Attribute }
func (c YGOCardREST) GetEffect() string       { return c.Effect }
func (c YGOCardREST) GetMonsterType() *string { return c.MonsterType }
func (c YGOCardREST) GetAttack() *uint32      { return c.Attack }
func (c YGOCardREST) GetDefense() *uint32     { return c.Defense }

func (c YGOCardREST) ToPB() *ygo.Card {
	return &ygo.Card{
		ID:          c.ID,
		Color:       c.Color,
		Name:        c.Name,
		Attribute:   c.Attribute,
		Effect:      c.Effect,
		MonsterType: util.PBStringValue(c.MonsterType),
		Attack:      util.PBUInt32Value(c.Attack),
		Defense:     util.PBUInt32Value(c.Defense),
	}
}

/*
YGOCardGRPC Struct and CardDescriptor conformance
*/
type YGOCardGRPC struct{ *ygo.Card }

func (c YGOCardGRPC) GetID() string        { return c.ID }
func (c YGOCardGRPC) GetColor() string     { return c.Color }
func (c YGOCardGRPC) GetName() string      { return c.Name }
func (c YGOCardGRPC) GetAttribute() string { return c.Attribute }
func (c YGOCardGRPC) GetEffect() string    { return c.Effect }
func (c YGOCardGRPC) GetMonsterType() *string {
	if c.MonsterType == nil {
		return nil
	}
	return &c.MonsterType.Value
}
func (c YGOCardGRPC) GetAttack() *uint32 {
	if c.Attack == nil {
		return nil
	}
	return &c.Attack.Value
}
func (c YGOCardGRPC) GetDefense() *uint32 {
	if c.Defense == nil {
		return nil
	}
	return &c.Defense.Value
}

/*
CardDescriptor helper functions
*/
func IsExtraDeckMonster(c YGOCard) bool {
	color := strings.ToUpper(c.GetColor())
	return strings.Contains(color, "FUSION") || strings.Contains(color, "SYNCHRO") || strings.Contains(color, "XYZ") || strings.Contains(color, "PENDULUM") || strings.Contains(color, "LINK")
}

// Uses new line as delimiter to split card effect. Materials are found in the first token.
func GetPotentialMaterialsAsString(c YGOCard) string {
	var effectTokens []string

	if !IsExtraDeckMonster(c) {
		return ""
	}

	color := strings.ToUpper(c.GetColor())
	if strings.Contains(color, "PENDULUM") && color != "PENDULUM-EFFECT" && color != "PENDULUM-NORMAL" {
		effectTokens = strings.SplitAfter(strings.SplitAfter(c.GetEffect(), "\n\nMonster Effect\n")[1], "\n")
	} else {
		effectTokens = strings.SplitAfter(c.GetEffect(), "\n")
	}

	if len(effectTokens) < 2 {
		return c.GetEffect()
	}
	return effectTokens[0]
}

func IsCardNameInTokens(c YGOCard, tokens []QuotedToken) bool {
	isFound := false

	for _, token := range tokens {
		CleanupToken(&token)

		if strings.EqualFold(c.GetName(), token) {
			isFound = true
			break
		}
	}

	return isFound
}

func (c YGOCards) SortCardsByName() {
	sort.SliceStable(c, func(i, j int) bool {
		return (c)[i].GetName() < (c)[j].GetName()
	})
}

// cleans up a quoted string found in card text so its easier to parse
func CleanupToken(t *QuotedToken) {
	*t = strings.TrimSpace(*t)
	*t = strings.ReplaceAll(*t, `".`, "")
	*t = strings.ReplaceAll(*t, `",`, "")
	*t = strings.ReplaceAll(*t, "'.", "")
	*t = strings.ReplaceAll(*t, "',", "")

	*t = strings.Trim(*t, "'")
	*t = strings.Trim(*t, `"`)
}
