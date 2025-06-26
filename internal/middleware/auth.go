package middleware

import (
	"net/http"
	"strings"

	"hospital-management-system/internal/auth"
	"hospital-management-system/internal/models"
	"hospital-management-system/pkg/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c, "Authorization header is required")
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.UnauthorizedResponse(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		token := tokenParts[1]
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			utils.UnauthorizedResponse(c, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)
		c.Set("claims", claims)

		c.Next()
	}
}

func RequireRole(allowedRoles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			utils.UnauthorizedResponse(c, "User role not found in context")
			c.Abort()
			return
		}

		role, ok := userRole.(models.UserRole)
		if !ok {
			utils.InternalErrorResponse(c, "Invalid user role format", nil)
			c.Abort()
			return
		}

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		utils.ForbiddenResponse(c, "Insufficient permissions")
		c.Abort()
	}
}

func RequireReceptionist() gin.HandlerFunc {
	return RequireRole(models.RoleReceptionist)
}

func RequireDoctor() gin.HandlerFunc {
	return RequireRole(models.RoleDoctor)
}

func RequireReceptionistOrDoctor() gin.HandlerFunc {
	return RequireRole(models.RoleReceptionist, models.RoleDoctor)
}

func GetUserFromContext(c *gin.Context) (uint, models.UserRole, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, "", gin.Error{Err: http.ErrMissingFile, Type: gin.ErrorTypePublic}
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		return 0, "", gin.Error{Err: http.ErrMissingFile, Type: gin.ErrorTypePublic}
	}

	id, ok := userID.(uint)
	if !ok {
		return 0, "", gin.Error{Err: http.ErrMissingFile, Type: gin.ErrorTypePublic}
	}

	role, ok := userRole.(models.UserRole)
	if !ok {
		return 0, "", gin.Error{Err: http.ErrMissingFile, Type: gin.ErrorTypePublic}
	}

	return id, role, nil
}
