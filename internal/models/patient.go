package models

import (
	"time"
	"gorm.io/gorm"
)

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

type BloodType string

const (
	BloodTypeAPos  BloodType = "A+"
	BloodTypeANeg  BloodType = "A-"
	BloodTypeBPos  BloodType = "B+"
	BloodTypeBNeg  BloodType = "B-"
	BloodTypeABPos BloodType = "AB+"
	BloodTypeABNeg BloodType = "AB-"
	BloodTypeOPos  BloodType = "O+"
	BloodTypeONeg  BloodType = "O-"
)

type Patient struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	PatientID       string         `json:"patient_id" gorm:"uniqueIndex;not null"`
	FirstName       string         `json:"first_name" gorm:"not null" binding:"required,min=2,max=50"`
	LastName        string         `json:"last_name" gorm:"not null" binding:"required,min=2,max=50"`
	Email           string         `json:"email" gorm:"uniqueIndex" binding:"omitempty,email"`
	Phone           string         `json:"phone" gorm:"not null" binding:"required,min=10,max=15"`
	DateOfBirth     time.Time      `json:"date_of_birth" gorm:"not null" binding:"required"`
	Gender          Gender         `json:"gender" gorm:"not null" binding:"required,oneof=male female other"`
	BloodType       BloodType      `json:"blood_type" binding:"omitempty,oneof=A+ A- B+ B- AB+ AB- O+ O-"`
	Address         string         `json:"address" gorm:"type:text"`
	EmergencyContact string        `json:"emergency_contact" gorm:"not null" binding:"required,min=10,max=15"`
	MedicalHistory  string         `json:"medical_history" gorm:"type:text"`
	Allergies       string         `json:"allergies" gorm:"type:text"`
	CurrentMedications string      `json:"current_medications" gorm:"type:text"`
	
	// System fields
	CreatedByID     uint           `json:"created_by_id" gorm:"not null"`
	CreatedBy       User           `json:"created_by" gorm:"foreignKey:CreatedByID"`
	LastUpdatedByID *uint          `json:"last_updated_by_id"`
	LastUpdatedBy   *User          `json:"last_updated_by,omitempty" gorm:"foreignKey:LastUpdatedByID"`
	IsActive        bool           `json:"is_active" gorm:"default:true"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Patient) TableName() string {
	return "patients"
}

type PatientResponse struct {
	ID                 uint          `json:"id"`
	PatientID          string        `json:"patient_id"`
	FirstName          string        `json:"first_name"`
	LastName           string        `json:"last_name"`
	Email              string        `json:"email"`
	Phone              string        `json:"phone"`
	DateOfBirth        time.Time     `json:"date_of_birth"`
	Age                int           `json:"age"`
	Gender             Gender        `json:"gender"`
	BloodType          BloodType     `json:"blood_type"`
	Address            string        `json:"address"`
	EmergencyContact   string        `json:"emergency_contact"`
	MedicalHistory     string        `json:"medical_history"`
	Allergies          string        `json:"allergies"`
	CurrentMedications string        `json:"current_medications"`
	CreatedBy          UserResponse  `json:"created_by"`
	LastUpdatedBy      *UserResponse `json:"last_updated_by,omitempty"`
	IsActive           bool          `json:"is_active"`
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`
}

func (p *Patient) ToResponse() PatientResponse {
	response := PatientResponse{
		ID:                 p.ID,
		PatientID:          p.PatientID,
		FirstName:          p.FirstName,
		LastName:           p.LastName,
		Email:              p.Email,
		Phone:              p.Phone,
		DateOfBirth:        p.DateOfBirth,
		Age:                p.CalculateAge(),
		Gender:             p.Gender,
		BloodType:          p.BloodType,
		Address:            p.Address,
		EmergencyContact:   p.EmergencyContact,
		MedicalHistory:     p.MedicalHistory,
		Allergies:          p.Allergies,
		CurrentMedications: p.CurrentMedications,
		CreatedBy:          p.CreatedBy.ToResponse(),
		IsActive:           p.IsActive,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
	}

	if p.LastUpdatedBy != nil {
		lastUpdatedBy := p.LastUpdatedBy.ToResponse()
		response.LastUpdatedBy = &lastUpdatedBy
	}

	return response
}

func (p *Patient) CalculateAge() int {
	now := time.Now()
	age := now.Year() - p.DateOfBirth.Year()

	if now.Month() < p.DateOfBirth.Month() || 
		(now.Month() == p.DateOfBirth.Month() && now.Day() < p.DateOfBirth.Day()) {
		age--
	}
	
	return age
} 