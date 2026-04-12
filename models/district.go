package models

type District struct {
	ID         uint   `gorm:"primarykey" json:"id"`
	Name       string `json:"name"`
	ProvinceID int    `json:"province_id"`
}
