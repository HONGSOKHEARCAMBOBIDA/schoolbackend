package models

import "time"

type ClassInfo struct {
	ID             int    `json:"id"`
	ClassTeacherID int    `json:"class_teacher_id"`
	Name           string `json:"name"`
	AcademicYearID int    `json:"academic_year_id"`
	AcademicYear   string `json:"academic_year"`
}

type UserDetail struct {
	ID             uint      `json:"id" gorm:"primarykey"`
	Name           string    `json:"name" gorm:"column:name"`
	Phone          string    `json:"phone"`
	Password       string    `json:"-"` // hide password
	IDCardNumber   string    `json:"id_card_number"`
	Gender         int       `json:"gender"`
	DOB            time.Time `json:"dob"`
	MaterialStatus string    `json:"marital_status"`
	Image          string    `json:"image"`

	RoleID   int    `json:"role_id"`
	RoleName string `json:"role_name"`

	VillageID   int    `json:"village_id"`
	VillageName string `json:"village_name"`

	CommuneID   int    `json:"commune_id"`
	CommuneName string `json:"commune_name"`

	DistrictID   int    `json:"district_id"`
	DistrictName string `json:"district_name"`

	ProvinceID   int    `json:"province_id"`
	ProvinceName string `json:"province_name"`

	IsActive int `json:"is_active"`

	Classes []ClassInfo `json:"classes" gorm:"-"`
}
