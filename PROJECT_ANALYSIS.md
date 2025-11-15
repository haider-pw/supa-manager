# Supa-Manager Project Analysis
**Date:** November 15, 2025
**Analyzed By:** Claude Code
**Purpose:** Complete technical documentation for implementing dynamic Supabase project provisioning

---

## 📁 Project Structure

```
supa-manager/
├── api/                    # API handlers (22 files)
│   ├── api.go             # Main router & middleware setup
│   ├── postPlatformProjects.go  # Project creation handler
│   ├── getProjectStatus.go      # Project status endpoint
│   ├── getProjectApi.go         # Project connection info
│   └── [18 other handlers]
├── conf/                  # Configuration management
│   ├── config.go          # Config struct & env loading
│   └── migrations.go      # Migration runner
├── database/              # Generated sqlc code (type-safe queries)
│   ├── db.go             # Database connection
│   ├── *.sql.go          # Generated query functions
│   └── models.go         # Database models
├── migrations/            # SQL migrations
│   └── 00_init.sql       # Initial schema
├── queries/               # SQL queries for sqlc
│   ├── accounts.sql
│   ├── organizations.sql
│   ├── projects.sql
│   └── organization_membership.sql
├── permisions/            # Authorization logic
│   └── org.go
├── utils/                 # Utility functions
│   └── sqlTypes.go
├── main.go                # Application entry point
├── go.mod                 # Go dependencies
├── Dockerfile             # Container build instructions
└── .env                   # Configuration (runtime)
```

---

## 🏗️ Architecture Overview

### **Current State**

```
┌─────────────┐
│   Browser   │
└──────┬──────┘
       │ HTTP
       ↓
┌─────────────┐
│   Studio    │  ← Supabase Studio (patched v1.24.04)
│   (Port     │     Next.js frontend
│   3000)     │
└──────┬──────┘
       │ REST API
       ↓
┌─────────────┐
│ Supa-Manager│  ← Go/Gin API Server
│   (Port     │     - Authentication (JWT + Argon2)
│   8080)     │     - Project Management
└──────┬──────┘     - Organization Management
       │
       ↓
┌─────────────┐
│  PostgreSQL │  ← Management Database
│   (Port     │     - User accounts
│   5432)     │     - Organizations
└─────────────┘     - Projects metadata
```

### **Missing: Per-Project Supabase Stacks**

```
┌──────────────────────────────────┐
│  Project: "staging"              │
│  Ref: upending-expectoration's   │
├──────────────────────────────────┤
│  ┌────────────┐  ┌────────────┐ │
│  │ PostgreSQL │  │  PostgREST │ │  ← NOT IMPLEMENTED
│  │            │  │            │ │
│  └────────────┘  └────────────┘ │
│  ┌────────────┐  ┌────────────┐ │
│  │   GoTrue   │  │    Kong    │ │
│  │            │  │            │ │
│  └────────────┘  └────────────┘ │
└──────────────────────────────────┘
```

---

## 🔑 Key Technologies

| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| **Language** | Go | 1.21 | Backend API |
| **Web Framework** | Gin | 1.10.0 | HTTP routing & middleware |
| **Database** | PostgreSQL | 15.1.0 | Management DB |
| **Database Driver** | pgx/v5 | 5.6.0 | Postgres client |
| **Query Generator** | sqlc | 2.0 | Type-safe SQL queries |
| **Auth** | JWT | jwt/v5 5.2.1 | Token-based auth |
| **Password Hashing** | Argon2 | 1.0.0 | Secure password storage |
| **Random Names** | babble | latest | Project ref generation |
| **Config** | envconfig | 1.4.0 | Environment variables |

---

## 📊 Database Schema

### **Tables**

#### 1. `accounts` - User Authentication
```sql
CREATE TABLE public.accounts (
    id            serial PRIMARY KEY,
    gotrue_id     text NOT NULL DEFAULT gen_random_uuid()::text,
    email         text NOT NULL,
    password_hash text NOT NULL,
    username      text NOT NULL,
    first_name    text,
    last_name     text,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now()
);
```

