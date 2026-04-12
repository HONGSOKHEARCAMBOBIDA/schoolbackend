package models

type StudentDetail struct {
	ID                   uint            `gorm:"primarykey" json:"id"`
	Code                 string          `json:"code" gorm:""`
	Name                 string          `json:"name"`
	Dob                  string          `json:"dob"`
	Gender               int             `josn:"gender"`
	Phone                string          `json:"phone"`
	IsActive             int             `json:"is_active"`
	IsPoor               int             `json:"is_poor" gorm:"column:is_poor"`
	Isdisability         int             `json:"is_disability" gorm:"column:is_disability"`
	MotherName           string          `json:"mother_name"`
	FatherName           string          `json:"father_name"`
	MotheOccupation      string          `json:"mother_occupation"`
	FotherOccupation     string          `json:"father_occupation"`
	VillageID            int             `json:"village_id"`
	VillageName          string          `json:"village_name"`
	CommuneID            int             `json:"commune_id"`
	CommuneName          string          `json:"commune_name"`
	DistrictID           int             `json:"district_id"`
	DistrictName         string          `json:"district_name"`
	ProvinceID           int             `json:"province_id"`
	ProvinceName         string          `json:"province_name"`
	TeacherName          string          `json:"teacher_name"`
	ClassName            string          `json:"class_name"`
	StudentClassID       int             `json:"student_class_id"`
	AcademicName         string          `json:"academic_name"`
	Disability           []DisabilityRes `json:"disabilities" gorm:"-"`
	StudentClassIsActive int             `json:"student_class_is_active"`
}
