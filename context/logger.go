package context

import (
	"context"
	"log"
)

const (
	loggerKey key = "logger"
)

func WithLogger(ctx context.Context, logger *log.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func Logger(ctx context.Context) *log.Logger {
	logger, ok := ctx.Value(loggerKey).(*log.Logger)
	if !ok {
		return nil
	}
	return logger
}
