package services

import (
	"errors"
	"time"

	"hospital-management-system/internal/models"
	"hospital-management-system/internal/repository"
)

type CreatePatientRequest struct {
	FirstName          string           `json:"first_name" binding:"required,min=2,max=50"`
	LastName           string           `json:"last_name" binding:"required,min=2,max=50"`
	Email              string           `json:"email" binding:"omitempty,email"`
	Phone              string           `json:"phone" binding:"required,min=10,max=15"`
	DateOfBirth        string           `json:"date_of_birth" binding:"required"` // Format: YYYY-MM-DD
	Gender             models.Gender    `json:"gender" binding:"required,oneof=male female other"`
	BloodType          models.BloodType `json:"blood_type" binding:"omitempty,oneof=A+ A- B+ B- AB+ AB- O+ O-"`
	Address            string           `json:"address"`
	EmergencyContact   string           `json:"emergency_contact" binding:"required,min=10,max=15"`
	MedicalHistory     string           `json:"medical_history"`
	Allergies          string           `json:"allergies"`
	CurrentMedications string           `json:"current_medications"`
}

type UpdatePatientRequest struct {
	FirstName          *string           `json:"first_name,omitempty" binding:"omitempty,min=2,max=50"`
	LastName           *string           `json:"last_name,omitempty" binding:"omitempty,min=2,max=50"`
	Email              *string           `json:"email,omitempty" binding:"omitempty,email"`
	Phone              *string           `json:"phone,omitempty" binding:"omitempty,min=10,max=15"`
	DateOfBirth        *string           `json:"date_of_birth,omitempty"`
	Gender             *models.Gender    `json:"gender,omitempty" binding:"omitempty,oneof=male female other"`
	BloodType          *models.BloodType `json:"blood_type,omitempty" binding:"omitempty,oneof=A+ A- B+ B- AB+ AB- O+ O-"`
	Address            *string           `json:"address,omitempty"`
	EmergencyContact   *string           `json:"emergency_contact,omitempty" binding:"omitempty,min=10,max=15"`
	MedicalHistory     *string           `json:"medical_history,omitempty"`
	Allergies          *string           `json:"allergies,omitempty"`
	CurrentMedications *string           `json:"current_medications,omitempty"`
}

type PatientListResponse struct {
	Patients   []models.PatientResponse `json:"patients"`
	Pagination PaginationResponse       `json:"pagination"`
}

type PaginationResponse struct {
	Total       int64 `json:"total"`
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalPages  int   `json:"total_pages"`
}

type PatientService struct {
	patientRepo repository.PatientRepository
	userRepo    repository.UserRepository
}

func NewPatientService(patientRepo repository.PatientRepository, userRepo repository.UserRepository) *PatientService {
	return &PatientService{
		patientRepo: patientRepo,
		userRepo:    userRepo,
	}
}

func (s *PatientService) CreatePatient(req CreatePatientRequest, createdByID uint) (*models.PatientResponse, error) {
	_, err := s.userRepo.GetByID(createdByID)
	if err != nil {
		return nil, errors.New("invalid user")
	}

	// Parse date of birth
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return nil, errors.New("invalid date of birth format, use YYYY-MM-DD")
	}

	// Check if date of birth is not in the future
	if dob.After(time.Now()) {
		return nil, errors.New("date of birth cannot be in the future")
	}

	// Create patient
	patient := &models.Patient{
		FirstName:          req.FirstName,
		LastName:           req.LastName,
		Email:              req.Email,
		Phone:              req.Phone,
		DateOfBirth:        dob,
		Gender:             req.Gender,
		BloodType:          req.BloodType,
		Address:            req.Address,
		EmergencyContact:   req.EmergencyContact,
		MedicalHistory:     req.MedicalHistory,
		Allergies:          req.Allergies,
		CurrentMedications: req.CurrentMedications,
		CreatedByID:        createdByID,
		IsActive:           true,
	}

	if err := s.patientRepo.Create(patient); err != nil {
		return nil, errors.New("failed to create patient")
	}

	// Retrieve the created patient with relations
	createdPatient, err := s.patientRepo.GetByID(patient.ID)
	if err != nil {
		return nil, errors.New("failed to retrieve created patient")
	}

	response := createdPatient.ToResponse()
	return &response, nil
}

