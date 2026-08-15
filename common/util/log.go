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

type loggerKeyType struct{}

var loggerKey = loggerKeyType{}

const (
	traceIDKey         = "trace_id"
	spanIDKey          = "span_id"
	flowKey            = "app.flow"
	originatingFlowKey = "app.origin"
	clientIDKey        = "client_id"
)

func RetrieveLogger(ctx context.Context) *slog.Logger {
	l := ctx.Value(loggerKey)
	if l == nil {
		slog.Warn("Using default slog as context does not have logger info")
		return slog.Default()
	}
	return l.(*slog.Logger)
}

func NewLogger(ctx context.Context, flow string, customAttributes ...slog.Attr) (*slog.Logger, context.Context) {
	traceID := traceFromContext(ctx)
	var originatingFlow, clientID string

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if flow := md.Get(flowMetaName); len(flow) > 0 && flow[0] != "" {
			originatingFlow = flow[0]
		}
		if id := md.Get(clientIDMetaName); len(id) > 0 && id[0] != "" {
			clientID = id[0]
		}
		if trace := md.Get(traceMetaName); len(trace) > 0 && trace[0] != "" {
			traceID = trace[0] // override trace if present
		}
	}

	defaults := []any{
		slog.String(traceIDKey, traceID),
		slog.String(spanIDKey, uuid.New().String()),
		slog.String(flowKey, flow),
	}
	if originatingFlow != "" {
		defaults = append(defaults, slog.String(originatingFlowKey, originatingFlow))
	}
	if clientID != "" {
		defaults = append(defaults, slog.String(clientIDKey, clientID))
	}
	for _, customAttribute := range customAttributes {
		defaults = append(defaults, customAttribute)
	}

	l := slog.With(defaults...)
	return l, context.WithValue(ctx, loggerKey, l)
}

func AddLoggerAttribute(ctx context.Context, customAttributes ...slog.Attr) (*slog.Logger, context.Context) {
	newAttributes := make([]any, len(customAttributes))
	for i, customAttribute := range customAttributes {
		newAttributes[i] = customAttribute
	}

	l := RetrieveLogger(ctx).With(newAttributes...)
	return l, context.WithValue(ctx, loggerKey, l)
}
