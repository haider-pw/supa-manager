# GitHub Actions Workflows

This directory contains GitHub Actions workflows for automating deployment and testing.

## Workflows

### 1. Deploy to Production (`deploy-production.yml`)

**Trigger:** Automatically runs when code is pushed to the `main` branch, or can be manually triggered.

**What it does:**
1. Connects to production server via SSH
2. **Updates all environment variables** from GitHub Secrets
3. Pulls latest code from GitHub
4. Creates a backup of current deployment
5. Detects which services changed (API or Studio)
6. Rebuilds only the changed services
7. Restarts the updated services
8. Performs health checks
9. Reports success or failure

**Smart Features:**
- Centralized environment management via GitHub Secrets
- Automatic .env file updates on every deployment
- Only rebuilds services that have code changes
- Skips code rebuild if only environment changed (just restarts)
- Creates timestamped backups before deployment
- Performs health checks after deployment
- Stops deployment if health checks fail

## Required GitHub Secrets

You need to configure these secrets in your GitHub repository settings. The workflow uses **GitHub Environments** for better organization.

### SSH Connection (3 secrets)

| Secret Name | Value | Description |
|------------|-------|-------------|
| `PRODUCTION_HOST` | `182.191.91.226` | Production server public IP |
| `PRODUCTION_USER` | `haider` | SSH username |
| `PRODUCTION_PASSWORD` | Your SSH password | SSH password (or use SSH key) |

### Production Secrets (3 secrets)

| Secret Name | Example Value | Description |
|------------|---------------|-------------|
| `DATABASE_PASSWORD` | `ijKSTr78nr751qI6TjJfYVNeZA/cmcZ7LHlk9I2wliM=` | PostgreSQL password (auto URL-encoded) |
| `JWT_SECRET` | `Ev8o5bRYicouOZaiTPHIEe64xG...` | Secret for signing JWT tokens |
| `ENCRYPTION_SECRET` | `wxIJwYzb5ghyBMUEdYzlliflIg...` | Secret for encrypting sensitive data |

### Domain Configuration (2 secrets)

| Secret Name | Value | Description |
|------------|-------|-------------|
| `DOMAIN_BASE` | `supamanage.buzz` | Base domain for projects |
| `DOMAIN_STUDIO_URL` | `https://studio.supamanage.buzz` | Studio frontend URL |

### Service Configuration (4 secrets)

| Secret Name | Value | Description |
|------------|-------|-------------|
| `SERVICE_VERSION_URL` | `https://supamanager.io/updates` | Version service URL |
| `POSTGRES_DISK_SIZE` | `10` | Default disk size for projects (GB) |
| `POSTGRES_DEFAULT_VERSION` | `15.1` | Default PostgreSQL version |
| `POSTGRES_DOCKER_IMAGE` | `supabase/postgres` | PostgreSQL Docker image |

### Provisioning Configuration (5 secrets)

| Secret Name | Value | Description |
|------------|-------|-------------|
| `PROVISIONING_ENABLED` | `false` | Enable/disable project provisioning |
| `PROVISIONING_DOCKER_HOST` | `unix:///var/run/docker.sock` | Docker host for provisioning |
| `PROVISIONING_PROJECTS_DIR` | `/root/projects` | Directory for project data |
| `PROVISIONING_BASE_POSTGRES_PORT` | `5433` | Starting port for PostgreSQL instances |
| `PROVISIONING_BASE_KONG_HTTP_PORT` | `54321` | Starting port for Kong API gateway |

### Studio Configuration (8 secrets)

| Secret Name | Value | Description |
|------------|-------|-------------|
| `PLATFORM_PG_META_URL` | `https://www.supamanage.buzz/pg` | pg-meta API URL |
| `NEXT_PUBLIC_SITE_URL` | `https://studio.supamanage.buzz` | Studio site URL |
| `NEXT_PUBLIC_SUPABASE_URL` | `https://www.supamanage.buzz` | Main API URL |
| `NEXT_PUBLIC_SUPABASE_ANON_KEY` | `aaa.bbb.ccc` | Anonymous access key |
| `NEXT_PUBLIC_GOTRUE_URL` | `https://www.supamanage.buzz/auth` | Auth service URL |
| `NEXT_PUBLIC_API_URL` | `https://www.supamanage.buzz` | API base URL |
| `NEXT_PUBLIC_API_ADMIN_URL` | `https://www.supamanage.buzz` | Admin API URL |
| `NEXT_PUBLIC_HCAPTCHA_SITE_KEY` | `10000000-ffff-ffff-ffff-000000000001` | hCaptcha site key |

### Other Configuration (1 secret)

| Secret Name | Value | Description |
|------------|-------|-------------|
| `ALLOW_SIGNUP` | `true` | Enable/disable user registration |

## Total: 26 Secrets Required

### How to Add Secrets

#### Option 1: Using GitHub Environments (Recommended)

1. Go to your GitHub repository
2. Click **Settings** → **Environments**
3. Click **New environment** and name it `production`
4. Add all 26 secrets listed above

#### Option 2: Using Repository Secrets

1. Go to your GitHub repository
2. Click **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Add each secret one by one

### Quick Setup Script

