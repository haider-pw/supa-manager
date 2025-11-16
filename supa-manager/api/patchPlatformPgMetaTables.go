package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UpdateTableRequest struct {
	RLSEnabled *bool   `json:"rls_enabled"`
	RLSForced  *bool   `json:"rls_forced"`
	Comment    *string `json:"comment"`
}

func (a *Api) patchPlatformPgMetaTables(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	projectRef := c.Param("ref")
	tableID := c.Query("id")

	var req UpdateTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// TODO: Phase 3/4 - Connect to actual project PostgreSQL and update table
	// For now, return a mock successful response

	a.logger.Info("Table update requested (stub)",
		"project", projectRef,
		"table_id", tableID,
		"rls_enabled", req.RLSEnabled,
		"rls_forced", req.RLSForced)

	// Return mock table object with updated values
	c.JSON(http.StatusOK, gin.H{
		"id":                 tableID,
		"schema":             "public",
		"name":               "users",
		"rls_enabled":        req.RLSEnabled,
		"rls_forced":         req.RLSForced,
		"replica_identity":   "DEFAULT",
		"bytes":              0,
		"size":               "0 bytes",
		"live_rows_estimate": 0,
		"dead_rows_estimate": 0,
		"comment":            req.Comment,
		"columns":            []interface{}{},
		"primary_keys":       []interface{}{},
		"relationships":      []interface{}{},
		"grants":             []interface{}{},
	})
}
