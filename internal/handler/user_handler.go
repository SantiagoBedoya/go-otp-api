package handler

import (
	"log"
	"net/http"

	"github.com/SantiagoBedoya/otp-api/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) GetUsers(ctx *gin.Context) {
	users, err := h.service.GetUsers()
	if err != nil {
		log.Printf("Error getting users: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
		return
	}
	ctx.JSON(http.StatusOK, users)
}
