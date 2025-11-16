```
# Resource Quotas & Limits System

## Overview

SupaManager provides configurable resource quotas to manage multi-tenant environments. Admins can set limits per project, enforce usage policies, and prevent resource abuse.

---

## Quota Types

### 1. Storage Quotas

| Quota | Description | Example |
|-------|-------------|---------|
| **Database Size** | PostgreSQL database size | 2 GB |
| **Storage Size** | Uploaded files size | 5 GB |
| **Backup Size** | Total backup storage | 10 GB |
| **Total Disk** | Combined disk usage | 15 GB |

### 2. Compute Quotas

| Quota | Description | Example |
|-------|-------------|---------|
| **CPU Limit** | Max CPU cores | 1.0 core |
| **Memory Limit** | Max RAM | 1 GB |

### 3. Network Quotas

| Quota | Description | Example |
|-------|-------------|---------|
| **Bandwidth** | Monthly data transfer | 50 GB/month |
| **Requests/Hour** | API requests per hour | 10,000 |
| **Connections** | Max concurrent DB connections | 50 |

### 4. Feature Limits

| Quota | Description | Example |
|-------|-------------|---------|
| **Max Backups** | Number of backups to keep | 7 |
| **Max Users** | Auth users (0 = unlimited) | 1,000 |
| **Max Tables** | Database tables | 200 |
| **Max File Size** | Individual file size limit | 50 MB |

---

## Quota Plans

### Default Plans

```go
// FREE PLAN
DatabaseSize:     500 MB
StorageSize:      1 GB
TotalDiskSize:    3 GB
CPULimit:         0.5 cores
MemoryLimit:      512 MB
BandwidthLimit:   10 GB/month
RequestsPerHour:  1,000
MaxBackups:       3
MaxUsers:         100
MaxFileSize:      10 MB

// STARTER PLAN
DatabaseSize:     2 GB
StorageSize:      5 GB
TotalDiskSize:    15 GB
CPULimit:         1.0 core
MemoryLimit:      1 GB
BandwidthLimit:   50 GB/month
RequestsPerHour:  10,000
MaxBackups:       7
MaxUsers:         1,000
MaxFileSize:      50 MB

// PRO PLAN
DatabaseSize:     10 GB
StorageSize:      50 GB
TotalDiskSize:    150 GB
CPULimit:         2.0 cores
MemoryLimit:      4 GB
BandwidthLimit:   500 GB/month
RequestsPerHour:  100,000
MaxBackups:       30
MaxUsers:         Unlimited
MaxFileSize:      500 MB

