# OpenMCF Backend

The backend service for OpenMCF, providing gRPC/Connect-RPC APIs for cloud resource management using Pulumi.

## Architecture

The backend is a Go-based gRPC server that:
- Manages cloud resources via MongoDB
- Executes infrastructure deployments using Pulumi
- Streams deployment logs in real-time
- Resolves cloud provider credentials from the database

## Key Components

### Services

- **CloudResourceService** - CRUD operations for cloud resources
- **StackUpdateService** - Pulumi deployment orchestration and streaming
- **CredentialService** - Cloud provider credential management

### Database Repositories

- **CloudResourceRepository** - Cloud resource persistence
- **StackUpdateRepository** - Deployment job tracking
- **StackUpdateStreamingResponseRepository** - Real-time log streaming
- **CredentialRepository** - Credential storage and retrieval

## Pulumi Backend Configuration

The backend uses Pulumi for infrastructure deployments. By default, it uses a **local file-based backend** that stores state in the container's filesystem.

### How It Works

1. **Entrypoint Initialization** - `app/entrypoint-unified.sh` runs on container startup
2. **Backend Login** - Executes `pulumi login --local --non-interactive` (or cloud if `PULUMI_ACCESS_TOKEN` is set)
3. **State Storage** - Pulumi state is stored in `/home/appuser/.pulumi/state`
4. **Persistence** - State is persisted via Docker volume `pulumi-state`

### Why Local Backend?

- **No external dependencies** - Works out of the box without Pulumi Cloud account
- **Data privacy** - All state data stays on your machine
- **Cost** - Free, no API limits
- **Simplicity** - No authentication tokens required

### Using Pulumi Cloud (Optional)

If you prefer using Pulumi Cloud for state management:

1. Get your access token from https://app.pulumi.com/account/tokens
2. Set environment variables in `docker-compose.yml`:

```yaml
environment:
  - PULUMI_ACCESS_TOKEN=pul-xxxxxxxxxxxxx
  - PULUMI_BACKEND_URL=https://api.pulumi.com
```

The entrypoint script automatically detects `PULUMI_ACCESS_TOKEN` and uses Pulumi Cloud instead of the local backend. No code changes needed.

## Environment Variables

### Required

- `MONGODB_URI` - MongoDB connection string (default: `mongodb://localhost:27017/openmcf`)
- `SERVER_PORT` - gRPC server port (default: `50051`)

### Pulumi Configuration

- `PULUMI_HOME` - Pulumi home directory (default: `/home/appuser/.pulumi`)
- `PULUMI_STATE_DIR` - State storage directory (default: `/home/appuser/.pulumi/state`)
- `PULUMI_CONFIG_PASSPHRASE` - Encryption passphrase for secrets (default: `openmcf-default-passphrase`)
- `PULUMI_SKIP_UPDATE_CHECK` - Disable update checks (default: `true`)
- `PULUMI_ACCESS_TOKEN` - _(Optional)_ Pulumi Cloud access token
- `PULUMI_BACKEND_URL` - _(Optional)_ Pulumi Cloud API URL

### Optional

- `CORS_ALLOWED_ORIGINS` - Comma-separated list of allowed CORS origins
- `ENABLE_CORS` - Enable/disable CORS middleware (default: `true`)

## Development

### Build

```bash
make generate  # Generate proto code
make build     # Build binary
```

### Run Locally

**Quick Start (Recommended):**

```bash
./start-local.sh
```

This script automatically configures Pulumi for local file-based storage (no passphrase) and starts the server.

**Manual Start:**

```bash
# Configure Pulumi for local backend
export PULUMI_CONFIG_PASSPHRASE=""
pulumi login --local

# Set environment variables
export MONGODB_URI=mongodb://localhost:27017/openmcf
export SERVER_PORT=50051
export PULUMI_SKIP_UPDATE_CHECK=true
export PULUMI_AUTOMATION_API_SKIP_VERSION_CHECK=true

# Run the server
make dev
```

**Prerequisites:**
- MongoDB running locally on port 27017
- Pulumi CLI installed

### Docker Build

The backend is built as part of the unified container image:

```bash
# From project root
docker build -f app/Dockerfile.unified -t openmcf:latest .
```

**Note:** There is no separate backend-only Docker image. The unified container includes MongoDB, backend, and frontend together.

## Troubleshooting

### PULUMI_CONFIG_PASSPHRASE Error

**Error:**
```
error: getting stack configuration: get stack secrets manager: passphrase must be set with PULUMI_CONFIG_PASSPHRASE or PULUMI_CONFIG_PASSPHRASE_FILE environment variables
```

**Cause:** Pulumi requires a passphrase for encrypting secrets. This is required even for local development.

**Solution:** Use the same passphrase as the Docker container and login to local backend:

```bash
export PULUMI_CONFIG_PASSPHRASE="openmcf-default-passphrase"
pulumi login --local
```

Or use the provided helper script:
```bash
./start-local.sh
```

### PULUMI_ACCESS_TOKEN Error

**Error:**
```
failed to initialize stack: exit status 255, stderr: error: PULUMI_ACCESS_TOKEN must be set for login during non-interactive CLI sessions
```

**Cause:** Pulumi is not configured with a backend (either local or cloud).

**Solution:** Ensure `entrypoint.sh` runs `pulumi login --local` on startup. This is the default behavior.

### Stack Lock Errors

**Error:**
```
error: the stack is currently locked by ...
```

**Cause:** A previous deployment was interrupted and left a lock file.

**Solution:** The backend automatically runs `pulumi cancel` before each deployment. If manual intervention is needed:

```bash
docker exec -it openmcf pulumi cancel --stack <stack-fqdn> --yes
```

### Out of Disk Space

**Error:**
```
no space left on device
```

**Cause:** Go build cache or Pulumi plugins filled the container's filesystem.

**Solution:** The backend uses Docker volumes for Go cache and Pulumi state. If space issues persist:

```bash
# Clean up old volumes
docker volume prune

# Or remove and recreate volumes
docker-compose down -v
docker-compose up -d
```

## API Documentation

The backend exposes Connect-RPC endpoints:

- `/org.openmcf.app.cloudresource.v1.CloudResourceCommandController/*`
- `/org.openmcf.app.cloudresource.v1.CloudResourceQueryController/*`
- `/org.openmcf.app.stackupdate.v1.StackUpdateCommandController/*`
- `/org.openmcf.app.stackupdate.v1.StackUpdateQueryController/*`
- `/org.openmcf.app.credential.v1.CredentialCommandController/*`
- `/org.openmcf.app.credential.v1.CredentialQueryController/*`

Health check endpoint:

- `GET /health` - Returns JSON with service status

## License

See project root LICENSE file.

