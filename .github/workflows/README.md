# GitHub Actions Workflows

This directory contains GitHub Actions workflows for automating deployment and testing.

## Quick Setup

**Want to get started fast?** See [SETUP.md](./SETUP.md) for step-by-step instructions!

## Workflows

### Deploy to Production (`deploy-production.yml`)

**Trigger:** Automatically runs when code is pushed to the `main` branch, or can be manually triggered.

**What it does:**
1. Connects to production server via SSH
2. **Updates all environment variables** from GitHub Secrets & Variables
3. Pulls latest code from GitHub
4. Creates timestamped backups
5. Detects which services changed (API or Studio)
6. Rebuilds only the changed services
7. Performs health checks
8. Reports success or failure

**Smart Features:**
- ✅ Separates secrets (encrypted) from variables (plain text)
- ✅ Centralized environment management
- ✅ Automatic .env file updates
- ✅ Smart rebuilds based on code changes
- ✅ Health checks with automatic failure detection
- ✅ Timestamped backups before each deployment

## Required Configuration

The workflow uses **GitHub Environments** with two types of configuration:

### Secrets (Encrypted) - 6 items

These are sensitive values that are encrypted and never visible:

| Secret Name | Example Value | Description |
|------------|---------------|-------------|
| `PRODUCTION_HOST` | `<your-server-ip>` | Production server public IP |
| `PRODUCTION_USER` | `<your-username>` | SSH username |
| `PRODUCTION_PASSWORD` | `<your-password>` | SSH password (consider using SSH keys) |
| `DATABASE_PASSWORD` | `<generated-secret>` | PostgreSQL database password |
| `JWT_SECRET` | `<generated-secret>` | Secret for JWT token signing (min 64 chars) |
| `ENCRYPTION_SECRET` | `<generated-secret>` | Secret for data encryption (min 64 chars) |

### Variables (Plain Text) - 20 items

These are configuration values that are visible (not sensitive):

| Variable Name | Example Value | Description |
|--------------|---------------|-------------|
| `DOMAIN_BASE` | `yourdomain.com` | Base domain |
| `DOMAIN_STUDIO_URL` | `https://studio.yourdomain.com` | Studio URL |
| `SERVICE_VERSION_URL` | `https://placeholder.local/updates` | Version service URL |
| `POSTGRES_DISK_SIZE` | `10` | Default disk size (GB) |
| `POSTGRES_DEFAULT_VERSION` | `15.1` | PostgreSQL version |
| `POSTGRES_DOCKER_IMAGE` | `supabase/postgres` | PostgreSQL image |
| `PROVISIONING_ENABLED` | `false` | Enable provisioning |
| `PROVISIONING_DOCKER_HOST` | `unix:///var/run/docker.sock` | Docker host |
| `PROVISIONING_PROJECTS_DIR` | `/root/projects` | Projects directory |
| `PROVISIONING_BASE_POSTGRES_PORT` | `5433` | Base PostgreSQL port |
| `PROVISIONING_BASE_KONG_HTTP_PORT` | `54321` | Base Kong port |
| `PLATFORM_PG_META_URL` | `https://www.yourdomain.com/pg` | pg-meta URL |
| `NEXT_PUBLIC_SITE_URL` | `https://studio.yourdomain.com` | Studio site URL |
| `NEXT_PUBLIC_SUPABASE_URL` | `https://www.yourdomain.com` | Main API URL |
| `NEXT_PUBLIC_SUPABASE_ANON_KEY` | `aaa.bbb.ccc` | Anon key (placeholder) |
| `NEXT_PUBLIC_GOTRUE_URL` | `https://www.yourdomain.com/auth` | Auth URL |
| `NEXT_PUBLIC_API_URL` | `https://www.yourdomain.com` | API URL |
| `NEXT_PUBLIC_API_ADMIN_URL` | `https://www.yourdomain.com` | Admin API URL |
| `NEXT_PUBLIC_HCAPTCHA_SITE_KEY` | `10000000-ffff-ffff-ffff-000000000001` | hCaptcha key |
| `ALLOW_SIGNUP` | `true` | Allow user registration |

## How to Configure

### Step 1: Create Production Environment

1. Go to repository **Settings** → **Environments**
2. Click **"New environment"**
3. Name it: `production`
4. Click **"Configure environment"**

### Step 2: Add Secrets

In the production environment page:
1. Scroll to **"Environment secrets"**
2. Click **"Add secret"** for each of the 6 secrets listed above
3. Copy values from your current `.env` files on the server

