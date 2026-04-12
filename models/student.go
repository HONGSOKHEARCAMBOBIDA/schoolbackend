package models

type Student struct {
	ID               uint   `gorm:"primarykey" json:"id"`
	Code             string `json:"code"`
	Name             string `json:"name"`
	Dob              string `json:"dob"`
	Gender           int    `josn:"gender"`
	Phone            string `json:"phone"`
	VillageID        int    `json:"village_id"`
	IsPoor           int    `json:"is_poor"`
	Isdisability     int    `json:"is_disability" gorm:"column:is_disability"`
	IsActive         int    `json:"is_active"`
	MotherName       string `json:"mother_name"`
	FatherName       string `json:"father_name"`
	MotheOccupation  string `json:"mother_occupation"`
	FotherOccupation string `json:"father_occupation"`
}

type StudentInput struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	Dob              string `json:"dob"`
	Gender           int    `json:"gender"`
	Phone            string `json:"phone"`
	VillageID        int    `json:"village_id"`
	IsActive         int    `json:"is_active"`
	IsPoor           int    `json:"is_poor"`
	Isdisability     int    `json:"is_disability"`
	MotherName       string `json:"mother_name"`
	FatherName       string `json:"father_name"`
	MotheOccupation  string `json:"mother_occupation"`
	FotherOccupation string `json:"father_occupation"`
	DisabilityIDs    []int  `json:"disability_id"`
	ClassId          int    `json:"class_id"`
	AcademicYearID   int    `json:"academic_year_id"`
}
