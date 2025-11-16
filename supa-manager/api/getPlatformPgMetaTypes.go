package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) getPlatformPgMetaTypes(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// Return empty types list for now
	// TODO: In Phase 3, connect to actual project's PostgreSQL database
	c.JSON(http.StatusOK, []interface{}{})
}

func (a *Api) getPlatformPgMetaPublications(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// Return empty publications list for now
	// TODO: In Phase 3, connect to actual project's PostgreSQL database
	c.JSON(http.StatusOK, []interface{}{})
}
