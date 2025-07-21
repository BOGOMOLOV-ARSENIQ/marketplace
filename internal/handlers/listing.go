package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/BOGOMOLOV-ARSENIQ/marketplace/internal/database"
	"github.com/BOGOMOLOV-ARSENIQ/marketplace/internal/models"
)

// создает новое объявление
func CreateListingHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	var req models.Listing
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.UserID = userID.(uint)

	result := database.DB.Create(&req)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create listing", "details": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, req) // возвращаем созданное объявление
}

// получает все объявления
// не требует авторизации
func GetListingsHandler(c *gin.Context) {
	var listings []models.Listing         // слайс для хранения объявлений
	result := database.DB.Find(&listings) // запрос для выдачи всех объявлений
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve listings"})
		return
	}
	c.JSON(http.StatusOK, listings)
}

// получает одно объявление по ID
// не требует авторизации
func GetListingByIDHandler(c *gin.Context) {
	id := c.Param("id") // получаем ID объявления

	var listing models.Listing
	result := database.DB.First(&listing, id) // запрашивает обьявление по id
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve listing"})
		}
		return
	}
	c.JSON(http.StatusOK, listing)
}

// обновляет объявление по ID
// требует авторизации и проверки, что пользователь владеет объявлением
func UpdateListingHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	id := c.Param("id")
	listingID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid listing ID"})
		return
	}

	var existingListing models.Listing
	result := database.DB.First(&existingListing, listingID) // запрашиваем объявление
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve listing"})
		}
		return
	}

	// проверяем владельца объявления
	if existingListing.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this listing"})
		return
	}

	var updatedFields models.Listing // поля, которые мы хотим обновить
	if err := c.ShouldBindJSON(&updatedFields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingListing.Title = updatedFields.Title
	existingListing.Description = updatedFields.Description
	existingListing.Price = updatedFields.Price
	existingListing.ImageURL = updatedFields.ImageURL

	result = database.DB.Save(&existingListing) // обновляем запись в бд
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update listing"})
		return
	}

	c.JSON(http.StatusOK, existingListing) // возвращаем объявление
}

// удаляет объявление по ID
// требует авторизации и проверки, что пользователь владеет объявлением
func DeleteListingHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	id := c.Param("id")
	listingID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid listing ID"})
		return
	}

	var existingListing models.Listing
	result := database.DB.First(&existingListing, listingID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve listing"})
		}
		return
	}

	// проверяем владельца объявления
	if existingListing.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this listing"})
		return
	}

	result = database.DB.Delete(&existingListing) // удаляем объявление
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete listing"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Listing deleted successfully"})
}
