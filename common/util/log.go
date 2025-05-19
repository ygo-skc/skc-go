package util

import (
	"context"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

func init() {
	slogOpts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, slogOpts)))
}

type contextKey string

const (
	loggerKey contextKey = "logger"

	traceIDKey         = "trace_id"
	spanIDKey          = "span_id"
	flowKey            = "app.flow"
	originatingFlowKey = "app.origin"
	clientIDKey        = "client.id"
)

func LoggerFromContext(ctx context.Context) *slog.Logger {
	if l := ctx.Value(loggerKey); l == nil {
		slog.Warn("Using default slog as context does not have logger info")
		return slog.Default()
	} else {
		return l.(*slog.Logger)
	}
}

func NewRequestSetup(ctx context.Context, flow string, customAttributes ...slog.Attr) (*slog.Logger, context.Context) {
	defaults := []any{
		slog.String(traceIDKey, traceFromContext(ctx)),
		slog.String(spanIDKey, uuid.New().String()),
		slog.String(flowKey, flow),
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if flow := md.Get(flowMetaName); len(flow) > 0 && flow[0] != "" {
			defaults = append(defaults, slog.String(originatingFlowKey, flow[0]))
		}
		if clientID := md.Get(clientIDMetaName); len(clientID) > 0 && clientID[0] != "" {
			defaults = append(defaults, slog.String(clientIDKey, clientID[0]))
		}
		if traceID := md.Get(traceMetaName); len(traceID) > 0 && traceID[0] != "" {
			defaults = append(defaults, slog.String(traceIDKey, traceID[0])) // overrides default
		}
	}

	for _, customAttribute := range customAttributes {
		defaults = append(defaults, customAttribute)
	}

	l := slog.With(defaults...)
	return l, context.WithValue(ctx, loggerKey, l)
}

func AddAttribute(ctx context.Context, customAttributes ...slog.Attr) (*slog.Logger, context.Context) {
	newAttributes := []any{}

	for _, customAttribute := range customAttributes {
		newAttributes = append(newAttributes, customAttribute)
	}

	l := LoggerFromContext(ctx)
	l = l.With(newAttributes...)
	return l, context.WithValue(ctx, loggerKey, l)
}
