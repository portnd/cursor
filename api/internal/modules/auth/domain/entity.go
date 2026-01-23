package domain

import (
	"time"
)

// User represents the user entity in the system
// This is the core domain model for authentication
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"` // Unique email for login
	Password  string    `json:"-" gorm:"not null"`                 // Password hash (never expose in JSON)
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName overrides the default table name
func (User) TableName() string {
	return "users"
}

// RegisterRequest is the DTO for user registration
type RegisterRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

// LoginRequest is the DTO for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse is the response structure for authentication endpoints
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// Usecase defines the authentication business logic interface
// This follows the Dependency Inversion Principle (Hexagonal Architecture)
type Usecase interface {
	Register(req *RegisterRequest) (*AuthResponse, error)
	Login(req *LoginRequest) (*AuthResponse, error)
}

// Repository defines the data access interface for user operations
// This is the PORT in Hexagonal Architecture
type Repository interface {
	CreateUser(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
}
