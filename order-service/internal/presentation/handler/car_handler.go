package handler

import (
	"net/http"

	"github.com/qxblvde/CarService/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/qxblvde/CarService/order-service/internal/application/service"
	"github.com/qxblvde/CarService/order-service/internal/presentation/dto"
)

type CarHandler struct {
	svc *service.CarService
}

func NewCarHandler(svc *service.CarService) *CarHandler {
	return &CarHandler{svc: svc}
}

func (h *CarHandler) Register(rg *gin.RouterGroup) {
	rg.GET("", h.list)
	rg.POST("", middleware.RequireRole(middleware.RoleManager, middleware.RoleWarehouseAdmin), h.create)
	rg.GET("/:id", h.getByID)
	rg.DELETE("/:id", middleware.RequireRole(middleware.RoleAdmin), h.delete)
}

func (h *CarHandler) list(c *gin.Context) {
	cars, err := h.svc.List(c.Request.Context(), c.Query("carModelId"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.CarsToResponse(cars))
}

func (h *CarHandler) create(c *gin.Context) {
	var req dto.CreateCarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	car, err := h.svc.Create(c.Request.Context(), service.CreateCarParams{
		CarModelID: req.CarModelID,
		VIN:        req.VIN,
		Year:       req.Year,
		Price:      req.Price,
	})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.CarToResponse(car))
}

func (h *CarHandler) getByID(c *gin.Context) {
	car, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.CarToResponse(car))
}

func (h *CarHandler) delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
