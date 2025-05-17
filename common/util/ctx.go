package util

import (
	"context"

	"google.golang.org/grpc/metadata"
)

const (
	clientIDMetaName = "client-id"
	flowMetaName     = "flow"
)

func ContextWithMetadata(ctx context.Context, clientID string, flow string) context.Context {
	md := metadata.New(map[string]string{
		clientIDMetaName: clientID,
		flowMetaName:     flow,
	})
	return metadata.NewOutgoingContext(ctx, md)
}
