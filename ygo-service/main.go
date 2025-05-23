package main

import (
	"os"
	"strings"

	"github.com/ygo-skc/skc-go/ygo-service/api"
	"github.com/ygo-skc/skc-go/ygo-service/db"

	_ "github.com/go-sql-driver/mysql"
	cUtil "github.com/ygo-skc/skc-go/common/util"
	_ "google.golang.org/grpc/encoding/gzip"
)

const (
	ENV_VARIABLE_NAME string = "YGO_SERVICE_DOT_ENV_FILE"
)

func init() {
	isCICD := os.Getenv("IS_CICD")
	if isCICD != "true" && !strings.HasSuffix(os.Args[0], ".test") {
		cUtil.ConfigureEnv(ENV_VARIABLE_NAME)
	}
}

func main() {
	db.EstablishDBConn()
	go api.RunService()
	select {}
}
