package controllers

import (
	"net/http"
	"schoolbackend/config"
	"schoolbackend/helper"
	"schoolbackend/models"

	"github.com/gin-gonic/gin"
)

func CreateStudentClass(c *gin.Context) {
	var input models.StudentClassInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Step 1: Find the last StudentClass for this student
	var last models.StudentClass
	if err := config.DB.
		Where("student_id = ?", input.StudentID).
		Order("id desc").
		First(&last).Error; err == nil {
		// Step 2: Update status to 0 if found
		last.IsActive = 0
		config.DB.Save(&last)
	}

	// Step 3: Create new record with status = 1
	newStudentClass := models.StudentClass{
		StudentID:      input.StudentID,
		ClassID:        input.ClassID,
		IsActive:       1,
		AcademicYearID: input.AcademicYearID,
	}

	if err := config.DB.Create(&newStudentClass).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"student_class": newStudentClass})
}

// PUT /studentclass/:id
func UpdateStudentClass(c *gin.Context) {
	id := c.Param("id")

	var input models.StudentClassInput
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

	var studentClass models.StudentClass
	if err := tx.First(&studentClass, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	studentClass.ClassID = input.ClassID

	if err := tx.Save(&studentClass).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated", "student_class": studentClass})
}

func GetStudentClassByStudentID(c *gin.Context) {
	studentID := c.Param("id") // or c.Query("student_id")

	var detail models.StudentClassDetail

	err := config.DB.Table("student_classes as sc").
		Select(`sc.id, 
                sc.student_id, s.name as student_name, 
                sc.class_id, c.name as class_name, 
                sc.academic_year_id, ay.year_name as academic_year_name, 
                sc.is_active as is_active`).
		Joins("left join students s on s.id = sc.student_id").
		Joins("left join classes c on c.id = sc.class_id").
		Joins("left join academic_years as ay on ay.id = sc.academic_year_id").
		Where("sc.student_id = ?", studentID).
		Order("sc.id desc").
		First(&detail).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No student class found"})
		return
	}

	c.JSON(http.StatusOK, detail)
}

func GetStudentClassbyClassIDandAcademicyearID(c *gin.Context) {
	var studentclass []models.StudentClassDetail
	academicYearID := c.Query("academic_year_id")
	classID := c.Query("class_id")

	userID, ok := helper.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var currenUser struct {
		ManageClass int
	}
	if err := config.DB.Table("users").
		Select("manage_class").
		Where("id = ?", userID).
		Scan(&currenUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	db := config.DB.Table("student_classes as sc").
		Select(`sc.id, 
                sc.student_id, s.name as student_name, 
                sc.class_id, c.name as class_name, 
                sc.academic_year_id, ay.year_name as academic_year_name, 
                sc.is_active as is_active`).
		Joins("LEFT JOIN students s ON s.id = sc.student_id").
		Joins("LEFT JOIN classes c ON c.id = sc.class_id").
		Joins("LEFT JOIN academic_years ay ON ay.id = sc.academic_year_id")

	if academicYearID != "" {
		db = db.Where("sc.academic_year_id = ?", academicYearID)
	}
	if classID != "" {
		db = db.Where("sc.class_id = ?", classID)
	}

	db = db.Where("sc.is_active = 1").Order("s.name ASC")

	switch currenUser.ManageClass {
	case 1:
		// teacher → restrict to his assigned classes
		db = db.Where("sc.class_id IN (?)",
			config.DB.Table("class_teachers AS ct").
				Select("ct.class_id").
				Where("ct.teacher_id = ?", userID),
		)

	case 2:
		// admin → no extra restriction

	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to view students"})
		return
	}

	if err := db.Scan(&studentclass).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch scores"})
		return
	}

	c.JSON(http.StatusOK, studentclass)
}
