package main

import (
	"log"
	"os"
	"task_management/src/configs"
	"task_management/src/controllers"
	_ "task_management/src/docs"
	"task_management/src/middleware"
	repository "task_management/src/repositorys"
	"task_management/src/routes"
	"task_management/src/services"
	"task_management/src/utils"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

// @title Task Management Service API
// @version 1.0
// @description A Tag Service API Task Management

// @host	localhost:8080
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

	// Add OAuth2Id column to users table if it doesn't exist
	err = utils.AddOAuth2IdToUsers(db)
	if err != nil {
		log.Fatalf("Failed to add oauth2_id to users: %v", err)
	}

	// Initialize session store with a secure secret key
	middleware.InitSessionStore([]byte(os.Getenv("SECRET_SESSION")))

	// Load JWT config
	jwtConfig := configs.LoadConfigJWT()

	// Load  Google OAuth2 Config
	configs.InitGoogleOAuth2Config()

	// Initialize services
	totpService := services.NewTOTPService()
	emailService := services.NewEmailService(emailConfig, totpService)
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)
	userService := services.NewUserService(userRepo, jwtConfig)
	userController := controllers.NewUserController(userService, totpService, emailService)
	authService := services.NewAuthService(authRepo, jwtConfig)
	authController := controllers.NewAuthController(authService)

	// Set up Gin router
	router := gin.Default()

	// swagger
	router.GET("/v1/api/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Apply session middleware
	router.Use(middleware.SessionMiddleware())

	// Register routes
	routes.RegisterRoutes(router, userController)
	routes.RegisterOAuth2Routes(router, authController)

	// Start the server
	port := os.Getenv("PORT_SERVER")
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
