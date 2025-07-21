package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/BOGOMOLOV-ARSENIQ/marketplace/internal/database"
	"github.com/BOGOMOLOV-ARSENIQ/marketplace/internal/handlers"
	"github.com/BOGOMOLOV-ARSENIQ/marketplace/internal/middleware"
	"github.com/BOGOMOLOV-ARSENIQ/marketplace/internal/models"
)

func main() {
	database.InitDB()

	router := gin.Default()

	_ = models.User{}

	router.POST("/register", handlers.RegisterHandler)
	router.POST("/login", handlers.LoginHandler)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Привет от Gin! (после реорганизации и инициализации БД)",
		})
	})

	authenticated := router.Group("/api")
	{
		authenticated.Use(middleware.AuthMiddleware())
		authenticated.GET("/protected", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			c.JSON(http.StatusOK, gin.H{"message": "Вы авторизованы!", "user_id": userID})
		})
	}

	router.Run(":8080")
}
