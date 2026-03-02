package http

import (
	"fmt"
	"net/http"

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
// Request Body: { "role": "PM" }
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
