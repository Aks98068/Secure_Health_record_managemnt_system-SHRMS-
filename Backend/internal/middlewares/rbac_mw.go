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
		// Get user role from context (set by AuthMiddleware)
		userRole := c.GetString("role") 
		
		// Validate that user role exists
		if userRole == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
				"code":  "AUTH_REQUIRED",
			})
			return
		}
		
		// Check if user role matches any of the allowed roles
		for _, allowedRole := range roles {
			if userRole == allowedRole {
				c.Next()
				return
			}
		}
		
		// Enhanced error response with more context
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error":        "Access denied",
			"code":         "INSUFFICIENT_PERMISSIONS",
			"user_role":    userRole,
			"required_roles": roles,
		})
	}
}
