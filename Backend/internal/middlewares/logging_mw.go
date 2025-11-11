package Middlewares

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware writes logs to terminal and file with rotation
func LoggingMiddleware() gin.HandlerFunc {
	CreateFolder() // ensure "logs" folder exists

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// End timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		status := c.Writer.Status()
		userAgent := c.Request.UserAgent()
		errMsg := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Optional: extract authenticated user info (from JWT middleware)
		userID := c.GetString("user_id")
		role := c.GetString("role")

		if raw != "" {
			path = path + "?" + raw
		}

		// Log format
		logline := fmt.Sprintf("[%s] | %3d | %13v | %-15s | %-6s %s | user:%s | role:%s | UA:%s | err:%s\n",
			end.Format("2006-01-02 15:04:05"),
			status,
			latency,
			clientIP,
			method,
			path,
			userID,
			role,
			userAgent,
			errMsg,
		)

		// Print to terminal
		fmt.Print(logline)

		// Write to daily log file
		WriteLogToFile(logline)
	}
}

func WriteLogToFile(logline string) {
	// File name: logs/2025-11-09.log
	filename := filepath.Join("logs", time.Now().Format("2006-01-02")+".log")

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Log file error: %v\n", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(logline); err != nil {
		fmt.Printf("Failed to write log: %v\n", err)
	}
}

// CreateFolder ensures "logs" directory exists
func CreateFolder() {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", os.ModePerm)
	}
}
