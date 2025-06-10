package model

// =======================
// YGO Product
// =======================
type YGOProduct interface {
	YGOResource
	GetLocale() string
	GetType() string
	GetSubType() string
	GetReleaseDate() string
	GetTotal() int
	GetRarityStats() map[string]int
	GetContent() []ProductContent
}
type YGOProductREST struct {
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

func (p YGOProductREST) GetID() string                  { return p.ID }
func (p YGOProductREST) GetName() string                { return p.Name }
func (p YGOProductREST) GetLocale() string              { return p.Locale }
func (p YGOProductREST) GetType() string                { return p.Type }
func (p YGOProductREST) GetSubType() string             { return p.SubType }
func (p YGOProductREST) GetReleaseDate() string         { return p.ReleaseDate }
func (p YGOProductREST) GetTotal() int                  { return p.Total }
func (p YGOProductREST) GetRarityStats() map[string]int { return p.RarityStats }
func (p YGOProductREST) GetContent() []ProductContent   { return p.Content }

// =======================
// Product Content
// =====================
type ProductContent struct {
	Card            YGOProduct `json:"card"`
	ProductPosition string     `json:"productPosition"`
	Rarities        []string   `json:"rarities"`
}

// =======================
// Product Summary
// =====================
type YGOProductSummary interface {
	YGOResource
	GetLocale() string
	GetType() string
	GetSubType() string
	GetReleaseDate() string
	GetTotal() int
}
type YGOProductSummaryREST struct {
	ID          string `json:"productId"`
	Locale      string `json:"productLocale"`
	Name        string `json:"productName"`
	Type        string `json:"productType"`
	SubType     string `json:"productSubType"`
	ReleaseDate string `json:"productReleaseDate"`
	Total       int    `json:"productTotal,omitempty"`
}

func (p YGOProductSummaryREST) GetID() string          { return p.ID }
func (p YGOProductSummaryREST) GetName() string        { return p.Name }
func (p YGOProductSummaryREST) GetLocale() string      { return p.Locale }
func (p YGOProductSummaryREST) GetType() string        { return p.Type }
func (p YGOProductSummaryREST) GetSubType() string     { return p.SubType }
func (p YGOProductSummaryREST) GetReleaseDate() string { return p.ReleaseDate }
func (p YGOProductSummaryREST) GetTotal() int          { return p.Total }
