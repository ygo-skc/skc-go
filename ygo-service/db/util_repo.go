package db

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/ygo-skc/skc-go/common/v3/model"
	"github.com/ygo-skc/skc-go/common/v3/util"
)

type UtilRepository interface {
	GetDBVersion(context.Context) (string, error)
}
type YGOUtilRepository struct{}

// Get version of MYSQL being used by SKC DB.
func (imp YGOCardRepository) GetDBVersion(ctx context.Context) (string, error) {
	var version string
	if err := skcDBConn.QueryRow(dbVersionQuery).Scan(&version); err != nil {
		util.RetrieveLogger(ctx).Error("Error getting SKC DB version", slog.Any("err", err))
		return version, &model.APIError{Message: genericError, StatusCode: http.StatusInternalServerError}
	}

	return version, nil
}
