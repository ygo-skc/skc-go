module github.com/ygo-skc/skc-go/skc-db-service

go 1.23.4

require (
	github.com/go-sql-driver/mysql v1.8.1
	github.com/ygo-skc/skc-go/common v0.0.0
	google.golang.org/grpc v1.70.0
	google.golang.org/protobuf v1.36.4
)

replace github.com/ygo-skc/skc-go/common => ../common

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241202173237-19429a94021a // indirect
)
