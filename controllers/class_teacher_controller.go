package controllers

import (
	"net/http"
	"schoolbackend/config"
	"schoolbackend/models"

	"github.com/gin-gonic/gin"
)

func CreateClassTeacher(c *gin.Context) {
	var input models.ClassTeacherInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var existing models.ClassTeacher
	err := config.DB.
		Where("class_id = ? AND academic_year_id = ? AND is_active = ?", input.ClassID, input.AcademicYearID, 1).
		First(&existing).Error

	if err == nil {
		// ❗ Already exists
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This class already has a teacher assigned for this academic year",
		})
		return
	}
	newClassTeacher := models.ClassTeacher{
		ClassID:        input.ClassID,
		AcademicYearID: input.AcademicYearID,
		TeacherID:      input.TeacherID,
		IsActive:       1,
	}

	if err := config.DB.Create(&newClassTeacher).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Class teacher created successfully", "data": newClassTeacher})
}
func UpdatestatusClassTeacher(c *gin.Context) {
	id := c.Param("id")
	var classteacher models.ClassTeacher

	// Find record
	if err := config.DB.Where("id = ?", id).First(&classteacher).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Class teacher not found"})
		return
	}

	// Update status
	classteacher.IsActive = 0

	if err := config.DB.Save(&classteacher).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status updated successfully",
		"data":    classteacher,
	})
}
