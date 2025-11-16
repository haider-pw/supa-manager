package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) getProjectHealth(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// Return healthy status for all services
	// TODO: Implement actual health checking in future phases
	c.JSON(http.StatusOK, gin.H{
		"services": []gin.H{
			{
				"name":   "auth",
				"status": "healthy",
			},
			{
				"name":   "realtime",
				"status": "healthy",
			},
			{
				"name":   "rest",
				"status": "healthy",
			},
			{
				"name":   "storage",
				"status": "healthy",
			},
			{
				"name":   "db",
				"status": "healthy",
			},
		},
	})
}
