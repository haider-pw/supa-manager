package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) getIntegrations(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		println(err.Error())
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// orgId := c.Param("id")
	// expand := c.Query("expand")

	// Return empty integrations array directly (UI expects array, not object)
	// TODO: Implement integrations (GitHub, Vercel, etc.) in future phases
	c.JSON(http.StatusOK, []interface{}{})
}
