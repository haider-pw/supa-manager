package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"supamanager.io/supa-manager/database"
	"supamanager.io/supa-manager/utils"
)

type ProjectCreationBody struct {
	CloudProvider                  string `json:"cloud_provider"`
	OrgId                          int32  `json:"org_id"`
	Name                           string `json:"name"`
	DbPass                         string `json:"db_pass"`
	DbRegion                       string `json:"db_region"`
	CustomSupabaseInternalRequests struct {
		Ami struct {
			SearchTags struct {
				TagPostgresVersion string `json:"tag:postgresVersion"`
			} `json:"search_tags"`
		} `json:"ami"`
	} `json:"custom_supabase_internal_requests"`
	DesiredInstanceSize string `json:"desired_instance_size"`
}

type ProjectCreationResponse struct {
	Id                       int32    `json:"id"`
	Ref                      string   `json:"ref"`
	Name                     string   `json:"name"`
	Status                   string   `json:"status"`
	OrganizationId           int32    `json:"organization_id"`
	CloudProvider            string   `json:"cloud_provider"`
	Region                   string   `json:"region"`
	InsertedAt               string   `json:"inserted_at"`
	Endpoint                 string   `json:"endpoint"`
	AnonKey                  string   `json:"anon_key"`
	ServiceKey               string   `json:"service_key"`
	IsBranchEnabled          bool     `json:"is_branch_enabled"`
	PreviewBranchRefs        []string `json:"preview_branch_refs"`
	IsPhysicalBackupsEnabled bool     `json:"is_physical_backups_enabled"`
	IsReadReplicasEnabled    bool     `json:"is_read_replicas_enabled"`
	DiskVolumeSizeGb         int32    `json:"disk_volume_size_gb"`
	SubscriptionId           string   `json:"subscription_id"`
}

func (a *Api) postPlatformProjects(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var createProject ProjectCreationBody
	if err := c.BindJSON(&createProject); err != nil {
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}

	proj, err := a.queries.CreateProject(c.Request.Context(), database.CreateProjectParams{
		ProjectRef:     utils.GenerateProjectRef(createProject.Name),
		ProjectName:    createProject.Name,
		OrganizationID: createProject.OrgId,
		JwtSecret:      uuid.New().String(),
		CloudProvider:  strings.ToUpper(createProject.CloudProvider),
		Region:         strings.ToUpper(createProject.DbRegion),
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	// Trigger async provisioning if enabled
	if a.provisioner != nil {
		go func() {
			// Use background context instead of request context to avoid cancellation
			ctx := context.Background()
			a.logger.Info(fmt.Sprintf("Starting async provisioning for project %s", proj.ProjectRef))
			// TODO: Uncomment when ProvisionProject is implemented
		// if err := a.provisioner.ProvisionProject(ctx, &proj); err != nil {
		if false {
				a.logger.Error(fmt.Sprintf("Provisioning failed for project %s: %v", proj.ProjectRef, err))
				// Update status to FAILED
				if _, updateErr := a.queries.UpdateProjectStatus(ctx, database.UpdateProjectStatusParams{
					ProjectRef: proj.ProjectRef,
					Status:     "FAILED",
				}); updateErr != nil {
					a.logger.Error(fmt.Sprintf("Failed to update status to FAILED: %v", updateErr))
				}
			} else {
				a.logger.Info(fmt.Sprintf("Provisioning completed successfully for project %s", proj.ProjectRef))
			}
		}()
	}

	// Get the real keys if available (for provisioned projects)
	anonKey := "a.b.c"
	serviceKey := "a.b.c"
	if proj.AnonKey.Valid && proj.AnonKey.String != "" {
		anonKey = proj.AnonKey.String
	}
	if proj.ServiceRoleKey.Valid && proj.ServiceRoleKey.String != "" {
		serviceKey = proj.ServiceRoleKey.String
	}

	c.JSON(http.StatusCreated, ProjectCreationResponse{
		Id:                       proj.ID,
		Ref:                      proj.ProjectRef,
		Name:                     proj.ProjectName,
		Status:                   proj.Status,
		OrganizationId:           proj.OrganizationID,
		CloudProvider:            proj.CloudProvider,
		Region:                   proj.Region,
		InsertedAt:               proj.CreatedAt.Time.Format("2006-01-02T15:04:05.999Z"),
		Endpoint:                 fmt.Sprintf("https://%s.%s", proj.ProjectRef, a.config.Domain.Base),
		AnonKey:                  anonKey,
		ServiceKey:               serviceKey,
		IsBranchEnabled:          false,
		PreviewBranchRefs:        []string{},
		IsPhysicalBackupsEnabled: false,
		IsReadReplicasEnabled:    false,
		DiskVolumeSizeGb:         0,
		SubscriptionId:           "wedontbill",
	})
}