// ENTERPRISE PLAN
All quotas:       Unlimited
```

### Custom Plans

Admins can create custom plans with specific quotas.

---

## Admin Configuration

### Global Settings Page

```
┌─────────────────────────────────────────────────────┐
│ Admin Settings > Resource Quotas                    │
├─────────────────────────────────────────────────────┤
│                                                     │
│ Default Plan for New Projects                      │
│ ┌─────────────────────────────────────────────┐    │
│ │ [Free ▼] [Starter] [Pro] [Enterprise]      │    │
│ └─────────────────────────────────────────────┘    │
│                                                     │
│ Quota Enforcement                                  │
│ ┌─────────────────────────────────────────────┐    │
│ │ ☑ Enable quota enforcement                  │    │
│ │                                             │    │
│ │ Warning Threshold: [80] %                   │    │
│ │ Block Threshold:   [100] %                  │    │
│ │                                             │    │
│ │ When Quota Exceeded:                        │    │
│ │ ☑ Notify admin via email                    │    │
│ │ ☑ Notify user via email                     │    │
│ │ ☐ Pause project automatically               │    │
│ │ ☑ Block new uploads                         │    │
│ │ ☑ Block new backups                         │    │
│ │                                             │    │
│ │ Grace Period: [24] hours                    │    │
│ └─────────────────────────────────────────────┘    │
│                                                     │
│ Monitoring                                         │
│ ┌─────────────────────────────────────────────┐    │
│ │ Check quotas every: [1] hour                │    │
│ │                                             │    │
│ │ Notifications:                              │    │
│ │ ☑ Email on warning                          │    │
│ │ ☑ Email on exceeded                         │    │
│ │                                             │    │
│ │ Slack Webhook (optional):                   │    │
│ │ [https://hooks.slack.com/...]               │    │
│ └─────────────────────────────────────────────┘    │
│                                                     │
│ [Save Settings]                                    │
└─────────────────────────────────────────────────────┘
```

### Custom Plan Editor

```
┌─────────────────────────────────────────────────────┐
│ Create Custom Plan                                  │
├─────────────────────────────────────────────────────┤
│                                                     │
│ Plan Name: [__________________________]             │
│                                                     │
│ Storage Quotas                                     │
│ ┌─────────────────────────────────────────────┐    │
│ │ Database Size:    [2] GB                    │    │
│ │ Storage Size:     [5] GB                    │    │
│ │ Backup Size:      [10] GB                   │    │
│ │ Total Disk:       [15] GB                   │    │
│ └─────────────────────────────────────────────┘    │
│                                                     │
│ Compute Quotas                                     │
│ ┌─────────────────────────────────────────────┐    │
│ │ CPU Limit:        [1.0] cores               │    │
│ │ Memory Limit:     [1024] MB                 │    │
│ └─────────────────────────────────────────────┘    │
│                                                     │
│ Network Quotas                                     │
│ ┌─────────────────────────────────────────────┐    │
│ │ Bandwidth:        [50] GB/month             │    │
│ │ Requests/Hour:    [10000]                   │    │
│ │ Max Connections:  [50]                      │    │
│ └─────────────────────────────────────────────┘    │
│                                                     │
│ Feature Limits                                     │
│ ┌─────────────────────────────────────────────┐    │
│ │ Max Backups:      [7]                       │    │
│ │ Max Users:        [1000] (0 = unlimited)    │    │
│ │ Max Tables:       [200] (0 = unlimited)     │    │
│ │ Max File Size:    [50] MB                   │    │
│ └─────────────────────────────────────────────┘    │
│                                                     │
│ [Cancel] [Save Plan]                               │
└─────────────────────────────────────────────────────┘
```

---

## Per-Project Quota Override

Admins can override quotas for specific projects:

```
┌─────────────────────────────────────────────────────┐
│ Project Settings > Resource Quotas                  │
├─────────────────────────────────────────────────────┤
│                                                     │
│ Current Plan: [Starter ▼]                          │
│                                                     │
│ Override Quotas (leave blank to use plan defaults) │
│ ┌─────────────────────────────────────────────┐    │
│ │ Database Size:    [_____] GB                │    │
│ │ Storage Size:     [_____] GB                │    │
│ │ Total Disk:       [_____] GB                │    │
│ │ CPU Limit:        [_____] cores             │    │
│ │ Memory Limit:     [_____] MB                │    │
│ │ Bandwidth:        [_____] GB/month          │    │
│ │ Max Backups:      [_____]                   │    │
│ └─────────────────────────────────────────────┘    │
│                                                     │
│ [Reset to Plan Defaults] [Save Overrides]          │
└─────────────────────────────────────────────────────┘
```

---

## User View - Quota Dashboard

Users see their current usage and limits:

```
┌─────────────────────────────────────────────────────┐
│ Resource Usage - Project: my-app                    │
├─────────────────────────────────────────────────────┤
│                                                     │
│ Storage                                            │
│ ┌─────────────────────────────────────────────┐    │
│ │ Database:  1.2 GB / 2.0 GB  [████████░░] 60%│    │
│ │ Storage:   3.8 GB / 5.0 GB  [████████░░] 76%│    │
│ │ Backups:   5.2 GB / 10 GB   [█████░░░░░] 52%│    │
│ │ ─────────────────────────────────────────   │    │
│ │ Total:     10.2 GB / 15 GB  [███████░░░] 68%│    │
│ └─────────────────────────────────────────────┘    │
│                                                     │
│ Compute                                            │
│ ┌─────────────────────────────────────────────┐    │
│ │ CPU:       0.65 / 1.0 cores [████████░░] 65%│    │
│ │ Memory:    780 MB / 1 GB    [████████░░] 76%│    │
│ └─────────────────────────────────────────────┘    │
│                                                     │
│ Network (This Month)                               │
│ ┌─────────────────────────────────────────────┐    │
│ │ Bandwidth: 32 GB / 50 GB    [██████░░░░] 64%│    │
│ │ Requests:  87.5K / 10K/hr   [████████░░] 88%│    │
│ │                            ⚠ Near limit      │    │
│ └─────────────────────────────────────────────┘    │
│                                                     │
│ Features                                           │
│ ┌─────────────────────────────────────────────┐    │
│ │ Backups:   5 / 7            [███████░░░] 71%│    │
│ │ Users:     487 / 1,000      [████░░░░░░] 49%│    │
│ │ Tables:    78 / 200         [████░░░░░░] 39%│    │
│ └─────────────────────────────────────────────┘    │
│                                                     │
│ [View Detailed Usage] [Upgrade Plan]               │
└─────────────────────────────────────────────────────┘
```

---

## Quota Enforcement Flow

### 1. Soft Limit (Warning - 80%)

```
User uploads file
      ↓
Check quota: 81% used
      ↓
⚠ Warning shown:
"Storage is 81% full (4.05 GB / 5 GB)"
"Consider upgrading your plan"
      ↓
Upload proceeds
      ↓
Email sent to user
```

### 2. Hard Limit (Block - 100%)

```
User uploads file
      ↓
Check quota: 100% used
      ↓
❌ Upload blocked:
"Storage limit exceeded (5.0 GB / 5.0 GB)"
"Please upgrade your plan or delete files"
      ↓
Upload rejected with 413 error
      ↓
Email sent to user & admin
```

### 3. Grace Period

```
Quota exceeded
      ↓
Grace period starts (24 hours)
      ↓
Warnings shown, but operations allowed
      ↓
After 24 hours
      ↓
Hard limits enforced
      ↓
Operations blocked
```

---

## Database Schema

### Quotas Table

```sql
CREATE TABLE project_quotas (
    project_id TEXT PRIMARY KEY,
    plan QuotaPlan NOT NULL DEFAULT 'FREE',

    -- Storage quotas (bytes)
    database_size BIGINT,
    storage_size BIGINT,
    backup_size BIGINT,
    total_disk_size BIGINT,

    -- Compute quotas
    cpu_limit FLOAT,
    memory_limit BIGINT,

    -- Network quotas
    bandwidth_limit BIGINT,
    requests_per_hour BIGINT,
    connections_limit INTEGER,

    -- Feature limits
    max_backups INTEGER,
    max_users INTEGER,
    max_tables INTEGER,
    max_file_size BIGINT,

    -- Metadata
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    updated_by TEXT
);
```

### Usage Tracking Table

```sql
CREATE TABLE project_usage (
    project_id TEXT PRIMARY KEY,

    -- Current usage (bytes)
    database_size BIGINT DEFAULT 0,
    storage_size BIGINT DEFAULT 0,
    backup_size BIGINT DEFAULT 0,
    total_disk_size BIGINT DEFAULT 0,

    -- Compute usage
    cpu_usage FLOAT DEFAULT 0,
    memory_usage BIGINT DEFAULT 0,

    -- Network usage
    bandwidth_used BIGINT DEFAULT 0,
    bandwidth_reset_at TIMESTAMPTZ,
    requests_this_hour BIGINT DEFAULT 0,
    requests_hour_start TIMESTAMPTZ,
    active_connections INTEGER DEFAULT 0,

    -- Feature usage
    backup_count INTEGER DEFAULT 0,
    user_count INTEGER DEFAULT 0,
    table_count INTEGER DEFAULT 0,

    -- Metadata
    last_updated TIMESTAMPTZ DEFAULT NOW()
);

-- Index for quick lookups
CREATE INDEX idx_usage_updated ON project_usage(last_updated);
```

### Quota Violations Log

```sql
CREATE TABLE quota_violations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id TEXT NOT NULL,
    quota_type TEXT NOT NULL, -- 'database', 'storage', 'bandwidth', etc.
    limit_value BIGINT NOT NULL,
    current_value BIGINT NOT NULL,
    exceeded_by BIGINT NOT NULL,
    violation_time TIMESTAMPTZ DEFAULT NOW(),
    notified BOOLEAN DEFAULT FALSE,
    resolved_at TIMESTAMPTZ
);

