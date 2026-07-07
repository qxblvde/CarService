package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const callerKey contextKey = "caller"

type Role string

const (
	RoleUser           Role = "user"
	RoleManager        Role = "manager"
	RoleWarehouseAdmin Role = "warehouse_admin"
	RoleAdmin          Role = "admin"
)

type Caller struct {
	ID    string
	Roles []string
}

func (c Caller) HasRole(r Role) bool {
	for _, role := range c.Roles {
		if role == string(r) {
			return true
		}
	}
	return false
}

func (c Caller) IsAdmin() bool   { return c.HasRole(RoleAdmin) }
func (c Caller) IsManager() bool { return c.HasRole(RoleManager) || c.IsAdmin() }

func CallerFromCtx(ctx context.Context) (Caller, bool) {
	c, ok := ctx.Value(callerKey).(Caller)
	return c, ok
}

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}

		sub, _ := claims["sub"].(string)
		var roles []string
		if rs, ok := claims["roles"].([]any); ok {
			for _, r := range rs {
				if s, ok := r.(string); ok {
					roles = append(roles, s)
				}
			}
		}

		caller := Caller{ID: sub, Roles: roles}
		ctx := context.WithValue(c.Request.Context(), callerKey, caller)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func RequireRole(roles ...Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		caller, ok := CallerFromCtx(c.Request.Context())
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		if caller.IsAdmin() {
			c.Next()
			return
		}
		for _, r := range roles {
			if caller.HasRole(r) {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	}
}
