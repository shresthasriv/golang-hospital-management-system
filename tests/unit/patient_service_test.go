package unit

import (
	"errors"
	"testing"
	"time"

	"hospital-management-system/internal/models"
	"hospital-management-system/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPatientRepository struct {
	mock.Mock
}

func (m *MockPatientRepository) Create(patient *models.Patient) error {
	args := m.Called(patient)
	return args.Error(0)
}

func (m *MockPatientRepository) GetByID(id uint) (*models.Patient, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Patient), args.Error(1)
}

func (m *MockPatientRepository) GetByPatientID(patientID string) (*models.Patient, error) {
	args := m.Called(patientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Patient), args.Error(1)
}

func (m *MockPatientRepository) Update(patient *models.Patient) error {
	args := m.Called(patient)
	return args.Error(0)
}

func (m *MockPatientRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPatientRepository) List(limit, offset int) ([]*models.Patient, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Patient), args.Error(1)
}

func (m *MockPatientRepository) Search(query string, limit, offset int) ([]*models.Patient, error) {
	args := m.Called(query, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Patient), args.Error(1)
}

func (m *MockPatientRepository) GetByCreatedBy(userID uint, limit, offset int) ([]*models.Patient, error) {
	args := m.Called(userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Patient), args.Error(1)
}

func (m *MockPatientRepository) Count() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPatientRepository) GenerateUniquePatientID() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func TestPatientService_CreatePatient_Success(t *testing.T) {
	mockPatientRepo := new(MockPatientRepository)
	mockUserRepo := new(MockUserRepository)
	patientService := services.NewPatientService(mockPatientRepo, mockUserRepo)

	createdBy := &models.User{
		ID:       1,
		Username: "receptionist1",
		Role:     models.RoleReceptionist,
		IsActive: true,
	}

	mockUserRepo.On("GetByID", uint(1)).Return(createdBy, nil)

	mockPatientRepo.On("Create", mock.AnythingOfType("*models.Patient")).Return(nil).Run(func(args mock.Arguments) {
		patient := args.Get(0).(*models.Patient)
		patient.ID = 1
	})

	createdPatient := &models.Patient{
		ID:               1,
		PatientID:        "PAT20240101001",
		FirstName:        "John",
		LastName:         "Doe",
		Phone:            "1234567890",
		DateOfBirth:      time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Gender:           models.GenderMale,
		EmergencyContact: "0987654321",
		CreatedByID:      1,
		CreatedBy:        *createdBy,
		IsActive:         true,
	}

	mockPatientRepo.On("GetByID", uint(1)).Return(createdPatient, nil)

	req := services.CreatePatientRequest{
		FirstName:        "John",
		LastName:         "Doe",
		Phone:            "1234567890",
		DateOfBirth:      "1990-01-01",
		Gender:           models.GenderMale,
		EmergencyContact: "0987654321",
	}

	response, err := patientService.CreatePatient(req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "John", response.FirstName)
	assert.Equal(t, "Doe", response.LastName)
	assert.Equal(t, models.GenderMale, response.Gender)
	mockUserRepo.AssertExpectations(t)
	mockPatientRepo.AssertExpectations(t)
}

