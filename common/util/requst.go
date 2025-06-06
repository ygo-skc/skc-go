package util

import (
	"context"
	"log/slog"
)

func InitRequest(ctx context.Context, apiName string, flow string, customAttributes ...slog.Attr) (*slog.Logger, context.Context) {
	return NewLogger(ContextWithMetadata(ctx, apiName, flow), flow, customAttributes...)
}
