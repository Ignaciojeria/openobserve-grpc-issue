package strategy

import (
	"multi-folder-components/app/shared/configuration"
	"log/slog"
	"os"

	otelslogjson "github.com/go-slog/otelslog"
)

func NoOpStdoutLogProvider(conf configuration.Conf) *slog.Logger {
	return slog.New(otelslogjson.NewHandler(slog.NewJSONHandler(os.Stdout, nil))).With(
		slog.String("env", conf.ENVIRONMENT),
		slog.String("version", conf.VERSION),
		slog.String("service", conf.PROJECT_NAME),
	)
}
