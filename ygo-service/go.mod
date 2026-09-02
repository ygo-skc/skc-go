module github.com/ygo-skc/skc-go/ygo-service

go 1.27

require (
	github.com/go-sql-driver/mysql v1.10.1
	github.com/ygo-skc/skc-go/common/v3 v3.3.1
	google.golang.org/grpc v1.83.2
)

replace github.com/ygo-skc/skc-go/common/v3 => ../common

require (
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260526163538-3dc84a4a5aaa // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)
