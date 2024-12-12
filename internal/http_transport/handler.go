package http_transport

import (
	"context"
	"test-work/internal/services"
)

type Handlers struct {
	ServiceLayer *services.ServiceLayer
}

func RegisterHandlers(ctx context.Context, httpServerApi *API, serviceLayer *services.ServiceLayer) error {
	handlers := Handlers{ServiceLayer: serviceLayer}
	handlers.RegisterRoutes(ctx, httpServerApi.Echo)

	return nil
}
