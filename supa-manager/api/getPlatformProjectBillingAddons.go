package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) getPlatformProjectBillingAddons(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		println(err.Error())
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// Return empty addons array
	// TODO: Implement billing addons (compute, storage, bandwidth upgrades) in future phases
	c.JSON(http.StatusOK, gin.H{
		"addons": []interface{}{},
		"available_addons": []interface{}{
			gin.H{
				"id":          "compute-upgrade",
				"name":        "Compute Upgrade",
				"description": "Additional compute capacity",
				"price":       0,
			},
			gin.H{
				"id":          "storage-upgrade",
				"name":        "Storage Upgrade",
				"description": "Additional storage capacity",
				"price":       0,
			},
		},
	})
}
