package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) getPlatformIntegrationAuthorization(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		println(err.Error())
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// integration := c.Param("integration") // e.g., "github"

	// Return empty authorization status
	// TODO: Implement GitHub OAuth integration in future phases
	c.JSON(http.StatusOK, gin.H{
		"authorized":    false,
		"authorization": nil,
		"metadata":      gin.H{},
	})
}
