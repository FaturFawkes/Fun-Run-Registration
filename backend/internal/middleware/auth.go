package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tau-tau-run/backend/internal/services"
)

// AuthMiddleware validates JWT tokens for protected routes
func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			RespondWithError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required. Please provide a valid token.", nil)
			c.Abort()
			return
		}

		// Check Bearer format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			RespondWithError(c, http.StatusUnauthorized, "INVALID_TOKEN_FORMAT", "Authorization header must be in format: Bearer <token>", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			RespondWithError(c, http.StatusForbidden, "TOKEN_EXPIRED", "Your session has expired. Please log in again.", nil)
			c.Abort()
			return
		}

		// Attach admin info to context
		c.Set("admin_id", claims.AdminID)
		c.Set("admin_email", claims.Email)

		c.Next()
	}
}

// GetAdminID retrieves admin ID from context
func GetAdminID(c *gin.Context) string {
	if adminID, exists := c.Get("admin_id"); exists {
		return adminID.(string)
	}
	return ""
}

// GetAdminEmail retrieves admin email from context
func GetAdminEmail(c *gin.Context) string {
	if email, exists := c.Get("admin_email"); exists {
		return email.(string)
	}
	return ""
}
