# Quick Setup Guide for GitHub Actions CI/CD

## Step 1: Create Production Environment

1. Go to: https://github.com/haider-pw/supa-manager/settings/environments
2. Click **"New environment"**
3. Name it: `production`
4. Click **"Configure environment"**

## Step 2: Add Secrets (Encrypted)

In the production environment, go to **"Environment secrets"** and add these **6 secrets**:

```bash
PRODUCTION_HOST=182.191.91.226
PRODUCTION_USER=haider
PRODUCTION_PASSWORD=Admin@786
DATABASE_PASSWORD=ijKSTr78nr751qI6TjJfYVNeZA/cmcZ7LHlk9I2wliM=
JWT_SECRET=Ev8o5bRYicouOZaiTPHIEe64xGCv52QQa4dp8XcLv1MEIomP1psZ6jjAebEbhepX
ENCRYPTION_SECRET=wxIJwYzb5ghyBMUEdYzlliflIg54n6bRPiYTH+lO08x/wUQ+bZVNS5xmNhtb7qQp
```

## Step 3: Add Variables (Plain Text)

In the production environment, go to **"Environment variables"** and add these **20 variables**:

```bash
DOMAIN_BASE=supamanage.buzz
DOMAIN_STUDIO_URL=https://studio.supamanage.buzz
SERVICE_VERSION_URL=https://placeholder.local/updates
POSTGRES_DISK_SIZE=10
POSTGRES_DEFAULT_VERSION=15.1
POSTGRES_DOCKER_IMAGE=supabase/postgres
PROVISIONING_ENABLED=false
PROVISIONING_DOCKER_HOST=unix:///var/run/docker.sock
PROVISIONING_PROJECTS_DIR=/root/projects
PROVISIONING_BASE_POSTGRES_PORT=5433
PROVISIONING_BASE_KONG_HTTP_PORT=54321
PLATFORM_PG_META_URL=https://www.supamanage.buzz/pg
NEXT_PUBLIC_SITE_URL=https://studio.supamanage.buzz
NEXT_PUBLIC_SUPABASE_URL=https://www.supamanage.buzz
NEXT_PUBLIC_SUPABASE_ANON_KEY=aaa.bbb.ccc
NEXT_PUBLIC_GOTRUE_URL=https://www.supamanage.buzz/auth
NEXT_PUBLIC_API_URL=https://www.supamanage.buzz
NEXT_PUBLIC_API_ADMIN_URL=https://www.supamanage.buzz
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

1. Visit: https://www.supamanage.buzz
2. Visit: https://studio.supamanage.buzz
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
- Verify `PRODUCTION_HOST` is correct: `182.191.91.226`
- Check that port 22 is accessible from the internet
- Verify SSH credentials are correct

### "Health checks failed"
- Services might need more time to start
- Check deployment logs for errors
- SSH into server and check container logs

## Security Note

**Secrets** (encrypted, never visible):
- SSH credentials
- Database password  
- JWT and encryption secrets

**Variables** (plain text, visible to everyone with repo access):
- URLs and domains
- Configuration values
- Port numbers

This separation follows security best practices!
