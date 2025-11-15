# Troubleshooting Guide

Solutions to common problems and issues with SupaManager.

---

## Installation Issues

### Docker Not Found

**Symptom:**
```bash
bash: docker: command not found
```

**Solution:**
```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Add user to docker group
sudo usermod -aG docker $USER

# Log out and back in
```

**Verify:**
```bash
docker --version
docker compose version
```

### Permission Denied

**Symptom:**
```
permission denied while trying to connect to the Docker daemon socket
```

**Solution:**
```bash
# Add user to docker group
sudo usermod -aG docker $USER

# Log out and back in, or:
newgrp docker

# Verify
docker ps
```

### Port Already in Use

**Symptom:**
```
Error: bind: address already in use
```

**Check what's using the port:**
```bash
sudo lsof -i :3000  # Studio
sudo lsof -i :8080  # API
sudo lsof -i :5432  # PostgreSQL
```

**Solutions:**

**Option 1:** Stop the conflicting service
```bash
sudo systemctl stop postgresql  # If system PostgreSQL is running
```

**Option 2:** Change ports in `docker-compose.yml`
```yaml
services:
  studio:
    ports:
      - "3001:3000"  # Change 3000 to 3001
```

---

## Build Issues

### Studio Build Fails

**Symptom:**
```
error: patch failed to apply
```

**Possible Causes:**
1. Wrong Studio version
2. Network issues during download
3. Incomplete download

**Solutions:**

**1. Use correct version:**
```bash
cd studio
./build.sh v1.24.04 supa-manager/studio:v1.24.04 .env
```

**2. Clean and retry:**
```bash
cd studio
rm -rf code-*
./build.sh v1.24.04 supa-manager/studio:v1.24.04 .env
```

**3. Check internet connection:**
```bash
curl -I https://github.com
```

### Docker Image Build Fails

**Symptom:**
```
ERROR: failed to solve: process "/bin/sh -c go build" did not complete successfully
```

**Solutions:**

**1. Clear Docker cache:**
```bash
docker builder prune -af
```

**2. Check disk space:**
```bash
df -h
```

**3. Rebuild with no cache:**
```bash
docker compose build --no-cache
```

---

## Service Start Issues

### Database Won't Start

**Symptom:**
```
database-1 exited with code 1
```

**Check logs:**
```bash
docker compose logs database
```

**Common Causes:**

**1. Data corruption:**
```bash
# Remove volumes and restart
docker compose down -v
docker compose up -d
```

**2. Permission issues:**
```bash
# Check volume permissions
docker volume inspect supabase-manager_db-data
```

**3. Port conflict:**
```bash
# Check if PostgreSQL is already running
sudo lsof -i :5432
```

### API Won't Start

**Symptom:**
```
supa-manager-1 exited with code 1
```

**Check logs:**
```bash
docker compose logs supa-manager
```

**Common Issues:**

**1. Database not ready:**
```bash
# Wait for database health check
docker compose ps
```

**2. Missing environment variables:**
```bash
# Check .env file exists
ls -la supa-manager/.env

# Verify required variables
cat supa-manager/.env | grep JWT_SECRET
```

**3. Database connection error:**
```bash
# Test connection
docker exec supabase-manager-supa-manager-1 \
  psql "$DATABASE_URL" -c "SELECT 1"
```

### Studio Won't Start

**Symptom:**
```
studio-1 exited with code 1
```

**Check logs:**
```bash
docker compose logs studio
```

**Common Issues:**

**1. API not accessible:**
```bash
# Test API from studio container
docker exec supabase-manager-studio-1 \
  curl http://supa-manager:8080/health
```

**2. Missing environment variables:**
```bash
# Check studio .env
cat studio/.env | grep PLATFORM_PG_META_URL
```

---

## Runtime Issues

### Can't Login to Studio

**Symptom:** Login fails with "Invalid credentials"

**Check:**

**1. Verify default account exists:**
```bash
docker exec -it supabase-manager-database-1 \
  psql -U postgres -d supabase \
  -c "SELECT email FROM accounts;"
```

