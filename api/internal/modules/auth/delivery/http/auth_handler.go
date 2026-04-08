package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
)

// AuthHandler handles HTTP requests for authentication
// This is the ADAPTER layer in Hexagonal Architecture
type AuthHandler struct {
	usecase domain.Usecase
}

// NewAuthHandler creates a new authentication handler
// Constructor following Dependency Injection pattern
func NewAuthHandler(usecase domain.Usecase) *AuthHandler {
	return &AuthHandler{
		usecase: usecase,
	}
}

// Register handles POST /auth/register
// Request Body: { "email": "user@example.com", "password": "password123", "confirm_password": "password123" }
// Response: { "token": "jwt_token", "user": { ... } }
func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest

	// Bind and validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Call usecase
	response, err := h.usecase.Register(&req)
	if err != nil {
		// Check for specific error types
		if err.Error() == "email already registered" {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Conflict",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Registration failed",
			"message": err.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data":    response,
	})
}

// Login handles POST /auth/login
// Request Body: { "email": "user@example.com", "password": "password123" }
// Response: { "token": "jwt_token", "user": { ... } }
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest

	// Bind and validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Call usecase
	response, err := h.usecase.Login(&req)
	if err != nil {
		// Always return generic error for security (don't reveal if email exists)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Authentication failed",
			"message": "Invalid email or password",
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    response,
	})
}

// GetMe handles GET /auth/me
// Returns the current authenticated user's profile
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Authentication required",
		})
		return
	}
	user, err := h.usecase.GetProfile(userID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch profile",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile retrieved successfully",
		"data":    user,
	})
}

// UpdateProfile handles PATCH /auth/me
// Request Body: { "display_name": "Optional", "tech_stack": ["Go", "Vue"] }
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Authentication required",
		})
		return
	}
	var req domain.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}
	user, err := h.usecase.UpdateProfile(userID, &req)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update profile",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"data":    user,
	})
}

// UploadAvatar handles POST /auth/me/avatar
// Request Body: { "avatar_data_url": "data:image/png;base64,..." }
// Stores the data-URL string directly in the database (max ~2 MB).
func (h *AuthHandler) UploadAvatar(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Authentication required",
		})
		return
	}
	var req domain.UpdateAvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}
	user, err := h.usecase.UpdateAvatar(userID, req.AvatarDataURL)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		} else if err.Error() == "avatar image too large (max 2 MB)" {
			status = http.StatusRequestEntityTooLarge
		}
		c.JSON(status, gin.H{
			"error":   "Failed to update avatar",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Avatar updated successfully",
		"data":    user,
	})
}

// UpdateTheme handles PATCH /auth/me/theme
// Request Body: { "theme": "dark" | "light" }
func (h *AuthHandler) UpdateTheme(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Authentication required",
		})
		return
	}
	var req domain.UpdateThemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}
	user, err := h.usecase.UpdateThemePreference(userID, req.Theme)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"error":   "Failed to update theme preference",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Theme preference updated",
		"data":    user,
	})
}

// ChangePassword handles PATCH /auth/me/password
// Request Body: { "current_password": "...", "new_password": "..." }
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Authentication required",
		})
		return
	}
	var req domain.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}
	if err := h.usecase.ChangePassword(userID, req.CurrentPassword, req.NewPassword); err != nil {
		if err.Error() == "current password is incorrect" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to change password",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}

// GetTeam handles GET /auth/users
// Retrieves all users (CEO only)
// Requires JWT authentication with userId in context
func (h *AuthHandler) GetTeam(c *gin.Context) {
	// Extract user ID from JWT context (set by auth middleware)
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Authentication required",
		})
		return
	}

	// Call usecase
	users, err := h.usecase.GetTeamMembers(userID)
	if err != nil {
		// Check for authorization errors
		if err.Error() == "unauthorized: only CEO can view team members" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch team members",
			"message": err.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Team members retrieved successfully",
		"data":    users,
	})
}

// ChangeRole handles PATCH /auth/users/:id/role
// Changes a user's role (CEO only)
// Request Body: { "role": "PRODUCT_OWNER" }
func (h *AuthHandler) ChangeRole(c *gin.Context) {
	// Extract user ID from JWT context
	requestingUserID := getUserIDFromContext(c)
	if requestingUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Authentication required",
		})
		return
	}

	// Get target user ID from URL parameter
	targetUserID := c.Param("id")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "User ID is required",
		})
		return
	}

	// Convert target user ID to uint
	var targetID uint
	if _, err := fmt.Sscanf(targetUserID, "%d", &targetID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid user ID format",
		})
		return
	}

	// Bind and validate request body
	var req domain.ChangeRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Call usecase
	if err := h.usecase.ChangeUserRole(requestingUserID, targetID, req.Role); err != nil {
		// Check for specific error types
		if err.Error() == "unauthorized: only CEO can change user roles" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": err.Error(),
			})
			return
		}

		if err.Error() == "target user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}

		if err.Error() == "cannot change your own role" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to change user role",
			"message": err.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "User role updated successfully",
	})
}

// DeleteUser handles DELETE /auth/users/:id (CEO only). Cannot delete yourself.
func (h *AuthHandler) DeleteUser(c *gin.Context) {
	requestingUserID := getUserIDFromContext(c)
	if requestingUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Authentication required",
		})
		return
	}

	targetUserID := c.Param("id")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "User ID is required",
		})
		return
	}

	var targetID uint
	if _, err := fmt.Sscanf(targetUserID, "%d", &targetID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid user ID format",
		})
		return
	}

	if err := h.usecase.DeleteUser(requestingUserID, targetID); err != nil {
		if err.Error() == "unauthorized: only CEO can delete users" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": err.Error(),
			})
			return
		}
		if err.Error() == "cannot delete your own account" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete user",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// ResetPassword handles PATCH /auth/users/:id/password (CEO only)
