package zlogger

import (
	"context"

	"go.uber.org/zap"
)

var (
	// CtxLoggerObject reference to the zap.Logger
	CtxLoggerObject = struct{ s string }{"logger"}
)

// FromContext logger object
func FromContext(ctx context.Context) *zap.Logger {
	if logObj := ctx.Value(CtxLoggerObject); logObj != nil {
		return logObj.(*zap.Logger)
	}
	return zap.L()
}

// WithLogger puts logger to context
func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, CtxLoggerObject, logger)
}
