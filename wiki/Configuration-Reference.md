# Configuration Reference

Complete guide to configuring SupaManager.

---

## Configuration Files

SupaManager uses environment variables for configuration, stored in `.env` files:

1. **`docker-compose.yml`** - Service orchestration
2. **`supa-manager/.env`** - Backend API configuration
3. **`studio/.env`** - Frontend UI configuration

---

## SupaManager API Configuration

**File:** `supa-manager/.env`

### Database Configuration

```bash
# PostgreSQL connection string
DATABASE_URL=postgres://postgres:postgres@database:5432/supabase
```

**Format:** `postgres://[user]:[password]@[host]:[port]/[database]`

**Docker Service Name:** Use `database` (container name) not `localhost`

**External Database:** If using external PostgreSQL:
```bash
DATABASE_URL=postgres://user:pass@external-db.example.com:5432/supabase
```

### Authentication

```bash
# Allow new user registration
ALLOW_SIGNUP=true

# JWT secret for token signing (CHANGE FOR PRODUCTION!)
JWT_SECRET=your-super-secret-jwt-key-min-32-chars-long

# Encryption secret for sensitive data (CHANGE FOR PRODUCTION!)
ENCRYPTION_SECRET=your-super-secret-encryption-key-min-32-chars
```

**Security Requirements:**
- `JWT_SECRET`: Minimum 32 characters, random string
- `ENCRYPTION_SECRET`: Minimum 32 characters, different from JWT_SECRET
- Use a secure random generator:
  ```bash
  openssl rand -base64 48
  ```

**ALLOW_SIGNUP Options:**
- `true` - Anyone can create accounts
- `false` - Only existing users can login (admin must create accounts)

### Project Defaults

```bash
# Default PostgreSQL version for new projects
POSTGRES_DEFAULT_VERSION=14.2

# Default disk size in GB
POSTGRES_DISK_SIZE=10

# Docker image for PostgreSQL
POSTGRES_DOCKER_IMAGE=supabase/postgres
```

**Supported PostgreSQL Versions:**
- 13.x
- 14.x (default)
- 15.x

### Domain Configuration

```bash
# Base domain for project URLs
DOMAIN_BASE=supamanager.io

# Studio frontend URL
DOMAIN_STUDIO_URL=http://localhost:3000

# DNS webhook for dynamic DNS updates
DOMAIN_DNS_HOOK_URL=http://localhost:8081
DOMAIN_DNS_HOOK_KEY=mysecretkey
```

**Project URLs:** Projects will be accessible at:
```
https://[project-ref].supamanager.io
```

Example: `https://flying-rocket.supamanager.io`

**DNS Webhook:** (Optional) Called when projects are created to update DNS records.

### Service Version URL

```bash
# External service for checking Supabase component versions
SERVICE_VERSION_URL=https://supamanager.io/updates
```

**Purpose:** Checks for available versions of:
- PostgreSQL
- PostgREST
- GoTrue
- Kong
- Other Supabase services

---

## Studio UI Configuration

**File:** `studio/.env`

### API Connection

```bash
# Server-side API connection (use container name)
PLATFORM_PG_META_URL=http://supa-manager:8080/pg

# Client-side API connection (use localhost)
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_API_ADMIN_URL=http://localhost:8080
```

**Important:**
- `PLATFORM_PG_META_URL` uses `supa-manager` (Docker internal)
- `NEXT_PUBLIC_*` uses `localhost` (accessed from browser)

See [Docker Networking](Docker-Networking) for details.

### Frontend URLs

```bash
# Studio site URL
NEXT_PUBLIC_SITE_URL=http://localhost:3000

# Supabase API URL
NEXT_PUBLIC_SUPABASE_URL=http://localhost:8080

# GoTrue authentication URL
NEXT_PUBLIC_GOTRUE_URL=http://localhost:8080/auth
```

### Supabase Keys

```bash
# Anonymous key (placeholder for now)
NEXT_PUBLIC_SUPABASE_ANON_KEY=aaa.bbb.ccc
```

