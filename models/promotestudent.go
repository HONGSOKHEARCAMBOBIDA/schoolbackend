package models

import "time"

type PromoteStudent struct {
	ID                 uint      `gorm:"primarykey" json:"id"`
	StudentID          int       `json:"student_id"`
	FromClassID        int       `json:"from_class_id"`
	ToClassID          int       `json:"to_class_id"`
	FromAcademicyearID int       `json:"from_academic_year_id"`
	ToAcademicyearID   int       `json:"to_academic_year_id"`
	PromoteBy          int       `json:"promote_by"`
	PromoteDate        time.Time `json:"promoted_date"`
}
type PromoteStudentRes struct {
	ID                   uint      `json:"id"`
	StudentID            int       `json:"student_id"`
	FromClassID          int       `json:"from_class_id"`
	ToClassID            int       `json:"to_class_id"`
	FromAcademicyearID   int       `json:"from_academicyear_id"`
	ToAcademicyearID     int       `json:"to_academicyear_id"`
	PromoteBy            int       `json:"promote_by"`
	PromoteDate          time.Time `json:"promote_date"`
	StudentName          string    `json:"student_name"`
	FromClassName        string    `json:"from_class_name"`
	ToClassName          string    `json:"to_class_name"`
	FromAcademicYearName string    `json:"from_academic_year_name"`
	ToAcademicYearName   string    `json:"to_academic_year_name"`
	PromoteName          string    `json:"promote_name"`
}
