package handlers

import (
	"net/http"
	"strconv"

	"hospital-management-system/internal/middleware"
	"hospital-management-system/internal/models"
	"hospital-management-system/internal/services"
	"hospital-management-system/pkg/utils"

	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	patientService *services.PatientService
}

func NewPatientHandler(patientService *services.PatientService) *PatientHandler {
	return &PatientHandler{
		patientService: patientService,
	}
}

func (h *PatientHandler) CreatePatient(c *gin.Context) {
	var req services.CreatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request data", err)
		return
	}

	userID, userRole, err := middleware.GetUserFromContext(c)
	if err != nil {
		utils.UnauthorizedResponse(c, "User context not found")
		return
	}

	if userRole != models.RoleReceptionist {
		utils.ForbiddenResponse(c, "Only receptionists can create patients")
		return
	}

	patient, err := h.patientService.CreatePatient(req, userID)
	if err != nil {
		utils.InternalErrorResponse(c, "Failed to create patient", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Patient created successfully", patient)
}

func (h *PatientHandler) GetPatient(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid patient ID", err)
		return
	}

	patient, err := h.patientService.GetPatientByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Patient not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Patient retrieved successfully", patient)
}

func (h *PatientHandler) GetPatientByPatientID(c *gin.Context) {
	patientID := c.Param("patient_id")
	if patientID == "" {
		utils.ValidationErrorResponse(c, "Patient ID is required", nil)
		return
	}

	patient, err := h.patientService.GetPatientByPatientID(patientID)
	if err != nil {
		utils.NotFoundResponse(c, "Patient not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Patient retrieved successfully", patient)
}

func (h *PatientHandler) UpdatePatient(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid patient ID", err)
		return
	}

	var req services.UpdatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request data", err)
		return
	}

	userID, userRole, err := middleware.GetUserFromContext(c)
	if err != nil {
		utils.UnauthorizedResponse(c, "User context not found")
		return
	}

	patient, err := h.patientService.UpdatePatient(uint(id), req, userID, userRole)
	if err != nil {
		utils.InternalErrorResponse(c, "Failed to update patient", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Patient updated successfully", patient)
}

func (h *PatientHandler) DeletePatient(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid patient ID", err)
		return
	}

	_, userRole, err := middleware.GetUserFromContext(c)
	if err != nil {
		utils.UnauthorizedResponse(c, "User context not found")
		return
	}

	if userRole != models.RoleReceptionist {
		utils.ForbiddenResponse(c, "Only receptionists can delete patients")
		return
	}

	err = h.patientService.DeletePatient(uint(id), userRole)
	if err != nil {
		utils.InternalErrorResponse(c, "Failed to delete patient", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Patient deleted successfully", nil)
}

func (h *PatientHandler) ListPatients(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	patients, err := h.patientService.ListPatients(page, pageSize)
	if err != nil {
		utils.InternalErrorResponse(c, "Failed to retrieve patients", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Patients retrieved successfully", patients)
}

func (h *PatientHandler) SearchPatients(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.ValidationErrorResponse(c, "Search query is required", nil)
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	patients, err := h.patientService.SearchPatients(query, page, pageSize)
	if err != nil {
		utils.InternalErrorResponse(c, "Failed to search patients", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Patient search completed successfully", patients)
}
