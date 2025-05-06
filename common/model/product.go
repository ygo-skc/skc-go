package model

type Product struct {
	ID          string           `json:"productId"`
	Locale      string           `json:"productLocale"`
	Name        string           `json:"productName"`
	Type        string           `json:"productType"`
	SubType     string           `json:"productSubType"`
	ReleaseDate string           `json:"productReleaseDate"`
	Total       int              `json:"productTotal,omitempty"`
	RarityStats map[string]int   `json:"productRarityStats,omitempty"`
	Content     []ProductContent `json:"productContent,omitempty"`
}

func (c Product) GetID() string   { return c.ID }
func (c Product) GetName() string { return c.Name }

type ProductContent struct {
	Card            YGOCardREST `json:"card"`
	ProductPosition string      `json:"productPosition"`
	Rarities        []string    `json:"rarities"`
}
