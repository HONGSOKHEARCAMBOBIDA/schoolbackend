package models

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type Scores struct {
	ID             uint            `gorm:"primarykey" json:"id"`
	StudentClassID int             `json:"student_class_id" gorm:"column:student_class_id"`
	ComponentID    int             `json:"component_id"`
	TypeExamID     int             `json:"type_exam_id"`
	Mark           decimal.Decimal `json:"mark"`
	ExamDate       time.Time       `json:"exam_date"`
	CreateBy       int             `json:"create_by"`
}

type ScoreInput struct {
	StudentClassID []int             `json:"student_class_id"`
	ComponentID    int               `json:"component_id"`
	TypeExamID     int               `json:"type_exam_id"`
	Mark           []decimal.Decimal `json:"mark"`
	ExamDate       string            `json:"exam_date"`
	CreateBy       int               `json:"create_by"`
}

func (s *ScoreInput) Validate() error {
	for i, m := range s.Mark {
		if m.LessThan(decimal.NewFromInt(0)) || m.GreaterThan(decimal.NewFromInt(10)) {
			return fmt.Errorf("ពិន្ទុមិនត្រូវតិចជាង០នឹងមិនត្រូវធំជាង10", i, m.String())
		}
	}
	return nil
}
