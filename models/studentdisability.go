package models

type StudentDisability struct {
	ID           uint `gorm:"primarykey" json:"id"`
	StudentID    int  `json:"student_id"`
	DisabilityID int  `json:"disability_id"`
}
