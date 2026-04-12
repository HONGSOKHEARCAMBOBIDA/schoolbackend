package models

type Commune struct {
	ID         uint   `gorm:"primarykey" json:"id"`
	Name       string `json:"name"`
	DistrictId int    `json:"district_id"`
}
