package db

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ygo-skc/skc-go/common/v2/model"
	"github.com/ygo-skc/skc-go/common/v2/util"
	"github.com/ygo-skc/skc-go/common/v2/ygo"
	"google.golang.org/grpc/status"
)

type UtilRepository interface {
	GetCardColorIDs(context.Context) (*ygo.CardColors, *status.Status)
}
type YGOUtilRepository struct{}

// Get version of MYSQL being used by SKC DB.
func (imp YGOCardRepository) GetDBVersion(ctx context.Context) (string, error) {
	var version string
	if err := skcDBConn.QueryRow(dbVersionQuery).Scan(&version); err != nil {
		util.RetrieveLogger(ctx).Error(fmt.Sprintf("Error getting SKC DB version - %v", err))
		return version, &model.APIError{Message: genericError, StatusCode: http.StatusInternalServerError}
	}

	return version, nil
}
