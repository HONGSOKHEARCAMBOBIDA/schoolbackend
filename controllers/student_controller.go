package controllers

import (
	"errors"
	"net/http"
	"schoolbackend/config"
	"schoolbackend/helper"
	"schoolbackend/models"
	"schoolbackend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// func SaveStudnent(c *gin.Context) {
// 	var input models.StudentInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	tx := config.DB.Begin()
// 	if tx.Error != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start transaction"})
// 		return
// 	}
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 		}
// 	}()
// 	var student models.Student
// 	if input.ID == 0 {
// 		student = models.Student{
// 			Name:             input.Name,
// 			Code:             utils.GenerateStudentCode(),
// 			Dob:              input.Dob,
// 			Gender:           input.Gender,
// 			Phone:            input.Phone,
// 			VillageID:        input.VillageID,
// 			IsActive:         1,
// 			IsPoor:           input.IsPoor,
// 			Isdisability:     input.Isdisability,
// 			MotherName:       input.MotherName,
// 			FatherName:       input.FatherName,
// 			MotheOccupation:  input.MotheOccupation,
// 			FotherOccupation: input.FotherOccupation,
// 		}
// 		if err := tx.Create(&student).Error; err != nil {
// 			tx.Rollback()
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{"message": "created"})
// 	} else {
// 		if err := tx.First(&student, input.ID).Error; err != nil {
// 			tx.Rollback()
// 			c.JSON(http.StatusFound, gin.H{"error": err.Error()})
// 			return
// 		}
// 		student.Name = input.Name
// 		student.Dob = input.Dob
// 		student.Gender = input.Gender
// 		student.Phone = input.Phone
// 		student.VillageID = input.VillageID
// 		student.IsActive = 1
// 		student.MotherName = input.MotherName
// 		student.FatherName = input.FatherName
// 		student.MotheOccupation = input.MotheOccupation
// 		student.FotherOccupation = input.FotherOccupation
// 		if err := tx.Save(&student).Error; err != nil {
// 			tx.Rollback()
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{"message": "updated"})

