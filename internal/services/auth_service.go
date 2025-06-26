package services

import (
	"errors"

	"hospital-management-system/internal/auth"
	"hospital-management-system/internal/models"
	"hospital-management-system/internal/repository"
	"hospital-management-system/pkg/utils"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	User  models.UserResponse `json:"user"`
	Token string              `json:"token"`
}

type RegisterRequest struct {
	Username  string          `json:"username" binding:"required,min=3,max=50"`
	Email     string          `json:"email" binding:"required,email"`
	Password  string          `json:"password" binding:"required,min=6"`
	FirstName string          `json:"first_name" binding:"required,min=2,max=50"`
	LastName  string          `json:"last_name" binding:"required,min=2,max=50"`
	Role      models.UserRole `json:"role" binding:"required,oneof=receptionist doctor"`
}

type AuthService struct {
	userRepo   repository.UserRepository
	jwtService *auth.JWTService
}

func NewAuthService(userRepo repository.UserRepository, jwtService *auth.JWTService) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *AuthService) Login(req LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if !utils.VerifyPassword(user.Password, req.Password) {
		return nil, errors.New("invalid username or password")
	}

	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	token, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &LoginResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

func (s *AuthService) Register(req RegisterRequest) (*models.UserResponse, error) {
	existingUser, err := s.userRepo.GetByUsername(req.Username)
	if err == nil && existingUser != nil {
		return nil, errors.New("username already exists")
	}

	existingUser, err = s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to process password")
	}

	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
		IsActive:  true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *AuthService) GetUserByID(id uint) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *AuthService) RefreshToken(tokenString string) (string, error) {
	return s.jwtService.RefreshToken(tokenString)
}

func (s *AuthService) ValidateToken(tokenString string) (*auth.Claims, error) {
	return s.jwtService.ValidateToken(tokenString)
}
