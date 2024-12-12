package http_transport

import (
	"github.com/labstack/echo/v4"
)

type API struct {
	Echo *echo.Echo
}

func NewAPI(e *echo.Echo) *API {
	return &API{Echo: e}
}