// No body required. System generates a random temporary password and returns it.
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	requestingUserID := getUserIDFromContext(c)
	if requestingUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Authentication required",
		})
		return
	}

	targetUserID := c.Param("id")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "User ID is required",
		})
		return
	}

	var targetID uint
	if _, err := fmt.Sscanf(targetUserID, "%d", &targetID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid user ID format",
		})
		return
	}

	tempPassword, err := h.usecase.ResetUserPassword(requestingUserID, targetID)
	if err != nil {
		if err.Error() == "unauthorized: only CEO can reset passwords" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": err.Error(),
			})
			return
		}
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to reset password",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Password reset successfully",
		"temp_password": tempPassword,
	})
}

// CreateUser handles POST /auth/users (CEO only)
// Request Body: { "email": "user@example.com", "password": "password123", "role": "ENGINEER" }
func (h *AuthHandler) CreateUser(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Authentication required",
		})
		return
	}

	var req domain.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	user, err := h.usecase.CreateUserAsAdmin(userID, &req)
	if err != nil {
		if err.Error() == "unauthorized: only CEO can create users" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": err.Error(),
			})
			return
		}
		if err.Error() == "email already registered" {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Conflict",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"data":    user,
	})
}

// ImportUsers handles POST /auth/users/import (CEO only)
// Request Body: { "users": [ { "email": "...", "password": "optional", "role": "ENGINEER" } ] }
// Max 500 users per request. If password is omitted, a temporary password is returned per created user.
func (h *AuthHandler) ImportUsers(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Authentication required",
		})
		return
	}

	var req domain.ImportUsersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	resp, err := h.usecase.ImportUsers(userID, &req)
	if err != nil {
		if err.Error() == "unauthorized: only CEO can import users" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": err.Error(),
			})
			return
		}
		if strings.HasPrefix(err.Error(), "too many users:") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Import failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Import completed",
		"data":    resp,
	})
}

// getUserIDFromContext extracts the user ID from the Gin context
// This assumes the auth middleware sets "userId" or "user_id" in the context
func getUserIDFromContext(c *gin.Context) uint {
	// Try common key patterns
	if val, exists := c.Get("userId"); exists {
		if id, ok := val.(float64); ok {
			return uint(id)
		}
		if id, ok := val.(uint); ok {
			return id
		}
		if id, ok := val.(int); ok {
			return uint(id)
		}
	}

	if val, exists := c.Get("user_id"); exists {
		if id, ok := val.(float64); ok {
			return uint(id)
		}
		if id, ok := val.(uint); ok {
			return id
		}
		if id, ok := val.(int); ok {
			return uint(id)
		}
	}

	return 0
}

// --- Team / Squad Management Handlers ---

// GetTeams handles GET /auth/teams (CEO only)
func (h *AuthHandler) GetTeams(c *gin.Context) {
	teams, err := h.usecase.GetAllTeams()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teams", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Teams retrieved successfully", "data": teams})
}

// CreateTeam handles POST /auth/teams (CEO only)
func (h *AuthHandler) CreateTeam(c *gin.Context) {
	var req domain.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	team, err := h.usecase.CreateTeam(req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create team", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Team created successfully", "data": team})
}

// UpdateTeam handles PATCH /auth/teams/:id (CEO only)
func (h *AuthHandler) UpdateTeam(c *gin.Context) {
	idStr := c.Param("id")
	var teamID uint
	if _, err := fmt.Sscanf(idStr, "%d", &teamID); err != nil || teamID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}
	var req domain.UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	team, err := h.usecase.UpdateTeam(teamID, req.Name)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update team", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Team updated successfully", "data": team})
}

// DeleteTeam handles DELETE /auth/teams/:id (CEO only)
func (h *AuthHandler) DeleteTeam(c *gin.Context) {
	idStr := c.Param("id")
	var teamID uint
	if _, err := fmt.Sscanf(idStr, "%d", &teamID); err != nil || teamID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}
	if err := h.usecase.DeleteTeam(teamID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete team", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Team deleted successfully"})
}

// AssignUserToTeam handles PATCH /auth/users/:id/assign-team (CEO only)
// Assigns (or unassigns when team_id is null) a user to a team.
func (h *AuthHandler) AssignUserToTeam(c *gin.Context) {
	idStr := c.Param("id")
	var targetUserID uint
	if _, err := fmt.Sscanf(idStr, "%d", &targetUserID); err != nil || targetUserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	var req domain.AssignUserToTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	callerID := getUserIDFromContext(c)
	if err := h.usecase.AssignUserToTeam(callerID, targetUserID, req.TeamID); err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "unauthorized") {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": "Failed to assign user", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User team assignment updated successfully"})
}

// GetTeamsFeature handles GET /auth/settings/teams-feature (any authenticated user)
func (h *AuthHandler) GetTeamsFeature(c *gin.Context) {
	enabled, err := h.usecase.GetTeamsFeatureEnabled()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get setting", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": domain.TeamsFeatureSettingResponse{Enabled: enabled}})
}

// SetTeamsFeature handles PUT /auth/settings/teams-feature (CEO/MANAGER only)
func (h *AuthHandler) SetTeamsFeature(c *gin.Context) {
	var req domain.SetTeamsFeatureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	callerID := getUserIDFromContext(c)
	if callerID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	if err := h.usecase.SetTeamsFeatureEnabled(callerID, req.Enabled); err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update setting", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Teams feature updated", "data": domain.TeamsFeatureSettingResponse{Enabled: req.Enabled}})
}
