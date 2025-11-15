# SupaManager Wiki

Welcome to the **SupaManager** documentation wiki! This comprehensive guide will help you understand, install, configure, and contribute to SupaManager.

## 📚 What is SupaManager?

SupaManager is a self-hosted management platform for Supabase instances. It provides a web-based interface (using the official Supabase Studio) to create and manage multiple Supabase projects from a single control panel.

**Project by:** [Harry Bairstow](https://twitter.com/TheHarryET)

---

## 🎯 Key Features

- **Multi-tenant Architecture** - Manage multiple Supabase projects from one interface
- **Organization Management** - Group projects under organizations with team access
- **Supabase Studio Integration** - Official Supabase UI with custom patches
- **RESTful API** - Complete management API for automation
- **Self-hosted** - Run on your own infrastructure
- **Docker-based** - Easy deployment with Docker Compose

---

## 📖 Documentation Structure

### Getting Started
- **[Installation Guide](Installation-Guide)** - Step-by-step setup instructions
- **[Quick Start](Quick-Start)** - Get running in 5 minutes
- **[First Steps](First-Steps)** - Create your first organization and project

### Configuration & Usage
- **[Configuration Reference](Configuration-Reference)** - All configuration options explained
- **[Environment Variables](Environment-Variables)** - Complete environment variable reference
- **[Database Schema](Database-Schema)** - Understanding the management database
- **[API Reference](API-Reference)** - Complete API endpoint documentation

### Architecture & Design
- **[Architecture Overview](Architecture-Overview)** - System design and components
- **[Service Communication](Service-Communication)** - How services interact
- **[Docker Networking](Docker-Networking)** - Understanding container networking
- **[Supabase Stack](Supabase-Stack)** - Complete Supabase service architecture

### Development
- **[Development Guide](Development-Guide)** - Setting up your dev environment
- **[Contributing](Contributing)** - How to contribute to the project
- **[Code Structure](Code-Structure)** - Understanding the codebase
- **[Running Tests](Running-Tests)** - Testing guide

### Operations
- **[Deployment](Deployment)** - Production deployment guide
- **[Monitoring](Monitoring)** - Health checks and monitoring
- **[Backup & Recovery](Backup-Recovery)** - Data backup strategies
- **[Troubleshooting](Troubleshooting)** - Common issues and solutions

### Reference
- **[CLI Commands](CLI-Commands)** - Docker and management commands
- **[FAQ](FAQ)** - Frequently asked questions
- **[Roadmap](Roadmap)** - Future plans and features
- **[Changelog](Changelog)** - Version history

---

## 🚀 Quick Links

### For Users
- [Installation Guide](Installation-Guide) - Get started with installation
- [Quick Start](Quick-Start) - 5-minute setup
- [API Reference](API-Reference) - API documentation

### For Developers
- [Development Guide](Development-Guide) - Setup dev environment
- [Architecture Overview](Architecture-Overview) - Understand the system
- [Contributing](Contributing) - Contribute to the project

### For Operators
- [Deployment](Deployment) - Production deployment
- [Monitoring](Monitoring) - Monitor your installation
- [Troubleshooting](Troubleshooting) - Solve common issues

---

## ⚠️ Current Status

> **Active Development Notice**
>
> SupaManager is currently in active development. The management API and Studio UI are functional, but **dynamic project provisioning is not yet implemented**. Projects can be created in the database but will not automatically spin up Supabase infrastructure.

### What Works Today
- ✅ User authentication and authorization
- ✅ Organization and team management
- ✅ Project metadata management
- ✅ Supabase Studio integration
- ✅ Complete REST API

### What's Coming Soon
- 🚧 Automatic Supabase project provisioning
- 🚧 Project lifecycle management (pause/resume/delete)
- 🚧 Real-time status monitoring
- 🚧 Resource usage tracking
- 🚧 Kubernetes deployment support

See the [Roadmap](Roadmap) for detailed future plans.

---

## 🆘 Need Help?

- **Issues & Bugs:** [GitHub Issues](https://github.com/YOUR_USERNAME/supabase-manager/issues)
- **Discord:** [Harry's Discord Server](https://discord.gg/4k5HRe6YEp)
- **Twitter:** [@TheHarryET](https://twitter.com/TheHarryET)
- **Troubleshooting:** Check the [Troubleshooting](Troubleshooting) page
- **FAQ:** See [Frequently Asked Questions](FAQ)

---

## 📄 License

SupaManager is licensed under the **GNU General Public License v3.0**.

This means you can:
- ✅ Use it commercially
- ✅ Modify the source code
- ✅ Distribute copies
- ✅ Use it privately

But you must:
- 📋 Disclose the source
- 📋 Include the license
- 📋 State changes made
- 📋 Keep the same license

See [LICENSE](https://github.com/YOUR_USERNAME/supabase-manager/blob/main/LICENSE) for full details.

---

## 🙏 Acknowledgments

- **Supabase Team** - For the amazing open-source Firebase alternative
- **Contributors** - Everyone who has contributed to this project
- **Community** - For feedback, bug reports, and feature requests

---

## 📚 External Resources

- [Supabase Official Docs](https://supabase.com/docs)
- [Supabase Self-hosting Guide](https://supabase.com/docs/guides/self-hosting)
- [Supabase GitHub](https://github.com/supabase/supabase)
- [Docker Documentation](https://docs.docker.com/)
- [Go Documentation](https://go.dev/doc/)

---

**Last Updated:** November 2025
