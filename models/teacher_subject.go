package models

type TeacherSubjectInput struct {
	UserID         int   `json:"user_id"`
	ClassSubjectID []int `json:"class_subject_id"`
	AcademicYearID int   `json:"academic_year_id"`
}
type TeacherSubject struct {
	ID             uint `gorm:"primarykey" json:"id"`
	UserID         int  `json:"user_id"`
	ClassSubjectID int  `json:"class_subject_id"`
	AcademicYearID int  `json:"academic_year_id"`
	IsActive       int  `json:"is_active"`
}
