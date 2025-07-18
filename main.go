package main

import (
	"net/http"

	"github.com/BOGOMOLOV-ARSENIQ/marketplace/internal/models"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Пример использования модели (пока неактивно, но для иллюстрации)
	_ = models.User{} // Просто для того, чтобы импорт models был использован и не вызывал ошибку компиляции

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Привет от Gin! (после реорганизации)",
		})
	})

	router.Run(":8080")
}