func (s *PatientService) GetPatientByID(id uint) (*models.PatientResponse, error) {
	patient, err := s.patientRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := patient.ToResponse()
	return &response, nil
}

func (s *PatientService) GetPatientByPatientID(patientID string) (*models.PatientResponse, error) {
	patient, err := s.patientRepo.GetByPatientID(patientID)
	if err != nil {
		return nil, err
	}

	response := patient.ToResponse()
	return &response, nil
}

func (s *PatientService) UpdatePatient(id uint, req UpdatePatientRequest, updatedByID uint, userRole models.UserRole) (*models.PatientResponse, error) {
	patient, err := s.patientRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	_, err = s.userRepo.GetByID(updatedByID)
	if err != nil {
		return nil, errors.New("invalid user")
	}

	if req.FirstName != nil {
		patient.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		patient.LastName = *req.LastName
	}
	if req.Email != nil {
		patient.Email = *req.Email
	}
	if req.Phone != nil {
		patient.Phone = *req.Phone
	}
	if req.DateOfBirth != nil {
		dob, err := time.Parse("2006-01-02", *req.DateOfBirth)
		if err != nil {
			return nil, errors.New("invalid date of birth format, use YYYY-MM-DD")
		}
		if dob.After(time.Now()) {
			return nil, errors.New("date of birth cannot be in the future")
		}
		patient.DateOfBirth = dob
	}
	if req.Gender != nil {
		patient.Gender = *req.Gender
	}
	if req.BloodType != nil {
		patient.BloodType = *req.BloodType
	}
	if req.Address != nil {
		patient.Address = *req.Address
	}
	if req.EmergencyContact != nil {
		patient.EmergencyContact = *req.EmergencyContact
	}

	// Only allow medical information updates for doctors
	if userRole == models.RoleDoctor {
		if req.MedicalHistory != nil {
			patient.MedicalHistory = *req.MedicalHistory
		}
		if req.Allergies != nil {
			patient.Allergies = *req.Allergies
		}
		if req.CurrentMedications != nil {
			patient.CurrentMedications = *req.CurrentMedications
		}
	}

	// Set last updated by
	patient.LastUpdatedByID = &updatedByID

	if err := s.patientRepo.Update(patient); err != nil {
		return nil, errors.New("failed to update patient")
	}

	// Retrieve updated patient with relations
	updatedPatient, err := s.patientRepo.GetByID(patient.ID)
	if err != nil {
		return nil, errors.New("failed to retrieve updated patient")
	}

	response := updatedPatient.ToResponse()
	return &response, nil
}

func (s *PatientService) DeletePatient(id uint, userRole models.UserRole) error {
	if userRole != models.RoleReceptionist {
		return errors.New("only receptionists can delete patients")
	}

	_, err := s.patientRepo.GetByID(id)
	if err != nil {
		return err
	}

	return s.patientRepo.Delete(id)
}

func (s *PatientService) ListPatients(page, pageSize int) (*PatientListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	patients, err := s.patientRepo.List(pageSize, offset)
	if err != nil {
		return nil, errors.New("failed to retrieve patients")
	}

	// Get total count
	total, err := s.patientRepo.Count()
	if err != nil {
		return nil, errors.New("failed to count patients")
	}

	// Convert to response format
	patientResponses := make([]models.PatientResponse, len(patients))
	for i, patient := range patients {
		patientResponses[i] = patient.ToResponse()
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &PatientListResponse{
		Patients: patientResponses,
		Pagination: PaginationResponse{
			Total:       total,
			CurrentPage: page,
			PageSize:    pageSize,
			TotalPages:  totalPages,
		},
	}, nil
}

func (s *PatientService) SearchPatients(query string, page, pageSize int) (*PatientListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	patients, err := s.patientRepo.Search(query, pageSize, offset)
	if err != nil {
		return nil, errors.New("failed to search patients")
	}

	patientResponses := make([]models.PatientResponse, len(patients))
	for i, patient := range patients {
		patientResponses[i] = patient.ToResponse()
	}

	return &PatientListResponse{
		Patients: patientResponses,
		Pagination: PaginationResponse{
			Total:       int64(len(patients)),
			CurrentPage: page,
			PageSize:    pageSize,
			TotalPages:  0,
		},
	}, nil
}
