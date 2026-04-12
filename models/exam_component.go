package models

type ExamComponent struct {
	ID             uint   `gorm:"primarykey" json:"id"`
	Name           string `json:"name"`
	ClassSubjectId int    `json:"class_subject_id"`
	IsActive       int    `json:"is_active"`
}

type ExamComponentInput struct {
	Name           string `json:"name"`
	ClassSubjectId int    `json:"class_subject_id"`
	IsActive       int    `json:"is_active"`
}
