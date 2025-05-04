package model

import (
	"sort"
	"strings"

	"github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/common/ygo"
)

type Card struct {
	ID          string  `db:"card_number" json:"cardID"`
	Color       string  `db:"card_color" json:"cardColor"`
	Name        string  `db:"card_name" json:"cardName"`
	Attribute   string  `db:"card_attribute" json:"cardAttribute"`
	Effect      string  `db:"card_effect" json:"cardEffect"`
	MonsterType *string `db:"monster_type" json:"monsterType,omitempty"`
	Attack      *uint32 `db:"monster_attack" json:"monsterAttack,omitempty"`
	Defense     *uint32 `db:"monster_defense" json:"monsterDefense,omitempty"`
}

func (c Card) IsExtraDeckMonster() bool {
	color := strings.ToUpper(c.Color)
	return strings.Contains(color, "FUSION") || strings.Contains(color, "SYNCHRO") || strings.Contains(color, "XYZ") || strings.Contains(color, "PENDULUM") || strings.Contains(color, "LINK")
}

// Uses new line as delimiter to split card effect. Materials are found in the first token.
func (card Card) GetPotentialMaterialsAsString() string {
	var effectTokens []string

	if !card.IsExtraDeckMonster() {
		return ""
	}

	color := strings.ToUpper(card.Color)
	if strings.Contains(color, "PENDULUM") && color != "PENDULUM-EFFECT" && color != "PENDULUM-NORMAL" {
		effectTokens = strings.SplitAfter(strings.SplitAfter(card.Effect, "\n\nMonster Effect\n")[1], "\n")
	} else {
		effectTokens = strings.SplitAfter(card.Effect, "\n")
	}

	if len(effectTokens) < 2 {
		return card.Effect
	}
	return effectTokens[0]
}

func (c Card) IsCardNameInTokens(tokens []QuotedToken) bool {
	isFound := false

	for _, token := range tokens {
		CleanupToken(&token)

		if strings.EqualFold(c.Name, token) {
			isFound = true
			break
		}
	}

	return isFound
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

type Cards []Card

func (cards Cards) SortCardsByName() {
	sort.SliceStable(cards, func(i, j int) bool {
		return (cards)[i].Name < (cards)[j].Name
	})
}

func (c Card) ToPB() *ygo.Card {
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
