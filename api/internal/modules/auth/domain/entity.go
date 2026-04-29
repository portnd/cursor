package domain

import (
	"strings"
	"time"

	"github.com/lib/pq"
)

// Team represents a squad/team that groups users and projects for data isolation
type Team struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Name            string    `json:"name" gorm:"type:varchar(100);uniqueIndex;not null"`
	CapitalBalance  float64   `json:"capital_balance" gorm:"column:capital_balance;type:decimal(15,2);default:0"`
	BonusPercentage float64   `json:"bonus_percentage" gorm:"type:decimal(5,2);default:0"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Users           []User    `json:"users,omitempty" gorm:"foreignKey:TeamID"`
}

func (Team) TableName() string { return "teams" }

// TeamTransactionType defines the type of a capital transaction
type TeamTransactionType string

const (
	TxInjection   TeamTransactionType = "INJECTION"
	TxBurn        TeamTransactionType = "BURN"
	TxBonusPayout TeamTransactionType = "BONUS_PAYOUT"
	TxAdjustment  TeamTransactionType = "ADJUSTMENT"
)

// TeamTransaction records every capital movement for a team
type TeamTransaction struct {
	ID        int64               `json:"id" gorm:"primaryKey;autoIncrement"`
	TeamID    uint                `json:"team_id" gorm:"not null;index"`
	Type      TeamTransactionType `json:"type" gorm:"type:varchar(20);not null"`
	Amount    float64             `json:"amount" gorm:"type:decimal(15,2);not null"`
	Reference string              `json:"reference" gorm:"type:text;default:''"`
	CreatedAt time.Time           `json:"created_at" gorm:"autoCreateTime"`
}

func (TeamTransaction) TableName() string { return "team_transactions" }

// TeamMonthlyCostResponse is the response DTO for CalculateTeamMonthlyCost
type TeamMonthlyCostResponse struct {
	TeamID           uint    `json:"team_id"`
	MemberCost       float64 `json:"member_cost"`     // Σ salary+SS for all team members
	SharedOverhead   float64 `json:"shared_overhead"` // (ExecExp + CompanyExp + mgr/support salaries) / numTeams
	TotalMonthlyCost float64 `json:"total_monthly_cost"`
	CapitalBalance   float64 `json:"capital_balance"`
	BonusPercentage  float64 `json:"bonus_percentage"`
	RunwayMonths     float64 `json:"runway_months"` // CapitalBalance / TotalMonthlyCost
}

// InjectCapitalRequest is the DTO for injecting capital into a team
type InjectCapitalRequest struct {
	Amount          float64 `json:"amount" binding:"required,gt=0"`
	BonusPercentage float64 `json:"bonus_percentage" binding:"min=0,max=100"`
	Note            string  `json:"note"`
}

// EditCapitalRequest is the DTO for directly setting the capital balance
type EditCapitalRequest struct {
	NewBalance      float64  `json:"new_balance" binding:"gte=0"`
	BonusPercentage *float64 `json:"bonus_percentage" binding:"omitempty,gte=0,lte=100"`
	Note            string   `json:"note"`
}

// CloseCycleResponse is the response DTO for CloseCycleAndPayout
type CloseCycleResponse struct {
	TeamID          uint    `json:"team_id"`
	BalanceBefore   float64 `json:"balance_before"`
	BonusPercentage float64 `json:"bonus_percentage"`
	BonusAmount     float64 `json:"bonus_amount"`
	BalanceAfter    float64 `json:"balance_after"`
}

// User represents the user entity in the system
// This is the core domain model for authentication and Sentinel role management
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"` // Unique email for login
	Password  string    `json:"-" gorm:"not null"`                 // Password hash (never expose in JSON)
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Sentinel enhancements
	Role        string         `json:"role" gorm:"type:varchar(20);not null;default:'ENGINEER'"` // CEO, MANAGER, PRODUCT_OWNER, ENGINEER, CHIEF_ENGINEER, SUPPORT
	HealthScore float64        `json:"health_score" gorm:"type:decimal(5,2);default:100.00"`     // Performance tracking (0-100)
	TechStack   pq.StringArray `json:"tech_stack" gorm:"type:text[]"`                            // Array of technologies
	FirstName   string         `json:"first_name" gorm:"type:varchar(100);default:''"`          // Optional given name
	LastName    string         `json:"last_name" gorm:"type:varchar(100);default:''"`           // Optional family name
	DisplayName string         `json:"display_name" gorm:"type:varchar(100)"`                    // Optional display name (enterprise profile)

	// Squad Model
	TeamID *uint `json:"team_id" gorm:"index"` // nullable FK → teams
	Team   *Team `json:"team,omitempty" gorm:"foreignKey:TeamID"`

	// Performance Reset: CEO can reset a dev's Rework Rate; rejections before this timestamp are excluded
	ReworkResetAt *time.Time `json:"rework_reset_at,omitempty" gorm:"column:rework_reset_at"`

	// Avatar: stores a data-URL (base64) or external URL for the user's profile picture
	AvatarURL string `json:"avatar_url" gorm:"type:text;default:''"`

	// ThemePreference: "dark" or "light" — persisted per user account
	ThemePreference string `json:"theme_preference" gorm:"type:varchar(10);not null;default:'dark'"`
}

