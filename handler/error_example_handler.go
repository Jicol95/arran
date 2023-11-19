package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func ErrorExampleHandler(c echo.Context) error {
	return echo.NewHTTPError(http.StatusInternalServerError)
}
