package models

type Role struct {
	ID          uint         `gorm:"primarykey" json:"id"`
	Name        string       `json:"name"`
	IsActive    int          `json:"is_active"`
	Permissions []Permission `gorm:"many2many:role_has_permissions"`
	CountUser   int          `json:"count_user"`
}
