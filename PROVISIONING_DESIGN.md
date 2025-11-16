# Supabase Project Provisioning Design

## Overview

This document describes the automated provisioning system for Supabase projects. Users can create, manage, and delete Supabase instances through the web UI with zero manual terminal work.

---

## Architecture

### High-Level Flow

```
User clicks "Create Project" in Studio UI
              ↓
POST /platform/projects (API)
              ↓
    Docker Provisioner
              ↓
     Project Directory Created
              ↓
   Templates Rendered with Config
              ↓
    Docker Containers Started
              ↓
    Health Checks Verified
              ↓
   Project Status: ACTIVE_HEALTHY
              ↓
User gets project endpoint URL
```

---

## Directory Structure

### Project Storage Layout

```
/var/lib/supamanager/projects/
├── abc123/                          # Project ID
│   ├── docker-compose.yml           # Rendered from template
│   ├── kong.yml                     # Kong API Gateway config
│   ├── vector.yml                   # Logging config
│   ├── .env                         # Project secrets
│   └── volumes/                     # Docker volume data
│       ├── db/                      # PostgreSQL data
│       └── storage/                 # Storage files
├── def456/
│   └── ...
└── ghi789/
    └── ...
```

### Template Directory

```
supa-manager/templates/
├── project-compose.tmpl.yml         # Docker Compose template
├── kong.tmpl.yml                    # Kong config template
└── vector.tmpl.yml                  # Vector logging template
```

---

## Provisioner Interface

### Core Interface

```go
type Provisioner interface {
    CreateProject(ctx, config) (*ProjectInfo, error)
    GetProjectInfo(ctx, projectID) (*ProjectInfo, error)
    UpdateProject(ctx, projectID, config) error
    PauseProject(ctx, projectID) error
    ResumeProject(ctx, projectID) error
    DeleteProject(ctx, projectID) error
    ListProjects(ctx) ([]*ProjectInfo, error)
    GetLogs(ctx, projectID, service, tail) ([]string, error)
    ExecuteCommand(ctx, projectID, service, cmd) (string, error)
}
```

### Implementation

- **DockerProvisioner**: Uses Docker SDK to manage containers
- Future: **KubernetesProvisioner** (if needed)

---

## Project Lifecycle

### 1. Creation Flow

```mermaid
graph TD
    A[User clicks Create] --> B[Generate Project ID]
    B --> C[Generate Secrets]
    C --> D[Allocate Ports]
    D --> E[Create Project Directory]
    E --> F[Render Templates]
    F --> G[Start Containers]
    G --> H{Health Checks}
    H -->|Pass| I[Status: ACTIVE]
    H -->|Fail| J[Status: FAILED]
    I --> K[Return Endpoint URL]
```

**Steps:**
1. User submits project name and configuration
2. Backend generates unique project ID (e.g., `abc123`)
3. Backend generates secure secrets (JWT, API keys, passwords)
4. Backend allocates available ports for API and DB
5. Create project directory structure
6. Render docker-compose.yml with project-specific values
7. Render Kong and Vector configs
8. Start all Docker containers
9. Wait for all services to become healthy
10. Update database with project info
11. Return project endpoint to user

### 2. Pause Flow

```
User clicks "Pause Project"
              ↓
    Stop all containers
    (keep volumes intact)
              ↓
  Status: PAUSED
```

### 3. Resume Flow

```
User clicks "Resume Project"
              ↓
  Start all containers
              ↓
   Wait for health checks
              ↓
Status: ACTIVE_HEALTHY
```

### 4. Delete Flow

```
User clicks "Delete Project"
              ↓
   Confirm deletion
              ↓
  Stop all containers
              ↓
 Remove all containers
              ↓
  Remove all volumes
              ↓
 Remove all networks
              ↓
 Delete project directory
              ↓
Update database status: DELETED
```

---

## Configuration

### Project Configuration (Input)

```go
type ProjectConfig struct {
    ProjectID      string  // Generated: "abc123"
    ProjectName    string  // User input: "My App"
    OrganizationID string  // Owner org
    Region         string  // "us-east" (future multi-region)

    // Ports (auto-allocated)
    DBPort         int     // 54321
    APIPort        int     // 54322
    StudioPort     int     // 54323 (optional)

    // Security (auto-generated)
    DBPassword     string
    JWTSecret      string
    AnonKey        string
    ServiceKey     string

    // Resources (future use)
    CPULimit       string  // "1.0"
    MemoryLimit    string  // "2GB"
    StorageLimit   string  // "10GB"
}
```

### Project Info (Output)

```go
type ProjectInfo struct {
    ProjectID     string
    Status        ProjectStatus  // CREATING, ACTIVE, PAUSED, etc.
    Endpoint      string         // https://abc123.supamanager.io
    DBEndpoint    string         // postgres://...@host:54321/postgres

    Containers    map[string]string  // service -> container ID
    HealthChecks  map[string]bool    // service -> healthy

    CPUUsage      float64
    MemoryUsage   uint64
    StorageUsage  uint64

    CreatedAt     string
    UpdatedAt     string
}
```

---

## Supabase Stack Services

Each project provisions these containers:

