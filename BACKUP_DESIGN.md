# Backup & Restore System Design

## Overview

SupaManager provides comprehensive backup and restore capabilities for all Supabase projects. Users can create manual or automated backups through the web UI with zero manual work.

---

## Backup Types

### 1. Full Backup (Recommended)
```
✓ PostgreSQL database (pg_dump)
✓ Storage files (all uploaded files)
✓ Project configuration (.env, configs)
✓ Docker compose files
```

**Use case**: Complete project snapshot for disaster recovery

### 2. Database Only
```
✓ PostgreSQL database (pg_dump)
✗ Storage files
✗ Configuration
```

**Use case**: Quick database snapshot before schema changes

### 3. Storage Only
```
✗ PostgreSQL database
✓ Storage files
✗ Configuration
```

**Use case**: Backup uploaded files separately

### 4. Incremental
```
✓ Changes since last full backup
✓ Much faster and smaller
```

**Use case**: Frequent automated backups

---

## Directory Structure

```
/var/lib/supamanager/
├── projects/
│   └── abc123/
│       ├── docker-compose.yml
│       └── volumes/
└── backups/
    └── abc123/
        ├── backup_20250315_020000_full.tar.gz
        ├── backup_20250316_020000_incremental.tar.gz
        └── backup_20250317_020000_full.tar.gz
```

---

## User Flow - Create Backup

### Via Web UI

1. **Navigate** to project settings
2. **Click** "Create Backup"
3. **Select** backup type:
   - Full Backup (Database + Storage + Config)
   - Database Only
   - Storage Only
4. **Click** "Create"
5. **Progress** shown in real-time (0-100%)
6. **Done!** Download button appears when complete

### Via API

```http
POST /platform/projects/{projectID}/backups
Authorization: Bearer <token>
Content-Type: application/json

{
  "type": "FULL",
  "compression": true,
  "encryption": true
}
```

**Response:**
```json
{
  "backup_id": "backup_abc123_20250315_020000",
  "status": "CREATING",
  "progress": 0,
  "estimated_time": "5 minutes"
}
```

---

## How Backups Work

### Database Backup Process

```bash
# 1. Connect to PostgreSQL container
docker exec abc123-db pg_dump -U postgres -Fc postgres > backup.dump

# 2. Compress (optional)
gzip backup.dump

# 3. Encrypt (optional)
openssl enc -aes-256-cbc -salt -in backup.dump.gz -out backup.dump.gz.enc

# 4. Store locally or upload to S3
mv backup.dump.gz.enc /var/lib/supamanager/backups/abc123/
```

### Storage Backup Process

```bash
# 1. Copy storage volume data
docker run --rm \
  -v abc123_storage:/source \
  -v /var/lib/supamanager/backups/abc123:/backup \
  alpine tar czf /backup/storage_backup.tar.gz -C /source .
```

### Full Backup Process

```bash
# 1. Database backup
pg_dump → backup_db.dump

# 2. Storage backup
tar czf storage_backup.tar.gz volumes/storage/

# 3. Configuration backup
cp .env docker-compose.yml kong.yml vector.yml → config/

# 4. Combine all into single archive
tar czf backup_20250315_020000_full.tar.gz db/ storage/ config/
```

---

## Automated Backups

### Schedule Configuration

Users can configure automated backups:

```json
{
  "enabled": true,
  "frequency": "daily",
  "time": "02:00",
  "type": "FULL",
  "retention": 7
}
```

**Frequencies:**
- **Hourly**: Every hour
- **Daily**: Once per day at specified time (UTC)
- **Weekly**: Once per week on Sunday

**Retention:** Keep N most recent backups, delete older ones

### Backup Schedule UI

```
┌─────────────────────────────────────────┐
│ Automated Backups                       │
├─────────────────────────────────────────┤
│ ☑ Enable automated backups              │
│                                         │
│ Frequency: [Daily ▼]                    │
│ Time (UTC): [02:00 ▼]                   │
│ Backup Type: [Full Backup ▼]            │
│ Retention: [7 ▼] backups                │
│                                         │
│ Next backup: March 16, 2025 02:00 UTC  │
│                                         │
│ [Save Schedule]                         │
└─────────────────────────────────────────┘
```

---

## Restore Process

### Via Web UI

