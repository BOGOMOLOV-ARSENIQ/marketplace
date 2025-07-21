package handlers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/BOGOMOLOV-ARSENIQ/marketplace/internal/database"
	"github.com/BOGOMOLOV-ARSENIQ/marketplace/internal/models"
	"github.com/BOGOMOLOV-ARSENIQ/marketplace/pkg/utils"
)

var JWTSecret = []byte(os.Getenv("JWT_SECRET"))

func init() {
	if len(JWTSecret) == 0 {
		log.Fatal("JWT_SECRET environment variable is not set or empty")
	}
}

// запрос на регистрацию
func RegisterHandler(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// проверка на сложность
	if !utils.ValidatePasswordComplexity(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least two of the following: letters, digits, or special characters."})
		return
	}

	// хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	// сохраняем пользователя в бд
	result := database.DB.Create(&user)
	if result.Error != nil {
		if database.IsUniqueConstraintError(result.Error) { // уже существует пользователь с таким логином
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user_id": user.ID})
}

// запрос на вход пользователя
func LoginHandler(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// ищем пользователя в БД
	result := database.DB.Where("username = ?", req.Username).First(&user) // Ищем по username
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		}
		return
	}

	// проверяем пароль
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// пароль верный, генерируем JWT-токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	// подписываем токен
	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// успешный вход, даем токен
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": tokenString})
}
