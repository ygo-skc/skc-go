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

type ProductContent struct {
	Card            Card     `json:"card"`
	ProductPosition string   `json:"productPosition"`
	Rarities        []string `json:"rarities"`
}