**2. Check API logs:**
```bash
docker compose logs supa-manager | grep -i auth
```

**3. Test API directly:**
```bash
curl -X POST http://localhost:8080/platform/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"haideritx@gmail.com","password":"NoAdmin@456"}'
```

**Solutions:**

**Create account manually:**
```bash
docker exec -it supabase-manager-database-1 psql -U postgres -d supabase
```
```sql
-- Check if account exists
SELECT * FROM accounts WHERE email = 'haideritx@gmail.com';

-- If not, restart supa-manager to run migrations
```

### Project Creation Fails

**Symptom:** "Failed to create new project: undefined"

**This is partially expected!** See [Current Status](Home#current-status).

**Check:**

**1. API logs:**
```bash
docker compose logs supa-manager | tail -20
```

**2. Database connection:**
```bash
docker exec -it supabase-manager-database-1 \
  psql -U postgres -d supabase -c "SELECT COUNT(*) FROM project;"
```

**3. Check if project was created:**
```bash
# View projects table
docker exec -it supabase-manager-database-1 \
  psql -U postgres -d supabase -c "SELECT project_ref, project_name, status FROM project;"
```

**Note:** Projects are created in database but provisioning is not implemented, so status shows "UNKNOWN".

### Studio Shows 404 Errors

**Symptom:** Various API endpoints return 404

**Common 404s (expected during development):**
- `/platform/organizations/{id}/usage` - Not implemented
- `/platform/notifications/summary` - Not implemented
- `/platform/projects/{ref}/settings` - Partially implemented

**Actual Problems:**

**Check API is running:**
```bash
curl http://localhost:8080/health
```

**Check routing:**
```bash
# Enable debug logs
docker compose logs supa-manager | grep -i route
```

---

## Network Issues

### Studio Can't Reach API

**Symptom:** "Network error" in Studio

**Check DNS resolution:**
```bash
# From studio container
docker exec supabase-manager-studio-1 nslookup supa-manager
```

**Should return:** 172.18.0.3 (or similar)

**Check network:**
```bash
docker network inspect supabase-manager_default
```

**Solution:** Ensure `PLATFORM_PG_META_URL` uses `supa-manager` not `localhost`:
```bash
# In studio/.env
PLATFORM_PG_META_URL=http://supa-manager:8080/pg
```

### Can't Access from Browser

**Symptom:** "This site can't be reached"

**Check:**

**1. Services are running:**
```bash
docker compose ps
```

**2. Ports are mapped:**
```bash
docker compose ps studio
# Should show 0.0.0.0:3000->3000/tcp
```

**3. Firewall allows connections:**
```bash
# Ubuntu/Debian
sudo ufw status
sudo ufw allow 3000/tcp
sudo ufw allow 8080/tcp
```

**4. Test from command line:**
```bash
curl http://localhost:3000
curl http://localhost:8080/health
```

---

## Database Issues

### Migrations Fail

**Symptom:**
```
ERROR: relation "migrations" does not exist
```

**Solution:**

**Restart API to run migrations:**
```bash
docker compose restart supa-manager
docker compose logs -f supa-manager
```

**Manually create migrations table:**
```bash
docker exec -it supabase-manager-database-1 psql -U postgres -d supabase
```
```sql
CREATE TABLE IF NOT EXISTS migrations (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    applied_at TIMESTAMPTZ DEFAULT NOW()
);
```

### Data Corruption

**Symptom:** Weird errors, inconsistent state

**Solution:**

**⚠️ This deletes all data!**
```bash
# Stop services
docker compose down

# Remove volumes
docker volume rm supabase-manager_db-data

# Start fresh
docker compose up -d
```

### Connection Pool Exhausted

**Symptom:** "too many clients already"

**Check connections:**
```bash
docker exec -it supabase-manager-database-1 psql -U postgres -d supabase
```
```sql
SELECT COUNT(*) FROM pg_stat_activity;
SELECT * FROM pg_stat_activity WHERE datname='supabase';
```

**Solution:**

**Increase max_connections:**
```sql
ALTER SYSTEM SET max_connections = 200;
SELECT pg_reload_conf();
```

---

## Performance Issues

### Slow Response Times

**Check:**

**1. Resource usage:**
```bash
docker stats
```

**2. Database performance:**
```sql
-- Slow queries
SELECT pid, query, state, now() - query_start AS duration
FROM pg_stat_activity
WHERE state != 'idle'
ORDER BY duration DESC;
```

**3. Container logs:**
```bash
docker compose logs --tail=100
```

**Solutions:**

**Add resource limits:**
```yaml
supa-manager:
  deploy:
    resources:
      limits:
        cpus: '2.0'
        memory: 1G
```

**Optimize database:**
```sql
VACUUM ANALYZE;
REINDEX DATABASE supabase;
```

### High Memory Usage

**Check:**
```bash
docker stats
free -h
```

**Solutions:**

**1. Restart services:**
```bash
docker compose restart
```

**2. Clear unused Docker resources:**
```bash
docker system prune -a
```

**3. Adjust container limits:**
```yaml
services:
  supa-manager:
    deploy:
      resources:
        limits:
          memory: 512M
```

---

## Debugging Tips

### Enable Verbose Logging

**API logs:**
```bash
# In supa-manager/.env (future feature)
LOG_LEVEL=debug
```

**Docker logs:**
```bash
# Follow logs in real-time
docker compose logs -f

# Specific service
docker compose logs -f supa-manager

# Last N lines
docker compose logs --tail=50 supa-manager
```

### Interactive Debugging

**Access containers:**
```bash
# API container
docker exec -it supabase-manager-supa-manager-1 sh

# Database
docker exec -it supabase-manager-database-1 psql -U postgres -d supabase

# Studio
docker exec -it supabase-manager-studio-1 sh
```

**Test connectivity:**
```bash
# From API to database
docker exec supabase-manager-supa-manager-1 \
  nc -zv database 5432

# From studio to API
docker exec supabase-manager-studio-1 \
  curl http://supa-manager:8080/health
```

### Inspect Network

```bash
# List networks
docker network ls

# Inspect network
docker network inspect supabase-manager_default

# Check DNS resolution
docker exec supabase-manager-studio-1 getent hosts supa-manager
```

---

## Getting Additional Help

### Before Asking for Help

Gather this information:

1. **System info:**
   ```bash
   uname -a
   docker --version
   docker compose version
   ```

2. **Service status:**
   ```bash
   docker compose ps
   ```

3. **Recent logs:**
   ```bash
   docker compose logs --tail=50
   ```

4. **Configuration:**
   ```bash
   # Sanitize sensitive data before sharing!
   cat docker-compose.yml
   ```

### Where to Get Help

1. **Check existing issues:** [GitHub Issues](https://github.com/YOUR_USERNAME/supabase-manager/issues)
2. **Search the wiki:** Use search function
3. **Ask in Discord:** [Harry's Discord](https://discord.gg/4k5HRe6YEp)
4. **Open new issue:** If problem persists

### Creating a Good Issue Report

Include:
- Description of the problem
- Steps to reproduce
- Expected vs actual behavior
- System information
- Relevant logs (sanitized)
- What you've already tried

---

## Known Issues

### Project Provisioning Not Implemented

**Symptom:** Projects show status "UNKNOWN" and no infrastructure is created

**Status:** This is expected. See [Roadmap](Roadmap) for implementation plans.

**Workaround:** None currently. This is the main feature being developed.

### Missing API Endpoints

**Symptom:** Some Studio features return 404

**Status:** Expected during development. These endpoints exist in the official Supabase API but are not yet implemented in SupaManager.

**Affected endpoints:**
- Organization usage statistics
- Notification summaries
- Some project settings

---

## Related Documentation

- [Configuration Reference](Configuration-Reference) - Check your configuration
- [Docker Networking](Docker-Networking) - Network issues
- [FAQ](FAQ) - Common questions

---

**Still stuck?** Ask in [Discord](https://discord.gg/4k5HRe6YEp) or open a [GitHub Issue](https://github.com/YOUR_USERNAME/supabase-manager/issues).
