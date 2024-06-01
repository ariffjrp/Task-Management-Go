package routes

import (
	"task_management/src/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterOAuth2Routes(router *gin.Engine, authController *controllers.AuthController) {
	router.GET("/v1/api/auth/login/oauth2/google", authController.GoogleLogin)
	router.GET("/v1/api/auth/login/oauth2/code/google", authController.GoogleCallback)
}