**Current Data:**
- 1 account exists: `haideritx@gmail.com` (username: `idekman`)

#### 2. `organizations` - Multi-tenancy
```sql
CREATE TABLE public.organizations (
    id         serial PRIMARY KEY,
    slug       text NOT NULL DEFAULT gen_random_uuid()::text,
    name       text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
```

#### 3. `organization_membership` - User-Org Relationship
```sql
CREATE TABLE public.organization_membership (
    organization_id int NOT NULL REFERENCES organizations(id),
    account_id      int NOT NULL REFERENCES accounts(id),
    role            text NOT NULL,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (organization_id, account_id)
);
```

#### 4. `project` - Supabase Project Metadata ⭐
```sql
CREATE TABLE public.project (
    id              serial PRIMARY KEY,
    project_ref     text NOT NULL,           -- e.g., "upending-expectoration's"
    project_name    text NOT NULL,           -- User-friendly name
    organization_id int NOT NULL REFERENCES organizations(id),
    status          text NOT NULL,           -- UNKNOWN, PROVISIONING, ACTIVE, FAILED, etc.
    cloud_provider  text NOT NULL DEFAULT 'k8s',
    region          text NOT NULL DEFAULT 'mars-1',
    jwt_secret      text NOT NULL,           -- UUID for project JWT signing
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now()
);
```

**Current Project:**
```
ID: 1
Ref: upending-expectoration's
Name: staging
Status: UNKNOWN  ← Stuck here, no infrastructure provisioned
Organization: 1
Cloud Provider: K8S
Region: MARS
```

---

## 🔌 API Endpoints

### **✅ Implemented & Working**

#### Authentication
```
POST /auth/token                     - Get JWT token (login)
POST /platform/signup                - Create account
POST /profile/password-check         - Verify password
```

#### Profile
```
GET /profile                         - Get user profile
GET /profile/permissions             - Get user permissions
```

#### Organizations
```
GET /organizations                   - List user's orgs
POST /platform/organizations         - Create organization
GET /organizations/:slug/members/reached-free-project-limit
GET /platform/organizations/:slug/billing/subscription
```

#### Projects
```
POST /platform/projects              - Create project ✅
GET /platform/projects               - List projects ✅
GET /platform/projects/:ref          - Get project details ✅
GET /projects/:ref/status            - Get project status ✅
GET /projects/:ref/api               - Get project API info
GET /projects/:ref/jwt-secret-update-status
```

### **❌ Not Implemented (Return 404)**

```
GET /platform/project/:ref/settings            - Project settings
PUT /platform/project/:ref/settings            - Update settings
GET /platform/organizations/:id/usage          - Usage statistics
GET /platform/notifications/summary            - Notification summary
GET /platform/stripe/invoices/overdue          - Billing (stub)
```

### **🔴 Critical Missing: Provisioning Logic**

No endpoints exist for:
- Starting project infrastructure
- Stopping/pausing projects
- Deleting projects
- Health checking project services
- Managing project lifecycle

---

## 🔍 Code Deep Dive

### **1. Project Creation Flow** (`api/postPlatformProjects.go`)

**Current Implementation:**
```go
func (a *Api) postPlatformProjects(c *gin.Context) {
    // 1. Authenticate user
    _, err := a.GetAccountFromRequest(c)

    // 2. Parse request body
    var createProject ProjectCreationBody
    c.BindJSON(&createProject)

    // 3. Generate project reference (babble = random words)
    projectRef := strings.ToLower(babble.NewBabbler().Babble())
    // Result: "upending-expectoration's"

    // 4. Insert into database
    proj, err := a.queries.CreateProject(c.Request.Context(), database.CreateProjectParams{
        ProjectRef:     projectRef,
        ProjectName:    createProject.Name,
        OrganizationID: createProject.OrgId,
        JwtSecret:      uuid.New().String(),
        CloudProvider:  strings.ToUpper(createProject.CloudProvider),
        Region:         strings.ToUpper(createProject.DbRegion),
    })

    // 5. Return response with STUB data
    c.JSON(http.StatusCreated, ProjectCreationResponse{
        ...
        Status: "UNKNOWN",  // ← Wrong! Should be "PROVISIONING"
        Endpoint: fmt.Sprintf("https://%s.%s", proj.ProjectRef, a.config.Domain.Base),
        AnonKey: "a.b.c",   // ← STUB!
        ServiceKey: "a.b.c" // ← STUB!
    })
}
```