//		}
//		tx.Commit()
//	}
func SaveStudent(c *gin.Context) {
	var input models.StudentInput
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

	var student models.Student
	if input.ID == 0 {
		// Create new student
		student = models.Student{
			Name:             input.Name,
			Code:             utils.GenerateStudentCode(),
			Dob:              input.Dob,
			Gender:           input.Gender,
			Phone:            input.Phone,
			VillageID:        input.VillageID,
			IsActive:         1,
			IsPoor:           input.IsPoor,
			Isdisability:     input.Isdisability,
			MotherName:       input.MotherName,
			FatherName:       input.FatherName,
			MotheOccupation:  input.MotheOccupation,
			FotherOccupation: input.FotherOccupation,
		}
		if err := tx.Create(&student).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Handle student_disabilities if needed
		if input.Isdisability == 1 && len(input.DisabilityIDs) > 0 {
			for _, did := range input.DisabilityIDs {
				if err := tx.Create(&models.StudentDisability{
					StudentID:    int(student.ID),
					DisabilityID: did,
				}).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
		}
		if input.ClassId != 0 {
			studentclass := models.StudentClass{
				StudentID:      int(student.ID),
				ClassID:        input.ClassId,
				IsActive:       1,
				AcademicYearID: input.AcademicYearID,
			}
			if err := tx.Create(&studentclass).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "created"})
	} else {
		// Update existing student
		if err := tx.First(&student, input.ID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		student.Name = input.Name
		student.Dob = input.Dob
		student.Gender = input.Gender
		student.Phone = input.Phone
		student.VillageID = input.VillageID
		student.IsActive = 1
		student.MotherName = input.MotherName
		student.FatherName = input.FatherName
		student.MotheOccupation = input.MotheOccupation
		student.FotherOccupation = input.FotherOccupation
		student.Isdisability = input.Isdisability

		if err := tx.Save(&student).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Clear old disabilities and insert new ones
		if input.Isdisability == 1 {
			if err := tx.Where("student_id = ?", student.ID).Delete(&models.StudentDisability{}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			for _, did := range input.DisabilityIDs {
				if err := tx.Create(&models.StudentDisability{
					StudentID:    int(student.ID),
					DisabilityID: did,
				}).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
		} else {
			// if student is not disabled, remove any existing entries
			tx.Where("student_id = ?", student.ID).Delete(&models.StudentDisability{})
		}

		c.JSON(http.StatusOK, gin.H{"message": "updated"})
	}

	tx.Commit()
}

func HandlStudent(c *gin.Context) {
	method := c.Request.Method

	if method == http.MethodGet {
		var students []models.StudentDetail

		// ---- get current user id ----
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		userIDFloat, ok := userIDInterface.(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID type"})
			return
		}
		userID := int(userIDFloat)

		// ---- check if user manages classes ----
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

		academicYearID := c.Query("academic_year_id")
		classID := c.Query("class_id")
		poorID := c.Query("is_poor")
		disabilityID := c.Query("is_disability")
		namekh := c.Query("name")
		suspendStudy := c.Query("SuspendStudy")
		chnagescool := c.Query("changeschool")
		stopstudy := c.Query("stopstudy")

		// ---- declare db ----
		db := config.DB.Table("students").Select(`
			students.id,
			students.code,
			students.name,
			students.dob,
			students.gender,
			students.phone,
			students.is_active,
			students.is_poor,
			students.is_disability,
			students.mother_name,
			students.father_name,
			students.mothe_occupation,
			students.fother_occupation,
			sc.id AS student_class_id,
			sc.class_id AS class_id,
			sc.academic_year_id AS academic_year_id,
			sc.is_active AS student_class_is_active,
			a.year_name AS academic_name,
			students.village_id, villages.name AS village_name,
			communes.id AS commune_id, communes.name AS commune_name,
			districts.id AS district_id, districts.name AS district_name,
			provinces.id AS province_id, provinces.name AS province_name,
			u.name AS teacher_name,
			c.name AS class_name
		`).
			Joins("LEFT JOIN villages ON villages.id = students.village_id").
			Joins("LEFT JOIN communes ON communes.id = villages.commune_id").
			Joins("LEFT JOIN districts ON districts.id = communes.district_id").
			Joins("LEFT JOIN provinces ON provinces.id = districts.province_id").
			Joins("INNER JOIN student_classes sc ON sc.student_id = students.id").
			Joins("INNER JOIN academic_years a ON a.id = sc.academic_year_id ").
			Joins("INNER JOIN classes c ON c.id = sc.class_id").
			Joins("INNER JOIN class_teachers ct ON ct.class_id = sc.class_id AND ct.academic_year_id = sc.academic_year_id AND ct.is_active =1").
			Joins("LEFT JOIN users u ON u.id = ct.teacher_id")

		// ---- filter by academic_year_id if given ----
		if academicYearID != "" {
			db = db.Where("sc.academic_year_id = ?", academicYearID)
		}

		// ---- filter by class_id if given ----
		if classID != "" {
			db = db.Where("sc.class_id = ?", classID)
		}
		if poorID != "" {
			db = db.Where("students.is_poor =?", poorID)
		}
		if disabilityID != "" {
			db = db.Where("students.is_disability", disabilityID)
		}
		if suspendStudy != "" {
			db = db.Where("sc.is_active =?", suspendStudy)
		}
		if chnagescool != "" {
			db = db.Where("sc.is_active =?", chnagescool)
		}
		if stopstudy != "" {
			db = db.Where("sc.is_active =?", stopstudy)
		}

		switch currenUser.ManageClass {
		case 1:
			// teacher → restrict to his assigned classes
			db = db.Where("sc.class_id IN (?)",
				config.DB.Table("class_teachers as ct").
					Select("ct.class_id").
					Where("ct.teacher_id =?", userID),
			)

		case 2:
			// admin → no extra restriction

		default:
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to view students"})
			return
		}

		// ---- optional filter by name ----
		if namekh != "" {
			db = db.Where("students.name LIKE ?", "%"+namekh+"%")
		}

		// ---- run query ----
		result := db.Scan(&students)
		for i, student := range students {
			var disabilities []models.DisabilityRes
			config.DB.Raw(`
		SELECT d.id, d.name
		FROM student_disabilities sd
		JOIN disability_res d ON d.id = sd.disability_id
		WHERE sd.student_id = ?
	`, student.ID).Scan(&disabilities)

			students[i].Disability = disabilities
		}
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch students"})
			return
		}

		for i := range students {
			students[i].Dob = helper.FormatDate(students[i].Dob)
		}

		c.JSON(http.StatusOK, gin.H{"students": students})

	} else if method == http.MethodPut {
		// ----- TOGGLE STATUS -----
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var last models.StudentClass
		if err := config.DB.Where("id = ?", id).First(&last).Error; err == nil {
			last.IsActive = 4
			config.DB.Save(&last)
		}

		c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})

	} else {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Only GET and PUT are allowed"})
	}
}
func SuspendStudies(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

	var last models.StudentClass
	err = config.DB.
		Where("id = ?", id).
		First(&last).Error

	if err != nil {
		// Not found or DB error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no active class found for this student"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Found -> update
	last.IsActive = 2 // 2 = Suspended
	if err := config.DB.Save(&last).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "student suspended successfully"})
}
func ChangeSchool(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

	var last models.StudentClass
	err = config.DB.
		Where("id = ?", id).
		First(&last).Error

	if err != nil {
		// Not found or DB error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no active class found for this student"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Found -> update
	last.IsActive = 3 // 2 = Suspended
	if err := config.DB.Save(&last).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "student suspended successfully"})
}
func Comeback(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

	var last models.StudentClass
	err = config.DB.
		Where("id = ?", id).
		First(&last).Error

	if err != nil {
		// Not found or DB error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no active class found for this student"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Found -> update
	last.IsActive = 1 // 2 = Suspended
	if err := config.DB.Save(&last).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "student suspended successfully"})
}
func Getstudent(c *gin.Context) {
	method := c.Request.Method

	if method == http.MethodGet {
		var students []models.StudentDetail

		// ---- get current user id ----
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		userIDFloat, ok := userIDInterface.(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID type"})
			return
		}
		userID := int(userIDFloat)

		// ---- check if user manages classes ----
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

		// ---- optional query params ----
		namekh := c.Query("name")

		// ---- declare db ----
		db := config.DB.Table("students").Select(`
			students.id,
			students.code,
			students.name,
			students.dob,
			students.gender,
			students.phone,
			students.is_active,
			students.mother_name,
			students.father_name,
			students.mothe_occupation,
			students.fother_occupation,
			students.village_id, villages.name AS village_name,
			communes.id AS commune_id, communes.name AS commune_name,
			districts.id AS district_id, districts.name AS district_name,
			provinces.id AS province_id, provinces.name AS province_name
		`).
			Joins("LEFT JOIN villages ON villages.id = students.village_id").
			Joins("LEFT JOIN communes ON communes.id = villages.commune_id").
			Joins("LEFT JOIN districts ON districts.id = communes.district_id").
			Joins("LEFT JOIN provinces ON provinces.id = districts.province_id").
			Joins("LEFT JOIN student_classes sc ON sc.student_id = students.id").
			Where("sc.id IS NULL OR sc.is_active IN (2,3,4)")

		// ---- optional filter by name ----
		if namekh != "" {
			db = db.Where("students.name LIKE ?", "%"+namekh+"%")
		}

		// ---- run query ----
		if err := db.Scan(&students).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch students"})
			return
		}

		// ---- response ----
		c.JSON(http.StatusOK, gin.H{"students": students})
	}
}
