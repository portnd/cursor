package http

import (
	"github.com/gin-gonic/gin"
	"github.com/portnd/the-sentinel-core/internal/modules/deployment/domain"
)

// RegisterRoutes mounts all deployment endpoints onto the provided router group.
// All routes require authentication (middleware applied by the caller in main.go).
//
// Access control summary:
//   POST   /deployment/requests              — any authenticated user except SUPPORT
//   GET    /deployment/requests              — engineers see their own; CEO/MANAGER/CHIEF_ENGINEER see all
//   GET    /deployment/requests/:id          — same scoping as list
//   PATCH  /deployment/requests/:id/pick     — CHIEF_ENGINEER only (start review)
//   PATCH  /deployment/requests/:id/approve  — CHIEF_ENGINEER only
//   PATCH  /deployment/requests/:id/reject   — CHIEF_ENGINEER only
//   PATCH  /deployment/requests/:id/deploy   — CHIEF_ENGINEER only (mark as deployed)
//   GET    /deployment/stats                 — all authenticated users
func RegisterRoutes(r *gin.RouterGroup, uc domain.Usecase) {
	h := newDeploymentHandler(uc)

	g := r.Group("/deployment")
	{
		g.POST("/requests", h.createRequest)
		g.GET("/requests", h.listRequests)
		g.GET("/requests/by-task/:task_id", h.getByTaskID) // must be before /:id
		g.GET("/requests/:id", h.getRequest)
		g.PATCH("/requests/:id/pick", h.pickForReview)
		g.PATCH("/requests/:id/approve", h.approveRequest)
		g.PATCH("/requests/:id/reject", h.rejectRequest)
		g.PATCH("/requests/:id/deploy", h.markDeployed)
		g.GET("/stats", h.getStats)
	}
}
