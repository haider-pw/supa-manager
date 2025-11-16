package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) getProjectSupervisor(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// Return empty supervisor status
	// TODO: Implement actual supervisor monitoring in future phases
	c.JSON(http.StatusOK, gin.H{
		"supervisor": gin.H{
			"status":  "running",
			"version": "1.0.0",
		},
	})
}
