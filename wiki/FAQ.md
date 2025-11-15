# Frequently Asked Questions (FAQ)

Common questions about SupaManager.

---

## General Questions

### What is SupaManager?

SupaManager is a self-hosted management platform for Supabase. It provides a web interface (using the official Supabase Studio) to create and manage multiple Supabase projects from a single control panel.

### Is this an official Supabase project?

No, SupaManager is a community project by [Harry Bairstow](https://twitter.com/TheHarryET). It integrates with the official Supabase Studio but is not affiliated with Supabase Inc.

### Is it production-ready?

**Not yet.** The management API and Studio UI work, but the core feature (dynamic project provisioning) is not yet implemented. Projects can be created in the database but won't automatically spin up Supabase infrastructure.

See the [Roadmap](Roadmap) for development status.

### What's the difference between this and Supabase Cloud?

| Feature | Supabase Cloud | SupaManager |
|---------|----------------|-------------|
| Hosting | Managed by Supabase | Self-hosted |
| Cost | Pay per project | Infrastructure costs only |
| Setup | Instant | Requires setup |
| Maintenance | Handled by Supabase | You manage |
| Customization | Limited | Full control |
| Support | Official support | Community support |

**Use Supabase Cloud if:** You want a managed service with official support.
**Use SupaManager if:** You need self-hosted, full control, or have specific compliance requirements.

---

## Installation & Setup

### What are the system requirements?

**Minimum:**
- 2 CPU cores
- 4GB RAM
- 20GB disk space
- Docker 20.10+
- Docker Compose v2+

**Recommended:**
- 4 CPU cores
- 8GB RAM
- 50GB SSD
- Ubuntu 22.04 or similar Linux

### Can I run this on Windows or Mac?

Yes, via Docker Desktop! SupaManager works on:
- Linux (recommended)
- macOS (via Docker Desktop)
- Windows (via Docker Desktop with WSL2)

Performance may be better on Linux for production use.

### Do I need to know Go or Next.js to use this?

**No!** You only need to know:
- Basic command line usage
- Docker basics
- How to edit `.env` files

**Development:** Yes, you'll need Go and Next.js knowledge to contribute.

### How long does setup take?

- **Studio build:** 5-10 minutes (first time only)
- **Services start:** 30-60 seconds
- **Total:** ~10-15 minutes from clone to running

---

## Features & Capabilities

### Can SupaManager create Supabase projects?

**Currently:** Projects are created in the database with metadata, but no actual Supabase infrastructure is provisioned. Status shows "UNKNOWN".

**Future:** Yes! Automatic provisioning of complete Supabase stack (12 services) is planned. See [Roadmap](Roadmap).

### What Supabase features are supported?

Currently, only the management layer is implemented:
- ✅ User authentication
- ✅ Organization management
- ✅ Project metadata
- ✅ Studio UI integration

**After provisioning is implemented:**
- ✅ PostgreSQL database
- ✅ Authentication (GoTrue)
- ✅ REST API (PostgREST)
- ✅ Real-time subscriptions
- ✅ Storage
- ✅ Edge Functions
- ✅ And more!

See [Supabase Stack](Supabase-Stack) for complete list.

### How many projects can I create?

**Current:** Unlimited in database (but they don't do anything yet).

**Future:** Limited by your server resources. Each project needs:
- ~2GB RAM
- ~20GB disk space
- Unique ports

**Example:** A server with 16GB RAM could run ~6-8 projects.

### Can I use custom domains?

**Planned:** Yes! Projects will be accessible at:
```
https://project-ref.yourdomain.com
```

You'll need to:
1. Configure `DOMAIN_BASE` in config
2. Set up wildcard DNS (*.yourdomain.com)
3. Configure SSL certificates

See [Configuration Reference](Configuration-Reference) for details.

---

## Technical Questions

### What database does SupaManager use?

**Management Database:** PostgreSQL 15.1 (stores users, organizations, projects)

**Project Databases:** Each Supabase project will have its own PostgreSQL instance (when provisioning is implemented).

### How are passwords stored?

Using **Argon2id** (industry-standard secure hashing):
- Memory: 64MB
- Iterations: 3
- Random salt per password
- Constant-time comparison

See [Architecture Overview](Architecture-Overview#security-model) for details.

### How does authentication work?

**JWT tokens** (JSON Web Tokens):
1. User logs in with email/password
2. Server validates credentials
3. Server generates JWT token signed with `JWT_SECRET`
4. Client includes token in subsequent requests
5. Server validates token signature

Tokens expire after the configured time (default: 1 hour).

### Is data encrypted?

**In Transit:** HTTPS when properly configured
**At Rest:** Database encryption when enabled
**Secrets:** Encrypted using `ENCRYPTION_SECRET`

**Production:** Enable PostgreSQL encryption and use HTTPS/TLS.

### Can I use an external PostgreSQL database?

Yes! Change `DATABASE_URL` in `supa-manager/.env`:

```bash
DATABASE_URL=postgres://user:password@external-host:5432/database
```

### How do I back up data?

**Management Database:**
```bash
# Backup
docker exec supabase-manager-database-1 \
  pg_dump -U postgres supabase > backup.sql

# Restore
cat backup.sql | docker exec -i supabase-manager-database-1 \
  psql -U postgres supabase
```

**Project Data:** (Future) Each project will have its own backup strategy.

See [Backup & Recovery](Backup-Recovery) for details.

---

## Security Questions

### Is SupaManager secure?

The codebase follows security best practices:
- ✅ Argon2 password hashing
- ✅ JWT authentication
- ✅ Parameterized SQL queries (sqlc)
- ✅ CORS configuration
- ✅ Input validation

**However:**
- ⚠️ Change default secrets before production!
- ⚠️ Enable HTTPS/TLS
- ⚠️ Keep software updated
- ⚠️ Follow security best practices

See [Deployment](Deployment) for production security checklist.

### Should I change the default secrets?

**YES! Immediately!** Generate new secrets:

```bash
# Generate JWT_SECRET
openssl rand -base64 48

# Generate ENCRYPTION_SECRET
openssl rand -base64 48
```

Update in `supa-manager/.env`:
```bash
JWT_SECRET=<generated-secret-1>
ENCRYPTION_SECRET=<generated-secret-2>
```

### Can I disable user registration?

Yes! In `supa-manager/.env`:

```bash
ALLOW_SIGNUP=false
```

This prevents new user registration. Existing users can still log in.

### How do I create admin accounts?

Currently, manually via database:

```bash
docker exec -it supabase-manager-database-1 psql -U postgres -d supabase
```

```sql
-- Insert account (password will be hashed)
-- Note: Manual password hashing required, or use registration endpoint once then disable
```

**Better:** Use the registration endpoint to create accounts, then disable registration.

---

## Troubleshooting Questions

### Why can't I login to Studio?

**Common causes:**
1. API not running: `docker compose ps`
2. Wrong credentials: Ensure you're using the account you created during signup
3. Database not initialized: `docker compose restart supa-manager`
4. Token expired: Clear browser cookies and log in again

See [Troubleshooting - Can't Login](Troubleshooting#cant-login-to-studio).

### Why does project creation show "undefined" error?

This is **partially expected** during development:
- Project IS created in database
- But provisioning is not implemented
- So no actual infrastructure is created
- Status shows "UNKNOWN"

Check if project exists:
```bash
docker exec -it supabase-manager-database-1 \
  psql -U postgres -d supabase \
  -c "SELECT project_ref, project_name FROM project;"
```

### Port 3000 is already in use, what do I do?

**Option 1:** Stop the conflicting service
```bash
sudo lsof -i :3000
# Kill the process or:
sudo systemctl stop <service>
```

**Option 2:** Change the port in `docker-compose.yml`
```yaml
studio:
  ports:
    - "3001:3000"  # Use port 3001 instead
```

Then access Studio at `http://localhost:3001`

### Services keep restarting, how do I fix it?

**Check logs:**
```bash
docker compose logs --tail=50
```

**Common causes:**
1. Database not ready - Wait for health check
2. Config error - Check `.env` files
3. Resource exhaustion - Check `docker stats`

See [Troubleshooting](Troubleshooting) for detailed solutions.

---

## Development Questions

### How can I contribute?

See the [Contributing Guide](Contributing) for:
- Setting up development environment
- Code standards
- Pull request process
- Issue reporting

### What technologies are used?

**Backend:**
- Go 1.21
- Gin framework
- PostgreSQL
- sqlc (type-safe queries)

**Frontend:**
- Next.js 12
- React
- TypeScript

**Infrastructure:**
- Docker
- Docker Compose
- (Future) Kubernetes

See [Architecture Overview](Architecture-Overview) for details.

### Where is the code?

GitHub repository: (Your repository URL)

Key directories:
- `supa-manager/` - Backend API (Go)
- `studio/` - Frontend UI (Next.js patches)
- `helm/` - Kubernetes charts
- `version-service/` - Version tracking
- `dns-service/` - DNS management

### How do I run in development mode?

**Backend:**
```bash
cd supa-manager
go run main.go
```

**Frontend:** Use the Studio image or build from source.

See [Development Guide](Development-Guide) for complete setup.

### Can I add new features?

Yes! Contributions are welcome. See:
1. [Roadmap](Roadmap) - What's planned
2. [Contributing](Contributing) - How to contribute
3. [GitHub Issues](https://github.com/YOUR_USERNAME/supabase-manager/issues) - Open issues

---

## Deployment Questions

### Can I use this in production?

**Not recommended yet** until provisioning is implemented. The management layer works, but you can't actually create functioning Supabase projects.

**Future:** Yes, with proper configuration and security hardening.

### How do I deploy to production?

See the [Deployment Guide](Deployment) for:
- Production configuration
- Security hardening
- SSL/TLS setup
- Monitoring
- Backup strategy

### Can I use Kubernetes instead of Docker Compose?

**Yes!** Helm charts are available in the `helm/` directory.

**Note:** This is for deploying SupaManager itself. Each Supabase project can also be deployed via Kubernetes (planned feature).

### How do I update SupaManager?

```bash
# Pull latest changes
git pull origin main

# Rebuild images
docker compose build

# Restart services
docker compose up -d
```

**Note:** Check changelog for breaking changes or migration requirements.

---

## Comparison Questions

### SupaManager vs Supabase CLI?

| Feature | Supabase CLI | SupaManager |
|---------|-------------|-------------|
| Purpose | Local development | Multi-project management |
| Interface | Command line | Web UI (Studio) |
| Projects | One at a time | Multiple projects |
| Use Case | Development | Production self-hosting |

**Both can be used together!** CLI for development, SupaManager for hosting.

### SupaManager vs Manual Docker Setup?

**Manual Setup:**
- More control
- More complex
- Requires manual configuration per project

**SupaManager:**
- Automated (when provisioning is done)
- Web interface
- Easier to manage multiple projects
- Less flexibility

**Choose SupaManager if:** You want ease of use and multi-project management.

---

## Licensing Questions

### What license is SupaManager under?

**GNU General Public License v3.0 (GPLv3)**

### Can I use this commercially?

**Yes!** The GPLv3 allows commercial use.

**Requirements:**
- Disclose source code
- Include license notice
- State changes made
- Keep same license

### Can I modify the code?

**Yes!** You're free to modify the source code.

**Requirements:**
- Keep the GPLv3 license
- Disclose your changes
- Make source code available

### Do I need to share my modifications?

**Yes, if you distribute** the modified software.

**No, if you only use it internally** (e.g., for your own company's infrastructure).

---

## Future Plans

### When will provisioning be implemented?

See the [Roadmap](Roadmap) for timeline and phases.

**Current Status:** Phase 1 complete (analysis)
**Next:** Phase 2 (design) and Phase 3 (implementation)

### What features are planned?

- ✅ Dynamic Supabase project provisioning
- ✅ Project lifecycle management (pause/resume/delete)
- ✅ Real-time status monitoring
- ✅ Resource usage tracking
- ✅ Kubernetes support
- ✅ Backup and restore
- ✅ Custom domains
- ✅ SSL certificate management

See [Roadmap](Roadmap) for complete list.

### Will there be a hosted version?

No official plans currently. This is designed to be self-hosted.

However, community members might offer hosting services.

---

## Getting Help

### Where can I get help?

1. **Check the wiki** - Most questions answered here
2. **Search GitHub Issues** - Someone may have asked already
3. **Discord** - [Harry's Discord Server](https://discord.gg/4k5HRe6YEp)
4. **Twitter** - [@TheHarryET](https://twitter.com/TheHarryET)
5. **Open an Issue** - For bugs or feature requests

### How do I report a bug?

1. **Check** [Known Issues](Troubleshooting#known-issues)
2. **Search** existing [GitHub Issues](https://github.com/YOUR_USERNAME/supabase-manager/issues)
3. **If new**, open an issue with:
   - Description
   - Steps to reproduce
   - Expected vs actual behavior
   - System info
   - Logs (sanitized)

### How do I request a feature?

Open a [GitHub Issue](https://github.com/YOUR_USERNAME/supabase-manager/issues) with:
- Clear description of the feature
- Use case / why it's needed
- Proposed implementation (optional)
- Willingness to contribute (if applicable)

---

## Still Have Questions?

- **Check other wiki pages** - [Home](Home)
- **Ask in Discord** - [Join here](https://discord.gg/4k5HRe6YEp)
- **Open an issue** - [GitHub Issues](https://github.com/YOUR_USERNAME/supabase-manager/issues)

---

**Didn't find your answer?** Please open an issue or ask in Discord - your question might help others too!
