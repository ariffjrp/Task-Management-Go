package controllers

import (
	"context"
	"log"
	"net/http"
	"task_management/src/entity"
	"task_management/src/middleware"
	"task_management/src/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService  services.UserService
	totpService  *services.TOTPService
	emailService *services.EmailService
}

func NewUserController(userService services.UserService, totpService *services.TOTPService, emailService *services.EmailService) *UserController {
	return &UserController{
		userService:  userService,
		totpService:  totpService,
		emailService: emailService,
	}
}

type RegisterUserRequest struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type VerifyOTPRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Create Tags	registers a new
// @Summary 	Register a new user
// @Description Register a new user with email verification
// @Tags 		User
// @Accept 		json
// @Produce 	json
// @Param 		request body RegisterUserRequest true "Registration Request"
// @Success 	200 {string} string "success"
// @Failure 	400 {object} string "Invalid request payload"
// @Failure 	500 {object} string "Internal server error"
// @Router 		/v1/api/auth/register [post]
func (c *UserController) RegisterUserHandler(ctx *gin.Context) {
	session := middleware.GetSession(ctx)
	var req RegisterUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	secret, err := c.totpService.GenerateSecret()
	if err != nil {
		log.Printf("Failed to generate TOTP secret: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate TOTP secret"})
		return
	}

	otp, err := c.totpService.GenerateOTP(secret)
	if err != nil {
		log.Printf("Failed to generate OTP: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
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

	session.Values["user"] = user
	session.Values["account"] = account
	session.Values["secret"] = secret

	if err := session.Save(ctx.Request, ctx.Writer); err != nil {
		log.Printf("Failed to save session: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	if err := c.emailService.SendVerificationEmail(req.Email, otp); err != nil {
		log.Printf("Failed to send verification email: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "OTP has been sent to your email",
	})
}

// Create Tags	Verifikasi
// @Summary Verify OTP
// @Description Verify OTP for user registration
// @Tags User
// @Accept json
// @Produce json
// @Param request body VerifyOTPRequest true "OTP Verification Request"
// @Success 200 {object} string "success"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 401 {object} string "Invalid OTP"
// @Failure 500 {object} string "Internal server error"
// @Router /v1/api/auth/verify-otp [post]
func (c *UserController) VerifyOTPHandler(ctx *gin.Context) {
	session := middleware.GetSession(ctx)
	var req VerifyOTPRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	secret, ok := session.Values["secret"].(string)
	if !ok {
		log.Printf("No secret found in session")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No secret found in session"})
		return
	}

	if err := c.totpService.VerifyOTP(secret, req.OTP); err != nil {
		log.Printf("Invalid OTP: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}

	user, ok := session.Values["user"].(*entity.User)
	if !ok {
		log.Printf("No user data found in session")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No user data found in session"})
		return
	}

	account, ok := session.Values["account"].(*entity.Account)
	if !ok {
		log.Printf("No account data found in session")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No account data found in session"})
		return
	}

	// Simpan data pengguna dan akun ke database
	registeredUser, registeredAccount, accessToken, refreshToken, err := c.userService.Register(ctx, user, account)
	if err != nil {
		log.Printf("Failed to register user: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	// Hapus data dari sesi setelah penyimpanan berhasil
	session.Values["user"] = nil
	session.Values["account"] = nil
	session.Values["secret"] = nil
	if err := session.Save(ctx.Request, ctx.Writer); err != nil {
		log.Printf("Failed to clear session: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear session"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "User verified and registered successfully",
		"user":         registeredUser,
		"account":      registeredAccount,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

// Create Tags	login
// @Summary Login a user
// @Description Login a user with email and password
// @Tags User
// @Accept json
// @Produce json
// @Param request body entity.UserLogin true "Login Request"
// @Success 200 {object} string "success"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 401 {object} string "Unauthorized"
// @Router /v1/api/auth/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var loginRequest entity.UserLogin
	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := c.userService.Login(context.Background(), &loginRequest)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
