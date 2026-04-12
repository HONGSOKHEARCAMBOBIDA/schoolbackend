package models

type ClassSubjectDetail struct {
	ID          int    `gorm:"primarykey" json:"id"`
	ClassID     int    `json:"class_id" gorm:"column:class_id"`
	ClassName   string `json:"class_name"`
	SubjectID   int    `json:"subject_id" gorm:"column:subject_id"`
	SubjectName string `json:"subject_name"`
}
type ClassSubject struct {
	ID      int `gorm:"primarykey" json:"id"`
	ClassID int `json:"class_id" gorm:"column:class_id"`

	SubjectID int `json:"subject_id" gorm:"column:subject_id"`
}

type ClassSubjectInput struct {
	ClassID   int   `json:"class_id" gorm:"column:class_id"`
	SubjectID []int `json:"subject_id" gorm:"column:subject_id"`
}

type SubjectNotAssign struct {
	SubjectID   int    `json:"subject_id"`
	SubjectName string `json:"subject_name"`
}
