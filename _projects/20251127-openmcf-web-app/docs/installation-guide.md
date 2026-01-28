# OpenMCF Web App - Installation Guide

**Last Updated:** December 11, 2025
**Status:** Internal Preview (Not yet public release)

---

## Overview

The OpenMCF Web App provides a unified web interface for managing cloud resources and deployments. Everything runs in a single Docker container including MongoDB, backend API, and frontend web UI.

## Architecture

```
┌─────────────────────────────────────────┐
│  Unified Docker Container              │
│                                         │
│  ┌───────────────────────────────────┐ │
│  │  MongoDB (port 27017)             │ │
│  │  - Data persistence               │ │
│  └───────────────────────────────────┘ │
│                                         │
│  ┌───────────────────────────────────┐ │
│  │  Backend API (port 50051)         │ │
│  │  - Connect-RPC server             │ │
│  │  - Pulumi deployments             │ │
│  └───────────────────────────────────┘ │
│                                         │
│  ┌───────────────────────────────────┐ │
│  │  Frontend (port 3000)             │ │
│  │  - Next.js web interface          │ │
│  └───────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

---

## Prerequisites

### Required

- **Docker Engine** - The container runtime
  - macOS: `brew install docker` or [Docker Desktop](https://docker.com/products/docker-desktop)
  - Linux: [Docker Engine Installation Guide](https://docs.docker.com/engine/install/)
  - Windows: [Docker Desktop for Windows](https://docs.docker.com/desktop/install/windows-install/)

### Verification

After installing Docker, verify it's working:

```bash
docker --version
docker info
```

---

## Installation

### Step 1: Install CLI

Install the OpenMCF CLI using Homebrew:

```bash
brew install plantonhq/tap/openmcf
```

Verify the installation:

```bash
openmcf version
```

### Step 2: Initialize Web App

Initialize the web app (this pulls the Docker image and sets up the container):

```bash
openmcf webapp init
```

This command will:
- ✅ Check Docker availability
- ✅ Pull the unified Docker image (~500MB)
- ✅ Create Docker volumes for data persistence
- ✅ Create the container with proper configuration
- ✅ Configure CLI to use the local backend

**Time required:** 2-5 minutes (depending on internet speed)

### Step 3: Start the Web App

Start the web app services:

```bash
openmcf webapp start
```

This command will:
- ✅ Start the container
- ✅ Wait for all services to be healthy
- ✅ Display access URLs

**Time required:** 30-60 seconds

### Step 4: Access the Web Interface

Once started, access the web app:

- **Frontend (Web UI):** http://localhost:3000
- **Backend API:** http://localhost:50051

---

## Daily Usage

### Starting the Web App

```bash
openmcf webapp start
```

### Stopping the Web App

```bash
openmcf webapp stop
```

Data is preserved when stopped. Start again anytime.

### Checking Status

```bash
openmcf webapp status
```

Shows container and service status.

### Viewing Logs

```bash
# View last 100 lines
openmcf webapp logs

# Follow logs in real-time
openmcf webapp logs -f

# Show more lines
openmcf webapp logs -n 500
```

Press `Ctrl+C` to stop following logs.

### Restarting

```bash
openmcf webapp restart
```

Useful after configuration changes or if services become unresponsive.

---

## Data Persistence

All data is stored in Docker volumes and persists across container restarts:

| Volume | Purpose | Location |
|--------|---------|----------|
| `openmcf-mongodb-data` | MongoDB database | `/data/db` |
| `openmcf-pulumi-state` | Pulumi state files | `/home/appuser/.pulumi` |
| `openmcf-go-cache` | Go build cache | `/home/appuser/go` |

### Backing Up Data

```bash
# Backup MongoDB
docker run --rm -v openmcf-mongodb-data:/data \
  -v $(pwd):/backup ubuntu tar czf /backup/mongodb-backup.tar.gz /data

# Backup Pulumi state
docker run --rm -v openmcf-pulumi-state:/data \
  -v $(pwd):/backup ubuntu tar czf /backup/pulumi-backup.tar.gz /data
```

---

## Uninstallation

### Keep Data (Recommended)

```bash
openmcf webapp uninstall
```

This removes the container but keeps data volumes. You can reinstall later with existing data.

### Complete Removal (Delete Everything)

```bash
openmcf webapp uninstall --purge-data
```

⚠️ **Warning:** This deletes all data including MongoDB database and Pulumi state. Cannot be undone!

---

## Troubleshooting

### Docker Not Found

**Error:**
```
❌ Error: Docker Engine is not installed or not running
```

**Solution:**
1. Install Docker (see Prerequisites section)
2. Verify: `docker info`

### Port Already in Use

**Error:**
```
Error response from daemon: driver failed programming external connectivity...
Bind for 0.0.0.0:3000 failed: port is already allocated
```

**Solution:**
1. Check what's using the port: `lsof -i :3000` or `lsof -i :50051`
2. Stop the conflicting service
3. Or modify ports in container creation (advanced)

### Services Not Starting

**Check logs:**
```bash
openmcf webapp logs -f
```

**Common issues:**
- MongoDB taking longer to initialize (wait 60-90 seconds)
- Insufficient disk space (MongoDB needs ~1GB)
- Docker resource limits (increase Docker memory/CPU allocation)

### Container Already Exists

**Error:**
```
⚠️ Container 'openmcf-webapp' already exists.
```

**Solution:**
```bash
# If you want to start existing container
openmcf webapp start

# If you want to start fresh
openmcf webapp uninstall
openmcf webapp init
```

---

## Configuration

### CLI Configuration

The CLI stores configuration in `~/.openmcf/config.yaml`:

```yaml
backend-url: http://localhost:50051
webapp-container-id: openmcf-webapp
webapp-version: latest
```

### Environment Variables

The container uses these environment variables (configured automatically):

| Variable | Value | Purpose |
|----------|-------|---------|
| `MONGODB_URI` | `mongodb://localhost:27017/openmcf` | MongoDB connection |
| `SERVER_PORT` | `50051` | Backend API port |
| `PORT` | `3000` | Frontend port |
| `PULUMI_HOME` | `/home/appuser/.pulumi` | Pulumi state location |

---

## Next Steps

Once the web app is running:

1. **Explore the Dashboard** - http://localhost:3000
2. **Create Cloud Resources** - Use the web interface to define infrastructure
3. **Deploy Resources** - Deploy to actual cloud providers
4. **Manage Credentials** - Store cloud provider credentials securely

---

## Getting Help

- View all commands: `openmcf webapp --help`
- Check service status: `openmcf webapp status`
- View logs: `openmcf webapp logs -f`

---

## System Requirements

| Resource | Minimum | Recommended |
|----------|---------|-------------|
| RAM | 2GB | 4GB+ |
| Disk Space | 2GB | 5GB+ |
| CPU | 2 cores | 4+ cores |
| Docker Version | 20.10+ | Latest |

---

## Security Notes

- The web app runs on localhost only (not exposed to network)
- MongoDB has no authentication (localhost only)
- Pulumi state is stored locally (not in cloud backend)
- For production use, additional security hardening is required

---

**Status:** This is an internal preview release. Not recommended for production use yet.


