package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) getPlatformProjectDatabases(c *gin.Context) {
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

	a.logger.Info("Fetching project databases", "project", projectRef)

	// Determine database host based on project configuration
	dbHost := "localhost"
	if project.DockerNetworkName.Valid && project.DockerNetworkName.String != "" {
		// Use Docker network name for internal communication
		dbHost = fmt.Sprintf("db.%s", projectRef)
	}

	// Determine database port
	dbPort := int32(5432)
	if project.PostgresPort.Valid {
		dbPort = project.PostgresPort.Int32
	}

	// Map project status to database status
	dbStatus := "UNKNOWN"
	switch project.Status {
	case "ACTIVE_HEALTHY":
		dbStatus = "ACTIVE_HEALTHY"
	case "PROVISIONING":
		dbStatus = "PROVISIONING"
	case "PAUSED":
		dbStatus = "PAUSED"
	case "INACTIVE":
		dbStatus = "INACTIVE"
	default:
		dbStatus = "UNKNOWN"
	}

	// Return real database information from project configuration
	c.JSON(http.StatusOK, []gin.H{
		{
			"id":             project.ID,
			"ref":            project.ProjectRef,
			"name":           "postgres",
			"host":           dbHost,
			"port":           dbPort,
			"user":           "postgres",
			"version":        "15.1.0.147",
			"status":         dbStatus,
			"inserted_at":    project.CreatedAt,
			"cloud_provider": project.CloudProvider,
			"region":         project.Region,
		},
	})
}
