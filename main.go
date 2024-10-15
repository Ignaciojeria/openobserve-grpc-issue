package main

import (
	_ "multi-folder-components/app/shared/configuration"
	"multi-folder-components/app/shared/constants"
	_ "embed"
	"log"
	"os"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	_ "multi-folder-components/app/shared/infrastructure/serverwrapper"
	_ "multi-folder-components/app/shared/infrastructure/healthcheck"
	_ "multi-folder-components/app/orders/adapter/in/api"
	_ "multi-folder-components/app/payments/adapter/in/api"
	_ "multi-folder-components/app/customers/adapter/in/api"
	_ "multi-folder-components/app/shared/infrastructure/observability"
	_ "multi-folder-components/app/shared/infrastructure/observability/strategy"
)

//go:embed .version
var version string

func main() {
	os.Setenv(constants.Version, version)
	if err := ioc.LoadDependencies(); err != nil {
		log.Fatal(err)
	}
}
