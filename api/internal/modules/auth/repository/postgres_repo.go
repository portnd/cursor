package repository

import (
	"fmt"

	"github.com/lib/pq"
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

// UpdateProfile updates display name and/or tech stack for the given user
func (r *postgresRepository) UpdateProfile(userID uint, displayName *string, techStack []string) error {
	updates := make(map[string]interface{})
	if displayName != nil {
		updates["display_name"] = *displayName
	}
	if techStack != nil {
		updates["tech_stack"] = pq.Array(techStack)
	}
	if len(updates) == 0 {
		return nil
	}
	result := r.db.Model(&domain.User{}).Where("id = ?", userID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update profile: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
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

// DeleteUser removes a user by ID (hard delete)
func (r *postgresRepository) DeleteUser(userID uint) error {
	result := r.db.Delete(&domain.User{}, userID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// UpdatePassword updates a user's password hash
func (r *postgresRepository) UpdatePassword(userID uint, hashedPassword string) error {
	result := r.db.Model(&domain.User{}).Where("id = ?", userID).Update("password", hashedPassword)
	if result.Error != nil {
		return fmt.Errorf("failed to update password: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *postgresRepository) ResetReworkRate(userID uint) error {
	result := r.db.Exec("UPDATE users SET rework_reset_at = NOW(), updated_at = NOW() WHERE id = ?", userID)
	if result.Error != nil {
		return fmt.Errorf("failed to reset rework rate: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// --- Team / Squad Management ---

// CreateTeam inserts a new team
func (r *postgresRepository) CreateTeam(team *domain.Team) error {
	if err := r.db.Create(team).Error; err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}
	return nil
}

// GetAllTeams returns all teams with their members preloaded
func (r *postgresRepository) GetAllTeams() ([]domain.Team, error) {
	var teams []domain.Team
	if err := r.db.Preload("Users").Order("created_at DESC").Find(&teams).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch teams: %w", err)
	}
	return teams, nil
}

// GetTeamByID returns a single team by ID
func (r *postgresRepository) GetTeamByID(id uint) (*domain.Team, error) {
	var team domain.Team
	if err := r.db.Preload("Users").First(&team, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("team not found")
		}
		return nil, fmt.Errorf("failed to find team: %w", err)
	}
	return &team, nil
}

// UpdateTeamName updates a team's name by ID
func (r *postgresRepository) UpdateTeamName(teamID uint, name string) error {
	result := r.db.Model(&domain.Team{}).Where("id = ?", teamID).Update("name", name)
	if result.Error != nil {
		return fmt.Errorf("failed to update team name: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("team not found")
	}
	return nil
}

// DeleteTeam removes a team by ID (users' team_id is set to NULL by DB cascade)
func (r *postgresRepository) DeleteTeam(teamID uint) error {
	// Unassign all users from this team first
	r.db.Model(&domain.User{}).Where("team_id = ?", teamID).Update("team_id", nil)
	result := r.db.Delete(&domain.Team{}, teamID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete team: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("team not found")
	}
	return nil
}

// AssignUserToTeam sets or clears a user's team_id
func (r *postgresRepository) AssignUserToTeam(userID uint, teamID *uint) error {
	result := r.db.Model(&domain.User{}).Where("id = ?", userID).Update("team_id", teamID)
	if result.Error != nil {
		return fmt.Errorf("failed to assign user to team: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// --- Team Finance ---

// UpdateTeamCapital updates capital_balance and optionally bonus_percentage atomically
func (r *postgresRepository) UpdateTeamCapital(teamID uint, newBalance float64, bonusPct *float64) error {
	updates := map[string]interface{}{
		"capital_balance": newBalance,
	}
	if bonusPct != nil {
		updates["bonus_percentage"] = *bonusPct
	}
	result := r.db.Model(&domain.Team{}).Where("id = ?", teamID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update team capital: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("team not found")
	}
	return nil
}

// CreateTeamTransaction inserts a new capital transaction record
func (r *postgresRepository) CreateTeamTransaction(tx *domain.TeamTransaction) error {
	if err := r.db.Create(tx).Error; err != nil {
		return fmt.Errorf("failed to create team transaction: %w", err)
	}
	return nil
}

// GetTeamTransactions returns all transactions for a team ordered by newest first
func (r *postgresRepository) GetTeamTransactions(teamID uint) ([]domain.TeamTransaction, error) {
	var txs []domain.TeamTransaction
	if err := r.db.Where("team_id = ?", teamID).Order("created_at DESC").Find(&txs).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch team transactions: %w", err)
	}
	return txs, nil
}
