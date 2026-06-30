module github.com/ygo-skc/skc-go/ygo-service

go 1.26

require (
	github.com/go-sql-driver/mysql v1.10.0
	github.com/ygo-skc/skc-go/common/v2 v2.1.6
	google.golang.org/grpc v1.82.0
	google.golang.org/protobuf v1.36.11
)

replace github.com/ygo-skc/skc-go/common/v2 => ../common

require (
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260414002931-afd174a4e478 // indirect
)
