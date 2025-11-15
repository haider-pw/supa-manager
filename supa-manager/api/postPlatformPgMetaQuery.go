package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type PgMetaQueryRequest struct {
	Query string `json:"query"`
}

func (a *Api) postPlatformPgMetaQuery(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		println(err.Error())
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	projectRef := c.Param("ref")
	key := c.Query("key")

	var req PgMetaQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// For now, return mock data based on the query key
	// TODO: In Phase 3, connect to actual project's PostgreSQL database
	switch key {
	case "schemas-list":
		// Return mock schema list
		c.JSON(http.StatusOK, []gin.H{
			{
				"id":    1,
				"name":  "public",
				"owner": "postgres",
			},
		})
	case "tables-list":
		// Return empty tables list for now
		c.JSON(http.StatusOK, []interface{}{})
	case "columns-list":
		// Return empty columns list
		c.JSON(http.StatusOK, []interface{}{})
	case "functions-list":
		// Return empty functions list
		c.JSON(http.StatusOK, []interface{}{})
	case "types-list":
		// Return empty types list
		c.JSON(http.StatusOK, []interface{}{})
	default:
		// For unknown keys, return empty result with proper format
		// executeSql wraps the response in { result: data }, so we just return the array
		a.logger.Info("Unknown pg-meta query key", "key", key, "project", projectRef)
		c.JSON(http.StatusOK, []gin.H{
			{
				"data": gin.H{
					"entities": []interface{}{},
					"count":    0,
				},
			},
		})
	}
}
