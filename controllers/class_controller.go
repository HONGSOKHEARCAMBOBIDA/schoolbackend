package controllers

import (
	"net/http"
	"schoolbackend/config"
	"schoolbackend/helper"
	"schoolbackend/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SaveClass(c *gin.Context) {
	var input models.ClassInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start transaction"})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var class models.Class
	if input.ID == 0 {
		// --- Create new class ---
		class = models.Class{
			Name:     input.Name,
			IsActive: 1,
		}
		if err := tx.Create(&class).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "class created successfully", "class": class})
	} else {
		// --- Update existing class ---
		if err := tx.First(&class, input.ID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "class not found"})
			return
		}

		class.Name = input.Name

		if err := tx.Save(&class).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "class updated successfully", "class": class})
	}

	tx.Commit()
}

func Handleclass(c *gin.Context) {
	method := c.Request.Method

	if method == http.MethodGet {
		userID, ok := helper.GetUserID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "login!"})
		}
		var currenUser struct {
			Role int64 `gorm:"column:role_id"`
		}
		if err := config.DB.Table("users").Select("role_id").Where("id =?", userID).Scan(&currenUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// ----- GET ALL -----
		var class []models.Class
		dbQuery := config.DB.Model(&models.Class{})
		if currenUser.Role != 1 {
			dbQuery = dbQuery.Where("is_active = ?", 1)
		}

		if err := dbQuery.Find(&class).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, class)

	} else if method == http.MethodPut {
		// ----- TOGGLE STATUS -----
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		// Update CASE: if 1 → 0, if 0 → 1
		result := config.DB.Model(&models.Class{}).
			Where("id = ?", id).
			Update("is_active", gorm.Expr("CASE WHEN is_active = 1 THEN 0 ELSE 1 END"))

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Moto model not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
	} else {
		// If method not allowed
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Only GET and PUT are allowed"})
	}
}
