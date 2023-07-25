package otelzap

import (
	"context"
	"go.uber.org/zap"
)

// L returns the global Logger
func L() *Logger {
	return &Logger{
		zap.L(),
	}
}

func S() *SugaredLogger {
	return L().Sugar()
}

// Ctx is a shortcut for L().Ctx(ctx).
func Ctx(ctx context.Context) *Logger {
	return L().Ctx(ctx)
}
