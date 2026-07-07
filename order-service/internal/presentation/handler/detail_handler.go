package handler

import (
	"net/http"

	"github.com/qxblvde/CarService/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/qxblvde/CarService/order-service/internal/application/service"
	"github.com/qxblvde/CarService/order-service/internal/presentation/dto"
)

type DetailHandler struct {
	svc *service.DetailService
}

func NewDetailHandler(svc *service.DetailService) *DetailHandler {
	return &DetailHandler{svc: svc}
}

func (h *DetailHandler) Register(rg *gin.RouterGroup) {
	rg.GET("", h.list)
	rg.POST("", middleware.RequireRole(middleware.RoleWarehouseAdmin), h.create)
	rg.GET("/:id", h.getByID)
}

func (h *DetailHandler) list(c *gin.Context) {
	details, err := h.svc.List(c.Request.Context())
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.DetailsToResponse(details))
}

func (h *DetailHandler) create(c *gin.Context) {
	var req dto.CreateDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	d, err := h.svc.Create(c.Request.Context(), service.CreateDetailParams{
		Name:                  req.Name,
		Type:                  req.Type,
		Price:                 req.Price,
		CompatibleCarModelIDs: req.CompatibleCarModelIDs,
	})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.DetailToResponse(d))
}

func (h *DetailHandler) getByID(c *gin.Context) {
	d, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.DetailToResponse(d))
}
