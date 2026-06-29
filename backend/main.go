package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"careerhubv2-backend/config"
	"careerhubv2-backend/middleware"
	"careerhubv2-backend/handlers"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading environment variables from OS")
	}

	config.InitDB()

	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(middleware.CORS())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", handlers.Health)

		auth := v1.Group("/auth")
		auth.Use(middleware.RateLimit())
		{
			auth.POST("/login/microsoft", handlers.LoginMicrosoft)
			auth.POST("/refresh", handlers.RefreshToken)
		}

		protected := v1.Group("/")
		protected.Use(middleware.Auth())
		{
			protected.GET("/users/me", handlers.GetMe)
			protected.GET("/categories", handlers.GetCategories)
			protected.GET("/jobs", handlers.GetJobs)
			protected.GET("/jobs/:id", handlers.GetJobByID)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	log.Printf("CareerHubV2 backend starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
