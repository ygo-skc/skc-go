package db

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/ygo-skc/skc-go/common/v3/util"
)

var (
	skcDBConn *sql.DB
)

const (
	maxPoolSize = 150
	maxIdle     = 50
)

// Connect to SKC database.
func EstablishDBConn() {
	config := mysql.NewConfig()
	config.User = util.EnvMap["SKC_DB_USERNAME"]
	config.Passwd = util.EnvMap["SKC_DB_PASSWORD"]
	config.Net = "tcp"
	config.Addr = util.EnvMap["SKC_DB_HOST"]
	config.DBName = util.EnvMap["SKC_DB_NAME"]

	config.Timeout = 2 * time.Second
	config.ReadTimeout = 2 * time.Second
	config.WriteTimeout = 2 * time.Second
	config.InterpolateParams = true

	connector, err := mysql.NewConnector(config)
	if err != nil {
		slog.Error("Failed to configure DB connection", slog.Any("err", err))
		os.Exit(1)
	}
	skcDBConn = sql.OpenDB(connector)

	skcDBConn.SetMaxOpenConns(maxPoolSize)
	skcDBConn.SetConnMaxLifetime(20 * time.Minute)

	skcDBConn.SetMaxIdleConns(maxIdle)
	skcDBConn.SetConnMaxIdleTime(10 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := skcDBConn.PingContext(ctx); err != nil {
		slog.Error("Failed to ping DB", slog.Any("err", err))
		os.Exit(1)
	}
}
