package models

type TeacherSubjectDetail struct {
	ID               int    `gorm:"primarykey" json:"id"`
	UserID           int    `json:"user_id"`
	UserName         string `json:"user_name"`
	AcademicYearID   int    `json:"academicyear_id"`
	AcademicYearName string `json:"academic_year_name"`
	ClassID          int    `json:"class_id"`
	ClassName        string `json:"class_name"`
	SubjectID        int    `json:"subject_id"`
	SubjectName      string `json:"subject_name"`
	IsActive         int    `json:"is_active" gorm:"column:is_active"`
}
