package models

type TypeExam struct {
	Id   uint   `gorm:"primarykey" json:"id"`
	Name string `json:"name"`
}
