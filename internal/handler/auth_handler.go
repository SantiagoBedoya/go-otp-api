package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/SantiagoBedoya/otp-api/internal/dto"
	"github.com/SantiagoBedoya/otp-api/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) ValidateOTP(ctx *gin.Context) {
	var data dto.OTPDto
	if err := ctx.ShouldBindJSON(&data); err != nil {
		log.Printf("Error binding otp data: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid data",
		})
		return
	}
	userID := ctx.GetString("userID")
	token, err := h.service.Validate2FA(userID, &data)
	if err != nil {
		log.Printf("Error validating 2FA: %v", err)
		if errors.Is(err, service.ErrInvalidPasscode) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"access_token": token,
	})
}

func (h *AuthHandler) Generate(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	user, err := h.service.GetUserByID(userID)
	if err != nil {
		log.Printf("Error getting user by ID: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
		return
	}
	imageBytes, err := h.service.Setup2FA(userID, user.Email)
	if err != nil {
		log.Printf("Error doing the setup of 2FA: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	ctx.Header("Content-Type", "image/png")
	_, _ = ctx.Writer.Write(imageBytes)
}

func (h *AuthHandler) SignIn(ctx *gin.Context) {
	var data dto.SignInDto
	if err := ctx.ShouldBindJSON(&data); err != nil {
		log.Printf("Error binding signIn data: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid data",
		})
		return
	}
	token, err := h.service.SignIn(&data)
	if err != nil {
		log.Printf("Error doing signIn: %v", err)
		if errors.Is(err, service.ErrInvalidPassword) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}
		if errors.Is(err, service.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"access_token": token,
	})
}

func (h *AuthHandler) SignUp(ctx *gin.Context) {
	var data dto.SignUpDto
	if err := ctx.ShouldBindJSON(&data); err != nil {
		log.Printf("Error binding signUp data: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid data",
		})
		return
	}
	user, err := h.service.SaveUser(&data)
	if err != nil {
		log.Printf("Error saving user: %v", err)
		if errors.Is(err, service.ErrEmailInUse) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
		return
	}
	ctx.JSON(http.StatusCreated, user)
}
