package provisioner

import (
	"context"
	"time"
)

// ResourceQuotas defines resource limits for a project
type ResourceQuotas struct {
	// Storage quotas
	DatabaseSize     int64 // PostgreSQL database size in bytes
	StorageSize      int64 // File storage size in bytes
	BackupSize       int64 // Total backup storage in bytes
	TotalDiskSize    int64 // Total disk usage (DB + Storage + Backups)

	// Compute quotas
	CPULimit         float64 // CPU cores (e.g., 1.0, 2.0)
	MemoryLimit      int64   // RAM in bytes

	// Network quotas
	BandwidthLimit   int64   // Monthly bandwidth in bytes
	RequestsPerHour  int64   // API requests per hour
	ConnectionsLimit int     // Max concurrent DB connections

	// Feature limits
	MaxBackups       int     // Maximum number of backups to keep
	MaxUsers         int     // Maximum auth users (0 = unlimited)
	MaxTables        int     // Maximum database tables (0 = unlimited)
	MaxFileSize      int64   // Maximum individual file size
}

// QuotaUsage tracks current resource usage
type QuotaUsage struct {
	ProjectID        string
	LastUpdated      time.Time

	// Current usage
	DatabaseSize     int64
	StorageSize      int64
	BackupSize       int64
	TotalDiskSize    int64

	CPUUsage         float64   // Current CPU usage
	MemoryUsage      int64     // Current memory usage

	BandwidthUsed    int64     // Bandwidth used this month
	RequestsThisHour int64     // Requests in current hour
	ActiveConnections int      // Current DB connections

	BackupCount      int
	UserCount        int
	TableCount       int
}

// QuotaStatus indicates if quotas are being exceeded
type QuotaStatus struct {
	Exceeded         bool
	Warnings         []string
	Errors           []string

	// Individual quota checks
	DatabaseQuota    QuotaCheck
	StorageQuota     QuotaCheck
	BackupQuota      QuotaCheck
	BandwidthQuota   QuotaCheck
	CPUQuota         QuotaCheck
	MemoryQuota      QuotaCheck
}

// QuotaCheck represents the status of a single quota
type QuotaCheck struct {
	Current     int64
	Limit       int64
	Used        float64   // Percentage (0-100)
	Exceeded    bool
	Warning     bool      // Above 80%
}

// QuotaPlan defines different quota tiers
type QuotaPlan string

const (
	PlanFree       QuotaPlan = "FREE"
	PlanStarter    QuotaPlan = "STARTER"
	PlanPro        QuotaPlan = "PRO"
	PlanEnterprise QuotaPlan = "ENTERPRISE"
	PlanCustom     QuotaPlan = "CUSTOM"
)

