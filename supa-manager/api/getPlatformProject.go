package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) getPlatformProject(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	projectRef := c.Param("ref")
	project, err := a.queries.GetProjectByRef(c, projectRef)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	// Extract database connection details from provisioned infrastructure
	dbHost := "localhost" // Default for local Docker setup
	dbPort := int32(0)
	dbName := "postgres"
	dbUser := "postgres"

	// If project is provisioned, use actual port
	if project.PostgresPort.Valid {
		dbPort = project.PostgresPort.Int32
	}

	// Format inserted_at timestamp
	insertedAt := ""
	if project.CreatedAt.Valid {
		insertedAt = project.CreatedAt.Time.Format("2006-01-02T15:04:05.999Z")
	}

	c.JSON(http.StatusOK, Project{
		Id:                       project.ID,
		Ref:                      project.ProjectRef,
		Name:                     project.ProjectName,
		Status:                   project.Status,
		OrganizationId:           project.OrganizationID,
		InsertedAt:               insertedAt,
		SubscriptionId:           "free-tier",
		CloudProvider:            project.CloudProvider,
		Region:                   project.Region,
		DiskVolumeSizeGb:         8, // Default volume size
		Size:                     "small",
		DbUserSupabase:           dbUser,
		DbPassSupabase:           "", // Never expose password in API
		DbDnsName:                dbHost,
		DbHost:                   dbHost,
		DbPort:                   dbPort,
		DbName:                   dbName,
		SslEnforced:              false,
		WalgEnabled:              false,
		InfraComputeSize:         "small",
		PreviewBranchRefs:        []interface{}{},
		IsBranchEnabled:          false,
		IsPhysicalBackupsEnabled: false,
	})
}
