package Middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// =======================
// Rate Limiting Middleware
// =======================

var (
	rateLimitIP   = make(map[string]int)
	rateLimitLock sync.Mutex
)

// Simple IP rate limiting: 30 requests per minute
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		rateLimitLock.Lock()
		rateLimitIP[ip]++
		count := rateLimitIP[ip]
		rateLimitLock.Unlock()

		if count > 30 {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Slow down.",
			})
			c.Abort()
			return
		}

		// Reset count after 1 minute
		go func(ip string) {
			time.Sleep(1 * time.Minute)
			rateLimitLock.Lock()
			rateLimitIP[ip] = 0
			rateLimitLock.Unlock()
		}(ip)

		c.Next()
	}
}

// ================================
// Brute-Force Protection Middleware
// ================================

var (
	loginAttempts     = make(map[string]int) // track attempts by IP
	userLoginAttempts = make(map[string]int) // track attempts by username/email
	mutex             sync.Mutex
)

// BruteforceMiddleware blocks repeated failed login attempts
func BruteforceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		username := c.PostForm("email")

		mutex.Lock()
		defer mutex.Unlock()

		// Check IP attempts
		if loginAttempts[ip] >= 5 {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many attempts from this IP. Try again after 5 minutes.",
			})
			c.Abort()
			return
		}

		// Check username/email attempts
		if username != "" && userLoginAttempts[username] >= 5 {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Account temporarily locked. Try again after 5 minutes.",
			})
			c.Abort()
			return
		}

		c.Next() // allow request to proceed

		// On failed login (HTTP 401), increase counters
		if c.Writer.Status() == http.StatusUnauthorized {
			loginAttempts[ip]++
			if username != "" {
				userLoginAttempts[username]++
			}

			// Reset counters after 5 minutes
			go func(ip, user string) {
				time.Sleep(5 * time.Minute)
				mutex.Lock()
				defer mutex.Unlock()
				loginAttempts[ip] = 0
				if user != "" {
					userLoginAttempts[user] = 0
				}
			}(ip, username)
		}
	}
}
