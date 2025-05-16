module github.com/ygo-skc/skc-go/ygo-service

go 1.24.2

require (
	github.com/go-sql-driver/mysql v1.9.2
	github.com/ygo-skc/skc-go/common v0.0.0
	google.golang.org/grpc v1.72.1
	google.golang.org/protobuf v1.36.6
)

replace github.com/ygo-skc/skc-go/common => ../common

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
)
