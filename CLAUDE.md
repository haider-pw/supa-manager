# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

SupaManager (by Harry Bairstow) is a system for managing self-hosted Supabase instances through the Supabase Studio UI. It provides a fully functioning API compatible with Supabase's production platform API, enabling users to manage multiple Supabase instances from a single interface.

This is a mono-repo containing:
- **supa-manager**: Main API service that mimics Supabase's platform API
- **version-service**: Service for tracking and distributing container version information
- **dns-example-service**: Template for implementing DNS record management hooks
- **studio**: Patched Supabase Studio frontend with custom API endpoints
- **helm/studio-chart**: Kubernetes Helm chart for deploying the entire stack

## Architecture

### Service Communication Flow

1. **Supabase Studio** (frontend) → **supa-manager** (main API)
   - Studio makes requests to `/platform/*`, `/projects/*`, `/auth/*`, etc.
   - API handles user authentication, organization management, and project lifecycle

2. **supa-manager** → **dns-example-service** (DNS hook)
   - When projects are created/deleted, supa-manager calls the DNS hook service
   - Hook updates DNS records to point project subdomains to the correct infrastructure

3. **supa-manager** → **version-service**
   - Queries latest container versions for Supabase services (GoTrue, PostgREST, etc.)
   - Used during project creation to deploy current versions

### Main API Structure (supa-manager)

The API is built with Gin framework and follows a handler-per-endpoint pattern:
- **api/api.go**: Core router setup, middleware, and route definitions
- **api/*.go**: Individual handler files (e.g., `getProfile.go`, `postPlatformProjects.go`)
- **database/**: sqlc-generated database access layer
- **queries/**: SQL query definitions for sqlc
- **migrations/**: Database schema migrations
- **conf/**: Configuration management and migration runner
- **utils/**: Shared utilities
- **permisions/**: Permission/authorization logic

### Database Schema

Core tables (defined in `supa-manager/migrations/00_init.sql`):
- **accounts**: User accounts with GoTrue ID, email, password hash
- **organizations**: Organizations/teams that own projects
- **organization_membership**: Links accounts to organizations with roles
- **projects**: Supabase project instances with configuration
- **migrations**: Tracks applied database migrations

Uses PostgreSQL with pgx/v5 driver and sqlc for type-safe queries.

## Development Commands

### Building & Running Services

**Main API (supa-manager):**
```bash
cd supa-manager
cp .env.example .env  # Configure environment variables
go run main.go        # Runs on :8080 by default
```

**Version Service:**
```bash
cd version-service
cp .env.example .env
go run main.go        # Runs on configured LISTEN_ADDRESS
```

**DNS Example Service:**
```bash
cd dns-example-service
cp .env.example .env
go run main.go
```

**Docker Compose (full stack):**
```bash
docker-compose up  # Starts database and studio
```

### Code Generation

**Generate database code with sqlc:**
```bash
cd supa-manager
sqlc generate  # Reads sqlc.yaml, generates code in database/
```

```bash
cd version-service
sqlc generate
```

sqlc configuration in `sqlc.yaml` points to:
- Queries: `./queries/*.sql`
- Schema: `./migrations/*.sql`
- Output: `./database/` (for supa-manager) or inline (for version-service)

### Studio Patching

The Supabase Studio frontend requires patching to work with supa-manager's API:

```bash
cd studio
./patch.sh [branch]                          # Apply patches to studio codebase
./build.sh [branch] [docker-tag] [.env-file] # Build patched Docker image
```

Patches are stored in `studio/patches/` and applied incrementally. To create new patches:
1. Run `./patch.sh` on a clean branch
2. Make changes to the studio code
3. Create patch file: `git diff > patches/XX-description.patch`

### Helm Deployment

```bash
cd helm/studio-chart
helm install studio-release . -f values.yaml
helm upgrade studio-release . -f values.yaml
```

Chart includes deployments for:
- Studio (patched frontend)
- Supa-manager API
- Version service
- DNS example service
- Container registry (for studio builds)

## Key Configuration

### supa-manager Environment Variables

- `DATABASE_URL`: PostgreSQL connection string
- `ALLOW_SIGNUP`: Enable/disable user registration
- `JWT_SECRET`: Secret for signing authentication tokens
- `ENCRYPTION_SECRET`: Secret for encrypting sensitive data
- `SERVICE_VERSION_URL`: URL of version service (default: https://supamanager.io/updates)
- `POSTGRES_*`: Default settings for created Supabase instances
- `DOMAIN_STUDIO_URL`: Studio frontend URL
- `DOMAIN_BASE`: Base domain for project URLs (e.g., `project-ref.supamanager.io`)
- `DOMAIN_DNS_HOOK_URL`: URL of DNS update service
- `DOMAIN_DNS_HOOK_KEY`: Authentication key for DNS service

### version-service Configuration

- `DATABASE_URL`: PostgreSQL connection string
- `PUSHING_ACCOUNTS`: Comma-separated GitHub usernames allowed to push version updates
- `LISTEN_ADDRESS`: Server listen address (default: `0.0.0.0:8081`)

Version updates require SSH signature verification from authorized GitHub accounts.

## Important Notes

- **Authentication**: Uses JWT tokens with GoTrue-compatible IDs. API validates tokens via `Authorization: Bearer <token>` header.
- **Migrations**: Automatically run on startup by checking `migrations` table and applying pending migrations from `migrations/` directory.
- **Project Lifecycle**: Creating a project triggers DNS hook, provisions database, and configures services.
- **sqlc Usage**: All database queries should be defined in `queries/*.sql` and generated with `sqlc generate`. Do not write raw SQL in Go code.
- **CORS**: API allows all origins (`*`) for Studio integration.
- **Permissions**: Organization membership and roles control access to projects and organizations.

## Testing

No test framework is currently configured. When adding tests:
- Use Go's standard `testing` package
- Run tests with `go test ./...` from service directory
- Consider adding integration tests for API endpoints
