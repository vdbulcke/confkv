package logger

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	traceIdKey = TraceID("trace_id")
	jobIdKey   = TraceID("job_id")
)

type TraceID string
type JobID string

func GetTraceID(ctx context.Context) string {
	val := ""
	if v := ctx.Value(traceIdKey); v != nil {
		return v.(string)
	}

	return val
}

func GetJobID(ctx context.Context) string {
	val := ""
	if v := ctx.Value(jobIdKey); v != nil {
		return v.(string)
	}

	return val
}

func ContextWithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIdKey, traceID)
}

func ContextWithJobID(ctx context.Context) context.Context {
	JobID := NewTraceID()
	return context.WithValue(ctx, jobIdKey, JobID)
}

func GetZapTraceID(ctx context.Context) zap.Field {
	id := GetTraceID(ctx)
	if id == "" {
		return zap.Skip()
	}
	return zap.String("trace_id", id)

}
func GetZapJobID(ctx context.Context) zap.Field {
	id := GetJobID(ctx)
	if id == "" {
		return zap.Skip()
	}
	return zap.String("job_id", id)
}

func GetLoggerWithContext(ctx context.Context, l *zap.Logger) *zap.Logger {

	newLogger := l.With(GetZapTraceID(ctx), GetZapJobID(ctx))

	return newLogger
}

func NewTraceID() string {
	return uuid.New().String()
}
