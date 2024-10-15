package api

import (
	"multi-folder-components/app/shared/infrastructure/observability"
	"multi-folder-components/app/shared/infrastructure/serverwrapper"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/labstack/echo/v4"
)

func init() {
	ioc.Registry(testingObs, serverwrapper.NewEchoWrapper, observability.NewObservability)
}
func testingObs(e serverwrapper.EchoWrapper, obs observability.Observability) {
	e.GET("/testing-obs", func(c echo.Context) error {
		ctx, span := obs.Tracer.Start(c.Request().Context(), "hello")
		defer span.End()
		obs.Logger.InfoContext(ctx, "hello mom", "a", "b")
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Unimplemented",
		})
	})
}
