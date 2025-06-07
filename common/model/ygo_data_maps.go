package model

import "github.com/ygo-skc/skc-go/common/ygo"

// =======================
// Map Keys Type
// =======================
type CardIDs []string
type CardNames []string
type ProductIDs []string
type ProductNames []string
type YGOResourceKey interface {
	CardIDs | CardNames | ProductIDs | ProductNames
}

// =======================
// Maps
// =======================
type CardDataMap map[string]YGOCard
type ProductDataMap map[string]YGOProduct
type ProductSummaryDataMap map[string]YGOProductSummary

// =======================
// Batch Resources
// =======================
type BatchCardIDs struct {
	CardIDs CardIDs `json:"cardIDs" validate:"required,ygocardids"`
}

type BatchProductIDs struct {
	ProductIDs ProductIDs `json:"productIDs"`
}

type BatchCardData[RK YGOResourceKey] struct {
	CardInfo         CardDataMap `json:"cardInfo"`
	UnknownResources RK          `json:"unknownResources"`
}

type BatchProductData[RK YGOResourceKey] struct {
	ProductInfo      ProductDataMap `json:"productInfo"`
	UnknownResources RK             `json:"unknownResources"`
}

type BatchProductSummaryData[RK YGOResourceKey] struct {
	ProductInfo      ProductSummaryDataMap `json:"productInfo"`
	UnknownResources RK                    `json:"unknownResources"`
}

type BatchData[RK YGOResourceKey] interface {
	BatchCardData[RK] | BatchProductData[RK] | BatchProductSummaryData[RK]
}

// =======================
// Data Map Key Funcs
// =======================
func FindMissingKeys[T CardIDs | CardNames | ProductIDs | ProductNames, R ygo.Card | ygo.ProductSummary](cards map[string]*R, cardIDs T) T {
	missingIDs := make(T, 0, 10)

	for _, cardID := range cardIDs {
		if _, containsKey := cards[cardID]; !containsKey {
			missingIDs = append(missingIDs, cardID)
		}
	}
	return missingIDs
}

func CardIDAsKey(c *ygo.Card) string   { return c.ID }
func CardNameAsKey(c *ygo.Card) string { return c.Name }

func ProductIDAsKey(p *ygo.ProductSummary) string   { return p.ID }
func ProductNameAsKey(p *ygo.ProductSummary) string { return p.Name }
