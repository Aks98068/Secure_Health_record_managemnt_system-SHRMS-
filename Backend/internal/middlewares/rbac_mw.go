package Middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RBAC Middleware
// This middleware checks whether the authenticated user has permission to
// access the requested route. Each route group is attached to allowed roles.
// When a user tries to access a route, their role (decoded from JWT) is checked
// against allowed roles.
//
// If the user’s role does not match the required roles (admin, doctor, patient),
// access is blocked with HTTP 403 Forbidden. This enforces strict separation
// between admin, doctor and patient privileges.

func RoleGuard(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("roles")
		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
	}
}
