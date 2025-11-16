package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"supamanager.io/supa-manager/utils"
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
	dbPort := int32(5432) // Default PostgreSQL port
	dbName := "postgres"
	dbUser := "postgres"

	// If project is provisioned, use actual port
	if project.PostgresPort.Valid {
		dbPort = project.PostgresPort.Int32
	}

	// Construct database hostname using domain pattern: db.{project-ref}.{base-domain}
	// Example: db.p1-7cda6e.supamanager.org (matches Supabase pattern)
	dbHost := "db." + project.ProjectRef + "." + a.config.Domain.Base
	dbDnsName := dbHost

	// Format inserted_at timestamp
	insertedAt := ""
	if project.CreatedAt.Valid {
		insertedAt = project.CreatedAt.Time.Format("2006-01-02T15:04:05.999Z")
	}

	// Generate encrypted connection string
	// Format: postgresql://user:password@host:port/database
	// Note: Using a placeholder password - in production this should come from provisioner
	dbPassword := "postgres" // TODO: Retrieve from secure storage
	connectionStringPlain := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Encrypt connection string using ENCRYPTION_SECRET
	connectionString, err := utils.EncryptAES(connectionStringPlain, a.config.EncryptionSecret)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Failed to encrypt connection string: %v", err))
		connectionString = "" // Don't expose plaintext on error
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
		DbDnsName:                dbDnsName,
		DbHost:                   dbHost,
		DbPort:                   dbPort,
		DbName:                   dbName,
		SslEnforced:              false,
		WalgEnabled:              false,
		InfraComputeSize:         "small",
		PreviewBranchRefs:        []interface{}{},
		IsBranchEnabled:          false,
		IsPhysicalBackupsEnabled: false,
		ConnectionString:         connectionString, // Encrypted connection string
	})
}
