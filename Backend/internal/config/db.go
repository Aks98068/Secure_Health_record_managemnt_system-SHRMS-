package Configs

import (
	"fmt"
	"log"

	Models "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Println("No .env file found, reading environment variables")
	// }

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		GetEnv("DB_USER", "root"),
		GetEnv("DB_PASSWORD", ""),
		GetEnv("DB_HOST", "127.0.0.1:3306"),
		GetEnv("DB_NAME", "shrms"))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&Models.User{}); err != nil {
		log.Fatal("Failed to migrate databse models:", err)
	}
	DB = db
	log.Println("Database connected and migrate")
}
