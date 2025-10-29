package client

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"time"

	"github.com/ygo-skc/skc-go/common/v2/health"
	"github.com/ygo-skc/skc-go/common/v2/ygo"
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
	slog.Info(fmt.Sprintf("Creating Card Service gRPC Client using SSL Server Name %s and Host %s",
		sslServerName,
		serviceHost,
	))

	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: false,
		ServerName:         sslServerName,
	})

	conn, err := grpc.NewClient(serviceHost,
		grpc.WithTransportCredentials(creds),
		grpc.WithDefaultCallOptions(
			grpc.UseCompressor("gzip"),
		),
		grpc.WithDefaultServiceConfig(`{
			"methodConfig": [{
				"name": [{"service": "ygo.CardService"}, {"service": "ygo.ProductService"}],
				"retryPolicy": {
					"MaxAttempts": 3,
					"InitialBackoff": "0.1s",
					"MaxBackoff": "1s",
					"BackoffMultiplier": 2.0,
					"RetryableStatusCodes": ["UNAVAILABLE", "DEADLINE_EXCEEDED"]
				}
			}]
		}`),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             1 * time.Second,
			PermitWithoutStream: true,
		}))

	if err != nil {
		return nil, err
	}

	return newYGOClientImpV1(conn), nil
}
