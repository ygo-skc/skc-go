package db

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/ygo-skc/skc-go/common/util"
)

var (
	skcDBConn  *sql.DB
	spaceRegex = regexp.MustCompile(`[ ]+`)
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
		log.Fatalln("Error occurred while trying to establish DB connection: ", err)
	}

	skcDBConn.SetMaxOpenConns(maxPoolSize)
	skcDBConn.SetConnMaxLifetime(1 * time.Hour)
	skcDBConn.SetConnMaxIdleTime(30 * time.Minute)
}
