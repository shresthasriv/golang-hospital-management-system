package handlers

import (
	"hospital-management-system/internal/services"
	"hospital-management-system/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request data", err)
		return
	}

	response, err := h.authService.Login(req)
	if err != nil {
		utils.UnauthorizedResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", response)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request data", err)
		return
	}

	response, err := h.authService.Register(req)
	if err != nil {
		if err.Error() == "username already exists" || err.Error() == "email already exists" {
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), err)
			return
		}
		utils.InternalErrorResponse(c, "Failed to register user", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", response)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.UnauthorizedResponse(c, "Authorization header required")
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		utils.UnauthorizedResponse(c, "Invalid authorization header format")
		return
	}

	newToken, err := h.authService.RefreshToken(tokenParts[1])
	if err != nil {
		utils.UnauthorizedResponse(c, "Invalid or expired token")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Token refreshed successfully", gin.H{
		"token": newToken,
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.InternalErrorResponse(c, "User ID not found in context", nil)
		return
	}

	id, ok := userID.(uint)
	if !ok {
		utils.InternalErrorResponse(c, "Invalid user ID format", nil)
		return
	}

	user, err := h.authService.GetUserByID(id)
	if err != nil {
		utils.NotFoundResponse(c, "User not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", user)
}
