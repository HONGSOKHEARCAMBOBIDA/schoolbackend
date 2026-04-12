package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type ScoreDetail struct {
	ID               uint            `gorm:"primarykey" json:"id"`
	StudentClassID   int             `json:"student_class_id" gorm:"column:student_class_id"`
	ClassID          int             `json:"class_id"`
	ClassName        string          `json:"class_name"`
	StudentID        int             `json:"student_id"`
	StudentName      string          `json:"student_name"`
	AcademicYearID   int             `json:"academic_year_id"`
	AcademicYearName string          `json:"academic_year_name"`
	ComponentID      int             `json:"component_id"`
	ComponentName    string          `json:"component_name"`
	TypeExamID       int             `json:"type_exam_id"`
	TypeExamName     string          `json:"type_exam_name"`
	SubjectID        int             `json:"subject_id"`
	SubjectName      string          `json:"subject_name"`
	Mark             decimal.Decimal `json:"mark"`
	ExamDate         time.Time       `json:"exam_date"`
	CreateBy         int             `json:"create_by"`
}