### Step 3: Add Variables

In the production environment page:
1. Scroll to **"Environment variables"**
2. Click **"Add variable"** for each of the 20 variables listed above
3. Use your actual domain and configuration values

### Get Current Values from Server

SSH into your server to get current production values:

```bash
ssh <your-user>@<your-server-ip>

# Get secrets
grep DATABASE_PASSWORD /opt/supamanage/.env
grep JWT_SECRET /opt/supamanage/supa-manager/.env
grep ENCRYPTION_SECRET /opt/supamanage/supa-manager/.env
```

**Security Note:** Never commit these values to git or share them publicly!

## How It Works

### On Every Push to Main:

1. **Environment Sync**: All `.env` files regenerated from GitHub
2. **Code Sync**: Pulls latest code
3. **Change Detection**: Checks what changed
4. **Smart Rebuild**: Only rebuilds modified services
5. **Health Checks**: Verifies services respond correctly
6. **Notification**: Reports success/failure

### Environment Updates Only:

If you only change GitHub Secrets/Variables (no code changes):
- `.env` files are regenerated
- Services are restarted
- No rebuilding (faster)

## Benefits

- 🔐 **Secure**: Secrets encrypted, variables visible
- ⚡ **Fast**: Only rebuilds what changed (2-5 min)
- 🎯 **Smart**: Detects code vs config changes
- 📦 **Safe**: Backups before each deployment
- 🛡️ **Reliable**: Health checks prevent bad deployments
- 📊 **Visible**: See all deployments in Actions tab
- 🔄 **Easy**: Update config without SSH

## Manual Deployment

Trigger a deployment manually:

1. Go to **Actions** tab
2. Select **"Deploy to Production"**
3. Click **"Run workflow"**
4. Select `main` branch
5. Click **"Run workflow"** button

## Monitoring

- All deployments logged in **Actions** tab
- Email notifications on failures (if configured)
- Each deployment shows:
  - Environment updates
  - Code changes
  - Services rebuilt
  - Health check results
  - Duration

## Rollback

If a deployment causes issues:

```bash
# SSH into server
ssh <your-user>@<your-server-ip>

# List backups
ls -la /opt/supamanage-backups/

# Restore from backup
BACKUP="20250117-123456"  # Your backup timestamp
cd /opt/supamanage
cp /opt/supamanage-backups/$BACKUP/.env .
cp /opt/supamanage-backups/$BACKUP/supa-manager.env supa-manager/.env
cp /opt/supamanage-backups/$BACKUP/studio.env studio/.env

# Or checkout previous commit
git log --oneline
git checkout <commit-sha>

# Restart services
docker compose -f docker-compose.prod.yml up -d
```

## Security Best Practices

### Current Setup ✅
- **Secrets**: Encrypted by GitHub, never exposed
- **Variables**: Plain text, visible to repo collaborators
- **Separation**: Only truly sensitive data is encrypted
- **No hardcoded credentials**: All sensitive values in GitHub Secrets

### Recommended Improvements

**Switch to SSH Keys** (more secure than passwords):

```bash
# Generate SSH key
ssh-keygen -t ed25519 -C "github-actions" -f github-actions-key

# Copy to server
ssh-copy-id -i github-actions-key.pub <your-user>@<your-server-ip>

# In GitHub secrets:
# Delete: PRODUCTION_PASSWORD
# Add: PRODUCTION_SSH_KEY (contents of private key)
```

Update workflow:
```yaml
with:
  host: ${{ secrets.PRODUCTION_HOST }}
  username: ${{ secrets.PRODUCTION_USER }}
  key: ${{ secrets.PRODUCTION_SSH_KEY }}  # Instead of password
```

## Troubleshooting

### "Missing environment variable"
- Verify all 6 secrets are configured
- Verify all 20 variables are configured
- Check for typos (case-sensitive!)

### "SSH connection failed"
- Verify public IP is correct
- Check port 22 is accessible
- Test SSH manually first

### "Health checks failed"
- Services may need more startup time
- Check logs: `docker compose -f docker-compose.prod.yml logs`
- Verify Nginx configuration

### Deployment succeeds but changes not visible
- Check correct branch was deployed (main)
- Verify services restarted: `docker compose ps`
- Clear browser cache

## Future Improvements

- [ ] Add automated testing before deployment
- [ ] Implement blue-green deployment
- [ ] Add Slack/Discord notifications
- [ ] Database migration automation
- [ ] Staging environment setup
- [ ] Smoke tests after deployment
- [ ] Auto-rollback on failure
