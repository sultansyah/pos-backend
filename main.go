package main

import (
	"log"
	"os"
	"post-backend/internal/category"
	"post-backend/internal/config"
	"post-backend/internal/middleware"
	"post-backend/internal/product"
	"post-backend/internal/setting"
	"post-backend/internal/token"
	"post-backend/internal/user"
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

	userRepository := user.NewUserRepository()
	userService := user.NewUserService(db, userRepository, tokenService)
	userHandler := user.NewUserHandler(userService)

	settingRepository := setting.NewSettingRepository()
	settingService := setting.NewSettingService(settingRepository, db)
	settingHandler := setting.NewSettingHandler(settingService)

	categoryRepository := category.NewCategoryRepository()
	categoryService := category.NewCategoryService(categoryRepository, db)
	categoryHandler := category.NewCategoryHandler(categoryService)

	productRepository := product.NewProductRepository()
	productService := product.NewProductService(productRepository, db)
	productHandler := product.NewProductHandler(productService)

	api.POST("/auth/login", userHandler.Login)
	api.POST("/auth/password", middleware.AuthMiddleware(tokenService), userHandler.UpdatePassword)

	api.GET("/settings", middleware.AuthMiddleware(tokenService), settingHandler.GetAll)

	api.POST("/categories", middleware.AuthMiddleware(tokenService), middleware.RoleMiddleware([]string{"admin"}), categoryHandler.Create)
	api.GET("/categories", categoryHandler.GetAll)
	api.GET("/categories/:id", categoryHandler.Get)
	api.DELETE("/categories/:id", middleware.AuthMiddleware(tokenService), middleware.RoleMiddleware([]string{"admin"}), categoryHandler.Delete)
	api.PUT("/categories/:id", middleware.AuthMiddleware(tokenService), middleware.RoleMiddleware([]string{"admin"}), categoryHandler.Update)

	api.GET("/products", productHandler.GetAll)
	api.GET("/products/id/:id", productHandler.Get)
	api.GET("/products/slug/:slug", productHandler.GetBySlug)
	api.POST("/products", middleware.AuthMiddleware(tokenService), middleware.RoleMiddleware([]string{"admin"}), productHandler.Insert)
	api.PUT("/products/:id", middleware.AuthMiddleware(tokenService), middleware.RoleMiddleware([]string{"admin"}), productHandler.Update)
	api.DELETE("/products/:id", middleware.AuthMiddleware(tokenService), middleware.RoleMiddleware([]string{"admin"}), productHandler.Delete)

	api.POST("/products/:id/images", middleware.AuthMiddleware(tokenService), middleware.RoleMiddleware([]string{"admin"}), productHandler.InsertImage)
	api.PUT("/products/:id/images/:imageId", middleware.AuthMiddleware(tokenService), middleware.RoleMiddleware([]string{"admin"}), productHandler.SetLogoImage)
	api.DELETE("/products/:id/images/:imageId", middleware.AuthMiddleware(tokenService), middleware.RoleMiddleware([]string{"admin"}), productHandler.DeleteImage)

	if err := router.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
