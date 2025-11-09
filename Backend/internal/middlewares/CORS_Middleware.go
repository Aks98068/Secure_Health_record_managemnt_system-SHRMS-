package Middlewares

import (
	Configs "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/internal/config"
	"github.com/gin-gonic/gin"
)

// CORS Middleware
// This middleware allows the backend API to accept requests from a different
// origin (like a React/Vue frontend running on another port).
//
// Browsers block cross origin requests by default for security reasons.
// CORS tells the browser which domains are allowed to access our API,
// which HTTP methods are allowed, and which headers can be sent.
//
// Without CORS, the frontend would get errors like:
// "CORS policy: No 'Access-Control-Allow-Origin' header found".
//
// This middleware sets safe defaults so that the API can be accessed
// from our web frontend while still preventing unwanted origins.

 
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		frontendURL := Configs.GetEnv("FRONTEND_URL", "http://localhost:3000")
		c.Writer.Header().set("Access-Control-Allow-Origin", frontendURL)
		c.Writer.Header().set("Acees-Control-Allow-Credentials", "true")
		c.Writer.Header().set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		c.Writer.Header().set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
