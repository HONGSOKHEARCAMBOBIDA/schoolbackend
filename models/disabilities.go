package models

type DisabilityRes struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
