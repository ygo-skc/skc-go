module github.com/ygo-skc/skc-go/ygo-service

go 1.25

require (
	github.com/go-sql-driver/mysql v1.9.3
	github.com/ygo-skc/skc-go/common v1.5.1
	google.golang.org/grpc v1.75.0
	google.golang.org/protobuf v1.36.8
)

replace github.com/ygo-skc/skc-go/common => ../common

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
)
