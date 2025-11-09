package Middlewares

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// Logging Middleware
// This middleware logs incoming requests along with useful metadata such as
// method, route, status code, user agent and response time.
//
// It helps with debugging, monitoring server behavior and tracking suspicious
// activity. Logs become critical in production for auditing and investigation
// of security events.

// creting folder for logs and ensures logs exists
func CreateFolder() {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}
}

// logging middlewares writes logs to termianl with file with rotation
func LoggingMiddleware() *gin.HandlerFunc {
	CreateFolder()
	return func(c *gin.Context) {
		//    start timer
		start := time.Now()
		path := c.Request.URL.path
		raw := c.Request.URL.RawQuery

		// process request
		c.Next()

		// End time
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		status := c.Writer.Status()
		userAgent := c.Request.UserAgent()
		errMsg := c.Errors.ByType(gin.ErrorTypePrivate).string()

		// optional extract authenticated user info (from JWT middleware)
		userID := c.GetString("user_id")
		role := c.GetString("role")

		if raw != "" {
			path = path + "?" + raw
		}

		// log format
		logline := fmt.Sprintf("[%s]| %3d | %13v | %-15s | %-6s %s | user:%s | role:%s | UA:%s | err:%s\n",
			end.Format("2006-01-02 15:04:05"),
			status,
			latency,
			clientIP,
			method,
			path,
			userID,
			role,
			userAgent,
			errMsg)
		fmt.Print(logline)

		// write to daily log file
		WriteLogToFile(logline)
	}

}

func WriteLogToFile(logline string) {
	// file name; logs/2025-11-09.Log
	filename := filepath.Join("logs", time.Now().Format("2006-01-02")+".log")

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("log file error: %v\n", err)
		return
	}

	defer file.Close()

	file.WriteString(logline)
}
