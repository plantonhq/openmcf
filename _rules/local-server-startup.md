---
description: Action rule to start backend, frontend, and optionally CLI for local development
alwaysApply: false
---

# Rule: Start Development Servers

## Purpose

Start the OpenMCF development environment with backend, frontend, and optionally CLI. This rule ensures all services are properly started in the correct order with their dependencies.

## Rule Type

**Action Rule** - Interactive rule that asks for user input before executing

## Service Architecture

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                          Local Development Setup                              │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                               │
│   ┌─────────────┐     ┌─────────────────┐     ┌─────────────────────┐       │
│   │  MongoDB    │────▶│    Backend      │◀────│     Frontend        │       │
│   │  :27017     │     │    :50051       │     │     :3000           │       │
│   └─────────────┘     └────────┬────────┘     └─────────────────────┘       │
│                                │                                              │
│                                ▼                                              │
│                       ┌─────────────────┐                                     │
│                       │   CLI (opt)     │                                     │
│                       │ openmcf │                                     │
│                       └─────────────────┘                                     │
│                                                                               │
└──────────────────────────────────────────────────────────────────────────────┘
```

## User Input Required

**Before starting, ask the user:**

> Do you want to start the CLI locally as well? (yes/no)
>
> - **yes**: Will configure CLI to connect to local backend (http://localhost:50051) and build it
> - **no**: Will only start backend and frontend servers

## Prerequisites

Before executing this rule, ensure the following are installed:

- **MongoDB**: `brew install mongodb-community` (macOS) or Docker
- **Go**: 1.24.7+
- **Node.js**: 20+
- **Yarn**: 3.x

## Execution Steps

### Step 1: Verify MongoDB is Running

Check if MongoDB is running, and start it if not:

```bash
# Check if MongoDB is already running
brew services list | grep mongodb-community

# If not running, start MongoDB
brew services start mongodb-community
```

Alternative using Docker:

```bash
docker run -d --name mongodb -p 27017:27017 mongo:latest
```

### Step 2: Start Backend Server (Always)

The backend runs on port `50051` and connects to MongoDB.

```bash
# Navigate to backend directory
cd app/backend

# Start in development mode (with hot reload)
go run ./cmd/server
```

**Environment Variables (defaults are usually fine):**

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `50051` | Backend API port |
| `MONGODB_URI` | `mongodb://localhost:27017` | MongoDB connection |
| `MONGODB_DATABASE` | `openmcf` | Database name |

### Step 3: Start Frontend Server (Always)

The frontend runs on port `3000` and connects to the backend.

```bash
# Navigate to frontend directory
cd app/frontend

# Install dependencies (if needed)
yarn install

# Start development server
yarn dev
```

**Environment Variables (defaults are usually fine):**

| Variable | Default | Description |
|----------|---------|-------------|
| `NEXT_PUBLIC_API_URL` | `http://localhost:50051` | Backend API URL |
| `PORT` | `3000` | Frontend port |

### Step 4: Configure and Build CLI (Conditional)

**Only if user selected "yes" to starting CLI locally:**

```bash
# Build CLI from project root
cd /Users/satishlakhani/scm/github.com/plantonhq/openmcf
make build_darwin

# Install to ~/bin
cp ./build/openmcf-darwin ~/bin/openmcf
chmod +x ~/bin/openmcf

# Configure CLI to use local backend
openmcf config set backend-url http://localhost:50051

# Verify configuration
openmcf config list
```

## Terminal Layout

When starting servers, use separate terminals for each service:

| Terminal | Service | Command | Working Directory |
|----------|---------|---------|-------------------|
| Terminal 1 | MongoDB | `brew services start mongodb-community` | Any |
| Terminal 2 | Backend | `go run ./cmd/server` | `app/backend` |
| Terminal 3 | Frontend | `yarn dev` | `app/frontend` |
| Terminal 4 | CLI (optional) | CLI commands | Any |

## Execution Commands Summary

```bash
# Terminal 1: Ensure MongoDB is running
brew services start mongodb-community

# Terminal 2: Start Backend
cd app/backend && go run ./cmd/server

# Terminal 3: Start Frontend
cd app/frontend && yarn install && yarn dev

# Terminal 4 (if CLI requested): Build and Configure CLI
cd /Users/satishlakhani/scm/github.com/plantonhq/openmcf
make build_darwin
cp ./build/openmcf-darwin ~/bin/openmcf
chmod +x ~/bin/openmcf
openmcf config set backend-url http://localhost:50051
```

## Verification

After starting, verify services are running:

1. **MongoDB**: `mongosh` should connect successfully
2. **Backend**: `curl http://localhost:50051` should respond (or use gRPC tools)
3. **Frontend**: Open `http://localhost:3000` in browser
4. **CLI**: `openmcf config list` shows `backend-url=http://localhost:50051`

## Stopping Services

```bash
# Stop Frontend: Ctrl+C in Terminal 3

# Stop Backend: Ctrl+C in Terminal 2

# Stop MongoDB
brew services stop mongodb-community
# Or if using Docker:
docker stop mongodb
```

## Troubleshooting

### Backend Won't Start

- Check MongoDB is running: `brew services list | grep mongodb`
- Check port 50051 is free: `lsof -i :50051`
- Check logs for connection errors

### Frontend Build Fails

- Clear node_modules: `rm -rf node_modules && yarn install`
- Clear Next.js cache: `rm -rf .next`
- Regenerate proto code: `cd ../backend && make generate`

### CLI Can't Connect to Backend

- Verify backend is running: `curl http://localhost:50051`
- Check CLI config: `openmcf config list`
- Reconfigure: `openmcf config set backend-url http://localhost:50051`

## Notes

- Backend must be started before frontend for full functionality
- The frontend can start without backend but API calls will fail
- CLI requires backend to be running for all API operations
- Use `docker-compose up` for a simpler all-in-one setup (production mode)
