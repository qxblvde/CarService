package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qxblvde/CarService/storage-service/internal/application/service"
	"github.com/qxblvde/CarService/storage-service/internal/domain"
)

type AssemblyOrderHandler struct {
	svc *service.AssemblyOrderService
}

func NewAssemblyOrderHandler(svc *service.AssemblyOrderService) *AssemblyOrderHandler {
	return &AssemblyOrderHandler{svc: svc}
}

func (h *AssemblyOrderHandler) Register(rg *gin.RouterGroup) {
	rg.POST("", h.create)
	rg.GET("", h.list)
	rg.GET("/:id", h.getByID)
	rg.POST("/:id/assemble", h.assemble)
	rg.POST("/:id/fail", h.fail)
	rg.DELETE("/:id", h.remove)
}

type createRequest struct {
	SourceOrderID string                 `json:"sourceOrderId" binding:"required"`
	OrderType     domain.SourceOrderType `json:"orderType"     binding:"required"`
}

func (h *AssemblyOrderHandler) create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sourceOrderId and orderType are required"})
		return
	}
	o, err := h.svc.Create(c.Request.Context(), req.SourceOrderID, req.OrderType)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, toResponse(o))
}

func (h *AssemblyOrderHandler) list(c *gin.Context) {
	orders, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := make([]assemblyOrderResponse, 0, len(orders))
	for _, o := range orders {
		resp = append(resp, toResponse(o))
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AssemblyOrderHandler) getByID(c *gin.Context) {
	o, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, toResponse(o))
}

func (h *AssemblyOrderHandler) assemble(c *gin.Context) {
	if err := h.svc.Assemble(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *AssemblyOrderHandler) fail(c *gin.Context) {
	if err := h.svc.Fail(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *AssemblyOrderHandler) remove(c *gin.Context) {
	if err := h.svc.Remove(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

type assemblyOrderResponse struct {
	ID            string `json:"id"`
	SourceOrderID string `json:"sourceOrderId"`
	OrderType     string `json:"orderType"`
	Status        string `json:"status"`
	Removed       bool   `json:"removed"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

func toResponse(o domain.AssemblyOrder) assemblyOrderResponse {
	return assemblyOrderResponse{
		ID:            o.ID,
		SourceOrderID: o.SourceOrderID,
		OrderType:     string(o.OrderType),
		Status:        string(o.Status),
		Removed:       o.Removed,
		CreatedAt:     o.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     o.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
