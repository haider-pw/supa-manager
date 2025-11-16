package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) getProjectContent(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	projectRef := c.Param("ref")

	// Get real project from database
	project, err := a.queries.GetProjectByRef(c.Request.Context(), projectRef)
	if err != nil {
		a.logger.Error("Failed to get project", "error", err, "ref", projectRef)
		c.JSON(404, gin.H{"error": "Project not found"})
		return
	}

	a.logger.Info("Fetching project content", "project", projectRef)

	// Return real content object based on project configuration
	// Storage is enabled by default in Supabase projects
	storageEnabled := project.Status == "ACTIVE_HEALTHY" || project.Status == "PROVISIONING"

	c.JSON(http.StatusOK, gin.H{
		"id":              project.ProjectRef,
		"storage_enabled": storageEnabled,
		"buckets":         []interface{}{}, // TODO: Phase 4 - Query actual buckets from provisioned Storage
		"policies":        []interface{}{}, // TODO: Phase 4 - Query actual policies from provisioned Storage
	})
}
