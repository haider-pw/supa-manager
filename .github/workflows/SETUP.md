# Quick Setup Guide for GitHub Actions CI/CD

## Step 1: Create Production Environment

1. Go to: https://github.com/haider-pw/supa-manager/settings/environments
2. Click **"New environment"**
3. Name it: `production`
4. Click **"Configure environment"**

## Step 2: Add Secrets (Encrypted)

In the production environment, go to **"Environment secrets"** and add these **6 secrets**:

| Secret Name | Where to Find Value | Example |
|------------|---------------------|---------|
| `PRODUCTION_HOST` | Your production server's public IP | `203.0.113.10` |
| `PRODUCTION_USER` | Your SSH username | `ubuntu` or `admin` |
| `PRODUCTION_PASSWORD` | Your SSH password | `YourSecurePassword123` |
| `DATABASE_PASSWORD` | From server: `/opt/supamanage/.env` | Generated secret |
| `JWT_SECRET` | From server: `/opt/supamanage/supa-manager/.env` | Generated secret (64+ chars) |
| `ENCRYPTION_SECRET` | From server: `/opt/supamanage/supa-manager/.env` | Generated secret (64+ chars) |

**To get values from your server:**
```bash
ssh <your-user>@<your-server-ip>

# Get the secret values
grep DATABASE_PASSWORD /opt/supamanage/.env
grep JWT_SECRET /opt/supamanage/supa-manager/.env
grep ENCRYPTION_SECRET /opt/supamanage/supa-manager/.env
```

## Step 3: Add Variables (Plain Text)

In the production environment, go to **"Environment variables"** and add these **20 variables**:

**Replace `yourdomain.com` with your actual domain!**

```bash
DOMAIN_BASE=yourdomain.com
DOMAIN_STUDIO_URL=https://studio.yourdomain.com
SERVICE_VERSION_URL=https://placeholder.local/updates
POSTGRES_DISK_SIZE=10
POSTGRES_DEFAULT_VERSION=15.1
POSTGRES_DOCKER_IMAGE=supabase/postgres
PROVISIONING_ENABLED=false
PROVISIONING_DOCKER_HOST=unix:///var/run/docker.sock
PROVISIONING_PROJECTS_DIR=/root/projects
PROVISIONING_BASE_POSTGRES_PORT=5433
PROVISIONING_BASE_KONG_HTTP_PORT=54321
PLATFORM_PG_META_URL=https://www.yourdomain.com/pg
NEXT_PUBLIC_SITE_URL=https://studio.yourdomain.com
NEXT_PUBLIC_SUPABASE_URL=https://www.yourdomain.com
NEXT_PUBLIC_SUPABASE_ANON_KEY=aaa.bbb.ccc
NEXT_PUBLIC_GOTRUE_URL=https://www.yourdomain.com/auth
NEXT_PUBLIC_API_URL=https://www.yourdomain.com
NEXT_PUBLIC_API_ADMIN_URL=https://www.yourdomain.com
NEXT_PUBLIC_HCAPTCHA_SITE_KEY=10000000-ffff-ffff-ffff-000000000001
ALLOW_SIGNUP=true
```

## Step 4: Merge the PR

1. Go to: https://github.com/haider-pw/supa-manager/pull/9
2. Review the changes
3. Click **"Merge pull request"**
4. Confirm merge

## Step 5: Watch the Deployment

1. Go to: https://github.com/haider-pw/supa-manager/actions
2. You'll see the deployment workflow running
3. Click on it to see real-time logs
4. Wait for it to complete (2-5 minutes)

## Step 6: Verify Production

1. Visit: https://www.yourdomain.com (your main domain)
2. Visit: https://studio.yourdomain.com (your studio domain)
3. Both should be working!

---

## Testing Future Deployments

Once setup is complete, every time you push to `main` branch, the deployment will run automatically.

To test it:
1. Make a small change to any file
2. Commit and push to `main`
3. Watch the Actions tab
4. See your changes go live automatically!

## Manual Deployment

You can also manually trigger a deployment:
1. Go to Actions tab
2. Select "Deploy to Production" workflow
3. Click "Run workflow"
4. Select `main` branch
5. Click "Run workflow" button

## Troubleshooting

### "Error: Missing environment variable"
- Make sure all 6 secrets and 20 variables are added
- Check for typos in variable names
- Variables are case-sensitive!

### "SSH connection failed"
- Verify `PRODUCTION_HOST` is your server's public IP
- Check that port 22 is accessible from the internet
- Test SSH manually: `ssh <user>@<ip>`
- Verify SSH credentials are correct

### "Health checks failed"
- Services might need more time to start
- Check deployment logs for errors
- SSH into server and check container logs:
  ```bash
  docker compose -f docker-compose.prod.yml logs
  ```

## Security Notes

**Secrets** (encrypted, never visible):
- SSH credentials
- Database password  
- JWT and encryption secrets

**Variables** (plain text, visible to everyone with repo access):
- URLs and domains
- Configuration values
- Port numbers

**Important:** Never commit secrets to git or share them publicly! Store them only in GitHub Secrets.

## Network Requirements

- Your production server must have a public IP accessible from the internet
- Port 22 (SSH) must be open for GitHub Actions to connect
- If behind a router/firewall, configure port forwarding for SSH

## Next Steps

After successful deployment:
1. Set up SSH keys instead of password (more secure)
2. Configure automated backups
3. Set up monitoring/alerts
4. Create a staging environment
5. Add automated tests before deployment
