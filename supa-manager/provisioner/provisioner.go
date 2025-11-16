package provisioner

import (
	"context"
)

// ProjectStatus represents the current state of a provisioned project
type ProjectStatus string

const (
	StatusCreating      ProjectStatus = "CREATING"
	StatusActive        ProjectStatus = "ACTIVE_HEALTHY"
	StatusUnhealthy     ProjectStatus = "ACTIVE_UNHEALTHY"
	StatusPaused        ProjectStatus = "PAUSED"
	StatusDeleting      ProjectStatus = "DELETING"
	StatusFailed        ProjectStatus = "FAILED"
)

// ProjectConfig contains all configuration needed to provision a new project
type ProjectConfig struct {
	ProjectID     string // Unique identifier (e.g., "abc123")
	ProjectName   string // Human-readable name
	OrganizationID string // Owner organization
	Region        string // Deployment region (for future multi-region support)

	// Database configuration
	DBPassword    string // PostgreSQL password
	DBPort        int    // Exposed PostgreSQL port

	// API configuration
	APIPort       int    // Kong API Gateway port
	StudioPort    int    // Studio UI port (optional, for per-project studio)

	// Security
	JWTSecret     string // JWT signing secret
	AnonKey       string // Anonymous API key
	ServiceKey    string // Service role API key

	// Dashboard
	DashboardUser string // Dashboard username
	DashboardPass string // Dashboard password

	// Resource limits (for future use)
	CPULimit      string // CPU limit (e.g., "1.0")
	MemoryLimit   string // Memory limit (e.g., "2GB")
	StorageLimit  string // Storage limit (e.g., "10GB")
}

// ProjectInfo contains runtime information about a provisioned project
type ProjectInfo struct {
	ProjectID     string
	ProjectName   string
	Status        ProjectStatus
	Endpoint      string // API endpoint URL
	DBEndpoint    string // Database connection string

	// Container IDs for management
	Containers    map[string]string // service name -> container ID

	// Health information
	HealthChecks  map[string]bool // service name -> healthy status

	// Resource usage (for monitoring)
	CPUUsage      float64
	MemoryUsage   uint64
	StorageUsage  uint64

	CreatedAt     string
	UpdatedAt     string
}

// Provisioner defines the interface for provisioning Supabase projects
// This interface allows for different implementations (Docker, K8s, etc.)
type Provisioner interface {
	// CreateProject provisions a new Supabase instance
	// Returns the project info or an error if provisioning fails
	CreateProject(ctx context.Context, config *ProjectConfig) (*ProjectInfo, error)

	// GetProjectInfo retrieves current information about a project
	GetProjectInfo(ctx context.Context, projectID string) (*ProjectInfo, error)

	// UpdateProject updates project configuration (e.g., scale resources)
	UpdateProject(ctx context.Context, projectID string, config *ProjectConfig) error

	// PauseProject stops all containers for a project without deleting data
	PauseProject(ctx context.Context, projectID string) error

	// ResumeProject restarts all containers for a paused project
	ResumeProject(ctx context.Context, projectID string) error

	// DeleteProject removes all containers and volumes for a project
	// This is a destructive operation and cannot be undone
	DeleteProject(ctx context.Context, projectID string) error

	// ListProjects returns information about all provisioned projects
	ListProjects(ctx context.Context) ([]*ProjectInfo, error)

	// GetLogs retrieves logs for a specific service in a project
	GetLogs(ctx context.Context, projectID string, service string, tail int) ([]string, error)

	// ExecuteCommand runs a command in a specific service container
	// Useful for database migrations, backups, etc.
	ExecuteCommand(ctx context.Context, projectID string, service string, cmd []string) (string, error)
}

// ProvisionerError represents an error that occurred during provisioning
type ProvisionerError struct {
	ProjectID string
	Operation string
	Err       error
}

func (e *ProvisionerError) Error() string {
	return "provisioning error for project " + e.ProjectID + " during " + e.Operation + ": " + e.Err.Error()
}
