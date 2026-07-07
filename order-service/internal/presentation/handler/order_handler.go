package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qxblvde/CarService/internal/middleware"
	"github.com/qxblvde/CarService/order-service/internal/application/persistence"
	"github.com/qxblvde/CarService/order-service/internal/application/service"
	"github.com/qxblvde/CarService/order-service/internal/presentation/dto"
)

type OrderHandler struct {
	svc *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) RegisterInStock(rg *gin.RouterGroup) {
	rg.GET("", middleware.RequireRole(middleware.RoleUser, middleware.RoleManager), h.listInStock)
	rg.POST("", middleware.RequireRole(middleware.RoleUser), h.createInStock)
	rg.GET("/:id", middleware.RequireRole(middleware.RoleUser, middleware.RoleManager), h.getInStock)
	rg.POST("/:id/approve", middleware.RequireRole(middleware.RoleManager), h.action(h.svc.ApproveInStockOrder))
	rg.POST("/:id/wait-payment", middleware.RequireRole(middleware.RoleManager), h.action(h.svc.InStockOrderWaitPayment))
	rg.POST("/:id/pay", middleware.RequireRole(middleware.RoleManager), h.action(h.svc.MarkInStockOrderPaid))
	rg.POST("/:id/ready-for-pickup", middleware.RequireRole(middleware.RoleWarehouseAdmin), h.action(h.svc.InStockOrderReadyForPickup))
	rg.POST("/:id/complete", middleware.RequireRole(middleware.RoleManager), h.action(h.svc.CompleteInStockOrder))
	rg.POST("/:id/cancel", middleware.RequireRole(middleware.RoleUser, middleware.RoleManager), h.cancelInStock)
}

func (h *OrderHandler) RegisterCustom(rg *gin.RouterGroup) {
	rg.GET("", middleware.RequireRole(middleware.RoleUser, middleware.RoleManager), h.listCustom)
	rg.POST("", middleware.RequireRole(middleware.RoleUser), h.createCustom)
	rg.GET("/:id", middleware.RequireRole(middleware.RoleUser, middleware.RoleManager), h.getCustom)
	rg.POST("/:id/approve", middleware.RequireRole(middleware.RoleManager), h.action(h.svc.ApproveCustomOrder))
	rg.POST("/:id/wait-payment", middleware.RequireRole(middleware.RoleManager), h.action(h.svc.CustomOrderWaitPayment))
	rg.POST("/:id/pay", middleware.RequireRole(middleware.RoleManager), h.action(h.svc.MarkCustomOrderPaid))
	rg.POST("/:id/waiting-delivery", middleware.RequireRole(middleware.RoleWarehouseAdmin), h.action(h.svc.CustomOrderWaitingDelivery))
	rg.POST("/:id/ready-for-pickup", middleware.RequireRole(middleware.RoleWarehouseAdmin), h.action(h.svc.CustomOrderReadyForPickup))
	rg.POST("/:id/complete", middleware.RequireRole(middleware.RoleManager), h.action(h.svc.CompleteCustomOrder))
	rg.POST("/:id/cancel", middleware.RequireRole(middleware.RoleUser, middleware.RoleManager), h.cancelCustom)
}

func (h *OrderHandler) action(fn func(context.Context, string) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := fn(c.Request.Context(), c.Param("id")); err != nil {
			writeError(c, err)
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func (h *OrderHandler) listInStock(c *gin.Context) {
	caller, _ := middleware.CallerFromCtx(c.Request.Context())
	filter := persistence.OrderFilter{
		ManagerID: c.Query("managerId"),
	}
	if caller.HasRole(middleware.RoleUser) && !caller.IsManager() {
		filter.ClientID = caller.ID
	} else {
		filter.ClientID = c.Query("clientId")
	}

	orders, err := h.svc.ListInStockOrders(c.Request.Context(), filter)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.InStockOrdersToResponse(orders))
}

func (h *OrderHandler) createInStock(c *gin.Context) {
	caller, _ := middleware.CallerFromCtx(c.Request.Context())
	var req dto.CreateInStockOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	o, err := h.svc.CreateInStockOrder(c.Request.Context(), service.CreateInStockOrderParams{
		ManagerID: req.ManagerID,
		ClientID:  caller.ID,
		CarID:     req.CarID,
	})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.InStockOrderToResponse(o))
}

func (h *OrderHandler) getInStock(c *gin.Context) {
	caller, _ := middleware.CallerFromCtx(c.Request.Context())
	o, err := h.svc.GetInStockOrder(c.Request.Context(), c.Param("id"))
	if err != nil {
		writeError(c, err)
		return
	}
	if caller.HasRole(middleware.RoleUser) && !caller.IsManager() && o.ClientID != caller.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	c.JSON(http.StatusOK, dto.InStockOrderToResponse(o))
}

func (h *OrderHandler) cancelInStock(c *gin.Context) {
	caller, _ := middleware.CallerFromCtx(c.Request.Context())
	if caller.HasRole(middleware.RoleUser) && !caller.IsManager() {
		o, err := h.svc.GetInStockOrder(c.Request.Context(), c.Param("id"))
		if err != nil {
			writeError(c, err)
			return
		}
		if o.ClientID != caller.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
	}
	if err := h.svc.CancelInStockOrder(c.Request.Context(), c.Param("id")); err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *OrderHandler) listCustom(c *gin.Context) {
	caller, _ := middleware.CallerFromCtx(c.Request.Context())
	filter := persistence.OrderFilter{
		ManagerID: c.Query("managerId"),
	}
	if caller.HasRole(middleware.RoleUser) && !caller.IsManager() {
		filter.ClientID = caller.ID
	} else {
		filter.ClientID = c.Query("clientId")
	}

	orders, err := h.svc.ListCustomOrders(c.Request.Context(), filter)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.CustomOrdersToResponse(orders))
}

func (h *OrderHandler) createCustom(c *gin.Context) {
	caller, _ := middleware.CallerFromCtx(c.Request.Context())
	var req dto.CreateCustomOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	o, err := h.svc.CreateCustomOrder(c.Request.Context(), service.CreateCustomOrderParams{
		ManagerID:  req.ManagerID,
		ClientID:   caller.ID,
		CarModelID: req.CarModelID,
		DetailIDs:  req.DetailIDs,
	})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.CustomOrderToResponse(o))
}

func (h *OrderHandler) getCustom(c *gin.Context) {
	caller, _ := middleware.CallerFromCtx(c.Request.Context())
	o, err := h.svc.GetCustomOrder(c.Request.Context(), c.Param("id"))
	if err != nil {
		writeError(c, err)
		return
	}
	if caller.HasRole(middleware.RoleUser) && !caller.IsManager() && o.ClientID != caller.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	c.JSON(http.StatusOK, dto.CustomOrderToResponse(o))
}

func (h *OrderHandler) cancelCustom(c *gin.Context) {
	caller, _ := middleware.CallerFromCtx(c.Request.Context())
	if caller.HasRole(middleware.RoleUser) && !caller.IsManager() {
		o, err := h.svc.GetCustomOrder(c.Request.Context(), c.Param("id"))
		if err != nil {
			writeError(c, err)
			return
		}
		if o.ClientID != caller.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
	}
	if err := h.svc.CancelCustomOrder(c.Request.Context(), c.Param("id")); err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
