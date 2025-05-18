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

func (imp YGOCardRepository) GetCardByID(ctx context.Context, cardID string) (*ygo.Card, *model.APIError) {
	logger := cUtil.LoggerFromContext(ctx)

	args := make([]interface{}, 1)
	args[0] = cardID
	c, err := queryCard(logger, cardByCardIDQuery, args)
	if err != nil && err.StatusCode == http.StatusNotFound {
		logger.Info("Card ID is not valid")
	} else {
		logger.Info("Card ID is valid")
	}
	return c, err
}

func (imp YGOCardRepository) GetCardsByIDs(
	ctx context.Context, cardIDs []string,
) (*ygo.Cards, *model.APIError) {
	logger := cUtil.LoggerFromContext(ctx)

	args, numCards := buildVariableQuerySubjects(cardIDs)

	query := fmt.Sprintf(cardsByCardIDsQuery, variablePlaceholders(numCards))
	if rows, err := skcDBConn.Query(query, args...); err != nil {
		return nil, handleQueryError(logger, err)
	} else {
		if cards, err := parseRowsForCards(ctx, rows); err != nil {
			return nil, err
		} else {
			return &ygo.Cards{
				CardInfo:         *cards,
				UnknownResources: model.FindMissingIDs(*cards, cardIDs),
			}, nil
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
		query = randomCardQuery
	} else {
		args, _ = buildVariableQuerySubjects(blacklistedCards)
		query = fmt.Sprintf(randomCardWithBlackListQuery, variablePlaceholders(numBlackListed))
	}

	c, err := queryCard(logger, query, args)
	logger.Info(fmt.Sprintf("Random card determined to be; ID: %s, Name: %s", c.ID, c.Name))
	return c, err
}
