package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) getProjectUpgradeEligibility(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		println(err.Error())
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// Return upgrade eligibility info
	// For self-hosted, upgrades are managed differently
	c.JSON(http.StatusOK, gin.H{
		"eligible":           false,
		"reason":             "self_hosted",
		"current_version":    "1.0.0",
		"target_version":     "1.0.0",
		"available_upgrades": []interface{}{},
	})
}
