package route

import (
	"database/sql"

	"github.com/SantiagoBedoya/otp-api/internal/handler"
	"github.com/SantiagoBedoya/otp-api/internal/middleware"
	"github.com/SantiagoBedoya/otp-api/internal/repository/mysql"
	"github.com/SantiagoBedoya/otp-api/internal/service"
	"github.com/gin-gonic/gin"
)

func InitializeUserRoutes(gin *gin.Engine, db *sql.DB) {
	repo := mysql.NewUserRepository(db)
	service := service.NewUserService(repo)
	handler := handler.NewUserHandler(service)

	router := gin.Group("/users", middleware.AuthMiddleware)
	{
		router.GET("", handler.GetUsers)
	}
}
