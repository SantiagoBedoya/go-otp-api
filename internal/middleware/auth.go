package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(ctx *gin.Context) {
	bearerToken := ctx.GetHeader("authorization")
	if bearerToken == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Bearer token is required",
		})
		return
	}
	parts := strings.Split(bearerToken, " ")
	if strings.ToLower(parts[0]) != "bearer" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid bearer token",
		})
		return
	}
	token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid bearer token",
		})
		return
	}
	if !token.Valid {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid bearer token",
		})
		return
	}
	userId, err := token.Claims.GetSubject()
	if err != nil {
		log.Printf("Error getting the subject: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	ctx.Set("userID", userId)
	ctx.Next()
}
