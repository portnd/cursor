package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

// B2BHandler handles Internal B2B Outsource Request endpoints.
type B2BHandler struct {
	usecase domain.SentinelUsecase
}

func NewB2BHandler(usecase domain.SentinelUsecase) *B2BHandler {
	return &B2BHandler{usecase: usecase}
}

// DTOs

type createB2BRequestReq struct {
	Title            string `json:"title" binding:"required"`
	Description      string `json:"description"`
	EstimatedMinutes int    `json:"estimated_minutes" binding:"required,min=1"`
	TargetTeamID     uint   `json:"target_team_id" binding:"required"`
}

type updateB2BRequestReq struct {
	Action          string `json:"action" binding:"required,oneof=COUNTER REJECT"` // COUNTER | REJECT
	ProposedMinutes int    `json:"proposed_minutes"`
	Reason          string `json:"reason"`
}

// CreateB2BRequest handles POST /sentinel/b2b/requests
// PM/CEO sends an outsource request to another team.
func (h *B2BHandler) CreateB2BRequest(c *gin.Context) {
	callerTeamID, ok := getCallerTeamID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "you must belong to a team to create a B2B request"})
		return
	}

	var req createB2BRequestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}

	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	b2bReq, err := h.usecase.CreateB2BRequest(req.Title, req.Description, req.EstimatedMinutes, callerTeamID, req.TargetTeamID, userID)
	if err != nil {
		if domain.IsBadRequest(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create B2B request", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": b2bReq, "message": "B2B request created"})
}

// GetB2BRequests handles GET /sentinel/b2b/requests?direction=inbound|outbound
func (h *B2BHandler) GetB2BRequests(c *gin.Context) {
	callerTeamID, ok := getCallerTeamID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "you must belong to a team to view B2B requests"})
		return
	}

	direction := c.DefaultQuery("direction", "inbound")
	reqs, err := h.usecase.GetB2BRequests(callerTeamID, direction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch B2B requests", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reqs})
}

// UpdateB2BRequest handles PATCH /sentinel/b2b/requests/:id
// Target team PM can counter-offer or reject; requester team PM can reject.
func (h *B2BHandler) UpdateB2BRequest(c *gin.Context) {
	callerTeamID, ok := getCallerTeamID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "you must belong to a team"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "invalid request ID"})
		return
	}

	var req updateB2BRequestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}

	var result *domain.B2BRequest
	switch req.Action {
	case "COUNTER":
		if req.ProposedMinutes <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "proposed_minutes must be greater than 0"})
			return
		}
		result, err = h.usecase.CounterOfferB2BRequest(id, callerTeamID, req.ProposedMinutes, req.Reason)
	case "REJECT":
		result, err = h.usecase.RejectB2BRequest(id, callerTeamID)
	}

	if err != nil {
		if domain.IsBadRequest(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update B2B request", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// AcceptB2BRequest handles POST /sentinel/b2b/requests/:id/accept
// Target team PM accepts the request; a Task is created in target team's project.
func (h *B2BHandler) AcceptB2BRequest(c *gin.Context) {
	callerTeamID, ok := getCallerTeamID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "you must belong to a team"})
		return
	}

	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "invalid request ID"})
		return
	}

	task, err := h.usecase.AcceptB2BRequest(id, callerTeamID, userID)
	if err != nil {
		if domain.IsBadRequest(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to accept B2B request", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task, "message": "Request accepted and task created"})
}

// getCallerTeamID extracts team_id from gin context (set by AuthMiddleware).
// Returns (0, false) if caller has no team.
func getCallerTeamID(c *gin.Context) (uint, bool) {
	raw, exists := c.Get("team_id")
	if !exists {
		return 0, false
	}
	teamIDPtr, ok := raw.(*uint)
	if !ok || teamIDPtr == nil {
		return 0, false
	}
	return *teamIDPtr, true
}
