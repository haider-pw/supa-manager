package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matthewhartstonge/argon2"
	"log/slog"
	"net/http"
	"supamanager.io/supa-manager/conf"
	"supamanager.io/supa-manager/database"
	"supamanager.io/supa-manager/provisioner"
	"time"
)

type Api struct {
	isHealthy   bool
	logger      *slog.Logger
	config      *conf.Config
	queries     *database.Queries
	pgPool      *pgxpool.Pool
	argon       argon2.Config
	provisioner provisioner.Provisioner
}

func CreateApi(logger *slog.Logger, config *conf.Config) (*Api, error) {
	conn, err := pgxpool.New(context.Background(), config.DatabaseUrl)
	if err != nil {
		logger.Error(fmt.Sprintf("Unable to connect to database: %v", err))
		return nil, err
	}

	if err := conf.EnsureMigrationsTableExists(conn); err != nil {
		logger.Error(fmt.Sprintf("Failed to ensure migrations table: %v", err))
		return nil, err
	}

	queries := database.New(conn)

	if success, err := conf.EnsureMigrations(conn, queries); err != nil || !success {
		logger.Error(fmt.Sprintf("Failed to run migrations: %v", err))
		return nil, err
	}

	// Initialize provisioner if enabled
	var prov provisioner.Provisioner
	if config.Provisioning.Enabled {
		// Use NewDockerProvisioner with projects directory and templates directory
		dockerProv, err := provisioner.NewDockerProvisioner(
			config.Provisioning.ProjectsDir,
			"./templates", // Template directory for docker-compose files
		)
		if err != nil {
			logger.Warn(fmt.Sprintf("Failed to initialize provisioner: %v", err))
			logger.Info("Continuing without provisioner - projects can be created but not provisioned")
		} else {
			prov = dockerProv
			logger.Info("Docker provisioner initialized and enabled")
		}
	} else {
		logger.Info("Provisioner is disabled")
	}

	return &Api{
		logger:      logger,
		config:      config,
		queries:     queries,
		pgPool:      conn,
		argon:       argon2.DefaultConfig(),
		provisioner: prov,
	}, nil
}

func (a *Api) GetAccountIdFromRequest(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}

	tokenString := authHeader[len("Bearer "):]
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.config.JwtSecret), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	return claims.Subject, nil
}

func (a *Api) GetAccountFromRequest(c *gin.Context) (*database.Account, error) {
	id, err := a.GetAccountIdFromRequest(c)
	if err != nil {
		return nil, err
	}

	if id == "" {
		return nil, errors.New("missing account ID")
	}

	account, err := a.queries.GetAccountByGoTrueID(c.Request.Context(), id)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (a *Api) ListenAddress() string {
	return ":8080"
}

func (a *Api) index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (a *Api) status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"is_healthy": a.isHealthy})
}

func (a *Api) telemetry(c *gin.Context) {
	c.AbortWithStatus(http.StatusNoContent)
}

const INDEX = ""

