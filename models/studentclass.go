package models

type StudentClass struct {
	ID             int `gorm:"primarykey" json:"id"`
	StudentID      int `json:"student_id"`
	ClassID        int `json:"class_id"`
	AcademicYearID int `json:"academic_year_id"`
	PromotionID    int `json:"promote_student_id"`
	IsActive       int `json:"is_active" gorm:"column:is_active"`
}

type StudentClassInput struct {
	StudentID      int `json:"student_id"`
	ClassID        int `json:"class_id"`
	AcademicYearID int `json:"academic_year_id"`
	IsActive       int `json:"is_active"`
}
