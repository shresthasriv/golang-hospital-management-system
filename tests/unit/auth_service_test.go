package unit

import (
	"errors"
	"testing"

	"hospital-management-system/internal/auth"
	"hospital-management-system/internal/models"
	"hospital-management-system/internal/services"
	"hospital-management-system/pkg/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) List(limit, offset int) ([]*models.User, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByRole(role models.UserRole) ([]*models.User, error) {
	args := m.Called(role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func TestAuthService_Login_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtService := auth.NewJWTService("test_secret")
	authService := services.NewAuthService(mockRepo, jwtService)

	hashedPassword, _ := utils.HashPassword("password123")
	user := &models.User{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  hashedPassword,
		FirstName: "Test",
		LastName:  "User",
		Role:      models.RoleReceptionist,
		IsActive:  true,
	}

	mockRepo.On("GetByUsername", "testuser").Return(user, nil)

	req := services.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	response, err := authService.Login(req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "testuser", response.User.Username)
	assert.NotEmpty(t, response.Token)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidUsername(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtService := auth.NewJWTService("test_secret")
	authService := services.NewAuthService(mockRepo, jwtService)

	mockRepo.On("GetByUsername", "nonexistent").Return(nil, errors.New("user not found"))

	req := services.LoginRequest{
		Username: "nonexistent",
		Password: "password123",
	}

	response, err := authService.Login(req)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "invalid username or password", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtService := auth.NewJWTService("test_secret")
	authService := services.NewAuthService(mockRepo, jwtService)

	hashedPassword, _ := utils.HashPassword("correct_password")
	user := &models.User{
		ID:       1,
		Username: "testuser",
		Password: hashedPassword,
		IsActive: true,
	}

	mockRepo.On("GetByUsername", "testuser").Return(user, nil)

	req := services.LoginRequest{
		Username: "testuser",
		Password: "wrong_password",
	}

	response, err := authService.Login(req)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "invalid username or password", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_InactiveUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtService := auth.NewJWTService("test_secret")
	authService := services.NewAuthService(mockRepo, jwtService)

	hashedPassword, _ := utils.HashPassword("password123")
	user := &models.User{
		ID:       1,
		Username: "testuser",
		Password: hashedPassword,
		IsActive: false,
	}

	mockRepo.On("GetByUsername", "testuser").Return(user, nil)

	req := services.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	response, err := authService.Login(req)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "user account is deactivated", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Register_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtService := auth.NewJWTService("test_secret")
	authService := services.NewAuthService(mockRepo, jwtService)

	mockRepo.On("GetByUsername", "newuser").Return(nil, errors.New("user not found"))
	mockRepo.On("GetByEmail", "new@example.com").Return(nil, errors.New("user not found"))
	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

	req := services.RegisterRequest{
		Username:  "newuser",
		Email:     "new@example.com",
		Password:  "password123",
		FirstName: "New",
		LastName:  "User",
		Role:      models.RoleDoctor,
	}

	response, err := authService.Register(req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "newuser", response.Username)
	assert.Equal(t, "new@example.com", response.Email)
	assert.Equal(t, models.RoleDoctor, response.Role)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Register_UsernameExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtService := auth.NewJWTService("test_secret")
	authService := services.NewAuthService(mockRepo, jwtService)

	existingUser := &models.User{
		Username: "existinguser",
	}

	mockRepo.On("GetByUsername", "existinguser").Return(existingUser, nil)

	req := services.RegisterRequest{
		Username:  "existinguser",
		Email:     "new@example.com",
		Password:  "password123",
		FirstName: "New",
		LastName:  "User",
		Role:      models.RoleDoctor,
	}

	response, err := authService.Register(req)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "username already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Register_EmailExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtService := auth.NewJWTService("test_secret")
	authService := services.NewAuthService(mockRepo, jwtService)

	existingUser := &models.User{
		Email: "existing@example.com",
	}

	mockRepo.On("GetByUsername", "newuser").Return(nil, errors.New("user not found"))
	mockRepo.On("GetByEmail", "existing@example.com").Return(existingUser, nil)

	req := services.RegisterRequest{
		Username:  "newuser",
		Email:     "existing@example.com",
		Password:  "password123",
		FirstName: "New",
		LastName:  "User",
		Role:      models.RoleDoctor,
	}

	response, err := authService.Register(req)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "email already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestAuthService_GetUserByID_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtService := auth.NewJWTService("test_secret")
	authService := services.NewAuthService(mockRepo, jwtService)

	user := &models.User{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Role:      models.RoleReceptionist,
		IsActive:  true,
	}

	mockRepo.On("GetByID", uint(1)).Return(user, nil)

	response, err := authService.GetUserByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, "testuser", response.Username)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_GetUserByID_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtService := auth.NewJWTService("test_secret")
	authService := services.NewAuthService(mockRepo, jwtService)

	mockRepo.On("GetByID", uint(999)).Return(nil, errors.New("user not found"))

	response, err := authService.GetUserByID(999)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "user not found", err.Error())
	mockRepo.AssertExpectations(t)
}
