package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/ygo-skc/skc-go/common/v2/util"
)

var (
	skcDBConn *sql.DB
)

const (
	maxPoolSize = 150
)

// Connect to SKC database.
func EstablishDBConn() {
	uri := "%s:%s@tcp(%s)/%s"
	dataSourceName := fmt.Sprintf(uri, util.EnvMap["SKC_DB_USERNAME"], util.EnvMap["SKC_DB_PASSWORD"], util.EnvMap["SKC_DB_HOST"],
		util.EnvMap["SKC_DB_NAME"])

	var err error
	if skcDBConn, err = sql.Open("mysql", dataSourceName); err != nil {
		slog.Error("Failed to establish DB connection", slog.Any("err", err))
		os.Exit(1)
	}

	skcDBConn.SetMaxOpenConns(maxPoolSize)
	skcDBConn.SetConnMaxLifetime(1 * time.Hour)
	skcDBConn.SetConnMaxIdleTime(30 * time.Minute)
}
