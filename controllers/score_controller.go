package controllers

import (
	"math"
	"net/http"
	"schoolbackend/config"
	"schoolbackend/helper"
	"schoolbackend/models"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func UpdateScore(c *gin.Context) {

	// Get score ID from URL
	scoreID := c.Param("id")
	var score models.Scores

	if err := config.DB.First(&score, scoreID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Score not found"})
		return
	}

	var input struct {
		Mark     decimal.Decimal `json:"mark"`
		ExamDate string          `json:"exam_date"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse exam date
	examDate, err := time.Parse("2006-01-02", input.ExamDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exam date"})
		return
	}

	// Update fields
	score.Mark = input.Mark
	score.ExamDate = examDate

	if err := config.DB.Save(&score).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, score)
}

func DeleteScore(c *gin.Context) {
	// Get score ID from URL
	scoreID := c.Param("id")
	var score models.Scores

	if err := config.DB.First(&score, scoreID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Score not found"})
		return
	}

	if err := config.DB.Delete(&score).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Score deleted successfully"})
}

func CreateScore(c *gin.Context) {
	userID, ok := helper.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input models.ScoreInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ពិន្ទុមិនត្រូវតិចជាង០នឹងមិនត្រូវធំជាង10"})
		return
	}

	// parse exam date
	examDate, err := time.Parse("2006-01-02", input.ExamDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exam date"})
		return
	}

	var scores []models.Scores
	for i, studentClassID := range input.StudentClassID {
		mark := input.Mark[i]

		scores = append(scores, models.Scores{
			StudentClassID: studentClassID,
			ComponentID:    input.ComponentID,
			TypeExamID:     input.TypeExamID,
			Mark:           mark,
			ExamDate:       examDate,
			CreateBy:       userID,
		})
	}

	if len(scores) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid scores to insert"})
		return
	}

	if err := config.DB.Create(&scores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, scores)
}

func GetScore(c *gin.Context) {
	var scores []models.ScoreDetail
	academicYearID := c.Query("academic_year_id")
	classID := c.Query("class_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	componentID := c.Query("component_id")
	typeExamID := c.Query("type_exam_id")

	// Pagination params
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Base query
	db := config.DB.Table("scores").Select(`
        scores.id,
        scores.student_class_id,
        scores.component_id,
        scores.type_exam_id,
        scores.mark,
        scores.exam_date,
        students.id as student_id,
        students.name as student_name,
        classes.id as class_id,
        classes.name as class_name,
        academic_years.id as academic_year_id,
        academic_years.year_name as academic_year_name,
        exam_components.name as component_name,
        type_exams.name as type_exam_name,
        subjects.id as subject_id,
        subjects.name as subject_name
    `).
		Joins("INNER JOIN student_classes ON student_classes.id = scores.student_class_id").
		Joins("INNER JOIN students ON students.id = student_classes.student_id").
		Joins("INNER JOIN classes ON classes.id = student_classes.class_id").
		Joins("INNER JOIN academic_years ON academic_years.id = student_classes.academic_year_id").
		Joins("INNER JOIN exam_components ON exam_components.id = scores.component_id").
		Joins("INNER JOIN class_subjects ON class_subjects.id = exam_components.class_subject_id").
		Joins("INNER JOIN subjects ON subjects.id = class_subjects.subject_id").
		Joins("INNER JOIN type_exams ON type_exams.id = scores.type_exam_id")

	// Filters
	if academicYearID != "" {
		db = db.Where("student_classes.academic_year_id = ?", academicYearID)
	}
	if classID != "" {
		db = db.Where("student_classes.class_id = ?", classID)
	}
	if startDate != "" && endDate != "" {
		db = db.Where("scores.exam_date BETWEEN ? AND ?", startDate, endDate)
	}
	if componentID != "" {
		db = db.Where("scores.component_id = ?", componentID)
	}
	if typeExamID != "" {
		db = db.Where("scores.type_exam_id = ?", typeExamID)
	}

	// Count total before applying limit/offset
	var total int64
	if err := db.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count scores"})
		return
	}

	// Apply pagination + order
	result := db.Order("students.name ASC").Limit(limit).Offset(offset).Scan(&scores)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch scores"})
		return
	}

	// Return paginated result
	c.JSON(http.StatusOK, gin.H{
		"data":  scores,
		"page":  page,
		"limit": limit,
		"total": total,
		"pages": (total + int64(limit) - 1) / int64(limit), // ceil division
	})
}

type StudentAverage struct {
	StudentID     int     `json:"student_id"`
	StudentName   string  `json:"student_name"`
	StudentGender int     `json:"student_gender"`
	ClassName     string  `json:"class_name"`
	AvgMark       float64 `json:"avg_mark"`
	Rank          int     `json:"rank"`
	Grade         string  `json:"grade"`
	Type          string  `json:"type"`
}

func GetAverageScore(c *gin.Context) {
	var results []StudentAverage

	academicYearID := c.Query("academic_year_id")
	classID := c.Query("class_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	type_exam_id := c.Query("type_exam_id")

	db := config.DB.Table("scores").
		Select("students.id AS student_id, students.name AS student_name,classes.name AS class_name,students.gender AS student_gender, AVG(CAST(scores.mark AS DECIMAL(10,2))) AS avg_mark").
		Joins("INNER JOIN student_classes ON student_classes.id = scores.student_class_id").
		Joins("INNER JOIN students ON students.id = student_classes.student_id").
		Joins("INNER JOIN classes ON classes.id = student_classes.class_id")

	if academicYearID != "" {
		db = db.Where("student_classes.academic_year_id = ?", academicYearID)
	}
	if classID != "" {
		db = db.Where("student_classes.class_id = ?", classID)
	}
	if startDate != "" && endDate != "" {
		db = db.Where("scores.exam_date BETWEEN ? AND ?", startDate, endDate)
	}
	if type_exam_id != "" {
		db = db.Where("scores.type_exam_id =?", type_exam_id)
	}

	db = db.Group("students.id, students.name, classes.name, students.gender")

	if err := db.Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// assign rank
	sort.Slice(results, func(i, j int) bool {
		return results[i].AvgMark > results[j].AvgMark
	})
	rank := 1
	for i := range results {
		if i > 0 && results[i].AvgMark < results[i-1].AvgMark {
			rank = i + 1
		}
		results[i].Rank = rank

		// grade logic
		if results[i].AvgMark < 5 {
			results[i].Grade = "F"
		} else {
			switch results[i].Rank {
			case 1:
				results[i].Grade = "A"
			case 2:
				results[i].Grade = "B"
			case 3:
				results[i].Grade = "C"
			case 4:
				results[i].Grade = "D"
			default:
				results[i].Grade = "E"
			}
		}
	}

	c.JSON(http.StatusOK, results)
}

func GetAverageScoreBigthan5(c *gin.Context) {
	var results []StudentAverage

	academicYearID := c.Query("academic_year_id")
	classID := c.Query("class_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	type_exam_id := c.Query("type_exam_id")
	student_name := c.Query("student_name")

	db := config.DB.Table("scores").
		Select("students.id AS student_id, students.name AS student_name, AVG(CAST(scores.mark AS DECIMAL(10,2))) AS avg_mark").
		Joins("INNER JOIN student_classes ON student_classes.id = scores.student_class_id").Where("student_classes.is_active =?", 1).
		Joins("INNER JOIN students ON students.id = student_classes.student_id")

	if academicYearID != "" {
		db = db.Where("student_classes.academic_year_id = ?", academicYearID)
	}
	if classID != "" {
		db = db.Where("student_classes.class_id = ?", classID)
	}
	if startDate != "" && endDate != "" {
		db = db.Where("scores.exam_date BETWEEN ? AND ?", startDate, endDate)
	}
	if type_exam_id != "" {
		db = db.Where("scores.type_exam_id = ?", type_exam_id)
	}
	if student_name != "" {
		db = db.Where("students.name LIKE ?", "%"+student_name+"%")
	}

	// group before having
	db = db.Group("students.id, students.name,students.gender")

	if err := db.Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// assign rank
	sort.Slice(results, func(i, j int) bool {
		return results[i].AvgMark > results[j].AvgMark
	})
	rank := 1
	for i := range results {
		if i > 0 && results[i].AvgMark < results[i-1].AvgMark {
			rank = i + 1
		}
		results[i].Rank = rank
		if results[i].AvgMark >= 5 {
			results[i].Type = "ជាប់"
		} else {
			results[i].Type = "ធ្លាក់"
		}

	}

	c.JSON(http.StatusOK, results)
}
func GetAnnualAverageScore(c *gin.Context) {
	var results []StudentAverage

	academicYearID := c.Query("academic_year_id")
	classID := c.Query("class_id")

	db := config.DB.Table("scores").
		Select(`
			students.id AS student_id, 
			students.name AS student_name,
			AVG(CASE WHEN scores.type_exam_id = 2 THEN CAST(scores.mark AS DECIMAL(10,2)) ELSE NULL END) AS avg_sem1,
			AVG(CASE WHEN scores.type_exam_id = 3 THEN CAST(scores.mark AS DECIMAL(10,2)) ELSE NULL END) AS avg_sem2


			
		`).
		Joins("INNER JOIN student_classes ON student_classes.id = scores.student_class_id").
		Joins("INNER JOIN students ON students.id = student_classes.student_id")

	if academicYearID != "" {
		db = db.Where("student_classes.academic_year_id = ?", academicYearID)
	}
	if classID != "" {
		db = db.Where("student_classes.class_id = ?", classID)
	}

	db = db.Group("students.id, students.name")
	type Temp struct {
		StudentID   uint    `json:"student_id"`
		StudentName string  `json:"student_name"`
		AvgSem1     float64 `json:"avg_sem1"`
		AvgSem2     float64 `json:"avg_sem2"`
	}

	var temps []Temp
	if err := db.Scan(&temps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// បម្លែងទៅជា Annual Average
	for _, t := range temps {
		annual := (t.AvgSem1 + (t.AvgSem2 * 2)) / 3
		annual = math.Round(annual*100) / 100
		results = append(results, StudentAverage{
			StudentID:   int(t.StudentID),
			StudentName: t.StudentName,
			AvgMark:     annual,
		})
	}

	// Rank
	sort.Slice(results, func(i, j int) bool {
		return results[i].AvgMark > results[j].AvgMark
	})
	rank := 1
	for i := range results {
		if i > 0 && results[i].AvgMark < results[i-1].AvgMark {
			rank = i + 1
		}
		results[i].Rank = rank
	}

	c.JSON(http.StatusOK, results)
}
