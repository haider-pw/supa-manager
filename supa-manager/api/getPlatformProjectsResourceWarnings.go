package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) getPlatformProjectsResourceWarnings(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		println(err.Error())
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// Return empty warnings list
	// TODO: Implement resource warnings (disk space, bandwidth, etc.) in Phase 3
	c.JSON(http.StatusOK, []interface{}{})
}
