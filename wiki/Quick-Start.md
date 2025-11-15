# Quick Start Guide

Get SupaManager up and running in **5 minutes**!

---

## Prerequisites

Before starting, ensure you have:
- ✅ Docker (version 20.10+)
- ✅ Docker Compose (v2+)
- ✅ Git
- ✅ 4GB+ RAM available
- ✅ 20GB+ disk space

Check your Docker installation:
```bash
docker --version
docker compose version
```

---

## Installation Steps

### 1. Clone the Repository

```bash
git clone https://github.com/YOUR_USERNAME/supabase-manager.git
cd supabase-manager
```

### 2. Build Studio Image

The Studio UI needs to be built with custom patches:

```bash
cd studio
./build.sh v1.24.04 supa-manager/studio:v1.24.04 .env
cd ..
```

**⏱️ This takes 5-10 minutes** on the first run.

**What's happening:**
- Downloads Supabase Studio v1.24.04 source code
- Applies custom patches for SupaManager integration
- Builds a Docker image

### 3. Start Services

```bash
docker compose up -d
```

This starts three services:
- **PostgreSQL** (port 5432) - Management database
- **SupaManager API** (port 8080) - Backend API
- **Studio UI** (port 3000) - Web interface

### 4. Wait for Services

Wait 30-60 seconds for services to initialize:

```bash
# Check service status
docker compose ps

# Follow logs (optional)
docker compose logs -f
```

You should see all services as "Up" or "Up (healthy)":
```
NAME                            STATUS
supabase-manager-database-1     Up (healthy)
supabase-manager-supa-manager-1 Up
supabase-manager-studio-1       Up
```

### 5. Open Studio

Open your browser and navigate to:
```
http://localhost:3000
```

---

## First Login

**Create Your Account:**

1. Navigate to http://localhost:3000
2. Click "Sign Up"
3. Enter your email and password
4. Start managing your Supabase projects

> **Note:** The first account you create will have admin privileges.

---

## Verify Installation

After logging in, you should see:

1. **Dashboard** - Overview of your organizations
2. **Navigation** - Organizations, Projects, Settings
3. **No projects yet** - Ready to create your first project

---

## Create Your First Project

1. Click **"New Project"** button
2. Fill in project details:
   - **Name:** My First Project
   - **Database Password:** (choose a strong password)
   - **Region:** Select a region
   - **Pricing Plan:** (select appropriate plan)
3. Click **"Create Project"**

> ⚠️ **Important:** Currently, projects are created in the database but provisioning is not yet implemented. The project will show status "UNKNOWN". See [Current Status](Home#current-status) for details.

---

## What's Next?

Now that you have SupaManager running:

- **[First Steps](First-Steps)** - Learn the basics of using SupaManager
- **[Configuration Reference](Configuration-Reference)** - Customize your setup
- **[API Reference](API-Reference)** - Integrate with the API
- **[Troubleshooting](Troubleshooting)** - If you encounter issues

---

## Stopping Services

To stop all services:

```bash
docker compose down
```

To stop and remove all data:

```bash
docker compose down -v
```

To restart services:

```bash
docker compose up -d
```

---

## Common Issues

### Port Already in Use

**Error:** `Bind for 0.0.0.0:3000 failed: port is already allocated`

**Solution:** Stop the service using that port or change ports in `docker-compose.yml`

Check what's using the port:
```bash
sudo lsof -i :3000
sudo lsof -i :8080
sudo lsof -i :5432
```

### Studio Build Fails

**Error:** Patches fail to apply

**Solution:** Ensure you're using the correct version:
```bash
cd studio
./build.sh v1.24.04 supa-manager/studio:v1.24.04 .env
```

### Cannot Connect to API

**Check API logs:**
```bash
docker compose logs supa-manager
```

**Restart the API:**
```bash
docker compose restart supa-manager
```

### Database Connection Error

**Check database health:**
```bash
docker compose ps database
```

**View database logs:**
```bash
docker compose logs database
```

---

## Getting Help

If you're stuck:
1. Check the [Troubleshooting](Troubleshooting) guide
2. Review [FAQ](FAQ) for common questions
3. Search [GitHub Issues](https://github.com/YOUR_USERNAME/supabase-manager/issues)
4. Ask in [Discord](https://discord.gg/4k5HRe6YEp)

---

## Next Steps

- **[First Steps Tutorial](First-Steps)** - Learn to use SupaManager
- **[Configuration Guide](Configuration-Reference)** - Customize your setup
- **[Architecture Overview](Architecture-Overview)** - Understand how it works
- **[Development Guide](Development-Guide)** - Contribute to the project

---

**Congratulations! 🎉** You now have SupaManager running locally.
