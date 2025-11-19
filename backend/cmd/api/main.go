package main

import (
	"log"
	"roadmap/internal/handler"
	"roadmap/internal/handler/middleware"
	"roadmap/internal/infrastructure/database"

	"github.com/gin-gonic/gin"
)

func main() {
	dbConfig := database.NewConfig()

	if err := database.RunMigrations(dbConfig.DSNForMigrate(), "./migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations applied successfully")

	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Database connection established")

	// Create router without default middleware (we'll add our own)
	router := gin.New()

	// Setup middleware (recovery, CORS, logging)
	middleware.SetupMiddleware(router)

	// Routes
	router.GET("/health", handler.HealthHandler)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
