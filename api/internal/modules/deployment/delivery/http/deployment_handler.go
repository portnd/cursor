package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/portnd/the-sentinel-core/internal/modules/deployment/domain"
)

type deploymentHandler struct {
	uc domain.Usecase
}

func newDeploymentHandler(uc domain.Usecase) *deploymentHandler {
	return &deploymentHandler{uc: uc}
}

// ─── helpers ─────────────────────────────────────────────────────────────────

func callerInfo(c *gin.Context) (uint, string) {
	rawID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	var callerID uint
	switch v := rawID.(type) {
	case float64:
		callerID = uint(v)
	case uint:
		callerID = v
	}
	roleStr, _ := role.(string)
	return callerID, roleStr
}

func respondErr(c *gin.Context, err error) {
	switch err {
	case domain.ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case domain.ErrForbidden:
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case domain.ErrInvalidStatus, domain.ErrAlreadyReviewing:
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case domain.ErrInvalidEnv, domain.ErrMissingTitle, domain.ErrMissingBranch:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// ─── handlers ────────────────────────────────────────────────────────────────

// POST /deployment/requests
func (h *deploymentHandler) createRequest(c *gin.Context) {
	callerID, callerRole := callerInfo(c)
	var in domain.CreateDeploymentRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.uc.CreateRequest(callerID, callerRole, in)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

// GET /deployment/requests
func (h *deploymentHandler) listRequests(c *gin.Context) {
	callerID, callerRole := callerInfo(c)
	status := c.Query("status")
	results, err := h.uc.ListRequests(callerID, callerRole, status)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results})
}

// GET /deployment/requests/:id
func (h *deploymentHandler) getRequest(c *gin.Context) {
	callerID, callerRole := callerInfo(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	result, err := h.uc.GetRequest(callerID, callerRole, uint(id))
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// PATCH /deployment/requests/:id/pick
func (h *deploymentHandler) pickForReview(c *gin.Context) {
	callerID, callerRole := callerInfo(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	result, err := h.uc.PickForReview(callerID, callerRole, uint(id))
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// PATCH /deployment/requests/:id/approve
func (h *deploymentHandler) approveRequest(c *gin.Context) {
	callerID, callerRole := callerInfo(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in domain.ReviewActionRequest
	_ = c.ShouldBindJSON(&in)
	result, err := h.uc.Approve(callerID, callerRole, uint(id), in)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// PATCH /deployment/requests/:id/reject
func (h *deploymentHandler) rejectRequest(c *gin.Context) {
	callerID, callerRole := callerInfo(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in domain.ReviewActionRequest
	_ = c.ShouldBindJSON(&in)
	result, err := h.uc.Reject(callerID, callerRole, uint(id), in)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// PATCH /deployment/requests/:id/deploy
func (h *deploymentHandler) markDeployed(c *gin.Context) {
	callerID, callerRole := callerInfo(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in domain.ReviewActionRequest
	_ = c.ShouldBindJSON(&in)
	result, err := h.uc.MarkDeployed(callerID, callerRole, uint(id), in)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GET /deployment/requests/by-task/:task_id
func (h *deploymentHandler) getByTaskID(c *gin.Context) {
	callerID, callerRole := callerInfo(c)
	taskIDStr := c.Param("task_id")
	result, err := h.uc.GetRequestByTaskID(callerID, callerRole, taskIDStr)
	if err != nil {
		if err == domain.ErrNotFound {
			// Missing deployment request for a task is an expected empty state.
			c.JSON(http.StatusOK, nil)
			return
		}
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GET /deployment/stats
func (h *deploymentHandler) getStats(c *gin.Context) {
	callerID, callerRole := callerInfo(c)
	stats, err := h.uc.GetStats(callerID, callerRole)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, stats)
}
