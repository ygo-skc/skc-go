package client

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"time"

	"github.com/ygo-skc/skc-go/common/ygo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

func CreateCardServiceClient(sslServerName string, serviceHost string) (*ygo.CardServiceClient, error) {
	slog.Info(fmt.Sprintf("Creating Card Service gRPC Client using SSL Server Name %s and Host %s",
		sslServerName,
		serviceHost))

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
				"name": [{"service": "ygo.CardService"}],
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

	CardServiceClient := ygo.NewCardServiceClient(conn)
	return &CardServiceClient, nil
}
