package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) getPlatformOrganizationUsage(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		println(err.Error())
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// orgId := c.Param("slug") // Organization ID

	// Return empty usage data
	// TODO: Implement actual usage tracking in Phase 3
	c.JSON(http.StatusOK, gin.H{
		"usage": gin.H{
			"storage":   0,
			"bandwidth": 0,
			"egress":    0,
		},
		"limits": gin.H{
			"storage":   10737418240, // 10 GB
			"bandwidth": 53687091200, // 50 GB
		},
		"billing_cycle_start": "2025-11-01T00:00:00Z",
		"billing_cycle_end":   "2025-12-01T00:00:00Z",
	})
}
