package models

type AcademicYear struct {
	ID       uint   `gorm:"primarykey" json:"id"`
	YearName string `json:"year_name" gorm:"column:year_name"`
	IsActive int    `json:"is_active" gorm:"column:is_active"`
}
type AcademicYearInput struct {
	ID       uint   `json:"id"`
	YearName string `json:"year_name" gorm:"column:year_name"`
	IsActive int    `json:"is_active" gorm:"column:is_active"`
}
