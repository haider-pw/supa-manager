package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpgradeStatusResponse struct {
	Status              string `json:"status"`
	IsEligible          bool   `json:"is_eligible"`
	CurrentVersion      string `json:"current_version"`
	LatestVersion       string `json:"latest_version"`
	InitiatedAt         string `json:"initiated_at,omitempty"`
	UpgradeAvailable    bool   `json:"upgrade_available"`
	RequiresDowntime    bool   `json:"requires_downtime"`
	EstimatedDuration   int    `json:"estimated_duration_minutes"`
}

func (a *Api) getProjectUpgradeStatus(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	projectRef := c.Param("ref")
	if projectRef == "" {
		c.JSON(400, gin.H{"error": "Project reference is required"})
		return
	}

	project, err := a.queries.GetProjectByRef(c.Request.Context(), projectRef)
	if err != nil {
		c.JSON(404, gin.H{"error": "Project not found"})
		return
	}

	response := UpgradeStatusResponse{
		Status:            "stable",
		IsEligible:        true,
		CurrentVersion:    "1.24.04",
		LatestVersion:     "1.24.04",
		UpgradeAvailable:  false,
		RequiresDowntime:  false,
		EstimatedDuration: 0,
	}

	if project.Status == "PROVISIONING" {
		response.Status = "provisioning"
		response.IsEligible = false
	}

	c.JSON(http.StatusOK, response)
}
