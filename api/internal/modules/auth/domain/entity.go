package domain

import (
	"time"

	"github.com/lib/pq"
)

// User represents the user entity in the system
// This is the core domain model for authentication and Sentinel role management
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"` // Unique email for login
	Password  string    `json:"-" gorm:"not null"`                 // Password hash (never expose in JSON)
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Sentinel enhancements
	Role        string         `json:"role" gorm:"type:varchar(20);not null;default:'DEV'"` // CEO, PM, or DEV
	HealthScore float64        `json:"health_score" gorm:"type:decimal(5,2);default:100.00"` // Performance tracking (0-100)
	TechStack   pq.StringArray `json:"tech_stack" gorm:"type:text[]"`                        // Array of technologies
}

// UserRole constants for type safety
const (
	RoleCEO = "CEO"
	RolePM  = "PM"
	RoleDEV = "DEV"
)

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

// ChangeRoleRequest is the DTO for changing user roles (CEO only)
type ChangeRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=CEO PM DEV"`
}

// Usecase defines the authentication business logic interface
// This follows the Dependency Inversion Principle (Hexagonal Architecture)
type Usecase interface {
	Register(req *RegisterRequest) (*AuthResponse, error)
	Login(req *LoginRequest) (*AuthResponse, error)
	// User Management (CEO only)
	GetTeamMembers(requestingUserID uint) ([]User, error)
	ChangeUserRole(requestingUserID uint, targetUserID uint, newRole string) error
}

// Repository defines the data access interface for user operations
// This is the PORT in Hexagonal Architecture
type Repository interface {
	CreateUser(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
	// User Management
	GetAllUsers() ([]User, error)
	UpdateUserRole(userID uint, newRole string) error
}
