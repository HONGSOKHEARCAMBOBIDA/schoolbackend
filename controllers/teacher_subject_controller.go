package controllers

import (
	"net/http"
	"schoolbackend/config"
	"schoolbackend/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateTeacherSubject(c *gin.Context) {
	var input models.TeacherSubjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(input.ClassSubjectID) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no classsubject assigned"})
		return
	}

	var teachersubjects []models.TeacherSubject

	tx := config.DB.Begin()

	for _, classsubjectid := range input.ClassSubjectID {
		var existing models.TeacherSubject

		err := tx.Where("user_id = ? AND class_subject_id = ? AND academic_year_id = ?",
			input.UserID, classsubjectid, input.AcademicYearID).
			First(&existing).Error

		if err == nil {
			// Found existing record → update is_active = 0
			existing.IsActive = 0
			if err := tx.Save(&existing).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else if err == gorm.ErrRecordNotFound {
			// No existing record → create new
			teachersubjects = append(teachersubjects, models.TeacherSubject{
				UserID:         input.UserID,
				ClassSubjectID: classsubjectid,
				AcademicYearID: input.AcademicYearID,
				IsActive:       1, // default active
			})
		} else {
			// Other DB error
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Create new records if any
	if len(teachersubjects) > 0 {
		if err := tx.Create(&teachersubjects).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusCreated, gin.H{
		"message": "teacher subjects processed successfully",
		"data":    teachersubjects,
	})
}

func GetClassandSubjectteachbyteacher(c *gin.Context) {
	var result []models.TeacherSubjectDetail
	userID := c.Query("user_id")
	academicyearID := c.Query("academic_year_id")
	classID := c.Query("class_id")
	if err := config.DB.Table("teacher_subjects as ts").
		Select(`ts.id,
		        ts.user_id,
		        ts.academic_year_id,
		        ts.class_subject_id,
				ts.is_active,
		        c.id as class_id,
		        c.name as class_name,
		        s.id as subject_id,
		        s.name as subject_name,
		        a.year_name as academic_year_name,
		        u.name as user_name`).
		Joins("JOIN class_subjects cs ON cs.id = ts.class_subject_id").
		Joins("JOIN academic_years a ON a.id = ts.academic_year_id").
		Joins("JOIN users u ON u.id = ts.user_id").
		Joins("JOIN classes c ON c.id = cs.class_id").
		Joins("JOIN subjects s ON s.id = cs.subject_id").
		Where("ts.user_id = ? AND ts.academic_year_id =? AND cs.class_id =?", userID, academicyearID, classID).
		Scan(&result).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
func GetClassandSubjectNotAssigned(c *gin.Context) {
	var result []models.TeacherSubjectDetail
	userID := c.Query("user_id")

	if err := config.DB.Table("class_subjects as cs").
		Select(`cs.id as class_subject_id,
		        c.id as class_id,
		        c.name as class_name,
		        s.id as subject_id,
		        s.name as subject_name,
		        a.id as academic_year_id,
		        a.year_name as academic_year_name`).
		Joins("JOIN classes c ON c.id = cs.class_id").
		Joins("JOIN subjects s ON s.id = cs.subject_id").
		Joins("JOIN academic_years a ON a.id = cs.academic_year_id"). // must exist in your schema
		Where("cs.id NOT IN (SELECT class_subject_id FROM teacher_subjects WHERE user_id = ?)", userID).
		Scan(&result).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func ChangestatusTeachersubject(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := config.DB.Model(&models.TeacherSubject{}).Where("id =?", id).Update("is_active", gorm.Expr("CASE WHEN is_active = 1 THEN 0 ELSE 1 END "))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, http.StatusOK)
}
func DeleteTeachersubject(c *gin.Context) {
	id := c.Param("id")
	result := config.DB.Delete(&models.TeacherSubject{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Data not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func GetTeachersubjectBYTeacherID(c *gin.Context) {
	var result []models.TeacherSubjectDetail
	userID := c.Query("user_id")
	if err := config.DB.Table("teacher_subjects as ts").
		Select(`ts.id,
		        ts.user_id,
		        ts.academic_year_id,
		        ts.class_subject_id,
				ts.is_active,
		        c.id as class_id,
		        c.name as class_name,
		        s.id as subject_id,
		        s.name as subject_name,
		        a.year_name as academic_year_name,
		        u.name as user_name`).
		Joins("JOIN class_subjects cs ON cs.id = ts.class_subject_id").
		Joins("JOIN academic_years a ON a.id = ts.academic_year_id").
		Joins("JOIN users u ON u.id = ts.user_id").
		Joins("JOIN classes c ON c.id = cs.class_id").
		Joins("JOIN subjects s ON s.id = cs.subject_id").
		Where("ts.user_id = ?", userID).
		Scan(&result).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
