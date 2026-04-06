package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

type projectFinanceHandler struct {
	usecase domain.ProjectFinanceUsecase
}

func NewProjectFinanceHandler(usecase domain.ProjectFinanceUsecase) *projectFinanceHandler {
	return &projectFinanceHandler{usecase: usecase}
}

// GetProjectCapital GET /sentinel/projects/:id/finance/capital
func (h *projectFinanceHandler) GetProjectCapital(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}
	resp, err := h.usecase.GetProjectCapital(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// GetProjectCapitals GET /sentinel/projects/finance/capital?project_ids=uuid1,uuid2
func (h *projectFinanceHandler) GetProjectCapitals(c *gin.Context) {
	idsCSV := strings.TrimSpace(c.Query("project_ids"))
	if idsCSV == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_ids query param is required"})
		return
	}
	parts := strings.Split(idsCSV, ",")
	projectIDs := make([]uuid.UUID, 0, len(parts))
	seen := make(map[uuid.UUID]struct{}, len(parts))
	for _, part := range parts {
		idStr := strings.TrimSpace(part)
		if idStr == "" {
			continue
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "project_ids must be comma-separated UUIDs"})
			return
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		projectIDs = append(projectIDs, id)
	}
	items, err := h.usecase.GetProjectCapitals(projectIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// InjectProjectCapital POST /sentinel/projects/:id/finance/inject
func (h *projectFinanceHandler) InjectProjectCapital(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}
	var req domain.InjectProjectCapitalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	project, err := h.usecase.InjectProjectCapital(projectID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": project})
}

// EditProjectCapital PUT /sentinel/projects/:id/finance/capital
func (h *projectFinanceHandler) EditProjectCapital(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}
	var req domain.EditProjectCapitalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	project, err := h.usecase.EditProjectCapital(projectID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": project})
}

// CloseProjectCycleAndPayout POST /sentinel/projects/:id/finance/close-cycle
func (h *projectFinanceHandler) CloseProjectCycleAndPayout(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}
	resp, err := h.usecase.CloseProjectCycleAndPayout(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// DeleteProjectTransaction DELETE /sentinel/projects/:id/finance/transactions/:txID
func (h *projectFinanceHandler) DeleteProjectTransaction(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}
	var txID int64
	if _, err := fmt.Sscan(c.Param("txID"), &txID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction id"})
		return
	}
	if err := h.usecase.DeleteProjectTransaction(txID, projectID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
