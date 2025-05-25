package db

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ygo-skc/skc-go/common/model"
	cUtil "github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/common/ygo"
)

type YGOCardRepository struct{}

// Get version of MYSQL being used by SKC DB.
func (imp YGOCardRepository) GetDBVersion(ctx context.Context) (string, error) {
	var version string
	if err := skcDBConn.QueryRow(dbVersionQuery).Scan(&version); err != nil {
		cUtil.LoggerFromContext(ctx).Error(fmt.Sprintf("Error getting SKC DB version - %v", err))
		return version, &model.APIError{Message: genericError, StatusCode: http.StatusInternalServerError}
	}

	return version, nil
}

// Get IDs for all card colors currently in database.
func (imp YGOCardRepository) GetCardColorIDs(ctx context.Context) (*ygo.CardColors, *model.APIError) {
	logger := cUtil.LoggerFromContext(ctx)
	logger.Info("Retrieving card color IDs from DB")
	cardColorIDs := make(map[string]uint32)

	if rows, err := skcDBConn.Query(cardColorIDsQuery); err != nil {
		return nil, handleQueryError(logger, err)
	} else {
		for rows.Next() {
			var colorId uint32
			var cardColor string

			if err := rows.Scan(&colorId, &cardColor); err != nil {
				return nil, handleRowParsingError(logger, err)
			}

			cardColorIDs[cardColor] = colorId
		}
		return &ygo.CardColors{Values: cardColorIDs}, nil
	}
}

func (imp YGOCardRepository) GetCardByID(ctx context.Context, cardID string) (*ygo.Card, *model.APIError) {
	logger := cUtil.LoggerFromContext(ctx)

	args := make([]interface{}, 1)
	args[0] = cardID
	query := fmt.Sprintf(cardByCardIDQuery, cardAttributes)

	c, err := queryCard(logger, query, args)
	if err != nil && err.StatusCode == http.StatusNotFound {
		logger.Info("Card ID is not valid")
	} else {
		logger.Info("Card ID is valid")
	}
	return c, err
}

func (imp YGOCardRepository) GetCardsByIDs(
	ctx context.Context, cardIDs model.CardIDs,
) (*ygo.Cards, *model.APIError) {
	logger := cUtil.LoggerFromContext(ctx)

	args, numCards := buildVariableQuerySubjects(cardIDs)
	query := fmt.Sprintf(cardsByCardIDsQuery, cardAttributes, variablePlaceholders(numCards))

	if rows, err := skcDBConn.Query(query, args...); err != nil {
		return nil, handleQueryError(logger, err)
	} else {
		if cards, err := parseRowsForCards(ctx, rows, model.CardIDAsKey); err != nil {
			return nil, err
		} else {
			return &ygo.Cards{
				CardInfo:         cards,
				UnknownResources: model.FindMissingKeys(cards, cardIDs),
			}, nil
		}
	}
}

// Uses card names to find instance of card
func (imp YGOCardRepository) GetCardsByNames(ctx context.Context, cardNames model.CardNames) (*ygo.Cards, *model.APIError) {
	logger := cUtil.LoggerFromContext(ctx)

	args, numCards := buildVariableQuerySubjects(cardNames)
	query := fmt.Sprintf(cardsByCardNamesQuery, cardAttributes, variablePlaceholders(numCards))

	if rows, err := skcDBConn.Query(query, args...); err != nil {
		return nil, handleQueryError(logger, err)
	} else {
		if cards, err := parseRowsForCards(ctx, rows, model.CardNameAsKey); err != nil {
			return nil, err
		} else {
			return &ygo.Cards{
				CardInfo:         cards,
				UnknownResources: model.FindMissingKeys(cards, cardNames),
			}, nil
		}
	}
}
func (imp YGOCardRepository) GetArchetypalCardsUsingCardName(ctx context.Context, archetypeName string) (*ygo.CardList, *model.APIError) {
	logger := cUtil.LoggerFromContext(ctx)
	logger.Info(fmt.Sprintf("Retrieving card data from DB for all cards that reference archetype %s in their name", archetypeName))
	searchTerm := `%` + archetypeName + `%`

	query := fmt.Sprintf(archetypalCardsUsingCardNameQuery, cardAttributes)
	if rows, err := skcDBConn.Query(query, searchTerm); err != nil {
		return nil, handleQueryError(logger, err)
	} else {
		if cards, err := parseRowsForCardList(ctx, rows); err != nil {
			return nil, err
		} else {
			return &ygo.CardList{Cards: cards}, err
		}
	}
}

func (imp YGOCardRepository) GetExplicitArchetypalInclusions(ctx context.Context, archetypeName string) (*ygo.CardList, *model.APIError) {
	logger := cUtil.LoggerFromContext(ctx)
	logger.Info(fmt.Sprintf("Retrieving cards that are explicitly considered part of archetype %s", archetypeName))

	subQuery := fmt.Sprintf(archetypeInclusionSubQuery, cardAttributes, archetypeName)
	query := fmt.Sprintf(archetypalCardsUsingCardTextQuery, subQuery, archetypeName)
	if rows, err := skcDBConn.Query(query); err != nil {
		return nil, handleQueryError(logger, err)
	} else {
		if cards, err := parseRowsForCardList(ctx, rows); err != nil {
			return nil, err
		} else {
			return &ygo.CardList{Cards: cards}, err
		}
	}
}
func (imp YGOCardRepository) GetExplicitArchetypalExclusions(ctx context.Context, archetypeName string) (*ygo.CardList, *model.APIError) {
	logger := cUtil.LoggerFromContext(ctx)
	logger.Info(fmt.Sprintf("Retrieving cards that are explicitly NOT considered part of archetype %s", archetypeName))

	subQuery := fmt.Sprintf(archetypeExclusionSubQuery, cardAttributes, archetypeName)
	query := fmt.Sprintf(nonArchetypalCardsUsingCardTextQuery, subQuery, archetypeName)
	if rows, err := skcDBConn.Query(query); err != nil {
		return nil, handleQueryError(logger, err)
	} else {
		if cards, err := parseRowsForCardList(ctx, rows); err != nil {
			return nil, err
		} else {
			return &ygo.CardList{Cards: cards}, err
		}
	}
}

func (imp YGOCardRepository) GetRandomCard(
	ctx context.Context,
	blacklistedCards []string,
) (*ygo.Card, *model.APIError) {
	logger := cUtil.LoggerFromContext(ctx)

	// pick correct query based on contents of blacklistedCards
	numBlackListed := len(blacklistedCards)
	var query string
	var args []interface{}
	if numBlackListed == 0 {
		query = fmt.Sprintf(randomCardQuery, cardAttributes)
	} else {
		args, _ = buildVariableQuerySubjects(blacklistedCards)
		query = fmt.Sprintf(randomCardWithBlackListQuery, cardAttributes, variablePlaceholders(numBlackListed))
	}

	c, err := queryCard(logger, query, args)
	logger.Info(fmt.Sprintf("Random card determined to be; ID: %s, Name: %s", c.ID, c.Name))
	return c, err
}
