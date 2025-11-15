# Supabase Project Architecture Requirements

## Overview
This document describes the complete architecture needed to provision a fully functional Supabase project. This is based on the official Supabase docker-compose configuration found in `studio/code-v1.24.04/docker/docker-compose.yml`.

## Complete Service Stack

A single Supabase project requires **12 interconnected services**:

### 1. PostgreSQL Database (db)
- **Image**: `supabase/postgres:15.1.1.41`
- **Port**: Dynamic (configurable via POSTGRES_PORT)
- **Purpose**: Main project database with extensions
- **Key Features**:
  - Logical replication enabled
  - Pre-configured roles (anon, authenticated, service_role, supabase_admin)
  - Extensions: pg_graphql, pgsodium, pgvector, etc.
  - Custom initialization scripts for realtime, webhooks, JWT
- **Environment Variables**:
  - `POSTGRES_PASSWORD`: Database password
  - `POSTGRES_DB`: Database name
  - `JWT_SECRET`: For generating auth tokens
  - `JWT_EXPIRY`: Token expiration time

### 2. Kong API Gateway (kong)
- **Image**: `kong:2.8.1`
- **Ports**: 8000 (HTTP), 8443 (HTTPS)
- **Purpose**: Reverse proxy routing all API requests
- **Key Features**:
  - DB-less mode (declarative config)
  - JWT authentication via key-auth plugin
  - CORS support
  - Routes all services through single endpoint
- **Configuration File**: `kong.yml` (declarative routing config)
- **Routing Map**:
  ```
  /auth/v1/*      -> http://auth:9999/*
  /rest/v1/*      -> http://rest:3000/*
  /realtime/v1/*  -> ws://realtime:4000/socket/*
  /storage/v1/*   -> http://storage:5000/*
  /functions/v1/* -> http://functions:9000/*
  /analytics/v1/* -> http://analytics:4000/*
  /pg/*           -> http://meta:8080/*
  /*              -> http://studio:3000/* (dashboard)
  ```

### 3. GoTrue Auth Service (auth)
- **Image**: `supabase/gotrue:v2.149.0`
- **Port**: 9999
- **Purpose**: Authentication and user management
- **Key Features**:
  - Email/password authentication
  - OAuth providers support
  - JWT token generation
  - User confirmation and password recovery
- **Environment Variables**:
  - `GOTRUE_DB_DATABASE_URL`: Connection to project database
  - `GOTRUE_JWT_SECRET`: Shared secret with other services
  - `GOTRUE_SITE_URL`: Frontend URL for redirects
  - `API_EXTERNAL_URL`: Public API URL
  - `SMTP_*`: Email configuration for confirmations

### 4. PostgREST API (rest)
- **Image**: `postgrest/postgrest:v12.0.1`
- **Port**: 3000
- **Purpose**: Automatic REST API from PostgreSQL schema
- **Key Features**:
  - Auto-generated REST endpoints from tables
  - Row-level security enforcement
  - JWT role-based access
- **Environment Variables**:
  - `PGRST_DB_URI`: Database connection
  - `PGRST_DB_SCHEMAS`: Exposed schemas (public, storage, graphql_public)
  - `PGRST_DB_ANON_ROLE`: Default anonymous role
  - `PGRST_JWT_SECRET`: For token verification

### 5. Realtime Service (realtime)
- **Image**: `supabase/realtime:v2.28.32`
- **Port**: 4000
- **Purpose**: WebSocket connections for real-time data
- **Key Features**:
  - PostgreSQL logical replication
  - Broadcast messages
  - Presence tracking
  - Change data capture (CDC)
- **Environment Variables**:
  - `DB_*`: PostgreSQL connection details
  - `API_JWT_SECRET`: Token verification
  - `SECRET_KEY_BASE`: Encryption key

### 6. Storage API (storage)
- **Image**: `supabase/storage-api:v1.0.6`
- **Port**: 5000
- **Purpose**: File storage and management
- **Key Features**:
  - S3-compatible API
  - File upload/download
  - Image transformation (via imgproxy)
  - Row-level security for buckets
- **Environment Variables**:
  - `ANON_KEY`, `SERVICE_KEY`: JWT tokens
  - `POSTGREST_URL`: For authorization checks
  - `DATABASE_URL`: Storage metadata
  - `FILE_SIZE_LIMIT`: Max upload size (default: 50MB)
  - `STORAGE_BACKEND`: file or s3

