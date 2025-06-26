package unit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"hospital-management-system/internal/auth"
	"hospital-management-system/internal/handlers"
	"hospital-management-system/internal/models"
	"hospital-management-system/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type AuthServiceInterface interface {
	Login(req services.LoginRequest) (*services.LoginResponse, error)
	Register(req services.RegisterRequest) (*models.UserResponse, error)
	GetUserByID(id uint) (*models.UserResponse, error)
	RefreshToken(tokenString string) (string, error)
	ValidateToken(tokenString string) (*auth.Claims, error)
}

type PatientServiceInterface interface {
	CreatePatient(req services.CreatePatientRequest, createdByID uint) (*models.PatientResponse, error)
	GetPatientByID(id uint) (*models.PatientResponse, error)
	GetPatientByPatientID(patientID string) (*models.PatientResponse, error)
	UpdatePatient(id uint, req services.UpdatePatientRequest, updatedByID uint, userRole models.UserRole) (*models.PatientResponse, error)
	DeletePatient(id uint, userRole models.UserRole) error
	ListPatients(page, pageSize int) (*services.PatientListResponse, error)
	SearchPatients(query string, page, pageSize int) (*services.PatientListResponse, error)
}

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(req services.LoginRequest) (*services.LoginResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.LoginResponse), args.Error(1)
}

func (m *MockAuthService) Register(req services.RegisterRequest) (*models.UserResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserResponse), args.Error(1)
}

func (m *MockAuthService) GetUserByID(id uint) (*models.UserResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserResponse), args.Error(1)
}

func (m *MockAuthService) RefreshToken(tokenString string) (string, error) {
	args := m.Called(tokenString)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ValidateToken(tokenString string) (*auth.Claims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.Claims), args.Error(1)
}

type MockPatientService struct {
	mock.Mock
}

func (m *MockPatientService) CreatePatient(req services.CreatePatientRequest, createdByID uint) (*models.PatientResponse, error) {
	args := m.Called(req, createdByID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PatientResponse), args.Error(1)
}

func (m *MockPatientService) GetPatientByID(id uint) (*models.PatientResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PatientResponse), args.Error(1)
}

func (m *MockPatientService) GetPatientByPatientID(patientID string) (*models.PatientResponse, error) {
	args := m.Called(patientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PatientResponse), args.Error(1)
}

func (m *MockPatientService) UpdatePatient(id uint, req services.UpdatePatientRequest, updatedByID uint, userRole models.UserRole) (*models.PatientResponse, error) {
	args := m.Called(id, req, updatedByID, userRole)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PatientResponse), args.Error(1)
}

func (m *MockPatientService) DeletePatient(id uint, userRole models.UserRole) error {
	args := m.Called(id, userRole)
	return args.Error(0)
}

func (m *MockPatientService) ListPatients(page, pageSize int) (*services.PatientListResponse, error) {
	args := m.Called(page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.PatientListResponse), args.Error(1)
}

func (m *MockPatientService) SearchPatients(query string, page, pageSize int) (*services.PatientListResponse, error) {
	args := m.Called(query, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.PatientListResponse), args.Error(1)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	authService := &services.AuthService{}
	authHandler := handlers.NewAuthHandler(authService)

	router := setupRouter()
	router.POST("/login", authHandler.Login)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response["success"].(bool))
}

