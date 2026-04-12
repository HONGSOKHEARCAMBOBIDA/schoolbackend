package models

type Province struct {
	ID   uint   `gorm:"primarykey" json:"id"`
	Name string `json:"name"`
}
