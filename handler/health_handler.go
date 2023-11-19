package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func HealthHandler(c echo.Context) error {
	type HealthResponse struct {
		Healthy bool      `json:"healthy"`
		Now     time.Time `json:"now"`
	}

	return c.JSON(http.StatusOK, HealthResponse{Healthy: true, Now: time.Now().UTC()})
}
