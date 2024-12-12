package http_transport

import (
	"context"
	"github.com/labstack/echo/v4"
	"test-work/internal/http_transport/handlers"
)

func (h *Handlers) RegisterRoutes(ctx context.Context, echo *echo.Echo) {
	v1 := echo.Group("/api/")
	v1.GET("counter/:bannerID", handlers.Banner(h.ServiceLayer))
	v1.POST("stats/:bannerID", handlers.GetStatistic(h.ServiceLayer))
}