1. **Navigate** to Backups page
2. **See list** of available backups
3. **Click** "Restore" on desired backup
4. **Warning**: "This will overwrite current data. Continue?"
5. **Click** "Yes, Restore"
6. **Project paused** automatically
7. **Restore runs** (5-15 minutes)
8. **Project resumed** when complete
9. **Done!** Project restored to backup point

### Restore Options

```
┌─────────────────────────────────────────┐
│ Restore Backup                          │
├─────────────────────────────────────────┤
│ Backup: backup_20250315_020000_full    │
│ Created: March 15, 2025 02:00 UTC      │
│ Size: 1.2 GB                            │
│                                         │
│ Restore Options:                        │
│ ☑ Database                              │
│ ☑ Storage Files                         │
│ ☑ Configuration                         │
│                                         │
│ ⚠ Warning: This will overwrite         │
│   current project data                  │
│                                         │
│ [Cancel] [Restore Project]              │
└─────────────────────────────────────────┘
```

---

## Backup Storage Options

### Local Storage (Default)

```
Location: /var/lib/supamanager/backups/
Pros:
  - Fast access
  - No additional costs
  - Simple setup
Cons:
  - Limited by disk space
  - No offsite redundancy
```

### S3 Storage (Optional)

```
Provider: AWS S3, MinIO, Wasabi, etc.
Pros:
  - Unlimited storage
  - Offsite redundancy
  - Automatic replication
  - Lower cost for cold storage
Cons:
  - Slower access
  - Transfer costs
  - Requires configuration
```

**S3 Configuration:**
```yaml
backup_storage:
  type: s3
  bucket: supamanager-backups
  region: us-east-1
  prefix: production/
  storage_class: STANDARD_IA  # Infrequent Access
```

---

## Download Backups

Users can download backups to their local machine:

### Via Web UI

1. Click **"Download"** next to backup
2. Generates temporary download URL
3. File downloads to user's computer
4. URL expires after 1 hour

### Via API

```http
GET /platform/projects/{projectID}/backups/{backupID}/download
Authorization: Bearer <token>
```

**Response:**
```json
{
  "download_url": "https://backups.supamanager.io/abc123/backup.tar.gz?token=...",
  "expires_at": "2025-03-15T03:00:00Z",
  "size": 1234567890
}
```

---

## Import/Export Projects

### Export Project (Migration)

**Use case**: Moving project to another SupaManager instance

```http
POST /platform/projects/{projectID}/export
```

**Creates:**
```
project_abc123_export.tar.gz
├── database.dump
├── storage/
├── config/
│   ├── .env
│   ├── docker-compose.yml
│   ├── kong.yml
│   └── vector.yml
└── metadata.json
```

### Import Project

```http
POST /platform/projects/import
Content-Type: multipart/form-data

{
  "file": project_abc123_export.tar.gz,
  "new_project_name": "My Imported App"
}
```

**Creates new project with all data from export**

---

## Backup Security

### Encryption

All backups are encrypted using AES-256:

```bash
openssl enc -aes-256-cbc -salt \
  -in backup.tar.gz \
  -out backup.tar.gz.enc \
  -pass pass:$ENCRYPTION_KEY
```

**Encryption key**: Derived from `ENCRYPTION_SECRET` in environment

### Access Control

- Only project owners can create/restore backups
- Organization admins can access all org project backups
- Service role can create automated backups

### Backup Verification

Each backup includes checksums:

```json
{
  "backup_id": "backup_abc123_20250315",
  "checksums": {
    "database": "sha256:abc123...",
    "storage": "sha256:def456...",
    "config": "sha256:ghi789..."
  }
}
```

Checksums verified during restore to ensure integrity.

---

## API Endpoints

### Create Backup
```http
POST /platform/projects/{projectID}/backups
```

### List Backups
```http
GET /platform/projects/{projectID}/backups
```

### Get Backup Info
```http
GET /platform/projects/{projectID}/backups/{backupID}
```

### Delete Backup
```http
DELETE /platform/projects/{projectID}/backups/{backupID}
```

### Restore Backup
```http
POST /platform/projects/{projectID}/backups/{backupID}/restore
```

### Download Backup
```http
GET /platform/projects/{projectID}/backups/{backupID}/download
```

### Set Backup Schedule
```http
PUT /platform/projects/{projectID}/backup-schedule
```

### Export Project
```http
POST /platform/projects/{projectID}/export
```

### Import Project
```http
POST /platform/projects/import
```

---

## Backup Management UI