-- Index for querying violations
CREATE INDEX idx_violations_project ON quota_violations(project_id);
CREATE INDEX idx_violations_time ON quota_violations(violation_time);
```

---

## API Endpoints

### Get Project Quotas
```http
GET /platform/projects/{projectID}/quotas
Authorization: Bearer <token>
```

**Response:**
```json
{
  "plan": "STARTER",
  "quotas": {
    "database_size": 2147483648,
    "storage_size": 5368709120,
    "total_disk_size": 16106127360,
    "cpu_limit": 1.0,
    "memory_limit": 1073741824,
    "bandwidth_limit": 53687091200,
    "max_backups": 7
  }
}
```

### Get Current Usage
```http
GET /platform/projects/{projectID}/usage
Authorization: Bearer <token>
```

**Response:**
```json
{
  "usage": {
    "database_size": 1288490189,
    "storage_size": 4080218931,
    "total_disk_size": 10956301107,
    "cpu_usage": 0.65,
    "memory_usage": 817889689,
    "bandwidth_used": 34359738368,
    "backup_count": 5,
    "user_count": 487,
    "table_count": 78
  },
  "status": {
    "exceeded": false,
    "warnings": [
      "Storage is 76% full (3.8 GB / 5 GB)",
      "Requests are 88% of hourly limit (8,800 / 10,000)"
    ]
  }
}
```

### Update Project Quotas (Admin Only)
```http
PUT /admin/projects/{projectID}/quotas
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "plan": "PRO",
  "overrides": {
    "database_size": 21474836480,
    "storage_size": 53687091200
  }
}
```

### Get Global Quota Settings (Admin Only)
```http
GET /admin/settings/quotas
Authorization: Bearer <admin-token>
```

### Update Global Quota Settings (Admin Only)
```http
PUT /admin/settings/quotas
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "default_plan": "FREE",
  "enforcement": {
    "warn_at_percent": 80,
    "block_at_percent": 100,
    "notify_admin": true,
    "notify_user": true,
    "grace_period_hours": 24
  }
}
```

---

## Quota Monitoring

### Background Job

Runs every hour to:
1. Calculate current usage for all projects
2. Update `project_usage` table
3. Check if quotas exceeded
4. Send notifications
5. Enforce hard limits

```go
func MonitorQuotas(ctx context.Context) {
    projects := GetAllProjects()

    for _, project := range projects {
        // 1. Calculate current usage
        usage := CalculateUsage(project.ID)

        // 2. Get quotas
        quotas := GetProjectQuotas(project.ID)

        // 3. Check limits
        status := CheckQuotas(usage, quotas)

        // 4. Take action
        if status.Exceeded {
            HandleQuotaViolation(project, status)
        } else if status.Warning {
            SendWarningNotification(project, status)
        }

        // 5. Update database
        SaveUsageStats(project.ID, usage)
    }
}
```

---

## Role-Based Access Control

### Roles

1. **Admin**
   - Configure global quota settings
   - Override quotas for any project
   - View all project usage
   - Enforce/disable quotas

2. **Organization Owner**
   - View quota usage for org projects
   - Request quota increases
   - Set quotas within org limits

3. **Project Owner**
   - View quota usage for their projects
   - Request quota increases
   - Delete old backups to free space

4. **Project Member**
   - View quota usage (read-only)

### Permissions Check

```go
func CanModifyQuotas(user *User, projectID string) bool {
    if user.IsAdmin() {
        return true
    }

    if user.IsOrgOwner(projectID) {
        return true
    }

    return false
}
```

---

## User Experience

### When Quota Warning (80%)

```
┌─────────────────────────────────────────┐
│ ⚠ Storage Usage Warning                 │
├─────────────────────────────────────────┤
│ Your project is using 4.0 GB of 5.0 GB │
│ storage (80% full).                     │
│                                         │
│ Consider:                               │
│ • Upgrading to Pro plan (50 GB)         │
│ • Deleting old backups                  │
│ • Removing unused files                 │
│                                         │
│ [Upgrade Plan] [Manage Storage]         │
└─────────────────────────────────────────┘
```

### When Quota Exceeded (100%)

```
┌─────────────────────────────────────────┐
│ ❌ Storage Limit Exceeded               │
├─────────────────────────────────────────┤
│ Your project has reached its storage   │
│ limit of 5.0 GB.                        │
│                                         │
│ New uploads are currently blocked.     │
│                                         │
│ To continue:                            │
│ • Upgrade to a larger plan              │
│ • Delete files to free up space         │
│                                         │
│ [Upgrade Now] [Manage Files]            │
└─────────────────────────────────────────┘
```

---

## Implementation Phases

### Phase 3: Basic Quotas
- ✅ Quota interface design (done)
- Database schema
- Storage quota enforcement
- Basic monitoring

### Phase 4: Advanced Quotas
- Compute quotas (CPU/Memory)
- Network quotas (Bandwidth/Requests)
- Feature limits (Users/Tables)
- Grace periods

### Phase 5: Complete System
- Plan management UI
- Admin configuration
- Detailed usage analytics
- Automated notifications
- Quota violation logs

---

## Future Enhancements

- **Pay-as-you-go**: Automatic billing for overages
- **Burst capacity**: Temporary quota increases
- **Predictive alerts**: "At current rate, will exceed in 7 days"
- **Usage analytics**: Charts and trends
- **Resource optimization**: Suggestions to reduce usage
- **Cost calculator**: Estimate costs based on usage

---

**Status**: Design Complete ✅
**Implementation**: Phase 3-5 🔨
```