### 7. ImgProxy (imgproxy)
- **Image**: `darthsim/imgproxy:v3.8.0`
- **Port**: 5001
- **Purpose**: On-the-fly image transformation
- **Key Features**:
  - Image resizing
  - Format conversion
  - WebP support
  - ETags for caching

### 8. Postgres Meta (meta)
- **Image**: `supabase/postgres-meta:v0.80.0`
- **Port**: 8080
- **Purpose**: Database metadata and schema management
- **Key Features**:
  - Table/column introspection
  - Schema modifications via API
  - Used by Studio for database management
- **Environment Variables**:
  - `PG_META_DB_*`: PostgreSQL connection details

### 9. Edge Functions Runtime (functions)
- **Image**: `supabase/edge-runtime:v1.45.2`
- **Port**: 9000
- **Purpose**: Serverless Deno functions
- **Key Features**:
  - Deno runtime for edge functions
  - Direct database access
  - HTTP triggers
- **Environment Variables**:
  - `JWT_SECRET`: Token verification
  - `SUPABASE_URL`, `SUPABASE_ANON_KEY`: For Supabase client
  - `SUPABASE_DB_URL`: Direct DB access

### 10. Analytics/Logging (analytics)
- **Image**: `supabase/logflare:1.4.0`
- **Port**: 4000
- **Purpose**: Centralized logging and analytics
- **Key Features**:
  - Log aggregation from all services
  - Query logs viewer in Studio
  - Postgres or BigQuery backend
- **Environment Variables**:
  - `DB_*`: PostgreSQL connection for logs storage
  - `LOGFLARE_API_KEY`: Authentication key
  - `LOGFLARE_SINGLE_TENANT`: true for self-hosted
  - `POSTGRES_BACKEND_URL`: Store logs in Postgres

### 11. Vector Log Collector (vector)
- **Image**: `timberio/vector:0.28.1-alpine`
- **Port**: 9001
- **Purpose**: Log collection from Docker containers
- **Key Features**:
  - Collects logs from all containers
  - Forwards to analytics service
- **Environment Variables**:
  - `LOGFLARE_API_KEY`: For forwarding logs
- **Volumes**:
  - Docker socket access for log collection

### 12. Studio Dashboard (studio)
- **Image**: `supabase/studio:20240422-5cf8f30` (or custom patched version)
- **Port**: 3000
- **Purpose**: Web-based management interface
- **Key Features**:
  - Table editor
  - SQL editor
  - Auth user management
  - Storage browser
  - API documentation
- **Environment Variables**:
  - `STUDIO_PG_META_URL`: Connection to pg-meta
  - `POSTGRES_PASSWORD`: Database access
  - `SUPABASE_URL`: Kong gateway URL
  - `SUPABASE_ANON_KEY`, `SUPABASE_SERVICE_KEY`: API keys
  - `LOGFLARE_API_KEY`: For logs viewer

## Service Dependencies

```
Database (db)
    ├─> Auth (depends on db)
    ├─> PostgREST (depends on db)
    ├─> Realtime (depends on db)
    ├─> Storage (depends on db + rest)
    ├─> Meta (depends on db)
    ├─> Analytics (depends on db)
    └─> Edge Functions (depends on db)

Vector (must start before db for logging)

Kong (depends on all services being ready)
    └─> Routes all traffic

Studio (depends on analytics + kong)
```

## Network Architecture

All services communicate on an internal Docker network. Kong is the only service exposed externally (ports 8000/8443), acting as the single entry point for all API requests.

```
External Request
    ↓
Kong Gateway (:8000)
    ↓
Routes based on path:
    /auth/v1/*      → GoTrue (:9999)
    /rest/v1/*      → PostgREST (:3000)
    /realtime/v1/*  → Realtime (:4000)
    /storage/v1/*   → Storage (:5000)
    /functions/v1/* → Functions (:9000)
    /analytics/v1/* → Analytics (:4000)
    /pg/*           → Meta (:8080)
    /*              → Studio (:3000)
```

## Required Environment Variables per Project

Each project needs these unique configuration values:

