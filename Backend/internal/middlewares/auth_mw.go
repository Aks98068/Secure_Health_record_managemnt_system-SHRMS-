package Middlewares

import (
	"net/http"

	Utils "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/internal/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT from HTTP-only cookie
func AuthMiddleware(jwt *Utils.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read token from secure cookie
		tokenStr, err := c.Cookie("auth_token")
		if err != nil || tokenStr == "" {
			c.Redirect(http.StatusSeeOther, "/login") // redirect to login if missing
			c.Abort()
			return
		}

		claims, err := jwt.VerifyToken(tokenStr)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/login") // redirect if invalid
			c.Abort()
			return
		}

		// Attach user info to context
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
