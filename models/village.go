package models

type Village struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	Name      string `json:"name"`
	CommuneID int    `json:"commune_id"`
}