### Backups Page

```
┌────────────────────────────────────────────────────────┐
│ Backups - Project: my-app (abc123)                    │
├────────────────────────────────────────────────────────┤
│ [Create Backup] [Configure Schedule]                   │
│                                                        │
│ ┌──────────────────────────────────────────────────┐  │
│ │ backup_20250317_020000_full           1.2 GB    │  │
│ │ Created: Mar 17, 2025 02:00 UTC                 │  │
│ │ Type: Full | Status: ✓ Completed                │  │
│ │ [Restore] [Download] [Delete]                   │  │
│ └──────────────────────────────────────────────────┘  │
│                                                        │
│ ┌──────────────────────────────────────────────────┐  │
│ │ backup_20250316_020000_full           1.1 GB    │  │
│ │ Created: Mar 16, 2025 02:00 UTC                 │  │
│ │ Type: Full | Status: ✓ Completed                │  │
│ │ [Restore] [Download] [Delete]                   │  │
│ └──────────────────────────────────────────────────┘  │
│                                                        │
│ ┌──────────────────────────────────────────────────┐  │
│ │ backup_20250315_020000_full           1.0 GB    │  │
│ │ Created: Mar 15, 2025 02:00 UTC                 │  │
│ │ Type: Full | Status: ✓ Completed                │  │
│ │ [Restore] [Download] [Delete]                   │  │
│ └──────────────────────────────────────────────────┘  │
│                                                        │
│ Total: 3 backups | 3.3 GB                             │
└────────────────────────────────────────────────────────┘
```

---

## Implementation Phases

### Phase 3: Basic Backups
- ✅ Backup interface design (done)
- Manual database backups
- Manual restore
- Local storage only

### Phase 4: Automated Backups
- Backup scheduler
- Automated daily backups
- Retention policy
- Email notifications

### Phase 5: Advanced Features
- S3 storage support
- Incremental backups
- Point-in-time recovery
- Backup encryption
- Import/Export
- Backup verification

---

## Storage Requirements

### Typical Sizes

| Project Size | Database | Storage | Full Backup |
|--------------|----------|---------|-------------|
| Small        | 50 MB    | 100 MB  | ~200 MB     |
| Medium       | 500 MB   | 2 GB    | ~3 GB       |
| Large        | 5 GB     | 50 GB   | ~60 GB      |

### Disk Space Planning

For 10 projects with 7-day retention:

```
10 projects × 3 GB average × 7 backups = 210 GB
```

**Recommendation**: Allocate 500 GB for backups

---

## Best Practices

### For Users

1. **Enable automated backups** - Don't rely on manual backups
2. **Test restores** - Verify backups work before you need them
3. **Download critical backups** - Keep offsite copies
4. **Monitor backup status** - Check for failures
5. **Use incremental** - For frequently changing databases

### For Administrators

1. **Monitor disk space** - Backups can fill disk quickly
2. **Set retention policies** - Delete old backups automatically
3. **Enable S3 storage** - For production deployments
4. **Verify backup integrity** - Run periodic restore tests
5. **Document restore procedures** - Train team on recovery

---

## Disaster Recovery

### Scenarios

**1. Accidental data deletion**
- Restore from most recent backup
- Downtime: ~10 minutes

**2. Database corruption**
- Restore from last known good backup
- May lose recent data
- Downtime: ~15 minutes

**3. Server failure**
- Restore on new server from S3 backup
- Downtime: ~1 hour

**4. Ransomware**
- Restore from pre-infection backup
- Verify backups are clean
- Downtime: ~2 hours

### Recovery Time Objective (RTO)

- **Small projects**: 10-15 minutes
- **Medium projects**: 30-45 minutes
- **Large projects**: 1-2 hours

### Recovery Point Objective (RPO)

- **With daily backups**: Up to 24 hours of data loss
- **With hourly backups**: Up to 1 hour of data loss
- **With continuous replication**: Near-zero data loss

---

## Future Enhancements

- **Continuous Backup**: Real-time replication
- **Multi-region Replication**: Backups in multiple regions
- **Point-in-time Recovery**: Restore to any second
- **Backup Analytics**: Usage trends, recommendations
- **Backup Testing**: Automated restore verification
- **Compressed Backups**: Better compression algorithms
- **Deduplicated Backups**: Save space with deduplication

---

**Status**: Design Complete ✅
**Implementation**: Phase 3-5 🔨
