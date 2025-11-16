package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (a *Api) getProjectAnalyticsEndpointUsage(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// Return empty analytics data
	// TODO: Implement actual analytics tracking in future phases
	c.JSON(http.StatusOK, gin.H{
		"data": []gin.H{},
		"total": 0,
	})
}

func (a *Api) getPlatformProjectAnalyticsEndpointUsage(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// Return mock analytics data for the dashboard
	now := time.Now()
	data := []gin.H{}

	// Generate 24 hours of mock data points
	for i := 23; i >= 0; i-- {
		timestamp := now.Add(time.Duration(-i) * time.Hour)
		data = append(data, gin.H{
			"timestamp": timestamp.Format(time.RFC3339),
			"count":     0,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  data,
		"total": 0,
	})
}
