package strategy

import (
	"multi-folder-components/app/shared/configuration"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc/credentials/insecure"
)

// newGRPCOpenObserveLoggerProvider configures the logger provider for OpenObserve.
func DatadogGRPCLogProvider(conf configuration.Conf) (*slog.Logger, error) {
	ctx, cancel := context.WithCancel(context.Background())

	var exporterOpts []otlploggrpc.Option
	exporterOpts = append(exporterOpts, otlploggrpc.WithEndpoint(conf.LoadFromSystem(OTEL_EXPORTER_OTLP_ENDPOINT)))
	if conf.LoadFromSystem(OTEL_EXPORTER_OTLP_INSECURE) == "true" {
		exporterOpts = append(exporterOpts, otlploggrpc.WithTLSCredentials(insecure.NewCredentials()))
	}
	// Create the exporter
	exporter, err := otlploggrpc.New(ctx, exporterOpts...)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("creating OTLP log exporter: %w", err)
	}

	// Set up the processor
	logProcessor := log.NewBatchProcessor(exporter)

	// Define the resource attributes
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(conf.PROJECT_NAME),
		semconv.ServiceVersionKey.String(conf.VERSION),
		semconv.DeploymentEnvironmentKey.String(conf.ENVIRONMENT),
	)

	// Create the LoggerProvider
	loggerProvider := log.NewLoggerProvider(
		log.WithResource(res),
		log.WithProcessor(logProcessor),
	)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
		defer shutdownCancel()
		if err := loggerProvider.Shutdown(shutdownCtx); err != nil {
			fmt.Println("Failed to shutdown logger provider:", err)
		}
		cancel()
	}()

	datadogHandler := newDatadogHandler(otelslog.NewHandler(
		"datadog",
		otelslog.WithLoggerProvider(loggerProvider)))

	return slog.New(datadogHandler).With(
		slog.String(ddEnvKey, conf.ENVIRONMENT),
		slog.String(ddVersionKey, conf.VERSION),
		slog.String(ddServiceKey, conf.PROJECT_NAME),
	), nil
}
