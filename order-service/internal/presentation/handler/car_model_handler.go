package handler

import (
	"net/http"

	"github.com/qxblvde/CarService/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/qxblvde/CarService/order-service/internal/application/service"
	"github.com/qxblvde/CarService/order-service/internal/presentation/dto"
)

type CarModelHandler struct {
	svc     *service.CarModelService
	details *service.DetailService
}

func NewCarModelHandler(svc *service.CarModelService, details *service.DetailService) *CarModelHandler {
	return &CarModelHandler{svc: svc, details: details}
}

func (h *CarModelHandler) Register(rg *gin.RouterGroup) {
	rg.GET("", h.list)
	rg.POST("", middleware.RequireRole(middleware.RoleManager), h.create)
	rg.GET("/:id", h.getByID)
	rg.DELETE("/:id", middleware.RequireRole(middleware.RoleAdmin), h.delete)
	rg.GET("/:id/details", h.listDetails)
}

func (h *CarModelHandler) list(c *gin.Context) {
	models, err := h.svc.List(c.Request.Context(), c.Query("brand"), c.QueryArray("detailId"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.CarModelsToResponse(models))
}

func (h *CarModelHandler) create(c *gin.Context) {
	var req dto.CreateCarModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	m, err := h.svc.Create(c.Request.Context(), service.CreateCarModelParams{
		Brand:          req.Brand,
		Model:          req.Model,
		BasePrice:      req.BasePrice,
		BodyType:       req.BodyType,
		FuelType:       req.FuelType,
		EnginePowerHP:  req.EnginePowerHP,
		EngineVolumeCC: req.EngineVolumeCC,
		Transmission:   req.Transmission,
		DriveType:      req.DriveType,
	})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.CarModelToResponse(m))
}

func (h *CarModelHandler) getByID(c *gin.Context) {
	m, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.CarModelToResponse(m))
}

func (h *CarModelHandler) delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *CarModelHandler) listDetails(c *gin.Context) {
	details, err := h.details.ListByModel(c.Request.Context(), c.Param("id"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.DetailsToResponse(details))
}