**Note:** This is a placeholder. Real keys will be generated per-project when provisioning is implemented.

### hCaptcha

```bash
# hCaptcha site key for bot protection
NEXT_PUBLIC_HCAPTCHA_SITE_KEY=10000000-ffff-ffff-ffff-000000000001
```

**Test Key:** The default is a hCaptcha test key (always passes).

**Production:** Get a real key from https://www.hcaptcha.com/

---

## Docker Compose Configuration

**File:** `docker-compose.yml`

### Service Configuration

```yaml
version: "3"
services:
  database:
    image: supabase/postgres:15.1.0.147
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres  # CHANGE FOR PRODUCTION
      POSTGRES_DB: supabase
```

**Production Changes:**
- Change `POSTGRES_PASSWORD`
- Use Docker secrets for sensitive data
- Enable SSL/TLS

### Port Mapping

```yaml
ports:
  - "HOST_PORT:CONTAINER_PORT"
```

**Default Ports:**
- `3000:3000` - Studio UI
- `8080:8080` - SupaManager API
- `5432:5432` - PostgreSQL

**Custom Ports:** Change if ports are already in use:
```yaml
ports:
  - "3001:3000"  # Studio on port 3001
  - "8081:8080"  # API on port 8081
```

### Resource Limits

```yaml
supa-manager:
  deploy:
    resources:
      limits:
        cpus: '1.0'
        memory: 512M
      reservations:
        cpus: '0.5'
        memory: 256M
```

### Health Checks

```yaml
database:
  healthcheck:
    test: ["CMD-SHELL", "pg_isready -U postgres"]
    interval: 5s
    timeout: 5s
    retries: 5
```

**Purpose:** Ensures services are ready before starting dependent services.

---

## Environment-Specific Configuration

### Development

**Characteristics:**
- Verbose logging
- Hot reload enabled
- Debug mode on
- ALLOW_SIGNUP=true
- Weak secrets OK

**Example:**
```bash
# .env.development
DATABASE_URL=postgres://postgres:postgres@localhost:5432/supabase
JWT_SECRET=dev-secret-not-secure
ALLOW_SIGNUP=true
LOG_LEVEL=debug
```

### Production

**Characteristics:**
- Minimal logging
- Strong secrets required
- ALLOW_SIGNUP=false (typically)
- SSL/TLS enabled
- Resource limits set

**Example:**
```bash
# .env.production
DATABASE_URL=postgres://user:$(cat /run/secrets/db_password)@db.prod.internal:5432/supabase
JWT_SECRET=$(cat /run/secrets/jwt_secret)
ALLOW_SIGNUP=false
LOG_LEVEL=info
```

**Production Checklist:**
- [ ] Change all default passwords
- [ ] Generate strong JWT_SECRET (48+ chars)
- [ ] Generate strong ENCRYPTION_SECRET
- [ ] Disable ALLOW_SIGNUP or add restrictions
- [ ] Set up SSL/TLS certificates
- [ ] Configure firewall rules
- [ ] Set up backups
- [ ] Configure monitoring
- [ ] Use Docker secrets for sensitive data

---

## Advanced Configuration

### Custom PostgreSQL Configuration

Mount custom `postgresql.conf`:

```yaml
database:
  volumes:
    - ./postgres-custom.conf:/etc/postgresql/postgresql.conf
  command: postgres -c config_file=/etc/postgresql/postgresql.conf
```

### Log Configuration

```bash
# Log level: debug, info, warn, error
LOG_LEVEL=info

# Log format: json, text
LOG_FORMAT=json
```

### CORS Configuration

Currently hardcoded in `supa-manager/api/api.go`:

```go
config := cors.DefaultConfig()
config.AllowOrigins = []string{"http://localhost:3000"}
config.AllowCredentials = true
```

**To modify:** Edit `api/api.go` and rebuild.

### Database Connection Pool

```bash
# Maximum connections (future)
DB_MAX_CONNECTIONS=100

# Idle connections (future)
DB_IDLE_CONNECTIONS=10
```

