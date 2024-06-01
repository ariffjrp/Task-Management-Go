package controllers

import (
	"net/http"
	service "task_management/src/services"
	"task_management/src/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) GoogleLogin(ctx *gin.Context) {
	oauthConfig := c.authService.GetGoogleOAuth2Config()
	url := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (c *AuthController) GoogleCallback(ctx *gin.Context) {
	state := ctx.Query("state")
	code := ctx.Query("code")

	user, err := c.authService.HandleGoogleCallback(ctx.Request.Context(), state, code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Buat JWT token
	token, err := utils.GenerateToken(user, c.authService.GetJWTConfig())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
