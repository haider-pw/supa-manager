package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) getProjectJwtSecretUpdateStatus(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	projectRef := c.Param("ref")
	_, err = a.queries.GetProjectByRef(c, projectRef)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	// TODO: Stub implementation - always returns "completed" status
	// In a real implementation, this should:
	// 1. Track JWT secret rotation operations in the database
	// 2. Return actual status ("pending", "in_progress", "completed", "failed")
	// 3. Store operation progress and timestamps
	// For now, this prevents 404 errors in Studio UI while we implement
	// full JWT secret rotation in Phase 3/4
	c.JSON(http.StatusOK, gin.H{
		"jwtSecretUpdateStatus": gin.H{
			"status":   "completed",
			"progress": 100,
		},
	})
}
