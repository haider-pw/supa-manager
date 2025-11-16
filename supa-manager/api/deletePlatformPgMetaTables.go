package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) deletePlatformPgMetaTables(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	projectRef := c.Param("ref")
	tableID := c.Query("id")
	cascade := c.Query("cascade") == "true"

	// TODO: Phase 3/4 - Connect to actual project PostgreSQL and delete table
	// For now, return a mock successful response

	a.logger.Info("Table deletion requested (stub)",
		"project", projectRef,
		"table_id", tableID,
		"cascade", cascade)

	// Return success response (Studio expects no content or a simple object)
	c.JSON(http.StatusOK, gin.H{
		"id":      tableID,
		"message": "Table deleted successfully",
	})
}
