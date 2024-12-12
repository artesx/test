package handlers

import (
	"fmt"
	"net/http"
	"test-work/internal/services"

	"github.com/labstack/echo/v4"
)

func Banner(serviceLayer *services.ServiceLayer) func(c echo.Context) error {
	return func(c echo.Context) error {
		bannerID := c.Param("bannerID")

		result, err := serviceLayer.BannerService.ProduceBanner(bannerID)

		if err != nil {
			return writeJSON(c, http.StatusInternalServerError, err.Error())
		}

		return writeJSON(c, http.StatusOK, result)
	}
}

//func CreateTrade(serviceLayer *services.ServiceLayer, ctx context.Context) func(c echo.Context) error {
//	return func(c echo.Context) error {
//
//		body := &domain.CreateTradeRequest{}
//
//		if err := c.Bind(body); err != nil {
//			fmt.Println("CREATE TRADE BIND BODY ERROR", err, body)
//			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
//		}
//
//		if err := c.Validate(body); err != nil {
//			fmt.Println("CREATE TRADE Validate BODY ERROR", err, body)
//			return err
//		}
//
//		result, err := serviceLayer.Trade.CreateTrade(ctx, body)
//
//		if err != nil {
//			return writeJSON(c, http.StatusInternalServerError, err.Error())
//		}
//
//		return writeJSON(c, http.StatusOK, result)
//	}
//}

func writeJSON(c echo.Context, httpCode int, data any) error {
	if err := c.JSON(httpCode, data); err != nil {
		return fmt.Errorf("error while render JSON | %w", err)
	}

	return nil
}
