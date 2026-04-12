package controllers

import (
	"net/http"
	"schoolbackend/config"
	"schoolbackend/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateExamComponent(c *gin.Context) {
	var input models.ExamComponentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	examcomponent := models.ExamComponent{
		Name:           input.Name,
		ClassSubjectId: input.ClassSubjectId,
		IsActive:       1,
	}
	if err := config.DB.Create(&examcomponent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, http.StatusCreated)
}
func GetExamComponent(c *gin.Context) {
	// Get query parameters
	classID := c.Query("class_id")
	subjectID := c.Query("subject_id")

	var examComponents []models.ExamComponent

	// Join with class_subject and filter by class_id and subject_id
	if err := config.DB.
		Joins("JOIN class_subjects ON class_subjects.id = exam_components.class_subject_id").
		Where("class_subjects.class_id = ? AND class_subjects.subject_id = ?", classID, subjectID).
		Find(&examComponents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, examComponents)
}

func UpdateExamComponent(c *gin.Context) {
	id := c.Param("id") // get exam component id from URL

	var input models.ExamComponentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find existing exam component
	var examcomponent models.ExamComponent
	if err := config.DB.First(&examcomponent, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "exam component not found"})
		return
	}

	// Update fields
	examcomponent.Name = input.Name

	// keep IsActive as is, or allow update if you want:
	// examcomponent.IsActive = input.IsActive

	if err := config.DB.Save(&examcomponent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "exam component updated successfully",
		"data":    examcomponent,
	})
}
func ChangeStatusExamcomponent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Update with CASE: if 1 → 0, if 0 → 1
	result := config.DB.Model(&models.ExamComponent{}).
		Where("id = ?", id).
		Update("is_active", gorm.Expr("CASE WHEN is_active = 1 THEN 0 ELSE 1 END"))

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(http.StatusOK, http.StatusOK)
}
