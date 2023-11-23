package handler

import (
	"fmt"
	"net/http"

	"github.com/jicol-95/arran/dal"
	"github.com/jicol-95/arran/domain"
	"github.com/labstack/echo/v4"
)

type ExampleResourceDataResponse struct {
	Data ExampleResourceResponse `json:"data"`
}

type ExampleResourceListDataResponse struct {
	Data []ExampleResourceResponse `json:"data"`
}

type ExampleResourceResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateExampleResourceRequest struct {
	Name string `json:"name"`
}

type UpdateExampleResourceRequest struct {
	Name string `json:"name"`
}

func ExampleResourcePostHandler(tm dal.TransactionManager, svc domain.ExampleResourceService) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger := c.Logger()
		logger.Info("Inserting example resource")

		req := new(CreateExampleResourceRequest)

		if err := c.Bind(req); err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		resource, err := svc.CreateExampleResource(req.Name)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		return c.JSON(
			http.StatusCreated,
			ExampleResourceDataResponse{
				Data: ExampleResourceResponse{
					ID:   resource.ID,
					Name: resource.Name,
				},
			},
		)
	}
}

func ExampleResourceUpdateByIdHandler(tm dal.TransactionManager, svc domain.ExampleResourceService) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger := c.Logger()
		id := c.Param("id")

		req := new(UpdateExampleResourceRequest)

		if err := c.Bind(req); err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		resource, err := svc.UpdateExampleResource(id, req.Name)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		return c.JSON(
			http.StatusCreated,
			ExampleResourceDataResponse{
				Data: ExampleResourceResponse{
					ID:   resource.ID,
					Name: resource.Name,
				},
			},
		)
	}
}

func ExampleResourceGetByIdHandler(tm dal.TransactionManager, svc domain.ExampleResourceService) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger := c.Logger()
		id := c.Param("id")

		logger.Info(fmt.Sprintf("Getting example resource by id: %s", id))

		resource, err := svc.GetExampleResourceById(id)

		if err != nil {
			logger.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		return c.JSON(
			http.StatusOK,
			ExampleResourceDataResponse{
				Data: ExampleResourceResponse{
					ID:   resource.ID,
					Name: resource.Name,
				},
			},
		)
	}
}

func ExampleResourceDeleteByIdHandler(tm dal.TransactionManager, svc domain.ExampleResourceService) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger := c.Logger()
		id := c.Param("id")

		logger.Info(fmt.Sprintf("Deleting example resource by id: %s", id))
		err := svc.DeleteExmpleResourceById(id)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		return c.NoContent(http.StatusOK)
	}
}
