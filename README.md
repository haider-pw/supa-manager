# SupaManager

**Originally created by [Harry Bairstow](https://twitter.com/TheHarryET)**
**Enhanced and maintained by [Syed Haider Hassan](https://github.com/haider-pw)**

Manage self-hosted Supabase instances using the Supabase Studio.

> **This is an enhanced fork** of the [original SupaManager project](https://github.com/TheHarryET/supa-manager) with improvements and comprehensive documentation.

## Key Improvements in This Fork

- ✅ **Comprehensive Documentation** - Complete wiki with 25+ detailed guides
- ✅ **Bug Fixes** - Fixed API endpoints, JWT handling, and database migrations
- ✅ **Security Improvements** - Removed exposed credentials, improved validation
- ✅ **Enhanced Developer Experience** - Added MCP server support, debugging guides
- ✅ **Better Error Handling** - Improved API response formats and error messages
- ✅ **Production Ready** - Deployment guides, monitoring, and backup procedures

📖 See [CHANGELOG.md](wiki/Changelog.md) for detailed improvements.

---

> [!WARNING]
> **Active Development Status**
> This project is in active development. Currently, the management API and Studio UI are functional, but **project provisioning is not yet implemented**. Projects can be created in the database but will not spin up actual Supabase infrastructure.

This project provides a management API to work with Supabase Studio, allowing you to manage multiple Supabase projects through a web interface.

> [!NOTE]
> We provide a patched version of Supabase Studio (v1.24.04) that works with this API. The build script is located in the `studio/` folder.

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Detailed Setup](#detailed-setup)
- [Accessing the Application](#accessing-the-application)
- [Project Structure](#project-structure)
- [Configuration](#configuration)
- [Current Status](#current-status)
- [Troubleshooting](#troubleshooting)
- [Documentation](#documentation)

---

## Prerequisites

Before you begin, ensure you have the following installed on your Ubuntu/Linux machine:

- **Docker** (version 20.10 or higher)
- **Docker Compose** (v2 or higher - comes with Docker Desktop)
- **Git** (for cloning the repository)

### Installing Docker

If you don't have Docker installed:

```bash
# Update package list
sudo apt-get update

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Add your user to the docker group (to run without sudo)
sudo usermod -aG docker $USER

# Log out and back in for group changes to take effect
```

Verify Docker is installed:
```bash
docker --version
docker compose version
```

---

## Quick Start

Get up and running in 5 minutes:

```bash
# 1. Clone the repository
git clone <repository-url>
cd supabase-manager

# 2. Build the Studio image (takes 5-10 minutes first time)
cd studio
./build.sh v1.24.04 supa-manager/studio:v1.24.04 .env
cd ..

# 3. Start all services
docker compose up -d

# 4. Wait for services to be healthy (30-60 seconds)
docker compose ps

# 5. Open Studio in your browser
# http://localhost:3000
```

**Note:** You'll need to sign up for an account on first access through the Studio UI.

---

## Detailed Setup

### Step 1: Clone the Repository

```bash
git clone <repository-url>
cd supabase-manager
```

### Step 2: Build the Patched Studio Image

The Studio UI requires a patched version to work with the management API.

```bash
# Navigate to the studio directory
cd studio

# Build the patched Studio image (this will take 5-10 minutes)
./build.sh v1.24.04 supa-manager/studio:v1.24.04 .env

# Return to the root directory
cd ..
```

**What this does:**
- Downloads Supabase Studio source code (v1.24.04)
- Applies custom patches to work with the supa-manager API
- Builds a Docker image tagged as `supa-manager/studio:v1.24.04`

### Step 3: Review Configuration (Optional)

The docker-compose.yml is already configured with sensible defaults. You can optionally review the configuration files:

**Main configuration files:**
- `docker-compose.yml` - Service definitions
- `supa-manager/.env` - Backend API configuration
- `studio/.env` - Studio frontend configuration

For production deployments, you should change:
- `JWT_SECRET` in `supa-manager/.env`
- `ENCRYPTION_SECRET` in `supa-manager/.env`
- Database passwords

### Step 4: Start Services

Start all services using Docker Compose:

```bash
# Start in detached mode (background)
docker compose up -d

# Or start with logs visible (useful for debugging)
docker compose up
```

This will start three services:
1. **PostgreSQL Database** (port 5432) - Management database
2. **Supa-Manager API** (port 8080) - Backend API
3. **Studio UI** (port 3000) - Web interface

### Step 5: Verify Services Are Running

Check that all services are healthy:

```bash
docker compose ps
```

You should see:
```
NAME                            STATUS
supabase-manager-database-1     Up (healthy)
supabase-manager-supa-manager-1 Up
supabase-manager-studio-1       Up
```

Check logs if any service has issues:
```bash
# View all logs
docker compose logs

# View logs for specific service
docker compose logs supa-manager
docker compose logs studio
docker compose logs database
```

---

## Accessing the Application

### Studio Web Interface

Open your browser and navigate to:
```
http://localhost:3000
```

**First Time Setup:**
1. Open the Studio UI at http://localhost:3000
2. Click "Sign Up" to create your admin account
3. Fill in your email and password
4. Start managing your Supabase projects

### API Endpoints

The management API is available at:
```
http://localhost:8080
```

**Example API calls:**

```bash
# Health check
curl http://localhost:8080/health

# Get API version (requires authentication)
curl http://localhost:8080/platform/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Database Access

PostgreSQL is exposed on the default port:
```
Host: localhost
Port: 5432
Database: supabase
Username: postgres
Password: postgres
```

Connect using `psql`:
```bash
docker exec -it supabase-manager-database-1 psql -U postgres -d supabase
```

---

## Project Structure

```
supabase-manager/
├── docker-compose.yml          # Service orchestration
├── supa-manager/               # Backend API (Go)
│   ├── api/                    # API handlers
│   ├── conf/                   # Configuration
│   ├── database/               # Database queries (sqlc)
│   ├── migrations/             # SQL migrations
│   ├── Dockerfile              # Backend container build
│   └── .env                    # Backend configuration
├── studio/                     # Frontend UI (Next.js)
│   ├── build.sh                # Studio build script
│   ├── patch.sh                # Patch application script
│   └── .env                    # Studio configuration
├── version-service/            # Version tracking service
├── dns-service/                # DNS management service
├── helm/                       # Kubernetes deployment charts
└── README.md                   # This file
```

---

## Configuration

### Environment Variables (supa-manager/.env)

Key configuration options:

```bash
# Database connection
DATABASE_URL=postgres://postgres:postgres@database:5432/supabase

# Authentication
ALLOW_SIGNUP=true
JWT_SECRET=secret                    # Change for production!
ENCRYPTION_SECRET=secret             # Change for production!

# Project defaults
POSTGRES_DISK_SIZE=10
POSTGRES_DEFAULT_VERSION=14.2
POSTGRES_DOCKER_IMAGE=supabase/postgres

# Studio integration
DOMAIN_STUDIO_URL=http://localhost:3000
DOMAIN_BASE=supamanager.io

# DNS webhook (for dynamic DNS updates)
DOMAIN_DNS_HOOK_URL=http://localhost:8081
DOMAIN_DNS_HOOK_KEY=mysecretkey
```

### Environment Variables (studio/.env)

```bash
# API connection
PLATFORM_PG_META_URL=http://supa-manager:8080/pg
NEXT_PUBLIC_API_URL=http://localhost:8080

# Frontend URLs
NEXT_PUBLIC_SITE_URL=http://localhost:3000
NEXT_PUBLIC_GOTRUE_URL=http://localhost:8080/auth
```

---

## Current Status

### ✅ What's Working

- User authentication (signup/login)
- Organization creation and management
- Project metadata creation
- Studio UI integration
- API endpoints for project management
- Database migrations

### ⚠️ What's Not Working (Yet)

- **Dynamic Supabase project provisioning** - Projects are created in the database but no actual Supabase infrastructure is spun up
- Project status tracking (projects remain in "UNKNOWN" status)
- Project lifecycle management (pause/resume/delete)
- Real JWT key generation (currently returns placeholder "a.b.c")

### 🚧 Roadmap

See the [PROJECT_ANALYSIS.md](PROJECT_ANALYSIS.md) and [SUPABASE_ARCHITECTURE.md](SUPABASE_ARCHITECTURE.md) documents for detailed information about the planned implementation.

**Upcoming features:**
- Phase 2: Design provisioning approach
- Phase 3: Implement Docker Compose provisioning for Supabase projects
- Phase 4: Project lifecycle management
- Phase 5: Kubernetes support

---

## Troubleshooting

### Services won't start

**Check if ports are already in use:**
```bash
sudo lsof -i :3000  # Studio
sudo lsof -i :8080  # API
sudo lsof -i :5432  # PostgreSQL
```

**Solution:** Stop services using those ports or change ports in docker-compose.yml

### Studio build fails

**Error:** "Patches failed to apply"

**Solution:** Ensure you're using the correct Studio version:
```bash
cd studio
./build.sh v1.24.04 supa-manager/studio:v1.24.04 .env
```

### Database migration errors

**Error:** "migrations table does not exist"

**Solution:** Restart the supa-manager service:
```bash
docker compose restart supa-manager
```

### Cannot login to Studio

**Check supa-manager logs:**
```bash
docker compose logs supa-manager
```

**Verify database is healthy:**
```bash
docker compose ps database
```

### "Failed to create new project: undefined"

This is expected! The project creation API works, but provisioning is not yet implemented. The project is created in the database but no Supabase infrastructure is started. See [Current Status](#current-status) above.

---

## Managing Services

### Stop services
```bash
docker compose down
```

### Stop and remove volumes (delete all data)
```bash
docker compose down -v
```

### Restart a specific service
```bash
docker compose restart supa-manager
docker compose restart studio
```

### View real-time logs
```bash
docker compose logs -f
```

### Rebuild after code changes
```bash
docker compose up -d --build
```

---

## Documentation

For more detailed technical documentation, see:

- **[Wiki Home](https://github.com/haider-pw/supa-manager/wiki)** - Complete documentation wiki
- **[CLAUDE.md](CLAUDE.md)** - Quick reference for development
- **[PROJECT_ANALYSIS.md](PROJECT_ANALYSIS.md)** - Complete codebase analysis
- **[SUPABASE_ARCHITECTURE.md](SUPABASE_ARCHITECTURE.md)** - Full Supabase service architecture
- **[PHASE_1_SUMMARY.md](PHASE_1_SUMMARY.md)** - Analysis summary and next steps

---

## About This Fork

### Upstream Project

This fork is based on Harry Bairstow's original SupaManager project:
- **Original Repository:** https://github.com/TheHarryET/supa-manager
- **Author:** [@TheHarryET](https://twitter.com/TheHarryET)

### Contributing

Contributions to this fork are welcome! See [Contributing Guide](https://github.com/haider-pw/supa-manager/wiki/Contributing) for details.

If you find improvements that would benefit the upstream project, consider contributing to Harry's original repository as well.

---

## License

**Original Project:**
Copyright (C) 2024 Harry Bairstow

**This Fork:**
Copyright (C) 2025 Syed Haider Hassan

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.

### Attribution

This project is a fork of [Harry Bairstow's SupaManager](https://github.com/TheHarryET/supa-manager).
All modifications and enhancements in this fork are also licensed under GPL v3.