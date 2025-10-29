package model

import (
	"github.com/ygo-skc/skc-go/common/v2/ygo"
)

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

func BatchCardDataFromProto[T CardIDs | CardNames](c *ygo.Cards, keyFn func(*ygo.Card) string) *BatchCardData[T] {
	batchCardData := make(CardDataMap, len(c.CardInfo))
	for _, v := range c.CardInfo {
		batchCardData[keyFn(v)] = YGOCardRESTFromProto(v)
	}
	return &BatchCardData[T]{CardInfo: batchCardData, UnknownResources: c.UnknownResources}
}

func BatchCardDataFromProductProto[T CardIDs | CardNames](p *ygo.Product, keyFn func(*ygo.Card) string) *BatchCardData[T] {
	batchCardData := make(CardDataMap, len(p.Items))
	for _, item := range p.Items {
		batchCardData[keyFn(item.Card)] = YGOCardRESTFromProto(item.Card)
	}
	return &BatchCardData[T]{CardInfo: batchCardData, UnknownResources: make([]string, 0)}
}

func BatchProductSummaryFromProductsProto[T ProductIDs](p *ygo.Products, keyFn func(*ygo.ProductSummary) string) *BatchProductSummaryData[T] {
	batchProductInfo := make(ProductSummaryDataMap, len(p.Products))
	for _, product := range p.Products {
		batchProductInfo[keyFn(product)] = YGOProductSummaryREST{
			ID:          product.ID,
			Locale:      product.Locale,
			Name:        product.Name,
			Type:        product.Type,
			SubType:     product.SubType,
			ReleaseDate: product.ReleaseDate,
		}
	}
	return &BatchProductSummaryData[T]{ProductInfo: batchProductInfo, UnknownResources: p.UnknownResources}
}
