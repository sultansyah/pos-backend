package main

import (
	"log"
	"os"
	"post-backend/internal/config"
	"post-backend/internal/token"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	dbConfig := config.DBConfig{
		User:     dbUser,
		Password: dbPassword,
		Host:     dbHost,
		Port:     dbPort,
		Name:     dbName,
	}

	db, err := config.InitDb(dbConfig)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5500"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.Static("/static/images", "./public/images")

	api := router.Group("/api/v1")

	tokenService := token.NewTokenService([]byte(jwtSecretKey))

	if err := router.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
