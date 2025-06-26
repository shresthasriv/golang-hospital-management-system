package models

import (
	"time"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleReceptionist UserRole = "receptionist"
	RoleDoctor      UserRole = "doctor"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null" binding:"required,min=3,max=50"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null" binding:"required,email"`
	Password  string         `json:"-" gorm:"not null" binding:"required,min=6"`
	FirstName string         `json:"first_name" gorm:"not null" binding:"required,min=2,max=50"`
	LastName  string         `json:"last_name" gorm:"not null" binding:"required,min=2,max=50"`
	Role      UserRole       `json:"role" gorm:"not null" binding:"required,oneof=receptionist doctor"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (User) TableName() string {
	return "users"
}

type UserResponse struct {
	ID        uint     `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Role      UserRole `json:"role"`
	IsActive  bool     `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
} 