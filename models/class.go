package models

type Class struct {
	ID       uint   `gorm:"primarykey" json:"id"`
	Name     string `json:"name"`
	IsActive int    `json:"is_active"`
}

type ClassInput struct {
	Name string `json:"name" binding:"required"`
	ID   int    `json:"id"`
}
