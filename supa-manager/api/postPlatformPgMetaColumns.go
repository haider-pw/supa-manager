package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateColumnRequest struct {
	TableID       int    `json:"tableId" binding:"required"`
	Name          string `json:"name" binding:"required"`
	Type          string `json:"type" binding:"required"`
	IsIdentity    *bool  `json:"isIdentity"`
	IsUnique      *bool  `json:"isUnique"`
	IsPrimaryKey  *bool  `json:"isPrimaryKey"`
	IsNullable    *bool  `json:"isNullable"`
	DefaultValue  string `json:"defaultValue"`
	Comment       string `json:"comment"`
	IdentityStart int    `json:"identityStart"`
	IdentityIncr  int    `json:"identityIncr"`
}

func (a *Api) postPlatformPgMetaColumns(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	projectRef := c.Param("ref")

	var req CreateColumnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// TODO: Phase 3/4 - Connect to actual project PostgreSQL and create column
	// For now, return a mock successful response

	a.logger.Info("Column creation requested (stub)",
		"project", projectRef,
		"table_id", req.TableID,
		"column_name", req.Name,
		"type", req.Type,
		"is_identity", req.IsIdentity,
		"is_primary_key", req.IsPrimaryKey)

	// Return mock column object
	c.JSON(http.StatusOK, gin.H{
		"table_id":        req.TableID,
		"schema":          "public",
		"table":           "users",
		"id":              "1",
		"ordinal_position": 1,
		"name":            req.Name,
		"default_value":   req.DefaultValue,
		"data_type":       req.Type,
		"format":          req.Type,
		"is_identity":     req.IsIdentity,
		"identity_generation": nil,
		"is_nullable":     req.IsNullable,
		"is_updatable":    true,
		"is_unique":       req.IsUnique,
		"enums":           []interface{}{},
		"comment":         req.Comment,
	})
}
