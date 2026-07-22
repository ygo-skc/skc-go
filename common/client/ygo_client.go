package client

import (
	"crypto/tls"
	"log/slog"
	"time"

	"github.com/ygo-skc/skc-go/common/v3/health"
	"github.com/ygo-skc/skc-go/common/v3/ygo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type YGOClientImpV1 struct {
	CardService    YGOCardClientImp
	ProductService YGOProductClientImp
	HealthService  YGOHealthClientImp
}

func newYGOClientImpV1(conn *grpc.ClientConn) *YGOClientImpV1 {
	return &YGOClientImpV1{
		CardService:    &YGOCardClientImpV1{client: ygo.NewCardServiceClient(conn)},
		ProductService: &YGOProductClientImpV1{client: ygo.NewProductServiceClient(conn)},
		HealthService:  &YGOHealthClientImpV1{client: health.NewHealthServiceClient(conn)},
	}
}

func NewYGOServiceClients(sslServerName string, serviceHost string) (*YGOClientImpV1, error) {
	slog.Info("Creating YGO service gRPC client",
		slog.String("ssl_server_name", sslServerName),
		slog.String("service_host", serviceHost),
	)

	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: false,
		ServerName:         sslServerName,
	})

	conn, err := grpc.NewClient(serviceHost,
		grpc.WithTransportCredentials(creds),

		// below fields can be high since im on same docker network. But since this is a shared lib, change them if that ever changes
		grpc.WithReadBufferSize(256<<10),
		grpc.WithWriteBufferSize(256<<10),
		grpc.WithInitialWindowSize(256<<10),
		grpc.WithInitialConnWindowSize(4*256*1024), // if above 3 change, change the middle value here too

		grpc.WithDefaultCallOptions(
			grpc.UseCompressor("gzip"),
		),
		grpc.WithDefaultServiceConfig(`{
			"methodConfig": [{
				"name": [{"service": ""}],
				"timeout": "6s",
				"retryPolicy": {
					"MaxAttempts": 3,
					"InitialBackoff": "0.1s",
					"MaxBackoff": "1s",
					"BackoffMultiplier": 2.0,
					"RetryableStatusCodes": ["UNKNOWN", "DEADLINE_EXCEEDED", "DATA_LOSS", "UNAVAILABLE"]
				}
			}]
		}`),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             1 * time.Second,
			PermitWithoutStream: false,
		}))

	if err != nil {
		return nil, err
	}

	return newYGOClientImpV1(conn), nil
}
