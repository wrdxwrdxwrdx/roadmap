package main

import (
	"log"
	"os"
	"roadmap/internal/handler"
	"roadmap/internal/handler/middleware"
	userhandler "roadmap/internal/handler/user"
	"roadmap/internal/infrastructure/database"
	jwtservice "roadmap/internal/pkg/jwt"
	userrepo "roadmap/internal/repository/user"
	userusecase "roadmap/internal/usecase/user"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func initDatabase() *database.Database {
	dbConfig := database.NewConfig()

	if err := database.RunMigrations(dbConfig.DSNForMigrate(), "./migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations applied successfully")

	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established")
	return db
}

func initJWT() *jwtservice.JWTService {
	secretKey := getEnv("JWT_SECRET_KEY", "your-secret-key-change-in-production")
	expiresInHoursStr := getEnv("JWT_EXPIRES_IN_HOURS", "24")

	expiresInHours, err := strconv.Atoi(expiresInHoursStr)
	if err != nil || expiresInHours <= 0 {
		log.Printf("Invalid JWT_EXPIRES_IN_HOURS value '%s', using default 24 hours", expiresInHoursStr)
		expiresInHours = 24
	}

	expiresIn := time.Duration(expiresInHours) * time.Hour

	log.Printf("JWT service initialized with expiration: %d hours", expiresInHours)
	return jwtservice.NewJWTService(secretKey, expiresIn)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	db := initDatabase()
	defer db.Close()

	router := gin.New()

	middleware.SetupMiddleware(router)

	userRepository := userrepo.NewUserRepository(db)

	jwtService := initJWT()

	createUserUseCase := userusecase.NewCreateUserUseCase(userRepository)
	registerUseCase := userusecase.NewRegisterUseCase(userRepository, jwtService)
	loginUseCase := userusecase.NewLoginUseCase(userRepository, jwtService)

	userHandler := userhandler.NewUserHandler(createUserUseCase, registerUseCase, loginUseCase)

	authMiddleware := middleware.AuthMiddleware(jwtService)

	api := router.Group("/api/v1")
	{
		api.GET("/health", handler.HealthHandler)
		userhandler.SetupUserRoutes(api, userHandler, authMiddleware)
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
