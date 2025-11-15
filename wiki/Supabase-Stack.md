# Supabase Stack Architecture

Complete breakdown of the 12 services required for each Supabase project.

---

## Overview

Each Supabase project is a complete stack of 12 interconnected microservices that work together to provide the full Supabase functionality.

```
                    External Clients (Browsers, Apps)
                                 │
                                 ▼
                          Kong API Gateway
                          (Single Entry Point)
                                 │
        ┌───────────────────────┼────────────────────────┐
        │                       │                        │
        ▼                       ▼                        ▼
    GoTrue (Auth)          PostgREST              Realtime
    Port 9999              Port 3000              Port 4000
        │                       │                        │
        └───────────────────────┼────────────────────────┘
                                ▼
                         PostgreSQL Database
                            Port 5432
```

---

## Service Breakdown

### 1. PostgreSQL Database

**Image:** `supabase/postgres:15.1.1.41`
**Port:** 5432 (dynamic per project)
**Purpose:** Main project database with Supabase extensions

**Key Features:**
- Pre-configured with Supabase extensions
- Roles: `anon`, `authenticated`, `service_role`
- Extensions: `pg_graphql`, `pgsodium`, `pgvector`, `postgis`
- Logical replication enabled for Realtime
- Custom authentication schema

**Environment Variables:**
```bash
POSTGRES_PASSWORD=project-specific-password
POSTGRES_DB=postgres
JWT_SECRET=project-jwt-secret
JWT_EXP=3600
```

**Resource Requirements:**
- **RAM:** 256MB minimum, 1GB recommended
- **Disk:** 10GB minimum, grows with data
- **CPU:** 0.5 cores minimum

**Initialization:**
- Creates auth schema
- Creates storage schema
- Creates realtime schema
- Sets up webhooks
- Configures JWT authentication

---

### 2. Kong API Gateway

**Image:** `kong:2.8.1`
**Ports:** 8000 (HTTP), 8443 (HTTPS)
**Purpose:** Unified API gateway and request router

**Routing Table:**
```
Path                  → Backend Service
────────────────────    ─────────────────
/auth/v1/*           → GoTrue :9999
/rest/v1/*           → PostgREST :3000
/realtime/v1/*       → Realtime :4000
/storage/v1/*        → Storage :5000
/functions/v1/*      → Functions :9000
/analytics/v1/*      → Analytics :4000
/pg/*                → Meta :8080
/*                   → Studio :3000
```

**Key Features:**
- JWT authentication middleware
- CORS handling
- Rate limiting
- Request transformation
- API key validation

**Configuration:**
- DB-less mode (declarative config)
- Custom `kong.yml` configuration
- Environment variable substitution

**Security:**
- Validates JWT tokens
- Enforces role-based access
- API key authentication
- Basic auth for Studio

---

### 3. GoTrue (Authentication)

**Image:** `supabase/gotrue:v2.149.0`
**Port:** 9999
**Purpose:** User authentication and management

**Features:**
- Email/password authentication
- Magic link authentication
- OAuth providers (Google, GitHub, etc.)
- JWT token generation
- User confirmation via email
- Password recovery
- Multi-factor authentication (MFA)

**Supported OAuth Providers:**
- Google
- GitHub
- GitLab
- Bitbucket
- Azure
- Facebook
- Twitter
- Apple
- And more...

**API Endpoints:**
```
POST   /signup       - Create new user
POST   /token        - Login / get JWT
POST   /verify       - Verify email
POST   /recover      - Password recovery
GET    /user         - Get user info
PUT    /user         - Update user
POST   /logout       - Logout
```

**Environment Variables:**
```bash
GOTRUE_DB_DATABASE_URL=postgres://...
GOTRUE_SITE_URL=https://project-ref.domain.com
GOTRUE_JWT_SECRET=shared-secret
API_EXTERNAL_URL=https://project-ref.domain.com
GOTRUE_SMTP_HOST=smtp.example.com  # For email
GOTRUE_EXTERNAL_GOOGLE_ENABLED=true
GOTRUE_EXTERNAL_GOOGLE_CLIENT_ID=...
```

