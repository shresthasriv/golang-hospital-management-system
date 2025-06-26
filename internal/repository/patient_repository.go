package repository

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"hospital-management-system/internal/models"

	"gorm.io/gorm"
)

type PatientRepository interface {
	Create(patient *models.Patient) error
	GetByID(id uint) (*models.Patient, error)
	GetByPatientID(patientID string) (*models.Patient, error)
	Update(patient *models.Patient) error
	Delete(id uint) error
	List(limit, offset int) ([]*models.Patient, error)
	Search(query string, limit, offset int) ([]*models.Patient, error)
	GetByCreatedBy(userID uint, limit, offset int) ([]*models.Patient, error)
	Count() (int64, error)
	GenerateUniquePatientID() (string, error)
}

type patientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) PatientRepository {
	return &patientRepository{db: db}
}

func (r *patientRepository) Create(patient *models.Patient) error {
	if patient.PatientID == "" {
		patientID, err := r.GenerateUniquePatientID()
		if err != nil {
			return err
		}
		patient.PatientID = patientID
	}

	if err := r.db.Create(patient).Error; err != nil {
		return err
	}
	return nil
}

func (r *patientRepository) GetByID(id uint) (*models.Patient, error) {
	var patient models.Patient
	if err := r.db.Preload("CreatedBy").Preload("LastUpdatedBy").
		Where("id = ? AND is_active = ?", id, true).First(&patient).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("patient not found")
		}
		return nil, err
	}
	return &patient, nil
}

func (r *patientRepository) GetByPatientID(patientID string) (*models.Patient, error) {
	var patient models.Patient
	if err := r.db.Preload("CreatedBy").Preload("LastUpdatedBy").
		Where("patient_id = ? AND is_active = ?", patientID, true).First(&patient).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("patient not found")
		}
		return nil, err
	}
	return &patient, nil
}

func (r *patientRepository) Update(patient *models.Patient) error {
	return r.db.Save(patient).Error
}

func (r *patientRepository) Delete(id uint) error {
	return r.db.Model(&models.Patient{}).Where("id = ?", id).Update("is_active", false).Error
}

func (r *patientRepository) List(limit, offset int) ([]*models.Patient, error) {
	var patients []*models.Patient
	query := r.db.Preload("CreatedBy").Preload("LastUpdatedBy").
		Where("is_active = ?", true).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&patients).Error; err != nil {
		return nil, err
	}

	return patients, nil
}

func (r *patientRepository) Search(query string, limit, offset int) ([]*models.Patient, error) {
	var patients []*models.Patient
	searchTerm := "%" + strings.ToLower(query) + "%"

	dbQuery := r.db.Preload("CreatedBy").Preload("LastUpdatedBy").
		Where("is_active = ?", true).
		Where("(LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ? OR LOWER(email) LIKE ? OR phone LIKE ? OR patient_id LIKE ?)",
			searchTerm, searchTerm, searchTerm, searchTerm, searchTerm).
		Order("created_at DESC")

	if limit > 0 {
		dbQuery = dbQuery.Limit(limit)
	}
	if offset > 0 {
		dbQuery = dbQuery.Offset(offset)
	}

	if err := dbQuery.Find(&patients).Error; err != nil {
		return nil, err
	}

	return patients, nil
}

func (r *patientRepository) GetByCreatedBy(userID uint, limit, offset int) ([]*models.Patient, error) {
	var patients []*models.Patient
	query := r.db.Preload("CreatedBy").Preload("LastUpdatedBy").
		Where("created_by_id = ? AND is_active = ?", userID, true).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&patients).Error; err != nil {
		return nil, err
	}

	return patients, nil
}

func (r *patientRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Patient{}).Where("is_active = ?", true).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *patientRepository) GenerateUniquePatientID() (string, error) {
	for attempts := 0; attempts < 10; attempts++ {
		now := time.Now()
		dateStr := now.Format("20060102")

		var count int64
		if err := r.db.Model(&models.Patient{}).Where("patient_id LIKE ?", "PAT"+dateStr+"%").Count(&count).Error; err != nil {
			return "", err
		}

		patientID := fmt.Sprintf("PAT%s%04d", dateStr, count+1)

		var existingPatient models.Patient
		if err := r.db.Where("patient_id = ?", patientID).First(&existingPatient).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return patientID, nil
			}
			return "", err
		}
	}

	return "", errors.New("failed to generate unique patient ID after multiple attempts")
}