// UserRole constants for type safety
const (
	RoleCEO          = "CEO"
	RoleManager      = "MANAGER"
	RoleProductOwner = "PRODUCT_OWNER"
	// RolePMLegacy is a legacy DB value treated like PRODUCT_OWNER (JWT middleware normalizes PM → PRODUCT_OWNER for requests).
	RolePMLegacy      = "PM"
	RoleEngineer      = "ENGINEER"
	RoleChiefEngineer = "CHIEF_ENGINEER"
	RoleSupport       = "SUPPORT"
)

// IsProductOwnerAssignableRole reports whether the user may be saved as a project PM/Product Owner when squads are disabled.
// MANAGER is allowed for backward compatibility with existing project selection flows.
func IsProductOwnerAssignableRole(role string) bool {
	switch strings.ToUpper(strings.TrimSpace(role)) {
	case RoleProductOwner, RolePMLegacy, RoleManager:
		return true
	default:
		return false
	}
}

// IsEngineerRole reports whether the role should receive engineer-equivalent permissions and KPIs.
func IsEngineerRole(role string) bool {
	switch strings.ToUpper(strings.TrimSpace(role)) {
	case RoleEngineer, RoleChiefEngineer:
		return true
	default:
		return false
	}
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

// UpdateUserAdminRequest is the DTO for CEO editing a user's profile fields
type UpdateUserAdminRequest struct {
	FirstName   *string `json:"first_name" binding:"omitempty,max=100"`
	LastName    *string `json:"last_name" binding:"omitempty,max=100"`
	DisplayName *string `json:"display_name" binding:"omitempty,max=100"`
	Email       *string `json:"email" binding:"omitempty,email,max=255"`
}

// ChangeRoleRequest is the DTO for changing user roles (CEO only)
type ChangeRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=CEO MANAGER PRODUCT_OWNER ENGINEER CHIEF_ENGINEER SUPPORT"`
}

// CreateUserRequest is the DTO for admin creating a single user (CEO only)
type CreateUserRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=8"`
	Role       string `json:"role" binding:"required,oneof=CEO MANAGER PRODUCT_OWNER ENGINEER CHIEF_ENGINEER SUPPORT"`
	FirstName  string `json:"first_name" binding:"omitempty,max=100"`
	LastName   string `json:"last_name" binding:"omitempty,max=100"`
	DisplayName string `json:"display_name" binding:"omitempty,max=100"`
}

// ImportUserItem is one row in a bulk user import (CEO only)
// Password is optional; if empty, a random temporary password is generated and returned.
type ImportUserItem struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password"`                                                                                 // optional; min 8 if provided
	Role        string `json:"role" binding:"omitempty,oneof=CEO MANAGER PRODUCT_OWNER ENGINEER CHIEF_ENGINEER SUPPORT"` // default ENGINEER
	FirstName   string `json:"first_name" binding:"omitempty,max=100"`
	LastName    string `json:"last_name" binding:"omitempty,max=100"`
	DisplayName string `json:"display_name" binding:"omitempty,max=100"`
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
	FirstName   *string  `json:"first_name"`
	LastName    *string  `json:"last_name"`
	TechStack   []string `json:"tech_stack"`
}

// ChangePasswordRequest is the DTO for changing own password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// UpdateAvatarRequest is the DTO for updating a user's avatar
type UpdateAvatarRequest struct {
	// AvatarDataURL must be a valid data-URL (e.g. "data:image/png;base64,...")
	// or an empty string to remove the avatar.
	AvatarDataURL string `json:"avatar_data_url" binding:"required"`
}

// UpdateThemeRequest is the DTO for updating a user's theme preference
type UpdateThemeRequest struct {
	Theme string `json:"theme" binding:"required,oneof=dark light"`
}

// CreateTeamRequest is the DTO for creating a new team (CEO only)
type CreateTeamRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

