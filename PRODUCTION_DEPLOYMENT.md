# Production Deployment Guide for SupaManager
**Domain:** supamanage.buzz

## Table of Contents
1. [Server Requirements](#server-requirements)
2. [DNS Configuration](#dns-configuration)
3. [Server Setup](#server-setup)
4. [SSL/TLS Certificates](#ssltls-certificates)
5. [Application Deployment](#application-deployment)
6. [Environment Configuration](#environment-configuration)
7. [Database Setup](#database-setup)
8. [Reverse Proxy (Nginx)](#reverse-proxy-nginx)
9. [Security Hardening](#security-hardening)
10. [Monitoring & Logging](#monitoring--logging)
11. [Backup Strategy](#backup-strategy)
12. [Maintenance](#maintenance)

---

## 1. Server Requirements

### Minimum Specifications:
- **CPU:** 4 cores (8+ recommended for multiple projects)
- **RAM:** 8GB minimum (16GB+ recommended)
- **Storage:** 100GB SSD minimum (500GB+ recommended)
- **OS:** Ubuntu 22.04 LTS (recommended)
- **Network:** Static IP address

### Required Software:
```bash
- Docker Engine 24.0+
- Docker Compose 2.20+
- Nginx 1.24+
- Certbot (for SSL)
- UFW (firewall)
```

---

## 2. DNS Configuration

### A. Main Domain Records (at your DNS provider):

```dns
# Main API and Studio
A       supamanage.buzz          → YOUR_SERVER_IP
A       www.supamanage.buzz      → YOUR_SERVER_IP
A       studio.supamanage.buzz   → YOUR_SERVER_IP

# Wildcard for project databases
A       *.supamanage.buzz        → YOUR_SERVER_IP

# Specific subdomains
A       db.*.supamanage.buzz     → YOUR_SERVER_IP  (if supported)
A       api.supamanage.buzz      → YOUR_SERVER_IP
```

### B. Example Project Subdomains:
```dns
# Each project will have:
db.project-abc123.supamanage.buzz    → YOUR_SERVER_IP
api.project-abc123.supamanage.buzz   → YOUR_SERVER_IP
```

### C. DNS Provider Configuration Steps:

**For Cloudflare:**
1. Add your domain to Cloudflare
2. Update nameservers at your registrar
3. Add A records as above
4. **Important:** Set Proxy status to "DNS only" (gray cloud) initially
5. Disable "Always Use HTTPS" initially (enable after SSL setup)

**For other providers (Namecheap, GoDaddy, etc.):**
1. Navigate to DNS Management
2. Add A records pointing to your server IP
3. Set TTL to 300 (5 minutes) for testing, increase to 3600 later

---

## 3. Server Setup

### A. Initial Server Configuration:

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install essential packages
sudo apt install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release \
    git \
    ufw \
    fail2ban \
    htop

# Set timezone
sudo timedatectl set-timezone UTC

# Set hostname
sudo hostnamectl set-hostname supamanage-prod
```

### B. Install Docker:

```bash
# Add Docker's official GPG key
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# Set up the repository
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker Engine
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# Start and enable Docker
sudo systemctl start docker
sudo systemctl enable docker

# Add your user to docker group
sudo usermod -aG docker $USER
newgrp docker

# Verify installation
docker --version
docker compose version
```

### C. Install Nginx:

```bash
sudo apt install -y nginx
sudo systemctl start nginx
sudo systemctl enable nginx
```

---

## 4. SSL/TLS Certificates

### A. Install Certbot:

```bash
sudo apt install -y certbot python3-certbot-nginx
```

### B. Obtain SSL Certificates:

**Option 1: Wildcard Certificate (Recommended)**
```bash
# Install DNS plugin for your provider (example: Cloudflare)
sudo apt install -y python3-certbot-dns-cloudflare

# Create Cloudflare credentials file
sudo mkdir -p /etc/letsencrypt
sudo nano /etc/letsencrypt/cloudflare.ini

# Add to cloudflare.ini:
dns_cloudflare_api_token = YOUR_CLOUDFLARE_API_TOKEN

# Secure the file
sudo chmod 600 /etc/letsencrypt/cloudflare.ini

# Obtain wildcard certificate
sudo certbot certonly \
  --dns-cloudflare \
  --dns-cloudflare-credentials /etc/letsencrypt/cloudflare.ini \
  -d supamanage.buzz \
  -d "*.supamanage.buzz"
```

**Option 2: Individual Certificates**
```bash
# For main domain
sudo certbot --nginx -d supamanage.buzz -d www.supamanage.buzz

# For studio
sudo certbot --nginx -d studio.supamanage.buzz
```

### C. Auto-renewal:

```bash
# Test renewal
sudo certbot renew --dry-run

# Certbot auto-renewal is set up automatically via systemd timer
sudo systemctl status certbot.timer
```

---

## 5. Application Deployment

### A. Clone Repository:

```bash
# Create application directory
sudo mkdir -p /opt/supamanage
sudo chown $USER:$USER /opt/supamanage
cd /opt/supamanage

# Clone repository
git clone https://github.com/haider-pw/supa-manager.git .

# Checkout production branch
git checkout update-v6
```

### B. Create Production docker-compose.yml:

```bash
nano /opt/supamanage/docker-compose.prod.yml
```

```yaml
version: "3.8"

networks:
  supamanage-network:
    driver: bridge

volumes:
  postgres-data:
  projects-data:

services:
  database:
    image: supabase/postgres:15.1.0.147
    container_name: supamanage-db
    restart: unless-stopped
    networks:
      - supamanage-network
    ports:
      - "127.0.0.1:5432:5432"  # Only accessible from localhost
    volumes:
      - postgres-data:/var/lib/postgresql/data
    environment:
      POSTGRES_PASSWORD: ${DATABASE_PASSWORD}
      POSTGRES_DB: supabase
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  supa-manager:
    build:
      context: ./supa-manager
      dockerfile: Dockerfile
    container_name: supamanage-api
    restart: unless-stopped
    networks:
      - supamanage-network
    ports:
      - "127.0.0.1:8080:8080"  # Only accessible from localhost (Nginx proxy)
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock  # For provisioning
      - projects-data:/root/projects
      - ./supa-manager/.env:/root/.env
    environment:
      - DATABASE_URL=postgres://postgres:${DATABASE_PASSWORD}@database:5432/supabase
    depends_on:
      database:
        condition: service_healthy

  studio:
    image: supa-manager/studio:v1.24.04
    container_name: supamanage-studio
    restart: unless-stopped
    networks:
      - supamanage-network
    ports:
      - "127.0.0.1:3000:3000"  # Only accessible from localhost (Nginx proxy)
    environment:
      - SUPABASE_URL=https://supamanage.buzz
      - SUPABASE_REST_URL=https://supamanage.buzz/rest/v1/
    depends_on:
      - supa-manager
```

---

## 6. Environment Configuration

### Create Production .env:

```bash
cd /opt/supamanage/supa-manager
cp .env.example .env.production
nano .env.production
```

```bash
# Database
DATABASE_URL=postgres://postgres:STRONG_PASSWORD_HERE@database:5432/supabase

# Authentication
ALLOW_SIGNUP=true
JWT_SECRET=YOUR_GENERATED_SECRET_HERE_64_CHARS
ENCRYPTION_SECRET=YOUR_GENERATED_SECRET_HERE_64_CHARS

# Service Version
SERVICE_VERSION_URL=https://supamanager.io/updates

# Postgres Defaults
POSTGRES_DISK_SIZE=10
POSTGRES_DEFAULT_VERSION=15.1
POSTGRES_DOCKER_IMAGE=supabase/postgres

# Domain Configuration
DOMAIN_STUDIO_URL=https://studio.supamanage.buzz
DOMAIN_BASE=supamanage.buzz

# DNS Hook (if you implement one)
DOMAIN_DNS_HOOK_URL=https://dns.supamanage.buzz
DOMAIN_DNS_HOOK_KEY=YOUR_DNS_HOOK_SECRET

# Provisioning
PROVISIONING_ENABLED=true
PROVISIONING_DOCKER_HOST=unix:///var/run/docker.sock
PROVISIONING_PROJECTS_DIR=/root/projects
PROVISIONING_BASE_POSTGRES_PORT=5433
PROVISIONING_BASE_KONG_HTTP_PORT=54321
```

### Generate Secrets:

```bash
# Generate JWT_SECRET (64 characters)
openssl rand -base64 48

# Generate ENCRYPTION_SECRET (64 characters)
openssl rand -base64 48

# Generate strong database password
openssl rand -base64 32
```

---

## 7. Database Setup

### A. Create Database Password File:

```bash
# Create .env file in root
cd /opt/supamanage
nano .env
```

```bash
DATABASE_PASSWORD=YOUR_STRONG_DB_PASSWORD_HERE
```

### B. Secure the file:

```bash
chmod 600 /opt/supamanage/.env
chmod 600 /opt/supamanage/supa-manager/.env.production
```

---

## 8. Reverse Proxy (Nginx)

### A. Create Nginx Configuration:

```bash
sudo nano /etc/nginx/sites-available/supamanage.buzz
```

```nginx
# Rate limiting
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;
limit_req_zone $binary_remote_addr zone=studio_limit:10m rate=30r/s;

# Main API
server {
    listen 80;
    listen [::]:80;
    server_name supamanage.buzz www.supamanage.buzz api.supamanage.buzz;

    # Redirect HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name supamanage.buzz www.supamanage.buzz api.supamanage.buzz;

    # SSL Configuration
    ssl_certificate /etc/letsencrypt/live/supamanage.buzz/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/supamanage.buzz/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # Security Headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # Rate limiting
    limit_req zone=api_limit burst=20 nodelay;

    # Client body size (for file uploads)
    client_max_body_size 50M;

    # Proxy to API
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;

        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
}

# Studio Frontend
server {
    listen 80;
    listen [::]:80;
    server_name studio.supamanage.buzz;

    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name studio.supamanage.buzz;

    ssl_certificate /etc/letsencrypt/live/supamanage.buzz/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/supamanage.buzz/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # Security Headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;

    # Rate limiting
    limit_req zone=studio_limit burst=50 nodelay;

    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}

# Project Database Wildcards (db.project-*.supamanage.buzz)
server {
    listen 80;
    listen [::]:80;
    server_name ~^db\.(?<project>[^.]+)\.supamanage\.buzz$;

    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name ~^db\.(?<project>[^.]+)\.supamanage\.buzz$;

    ssl_certificate /etc/letsencrypt/live/supamanage.buzz/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/supamanage.buzz/privkey.pem;

    # PostgreSQL connections handled at TCP level (not HTTP)
    # This requires stream module configuration (see below)
    return 444;  # Close connection
}
```

### B. PostgreSQL Stream Configuration:

```bash
sudo nano /etc/nginx/modules-available/stream.conf
```

```nginx
stream {
    # Log format for TCP connections
    log_format tcp_log '$remote_addr [$time_local] $protocol $status '
                       '$bytes_sent $bytes_received $session_time';

    # Project database routing (TCP proxy)
    # Note: This requires dynamic routing based on SNI
    # For simplicity, route all PostgreSQL traffic through single port

    upstream postgres_backend {
        # Dynamically route based on project
        # This is simplified - you'll need custom logic for multi-project routing
        server 127.0.0.1:5432;
    }

    server {
        listen 5432;
        proxy_pass postgres_backend;
        proxy_timeout 3600s;
        proxy_connect_timeout 3s;
    }
}
```

### C. Enable Site:

```bash
# Test configuration
sudo nginx -t

# Enable site
sudo ln -s /etc/nginx/sites-available/supamanage.buzz /etc/nginx/sites-enabled/

# Reload Nginx
sudo systemctl reload nginx
```

---

## 9. Security Hardening

### A. Firewall (UFW):

```bash
# Reset firewall
sudo ufw --force reset

# Default policies
sudo ufw default deny incoming
sudo ufw default allow outgoing

# Allow SSH
sudo ufw allow 22/tcp

# Allow HTTP/HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Allow PostgreSQL (only if needed externally)
# sudo ufw allow 5432/tcp

# Enable firewall
sudo ufw enable

# Check status
sudo ufw status verbose
```

### B. Fail2Ban:

```bash
# Configure Fail2Ban for SSH
sudo nano /etc/fail2ban/jail.local
```

```ini
[DEFAULT]
bantime = 3600
findtime = 600
maxretry = 5

[sshd]
enabled = true
port = 22
logpath = /var/log/auth.log

[nginx-http-auth]
enabled = true
port = http,https
logpath = /var/log/nginx/error.log
```

```bash
# Start Fail2Ban
sudo systemctl start fail2ban
sudo systemctl enable fail2ban
```

### C. Docker Security:

```bash
# Restrict Docker socket permissions
sudo chmod 660 /var/run/docker.sock

# Enable Docker content trust
echo "export DOCKER_CONTENT_TRUST=1" >> ~/.bashrc
source ~/.bashrc
```

---

## 10. Monitoring & Logging

### A. Docker Logging:

```bash
# Configure Docker daemon logging
sudo nano /etc/docker/daemon.json
```

```json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
```

```bash
sudo systemctl restart docker
```

### B. System Monitoring:

```bash
# Install monitoring tools
sudo apt install -y netdata

# Access at http://YOUR_SERVER_IP:19999
# Secure with Nginx reverse proxy
```

### C. Application Logs:

```bash
# View logs
docker compose -f docker-compose.prod.yml logs -f

# View specific service
docker compose -f docker-compose.prod.yml logs -f supa-manager

# Save logs to file
docker compose -f docker-compose.prod.yml logs > /var/log/supamanage/app.log
```

---

## 11. Backup Strategy

### A. Database Backups:

```bash
# Create backup script
sudo nano /opt/supamanage/backup.sh
```

```bash
#!/bin/bash
BACKUP_DIR="/opt/supamanage/backups"
DATE=$(date +%Y%m%d_%H%M%S)
DB_PASSWORD=$(grep DATABASE_PASSWORD /opt/supamanage/.env | cut -d '=' -f2)

# Create backup directory
mkdir -p $BACKUP_DIR

# Backup main database
docker exec supamanage-db pg_dump -U postgres supabase | gzip > $BACKUP_DIR/supabase_$DATE.sql.gz

# Keep only last 7 days of backups
find $BACKUP_DIR -name "supabase_*.sql.gz" -mtime +7 -delete

echo "Backup completed: supabase_$DATE.sql.gz"
```

```bash
# Make executable
chmod +x /opt/supamanage/backup.sh

# Add to crontab (daily at 2 AM)
crontab -e
```

Add:
```
0 2 * * * /opt/supamanage/backup.sh >> /var/log/supamanage/backup.log 2>&1
```

### B. Upload Backups to Cloud:

```bash
# Install rclone for cloud backup
curl https://rclone.org/install.sh | sudo bash

# Configure with your cloud provider
rclone config

# Add to backup script
rclone copy $BACKUP_DIR remote:supamanage-backups/
```

---

## 12. Maintenance

### A. Start Application:

```bash
cd /opt/supamanage
docker compose -f docker-compose.prod.yml up -d
```

### B. Stop Application:

```bash
docker compose -f docker-compose.prod.yml down
```

### C. Update Application:

```bash
cd /opt/supamanage
git pull origin update-v6
docker compose -f docker-compose.prod.yml down
docker compose -f docker-compose.prod.yml build
docker compose -f docker-compose.prod.yml up -d
```

### D. View Status:

```bash
docker compose -f docker-compose.prod.yml ps
docker stats
```

### E. Health Checks:

```bash
# Check API
curl https://supamanage.buzz/health

# Check Studio
curl https://studio.supamanage.buzz

# Check database
docker exec supamanage-db pg_isready
```

---

## 13. Post-Deployment Checklist

- [ ] DNS records configured and propagated
- [ ] SSL certificates installed and auto-renewal working
- [ ] Firewall configured and enabled
- [ ] Docker containers running and healthy
- [ ] Nginx proxy working
- [ ] Can access Studio at https://studio.supamanage.buzz
- [ ] Can access API at https://supamanage.buzz
- [ ] Database backups configured and tested
- [ ] Monitoring setup complete
- [ ] Create first admin account
- [ ] Test project creation
- [ ] Test database provisioning
- [ ] Test encrypted connection strings
- [ ] Documentation updated

---

## 14. Troubleshooting

### Check Logs:
```bash
# Nginx logs
sudo tail -f /var/log/nginx/error.log

# Docker logs
docker compose -f docker-compose.prod.yml logs -f

# System logs
sudo journalctl -f
```

### Common Issues:

**SSL Certificate Issues:**
```bash
sudo certbot certificates
sudo certbot renew --force-renewal
```

**Port Conflicts:**
```bash
sudo lsof -i :80
sudo lsof -i :443
sudo lsof -i :5432
```

**Docker Issues:**
```bash
docker system prune -a
docker volume prune
```

---

## 15. Support & Updates

- GitHub: https://github.com/haider-pw/supa-manager
- Issues: https://github.com/haider-pw/supa-manager/issues
- Documentation: Update as you customize

---

**Created:** 2025-11-16
**Domain:** supamanage.buzz
**Version:** 1.0
**Status:** Production Ready
