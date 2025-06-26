package unit

import (
	"testing"

	"hospital-management-system/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword_Success(t *testing.T) {
	// Setup
	password := "test_password_123"

	// Execute
	hashedPassword, err := utils.HashPassword(password)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, password, hashedPassword)
	assert.True(t, len(hashedPassword) > 50) // bcrypt hashes are typically 60 characters
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	// Setup
	password := ""

	// Execute
	hashedPassword, err := utils.HashPassword(password)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
}

func TestHashPassword_DifferentPasswords_DifferentHashes(t *testing.T) {
	// Setup
	password1 := "password123"
	password2 := "password456"

	// Execute
	hash1, err1 := utils.HashPassword(password1)
	hash2, err2 := utils.HashPassword(password2)

	// Assert
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, hash1, hash2)
}

func TestHashPassword_SamePassword_DifferentHashes(t *testing.T) {
	// Setup
	password := "same_password"

	// Execute
	hash1, err1 := utils.HashPassword(password)
	hash2, err2 := utils.HashPassword(password)

	// Assert
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, hash1, hash2) // bcrypt includes salt, so same password gives different hashes
}

func TestVerifyPassword_Success(t *testing.T) {
	// Setup
	password := "test_password_123"
	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)

	// Execute
	isValid := utils.VerifyPassword(hashedPassword, password)

	// Assert
	assert.True(t, isValid)
}

func TestVerifyPassword_WrongPassword(t *testing.T) {
	// Setup
	password := "correct_password"
	wrongPassword := "wrong_password"
	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)

	// Execute
	isValid := utils.VerifyPassword(hashedPassword, wrongPassword)

	// Assert
	assert.False(t, isValid)
}

func TestVerifyPassword_EmptyPassword(t *testing.T) {
	// Setup
	password := "test_password"
	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)

	// Execute
	isValid := utils.VerifyPassword(hashedPassword, "")

	// Assert
	assert.False(t, isValid)
}

func TestVerifyPassword_EmptyHash(t *testing.T) {
	// Setup
	password := "test_password"

	// Execute
	isValid := utils.VerifyPassword("", password)

	// Assert
	assert.False(t, isValid)
}

func TestVerifyPassword_InvalidHash(t *testing.T) {
	// Setup
	password := "test_password"
	invalidHash := "invalid_hash_format"

	// Execute
	isValid := utils.VerifyPassword(invalidHash, password)

	// Assert
	assert.False(t, isValid)
}

func TestVerifyPassword_MultiplePasswords(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "short password",
			password: "123",
		},
		{
			name:     "long password",
			password: "this_is_a_very_long_password_with_many_characters_123456789",
		},
		{
			name:     "special characters",
			password: "p@ssw0rd!#$%^&*()",
		},
		{
			name:     "unicode characters",
			password: "пароль123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Hash the password
			hashedPassword, err := utils.HashPassword(tt.password)
			assert.NoError(t, err)

			// Verify correct password
			assert.True(t, utils.VerifyPassword(hashedPassword, tt.password))

			// Verify wrong password
			assert.False(t, utils.VerifyPassword(hashedPassword, tt.password+"_wrong"))
		})
	}
}

func TestPasswordHashingIntegration(t *testing.T) {
	// Setup
	passwords := []string{
		"password123",
		"admin",
		"very_secure_password_2024",
		"P@ssw0rd!",
	}

	var hashedPasswords []string

	// Hash all passwords
	for _, password := range passwords {
		hashed, err := utils.HashPassword(password)
		assert.NoError(t, err)
		hashedPasswords = append(hashedPasswords, hashed)
	}

	// Verify all passwords with their corresponding hashes
	for i, password := range passwords {
		assert.True(t, utils.VerifyPassword(hashedPasswords[i], password))
	}

	// Verify passwords don't match wrong hashes
	for i, password := range passwords {
		for j, hash := range hashedPasswords {
			if i != j {
				assert.False(t, utils.VerifyPassword(hash, password))
			}
		}
	}
}
