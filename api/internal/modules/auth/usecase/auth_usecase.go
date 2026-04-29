package usecase

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	"golang.org/x/crypto/bcrypt"
)

// MaxImportUsers is the maximum number of users allowed in a single bulk import
const MaxImportUsers = 500

// authUsecase implements the business logic for authentication
// This is the CORE of Hexagonal Architecture - pure business logic
type authUsecase struct {
	repo      domain.Repository
	jwtSecret []byte
}

// NewAuthUsecase creates a new authentication usecase instance
// Follows Dependency Injection pattern
func NewAuthUsecase(repo domain.Repository) domain.Usecase {
	// Load JWT secret from environment.
	// Keep legacy fallback for compatibility, but surface a loud warning so it is not missed in non-dev environments.
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default_jwt_secret_change_in_production"
		if env := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV"))); env == "production" || env == "prod" {
			log.Println("[SECURITY][CRITICAL] JWT_SECRET is not set in production; using insecure default secret")
		} else {
			log.Println("[SECURITY][WARN] JWT_SECRET is not set; using development fallback secret")
		}
	}

	return &authUsecase{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
	}
}

// Register handles user registration
// Business Rules:
// 1. Check if email already exists
// 2. Hash the password using bcrypt
// 3. Create user in database
// 4. Generate JWT token
func (u *authUsecase) Register(req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	// Check if user with this email already exists
	existingUser, _ := u.repo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("email already registered")
	}

	// Hash the password using bcrypt (cost factor: 12)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user entity
	user := &domain.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	// Save to database
	if err := u.repo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token with role
	token, err := u.generateJWT(user.ID, user.Email, user.Role, user.TeamID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Return response with token and user info
	return &domain.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

// Login handles user login
// Business Rules:
// 1. Find user by email
// 2. Verify password using bcrypt.CompareHashAndPassword
// 3. Generate JWT token if password is correct
func (u *authUsecase) Login(req *domain.LoginRequest) (*domain.AuthResponse, error) {
	// Find user by email
	user, err := u.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Compare password hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate JWT token with role
	token, err := u.generateJWT(user.ID, user.Email, user.Role, user.TeamID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Return response with token and user info
	return &domain.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

// generateJWT creates a JWT token with user claims
// Token includes: user_id, email, role, team_id, issued_at, expires_at (24 hours)
// Algorithm: HS256 (HMAC with SHA-256)
func (u *authUsecase) generateJWT(userID uint, email, role string, teamID *uint) (string, error) {
	// Define JWT claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"iat":     time.Now().Unix(),                     // Issued At
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Expires in 24 hours
	}
	if teamID != nil {
		claims["team_id"] = *teamID
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString(u.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// GetProfile returns the current user's profile by ID (any authenticated user)
func (u *authUsecase) GetProfile(userID uint) (*domain.User, error) {
	user, err := u.repo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// UpdateProfile updates display name and/or tech stack for the current user
func (u *authUsecase) UpdateProfile(userID uint, req *domain.UpdateProfileRequest) (*domain.User, error) {
	if _, err := u.repo.FindByID(userID); err != nil {
		return nil, fmt.Errorf("user not found")
	}
	var displayName *string
	if req.DisplayName != nil {
		trimmed := strings.TrimSpace(*req.DisplayName)
		displayName = &trimmed
	}
	var firstName *string
	if req.FirstName != nil {
		trimmed := strings.TrimSpace(*req.FirstName)
		firstName = &trimmed
	}
	var lastName *string
	if req.LastName != nil {
		trimmed := strings.TrimSpace(*req.LastName)
		lastName = &trimmed
	}
	if displayName == nil && firstName != nil && lastName != nil {
		joined := strings.TrimSpace(strings.TrimSpace(*firstName) + " " + strings.TrimSpace(*lastName))
		displayName = &joined
	}
	if err := u.repo.UpdateProfile(userID, displayName, firstName, lastName, req.TechStack); err != nil {
		return nil, err
	}
	return u.repo.FindByID(userID)
}

// UpdateAvatar updates the user's avatar (stored as a data-URL string, max ~2 MB encoded)
func (u *authUsecase) UpdateAvatar(userID uint, avatarDataURL string) (*domain.User, error) {
	if _, err := u.repo.FindByID(userID); err != nil {
		return nil, fmt.Errorf("user not found")
	}
	// Guard against absurdly large payloads (2 MB base64 ≈ ~2.7 MB string)
	const maxLen = 3 * 1024 * 1024
	if len(avatarDataURL) > maxLen {
		return nil, fmt.Errorf("avatar image too large (max 2 MB)")
	}
	if err := u.repo.UpdateAvatar(userID, avatarDataURL); err != nil {
		return nil, err
	}
	return u.repo.FindByID(userID)
}

// UpdateThemePreference persists the user's chosen UI theme ("dark" or "light")
func (u *authUsecase) UpdateThemePreference(userID uint, theme string) (*domain.User, error) {
	if _, err := u.repo.FindByID(userID); err != nil {
		return nil, fmt.Errorf("user not found")
	}
	if err := u.repo.UpdateThemePreference(userID, theme); err != nil {
		return nil, err
	}
	return u.repo.FindByID(userID)
}

// ChangePassword changes the current user's password after verifying current password
func (u *authUsecase) ChangePassword(userID uint, currentPassword, newPassword string) error {
	user, err := u.repo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return fmt.Errorf("current password is incorrect")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	return u.repo.UpdatePassword(userID, string(hashedPassword))
}

// GetTeamMembers retrieves all users (CEO, MANAGER, or Product Owner — for team management and task assignment).
func (u *authUsecase) GetTeamMembers(requestingUserID uint) ([]domain.User, error) {
	// Get requesting user to check their role
	requestingUser, err := u.repo.FindByID(requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: user not found")
	}

	// CEO, MANAGER, and Product Owner can list users (e.g. for task assignment)
	if requestingUser.Role != domain.RoleCEO && requestingUser.Role != domain.RoleManager && requestingUser.Role != domain.RoleProductOwner {
		return nil, fmt.Errorf("unauthorized: only CEO, MANAGER, or Product Owner can view team members")
	}

	// Fetch all users
	users, err := u.repo.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch team members: %w", err)
	}

	return users, nil
}

// ChangeUserRole changes a user's role (CEO only)
// Business Rules:
// 1. Only users with role 'CEO' can change roles
// 2. New role must be one of: CEO, MANAGER, PRODUCT_OWNER, ENGINEER, CHIEF_ENGINEER, SUPPORT
// 3. Cannot change own role (optional safeguard)
func (u *authUsecase) ChangeUserRole(requestingUserID uint, targetUserID uint, newRole string) error {
	// Get requesting user to check their role
	requestingUser, err := u.repo.FindByID(requestingUserID)
	if err != nil {
		return fmt.Errorf("unauthorized: user not found")
	}

	// Check if user is CEO
	if requestingUser.Role != domain.RoleCEO {
		return fmt.Errorf("unauthorized: only CEO can change user roles")
	}

	// Validate new role
	if newRole != domain.RoleCEO && newRole != domain.RoleManager && newRole != domain.RoleProductOwner && newRole != domain.RoleEngineer && newRole != domain.RoleChiefEngineer && newRole != domain.RoleSupport {
		return fmt.Errorf("invalid role: must be one of CEO, MANAGER, PRODUCT_OWNER, ENGINEER, CHIEF_ENGINEER, or SUPPORT")
	}

	// Optional: Prevent CEO from changing their own role
	if requestingUserID == targetUserID {
		return fmt.Errorf("cannot change your own role")
	}

	// Check if target user exists
	targetUser, err := u.repo.FindByID(targetUserID)
	if err != nil {
		return fmt.Errorf("target user not found")
	}

	// Update role
	if err := u.repo.UpdateUserRole(targetUser.ID, newRole); err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}

	return nil
}

// CreateUserAsAdmin creates a single user (CEO only).
// Business rules: caller must be CEO; email must be unique; password hashed with bcrypt.
func (u *authUsecase) CreateUserAsAdmin(requestingUserID uint, req *domain.CreateUserRequest) (*domain.User, error) {
	requestingUser, err := u.repo.FindByID(requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: user not found")
	}
	if requestingUser.Role != domain.RoleCEO {
		return nil, fmt.Errorf("unauthorized: only CEO can create users")
	}

	existingUser, _ := u.repo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Email:       strings.TrimSpace(req.Email),
		Password:    string(hashedPassword),
		Role:        req.Role,
		FirstName:   strings.TrimSpace(req.FirstName),
		LastName:    strings.TrimSpace(req.LastName),
		DisplayName: strings.TrimSpace(req.DisplayName),
	}
	if err := u.repo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

// ImportUsers creates multiple users in one request (CEO only).
// Optional password per row; if empty, a random temporary password is generated and returned.
// Duplicate emails are skipped; invalid rows are reported in results.
func (u *authUsecase) ImportUsers(requestingUserID uint, req *domain.ImportUsersRequest) (*domain.ImportUsersResponse, error) {
	requestingUser, err := u.repo.FindByID(requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: user not found")
	}
	if requestingUser.Role != domain.RoleCEO {
		return nil, fmt.Errorf("unauthorized: only CEO can import users")
	}

	if len(req.Users) == 0 {
		return &domain.ImportUsersResponse{Total: 0, Results: []domain.ImportUserResult{}}, nil
	}
	if len(req.Users) > MaxImportUsers {
		return nil, fmt.Errorf("too many users: maximum %d per import", MaxImportUsers)
	}

	results := make([]domain.ImportUserResult, 0, len(req.Users))
	var created, skipped, errCount int

	for _, item := range req.Users {
		email := strings.TrimSpace(strings.ToLower(item.Email))
		role := item.Role
		if role == "" {
			role = domain.RoleEngineer
		}

		// Validate password if provided
		password := item.Password
		if password != "" && len(password) < 8 {
			errCount++
			results = append(results, domain.ImportUserResult{
				Email:   email,
				Status:  "error",
				Message: "password must be at least 8 characters",
			})
			continue
		}
		if password == "" {
			var errGen error
			password, errGen = generateTempPassword()
			if errGen != nil {
				errCount++
				results = append(results, domain.ImportUserResult{
					Email:   email,
					Status:  "error",
					Message: "failed to generate temporary password",
				})
				continue
			}
		}

		existingUser, _ := u.repo.FindByEmail(email)
		if existingUser != nil {
			skipped++
			results = append(results, domain.ImportUserResult{
				Email:   email,
				Status:  "skipped",
				Message: "email already registered",
			})
			continue
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			errCount++
			results = append(results, domain.ImportUserResult{
				Email:   email,
				Status:  "error",
				Message: "failed to hash password",
			})
			continue
		}

		user := &domain.User{
			Email:       email,
			Password:    string(hashedPassword),
			Role:        role,
			FirstName:   strings.TrimSpace(item.FirstName),
			LastName:    strings.TrimSpace(item.LastName),
			DisplayName: strings.TrimSpace(item.DisplayName),
		}
		if err := u.repo.CreateUser(user); err != nil {
			errCount++
			results = append(results, domain.ImportUserResult{
				Email:   email,
				Status:  "error",
				Message: err.Error(),
			})
			continue
		}

		created++
		res := domain.ImportUserResult{
			Email:   email,
			Status:  "created",
			User:    user,
		}
		if item.Password == "" {
			res.TempPassword = password
		}
		results = append(results, res)
	}

	return &domain.ImportUsersResponse{
		Total:   len(req.Users),
		Created: created,
		Skipped: skipped,
		Errors:  errCount,
		Results: results,
	}, nil
}

// UpdateUserAdmin allows CEO to edit a user's profile fields (first_name, last_name, display_name, email)
func (u *authUsecase) UpdateUserAdmin(requestingUserID uint, targetUserID uint, req *domain.UpdateUserAdminRequest) (*domain.User, error) {
	requestingUser, err := u.repo.FindByID(requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: user not found")
	}
	if requestingUser.Role != domain.RoleCEO {
		return nil, fmt.Errorf("unauthorized: only CEO can edit users")
	}
	targetUser, err := u.repo.FindByID(targetUserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// If email is changing, check for duplicates
	if req.Email != nil {
		trimmed := strings.TrimSpace(*req.Email)
		if trimmed != targetUser.Email {
			existing, _ := u.repo.FindByEmail(trimmed)
			if existing != nil {
				return nil, fmt.Errorf("email already registered")
			}
		}
	}

	// Trim string fields
	var firstName, lastName, displayName, email *string
	if req.FirstName != nil {
		v := strings.TrimSpace(*req.FirstName)
		firstName = &v
	}
	if req.LastName != nil {
		v := strings.TrimSpace(*req.LastName)
		lastName = &v
	}
	if req.DisplayName != nil {
		v := strings.TrimSpace(*req.DisplayName)
		displayName = &v
	}
	if req.Email != nil {
		v := strings.TrimSpace(*req.Email)
		email = &v
	}

	if err := u.repo.UpdateUserAdmin(targetUserID, firstName, lastName, displayName, email); err != nil {
		return nil, err
	}

	return u.repo.FindByID(targetUserID)
}

// DeleteUser removes a user (CEO only). Cannot delete yourself.
func (u *authUsecase) DeleteUser(requestingUserID uint, targetUserID uint) error {
	requestingUser, err := u.repo.FindByID(requestingUserID)
	if err != nil {
		return fmt.Errorf("unauthorized: user not found")
	}
	if requestingUser.Role != domain.RoleCEO {
		return fmt.Errorf("unauthorized: only CEO can delete users")
	}
	if requestingUserID == targetUserID {
		return fmt.Errorf("cannot delete your own account")
	}
	if _, err := u.repo.FindByID(targetUserID); err != nil {
		return fmt.Errorf("user not found")
	}
	return u.repo.DeleteUser(targetUserID)
}

// ResetUserPassword generates a temporary password for a user (CEO only) and returns it.
func (u *authUsecase) ResetUserPassword(requestingUserID uint, targetUserID uint) (string, error) {
	requestingUser, err := u.repo.FindByID(requestingUserID)
	if err != nil {
		return "", fmt.Errorf("unauthorized: user not found")
	}
	if requestingUser.Role != domain.RoleCEO {
		return "", fmt.Errorf("unauthorized: only CEO can reset passwords")
	}
	if _, err := u.repo.FindByID(targetUserID); err != nil {
		return "", fmt.Errorf("user not found")
	}
	tempPassword, err := generateTempPassword()
	if err != nil {
		return "", fmt.Errorf("failed to generate password: %w", err)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempPassword), 12)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	if err := u.repo.UpdatePassword(targetUserID, string(hashedPassword)); err != nil {
		return "", err
	}
	return tempPassword, nil
}

// generateTempPassword returns a random 12-character password (alphanumeric)
func generateTempPassword() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 12
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}
	return string(b), nil
}

const teamsFeatureEnabledKey = "teams_feature_enabled"

// GetAllTeams returns all teams; returns empty list when teams feature is disabled
func (u *authUsecase) GetAllTeams() ([]domain.Team, error) {
	enabled, err := u.repo.GetAppSetting(teamsFeatureEnabledKey)
	if err != nil {
		return nil, err
	}
	if enabled != "true" {
		return []domain.Team{}, nil
	}
	return u.repo.GetAllTeams()
}

// GetTeamsFeatureEnabled returns whether the teams (squads) feature is enabled
func (u *authUsecase) GetTeamsFeatureEnabled() (bool, error) {
	val, err := u.repo.GetAppSetting(teamsFeatureEnabledKey)
	if err != nil {
		return true, err
	}
	return val != "false", nil
}

// SetTeamsFeatureEnabled enables or disables the teams feature (CEO only)
func (u *authUsecase) SetTeamsFeatureEnabled(requestingUserID uint, enabled bool) error {
	requester, err := u.repo.FindByID(requestingUserID)
	if err != nil {
		return fmt.Errorf("unauthorized: user not found")
	}
	if requester.Role != domain.RoleCEO && requester.Role != domain.RoleManager {
		return fmt.Errorf("unauthorized: only CEO or MANAGER can change teams feature")
	}
	value := "true"
	if !enabled {
		value = "false"
	}
	return u.repo.SetAppSetting(teamsFeatureEnabledKey, value)
}

// CreateTeam creates a new squad/team (CEO only)
func (u *authUsecase) CreateTeam(name string) (*domain.Team, error) {
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("team name cannot be empty")
	}
	team := &domain.Team{Name: strings.TrimSpace(name)}
	if err := u.repo.CreateTeam(team); err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}
	return team, nil
}

