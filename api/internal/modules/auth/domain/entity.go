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
	DisplayName string         `json:"display_name" gorm:"type:varchar(100)"`                // Optional display name (enterprise profile)
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

// CreateUserRequest is the DTO for admin creating a single user (CEO only)
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"required,oneof=CEO PM DEV"`
}

// ImportUserItem is one row in a bulk user import (CEO only)
// Password is optional; if empty, a random temporary password is generated and returned.
type ImportUserItem struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password"`                    // optional; min 8 if provided
	Role     string `json:"role" binding:"omitempty,oneof=CEO PM DEV"` // default DEV
}

// ImportUsersRequest is the request body for bulk user import
type ImportUsersRequest struct {
	Users []ImportUserItem `json:"users" binding:"required,dive"`
}

// ImportUserResult represents the outcome of one user in a bulk import
type ImportUserResult struct {
	Email        string `json:"email"`
	Status       string `json:"status"` // "created" | "skipped" (duplicate) | "error"
	Message      string `json:"message,omitempty"`
	User         *User  `json:"user,omitempty"`
	TempPassword string `json:"temp_password,omitempty"` // only when password was auto-generated
}

// ImportUsersResponse is the response for bulk import
type ImportUsersResponse struct {
	Total   int                `json:"total"`
	Created int                `json:"created"`
	Skipped int                `json:"skipped"`
	Errors  int                `json:"errors"`
	Results []ImportUserResult `json:"results"`
}

// UpdateProfileRequest is the DTO for updating own profile (any authenticated user)
type UpdateProfileRequest struct {
	DisplayName *string  `json:"display_name"`
	TechStack   []string `json:"tech_stack"`
}

// ChangePasswordRequest is the DTO for changing own password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// Usecase defines the authentication business logic interface
// This follows the Dependency Inversion Principle (Hexagonal Architecture)
type Usecase interface {
	Register(req *RegisterRequest) (*AuthResponse, error)
	Login(req *LoginRequest) (*AuthResponse, error)
	// Profile (any authenticated user)
	GetProfile(userID uint) (*User, error)
	UpdateProfile(userID uint, req *UpdateProfileRequest) (*User, error)
	ChangePassword(userID uint, currentPassword, newPassword string) error
	// User Management (CEO only)
	GetTeamMembers(requestingUserID uint) ([]User, error)
	ChangeUserRole(requestingUserID uint, targetUserID uint, newRole string) error
	// Admin user creation (CEO only)
	CreateUserAsAdmin(requestingUserID uint, req *CreateUserRequest) (*User, error)
	ImportUsers(requestingUserID uint, req *ImportUsersRequest) (*ImportUsersResponse, error)
	// Delete user (CEO only; cannot delete self)
	DeleteUser(requestingUserID uint, targetUserID uint) error
	// Reset user password (CEO only). Returns the new temporary password.
	ResetUserPassword(requestingUserID uint, targetUserID uint) (tempPassword string, err error)
}

// Repository defines the data access interface for user operations
// This is the PORT in Hexagonal Architecture
type Repository interface {
	CreateUser(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
	// Profile
	UpdateProfile(userID uint, displayName *string, techStack []string) error
	// User Management
	GetAllUsers() ([]User, error)
	UpdateUserRole(userID uint, newRole string) error
	DeleteUser(userID uint) error
	UpdatePassword(userID uint, hashedPassword string) error
}
