package util

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

const (
	clientIDMetaName = "client-id"
	flowMetaName     = "flow"
	traceMetaName    = "trace-id"
)

func ContextWithMetadata(ctx context.Context, clientID string, flow string) context.Context {
	md := metadata.New(map[string]string{
		clientIDMetaName: clientID,
		flowMetaName:     flow,
		traceMetaName:    uuid.New().String(),
	})
	return metadata.NewOutgoingContext(ctx, md)
}