---

### 4. PostgREST (REST API)

**Image:** `postgrest/postgrest:v12.0.1`
**Port:** 3000
**Purpose:** Automatic REST API from PostgreSQL schema

**How It Works:**
1. Inspects your database schema
2. Generates REST endpoints automatically
3. Enforces row-level security (RLS)
4. Provides filtering, sorting, pagination

**Example:**
```sql
-- Create a table
CREATE TABLE todos (
    id SERIAL PRIMARY KEY,
    task TEXT,
    done BOOLEAN DEFAULT false
);

-- Automatic API endpoints:
GET    /todos              # List all
POST   /todos              # Create new
GET    /todos?id=eq.1      # Filter by id
PATCH  /todos?id=eq.1      # Update
DELETE /todos?id=eq.1      # Delete
```

**Features:**
- Automatic API generation
- Row-level security (RLS)
- JWT role-based access
- Complex queries via URL params
- Full-text search
- Computed columns
- Stored procedures as RPC

**Query Examples:**
```bash
# Filter
GET /todos?done=eq.false

# Sort
GET /todos?order=id.desc

# Limit
GET /todos?limit=10

# Join
GET /users?select=*,posts(*)

# Full-text search
GET /articles?title=fts.database
```

---

### 5. Realtime Service

**Image:** `supabase/realtime:v2.28.32`
**Port:** 4000
**Purpose:** WebSocket subscriptions and real-time data

**Features:**
1. **Database Changes** - Listen to INSERT/UPDATE/DELETE
2. **Broadcast** - Send messages between clients
3. **Presence** - Track who's online
4. **Postgres CDC** - Change data capture

**Usage Example:**
```javascript
// Subscribe to database changes
supabase
  .channel('todos')
  .on('postgres_changes',
    { event: '*', schema: 'public', table: 'todos' },
    (payload) => console.log(payload)
  )
  .subscribe()

// Broadcast messages
channel.send({
  type: 'broadcast',
  event: 'cursor-pos',
  payload: { x: 100, y: 200 }
})

// Track presence
channel.track({ user: 'john', status: 'online' })
```

**Architecture:**
- Uses PostgreSQL logical replication
- WebSocket connections
- Elixir/Phoenix framework
- Horizontal scaling support

---

### 6. Storage API

**Image:** `supabase/storage-api:v1.0.6`
**Port:** 5000
**Purpose:** Object storage (file uploads/downloads)

**Features:**
- File uploads
- Public and private buckets
- Row-level security (RLS) for files
- Image transformation
- Resumable uploads
- CDN integration

**Buckets:**
```javascript
// Create bucket
await supabase.storage.createBucket('avatars', {
  public: true
})

// Upload file
await supabase.storage
  .from('avatars')
  .upload('user-1.png', file)

// Get public URL
const { data } = supabase.storage
  .from('avatars')
  .getPublicUrl('user-1.png')
```

**Storage Backends:**
- Local filesystem (default)
- S3-compatible storage (MinIO, AWS S3, etc.)

**Image Transformation:**
Via ImgProxy integration:
```
/storage/v1/render/image/public/avatar.png?width=200&height=200
```

---

### 7. ImgProxy

**Image:** `darthsim/imgproxy:v3.8.0`
**Port:** 5001
**Purpose:** On-the-fly image processing

**Features:**
- Resize images
- Crop images
- Format conversion (WebP, AVIF)
- Quality adjustment
- Watermarks
- Blur effects

**Usage:**
```
/render/image/authenticated/bucket/image.jpg
  ?width=300
  &height=200
  &resize=fill
  &quality=80
  &format=webp
```

**Performance:**
- Fast processing (libvips)
- ETags for caching
- Source image caching
- Optimized for CDN

---

