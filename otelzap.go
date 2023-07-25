package otelzap

import (
	otel "github.com/agoda-com/opentelemetry-logs-go/logs"
	"go.uber.org/zap/zapcore"
)

// NewOtelCore creates new OpenTelemetry Core to export logs in OTLP format
func NewOtelCore(loggerProvider otel.LoggerProvider) zapcore.Core {

	logger := loggerProvider.Logger(
		instrumentationScope.Name,
		otel.WithInstrumentationVersion(instrumentationScope.Version),
	)

	return otlpCore{
		logger: logger,
	}
}