**🔴 What's Missing:**
1. No infrastructure provisioning
2. No async processing
3. Status always "UNKNOWN"
4. Fake API keys returned
5. No error handling for provisioning failures

---

### **2. Status Endpoint** (`api/getProjectStatus.go`)

**Current Implementation:**
```go
func (a *Api) getProjectStatus(c *gin.Context) {
    projectRef := c.Param("ref")
    project, err := a.queries.GetProjectByRef(c, projectRef)

    c.JSON(http.StatusOK, gin.H{"status": project.Status})
}
```

**🔴 What's Missing:**
1. No actual health checks
2. Just returns DB status field
3. No service status (postgres, postgrest, gotrue, etc.)
4. No container status checking

---

### **3. Configuration** (`conf/config.go`)

**Config Structure:**
```go
type Config struct {
    DatabaseUrl       string           // Management DB connection
    Port              int              // API port (8080)
    EncryptionSecret  string           // For sensitive data
    JwtSecret         string           // For signing user JWTs
    AllowSignup       bool             // Enable/disable signup
    ServiceVersionUrl string           // Service version endpoint
    Domain            DomainSettings   // Domain configuration
    Postgres          PostgresSettings // Postgres defaults
}

type PostgresSettings struct {
    DiskSize       int     // Default: 10GB
    DefaultVersion string  // Default: "14.2"
    DockerImage    string  // Default: "supabase/postgres"
}

type DomainSettings struct {
    StudioUrl  string   // Studio frontend URL
    Base       string   // Base domain for projects
    DnsHookUrl *string  // Dynamic DNS webhook
    DnsHookKey *string  // DNS webhook auth key
}
```

**Current Values (.env):**
```bash
DATABASE_URL=postgres://postgres:postgres@database:5432/supabase
ALLOW_SIGNUP=true
JWT_SECRET=secret
ENCRYPTION_SECRET=secret
SERVICE_VERSION_URL=https://supamanager.io/updates
POSTGRES_DISK_SIZE=10
POSTGRES_DEFAULT_VERSION=14.2
POSTGRES_DOCKER_IMAGE=supabase/postgres
DOMAIN_STUDIO_URL=http://localhost:3000
DOMAIN_BASE=supamanager.io
DOMAIN_DNS_HOOK_URL=http://localhost:8081
DOMAIN_DNS_HOOK_KEY=mysecretkey
```

---

## 🎯 What Works vs What Doesn't

### ✅ **Working Components**

1. **Authentication System**
   - User signup/login
   - JWT token generation & validation
   - Argon2 password hashing
   - Session management

2. **Database Layer**
   - PostgreSQL connection pooling
   - Type-safe queries via sqlc
   - Migration system
   - ACID transactions

3. **API Server**
   - Gin router with CORS
   - Request validation
   - Error handling
   - JSON serialization

4. **Multi-tenancy**
   - Organizations
   - Organization membership
   - Project-organization linkage

5. **Project Metadata**
   - Creating project records
   - Storing project configuration
   - Generating unique project refs

### ❌ **Missing / Not Working**

1. **🚨 CRITICAL: No Provisioning System**
   - No Docker client integration
   - No Kubernetes client
   - No container orchestration
   - No service lifecycle management

2. **Infrastructure Management**
   - Can't spin up PostgreSQL for projects
   - Can't start PostgREST
   - Can't start GoTrue
   - Can't configure Kong gateway

3. **Status Tracking**
   - Status never changes from "UNKNOWN"
   - No health checks
   - No service monitoring
   - No error reporting