### 8. Postgres Meta (pg-meta)

**Image:** `supabase/postgres-meta:v0.80.0`
**Port:** 8080
**Purpose:** Database metadata API for Studio

**Capabilities:**
- List tables, columns, types
- Create/modify schema
- Manage RLS policies
- Execute SQL queries
- Introspect database structure

**Used By:**
- Studio table editor
- Studio SQL editor
- Schema migrations
- API documentation generation

**API Endpoints:**
```
GET  /tables        # List tables
POST /tables        # Create table
GET  /columns       # List columns
POST /query         # Execute SQL
GET  /policies      # List RLS policies
```

---

### 9. Edge Functions Runtime

**Image:** `supabase/edge-runtime:v1.45.2`
**Port:** 9000
**Purpose:** Serverless Deno functions

**Features:**
- Deno runtime (TypeScript/JavaScript)
- HTTP triggers
- Direct database access
- Import npm packages
- Environment variables
- Secrets management

**Example Function:**
```typescript
// functions/hello/index.ts
Deno.serve(async (req) => {
  const { name } = await req.json()

  return new Response(
    JSON.stringify({ message: `Hello ${name}!` }),
    { headers: { 'Content-Type': 'application/json' } }
  )
})
```

**Invoke:**
```bash
curl -X POST \
  https://project-ref.supabase.co/functions/v1/hello \
  -H 'Authorization: Bearer ANON_KEY' \
  -d '{"name":"World"}'
```

**Use Cases:**
- Custom API endpoints
- Webhooks processing
- Scheduled tasks (cron)
- Data processing
- Third-party integrations

---

### 10. Analytics (Logflare)

**Image:** `supabase/logflare:1.4.0`
**Port:** 4000
**Purpose:** Centralized logging and analytics

**Collects:**
- API request logs
- Database query logs
- Auth events
- Edge function logs
- Error logs

**Features:**
- Real-time log streaming
- Log search and filtering
- Query performance analytics
- Usage statistics
- Error tracking

**Backends:**
- PostgreSQL (default)
- Google BigQuery

**Studio Integration:**
- Logs viewer
- Query performance
- API usage charts

---

### 11. Vector (Log Collector)

**Image:** `timberio/vector:0.28.1-alpine`
**Port:** 9001
**Purpose:** Collect logs from all containers

**How It Works:**
1. Connects to Docker socket
2. Collects logs from all containers
3. Processes and enriches logs
4. Forwards to Analytics service

**Configuration:**
```yaml
sources:
  docker:
    type: docker_logs

transforms:
  parse:
    type: remap
    inputs: [docker]

sinks:
  logflare:
    type: http
    inputs: [parse]
    uri: http://analytics:4000
```

**Benefits:**
- Centralized logging
- Structured logs
- Log processing
- Multiple output destinations

---

### 12. Studio Dashboard

**Image:** `supabase/studio:20240422-5cf8f30`
**Port:** 3000
**Purpose:** Web-based project management interface

**Features:**
- **Table Editor** - Visual table management
- **SQL Editor** - Write and execute SQL
- **Auth Management** - User management
- **Storage Browser** - File management
- **API Documentation** - Auto-generated docs
- **Database Logs** - Query logs
- **Settings** - Project configuration

**Used For:**
- Database schema management
- User administration
- File uploads
- Monitoring
- Testing APIs

**Note:** In SupaManager, each project would have its own Studio instance, or use a shared Studio that connects to different projects.

---

## Service Dependencies

```
Prerequisites:
└─ Vector (starts first for logging)

Core Database:
└─ PostgreSQL
    ├─ GoTrue (auth)
    ├─ PostgREST (rest api)
    ├─ Realtime (subscriptions)
    ├─ Storage (file storage)
    ├─ Meta (db metadata)
    ├─ Edge Functions (serverless)
    └─ Analytics (logging)

Storage Stack:
└─ Storage API
    └─ ImgProxy (image processing)

API Gateway:
└─ Kong (waits for all services)
    └─ Routes to all backends

Management UI:
└─ Studio
    ├─ Depends on Analytics
    └─ Connects via Kong
```

