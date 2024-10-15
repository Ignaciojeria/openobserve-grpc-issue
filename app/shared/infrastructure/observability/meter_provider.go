package observability

import (
	"multi-folder-components/app/shared/configuration"
	"multi-folder-components/app/shared/infrastructure/observability/strategy"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	otelmeter "go.opentelemetry.io/otel/metric"
)

func init() {
	ioc.Registry(
		newMeterProvider,
		configuration.NewConf,
	)
}

func newMeterProvider(conf configuration.Conf) (otelmeter.Meter, error) {
	// Get the observability strategy
	observabilityStrategyKey := conf.LoadFromSystem(strategy.OBSERVABILITY_STRATEGY)
	switch observabilityStrategyKey {
	default:
		return strategy.NoOpMeterProvider(conf)
	}
}
