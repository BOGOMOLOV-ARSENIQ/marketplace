package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/BOGOMOLOV-ARSENIQ/marketplace/internal/database"
	"github.com/BOGOMOLOV-ARSENIQ/marketplace/internal/handlers"
	"github.com/BOGOMOLOV-ARSENIQ/marketplace/internal/models"
)

func main() {
	database.InitDB()

	router := gin.Default()

	// Пример использования модели (пока неактивно, но для иллюстрации)
	_ = models.User{}

	// Определение маршрутов для аутентификации
	router.POST("/register", handlers.RegisterHandler) // POST запрос на /register будет обрабатывать RegisterHandler
	router.POST("/login", handlers.LoginHandler)       // POST запрос на /login будет обрабатывать LoginHandler

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Привет от Gin! (после реорганизации и инициализации БД)",
		})
	})

	router.Run(":8080")
}
