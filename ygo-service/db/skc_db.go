package db

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ygo-skc/skc-go/common/model"
	cUtil "github.com/ygo-skc/skc-go/common/util"
)

type CardRepository interface {
	GetDBVersion(context.Context) (string, error)

	GetCardByID(context.Context, string) (model.YGOCardREST, *model.APIError)
	GetCardsByIDs(context.Context, []string) (model.BatchCardData[model.CardIDs], *model.APIError)
}
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

// Leverages GetCardsByIDs to get information on a specific card using its identifier
func (imp YGOCardRepository) GetCardByID(ctx context.Context, cardID string) (model.YGOCardREST, *model.APIError) {
	logger := cUtil.LoggerFromContext(ctx)

	if results, err := imp.GetCardsByIDs(ctx, []string{cardID}); err != nil {
		return model.YGOCardREST{}, err
	} else {
		if card, exists := results.CardInfo[cardID]; !exists {
			logger.Warn("Card ID isn't valid")
			return model.YGOCardREST{}, &model.APIError{Message: fmt.Sprintf("No results found when querying by card ID %s", cardID), StatusCode: http.StatusNotFound}
		} else {
			logger.Info("Card ID is valid")
			return card.(model.YGOCardREST), nil
		}
	}
}

func (imp YGOCardRepository) GetCardsByIDs(ctx context.Context, cardIDs []string) (model.BatchCardData[model.CardIDs], *model.APIError) {
	logger := cUtil.LoggerFromContext(ctx)

	args, numCards := buildVariableQuerySubjects(cardIDs)
	cardData := make(model.CardDataMap, numCards) // used to store results

	query := fmt.Sprintf(cardsByCardIDsQuery, variablePlaceholders(numCards))

	if rows, err := skcDBConn.Query(query, args...); err != nil {
		return model.BatchCardData[model.CardIDs]{}, handleQueryError(logger, err)
	} else {
		if cards, err := parseRowsForCard(ctx, rows); err != nil {
			return model.BatchCardData[model.CardIDs]{}, err
		} else {
			for _, card := range cards {
				cardData[card.ID] = card
			}
		}
	}

	cd := model.BatchCardData[model.CardIDs]{CardInfo: cardData, UnknownResources: cardData.FindMissingIDs(cardIDs)}
	logger.Info(fmt.Sprintf("The following Card ID's were invalid: %v", cd.UnknownResources))
	return cd, nil
}
