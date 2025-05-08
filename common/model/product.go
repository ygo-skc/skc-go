package model

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

func (c YGOProductREST) GetID() string                  { return c.ID }
func (c YGOProductREST) GetName() string                { return c.Name }
func (c YGOProductREST) GetLocale() string              { return c.Locale }
func (c YGOProductREST) GetType() string                { return c.Type }
func (c YGOProductREST) GetSubType() string             { return c.SubType }
func (c YGOProductREST) GetReleaseDate() string         { return c.ReleaseDate }
func (c YGOProductREST) GetTotal() int                  { return c.Total }
func (c YGOProductREST) GetRarityStats() map[string]int { return c.RarityStats }
func (c YGOProductREST) GetContent() []ProductContent   { return c.Content }

type ProductContent struct {
	Card            YGOProduct `json:"card"`
	ProductPosition string     `json:"productPosition"`
	Rarities        []string   `json:"rarities"`
}
