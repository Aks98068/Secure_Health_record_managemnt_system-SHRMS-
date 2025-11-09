package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	Configs "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/internal/config"
	Middlewares "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/internal/middlewares"
	Models "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	// production mode
	gin.SetMode(gin.ReleaseMode)

	// initializing database
	Configs.InitDB()

	router := gin.New()

	// security middlewares and global middlewares
	router.Use(gin.Recovery())
	router.Use(Middlewares.CORSMiddleware())
	router.Use(Middlewares.LoggingMiddleware())

	// serve frontend static files
	router.Static("/assets", "./Fronted/assets")
	router.LoadHTMLGlob("./Frontend/templates/*.html")
	

	// public frontend pages and home page also
	router.GET("/", func(c *gin.Context){
		c.HTML(200, "index.html", nil)
	})



	//  secure API Routes
	router.SetupRoutes(router)

	srv := &http.Server{
		Addr:    ":443",
		Handler: router,
	}

	go func() {
		log.Println("starting HTTPS Server on port 443")
		if err := srv.LitenAndServeTLS("github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/Backend/certs/server.crt", "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/Backend/certs/server.key"); err != nil && err != http.ErrServerClosed {
			log.Fatal("HTTPS server failed;", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.second)
	defer cancel()
	if err := srv.shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}
	log.Println("server exited gracefully")

}
