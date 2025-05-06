package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ygo-skc/skc-go/common/util"
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
	skcDBConn.SetConnMaxIdleTime(10 * time.Minute)
}
