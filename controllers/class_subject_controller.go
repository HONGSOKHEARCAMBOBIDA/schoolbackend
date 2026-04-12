package controllers

import (
	"net/http"
	"schoolbackend/config"
	"schoolbackend/models"

	"github.com/gin-gonic/gin"
)

func DeleteClassSubjectByID(c *gin.Context) {
	id := c.Param("id")

	result := config.DB.Delete(&models.ClassSubject{}, id)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "ClassSubject not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ClassSubject deleted successfully"})
}

// CreateClassSubject creates multiple class-subject relationships
func CreateClassSubject(c *gin.Context) {
	var input models.ClassSubjectInput

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if SubjectID is empty
	if len(input.SubjectID) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No subjects provided"})
		return
	}

	var classsubjects []models.ClassSubject

	// Prevent duplicates
	for _, subjectid := range input.SubjectID {
		var exists int64
		config.DB.Model(&models.ClassSubject{}).
			Where("class_id = ? AND subject_id = ?", input.ClassID, subjectid).
			Count(&exists)
		if exists == 0 {
			classsubjects = append(classsubjects, models.ClassSubject{
				ClassID:   input.ClassID,
				SubjectID: subjectid,
			})
		}
	}

	// Check if there is anything to insert
	if len(classsubjects) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "All subjects already assigned"})
		return
	}

	// Use transaction for safety
	tx := config.DB.Begin()
	if err := tx.Create(&classsubjects).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tx.Commit()

	c.JSON(http.StatusCreated, classsubjects)
}
func GetClassSubjectsNotAssigntoTeacher(c *gin.Context) {
	var result []models.ClassSubjectDetail

	userID := c.Query("user_id")
	classID := c.Query("class_id")
	academicyearID := c.Query("academic_year_id")

	if err := config.DB.Table("class_subjects as cs").
		Select("cs.id, cs.subject_id, s.name as subject_name, cs.class_id, c.name as class_name").
		Joins("JOIN subjects s ON s.id = cs.subject_id").
		Joins("JOIN classes c ON c.id = cs.class_id").
		Where("cs.class_id = ?", classID).
		Where("cs.id NOT IN (SELECT class_subject_id FROM teacher_subjects WHERE user_id = ? AND academic_year_id = ? )", userID, academicyearID).
		Find(&result).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}

func GetClassSubjects(c *gin.Context) {

	var results []models.ClassSubjectDetail
	classID := c.Query("class_id")

	db := config.DB.Table("class_subjects as cs").
		Select("cs.id, cs.class_id, c.name as class_name, cs.subject_id, s.name as subject_name").
		Joins("left join classes c on c.id = cs.class_id").
		Joins("left join subjects s on s.id = cs.subject_id")

	if classID != "" {
		db = db.Where("cs.class_id = ?", classID)
	}

	if err := db.Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
func GetClassSubjectsNOtinexamecompone(c *gin.Context) {
	var results []models.ClassSubjectDetail
	classID := c.Query("class_id")

	db := config.DB.Table("class_subjects as cs").
		Select("cs.id, cs.class_id, c.name as class_name, cs.subject_id, s.name as subject_name").
		Joins("left join classes c on c.id = cs.class_id").
		Joins("left join subjects s on s.id = cs.subject_id").
		Joins("left join exam_components ec on ec.class_subject_id = cs.id").
		Where("ec.class_subject_id IS NULL") // only get class_subjects not in exam_component

	if classID != "" {
		db = db.Where("cs.class_id = ?", classID)
	}

	if err := db.Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func GetSubjectNotAssignedToClass(c *gin.Context) {
	var result []models.SubjectNotAssign

	classID := c.Query("class_id")

	// if no class_id provided
	if classID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "class_id is required"})
		return
	}

	// Query subjects not yet assigned to this class
	if err := config.DB.Table("subjects as s").
		Select("s.id as subject_id, s.name as subject_name").
		Where("s.id NOT IN (?)",
			config.DB.Table("class_subjects").Select("subject_id").Where("class_id = ?", classID),
		).
		Scan(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
