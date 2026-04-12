package models

type StudentClassDetail struct {
	ID               int    `gorm:"primarykey" json:"id"`
	StudentID        int    `json:"student_id"`
	StudentName      string `json:"student_name"`
	ClassID          int    `json:"class_id"`
	ClassName        string `json:"class_name"`
	AcademicYearID   int    `json:"academic_year_id"`
	AcademicYearName string `json:"academic_year_name"`
	IsActive         int    `json:"is_active"`
}
