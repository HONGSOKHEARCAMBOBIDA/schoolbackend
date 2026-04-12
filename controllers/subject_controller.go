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

func SaveSubject(c *gin.Context) {
	var input models.SubjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "faild"})
		return
	}
	defer func() {
		//defer means this function will run at the end of the request (when function exits).
		if r := recover(); r != nil {
			//	recover() is used to catch panics (unexpected crashes).
			tx.Rollback()
		}

	}()
	var subject models.Subject
	if input.ID == 0 {
		subject = models.Subject{
			Name:     input.Name,
			IsActive: 1,
		}
		if err := tx.Create(&subject).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.SecureJSON(http.StatusCreated, subject)
	} else {
		if err := tx.First(&subject, input.ID).Error; err != nil {
			tx.Rollback()
			c.SecureJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		subject.Name = input.Name
		subject.IsActive = 1
		if err := tx.Save(&subject).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "updated"})

	}
	tx.Commit()
}
func Getsubject(c *gin.Context) {
	var subjects []models.Subject

	userID, ok := helper.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login!"})
		return // 🔹 need to return here
	}

	var currenUser struct {
		Role int64 `gorm:"column:role_id"`
	}
	if err := config.DB.Table("users").Select("role_id").Where("id = ?", userID).Scan(&currenUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	db := config.DB.Model(&models.Subject{})
	if currenUser.Role != 1 {
		db = db.Where("is_active = ?", 1) // assign to outer db
	}

	if err := db.Find(&subjects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subjects)
}

func ChangestatusSubject(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := config.DB.Model(&models.Subject{}).Where("id =?", id).Update("is_active", gorm.Expr("CASE WHEN is_active =1 THEN 0 ELSE 1 END"))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Faild"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, http.StatusOK)
}
