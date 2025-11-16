package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) getProjectCustomHostname(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		println(err.Error())
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// ref := c.Param("ref")

	// Return empty custom hostname configuration
	// TODO: Implement custom domain management in future phases
	c.JSON(http.StatusOK, gin.H{
		"customHostname": "",
		"status":         "not_configured",
		"data":           gin.H{},
	})
}
