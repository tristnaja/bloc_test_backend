package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tristnaja/bloc_test_backend/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	log.SetPrefix("bcc_ws: ")
	log.SetFlags(log.LstdFlags)

	err := loadEnv(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Construct DSN (Data Source Name) from Env
	// Example format: user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(&model.User{})

	authController := &AuthController{DB: db}

	err, authMiddleware := AuthMiddleware()
	if err != nil {
		log.Fatalf("Failed to initialize middleware: %v", err)
	}

	r := gin.Default()

	// --- Public Routes ---
	r.POST("/register", authController.Register)
	r.POST("/login", authController.Login)

	// --- Protected Routes ---
	// Create a group that uses the JWT middleware
	api := r.Group("/api")
	api.Use(authMiddleware)
	{
		api.GET("/profile", authController.GetProfile)
		// Add more protected routes here
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on port %s\n", port)
	r.Run(":" + port)
}
