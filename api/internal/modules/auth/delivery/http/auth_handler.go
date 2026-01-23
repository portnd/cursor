package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/komgrip/starter-kit/internal/modules/auth/domain"
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
