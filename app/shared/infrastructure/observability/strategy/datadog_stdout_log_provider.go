package strategy

import (
	"multi-folder-components/app/shared/configuration"
	"context"
	"log/slog"
	"os"
	"strconv"

	"go.opentelemetry.io/otel/trace"
)

func DatadogStdoutLogProvider(conf configuration.Conf) *slog.Logger {
	baseHandler := slog.NewJSONHandler(os.Stdout, nil)

	datadogHandler := newDatadogHandler(baseHandler)

	return slog.New(datadogHandler).With(
		slog.String(ddEnvKey, conf.ENVIRONMENT),
		slog.String(ddVersionKey, conf.VERSION),
		slog.String(ddServiceKey, conf.PROJECT_NAME),
	)
}

type DatadogHandler struct {
	baseHandler slog.Handler
}

func newDatadogHandler(baseHandler slog.Handler) *DatadogHandler {
	return &DatadogHandler{baseHandler: baseHandler}
}

func (h *DatadogHandler) Handle(ctx context.Context, record slog.Record) error {
	if spanContext := trace.SpanContextFromContext(ctx); spanContext.IsValid() {
		record.AddAttrs(
			slog.String(ddTraceIDKey, convertTraceID(spanContext.TraceID().String())),
			slog.String(ddSpanIDKey, convertTraceID(spanContext.SpanID().String())),
		)
	}

	return h.baseHandler.Handle(ctx, record)
}

func (h *DatadogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &DatadogHandler{baseHandler: h.baseHandler.WithAttrs(attrs)}
}

func (h *DatadogHandler) WithGroup(name string) slog.Handler {
	return &DatadogHandler{baseHandler: h.baseHandler.WithGroup(name)}
}

func (h *DatadogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.baseHandler.Enabled(ctx, level)
}

func convertTraceID(id string) string {
	if len(id) < 16 {
		return ""
	}
	if len(id) > 16 {
		id = id[16:]
	}
	intValue, err := strconv.ParseUint(id, 16, 64)
	if err != nil {
		return ""
	}
	return strconv.FormatUint(intValue, 10)
}
