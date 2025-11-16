package provisioner

import (
	"context"
	"time"
)

// BackupType defines the type of backup
type BackupType string

const (
	BackupTypeFull        BackupType = "FULL"         // Database + Storage + Config
	BackupTypeDatabase    BackupType = "DATABASE"     // PostgreSQL only
	BackupTypeStorage     BackupType = "STORAGE"      // Storage files only
	BackupTypeIncremental BackupType = "INCREMENTAL"  // Changes since last backup
)

// BackupConfig contains backup configuration
type BackupConfig struct {
	ProjectID      string
	BackupType     BackupType
	Compression    bool          // Compress backup (gzip)
	Encryption     bool          // Encrypt backup
	Retention      int           // Keep N backups (older deleted)
	AutoCleanup    bool          // Automatically delete old backups
	S3Upload       bool          // Upload to S3
	S3Bucket       string        // S3 bucket name
	S3Prefix       string        // S3 key prefix
}

// BackupInfo contains metadata about a backup
type BackupInfo struct {
	BackupID       string        // Unique backup identifier
	ProjectID      string
	ProjectName    string
	BackupType     BackupType
	Size           int64         // Backup size in bytes
	Compressed     bool
	Encrypted      bool
	FilePath       string        // Local file path
	S3Key          string        // S3 key (if uploaded)
	Status         string        // CREATING, COMPLETED, FAILED
	ErrorMessage   string
	CreatedAt      time.Time
	CompletedAt    time.Time
	ExpiresAt      time.Time     // Auto-deletion date
}

// RestoreConfig contains restore configuration
type RestoreConfig struct {
	ProjectID      string
	BackupID       string
	RestoreType    BackupType    // What to restore
	OverwriteData  bool          // Overwrite existing data
	PointInTime    *time.Time    // Restore to specific point (if available)
	StopProject    bool          // Stop project before restore
}

// RestoreInfo contains metadata about a restore operation
type RestoreInfo struct {
	RestoreID      string
	ProjectID      string
	BackupID       string
	Status         string        // RESTORING, COMPLETED, FAILED
	Progress       float64       // 0-100%
	ErrorMessage   string
	StartedAt      time.Time
	CompletedAt    time.Time
}

// BackupSchedule defines automated backup schedule
type BackupSchedule struct {
	ProjectID      string
	Enabled        bool
	Frequency      string        // "hourly", "daily", "weekly"
	Time           string        // "02:00" (UTC)
	BackupType     BackupType
	Retention      int           // Keep N backups
}

// BackupProvisioner extends Provisioner with backup/restore capabilities
type BackupProvisioner interface {
	Provisioner

	// CreateBackup creates a new backup of a project
	// Returns backup info with status CREATING, use GetBackupInfo to check completion
	CreateBackup(ctx context.Context, config *BackupConfig) (*BackupInfo, error)

	// GetBackupInfo retrieves information about a backup
	GetBackupInfo(ctx context.Context, backupID string) (*BackupInfo, error)

	// ListBackups returns all backups for a project
	ListBackups(ctx context.Context, projectID string) ([]*BackupInfo, error)

	// DeleteBackup removes a backup file
	DeleteBackup(ctx context.Context, backupID string) error

	// RestoreBackup restores a project from a backup
	// Returns restore info with status RESTORING, use GetRestoreInfo to check progress
	RestoreBackup(ctx context.Context, config *RestoreConfig) (*RestoreInfo, error)

	// GetRestoreInfo retrieves information about a restore operation
	GetRestoreInfo(ctx context.Context, restoreID string) (*RestoreInfo, error)

	// DownloadBackup prepares a backup for download and returns a download URL
	DownloadBackup(ctx context.Context, backupID string, expiresIn time.Duration) (string, error)

	// SetBackupSchedule configures automated backups for a project
	SetBackupSchedule(ctx context.Context, schedule *BackupSchedule) error

	// GetBackupSchedule retrieves the backup schedule for a project
	GetBackupSchedule(ctx context.Context, projectID string) (*BackupSchedule, error)

	// ExportProject exports a complete project (for migration)
	// Includes database, storage, config, and docker-compose files
	ExportProject(ctx context.Context, projectID string, outputPath string) error

	// ImportProject imports a project from an export
	// Creates a new project from exported files
	ImportProject(ctx context.Context, exportPath string, newProjectID string) error
}

// BackupStorage defines where backups are stored
type BackupStorage interface {
	// Upload uploads a backup file
	Upload(ctx context.Context, filePath string, key string) error

	// Download downloads a backup file
	Download(ctx context.Context, key string, filePath string) error

	// Delete deletes a backup file
	Delete(ctx context.Context, key string) error

	// List lists all backups
	List(ctx context.Context, prefix string) ([]string, error)

	// GetDownloadURL generates a temporary download URL
	GetDownloadURL(ctx context.Context, key string, expiresIn time.Duration) (string, error)
}

// LocalBackupStorage implements BackupStorage using local filesystem
type LocalBackupStorage struct {
	baseDir string // /var/lib/supamanager/backups/
}

// S3BackupStorage implements BackupStorage using S3
type S3BackupStorage struct {
	bucket string
	region string
	// S3 client will be added in Phase 5
}
