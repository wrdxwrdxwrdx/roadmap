package main

import (
	"log"
	"roadmap/internal/handler"
	"roadmap/internal/handler/middleware"
	userhandler "roadmap/internal/handler/user"
	"roadmap/internal/infrastructure/database"
	userrepo "roadmap/internal/repository/user"
	userusecase "roadmap/internal/usecase/user"

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

	router := gin.New()

	middleware.SetupMiddleware(router)

	userRepository := userrepo.NewUserRepository(db)

	createUserUseCase := userusecase.NewCreateUserUseCase(userRepository)

	userHandler := userhandler.NewUserHandler(createUserUseCase)

	api := router.Group("/api/v1")
	{
		userhandler.SetupUserRoutes(api, userHandler)
	}

	router.GET("/health", handler.HealthHandler)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