| Service | Image | Port | Purpose |
|---------|-------|------|---------|
| **db** | supabase/postgres:15.1.0.147 | 5432 | PostgreSQL database |
| **kong** | kong:2.8.1 | 8000 | API Gateway |
| **auth** | supabase/gotrue:v2.99.0 | 9999 | Authentication (GoTrue) |
| **rest** | postgrest/postgrest:v11.2.2 | 3000 | REST API (PostgREST) |
| **realtime** | supabase/realtime:v2.25.50 | 4000 | WebSocket (Realtime) |
| **storage** | supabase/storage-api:v0.43.11 | 5000 | File storage |
| **imgproxy** | darthsim/imgproxy:v3.8.0 | 5001 | Image transformation |
| **meta** | supabase/postgres-meta:v0.68.0 | 8080 | Database management |
| **analytics** | supabase/logflare:1.4.0 | 4000 | Logging & Analytics |
| **vector** | timberio/vector:0.28.1-alpine | - | Log aggregation |

**Total: 10 containers per project**

---

## Port Allocation Strategy

### Port Ranges

- **API Ports**: 54000-54999 (Kong gateway)
- **DB Ports**: 55000-55999 (PostgreSQL)
- **Studio Ports**: 56000-56999 (Optional per-project Studio)

### Allocation Logic

```go
func AllocatePort(projectID string, service string) int {
    // Check database for used ports
    usedPorts := GetUsedPorts()

    // Find next available port in range
    basePort := GetBasePortForService(service)
    for port := basePort; port < basePort+1000; port++ {
        if !usedPorts[port] {
            return port
        }
    }

    return -1 // No ports available
}
```

---

## Security

### Secret Generation

All secrets are generated using `crypto/rand` for cryptographic security:

```go
func GenerateSecret(length int) (string, error) {
    bytes := make([]byte, length)
    _, err := rand.Read(bytes)
    if err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes), nil
}
```

### Secret Requirements

- **JWT Secret**: 64 bytes (512 bits)
- **Database Password**: 32 bytes (256 bits)
- **API Keys**: 32 bytes (256 bits)

### Storage

- Secrets stored in project's `.env` file
- File permissions: `0600` (owner read/write only)
- Encrypted in database (using ENCRYPTION_SECRET)

---

## Health Monitoring

### Service Health Checks

Each service has a health check endpoint:

```go
type HealthCheck struct {
    Service   string
    Endpoint  string
    Healthy   bool
    LastCheck time.Time
    Error     string
}
```

### Overall Project Health

```
ACTIVE_HEALTHY   - All services healthy
ACTIVE_UNHEALTHY - Some services unhealthy
CREATING         - Still provisioning
PAUSED           - Intentionally stopped
FAILED           - Provisioning failed
```

---

## API Integration

### Create Project Endpoint

```http
POST /platform/projects
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "My App",
  "region": "us-east",
  "organization_id": "org_123"
}
```

**Response:**

```json
{
  "id": "abc123",
  "name": "My App",
  "status": "CREATING",
  "endpoint": "https://abc123.supamanager.io",
  "database_url": "postgres://postgres:***@localhost:54321/postgres",
  "api_keys": {
    "anon": "eyJ...",
    "service_role": "eyJ..."
  }
}
```

### Get Project Status

```http
GET /platform/projects/abc123
```

**Response:**

```json
{
  "id": "abc123",
  "status": "ACTIVE_HEALTHY",
  "services": {
    "db": "healthy",
    "kong": "healthy",
    "auth": "healthy",
    "rest": "healthy",
    "realtime": "healthy",
    "storage": "healthy"
  },
  "resource_usage": {
    "cpu": 0.45,
    "memory": 1.2,
    "storage": 5.7
  }
}
```

---

## Implementation Phases

### ✅ Phase 1: Analysis (Complete)
- Analyzed codebase
- Identified provisioning requirements
- Documented architecture

### ✅ Phase 2: Design (Complete)
- Created docker-compose templates
- Designed Provisioner interface
- Planned directory structure

### 🔨 Phase 3: Core Implementation (Next)
1. Implement Docker SDK integration
2. Implement CreateProject provisioning
3. Implement health monitoring
4. Test single project creation

### 🔮 Phase 4: Lifecycle Management
1. Implement Pause/Resume
2. Implement Delete
3. Test full lifecycle
4. Add project listing

### 🔮 Phase 5: Production Ready
1. Add error recovery
2. Add resource monitoring
3. Add backup/restore
4. Performance optimization

---

## Testing Strategy

### Unit Tests
- Template rendering
- Secret generation
- Port allocation
- Configuration validation

### Integration Tests
- Full project creation
- Service health checks
- Project lifecycle operations

### End-to-End Tests
1. User creates project via UI
2. Backend provisions containers
3. All services become healthy
4. API endpoints are accessible
5. User can connect to database
6. User can pause/resume project
7. User can delete project

---

## Future Enhancements

### Multi-Region Support
- Deploy projects to different regions
- Automatic region selection based on user location

### Auto-Scaling
- Scale resources based on usage
- Automatic container restarts on failure

### Backup/Restore
- Automated daily backups
- Point-in-time recovery

### Monitoring Dashboard
- Real-time metrics
- Alerting on failures
- Usage analytics

### Custom Domains
- User can configure custom domain
- Automatic SSL certificate provisioning

---

## Resources

- Docker SDK: https://docs.docker.com/engine/api/sdk/
- Docker Compose Spec: https://docs.docker.com/compose/compose-file/
- Supabase Self-Hosting: https://supabase.com/docs/guides/self-hosting

---

**Status**: Phase 2 Complete ✅
**Next**: Phase 3 - Core Implementation 🔨
