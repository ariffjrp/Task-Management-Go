package main

import (
	"log"
	"os"
	"task_management/src/configs"
	"task_management/src/controllers"
	"task_management/src/middleware"
	repository "task_management/src/repositorys"
	"task_management/src/routes"
	"task_management/src/services"
	"task_management/src/utils"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env Mail")
	}

	// Load environment variables and configurations
	db, err := configs.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	emailConfig := configs.LoadConfigEmail()

	// Migrate database
	utils.MigrateDB(db)

	// Initialize session store with a secure secret key
	middleware.InitSessionStore([]byte(os.Getenv("SECRET_SESSION")))

	// Load JWT config
	jwtConfig := configs.LoadConfigJWT()

	// Initialize services
	totpService := services.NewTOTPService()
	emailService := services.NewEmailService(emailConfig, totpService)
	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo, jwtConfig)
	userController := controllers.NewUserController(userService, totpService, emailService)

	// Set up Gin router
	router := gin.Default()

	// Apply session middleware
	router.Use(middleware.SessionMiddleware())

	// Register routes
	routes.RegisterRoutes(router, userController)

	// Start the server
	port := os.Getenv("PORT_SERVER")
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
