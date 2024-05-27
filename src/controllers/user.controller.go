package controllers

import (
	"context"
	"net/http"
	"task_management/src/entity"
	"task_management/src/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

type RegisterUserRequest struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (c *UserController) RegisterUserHandler(ctx *gin.Context) {
	var req RegisterUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	user := &entity.User{
		Email:    req.Email,
		Password: req.Password,
	}

	account := &entity.Account{
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Phone:     req.Phone,
	}

	registeredUser, err := c.userService.Register(context.Background(), user, account)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	ctx.JSON(http.StatusOK, registeredUser)
}