func (a *Api) Router() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", a.index)
	r.GET("/status", a.status)

	profile := r.Group("/profile")
	{
		profile.GET(INDEX, a.getProfile)
		profile.GET("/permissions", a.getProfilePermissions)
		profile.POST("/password-check", a.postPasswordCheck)
	}

	organization := r.Group("/organizations")
	{
		organization.GET(INDEX, a.getOrganizations)

		specificOrganization := organization.Group("/:slug")
		{
			members := specificOrganization.Group("/members")
			{
				members.GET("/reached-free-project-limit", a.getOrganizationMembersReachedFreeProjectLimit)
			}
		}
	}

	projects := r.Group("/projects")
	{
		specificProject := projects.Group("/:ref")
		{
			specificProject.GET("/status", a.getProjectStatus)
			specificProject.GET("/jwt-secret-update-status", a.getProjectJwtSecretUpdateStatus)
			specificProject.GET("/api", a.getProjectApi)
			specificProject.GET("/upgrade/status", a.getProjectUpgradeStatus)
			specificProject.GET("/health", a.getProjectHealth)
			specificProject.GET("/supervisor", a.getProjectSupervisor)
			// TODO: Uncomment when implementing provisioning
			// specificProject.POST("/pause", a.postProjectPause)
			// specificProject.POST("/resume", a.postProjectResume)
			// specificProject.DELETE(INDEX, a.deleteProject)

			// Analytics routes
			analytics := specificProject.Group("/analytics/endpoints")
			{
				analytics.GET("/usage.api-counts", a.getProjectAnalyticsEndpointUsage)
				analytics.GET("/usage.api-requests-count", a.getProjectAnalyticsEndpointUsage)
			}
		}
	}

	// Singular /project routes (some Studio UI calls use singular)
	project := r.Group("/project")
	{
		specificProject := project.Group("/:ref")
		{
			specificProject.GET("/status", a.getProjectStatus)
			specificProject.GET("/jwt-secret-update-status", a.getProjectJwtSecretUpdateStatus)
			specificProject.GET("/api", a.getProjectApi)
			specificProject.GET("/health", a.getProjectHealth)
			specificProject.GET("/supervisor", a.getProjectSupervisor)
		}
	}

	// Props routes (used by Studio UI)
	props := r.Group("/props")
	{
		propsProject := props.Group("/project")
		{
			specificProject := propsProject.Group("/:ref")
			{
				specificProject.GET("/jwt-secret-update-status", a.getPropsProjectJwtSecretUpdateStatus)
			}
		}
	}

	gotrue := r.Group("/auth")
	{
		gotrue.POST("/token", a.postGotrueToken)
	}

	platform := r.Group("/platform")
	{
		platform.POST("/signup", a.postPlatformSignup)
		platform.GET("/notifications", a.getPlatformNotifications)
		platform.GET("/notifications/summary", a.getPlatformNotificationsSummary)
		platform.GET("/stripe/invoices/overdue", a.getPlatformOverdueInvoices)
		platform.GET("/projects-resource-warnings", a.getPlatformProjectsResourceWarnings)

		// pg-meta routes for database metadata queries
		platformPgMeta := platform.Group("/pg-meta")
		{
			specificProject := platformPgMeta.Group("/:ref")
			{
				specificProject.POST("/query", a.postPlatformPgMetaQuery)
				specificProject.GET("/types", a.getPlatformPgMetaTypes)
				specificProject.GET("/publications", a.getPlatformPgMetaPublications)
			}
		}

		platformProjects := platform.Group("/projects")
		{
			platformProjects.GET(INDEX, a.getPlatformProjects)
			platformProjects.POST(INDEX, a.postPlatformProjects)
			specificProject := platformProjects.Group("/:ref")
			{
				specificProject.GET(INDEX, a.getPlatformProject)
				specificProject.GET("/settings", a.getPlatformProjectSettings)
				specificProject.GET("/billing/addons", a.getPlatformProjectBillingAddons)

				// Analytics routes
				analytics := specificProject.Group("/analytics/endpoints")
				{
					analytics.GET("/usage.api-counts", a.getPlatformProjectAnalyticsEndpointUsage)
					analytics.GET("/usage.api-requests-count", a.getPlatformProjectAnalyticsEndpointUsage)
				}
			}
		}

		// Singular /project routes (some Studio UI calls use singular)
		platformProject := platform.Group("/project")
		{
			specificProject := platformProject.Group("/:ref")
			{
				specificProject.GET(INDEX, a.getPlatformProject)
				specificProject.GET("/settings", a.getPlatformProjectSettings)
				specificProject.GET("/billing/addons", a.getPlatformProjectBillingAddons)
			}
		}

		platformOrganizations := platform.Group("/organizations")
		{
			platformOrganizations.POST(INDEX, a.postPlatformOrganizations)
			specificOrganization := platformOrganizations.Group("/:slug")
			{
				specificOrganization.GET("/billing/subscription", a.getPlatformOrganizationSubscription)
				specificOrganization.GET("/usage", a.getPlatformOrganizationUsage)
			}
		}

		platform.GET("/integrations/:integration/connections", a.getIntegrationConnections)
		platform.GET("/integrations/:integration/authorization", a.getPlatformIntegrationAuthorization)
		platform.GET("/integrations/:integration/repositories", a.getPlatformIntegrationRepositories)
	}

	// Integrations routes (organization level)
	integrations := r.Group("/integrations")
	{
		integrations.GET("/:id", a.getIntegrations)
	}

	configcat := r.Group("/configcat")
	{
		configcat.GET("/configuration-files/:key/config_v5.json", a.getConfigCatConfiguration)
	}

	v1 := r.Group("/v1")
	{
		v1Projects := v1.Group("/projects")
		{
			specificProject := v1Projects.Group("/:ref")
			{
				specificProject.GET("/custom-hostname", a.getProjectCustomHostname)
				specificProject.GET("/upgrade/eligibility", a.getProjectUpgradeEligibility)
			}
		}
	}

	return r
}