---

## Resource Requirements per Project

### Minimum Configuration

```yaml
CPU:     2 cores
RAM:     2GB
Disk:    20GB
Network: 1Gbps
```

**Service Breakdown:**
| Service | CPU | RAM | Disk |
|---------|-----|-----|------|
| PostgreSQL | 0.5 | 512MB | 10GB |
| Kong | 0.2 | 128MB | 100MB |
| GoTrue | 0.1 | 64MB | 10MB |
| PostgREST | 0.2 | 128MB | 10MB |
| Realtime | 0.2 | 256MB | 10MB |
| Storage | 0.1 | 128MB | 5GB |
| ImgProxy | 0.1 | 128MB | 500MB |
| Meta | 0.1 | 64MB | 10MB |
| Functions | 0.2 | 256MB | 1GB |
| Analytics | 0.2 | 256MB | 2GB |
| Vector | 0.1 | 64MB | 500MB |
| Studio | 0.1 | 128MB | 500MB |

### Production Configuration

```yaml
CPU:     4 cores
RAM:     8GB
Disk:    100GB
Network: 10Gbps
```

**With:**
- Increased resource limits
- Autoscaling enabled
- Multiple replicas
- Backup volumes
- Monitoring

---

## Network Architecture

Each project has its own isolated Docker network:

```
Project: flying-rocket-network
├─ db (PostgreSQL) - 172.19.0.2
├─ kong - 172.19.0.3
├─ auth - 172.19.0.4
├─ rest - 172.19.0.5
├─ realtime - 172.19.0.6
├─ storage - 172.19.0.7
├─ imgproxy - 172.19.0.8
├─ meta - 172.19.0.9
├─ functions - 172.19.0.10
├─ analytics - 172.19.0.11
├─ vector - 172.19.0.12
└─ studio - 172.19.0.13

External Access:
├─ Kong HTTP: localhost:54321
├─ Kong HTTPS: localhost:54322
└─ Studio: localhost:54323
```

**Benefits:**
- Service isolation
- Network security
- DNS resolution
- No port conflicts

---

## Environment Variables per Project

Each project needs unique configuration:

**Secrets:**
```bash
JWT_SECRET=unique-per-project-32-chars-min
ANON_KEY=jwt-token-with-anon-role
SERVICE_ROLE_KEY=jwt-token-with-service-role
POSTGRES_PASSWORD=unique-strong-password
LOGFLARE_API_KEY=random-string
DASHBOARD_USERNAME=admin
DASHBOARD_PASSWORD=strong-password
```

**Project-Specific:**
```bash
PROJECT_REF=flying-rocket
POSTGRES_DB=postgres
POSTGRES_PORT=54320  # Unique per project
API_EXTERNAL_URL=https://flying-rocket.supamanager.io
SITE_URL=https://app.example.com
```

---

## Health Monitoring

Each service should have health checks:

```yaml
db:
  healthcheck:
    test: ["CMD", "pg_isready", "-U", "postgres"]
    interval: 10s
    timeout: 5s
    retries: 5

kong:
  healthcheck:
    test: ["CMD", "kong", "health"]
    interval: 10s

auth:
  healthcheck:
    test: ["CMD", "wget", "--spider", "http://localhost:9999/health"]
    interval: 10s
```

**Status Monitoring:**
- STARTING: Containers starting
- HEALTHY: All health checks passing
- UNHEALTHY: One or more services failing
- STOPPED: Containers stopped

---

## Related Documentation

- [Architecture Overview](Architecture-Overview) - Overall system design
- [Configuration Reference](Configuration-Reference) - Configuration options
- [Docker Networking](Docker-Networking) - Container networking

---

**Next:** Understand the [Development Guide](Development-Guide) to contribute.
