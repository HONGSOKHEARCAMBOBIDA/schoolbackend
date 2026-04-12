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

func SaveAcademicyear(c *gin.Context) {
	var input models.AcademicYearInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "faild transaction"})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}

	}()
	var academicyear models.AcademicYear
	if input.ID == 0 {
		// create acedemicyear
		academicyear = models.AcademicYear{
			YearName: input.YearName,
			IsActive: 1,
		}
		if err := tx.Create(&academicyear).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, http.StatusCreated)
	} else {
		// update acedemicyear
		if err := tx.First(&academicyear, input.ID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusFound, gin.H{"error": err.Error()})
			return
		}
		academicyear.YearName = input.YearName
		if err := tx.Save(&academicyear).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, http.StatusOK)

	}
	tx.Commit()
}
func Handleacademicyear(c *gin.Context) {
	method := c.Request.Method

	if method == http.MethodGet {
		var academicyear []models.AcademicYear

		userID, ok := helper.GetUserID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "login!"})
			return
		}

		var currenUser struct {
			Role int64 `gorm:"column:role_id"`
		}
		if err := config.DB.Table("users").Select("role_id").Where("id = ?", userID).Scan(&currenUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		db := config.DB.Model(&models.AcademicYear{})
		if currenUser.Role != 1 {
			db = db.Where("is_active = ?", 1)
		}

		if err := db.Order("id desc").Find(&academicyear).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, academicyear)

	} else if method == http.MethodPut {
		// ----- TOGGLE STATUS -----
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		result := config.DB.Model(&models.AcademicYear{}).
			Where("id = ?", id).
			Update("is_active", gorm.Expr("CASE WHEN is_active = 1 THEN 0 ELSE 1 END"))

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "AcademicYear model not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})

	} else {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Only GET and PUT are allowed"})
	}
}
