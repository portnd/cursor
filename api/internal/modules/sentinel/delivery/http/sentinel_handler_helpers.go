// Shared helpers: JWT context, task id/code resolution, role checks, string contains, readAll.
package http

import (
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

// callerCtx extracts CallerContext (role + team_id + user_id) from the Gin context (set by AuthMiddleware).
func callerCtx(c *gin.Context) domain.CallerContext {
	role, _ := c.Get("role")
	teamID, _ := c.Get("team_id")
	ctx := domain.CallerContext{UserID: getUserIDFromContext(c)}
	if r, ok := role.(string); ok {
		ctx.Role = r
	}
	if t, ok := teamID.(*uint); ok {
		ctx.TeamID = t
	}
	return ctx
}

func isPrivilegedTaskRole(role string) bool {
	switch strings.ToUpper(strings.TrimSpace(role)) {
	case domain.RoleCEO, domain.RoleProductOwner, domain.RoleManager:
		return true
	default:
		return false
	}
}

func canUserCreateSubtask(parent *domain.Task, userID uint, role string) bool {
	if isPrivilegedTaskRole(role) {
		return true
	}
	if parent.CreatedBy != nil && *parent.CreatedBy == userID {
		return true
	}
	if parent.AssignedTo != nil && *parent.AssignedTo == userID {
		return true
	}
	return false
}

func (h *SentinelHandler) resolveTaskIDOrCode(c *gin.Context) (*domain.Task, error) {
	idStr := c.Param("id")
	if idStr == "" {
		return nil, fmt.Errorf("task id or code is required")
	}
	return h.usecase.GetTaskByIDOrCode(idStr)
}

// Helper to extract UserID safely from context
// The auth middleware sets "user_id" as float64 (from JSON)
func getUserIDFromContext(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}

	// JWT claims store numbers as float64
	if uid, ok := userID.(float64); ok {
		return uint(uid)
	}

	// Fallback for other types
	if uid, ok := userID.(uint); ok {
		return uid
	}

	if uid, ok := userID.(int); ok {
		return uint(uid)
	}

	return 0
}

// getUserRoleFromContext extracts the user's role from JWT context
func getUserRoleFromContext(c *gin.Context) string {
	role, exists := c.Get("role")
	if !exists {
		return ""
	}

	if roleStr, ok := role.(string); ok {
		return roleStr
	}

	return ""
}

// Helper to check if error message contains a substring
func contains(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr || (len(str) > len(substr) &&
		(str[:len(substr)] == substr || contains(str[1:], substr))))
}

// readAll reads all bytes from a ReadCloser (alias for io.ReadAll used in handlers).
func readAll(r interface{ Read([]byte) (int, error) }) ([]byte, error) {
	return io.ReadAll(r)
}
