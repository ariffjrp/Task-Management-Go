package main

import (
	"task_management/src/configs"
	"task_management/src/controllers"
	"task_management/src/repositorys"
	"task_management/src/routes"
	"task_management/src/services"
	"task_management/src/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	router := gin.Default()

	// Initialize DB connection
	db, err := configs.ConnectDB()
	if err != nil {
		panic("failed to connect to the database")
	}

	// Migrate the schema
	utils.MigrateDB(db)

	// Initialize repositories and services
	userRepo := repositorys.NewUserRepository(db)
	userService := services.NewUserService(userRepo)

	// Initialize controllers
	userController := controllers.NewUserController(userService)

	// Register user routes
	routes.RegisterRoutes(router, userController)

	// Start the server
	router.Run(":8080")
}
