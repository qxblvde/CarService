package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qxblvde/CarService/order-service/internal/infrastructure/carstock"
)

type StockCarHandler struct {
	cars *carstock.Client
}

func NewStockCarHandler(cars *carstock.Client) *StockCarHandler {
	return &StockCarHandler{cars: cars}
}

func (h *StockCarHandler) Register(rg *gin.RouterGroup) {
	rg.GET("", h.list)
	rg.GET("/:id", h.get)
}

type stockCarResponse struct {
	ID    string `json:"id"`
	Brand string `json:"brand"`
	Model string `json:"model"`
	Year  int    `json:"year"`
	Price int64  `json:"price"`
}

func toStockCarResponse(c carstock.Car) stockCarResponse {
	return stockCarResponse{
		ID:    c.ID,
		Brand: c.Brand,
		Model: c.Model,
		Year:  c.Year,
		Price: c.Price,
	}
}

func (h *StockCarHandler) list(c *gin.Context) {
	cars, err := h.cars.List(c.Request.Context())
	if err != nil {
		writeStorageError(c, err)
		return
	}

	resp := make([]stockCarResponse, 0, len(cars))
	for _, car := range cars {
		resp = append(resp, toStockCarResponse(car))
	}
	c.JSON(http.StatusOK, resp)
}

func (h *StockCarHandler) get(c *gin.Context) {
	car, err := h.cars.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		writeStorageError(c, err)
		return
	}
	c.JSON(http.StatusOK, toStockCarResponse(car))
}

func writeStorageError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, carstock.ErrCarNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "car not found"})
	case errors.Is(err, carstock.ErrUnavailable):
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "storage service unavailable"})
	default:
		c.JSON(http.StatusBadGateway, gin.H{"error": "storage request failed"})
	}
}
