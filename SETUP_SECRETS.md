# Setting Up Environment Variables

Quick guide to set up your environment variables securely.

---

## For Local Development

### 1. Copy the Example File

```bash
cp .env.example .env
```

### 2. Use Default Values (Local Only)

The defaults in `.env` are suitable for local development:

```bash
POSTGRES_PASSWORD=postgres
JWT_SECRET=change-me-in-production-min-32-chars
ENCRYPTION_SECRET=change-me-in-production-min-32-chars
DOMAIN_DNS_HOOK_KEY=change-me-in-production
```

⚠️ **These defaults are intentionally weak and should NEVER be used in production!**

### 3. Start Services

```bash
docker compose up -d
```

---

## For Production Deployment

### 1. Copy the Example File

```bash
cp .env.example .env
```

### 2. Generate Strong Secrets

```bash
# Generate JWT_SECRET (48 characters)
echo "JWT_SECRET=$(openssl rand -base64 48)" >> .env.production

# Generate ENCRYPTION_SECRET (48 characters, must be different!)
echo "ENCRYPTION_SECRET=$(openssl rand -base64 48)" >> .env.production

# Generate DOMAIN_DNS_HOOK_KEY (32 characters)
echo "DOMAIN_DNS_HOOK_KEY=$(openssl rand -base64 32)" >> .env.production

# Generate strong database password
echo "POSTGRES_PASSWORD=$(openssl rand -base64 32)" >> .env.production
```

### 3. Edit .env File

```bash
nano .env  # or vim, code, etc.
```

Replace the placeholder values:

```bash
# Database Configuration
POSTGRES_PASSWORD=<paste-generated-password>

# Security Secrets
JWT_SECRET=<paste-generated-jwt-secret>
ENCRYPTION_SECRET=<paste-generated-encryption-secret>

# DNS Webhook Key
DOMAIN_DNS_HOOK_KEY=<paste-generated-hook-key>
```

### 4. Secure the File

```bash
# Set restrictive permissions
chmod 600 .env

# Verify it's ignored by git
git status  # Should not show .env
```

### 5. Start Services

```bash
docker compose up -d
```

---

## Verifying Configuration

### Check Environment Variables Are Loaded

```bash
# View environment variables (excluding secrets)
docker compose config | grep -v -E "(SECRET|PASSWORD|KEY)"

# Check specific service
docker exec supabase-manager-supa-manager-1 env | grep JWT_SECRET
# Should show: JWT_SECRET=your-actual-secret-here
```

### Test Application

```bash
# Health check
curl http://localhost:8080/health

# Should return: {"status":"ok","timestamp":"..."}
```

---

## Troubleshooting

### Environment Variables Not Loading

**Issue:** Services use default values instead of .env values

**Solution:**
1. Ensure `.env` is in the project root (same directory as `docker-compose.yml`)
2. Restart services: `docker compose down && docker compose up -d`
3. Check file name is exactly `.env` (not `.env.txt` or `.env.local`)

### Permission Denied Errors

**Issue:** Cannot read .env file

**Solution:**
```bash
# Add read permission for your user
chmod 600 .env
chown $USER:$USER .env
```

### Secrets Not Strong Enough

**Issue:** JWT_SECRET too short

**Solution:**
```bash
# Generate longer secret
openssl rand -base64 64
```

Minimum requirements:
- JWT_SECRET: 32 characters
- ENCRYPTION_SECRET: 32 characters
- POSTGRES_PASSWORD: 12 characters

---

## Different Environments

### Development (.env)

```bash
POSTGRES_PASSWORD=postgres
JWT_SECRET=dev-secret-not-for-production
ENCRYPTION_SECRET=dev-encryption-not-for-production
DOMAIN_DNS_HOOK_KEY=dev-hook-key
```

### Staging (.env.staging)

```bash
POSTGRES_PASSWORD=$(openssl rand -base64 32)
JWT_SECRET=$(openssl rand -base64 48)
ENCRYPTION_SECRET=$(openssl rand -base64 48)
DOMAIN_DNS_HOOK_KEY=$(openssl rand -base64 32)
```

### Production (.env.production)

```bash
# Use strongest secrets
# Store backup in secure password manager
# Never share or commit to git
```

### Using Different Files

```bash
# Development (default)
docker compose up -d

# Staging
docker compose --env-file .env.staging up -d

# Production
docker compose --env-file .env.production up -d
```

---

## Security Reminders

✅ **DO:**
- Use `.env` for local development
- Generate strong secrets for production
- Keep .env files in .gitignore
- Back up secrets in password manager
- Use different secrets per environment

❌ **DON'T:**
- Commit .env files to git
- Use default/weak secrets in production
- Share .env files via email/chat
- Reuse secrets across environments
- Log secrets to console

---

## Need Help?

- [Full Security Guide](SECURITY.md)
- [Configuration Reference](wiki/Configuration-Reference.md)
- [FAQ](wiki/FAQ.md)
- [Discord](https://discord.gg/4k5HRe6YEp)

---

**Next Step:** [Start the application](README.md#quick-start)
