package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var JWTSecret = []byte(os.Getenv("JWT_SECRET"))

func init() {
	if len(JWTSecret) == 0 {
		log.Fatal("JWT_SECRET environment variable is not set or empty")
	}
}

// генерирует токен для пользователя
func GenerateToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// извлекаем заголовок
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// разделяем на "Bearer" и токен
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}
		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// проверяем метод подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return JWTSecret, nil // секретный ключ для проверки подписи
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token", "details": err.Error()})
			return
		}

		// извлекаем клеймы данные из токна
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// получаем и сохраняем user_id
		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
			return
		}
		c.Set("userID", uint(userID)) // Сохраняем ID пользователя в контексте Gin для последующих обработчиков

		// идем дальше по цепочке middleware/обработчиков
		c.Next()
	}
}