### Secrets (Generated)
- `JWT_SECRET`: Minimum 32 characters, used for signing tokens
- `ANON_KEY`: JWT with role "anon", generated from JWT_SECRET
- `SERVICE_ROLE_KEY`: JWT with role "service_role", generated from JWT_SECRET
- `POSTGRES_PASSWORD`: Strong password for database
- `LOGFLARE_API_KEY`: Random string for analytics auth
- `DASHBOARD_USERNAME`, `DASHBOARD_PASSWORD`: Basic auth for Studio

### Project-Specific
- `POSTGRES_DB`: Unique database name per project
- `POSTGRES_PORT`: Unique port per project (to avoid conflicts)
- `PROJECT_REF`: Short identifier (e.g., "upending-expectoration")
- `API_EXTERNAL_URL`: Public URL for the project (e.g., https://upending-expectoration.supamanager.io)
- `SITE_URL`: Frontend URL for auth redirects

### Static Configuration
- `POSTGRES_HOST`: Always "db" (container name)
- `KONG_HTTP_PORT`: 8000 (internal)
- `PGRST_DB_SCHEMAS`: "public,storage,graphql_public"
- `JWT_EXPIRY`: 3600 (1 hour)

## Storage Requirements

### Volumes per Project
- **Database data**: Persistent volume for PostgreSQL data (~10GB minimum)
- **Storage files**: Persistent volume for uploaded files (~10GB minimum)
- **DB config**: Named volume for pgsodium keys
- **Function code**: Mounted from host or persistent volume

## Port Allocation Strategy

Each project needs unique ports to avoid conflicts:

1. **PostgreSQL**: Base port + project_id (e.g., 5432, 5433, 5434...)
2. **Kong HTTP**: Base port + project_id (e.g., 54321, 54322, 54323...)
3. **Kong HTTPS**: Base port + project_id (e.g., 54322, 54323, 54324...)

All internal ports (auth:9999, rest:3000, etc.) can remain the same because they're isolated in separate Docker networks per project.

## JWT Key Generation

The ANON_KEY and SERVICE_ROLE_KEY are JWTs generated with this payload:

**ANON_KEY:**
```json
{
  "role": "anon",
  "iss": "supabase",
  "iat": 1641769200,
  "exp": 1799535600
}
```

**SERVICE_ROLE_KEY:**
```json
{
  "role": "service_role",
  "iss": "supabase",
  "iat": 1641769200,
  "exp": 1799535600
}
```

Both signed with `JWT_SECRET` using HS256 algorithm.

## Provisioning Approaches

### Option A: Docker Compose (Recommended for MVP)

**Pros:**
- Simpler implementation
- Easier debugging
- No K8s cluster required
- Works on single server

**Cons:**
- Limited scalability
- Single point of failure
- Manual resource management

**Implementation:**
1. Generate docker-compose.yml from template
2. Use unique project name for network isolation
3. Use Docker SDK for Go to manage containers
4. Store compose file in database or filesystem
5. Track container states

### Option B: Kubernetes

**Pros:**
- Better scalability
- Built-in health checks
- Resource limits and autoscaling
- High availability

**Cons:**
- Complex implementation
- Requires K8s cluster
- Steeper learning curve

**Implementation:**
1. Use Helm chart (already exists in `/helm/studio-chart/`)
2. Kubernetes Go client
3. Custom Resource Definitions (CRDs)
4. Operator pattern for project lifecycle

## Recommendations for supa-manager

Based on the codebase analysis, here are the key changes needed:

1. **Add Docker SDK dependency**:
   ```go
   github.com/docker/docker v24.0.0
   github.com/docker/compose/v2 v2.20.0
   ```

2. **Create template engine** for docker-compose generation

3. **Implement provisioning service**:
   - Generate JWT secrets
   - Create docker-compose.yml from template
   - Start containers via Docker API
   - Monitor container health
   - Update project status in database

4. **Add status tracking**:
   - PROVISIONING → Creating containers
   - STARTING → Waiting for health checks
   - ACTIVE → All services healthy
   - PAUSED → Containers stopped
   - FAILED → Provisioning error

5. **Implement lifecycle endpoints**:
   - POST `/platform/projects/{ref}/pause` → docker compose stop
   - POST `/platform/projects/{ref}/resume` → docker compose start
   - DELETE `/platform/projects/{ref}` → docker compose down -v

6. **Store project metadata**:
   - Docker network name
   - Container IDs
   - Port mappings
   - Volume paths
   - compose file location

## Next Steps

See Phase 2 in the TODO list for designing and implementing the provisioning system.
