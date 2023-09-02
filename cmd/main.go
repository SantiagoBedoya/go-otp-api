package main

import (
	"log"
	"net/http"

	"github.com/SantiagoBedoya/otp-api/internal/repository/mysql"
	"github.com/SantiagoBedoya/otp-api/internal/route"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}
	router := gin.Default()

	db, err := mysql.NewMySQLConn()
	if err != nil {
		log.Fatalf("Error connecting to MySQL: %+v", err)
	}

	router.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	route.InitializeUserRoutes(router, db)
	route.InitializeAuthRoutes(router, db)

	router.Run(":8080")
}
