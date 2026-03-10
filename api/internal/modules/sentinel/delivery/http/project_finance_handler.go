package http

import (
	"fmt"
	"net/http"

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
