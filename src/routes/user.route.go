package routes

import (
	"task_management/src/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, userController *controllers.UserController) {
	router.POST("/v1/api/auth/register", userController.RegisterUserHandler)
	router.POST("/v1/api/auth/verify-otp", userController.VerifyOTPHandler)
	router.POST("/v1/api/auth/login", userController.Login)
}