4. **Project Operations**
   - Can't pause/resume projects
   - Can't delete projects
   - Can't scale resources
   - Can't backup/restore

5. **Networking**
   - No port allocation
   - No network isolation
   - No load balancing
   - No DNS management

---

## 🔧 Dependencies Analysis

### **Core Dependencies**

```go
github.com/gin-gonic/gin v1.10.0           // Web framework ✅
github.com/jackc/pgx/v5 v5.6.0             // PostgreSQL driver ✅
github.com/golang-jwt/jwt/v5 v5.2.1        // JWT auth ✅
github.com/matthewhartstonge/argon2 v1.0.0 // Password hashing ✅
github.com/tjarratt/babble v0.0.0          // Random name generation ✅
github.com/google/uuid v1.6.0              // UUID generation ✅
```

### **⚠️ Missing Dependencies for Provisioning**

```go
// Docker provisioning
github.com/docker/docker v24.0.0+incompatible
github.com/docker/go-connections v0.4.0

// OR Kubernetes provisioning
k8s.io/client-go v0.28.0
k8s.io/api v0.28.0
k8s.io/apimachinery v0.28.0

// Template rendering
text/template (stdlib) ✅
```

---

## 📝 SQL Queries (sqlc)

### **Projects** (`queries/projects.sql`)

```sql
-- name: GetProjectsForAccountId :many
SELECT p.*
FROM organization_membership om
JOIN project p on om.organization_id = p.organization_id
WHERE account_id = $1;

-- name: CreateProject :one
INSERT INTO project (project_ref, project_name, organization_id, status, jwt_secret, cloud_provider, region)
VALUES ($1, $2, $3, 'UNKNOWN', $4, $5, $6)
RETURNING *;

-- name: GetProjectByRef :one
SELECT * FROM project WHERE project_ref = $1;
```

**🔴 Missing Queries:**
```sql
-- name: UpdateProjectStatus :one
UPDATE project SET status = $2, updated_at = now()
WHERE project_ref = $1 RETURNING *;

-- name: DeleteProject :exec
DELETE FROM project WHERE project_ref = $1;

-- name: GetProjectsByStatus :many
SELECT * FROM project WHERE status = $1;
```

---

## 🎬 Actual Execution Flow (Current)

### **User Creates Project "staging"**

```
1. User fills form in Studio UI
   ↓
2. POST /platform/projects
   {
     "name": "staging",
     "db_pass": "NoAdmin@456@786",
     "cloud_provider": "K8S",
     "db_region": "MARS",
     "desired_instance_size": "micro"
   }
   ↓
3. Backend validates JWT token
   ↓
4. Generate project_ref: "upending-expectoration's"
   ↓
5. Insert into DB with status="UNKNOWN"
   ↓
6. Return response with fake data
   {
     "status": "UNKNOWN",
     "anon_key": "a.b.c",
     "endpoint": "https://upending-expectoration's.supamanager.io"
   }
   ↓
7. ❌ NOTHING ELSE HAPPENS
   - No containers created
   - No services started
   - Status stuck on "UNKNOWN"
   - Studio polls /projects/:ref/status forever
```

---

## 🎯 Required Supabase Services Per Project

### **Minimal Stack**

```yaml
services:
  postgres:          # Project database
  postgrest:         # REST API
  gotrue:            # Authentication
  kong:              # API Gateway
```

### **Full Stack (Optional)**

```yaml
  realtime:          # WebSocket server
  storage:           # File storage (MinIO)
  imgproxy:          # Image optimization
  meta:              # Database metadata API
  functions:         # Edge functions (Deno)
  analytics:         # Log aggregation
```

### **Service Communication**

```
Kong (Gateway)
  ├→ PostgREST → PostgreSQL
  ├→ GoTrue → PostgreSQL
  ├→ Realtime → PostgreSQL
  ├→ Storage → PostgreSQL + MinIO
  └→ Meta → PostgreSQL
```

---

## 🚀 Required Implementation Steps

### **Phase 1: Analysis** ✅ (CURRENT)

