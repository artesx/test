package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"test-work/internal/services"
	"test-work/internal/services/banner"

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

func GetStatistic(serviceLayer *services.ServiceLayer) func(c echo.Context) error {
	return func(c echo.Context) error {

		body := &banner.GetStatisticBody{}

		if err := c.Bind(body); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		bannerParam := c.Param("bannerID")

		bannerID, err := strconv.ParseInt(bannerParam, 10, 64)

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		result, err := serviceLayer.BannerService.GetStatistic(uint64(bannerID), body)
		fmt.Println(err)

		if err != nil {
			return writeJSON(c, http.StatusInternalServerError, err.Error())
		}

		return writeJSON(c, http.StatusOK, result)
	}
}

func writeJSON(c echo.Context, httpCode int, data any) error {
	if err := c.JSON(httpCode, data); err != nil {
		return fmt.Errorf("error while render JSON | %w", err)
	}

	return nil
}
