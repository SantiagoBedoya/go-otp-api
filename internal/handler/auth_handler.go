package handler

import (
	"bytes"
	"errors"
	"image/png"
	"log"
	"net/http"

	"github.com/SantiagoBedoya/otp-api/internal/dto"
	"github.com/SantiagoBedoya/otp-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
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
	user, err := h.service.GetUserByID(userID)
	if err != nil {
		log.Printf("Error getting user by ID: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
		return
	}
	isValid := totp.Validate(data.Code, user.Secret2FA)
	if !isValid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid passcode",
		})
		return
	}

	err = h.service.SetUser2FAValid(userID)
	if err != nil {
		log.Printf("Error setting valid user 2FA: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func (h *AuthHandler) Generate(ctx *gin.Context) {
	user, err := h.service.GetUserByID(ctx.GetString("userID"))
	if err != nil {
		log.Printf("Error getting user by ID: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
		return
	}
	key, err := totp.Generate(totp.GenerateOpts{
		AccountName: user.Email,
		Issuer:      "otp-api",
	})
	if err != nil {
		log.Printf("Error generating totp key: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
		return
	}
	err = h.service.SetUser2FA(ctx.GetString("userID"), key.Secret())
	if err != nil {
		log.Printf("Error saving totp key: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
		return
	}
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		panic(err)
	}
	if err := png.Encode(&buf, img); err != nil {
		log.Printf("Error decode qr code image: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	ctx.Header("Content-Type", "image/png")
	_, _ = ctx.Writer.Write(buf.Bytes())
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
