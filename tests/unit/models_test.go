package unit

import (
	"testing"
	"time"

	"hospital-management-system/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestPatient_CalculateAge(t *testing.T) {
	tests := []struct {
		name        string
		dateOfBirth time.Time
		expectedAge int
	}{
		{
			name:        "25 years old - birthday already passed this year",
			dateOfBirth: time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC),
			expectedAge: 25,
		},
		{
			name:        "24 years old - birthday not yet this year",
			dateOfBirth: time.Date(1999, 12, 31, 0, 0, 0, 0, time.UTC),
			expectedAge: 24,
		},
		{
			name:        "newborn - born this year",
			dateOfBirth: time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC),
			expectedAge: 0,
		},
		{
			name:        "50 years old",
			dateOfBirth: time.Date(1974, 6, 15, 0, 0, 0, 0, time.UTC),
			expectedAge: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patient := &models.Patient{
				DateOfBirth: tt.dateOfBirth,
			}

			age := patient.CalculateAge()

			currentYear := time.Now().Year()
			expectedAge := currentYear - tt.dateOfBirth.Year()

			now := time.Now()
			if now.Month() < tt.dateOfBirth.Month() ||
				(now.Month() == tt.dateOfBirth.Month() && now.Day() < tt.dateOfBirth.Day()) {
				expectedAge--
			}

			assert.Equal(t, expectedAge, age)
		})
	}
}

func TestPatient_ToResponse(t *testing.T) {
	createdBy := models.User{
		ID:        1,
		Username:  "receptionist1",
		Email:     "receptionist@hospital.com",
		FirstName: "John",
		LastName:  "Smith",
		Role:      models.RoleReceptionist,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	updatedBy := models.User{
		ID:        2,
		Username:  "doctor1",
		Email:     "doctor@hospital.com",
		FirstName: "Jane",
		LastName:  "Doe",
		Role:      models.RoleDoctor,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	patient := models.Patient{
		ID:                 1,
		PatientID:          "PAT20240101001",
		FirstName:          "Alice",
		LastName:           "Johnson",
		Email:              "alice@example.com",
		Phone:              "1234567890",
		DateOfBirth:        time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC),
		Gender:             models.GenderFemale,
		BloodType:          models.BloodTypeAPos,
		Address:            "123 Main St",
		EmergencyContact:   "0987654321",
		MedicalHistory:     "No major illnesses",
		Allergies:          "Peanuts",
		CurrentMedications: "None",
		CreatedByID:        1,
		CreatedBy:          createdBy,
		LastUpdatedByID:    &[]uint{2}[0],
		LastUpdatedBy:      &updatedBy,
		IsActive:           true,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	response := patient.ToResponse()

	assert.Equal(t, patient.ID, response.ID)
	assert.Equal(t, patient.PatientID, response.PatientID)
	assert.Equal(t, patient.FirstName, response.FirstName)
	assert.Equal(t, patient.LastName, response.LastName)
	assert.Equal(t, patient.Email, response.Email)
	assert.Equal(t, patient.Phone, response.Phone)
	assert.Equal(t, patient.DateOfBirth, response.DateOfBirth)
	assert.Equal(t, patient.Gender, response.Gender)
	assert.Equal(t, patient.BloodType, response.BloodType)
	assert.Equal(t, patient.Address, response.Address)
	assert.Equal(t, patient.EmergencyContact, response.EmergencyContact)
	assert.Equal(t, patient.MedicalHistory, response.MedicalHistory)
	assert.Equal(t, patient.Allergies, response.Allergies)
	assert.Equal(t, patient.CurrentMedications, response.CurrentMedications)
	assert.Equal(t, patient.IsActive, response.IsActive)
	assert.Equal(t, patient.CreatedAt, response.CreatedAt)
	assert.Equal(t, patient.UpdatedAt, response.UpdatedAt)

	expectedAge := patient.CalculateAge()
	assert.Equal(t, expectedAge, response.Age)

	assert.Equal(t, createdBy.ID, response.CreatedBy.ID)
	assert.Equal(t, createdBy.Username, response.CreatedBy.Username)
	assert.Equal(t, createdBy.Role, response.CreatedBy.Role)

	assert.NotNil(t, response.LastUpdatedBy)
	assert.Equal(t, updatedBy.ID, response.LastUpdatedBy.ID)
	assert.Equal(t, updatedBy.Username, response.LastUpdatedBy.Username)
	assert.Equal(t, updatedBy.Role, response.LastUpdatedBy.Role)
}

func TestPatient_ToResponse_NoLastUpdatedBy(t *testing.T) {
	createdBy := models.User{
		ID:        1,
		Username:  "receptionist1",
		Email:     "receptionist@hospital.com",
		FirstName: "John",
		LastName:  "Smith",
		Role:      models.RoleReceptionist,
		IsActive:  true,
	}

	patient := models.Patient{
		ID:              1,
		PatientID:       "PAT20240101001",
		FirstName:       "Alice",
		LastName:        "Johnson",
		CreatedByID:     1,
		CreatedBy:       createdBy,
		LastUpdatedByID: nil,
		LastUpdatedBy:   nil,
		IsActive:        true,
	}

	response := patient.ToResponse()

	assert.Equal(t, patient.ID, response.ID)
	assert.Equal(t, createdBy.ID, response.CreatedBy.ID)
	assert.Nil(t, response.LastUpdatedBy)
}

func TestUser_ToResponse(t *testing.T) {
	user := models.User{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      models.RoleDoctor,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	response := user.ToResponse()

	assert.Equal(t, user.ID, response.ID)
	assert.Equal(t, user.Username, response.Username)
	assert.Equal(t, user.Email, response.Email)
	assert.Equal(t, user.FirstName, response.FirstName)
	assert.Equal(t, user.LastName, response.LastName)
	assert.Equal(t, user.Role, response.Role)
	assert.Equal(t, user.IsActive, response.IsActive)
	assert.Equal(t, user.CreatedAt, response.CreatedAt)
	assert.Equal(t, user.UpdatedAt, response.UpdatedAt)
}

func TestPatient_TableName(t *testing.T) {
	patient := models.Patient{}
	assert.Equal(t, "patients", patient.TableName())
}

func TestUser_TableName(t *testing.T) {
	user := models.User{}
	assert.Equal(t, "users", user.TableName())
}

func TestUserRole_Constants(t *testing.T) {
	assert.Equal(t, models.UserRole("receptionist"), models.RoleReceptionist)
	assert.Equal(t, models.UserRole("doctor"), models.RoleDoctor)
}

func TestGender_Constants(t *testing.T) {
	assert.Equal(t, models.Gender("male"), models.GenderMale)
	assert.Equal(t, models.Gender("female"), models.GenderFemale)
	assert.Equal(t, models.Gender("other"), models.GenderOther)
}

func TestBloodType_Constants(t *testing.T) {
	assert.Equal(t, models.BloodType("A+"), models.BloodTypeAPos)
	assert.Equal(t, models.BloodType("A-"), models.BloodTypeANeg)
	assert.Equal(t, models.BloodType("B+"), models.BloodTypeBPos)
	assert.Equal(t, models.BloodType("B-"), models.BloodTypeBNeg)
	assert.Equal(t, models.BloodType("AB+"), models.BloodTypeABPos)
	assert.Equal(t, models.BloodType("AB-"), models.BloodTypeABNeg)
	assert.Equal(t, models.BloodType("O+"), models.BloodTypeOPos)
	assert.Equal(t, models.BloodType("O-"), models.BloodTypeONeg)
}
