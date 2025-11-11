package Middlewares

import (
	"net/http"
	"strings"

	Utils "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/internal/utils"
	"github.com/gin-gonic/gin"
)

// Authentication Middleware
// This middleware validates the JWT sent by the client in the Authorization
// header. It ensures the user is logged in and the token is valid, not expired,
// not tampered with, and signed using the server’s private key.
//
// If verification succeeds, the middleware extracts user data (UUID, role,
// email) and attaches it to the request context so protected routes can access
// it. All protected role-based API endpoints must pass through this middleware.

func AuthMiddleware(jwt *Utils.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		autheader := c.GetHeader("Authorization")
		if autheader == "" || !strings.HasPrefix(autheader, "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			return
		}
		token := strings.TrimPrefix(autheader, "Bearer")
		claims, err := jwt.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()

	}
}
