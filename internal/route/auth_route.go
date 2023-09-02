package route

import (
	"database/sql"

	"github.com/SantiagoBedoya/otp-api/internal/handler"
	"github.com/SantiagoBedoya/otp-api/internal/middleware"
	"github.com/SantiagoBedoya/otp-api/internal/repository/mysql"
	"github.com/SantiagoBedoya/otp-api/internal/service"
	"github.com/gin-gonic/gin"
)

func InitializeAuthRoutes(gin *gin.Engine, db *sql.DB) {

	repo := mysql.NewUserRepository(db)
	service := service.NewAuthService(repo)
	handler := handler.NewAuthHandler(service)

	router := gin.Group("/auth")
	{
		router.POST("/sign-up", handler.SignUp)
		router.POST("/sign-in", handler.SignIn)
		router.GET("/otp/generate", middleware.AuthMiddleware, handler.Generate)
		router.POST("/otp/validate", middleware.AuthMiddleware, handler.ValidateOTP)
	}
}