- [x] Understand codebase structure
- [x] Document database schema
- [x] Map API endpoints
- [x] Identify missing components
- [ ] Search for any provisioning code
- [ ] Document Supabase architecture

### **Phase 2: Design**

- [ ] Choose Docker vs K8S approach
- [ ] Design docker-compose template
- [ ] Plan port allocation strategy
- [ ] Design network isolation
- [ ] Define status state machine

### **Phase 3: Implementation**

- [ ] Add Docker SDK dependency
- [ ] Create provisioning package
- [ ] Implement template rendering
- [ ] Add async job processing
- [ ] Update project creation logic
- [ ] Implement status updates
- [ ] Add health checks

### **Phase 4: Testing**

- [ ] Test project creation
- [ ] Test service startup
- [ ] Test Studio connection
- [ ] Test project deletion
- [ ] Load testing

### **Phase 5: Enhancement**

- [ ] Implement pause/resume
- [ ] Add backup/restore
- [ ] Implement scaling
- [ ] Add monitoring

---

## 🔍 Key Findings

### **1. Project is Early Stage**
- Core API framework is solid
- Database layer is well-designed
- Authentication works properly
- **BUT**: No actual provisioning implemented

### **2. No Container Orchestration**
```bash
$ docker ps | grep "upending"
# Returns nothing!
```
Projects are metadata-only.

### **3. Hardcoded Stubs**
```go
AnonKey:    "a.b.c",      // Should be real JWT
ServiceKey: "a.b.c",      // Should be real JWT
Status:     "UNKNOWN",    // Should track real state
```

### **4. Config is K8S-Oriented**
```go
cloud_provider text NOT NULL default 'k8s'
```
But no K8S client exists.

### **5. Missing Error Handling**
No handling for:
- Provisioning failures
- Service crashes
- Resource exhaustion
- Network errors

---

## 📋 Next Steps (Phase 1 Remaining)

1. **Search for any provisioning code**
   ```bash
   grep -r "docker" supa-manager/
   grep -r "kubernetes\|k8s" supa-manager/
   grep -r "container" supa-manager/
   ```

2. **Check for external scripts**
   ```bash
   find . -name "*.sh" -o -name "*.py"
   ```

3. **Review VERSION_SERVICE_URL**
   - What does https://supamanager.io/updates return?
   - Service version management?

4. **Document Supabase components**
   - Required environment variables
   - Inter-service dependencies
   - Network requirements

---

## 💡 Recommendations

### **Quick Win: Docker Compose Approach**

**Pros:**
- ✅ Simpler to implement
- ✅ No K8S cluster needed
- ✅ Faster development
- ✅ Good for MVP/POC

**Cons:**
- ❌ Limited to single-server
- ❌ Manual scaling
- ❌ No advanced orchestration

### **Long-term: Kubernetes**

**Pros:**
- ✅ Production-ready
- ✅ Auto-scaling
- ✅ Self-healing
- ✅ Multi-node

**Cons:**
- ❌ Complex setup
- ❌ Steep learning curve
- ❌ Requires cluster

### **Suggested Path:**
1. Implement Docker Compose provisioning (Phase 2-3)
2. Get it working end-to-end (Phase 4)
3. Add K8S support later (Phase 5+)

---

## 📚 Reference Material

### **Supabase Official Docs**
- Self-hosting: https://supabase.com/docs/guides/self-hosting
- Docker setup: https://github.com/supabase/supabase/tree/master/docker

### **Relevant Technologies**
- Docker SDK for Go: https://docs.docker.com/engine/api/sdk/
- sqlc: https://docs.sqlc.dev/
- Gin framework: https://gin-gonic.com/docs/

---

## 🔗 Related Documentation

For detailed information about the complete Supabase service architecture, see:
- **[SUPABASE_ARCHITECTURE.md](/home/haider/supabase-manager/SUPABASE_ARCHITECTURE.md)** - Complete breakdown of all 12 services, dependencies, networking, and provisioning strategies

---

**End of Phase 1 Analysis**
