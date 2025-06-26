package unit

import (
	"testing"
	"time"

	"hospital-management-system/internal/auth"
	"hospital-management-system/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWTService_GenerateToken_Success(t *testing.T) {
	jwtService := auth.NewJWTService("test_secret_key")
	user := &models.User{
		ID:       1,
		Username: "testuser",
		Role:     models.RoleDoctor,
	}

	token, err := jwtService.GenerateToken(user)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := jwtService.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Username, claims.Username)
	assert.Equal(t, user.Role, claims.Role)
	assert.Equal(t, "hospital-management-system", claims.Issuer)
	assert.Equal(t, user.Username, claims.Subject)
}

func TestJWTService_ValidateToken_Success(t *testing.T) {
	jwtService := auth.NewJWTService("test_secret_key")
	user := &models.User{
		ID:       1,
		Username: "testuser",
		Role:     models.RoleReceptionist,
	}

	token, err := jwtService.GenerateToken(user)
	assert.NoError(t, err)

	claims, err := jwtService.ValidateToken(token)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Username, claims.Username)
	assert.Equal(t, user.Role, claims.Role)
	assert.True(t, claims.ExpiresAt.After(time.Now()))
	assert.True(t, claims.IssuedAt.Before(time.Now().Add(time.Minute)))
}

func TestJWTService_ValidateToken_InvalidToken(t *testing.T) {
	jwtService := auth.NewJWTService("test_secret_key")

	claims, err := jwtService.ValidateToken("invalid_token")

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTService_ValidateToken_EmptyToken(t *testing.T) {
	jwtService := auth.NewJWTService("test_secret_key")

	claims, err := jwtService.ValidateToken("")

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTService_ValidateToken_WrongSecret(t *testing.T) {
	jwtService1 := auth.NewJWTService("secret1")
	jwtService2 := auth.NewJWTService("secret2")

	user := &models.User{
		ID:       1,
		Username: "testuser",
		Role:     models.RoleDoctor,
	}

	token, err := jwtService1.GenerateToken(user)
	assert.NoError(t, err)

	claims, err := jwtService2.ValidateToken(token)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTService_ValidateToken_ExpiredToken(t *testing.T) {
	jwtService := auth.NewJWTService("test_secret_key")

	claims := &auth.Claims{
		UserID:   1,
		Username: "testuser",
		Role:     models.RoleDoctor,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			Issuer:    "hospital-management-system",
			Subject:   "testuser",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("test_secret_key"))
	assert.NoError(t, err)

	// Execute
	validatedClaims, err := jwtService.ValidateToken(tokenString)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, validatedClaims)
	assert.Contains(t, err.Error(), "token is expired")
}

func TestJWTService_RefreshToken_Success(t *testing.T) {
	jwtService := auth.NewJWTService("test_secret_key")

	user := &models.User{
		ID:       1,
		Username: "testuser",
		Role:     models.RoleReceptionist,
	}

	originalToken, err := jwtService.GenerateToken(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, originalToken)

	time.Sleep(time.Second * 1)

	newToken, err := jwtService.RefreshToken(originalToken)
	assert.NoError(t, err)
	assert.NotEmpty(t, newToken)

	assert.NotEqual(t, originalToken, newToken, "New token should be different from original token")

	originalClaims, err := jwtService.ValidateToken(originalToken)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, originalClaims.UserID)

	newClaims, err := jwtService.ValidateToken(newToken)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, newClaims.UserID)
}

func TestJWTService_RefreshToken_InvalidToken(t *testing.T) {
	jwtService := auth.NewJWTService("test_secret_key")

	newToken, err := jwtService.RefreshToken("invalid_token")

	assert.Error(t, err)
	assert.Empty(t, newToken)
}

func TestJWTService_RefreshToken_ExpiredToken(t *testing.T) {
	jwtService := auth.NewJWTService("test_secret_key")

	claims := &auth.Claims{
		UserID:   1,
		Username: "testuser",
		Role:     models.RoleDoctor,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			Issuer:    "hospital-management-system",
			Subject:   "testuser",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("test_secret_key"))
	assert.NoError(t, err)

	newToken, err := jwtService.RefreshToken(tokenString)

	assert.Error(t, err)
	assert.Empty(t, newToken)
}

func TestJWTService_TokenExpirationTime(t *testing.T) {
	jwtService := auth.NewJWTService("test_secret_key")
	user := &models.User{
		ID:       1,
		Username: "testuser",
		Role:     models.RoleDoctor,
	}

	beforeGeneration := time.Now()

	token, err := jwtService.GenerateToken(user)
	assert.NoError(t, err)

	afterGeneration := time.Now()

	claims, err := jwtService.ValidateToken(token)
	assert.NoError(t, err)

	expectedMinExpiration := beforeGeneration.Add(24 * time.Hour).Add(-time.Minute)
	expectedMaxExpiration := afterGeneration.Add(24 * time.Hour).Add(time.Minute)

	assert.True(t, claims.ExpiresAt.After(expectedMinExpiration))
	assert.True(t, claims.ExpiresAt.Before(expectedMaxExpiration))
}

func TestJWTService_DifferentUsers(t *testing.T) {
	jwtService := auth.NewJWTService("test_secret_key")

	user1 := &models.User{
		ID:       1,
		Username: "doctor1",
		Role:     models.RoleDoctor,
	}

	user2 := &models.User{
		ID:       2,
		Username: "receptionist1",
		Role:     models.RoleReceptionist,
	}

	token1, err1 := jwtService.GenerateToken(user1)
	token2, err2 := jwtService.GenerateToken(user2)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, token1, token2)

	claims1, err1 := jwtService.ValidateToken(token1)
	claims2, err2 := jwtService.ValidateToken(token2)

	assert.NoError(t, err1)
	assert.NoError(t, err2)

	assert.Equal(t, user1.ID, claims1.UserID)
	assert.Equal(t, user1.Username, claims1.Username)
	assert.Equal(t, user1.Role, claims1.Role)

	assert.Equal(t, user2.ID, claims2.UserID)
	assert.Equal(t, user2.Username, claims2.Username)
	assert.Equal(t, user2.Role, claims2.Role)
}
