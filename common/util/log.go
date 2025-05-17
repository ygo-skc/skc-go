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
)

func LoggerFromContext(ctx context.Context) *slog.Logger {
	if l := ctx.Value(loggerKey); l == nil {
		slog.Warn("Using default slog as context does not have logger info")
		return slog.Default()
	} else {
		return l.(*slog.Logger)
	}
}

func NewRequestSetup(ctx context.Context, operation string, customAttributes ...slog.Attr) (*slog.Logger, context.Context) {
	defaults := []any{slog.String("requestID", uuid.New().String()), slog.String("operation", operation)}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if clientID := md.Get(clientIDMetaName); len(clientID) > 0 && clientID[0] != "" {
			defaults = append(defaults, slog.String("clientID", clientID[0]))
		}
		if flow := md.Get(flowMetaName); len(flow) > 0 && flow[0] != "" {
			defaults = append(defaults, slog.String("flow", flow[0]))
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
