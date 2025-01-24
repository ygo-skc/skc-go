package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ygo-skc/skc-go/skc-db/util"
)

const (
	minPoolSize = 20
	maxPoolSize = 30
)

// Connect to SKC database.
func EstablishDBConn() {
	uri := "%s:%s@tcp(%s)/%s"
	dataSourceName := fmt.Sprintf(uri, util.EnvMap["SKC_DB_USERNAME"], util.EnvMap["SKC_DB_PASSWORD"], util.EnvMap["SKC_DB_HOST"], util.EnvMap["SKC_DB_NAME"])

	var err error
	if skcDBConn, err = sql.Open("mysql", dataSourceName); err != nil {
		log.Fatalln("Error occurred while trying to establish DB connection: ", err)
	}

	skcDBConn.SetMaxIdleConns(minPoolSize)
	skcDBConn.SetMaxOpenConns(maxPoolSize)
}
