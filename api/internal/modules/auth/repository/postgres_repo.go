package repository

import (
	"fmt"

	"github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	"gorm.io/gorm"
)

// postgresRepository is the ADAPTER that implements the Repository PORT
// This is the concrete implementation of the data access layer
type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new instance of the PostgreSQL repository
// Constructor function following Dependency Injection pattern
func NewPostgresRepository(db *gorm.DB) domain.Repository {
	return &postgresRepository{
		db: db,
	}
}

// CreateUser inserts a new user into the database
// The password should already be hashed by the usecase layer
func (r *postgresRepository) CreateUser(user *domain.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// FindByEmail retrieves a user by their email address
// Returns gorm.ErrRecordNotFound if user doesn't exist
func (r *postgresRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	return &user, nil
}

// FindByID retrieves a user by their ID
// Returns gorm.ErrRecordNotFound if user doesn't exist
func (r *postgresRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}
	return &user, nil
}

// GetAllUsers retrieves all users from the database
// Password hash is automatically hidden by the `json:"-"` tag in the User struct
func (r *postgresRepository) GetAllUsers() ([]domain.User, error) {
	var users []domain.User
	if err := r.db.Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	return users, nil
}

// UpdateUserRole updates a user's role
// Validates that the role is one of the allowed values
func (r *postgresRepository) UpdateUserRole(userID uint, newRole string) error {
	// Update only the role field
	result := r.db.Model(&domain.User{}).Where("id = ?", userID).Update("role", newRole)
	
	if result.Error != nil {
		return fmt.Errorf("failed to update user role: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	
	return nil
}
