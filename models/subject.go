package models

type Subject struct {
	ID       int    `gorm:"primarykey" json:"id"`
	Name     string `json:"name"`
	IsActive int    `json:"is_active"`
}

type SubjectInput struct {
	ID   int    `gorm:"primarykey" json:"id"`
	Name string `json:"name"`
}
