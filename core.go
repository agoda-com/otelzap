package otelzap

import (
	otel "github.com/agoda-com/opentelemetry-logs-go/logs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap/zapcore"
)

const (
	instrumentationName = "github.com/agoda-com/otelzap"
)

// This class provide interface for OTLP logger
type otlpCore struct {
	logger otel.Logger

	fields []zapcore.Field
}

var instrumentationScope = instrumentation.Scope{
	Name:      instrumentationName,
	Version:   Version(),
	SchemaURL: semconv.SchemaURL,
}

func (otlpCore) Enabled(zapcore.Level) bool {
	return true
}
func (c otlpCore) With(f []zapcore.Field) zapcore.Core {

	fields := c.fields

	for _, fld := range f {
		fields = append(fields, fld)
	}

	return otlpCore{
		logger: c.logger,
		fields: fields,
	}
}

// Check OTLP zap extension method to check if logger is enabled
func (c otlpCore) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return checked.AddCore(entry, c)
	}
	return checked
}

func (c otlpCore) Sync() error {
	return nil
}

func (c otlpCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {

	var attributes []attribute.KeyValue
	var spanCtx *trace.SpanContext

	// add common zap log fields as attributes
	for _, s := range c.fields {
		if s.Key == "context" {
			if ctxValue, ok := s.Interface.(trace.SpanContext); ok {
				spanCtx = &ctxValue
			}
		} else {
			attributes = append(attributes, otelAttribute(s)...)
		}
	}
	// add zap log fields as attributes
	for _, s := range fields {
		attributes = append(attributes, otelAttribute(s)...)
	}

	if ent.Level > zapcore.InfoLevel {
		callerString := ent.Caller.String()

		if len(callerString) > 0 {
			attributes = append(attributes, semconv.ExceptionType(callerString))
		}

		if len(ent.Stack) > 0 {
			attributes = append(attributes, semconv.ExceptionStacktrace(ent.Stack))
		}
	}

	severityString := ent.Level.String()
	severity := otelLevel(ent.Level)

	var traceID *trace.TraceID = nil
	var spanID *trace.SpanID = nil
	var traceFlags *trace.TraceFlags = nil
	if spanCtx != nil {
		tid := spanCtx.TraceID()
		sid := spanCtx.SpanID()
		tf := spanCtx.TraceFlags()
		traceID = &tid
		spanID = &sid
		traceFlags = &tf
	}

	lrc := otel.LogRecordConfig{
		Timestamp:            &ent.Time,
		ObservedTimestamp:    ent.Time,
		TraceId:              traceID,
		SpanId:               spanID,
		TraceFlags:           traceFlags,
		SeverityText:         &severityString,
		SeverityNumber:       &severity,
		Body:                 &ent.Message,
		Resource:             nil,
		InstrumentationScope: &instrumentationScope,
		Attributes:           &attributes,
	}

	r := otel.NewLogRecord(lrc)

	c.logger.Emit(r)

	return nil
}
