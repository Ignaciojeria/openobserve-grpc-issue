package observability

import (
	"multi-folder-components/app/shared/configuration"
	"multi-folder-components/app/shared/infrastructure/observability/strategy"
	"log/slog"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(
		newLoggerProvider,
		configuration.NewConf,
	)
}

func newLoggerProvider(conf configuration.Conf) (*slog.Logger, error) {
	// Get the observability strategy
	observabilityStrategyKey := conf.LoadFromSystem(strategy.OBSERVABILITY_STRATEGY)
	switch observabilityStrategyKey {
	case "openobserve":
		return strategy.OpenObserveGRPCLogProvider(conf)
	case "datadog":
		return strategy.DatadogGRPCLogProvider(conf)
	default:
		return strategy.NoOpStdoutLogProvider(conf), nil
	}
}
