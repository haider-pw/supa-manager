# Security Guidelines

## Environment Variables and Secrets

### Important Security Notice

⚠️ **Never commit sensitive credentials to version control!**

This project uses environment variables for sensitive configuration. These values should **never** be hardcoded in `docker-compose.yml` or any other file that is committed to git.

---

## Quick Start Security Setup

### 1. Copy the Example File

```bash
cp .env.example .env
```

### 2. Generate Secure Secrets

For production, generate strong random secrets:

```bash
# Generate JWT_SECRET
openssl rand -base64 48

# Generate ENCRYPTION_SECRET
openssl rand -base64 48

# Generate DOMAIN_DNS_HOOK_KEY
openssl rand -base64 32
```

### 3. Update .env File

Edit `.env` and replace the placeholder values:

```bash
# Database Configuration
POSTGRES_PASSWORD=your-secure-database-password

# Security Secrets (CHANGE THESE!)
JWT_SECRET=<paste-generated-secret-here>
ENCRYPTION_SECRET=<paste-different-generated-secret-here>

# DNS Webhook Key
DOMAIN_DNS_HOOK_KEY=<paste-generated-hook-key-here>
```

---

## Environment Variables Reference

### Required Secrets

| Variable | Purpose | Minimum Length | Production Required |
|----------|---------|----------------|---------------------|
| `JWT_SECRET` | Signs JWT authentication tokens | 32 characters | ✅ Yes |
| `ENCRYPTION_SECRET` | Encrypts sensitive data in database | 32 characters | ✅ Yes |
| `POSTGRES_PASSWORD` | Database password | 12 characters | ✅ Yes |
| `DOMAIN_DNS_HOOK_KEY` | DNS webhook authentication | 16 characters | ⚠️ If using DNS hooks |

### Default Values

The `docker-compose.yml` file uses this syntax:

```yaml
JWT_SECRET: ${JWT_SECRET:-change-me-in-production-min-32-chars}
```

**Explanation:**
- `${JWT_SECRET}` - Reads from .env file
- `:-change-me...` - Uses this default if not set (for local development only)

**Production:** The defaults are **intentionally weak** to force you to set proper values!

---

## Security Best Practices

### ✅ DO

1. **Generate unique secrets** for each environment (dev, staging, production)
2. **Use strong passwords** - minimum 32 characters for secrets
3. **Rotate secrets regularly** - every 90 days for production
4. **Keep .env file secure** - limit file permissions: `chmod 600 .env`
5. **Use different secrets** for JWT_SECRET and ENCRYPTION_SECRET
6. **Back up secrets securely** - use encrypted password manager

### ❌ DON'T

1. **Never commit .env files** - they are in .gitignore for a reason
2. **Never use default secrets in production** - "secret", "password", etc.
3. **Never share secrets in chat/email** - use secure sharing tools
4. **Never reuse secrets** across environments
5. **Never hardcode secrets** in source code
6. **Never log secrets** to console or files

---

## Local Development

For local development, you can use weak secrets:

```bash
# .env (local development)
POSTGRES_PASSWORD=postgres
JWT_SECRET=local-dev-secret-not-for-production
ENCRYPTION_SECRET=local-dev-encryption-not-for-production
DOMAIN_DNS_HOOK_KEY=local-dev-hook-key
```

**Note:** Never use these values in any environment accessible from the internet!

---

## Production Deployment

### Required Changes for Production

1. **Strong Secrets**
   ```bash
   # Generate with:
   openssl rand -base64 48
   ```

2. **Secure Database Password**
   ```bash
   # Generate with:
   openssl rand -base64 32
   ```

3. **File Permissions**
   ```bash
   chmod 600 .env
   chown root:root .env  # Or your service user
   ```

4. **Environment Isolation**
   - Different secrets for dev, staging, production
   - Never copy production .env to development machines

5. **Secret Management**
   - Consider using Docker secrets (Swarm mode)
   - Or Kubernetes secrets
   - Or external secret manager (Vault, AWS Secrets Manager, etc.)

---

## Docker Secrets (Alternative for Production)

For production, consider using Docker secrets instead of .env files:

```yaml
# docker-compose.yml
services:
  supa-manager:
    secrets:
      - jwt_secret
      - encryption_secret
    environment:
      JWT_SECRET_FILE: /run/secrets/jwt_secret
      ENCRYPTION_SECRET_FILE: /run/secrets/encryption_secret

secrets:
  jwt_secret:
    external: true
  encryption_secret:
    external: true
```

Create secrets:
```bash
echo "your-jwt-secret" | docker secret create jwt_secret -
echo "your-encryption-secret" | docker secret create encryption_secret -
```

---

## Rotating Secrets

### JWT_SECRET Rotation

1. Generate new secret
2. Update .env file
3. Restart supa-manager service
4. All users will need to re-login (tokens invalidated)

```bash
# Generate new secret
NEW_JWT_SECRET=$(openssl rand -base64 48)

# Update .env
sed -i "s/JWT_SECRET=.*/JWT_SECRET=$NEW_JWT_SECRET/" .env

# Restart
docker compose restart supa-manager
```

### ENCRYPTION_SECRET Rotation

⚠️ **More complex** - requires data re-encryption!

1. Create migration script to:
   - Decrypt data with old secret
   - Re-encrypt with new secret
2. Test migration on backup
3. Apply to production
4. Update .env
5. Restart services

**Not recommended without proper planning!**

### POSTGRES_PASSWORD Rotation

```bash
# 1. Update password in database
docker exec -it supabase-manager-database-1 psql -U postgres -c "ALTER USER postgres PASSWORD 'new-password';"

# 2. Update .env
echo "POSTGRES_PASSWORD=new-password" >> .env

# 3. Restart all services
docker compose restart
```

---

## Checking for Exposed Secrets

### Check Git History

```bash
# Search for potential secrets in git history
git log -p | grep -E "(secret|password|key)" -i

# Check if .env was ever committed
git log --all --full-history -- "**/.env"
```

### If Secrets Were Committed

1. **Rotate all exposed secrets immediately**
2. **Use git-filter-repo or BFG Repo-Cleaner** to remove from history
3. **Force push** (⚠️ breaks others' clones)
4. **Notify team members** to re-clone

```bash
# Remove from history (example with BFG)
bfg --delete-files .env
git reflog expire --expire=now --all
git gc --prune=now --aggressive
git push origin --force --all
```

---

## Security Checklist

### Before Deploying to Production

- [ ] Generated strong secrets using `openssl rand -base64 48`
- [ ] Updated all secrets in .env file
- [ ] Verified .env is in .gitignore
- [ ] Set secure file permissions: `chmod 600 .env`
- [ ] Different secrets than development
- [ ] Backed up secrets in secure password manager
- [ ] Documented secret rotation procedure
- [ ] Set up monitoring for unauthorized access
- [ ] Enabled HTTPS/TLS
- [ ] Configured firewall rules
- [ ] Disabled ALLOW_SIGNUP (if appropriate)

---

## Reporting Security Issues

If you discover a security vulnerability:

1. **Do NOT open a public issue**
2. **Email:** security@example.com (or repository owner)
3. **Include:**
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

We will respond within 48 hours and work with you to address the issue.

---

## Additional Resources

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Docker Security Best Practices](https://docs.docker.com/engine/security/)
- [Secret Management Guide](https://www.vaultproject.io/docs)
- [Password Generation](https://bitwarden.com/password-generator/)

---

**Remember:** Security is everyone's responsibility. When in doubt, ask!