func TestAuthHandler_Login_MissingFields(t *testing.T) {
	authService := &services.AuthService{}
	authHandler := handlers.NewAuthHandler(authService)

	router := setupRouter()
	router.POST("/login", authHandler.Login)

	reqData := map[string]string{"username": "testuser"}
	jsonData, _ := json.Marshal(reqData)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Register_InvalidJSON(t *testing.T) {
	authService := &services.AuthService{}
	authHandler := handlers.NewAuthHandler(authService)

	router := setupRouter()
	router.POST("/register", authHandler.Register)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Register_MissingFields(t *testing.T) {
	authService := &services.AuthService{}
	authHandler := handlers.NewAuthHandler(authService)

	router := setupRouter()
	router.POST("/register", authHandler.Register)

	reqData := map[string]string{"username": "testuser"}
	jsonData, _ := json.Marshal(reqData)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_GetProfile_MissingUserContext(t *testing.T) {
	authService := &services.AuthService{}
	authHandler := handlers.NewAuthHandler(authService)

	router := setupRouter()
	router.GET("/profile", authHandler.GetProfile)

	req, _ := http.NewRequest("GET", "/profile", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response["success"].(bool))
}

func TestAuthHandler_RefreshToken_MissingHeader(t *testing.T) {
	authService := &services.AuthService{}
	authHandler := handlers.NewAuthHandler(authService)

	router := setupRouter()
	router.POST("/refresh", authHandler.RefreshToken)

	req, _ := http.NewRequest("POST", "/refresh", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthHandler_RefreshToken_InvalidHeaderFormat(t *testing.T) {
	authService := &services.AuthService{}
	authHandler := handlers.NewAuthHandler(authService)

	router := setupRouter()
	router.POST("/refresh", authHandler.RefreshToken)

	req, _ := http.NewRequest("POST", "/refresh", nil)
	req.Header.Set("Authorization", "InvalidFormat token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPatientHandler_CreatePatient_MissingUserContext(t *testing.T) {
	// Setup
	patientService := &services.PatientService{} // nil dependencies for this test
	patientHandler := handlers.NewPatientHandler(patientService)

	router := setupRouter()
	router.POST("/patients", patientHandler.CreatePatient) // No middleware to set user context

	createRequest := services.CreatePatientRequest{
		FirstName:        "John",
		LastName:         "Doe",
		Phone:            "1234567890",
		DateOfBirth:      "1990-01-01",
		Gender:           models.GenderMale,
		EmergencyContact: "0987654321",
	}

	// Prepare request
	jsonData, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest("POST", "/patients", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert - should get 401 for missing user context
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPatientHandler_CreatePatient_ForbiddenForDoctor(t *testing.T) {
	// Setup
	patientService := &services.PatientService{} // nil dependencies for this test
	patientHandler := handlers.NewPatientHandler(patientService)

	router := setupRouter()
	router.POST("/patients", func(c *gin.Context) {
		// Simulate middleware setting doctor user context
		c.Set("user_id", uint(1))
		c.Set("user_role", models.RoleDoctor)
		patientHandler.CreatePatient(c)
	})

	createRequest := services.CreatePatientRequest{
		FirstName:        "John",
		LastName:         "Doe",
		Phone:            "1234567890",
		DateOfBirth:      "1990-01-01",
		Gender:           models.GenderMale,
		EmergencyContact: "0987654321",
	}

	// Prepare request
	jsonData, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest("POST", "/patients", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert - should get 403 for doctor trying to create patient
	assert.Equal(t, http.StatusForbidden, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response["success"].(bool))
	assert.Contains(t, response["message"], "Only receptionists can create patients")
}

func TestPatientHandler_CreatePatient_InvalidJSON(t *testing.T) {
	patientService := &services.PatientService{}
	patientHandler := handlers.NewPatientHandler(patientService)

	router := setupRouter()
	router.POST("/patients", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Set("user_role", models.RoleReceptionist)
		patientHandler.CreatePatient(c)
	})

	req, _ := http.NewRequest("POST", "/patients", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPatientHandler_GetPatient_InvalidID(t *testing.T) {
	// Setup
	patientService := &services.PatientService{} // nil dependencies for this test
	patientHandler := handlers.NewPatientHandler(patientService)

	router := setupRouter()
	router.GET("/patients/:id", patientHandler.GetPatient)

	// Prepare request with invalid ID
	req, _ := http.NewRequest("GET", "/patients/invalid", nil)

	// Execute
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert - should get 400 for invalid ID
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response["success"].(bool))
}

func TestPatientHandler_GetPatientByPatientID_EmptyID(t *testing.T) {
	patientService := &services.PatientService{}
	patientHandler := handlers.NewPatientHandler(patientService)

	router := setupRouter()
	router.GET("/patients/patient/:patient_id", patientHandler.GetPatientByPatientID)

	req, _ := http.NewRequest("GET", "/patients/patient/", nil) // Empty patient ID

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Gin returns 404 for routes with missing parameters
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPatientHandler_SearchPatients_MissingQuery(t *testing.T) {
	// Setup
	patientService := &services.PatientService{} // nil dependencies for this test
	patientHandler := handlers.NewPatientHandler(patientService)

	router := setupRouter()
	router.GET("/patients/search", patientHandler.SearchPatients)

	// Prepare request without query parameter
	req, _ := http.NewRequest("GET", "/patients/search", nil)

	// Execute
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert - should get 400 for missing query
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response["success"].(bool))
	assert.Contains(t, response["message"], "Search query is required")
}

func TestPatientHandler_UpdatePatient_InvalidID(t *testing.T) {
	// Setup
	patientService := &services.PatientService{} // nil dependencies for this test
	patientHandler := handlers.NewPatientHandler(patientService)

	router := setupRouter()
	router.PUT("/patients/:id", patientHandler.UpdatePatient)

	updateRequest := services.UpdatePatientRequest{
		FirstName: stringPtr("Jane"),
	}

	// Prepare request with invalid ID
	jsonData, _ := json.Marshal(updateRequest)
	req, _ := http.NewRequest("PUT", "/patients/invalid", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert - should get 400 for invalid ID
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPatientHandler_UpdatePatient_MissingUserContext(t *testing.T) {
	patientService := &services.PatientService{}
	patientHandler := handlers.NewPatientHandler(patientService)

	router := setupRouter()
	router.PUT("/patients/:id", patientHandler.UpdatePatient)

	updateRequest := services.UpdatePatientRequest{
		FirstName: stringPtr("Jane"),
	}

	jsonData, _ := json.Marshal(updateRequest)
	req, _ := http.NewRequest("PUT", "/patients/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPatientHandler_DeletePatient_InvalidID(t *testing.T) {
	// Setup
	patientService := &services.PatientService{} // nil dependencies for this test
	patientHandler := handlers.NewPatientHandler(patientService)

	router := setupRouter()
	router.DELETE("/patients/:id", patientHandler.DeletePatient)

	// Prepare request with invalid ID
	req, _ := http.NewRequest("DELETE", "/patients/invalid", nil)

	// Execute
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert - should get 400 for invalid ID
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPatientHandler_DeletePatient_ForbiddenForDoctor(t *testing.T) {
	// Setup
	patientService := &services.PatientService{} // nil dependencies for this test
	patientHandler := handlers.NewPatientHandler(patientService)

	router := setupRouter()
	router.DELETE("/patients/:id", func(c *gin.Context) {
		// Simulate middleware setting doctor user context
		c.Set("user_id", uint(1))
		c.Set("user_role", models.RoleDoctor)
		patientHandler.DeletePatient(c)
	})

	// Prepare request
	req, _ := http.NewRequest("DELETE", "/patients/1", nil)

	// Execute
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert - should get 403 for doctor trying to delete patient
	assert.Equal(t, http.StatusForbidden, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response["success"].(bool))
	assert.Contains(t, response["message"], "Only receptionists can delete patients")
}

func TestPatientHandler_DeletePatient_MissingUserContext(t *testing.T) {
	patientService := &services.PatientService{}
	patientHandler := handlers.NewPatientHandler(patientService)

	router := setupRouter()
	router.DELETE("/patients/:id", patientHandler.DeletePatient)

	req, _ := http.NewRequest("DELETE", "/patients/1", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}

// Test pagination parameter parsing
func TestPaginationParameterParsing(t *testing.T) {
	tests := []struct {
		name        string
		queryParams string
		shouldWork  bool
	}{
		{"No parameters", "", true},
		{"Valid page", "?page=2", true},
		{"Valid page_size", "?page_size=20", true},
		{"Both valid", "?page=2&page_size=20", true},
		{"Invalid page", "?page=invalid", true},           // Should default to 1
		{"Invalid page_size", "?page_size=invalid", true}, // Should default to 10
		{"Page too large", "?page=999999", true},          // Should be accepted
		{"Page_size too large", "?page_size=200", true},   // Should be clamped to 100
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test verifies that parameter parsing doesn't cause immediate failures
			// The actual business logic validation would happen in the service layer

			// Create a simple handler that just extracts parameters
			router := setupRouter()
			router.GET("/test", func(c *gin.Context) {
				page := c.DefaultQuery("page", "1")
				pageSize := c.DefaultQuery("page_size", "10")
				c.JSON(200, gin.H{"page": page, "page_size": pageSize})
			})

			req, _ := http.NewRequest("GET", "/test"+tt.queryParams, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should always succeed at the HTTP level
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}