func TestPatientService_CreatePatient_InvalidUser(t *testing.T) {
	mockPatientRepo := new(MockPatientRepository)
	mockUserRepo := new(MockUserRepository)
	patientService := services.NewPatientService(mockPatientRepo, mockUserRepo)

	mockUserRepo.On("GetByID", uint(999)).Return(nil, errors.New("user not found"))

	req := services.CreatePatientRequest{
		FirstName:        "John",
		LastName:         "Doe",
		Phone:            "1234567890",
		DateOfBirth:      "1990-01-01",
		Gender:           models.GenderMale,
		EmergencyContact: "0987654321",
	}

	response, err := patientService.CreatePatient(req, 999)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "invalid user", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestPatientService_CreatePatient_InvalidDateFormat(t *testing.T) {
	mockPatientRepo := new(MockPatientRepository)
	mockUserRepo := new(MockUserRepository)
	patientService := services.NewPatientService(mockPatientRepo, mockUserRepo)

	createdBy := &models.User{
		ID:       1,
		Username: "receptionist1",
		Role:     models.RoleReceptionist,
		IsActive: true,
	}

	mockUserRepo.On("GetByID", uint(1)).Return(createdBy, nil)

	req := services.CreatePatientRequest{
		FirstName:        "John",
		LastName:         "Doe",
		Phone:            "1234567890",
		DateOfBirth:      "invalid-date",
		Gender:           models.GenderMale,
		EmergencyContact: "0987654321",
	}

	response, err := patientService.CreatePatient(req, 1)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid date of birth format")
	mockUserRepo.AssertExpectations(t)
}

func TestPatientService_CreatePatient_FutureDateOfBirth(t *testing.T) {
	mockPatientRepo := new(MockPatientRepository)
	mockUserRepo := new(MockUserRepository)
	patientService := services.NewPatientService(mockPatientRepo, mockUserRepo)

	createdBy := &models.User{
		ID:       1,
		Username: "receptionist1",
		Role:     models.RoleReceptionist,
		IsActive: true,
	}

	mockUserRepo.On("GetByID", uint(1)).Return(createdBy, nil)

	futureDate := time.Now().AddDate(1, 0, 0).Format("2006-01-02")
	req := services.CreatePatientRequest{
		FirstName:        "John",
		LastName:         "Doe",
		Phone:            "1234567890",
		DateOfBirth:      futureDate,
		Gender:           models.GenderMale,
		EmergencyContact: "0987654321",
	}

	response, err := patientService.CreatePatient(req, 1)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "date of birth cannot be in the future", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestPatientService_GetPatientByID_Success(t *testing.T) {
	mockPatientRepo := new(MockPatientRepository)
	mockUserRepo := new(MockUserRepository)
	patientService := services.NewPatientService(mockPatientRepo, mockUserRepo)

	patient := &models.Patient{
		ID:        1,
		PatientID: "PAT20240101001",
		FirstName: "John",
		LastName:  "Doe",
	}

	mockPatientRepo.On("GetByID", uint(1)).Return(patient, nil)

	response, err := patientService.GetPatientByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, "John", response.FirstName)
	mockPatientRepo.AssertExpectations(t)
}

func TestPatientService_GetPatientByID_NotFound(t *testing.T) {
	mockPatientRepo := new(MockPatientRepository)
	mockUserRepo := new(MockUserRepository)
	patientService := services.NewPatientService(mockPatientRepo, mockUserRepo)

	mockPatientRepo.On("GetByID", uint(999)).Return(nil, errors.New("patient not found"))

	response, err := patientService.GetPatientByID(999)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "patient not found", err.Error())
	mockPatientRepo.AssertExpectations(t)
}

func TestPatientService_GetPatientByPatientID_Success(t *testing.T) {
	mockPatientRepo := new(MockPatientRepository)
	mockUserRepo := new(MockUserRepository)
	patientService := services.NewPatientService(mockPatientRepo, mockUserRepo)

	patient := &models.Patient{
		ID:        1,
		PatientID: "PAT20240101001",
		FirstName: "John",
		LastName:  "Doe",
	}

	mockPatientRepo.On("GetByPatientID", "PAT20240101001").Return(patient, nil)

	response, err := patientService.GetPatientByPatientID("PAT20240101001")

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "PAT20240101001", response.PatientID)
	assert.Equal(t, "John", response.FirstName)
	mockPatientRepo.AssertExpectations(t)
}

