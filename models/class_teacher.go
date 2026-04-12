package models

type ClassTeacherInput struct {
	ClassID        int `json:"class_id"`
	AcademicYearID int `json:"academic_year_id"`
	TeacherID      int `json:"teacher_id"`
	IsActive       int `json:"is_active"`
}
type ClassTeacher struct {
	ID             int `gorm:"primarykey" json:"id"`
	ClassID        int `json:"class_id"`
	AcademicYearID int `json:"academic_year_id"`
	TeacherID      int `json:"teacher_id"`
	IsActive       int `json:"is_active"`
}
