package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateTableRequest struct {
	Name    string `json:"name" binding:"required"`
	Schema  string `json:"schema"`
	Comment string `json:"comment"`
}

func (a *Api) postPlatformPgMetaTables(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	projectRef := c.Param("ref")

	var req CreateTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// TODO: Phase 3/4 - Connect to actual project PostgreSQL and create table
	// For now, return a mock successful response
	// Studio expects a table object back with these fields

	a.logger.Info("Table creation requested (stub)",
		"project", projectRef,
		"table_name", req.Name,
		"schema", req.Schema,
		"comment", req.Comment)

	// Return mock table object
	c.JSON(http.StatusOK, gin.H{
		"id":                   1,
		"schema":               req.Schema,
		"name":                 req.Name,
		"rls_enabled":          false,
		"rls_forced":           false,
		"replica_identity":     "DEFAULT",
		"bytes":                0,
		"size":                 "0 bytes",
		"live_rows_estimate":   0,
		"dead_rows_estimate":   0,
		"comment":              req.Comment,
		"columns":              []interface{}{},
		"primary_keys":         []interface{}{},
		"relationships":        []interface{}{},
		"grants":               []interface{}{},
	})
}
