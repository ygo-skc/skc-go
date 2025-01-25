package main

import (
	"os"
	"strings"

	"github.com/ygo-skc/skc-go/skc-db/db"

	_ "github.com/go-sql-driver/mysql"
	cUtil "github.com/ygo-skc/skc-go/common/util"
)

const (
	ENV_VARIABLE_NAME string = "SKC_API_GO_DOT_ENV_FILE"
)

func init() {
	isCICD := os.Getenv("IS_CICD")
	if isCICD != "true" && !strings.HasSuffix(os.Args[0], ".test") {
		cUtil.ConfigureEnv(ENV_VARIABLE_NAME)
	}
}

func main() {
	db.EstablishDBConn()
	go listen()
	select {}
}
