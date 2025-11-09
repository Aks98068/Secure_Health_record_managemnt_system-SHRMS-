package Middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/tools/go/analysis/passes/defers"
)


// Rate Limiting Middleware
// This middleware controls how many requests a single client (IP) can make
// within a specific time window. It helps protect the application from
// brute-force attempts on login and register routes, prevents spam or
// automated abuse, and reduces unnecessary load on the server.
//
// When a client exceeds the allowed number of requests, the middleware returns
// HTTP 429 (Too Many Requests) and blocks further attempts until the window resets.
//
// We apply stricter limits to public authentication routes (login/register)
// and more relaxed or no limits on authenticated, role-protected routes.


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
			c.JSON(429, gin.H{
				"error": "Too many requests. Slow down.",
			})
			c.Abort()
			return
		}

		// Reset after 1 minute
		go func() {
			time.Sleep(1 * time.Minute)
			rateLimitLock.Lock()
			rateLimitIP[ip] = 0
			rateLimitLock.Unlock()
		}()

		c.Next()
	}
}

 
// Brute-Force Protection Middleware
// This middleware blocks repeated failed login attempts from the same IP or
// email within a short time window. It prevents automated bots from guessing
// passwords by limiting how many failed attempts are allowed.
//
// When too many failed attempts occur, the login route is temporarily locked
// for that IP/email combination. This protects all login and registration
// endpoints from credential stuffing and brute-force attacks.

// BruteforceMiddleware for blockings  repeated attempts (IP + username) 
func BruteforceMiddleware() gin.HandlerFunc{
	return func(c *gin.Context){
		ip:= c. ClientIP()
		username:= c.PostForm("email")


		mutex.Lock()
		defers mutex.Unlock()


		// check ip attempts
		if loginAttempts[ip] >=5{
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many attepmpts from this IP. try again after 5 minutes.",
			})
			c.Abort()
			return 
		}

		// check username attempts
		if username != "" && userLoginAttempts[username] >=5{
			c.JSON(http.StatusTooManyRequests, g.H{
				"error":"Account temporarily locked. Try again after 5 minutes.",
			})
			c.Abort()
			return
		}

		// Allow request
		c.Next()

		// on failed login , increase counters
		if c.Writer.Status()==http.StatusUnauthorized{
			loginAttempts[ip]++
			if username != ""{
userLoginAttempts[username]++
			}

			// rest countrs after 5 minutes
			go func(ip, user string){
				time.Sleep(5 *time.Minute)
				mutex.Lock()
				defer mutex.Unlock()
				loginAttempts[ip] = 0
				if user !=""{
					userLoginAttempts[user]=0
				}
			}(ip, username)
		}
	}
}