// UpdateTeamRequest is the DTO for updating a team name (CEO only)
type UpdateTeamRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

// AssignUserToTeamRequest is the DTO for assigning a user to a team (CEO only)
type AssignUserToTeamRequest struct {
	TeamID *uint `json:"team_id"` // null = remove from team
}

// TeamsFeatureSettingResponse is the response for GET /auth/settings/teams-feature
type TeamsFeatureSettingResponse struct {
	Enabled bool `json:"enabled"`
}

// SetTeamsFeatureRequest is the request body for PUT /auth/settings/teams-feature (CEO only)
type SetTeamsFeatureRequest struct {
	Enabled bool `json:"enabled"`
}

// Usecase defines the authentication business logic interface
// This follows the Dependency Inversion Principle (Hexagonal Architecture)
type Usecase interface {
	Register(req *RegisterRequest) (*AuthResponse, error)
	Login(req *LoginRequest) (*AuthResponse, error)
	// Profile (any authenticated user)
	GetProfile(userID uint) (*User, error)
	UpdateProfile(userID uint, req *UpdateProfileRequest) (*User, error)
	UpdateAvatar(userID uint, avatarDataURL string) (*User, error)
	UpdateThemePreference(userID uint, theme string) (*User, error)
	ChangePassword(userID uint, currentPassword, newPassword string) error
	// User Management (CEO only)
	GetTeamMembers(requestingUserID uint) ([]User, error)
	ChangeUserRole(requestingUserID uint, targetUserID uint, newRole string) error
	// Admin user creation (CEO only)
	CreateUserAsAdmin(requestingUserID uint, req *CreateUserRequest) (*User, error)
	ImportUsers(requestingUserID uint, req *ImportUsersRequest) (*ImportUsersResponse, error)
	// Update user profile (CEO only)
	UpdateUserAdmin(requestingUserID uint, targetUserID uint, req *UpdateUserAdminRequest) (*User, error)
	// Delete user (CEO only; cannot delete self)
	DeleteUser(requestingUserID uint, targetUserID uint) error
	// Reset user password (CEO only). Returns the new temporary password.
	ResetUserPassword(requestingUserID uint, targetUserID uint) (tempPassword string, err error)
	// Squad / Team management (CEO only)
	GetAllTeams() ([]Team, error)
	GetTeamsFeatureEnabled() (bool, error)
	SetTeamsFeatureEnabled(requestingUserID uint, enabled bool) error
	CreateTeam(name string) (*Team, error)
	UpdateTeam(teamID uint, name string) (*Team, error)
	DeleteTeam(teamID uint) error
	AssignUserToTeam(requestingUserID uint, targetUserID uint, teamID *uint) error
}

// Repository defines the data access interface for user operations
// This is the PORT in Hexagonal Architecture
type Repository interface {
	CreateUser(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
	// Profile
	UpdateProfile(userID uint, displayName, firstName, lastName *string, techStack []string) error
	UpdateAvatar(userID uint, avatarURL string) error
	UpdateThemePreference(userID uint, theme string) error
	// User Management
	GetAllUsers() ([]User, error)
	UpdateUserRole(userID uint, newRole string) error
	UpdateUserAdmin(userID uint, firstName, lastName, displayName, email *string) error
	DeleteUser(userID uint) error
	UpdatePassword(userID uint, hashedPassword string) error
	ResetReworkRate(userID uint) error // CEO: set rework_reset_at = NOW() to clear rework history
	// Team / Squad management
	CreateTeam(team *Team) error
	GetAllTeams() ([]Team, error)
	GetTeamByID(id uint) (*Team, error)
	UpdateTeamName(teamID uint, name string) error
	DeleteTeam(teamID uint) error
	AssignUserToTeam(userID uint, teamID *uint) error
	// Team Finance
	UpdateTeamCapital(teamID uint, newBalance float64, bonusPct *float64) error
	CreateTeamTransaction(tx *TeamTransaction) error
	GetTeamTransactions(teamID uint) ([]TeamTransaction, error)
	// App settings (feature flags)
	GetAppSetting(key string) (string, error)
	SetAppSetting(key, value string) error
}

// TeamFinanceUsecase defines the business logic for the Internal VC model
type TeamFinanceUsecase interface {
	CalculateTeamMonthlyCost(teamID uint) (*TeamMonthlyCostResponse, error)
	InjectCapital(teamID uint, req *InjectCapitalRequest) (*Team, error)
	EditCapital(teamID uint, req *EditCapitalRequest) (*Team, error)
	CloseCycleAndPayout(teamID uint) (*CloseCycleResponse, error)
}