// UpdateTeam updates a team's name by ID (CEO only)
func (u *authUsecase) UpdateTeam(teamID uint, name string) (*domain.Team, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("team name cannot be empty")
	}
	if _, err := u.repo.GetTeamByID(teamID); err != nil {
		return nil, fmt.Errorf("team not found")
	}
	if err := u.repo.UpdateTeamName(teamID, name); err != nil {
		return nil, fmt.Errorf("failed to update team: %w", err)
	}
	return u.repo.GetTeamByID(teamID)
}

// DeleteTeam removes a team by ID (CEO only)
func (u *authUsecase) DeleteTeam(teamID uint) error {
	if _, err := u.repo.GetTeamByID(teamID); err != nil {
		return fmt.Errorf("team not found")
	}
	return u.repo.DeleteTeam(teamID)
}

// AssignUserToTeam assigns (or unassigns when teamID is nil) a user to a team (CEO only)
func (u *authUsecase) AssignUserToTeam(requestingUserID uint, targetUserID uint, teamID *uint) error {
	requester, err := u.repo.FindByID(requestingUserID)
	if err != nil {
		return fmt.Errorf("unauthorized: user not found")
	}
	if requester.Role != domain.RoleCEO {
		return fmt.Errorf("unauthorized: only CEO can assign users to teams")
	}
	if _, err := u.repo.FindByID(targetUserID); err != nil {
		return fmt.Errorf("target user not found")
	}
	if teamID != nil {
		if _, err := u.repo.GetTeamByID(*teamID); err != nil {
			return fmt.Errorf("team not found")
		}
	}
	return u.repo.AssignUserToTeam(targetUserID, teamID)
}
