package strategy

import (
	"multi-folder-components/app/shared/configuration"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/otel"
	otelmeter "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// NoOpMeterProvider configures a meter provider that does not export metrics.
func NoOpMeterProvider(conf configuration.Conf) (otelmeter.Meter, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Create a resource with the service and environment attributes.
	meterProvider := metric.NewMeterProvider(
		metric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(conf.PROJECT_NAME),
			semconv.DeploymentEnvironmentKey.String(conf.ENVIRONMENT),
		)),
		// No exporter or reader is configured, so no metrics will be exported.
	)

	// Register the meter provider as the global provider.
	otel.SetMeterProvider(meterProvider)

	go func() {
		// Handle shutdown signal for clean exit
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
		defer shutdownCancel()
		if err := meterProvider.Shutdown(shutdownCtx); err != nil {
			fmt.Println("Failed to shutdown meter provider:", err)
		}
		cancel()
	}()

	// No metrics are exported, but the Meter instance can still be used.
	return meterProvider.Meter("noop-observability"), nil
}
