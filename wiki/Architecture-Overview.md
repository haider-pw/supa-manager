# Architecture Overview

Understanding the SupaManager system design and components.

---

## System Architecture

SupaManager consists of three main layers:

```
┌─────────────────────────────────────────────────────────┐
│                        Browser                          │
│                    (User Interface)                     │
└────────────────────┬────────────────────────────────────┘
                     │ HTTP/HTTPS
                     ▼
┌─────────────────────────────────────────────────────────┐
│                    Studio UI (Next.js)                  │
│                      Port 3000                          │
│  • Project Management Interface                         │
│  • Authentication UI                                    │
│  • Database Management                                  │
└────────────────────┬────────────────────────────────────┘
                     │ REST API
                     ▼
┌─────────────────────────────────────────────────────────┐
│               SupaManager API (Go/Gin)                  │
│                      Port 8080                          │
│  • User Authentication (JWT)                            │
│  • Organization Management                              │
│  • Project Metadata Management                          │
│  • API Gateway                                          │
└────────────────────┬────────────────────────────────────┘
                     │ SQL
                     ▼
┌─────────────────────────────────────────────────────────┐
│              PostgreSQL Database                        │
│                      Port 5432                          │
│  • User accounts & authentication                       │
│  • Organizations & memberships                          │
│  • Project metadata                                     │
│  • Migrations tracking                                  │
└─────────────────────────────────────────────────────────┘
```

---

## Core Components

### 1. Studio UI (Frontend)

**Technology:** Next.js 12 (React framework)
**Port:** 3000
**Purpose:** Web-based management interface

**Key Features:**
- Project creation and management
- Database schema editor
- SQL query editor
- Authentication management
- User interface for all API operations

**How it works:**
- Patched version of official Supabase Studio
- Custom patches integrate with SupaManager API
- Server-side rendering with Next.js
- Client-side state management

**Environment Variables:**
```bash
PLATFORM_PG_META_URL=http://supa-manager:8080/pg  # Server-side API
NEXT_PUBLIC_API_URL=http://localhost:8080         # Client-side API
```

### 2. SupaManager API (Backend)

**Technology:** Go 1.21 with Gin framework
**Port:** 8080
**Purpose:** Management API and business logic

**Key Packages:**
- `github.com/gin-gonic/gin` - HTTP router
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/golang-jwt/jwt/v5` - JWT authentication
- `golang.org/x/crypto/argon2` - Password hashing

**API Structure:**
```
api/
├── api.go                      # Router setup & middleware
├── postPlatformProjects.go     # Create project
├── getPlatformProjects.go      # List projects
├── getPlatformProject.go       # Get single project
├── deletePlatformProject.go    # Delete project
├── getProjectStatus.go         # Project status
├── getProjectApi.go            # Project connection info
└── [15 more handlers...]       # Other endpoints
```

**Key Features:**
- RESTful API design
- JWT-based authentication
- Role-based access control (RBAC)
- Organization-scoped resources
- Type-safe database queries (sqlc)
- Automatic migrations

### 3. PostgreSQL Database

**Technology:** PostgreSQL 15.1
**Port:** 5432
**Purpose:** Management data storage

**Schema Overview:**
```sql
accounts                    -- User accounts
├── id (primary key)
├── email
├── password (argon2)
└── created_at

organizations              -- Organizations/teams
├── id (primary key)
├── name
├── slug
└── created_at

organization_membership    -- User-org relationships
├── organization_id
├── account_id
└── role (owner/member)

project                    -- Project metadata
├── id (primary key)
├── project_ref            -- Unique identifier
├── project_name
├── organization_id
├── status                 -- UNKNOWN/PROVISIONING/ACTIVE/PAUSED
├── cloud_provider         -- k8s/aws/gcp
├── region
└── jwt_secret
```

---

## Data Flow

### Authentication Flow

```
1. User enters credentials in Studio
   ↓
2. Studio POST /platform/auth/login
   ↓
3. API validates credentials (argon2)
   ↓
4. API generates JWT token
   ↓
5. Token returned to Studio
   ↓
6. Studio stores token in memory/cookies
   ↓
7. Subsequent requests include: Authorization: Bearer <token>
```

### Project Creation Flow (Current)

```
1. User clicks "New Project" in Studio
   ↓
2. Studio POST /platform/projects
   ├── project_name
   ├── org_id
   ├── cloud_provider
   └── region
   ↓
3. API generates project_ref (random words)
   ↓
4. API creates database record
   ├── status: "UNKNOWN"
   ├── jwt_secret: UUID
   └── returns fake keys (a.b.c)
   ↓
5. ⚠️ NO PROVISIONING HAPPENS
   ↓
6. Project exists in DB but no infrastructure
```

### Project Creation Flow (Future)

```
1. User clicks "New Project" in Studio
   ↓
2. Studio POST /platform/projects
   ↓
3. API creates database record
   ↓
4. API triggers provisioner service
   ├── Generate JWT keys (anon_key, service_key)
   ├── Create docker-compose.yml from template
   ├── Allocate unique ports
   ├── Create Docker network
   ├── Start 12 Supabase services
   └── Monitor health checks
   ↓
5. Update project status
   ├── PROVISIONING → Creating infrastructure
   ├── STARTING → Waiting for services
   └── ACTIVE → All services healthy
   ↓
