package handler

import (
	"net/http"

	"github.com/qxblvde/CarService/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/qxblvde/CarService/order-service/internal/application/persistence"
	"github.com/qxblvde/CarService/order-service/internal/application/service"
	"github.com/qxblvde/CarService/order-service/internal/presentation/dto"
)

type TestDriveHandler struct {
	svc *service.TestDriveService
}

func NewTestDriveHandler(svc *service.TestDriveService) *TestDriveHandler {
	return &TestDriveHandler{svc: svc}
}

func (h *TestDriveHandler) Register(rg *gin.RouterGroup) {
	rg.GET("", middleware.RequireRole(middleware.RoleUser, middleware.RoleManager), h.list)
	rg.POST("", middleware.RequireRole(middleware.RoleUser), h.create)
	rg.GET("/:id", middleware.RequireRole(middleware.RoleUser, middleware.RoleManager), h.getByID)
	rg.DELETE("/:id", middleware.RequireRole(middleware.RoleUser, middleware.RoleManager), h.cancel)
}

func (h *TestDriveHandler) list(c *gin.Context) {
	reqs, err := h.svc.List(c.Request.Context(), persistence.TestDriveFilter{
		ClientID: c.Query("clientId"),
		CarID:    c.Query("carId"),
	})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.TestDrivesToResponse(reqs))
}

func (h *TestDriveHandler) create(c *gin.Context) {
	var req dto.CreateTestDriveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	td, err := h.svc.Create(c.Request.Context(), service.CreateTestDriveParams{
		ClientID:      req.ClientID,
		CarID:         req.CarID,
		ScheduledTime: req.ScheduledTime,
	})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.TestDriveToResponse(td))
}

func (h *TestDriveHandler) getByID(c *gin.Context) {
	td, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.TestDriveToResponse(td))
}

func (h *TestDriveHandler) cancel(c *gin.Context) {
	if err := h.svc.Cancel(c.Request.Context(), c.Param("id")); err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
