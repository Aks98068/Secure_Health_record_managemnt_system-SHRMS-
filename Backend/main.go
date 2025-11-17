package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time" // Added missing import

	Configs "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/internal/config"
	Middlewares "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/internal/middlewares"
	Utils "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/internal/utils"
	Routes "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Set to debug mode for development, change to release for production
	gin.SetMode(gin.ReleaseMode) // Changed from ReleaseMode for better debugging

	// Initialize database
	Configs.InitDB()

	// Initialize JWT manager with proper parameters
	jwtManager, err := Utils.NewJWTManager(
		"./certs/jwt_private.pem", // Private key path
		"./certs/jwt_public.pem",  // Public key path
		24*time.Hour,              // Token expiration (24 hours)
	)
	if err != nil {
		log.Fatal("Failed to initialize JWT manager:", err)
	}

	// Create router with optimized middleware order
	router := gin.New()

	// Global middlewares applied in optimal order
	router.Use(gin.Recovery()) // Recovery should be first   // CORS for cross-origin requests
	router.Use(Middlewares.LoggingMiddleware())
	router.Use(Middlewares.CORSMiddleware()) // Logging for all requests

	// Serve frontend static files (fixed typo: was "Fronted")
	router.Static("/assets", "../Frontend/assets")
	router.LoadHTMLGlob("../Frontend/templates/*.html")

	// Public frontend pages
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Secure Health Records Management System",
		})
	})

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "login",
		})
	})

	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"tittle": "register"})
	})

	

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
		})
	})

	// Setup API routes with JWT manager
	Routes.SetupRoutes(router, jwtManager)

	// Server configuration with proper TLS setup
	srv := &http.Server{
		Addr:         "192.168.1.64:443",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start HTTPS server
	go func() {
		log.Println("Starting HTTPS server on https://192.168.1.64:443...")
		// Fixed method name: was "LitenAndServeTLS", should be "ListenAndServeTLS"
		// Fixed certificate paths to be relative
		if err := srv.ListenAndServeTLS("./certs/server.crt", "./certs/server.key"); err != nil && err != http.ErrServerClosed {
			log.Fatal("HTTPS server failed:", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	// Fixed: was "time.second", should be "time.Second"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited gracefully")
}
