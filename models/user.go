package models

import "time"

type User struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name_kh"`
	Phone          string    `json:"phone"`
	Password       string    `json:"password"`
	Image          string    `json:"image"`
	RoleID         int       `json:"role_id" gorm:"column:role_id"` // Capitalized for export
	VillageID      int       `json:"village_id"`
	IDCardNumber   string    `json:"id_card_number"`
	ManageClass    int       `json:"manage_class"`
	Gender         int       `json:"gender"`
	DOB            time.Time `json:"dob"`                                           // Use time.Time
	MaterialStatus int       `json:"material_status" gorm:"column:material_status"` // Changed to string (was date but unclear)
	IsActive       int       `json:"is_active"`
}
