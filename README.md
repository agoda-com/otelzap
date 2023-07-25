# otelzap
Zap logger with OpenTelemetry support

Configure OTLP exporter
```go
loggerProvider := sdk.NewLoggerProvider(
	sdk.WithBatcher(otlplogs.New(ctx, otlplogshttp.NewClient())),
	sdk.WithResource(newResource()), 
	)
otel.SetLoggerProvider(loggerProvider)	
```

Configure otelzap logger
```go
  zapOtlpCore := otelzap.NewOtlpCore(logsProvider)
```

Send logs with tracing information

```go
var ctx context.Context = ... // should be instrumented with opentelemetry instrumentation

otelzap.Ctx(ctx).Info("My message with trace context")
```
