package otelzap

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a thin wrapper for zap.Logger that adds Ctx method.
type Logger struct {
	*zap.Logger
}

const contextKey = "context"

type SugaredLogger struct {
	*zap.SugaredLogger
}

func (l *Logger) Sugar() *SugaredLogger {
	return &SugaredLogger{
		SugaredLogger: l.Logger.Sugar(),
	}
}

func (l *Logger) Ctx(ctx context.Context) *Logger {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return l.With(zap.Reflect(contextKey, span.SpanContext()))
	}
	return l
}

func (l *SugaredLogger) Ctx(ctx context.Context) *SugaredLogger {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return &SugaredLogger{
			SugaredLogger: l.With(zap.Reflect(contextKey, span.SpanContext())),
		}
	}
	return &SugaredLogger{
		SugaredLogger: l.SugaredLogger,
	}
}

func (l *Logger) With(fields ...zapcore.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(fields...),
	}
}