You can use this script to get the current values from your production server:

```bash
# SSH into production server
ssh haider@182.191.91.226

# Extract current environment values
echo "=== Current Production Configuration ==="
echo ""
echo "DATABASE_PASSWORD=$(grep DATABASE_PASSWORD /opt/supamanage/.env | cut -d'=' -f2)"
echo "JWT_SECRET=$(grep JWT_SECRET /opt/supamanage/supa-manager/.env | cut -d'=' -f2)"
echo "ENCRYPTION_SECRET=$(grep ENCRYPTION_SECRET /opt/supamanage/supa-manager/.env | cut -d'=' -f2)"
echo ""
echo "Copy these values to GitHub Secrets"
```

## Manual Deployment

You can manually trigger a deployment:

1. Go to **Actions** tab in GitHub
2. Select **Deploy to Production** workflow
3. Click **Run workflow**
4. Select the `main` branch
5. Click **Run workflow** button

## Environment Updates

When you update any environment variable in GitHub Secrets:
1. The workflow will automatically update the .env files on the server
2. Services will be restarted with new configuration
3. No manual SSH access required!

**Example:** To change `ALLOW_SIGNUP` from `true` to `false`:
1. Update the `ALLOW_SIGNUP` secret in GitHub
2. Manually trigger the workflow (or push to main)
3. The API will restart with the new setting

## Monitoring Deployments

- All deployments are logged in the **Actions** tab
- You'll receive email notifications on deployment failures (if configured)
- Each deployment shows:
  - Environment variables updated
  - Which commit was deployed
  - Which services were rebuilt
  - Health check results
  - Deployment duration

## Security Notes

### Current Setup (Password-based)
- Currently uses password authentication
- All secrets stored as encrypted GitHub secrets
- Only accessible to workflow runs
- Secrets are never logged or exposed

### Network Configuration
- **Public IP**: 182.191.91.226 (used by GitHub Actions)
- **Private IP**: 192.168.100.14 (internal network)
- Production server is accessible from internet via public IP

### Recommended: Switch to SSH Keys (More Secure)

For better security, you should switch to SSH key authentication:

```bash
# On your local machine, generate SSH key for GitHub Actions
ssh-keygen -t ed25519 -C "github-actions@supamanage" -f github-actions-key

# Copy public key to production server (use public IP)
ssh-copy-id -i github-actions-key.pub haider@182.191.91.226

# In GitHub secrets, replace PRODUCTION_PASSWORD with:
# PRODUCTION_SSH_KEY = contents of github-actions-key (private key)
```

Then update the workflow to use `key` instead of `password`:

```yaml
- name: Deploy to Production Server
  uses: appleboy/ssh-action@v1.0.3
  with:
    host: ${{ secrets.PRODUCTION_HOST }}
    username: ${{ secrets.PRODUCTION_USER }}
    key: ${{ secrets.PRODUCTION_SSH_KEY }}  # Use SSH key instead
    port: 22
    # ... rest of configuration
```

## Rollback Procedure

If a deployment causes issues, you can rollback:

```bash
# SSH into production server (use public IP from internet, or private IP from local network)
ssh haider@182.191.91.226

# List available backups
ls -la /opt/supamanage-backups/

# Restore from backup
BACKUP_DATE="20250117-123456"  # Replace with your backup timestamp
cd /opt/supamanage
cp /opt/supamanage-backups/$BACKUP_DATE/.env .
cp /opt/supamanage-backups/$BACKUP_DATE/supa-manager.env supa-manager/.env
cp /opt/supamanage-backups/$BACKUP_DATE/studio.env studio/.env

# Checkout previous commit (optional)
git log --oneline  # Find the commit you want to rollback to
git checkout <previous-commit-sha>

# Rebuild and restart services
docker compose -f docker-compose.prod.yml build
docker compose -f docker-compose.prod.yml up -d
```

**Note:** If you're on the same local network as the server, you can also use the private IP `192.168.100.14` for SSH access.

## Troubleshooting

### Deployment fails with "Permission denied"
- Check that SSH credentials in GitHub secrets are correct
- Verify the production server is accessible from GitHub's IP ranges
- Make sure SSH is accessible on port 22 from the internet

### Deployment fails with "Connection refused" or "Timeout"
- Verify the public IP (182.191.91.226) is correct
- Check that port 22 is open in firewall/router for SSH
- Ensure port forwarding is configured if behind NAT

### Deployment succeeds but services aren't updated
- Check Docker logs: `docker compose -f docker-compose.prod.yml logs`
- Verify the code was actually pulled: `git log -1` on production server

### Health checks fail
- Services might need more time to start
- Check service logs for errors
- Verify Nginx configuration is correct

### Environment variables not updating
- Check that all required secrets are set in GitHub
- Verify the secret names match exactly (case-sensitive)
- Check deployment logs for errors during .env file creation

## Future Improvements

- [ ] Add automated testing before deployment
- [ ] Add Slack/Discord notifications
- [ ] Implement blue-green deployment
- [ ] Add database migration automation
- [ ] Set up staging environment deployment
- [ ] Add smoke tests after deployment
- [ ] Implement automatic rollback on health check failure
