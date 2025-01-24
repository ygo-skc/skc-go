package main

import (
	"github.com/ygo-skc/skc-go/skc-db/db"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db.EstablishDBConn()
	go listen()
	select {}
}