func TestPatientService_UpdatePatient_Success_Receptionist(t *testing.T) {
	mockPatientRepo := new(MockPatientRepository)
	mockUserRepo := new(MockUserRepository)
	patientService := services.NewPatientService(mockPatientRepo, mockUserRepo)

	existingPatient := &models.Patient{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "1234567890",
	}

	updatedBy := &models.User{
		ID:       2,
		Username: "receptionist1",
		Role:     models.RoleReceptionist,
		IsActive: true,
	}

	mockPatientRepo.On("GetByID", uint(1)).Return(existingPatient, nil)
	mockUserRepo.On("GetByID", uint(2)).Return(updatedBy, nil)
	mockPatientRepo.On("Update", mock.AnythingOfType("*models.Patient")).Return(nil)
	mockPatientRepo.On("GetByID", uint(1)).Return(existingPatient, nil) // For returning updated patient

	newFirstName := "Jane"
	req := services.UpdatePatientRequest{
		FirstName: &newFirstName,
	}

	response, err := patientService.UpdatePatient(1, req, 2, models.RoleReceptionist)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	mockPatientRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestPatientService_UpdatePatient_Success_Doctor_MedicalFields(t *testing.T) {
	mockPatientRepo := new(MockPatientRepository)
	mockUserRepo := new(MockUserRepository)
	patientService := services.NewPatientService(mockPatientRepo, mockUserRepo)

	existingPatient := &models.Patient{
		ID:             1,
		FirstName:      "John",
		LastName:       "Doe",
		MedicalHistory: "None",
	}

	updatedBy := &models.User{
		ID:       2,
		Username: "doctor1",
		Role:     models.RoleDoctor,
		IsActive: true,
	}

	mockPatientRepo.On("GetByID", uint(1)).Return(existingPatient, nil)
	mockUserRepo.On("GetByID", uint(2)).Return(updatedBy, nil)
	mockPatientRepo.On("Update", mock.AnythingOfType("*models.Patient")).Return(nil)
	mockPatientRepo.On("GetByID", uint(1)).Return(existingPatient, nil)	

	newMedicalHistory := "Diabetes"
	req := services.UpdatePatientRequest{
		MedicalHistory: &newMedicalHistory,
	}

	response, err := patientService.UpdatePatient(1, req, 2, models.RoleDoctor)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	mockPatientRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestPatientService_DeletePatient_Success_Receptionist(t *testing.T) {
	mockPatientRepo := new(MockPatientRepository)
	mockUserRepo := new(MockUserRepository)
	patientService := services.NewPatientService(mockPatientRepo, mockUserRepo)

	patient := &models.Patient{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
	}

	mockPatientRepo.On("GetByID", uint(1)).Return(patient, nil)
	mockPatientRepo.On("Delete", uint(1)).Return(nil)

	err := patientService.DeletePatient(1, models.RoleReceptionist)

	assert.NoError(t, err)
	mockPatientRepo.AssertExpectations(t)
}

func TestPatientService_DeletePatient_Forbidden_Doctor(t *testing.T) {
	mockPatientRepo := new(MockPatientRepository)
	mockUserRepo := new(MockUserRepository)
	patientService := services.NewPatientService(mockPatientRepo, mockUserRepo)

	err := patientService.DeletePatient(1, models.RoleDoctor)

	assert.Error(t, err)
	assert.Equal(t, "only receptionists can delete patients", err.Error())
}

func TestPatientService_ListPatients_Success(t *testing.T) {
	mockPatientRepo := new(MockPatientRepository)
	mockUserRepo := new(MockUserRepository)
	patientService := services.NewPatientService(mockPatientRepo, mockUserRepo)

	patients := []*models.Patient{
		{ID: 1, FirstName: "John", LastName: "Doe"},
		{ID: 2, FirstName: "Jane", LastName: "Smith"},
	}

	mockPatientRepo.On("List", 10, 0).Return(patients, nil)
	mockPatientRepo.On("Count").Return(int64(2), nil)

	response, err := patientService.ListPatients(1, 10)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response.Patients, 2)
	assert.Equal(t, int64(2), response.Pagination.Total)
	assert.Equal(t, 1, response.Pagination.CurrentPage)
	assert.Equal(t, 10, response.Pagination.PageSize)
	mockPatientRepo.AssertExpectations(t)
}

func TestPatientService_SearchPatients_Success(t *testing.T) {
	mockPatientRepo := new(MockPatientRepository)
	mockUserRepo := new(MockUserRepository)
	patientService := services.NewPatientService(mockPatientRepo, mockUserRepo)

	patients := []*models.Patient{
		{ID: 1, FirstName: "John", LastName: "Doe"},
	}

	mockPatientRepo.On("Search", "John", 10, 0).Return(patients, nil)

	response, err := patientService.SearchPatients("John", 1, 10)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response.Patients, 1)
	assert.Equal(t, "John", response.Patients[0].FirstName)
	mockPatientRepo.AssertExpectations(t)
}
