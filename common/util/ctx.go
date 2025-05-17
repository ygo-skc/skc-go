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

	traceCtxKey = "traceKey"
)

func traceFromContext(ctx context.Context) string {
	if t := ctx.Value(traceCtxKey); t == nil {
		return uuid.New().String()
	} else {
		return t.(string)
	}
}

func ContextWithMetadata(ctx context.Context, clientID string, flow string) context.Context {
	traceID := traceFromContext(ctx)

	md := metadata.New(map[string]string{
		clientIDMetaName: clientID,
		flowMetaName:     flow,
		traceMetaName:    traceID,
	})
	return context.WithValue(metadata.NewOutgoingContext(ctx, md), traceCtxKey, traceID)
}