// GetDefaultQuotas returns default quotas for a plan
func GetDefaultQuotas(plan QuotaPlan) ResourceQuotas {
	switch plan {
	case PlanFree:
		return ResourceQuotas{
			DatabaseSize:     500 * 1024 * 1024,      // 500 MB
			StorageSize:      1 * 1024 * 1024 * 1024, // 1 GB
			BackupSize:       2 * 1024 * 1024 * 1024, // 2 GB
			TotalDiskSize:    3 * 1024 * 1024 * 1024, // 3 GB
			CPULimit:         0.5,
			MemoryLimit:      512 * 1024 * 1024,      // 512 MB
			BandwidthLimit:   10 * 1024 * 1024 * 1024, // 10 GB/month
			RequestsPerHour:  1000,
			ConnectionsLimit: 10,
			MaxBackups:       3,
			MaxUsers:         100,
			MaxTables:        50,
			MaxFileSize:      10 * 1024 * 1024,       // 10 MB
		}

	case PlanStarter:
		return ResourceQuotas{
			DatabaseSize:     2 * 1024 * 1024 * 1024,   // 2 GB
			StorageSize:      5 * 1024 * 1024 * 1024,   // 5 GB
			BackupSize:       10 * 1024 * 1024 * 1024,  // 10 GB
			TotalDiskSize:    15 * 1024 * 1024 * 1024,  // 15 GB
			CPULimit:         1.0,
			MemoryLimit:      1024 * 1024 * 1024,       // 1 GB
			BandwidthLimit:   50 * 1024 * 1024 * 1024,  // 50 GB/month
			RequestsPerHour:  10000,
			ConnectionsLimit: 50,
			MaxBackups:       7,
			MaxUsers:         1000,
			MaxTables:        200,
			MaxFileSize:      50 * 1024 * 1024,         // 50 MB
		}

	case PlanPro:
		return ResourceQuotas{
			DatabaseSize:     10 * 1024 * 1024 * 1024,  // 10 GB
			StorageSize:      50 * 1024 * 1024 * 1024,  // 50 GB
			BackupSize:       100 * 1024 * 1024 * 1024, // 100 GB
			TotalDiskSize:    150 * 1024 * 1024 * 1024, // 150 GB
			CPULimit:         2.0,
			MemoryLimit:      4 * 1024 * 1024 * 1024,   // 4 GB
			BandwidthLimit:   500 * 1024 * 1024 * 1024, // 500 GB/month
			RequestsPerHour:  100000,
			ConnectionsLimit: 200,
			MaxBackups:       30,
			MaxUsers:         0, // Unlimited
			MaxTables:        0, // Unlimited
			MaxFileSize:      500 * 1024 * 1024,        // 500 MB
		}

	case PlanEnterprise:
		return ResourceQuotas{
			DatabaseSize:     0, // Unlimited
			StorageSize:      0, // Unlimited
			BackupSize:       0, // Unlimited
			TotalDiskSize:    0, // Unlimited
			CPULimit:         0, // Unlimited
			MemoryLimit:      0, // Unlimited
			BandwidthLimit:   0, // Unlimited
			RequestsPerHour:  0, // Unlimited
			ConnectionsLimit: 0, // Unlimited
			MaxBackups:       0, // Unlimited
			MaxUsers:         0, // Unlimited
			MaxTables:        0, // Unlimited
			MaxFileSize:      0, // Unlimited
		}

	default:
		return GetDefaultQuotas(PlanFree)
	}
}

// QuotaProvisioner extends Provisioner with quota management
type QuotaProvisioner interface {
	Provisioner

	// SetProjectQuotas sets resource quotas for a project
	SetProjectQuotas(ctx context.Context, projectID string, quotas ResourceQuotas) error

	// GetProjectQuotas retrieves current quotas for a project
	GetProjectQuotas(ctx context.Context, projectID string) (*ResourceQuotas, error)

	// GetQuotaUsage retrieves current resource usage
	GetQuotaUsage(ctx context.Context, projectID string) (*QuotaUsage, error)

	// GetQuotaStatus checks if project is within quotas
	GetQuotaStatus(ctx context.Context, projectID string) (*QuotaStatus, error)

	// EnforceQuotas checks and enforces quotas (called before operations)
	// Returns error if operation would exceed quotas
	EnforceQuotas(ctx context.Context, projectID string, operation string, size int64) error

	// UpdateQuotaUsage recalculates current usage
	// Called periodically by background job
	UpdateQuotaUsage(ctx context.Context, projectID string) error
}

// QuotaEnforcement defines what happens when quotas are exceeded
type QuotaEnforcement struct {
	// Soft limits - warn but allow
	WarnAtPercent    float64 // Warn when usage exceeds this % (e.g., 80)

	// Hard limits - block operations
	BlockAtPercent   float64 // Block when usage exceeds this % (e.g., 100)

	// Actions when exceeded
	NotifyAdmin      bool    // Email admin when quota exceeded
	NotifyUser       bool    // Email user when quota exceeded
	PauseProject     bool    // Pause project when quota exceeded
	BlockUploads     bool    // Block file uploads when storage quota exceeded
	BlockBackups     bool    // Block new backups when backup quota exceeded
	BlockNewUsers    bool    // Block new user signups when user quota exceeded
}

// AdminQuotaSettings defines global quota configuration
type AdminQuotaSettings struct {
	// Default quotas for new projects
	DefaultPlan      QuotaPlan

	// Enforcement settings
	Enforcement      QuotaEnforcement

	// Monitoring
	CheckInterval    time.Duration // How often to check quotas

	// Grace period
	GracePeriod      time.Duration // Time before enforcing hard limits

	// Notifications
	EmailOnWarning   bool
	EmailOnExceeded  bool
	SlackWebhook     string // Optional Slack notifications
}
