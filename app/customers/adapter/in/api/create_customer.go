package api

import (
	"multi-folder-components/app/shared/infrastructure/serverwrapper"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/labstack/echo/v4"
)

func init() {
	ioc.Registry(createCustomer, serverwrapper.NewEchoWrapper)
}
func createCustomer(e serverwrapper.EchoWrapper) {
	e.POST("/insert-your-custom-pattern-here", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Unimplemented",
		})
	})
}