---

## Configuration Validation

### Checking Configuration

```bash
# View loaded environment variables
docker compose config

# Check specific service config
docker compose config supa-manager
```

### Testing Configuration

```bash
# Test database connection
docker exec supabase-manager-supa-manager-1 \
  psql "$DATABASE_URL" -c "SELECT 1"

# Test API health
curl http://localhost:8080/health

# Test Studio
curl http://localhost:3000
```

---

## Configuration Best Practices

### Security

1. **Never commit secrets to Git**
   ```bash
   # .gitignore
   .env
   .env.local
   .env.production
   ```

2. **Use Docker secrets in production**
   ```yaml
   secrets:
     jwt_secret:
       external: true
   ```

3. **Rotate secrets regularly**
   - JWT_SECRET: Every 90 days
   - ENCRYPTION_SECRET: Requires data re-encryption
   - Database passwords: Every 90 days

### Performance

1. **Database connection pooling**
   - Set appropriate pool size
   - Monitor connection usage

2. **Resource limits**
   - Prevent resource exhaustion
   - Ensure fair resource sharing

3. **Health checks**
   - Enable for all services
   - Set appropriate intervals

### Maintenance

1. **Environment-specific configs**
   - `.env.development`
   - `.env.staging`
   - `.env.production`

2. **Configuration as code**
   - Version control docker-compose.yml
   - Document all changes
   - Use variables for differences

3. **Configuration backup**
   - Back up .env files (encrypted)
   - Document configuration choices
   - Keep configuration changelog

---

## Troubleshooting Configuration

### Issue: Services can't communicate

**Symptom:** Studio can't reach API

**Check:**
```bash
# Verify network
docker network ls
docker network inspect supabase-manager_default

# Check DNS resolution
docker exec supabase-manager-studio-1 nslookup supa-manager
```

**Solution:** Ensure using correct hostnames in config.

### Issue: Database connection fails

**Symptom:** API logs show connection errors

**Check:**
```bash
# Verify DATABASE_URL format
docker compose exec supa-manager env | grep DATABASE_URL

# Test connection
docker exec supabase-manager-supa-manager-1 \
  psql "$DATABASE_URL" -c "SELECT version()"
```

**Solution:** Verify DATABASE_URL uses `database` not `localhost`.

### Issue: JWT token errors

**Symptom:** Authentication fails

**Check:**
```bash
# Verify JWT_SECRET is set
docker compose exec supa-manager env | grep JWT_SECRET
```

**Solution:** Ensure JWT_SECRET is minimum 32 characters.

---

## Configuration Templates

### Minimal Production Setup

```bash
# supa-manager/.env
DATABASE_URL=postgres://supabase:CHANGE_ME@database:5432/supabase
ALLOW_SIGNUP=false
JWT_SECRET=CHANGE_ME_MIN_32_CHARS_USE_OPENSSL_RAND
ENCRYPTION_SECRET=CHANGE_ME_DIFFERENT_FROM_JWT_SECRET
POSTGRES_DEFAULT_VERSION=15.1
POSTGRES_DISK_SIZE=20
POSTGRES_DOCKER_IMAGE=supabase/postgres
DOMAIN_STUDIO_URL=https://studio.yourdomain.com
DOMAIN_BASE=yourdomain.com
```

### High-Availability Setup

```yaml
# docker-compose.yml
services:
  database:
    image: supabase/postgres:15.1.0.147
    volumes:
      - db-data:/var/lib/postgresql/data
    deploy:
      replicas: 1
      restart_policy:
        condition: on-failure
        max_attempts: 3
    healthcheck:
      test: ["CMD", "pg_isready"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  db-data:
    driver: local
```

---

## Related Documentation

- [Docker Networking](Docker-Networking) - Container communication
- [Deployment](Deployment) - Production deployment
- [Troubleshooting](Troubleshooting) - Common issues

---

**Next:** Learn about [Docker Networking](Docker-Networking)
