package controllers

import (
	"net/http"
	"schoolbackend/config"
	"schoolbackend/helper"
	"schoolbackend/models"
	"time"

	"github.com/gin-gonic/gin"
)

func PromoteStudent(c *gin.Context) {
	// Step 1: Check user authentication
	userID, ok := helper.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Step 2: Bind JSON input
	var input models.PromoteStudent
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Step 3: Validate promote date

	// Step 4: Create promotion record
	promotion := models.PromoteStudent{
		StudentID:          input.StudentID,
		FromClassID:        input.FromClassID,
		ToClassID:          input.ToClassID,
		FromAcademicyearID: input.FromAcademicyearID,
		ToAcademicyearID:   input.FromAcademicyearID + 1,
		PromoteBy:          userID,
		PromoteDate:        time.Now().Truncate(24 * time.Hour),
	}

	if err := config.DB.Create(&promotion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Step 5: Update old student_class (set inactive)
	var last models.StudentClass
	if err := config.DB.
		Where("student_id = ? AND is_active = 1", input.StudentID).
		Order("id desc").
		First(&last).Error; err == nil {
		last.IsActive = 0
		last.PromotionID = int(promotion.ID)
		config.DB.Save(&last)
	}

	// Step 6: Insert new student_class (active)
	newStudentClass := models.StudentClass{
		StudentID:      input.StudentID,
		ClassID:        input.ToClassID,
		AcademicYearID: input.FromAcademicyearID + 1,
		PromotionID:    int(promotion.ID),
		IsActive:       1,
	}

	if err := config.DB.Create(&newStudentClass).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Step 7: Response
	c.JSON(http.StatusOK, gin.H{
		"promotion":     promotion,
		"student_class": newStudentClass,
	})
}

func DeletePromotion(c *gin.Context) {
	// Step 1: Check user authentication
	_, ok := helper.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Step 2: Get promotion id from URL
	id := c.Param("id")

	// Step 3: Find promotion
	var promotion models.PromoteStudent
	if err := config.DB.First(&promotion, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Promotion not found"})
		return
	}

	// Step 4: Find related student_class
	var studentClasses []models.StudentClass
	if err := config.DB.Where("promotion_id = ?", promotion.ID).Find(&studentClasses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Step 5: Handle student_class records
	for _, sc := range studentClasses {
		if sc.IsActive == 1 {
			// delete active student_class
			if err := config.DB.Delete(&sc).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			// update inactive student_class to active
			sc.IsActive = 1
			if err := config.DB.Save(&sc).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	// Step 6: Delete promotion
	if err := config.DB.Delete(&promotion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Step 7: Response
	c.JSON(http.StatusOK, gin.H{"message": "Promotion deleted successfully"})
}
func GetPromote(c *gin.Context) {
	var promote []models.PromoteStudentRes
	academicYearID := c.Query("academic_year_id")
	toacademicYearID := c.Query("to_academic_year_id")
	classID := c.Query("class_id")
	toClassId := c.Query("to_class_id")
	namekh := c.Query("name")
	db := config.DB.Table("promote_students").
		Select(`
			promote_students.id,
			promote_students.student_id,
			promote_students.from_class_id,	
			promote_students.to_class_id,	
			promote_students.from_academicyear_id,
			promote_students.to_academicyear_id,	
			promote_students.promote_by,
			promote_students.promote_date,
			students.name AS student_name,
			from_classes.name AS from_class_name,
			to_classes.name AS to_class_name,
			from_academic_years.year_name AS from_academic_year_name,
			to_academic_years.year_name AS to_academic_year_name,
			users.name AS promote_name
		`).
		Joins("INNER JOIN students ON students.id = promote_students.student_id").
		Joins("INNER JOIN classes AS from_classes ON from_classes.id = promote_students.from_class_id").
		Joins("INNER JOIN classes AS to_classes ON to_classes.id = promote_students.to_class_id").
		Joins("INNER JOIN academic_years AS from_academic_years ON from_academic_years.id = promote_students.from_academicyear_id").
		Joins("INNER JOIN academic_years AS to_academic_years ON to_academic_years.id = promote_students.to_academicyear_id").
		Joins("INNER JOIN users ON users.id = promote_students.promote_by")
	if academicYearID != "" {
		db = db.Where("promote_students.from_academicyear_id =?", academicYearID)
	}
	if classID != "" {
		db = db.Where("promote_students.from_class_id =?", classID)
	}
	if namekh != "" {
		db = db.Where("students.name LIKE ?", "%"+namekh+"%")
	}
	if toacademicYearID != "" {
		db = db.Where("promote_students.to_academicyear_id =?", toacademicYearID)
	}
	if toClassId != "" {
		db = db.Where("promote_students.to_class_id =?", toClassId)
	}
	result := db.Scan(&promote)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "faild"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": promote})
}
