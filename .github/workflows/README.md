# GitHub Actions Workflows

This directory contains GitHub Actions workflows for automating deployment and testing.

## Workflows

### 1. Deploy to Production (`deploy-production.yml`)

**Trigger:** Automatically runs when code is pushed to the `main` branch, or can be manually triggered.

**What it does:**
1. Connects to production server via SSH
2. Pulls latest code from GitHub
3. Creates a backup of current deployment
4. Detects which services changed (API or Studio)
5. Rebuilds only the changed services
6. Restarts the updated services
7. Performs health checks
8. Reports success or failure

**Smart Features:**
- Only rebuilds services that have code changes
- Skips deployment if no changes detected
- Creates timestamped backups before deployment
- Performs health checks after deployment
- Stops deployment if health checks fail

## Required GitHub Secrets

You need to configure these secrets in your GitHub repository settings:

1. **PRODUCTION_HOST**: The production server public IP address
   - Value: `182.191.91.226`

2. **PRODUCTION_USER**: SSH username for production server
   - Value: `haider`

3. **PRODUCTION_PASSWORD**: SSH password for production server
   - Value: Your SSH password

### How to Add Secrets

1. Go to your GitHub repository
2. Click **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Add each secret with the name and value listed above

## Manual Deployment

You can manually trigger a deployment:

1. Go to **Actions** tab in GitHub
2. Select **Deploy to Production** workflow
3. Click **Run workflow**
4. Select the `main` branch
5. Click **Run workflow** button

## Monitoring Deployments

- All deployments are logged in the **Actions** tab
- You'll receive email notifications on deployment failures (if configured)
- Each deployment shows:
  - Which commit was deployed
  - Which services were rebuilt
  - Health check results
  - Deployment duration

## Security Notes

### Current Setup (Password-based)
- Currently uses password authentication
- Password stored as encrypted GitHub secret
- Only accessible to workflow runs

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

# Navigate to installation directory
cd /opt/supamanage

# Checkout previous commit
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

## Future Improvements

- [ ] Add automated testing before deployment
- [ ] Add Slack/Discord notifications
- [ ] Implement blue-green deployment
- [ ] Add database migration automation
- [ ] Set up staging environment deployment
