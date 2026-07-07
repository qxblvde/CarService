package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	secret string
}

func NewAuthHandler(secret string) *AuthHandler {
	return &AuthHandler{secret: secret}
}

func (h *AuthHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/login", h.login)
}

type loginRequest struct {
	UserID string   `json:"userId" binding:"required"`
	Roles  []string `json:"roles"  binding:"required"`
}

func (h *AuthHandler) login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId and roles are required"})
		return
	}

	claims := jwt.MapClaims{
		"sub":   req.UserID,
		"roles": req.Roles,
		"exp":   time.Now().Add(5 * time.Minute).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(h.secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not sign token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": signed})
}
