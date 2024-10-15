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
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	otelmeter "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// NewGRPCOpenObserveMeterProvider configures the meter provider for OpenObserve.
func NewGRPCOpenObserveMeterProvider(conf configuration.Conf) (otelmeter.Meter, error) {

	var exporterOpts []otlpmetricgrpc.Option

	exporterOpts = append(exporterOpts, otlpmetricgrpc.WithEndpoint(conf.LoadFromSystem(OTEL_EXPORTER_OTLP_ENDPOINT)))
	if conf.LoadFromSystem(OTEL_EXPORTER_OTLP_INSECURE) == "true" {
		exporterOpts = append(exporterOpts, otlpmetricgrpc.WithTLSCredentials(insecure.NewCredentials()))
	}
	exporterOpts = append(exporterOpts, otlpmetricgrpc.WithDialOption(grpc.WithUnaryInterceptor(func(
		ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption) error {
		md := metadata.New(map[string]string{
			"Authorization": conf.LoadFromSystem(OPENOBSERVE_AUTHORIZATION),
			"organization":  conf.LoadFromSystem(OPENOBSERVE_ORGANIZATION),
			"stream-name":   conf.LoadFromSystem(OPENOBSERVE_STREAM_NAME),
		})
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	})))
	ctx, cancel := context.WithCancel(context.Background())
	exporter, err := otlpmetricgrpc.New(ctx, exporterOpts...)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("creating OTLP metric exporter: %w", err)
	}
	// Set up graceful shutdown.
	meterProvider := metric.NewMeterProvider(
		metric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(conf.PROJECT_NAME),
			semconv.DeploymentEnvironmentKey.String(conf.ENVIRONMENT),
		)),
		metric.WithReader(metric.NewPeriodicReader(exporter,
			// Default interval for exporting metrics.
			metric.WithInterval(5*time.Second))),
	)

	// Register the meter provider as the global provider.
	otel.SetMeterProvider(meterProvider)

	go func() {
		// Wait for termination signal.
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

	return meterProvider.Meter("openobserve"), nil
}
