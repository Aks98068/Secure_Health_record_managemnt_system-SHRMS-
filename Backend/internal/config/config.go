package Configs

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv          string
	ServerPort      string
	DBHost          string
	DBUser          string
	DBPass          string
	DBName          string
	DBPort          string
	JWTPrivateKey   string
	JWTPublicKey    string
	AccessTokenExp  time.Duration
	RefreshTokenExp time.Duration
	FrontendOrigin  string
	MaxUploadSize   int64
}

func LoadConfig() *Config {
	_ = godotenv.Load("github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/configs/.env")

	cfg := &Config{
		ServerPort:      GetEnv("PORT", "8080"),
		AppEnv:          GetEnv("APP_ENV", "development"),
		DBHost:          GetEnv("DB_HOST", "localhost"),
		DBUser:          GetEnv("DB_USER", "root"),
		DBPass:          GetEnv("DB_PASS", "changeme"),
		DBName:          GetEnv("DB_NAME", "shrms_db"),
		DBPort:          GetEnv("DB_PORT", "3306"),
		JWTPrivateKey:   GetEnv("JWT_PRIVATE_KEY_PATH", "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/certs/jwt_private.pem"),
		JWTPublicKey:    GetEnv("JWT_PUBLIC_KEY_PATH", "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/certs/jwt_public.pem"),
		AccessTokenExp:  time.Duration(getInt("JWT_ACCESS_EXP", 15)) * time.Minute,
		RefreshTokenExp: time.Duration(getInt("JWT_REFRESH_EXP", 7)) * 24 * time.Hour,
		FrontendOrigin:  GetEnv("FRONTEND_ORIGIN", "http://localhost:3000"),
		MaxUploadSize:   int64(getInt("MAX_UPLOAD_MB", 100)) * 1024 * 1024,
	}
	log.Printf("Config loaded for %s environment", cfg.AppEnv)
	return cfg
}

func GetEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		val, err := strconv.Atoi(v)
		if err == nil {
			return val
		}
	}
	return def
}