6. Return real connection info to user
```

---

## Service Communication

### Internal Docker Network

All services communicate via Docker's internal DNS:

```
Container Name              Internal Hostname    Port
─────────────────────────  ──────────────────  ─────
supabase-manager-studio-1      studio          3000
supabase-manager-supa-manager-1 supa-manager   8080
supabase-manager-database-1    database        5432
```

**Example:**
- Studio server-side code calls: `http://supa-manager:8080/pg`
- API connects to database: `postgres://postgres@database:5432/supabase`

### External Access

Only specific ports are exposed to the host:

```
Host Port    →    Container
─────────         ─────────────────────
3000         →    studio:3000
8080         →    supa-manager:8080
5432         →    database:5432
```

Browsers and external clients use `localhost:PORT` to access services.

---

## Security Model

### Authentication

**Method:** JWT (JSON Web Tokens)
**Algorithm:** HS256
**Secret:** Configured in `JWT_SECRET` environment variable

**Token Payload:**
```json
{
  "account_id": 123,
  "email": "user@example.com",
  "exp": 1234567890
}
```

### Password Storage

**Algorithm:** Argon2id
**Parameters:**
- Memory: 64MB
- Iterations: 3
- Parallelism: 2
- Salt: Random 16 bytes per password

### Authorization

**Model:** Organization-based RBAC

```
User → Organization Membership → Projects
        ↓
     role: owner | member
```

**Access Control:**
- Users can only access projects in their organizations
- Organization owners have full control
- Members have limited permissions (configured)

### API Security

**Headers Required:**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**CORS Configuration:**
- Configured for Studio origin
- Credentials allowed
- Specific methods whitelisted

---

## Database Design

### Type-Safe Queries (sqlc)

All database queries are generated from SQL files:

```
queries/
├── accounts.sql              -- name: GetAccountByEmail :one
├── organizations.sql         -- name: CreateOrganization :one
├── projects.sql             -- name: CreateProject :one
└── organization_membership.sql
```

**Benefits:**
- Compile-time type safety
- No ORM overhead
- Direct SQL control
- Auto-generated Go code

### Migrations

**System:** Custom migration runner

```
migrations/
└── 00_init.sql              -- Initial schema

Migration tracking:
CREATE TABLE migrations (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    applied_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

## Planned Architecture (Future)

### With Provisioning System

```
┌──────────────────────────────────────────────────────┐
│                   Studio UI                          │
└───────────────────┬──────────────────────────────────┘
                    ▼
┌──────────────────────────────────────────────────────┐
│              SupaManager API                         │
│  ┌────────────────────────────────────────────────┐  │
│  │         Provisioner Service (New)              │  │
│  │  • Template generator                          │  │
│  │  • Docker SDK integration                      │  │
│  │  • Health monitoring                           │  │
│  │  • Status tracking                             │  │
│  └────────────────────────────────────────────────┘  │
└───────────────────┬──────────────────────────────────┘
                    ▼
        ┌───────────────────────┐
        │  Docker Engine        │
        │  ┌─────────────────┐  │
        │  │ Project 1       │  │
        │  │ (12 services)   │  │
        │  └─────────────────┘  │
        │  ┌─────────────────┐  │
        │  │ Project 2       │  │
        │  │ (12 services)   │  │
        │  └─────────────────┘  │
        └───────────────────────┘
```

Each project will run a complete Supabase stack:
1. PostgreSQL (project database)
2. Kong (API gateway)
3. GoTrue (auth)
4. PostgREST (REST API)
5. Realtime (WebSockets)
6. Storage (file storage)
7. ImgProxy (image processing)
8. Meta (DB metadata)
9. Functions (edge functions)
10. Analytics (logging)
11. Vector (log collection)
12. Studio (management UI)

See [Supabase Stack](Supabase-Stack) for complete details.

---

## Technology Stack Summary

| Component | Technology | Purpose |
|-----------|-----------|---------|
| Frontend | Next.js 12 | Web interface |
| Backend | Go 1.21 + Gin | REST API |
| Database | PostgreSQL 15.1 | Data storage |
| Query Layer | sqlc | Type-safe queries |
| Authentication | JWT + Argon2 | User auth |
| Containerization | Docker | Deployment |
| Orchestration | Docker Compose | Service management |

---

## Performance Considerations

### Current System

**Lightweight:** Minimal resource usage
- API: ~100MB RAM
- Studio: ~200MB RAM
- Database: ~50MB RAM
- Total: ~350MB RAM

**Scalability:** Limited by single server

### With Provisioning

**Resource Requirements per Project:**
- ~2GB RAM for 12 services
- ~20GB disk space
- Unique ports required

**Scaling Strategy:**
- Horizontal: Multiple SupaManager instances
- Vertical: Larger server for more projects
- Future: Kubernetes for distributed deployment

---

## Related Documentation

- [Supabase Stack](Supabase-Stack) - Complete service architecture
- [Docker Networking](Docker-Networking) - Container communication
- [Service Communication](Service-Communication) - API interactions
- [Database Schema](Database-Schema) - Complete schema reference

---

**Next:** Learn about the [Supabase Stack](Supabase-Stack) that each project will run.
