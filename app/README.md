# OpenMCF Web App

<!-- Auto-release test: App change triggers v{semver}-app-{YYYYMMDD}.{N} tag format -->

Web interface for managing cloud resources and deployments with OpenMCF.

---

## Architecture

This directory contains both backend and frontend applications that together provide a complete web-based infrastructure management platform.

### Components

```
app/
├── backend/           # Go backend (Connect-RPC API)
├── frontend/          # Next.js frontend (React + TypeScript)
├── Dockerfile.unified # Production: Single container with all services
├── supervisord.conf   # Process manager configuration
└── entrypoint-unified.sh  # Container startup script
```

### Deployment Modes

#### Production/Testing (Unified Container)
- **Use case:** End users and integration testing
- **Setup:** Single Docker container with MongoDB, backend, and frontend
- **Managed by:** CLI (`planton webapp init/start/stop`) or `docker-compose up`
- **Image:** `ghcr.io/plantonhq/openmcf:latest`
- **Benefits:** Simple deployment, consistent environment, all services pre-configured

#### Development (Local Services)
- **Use case:** Contributors and active development
- **Setup:** Run MongoDB locally + backend via `go run` + frontend via `yarn dev`
- **Benefits:** Hot reload, faster iteration, easier debugging, native tooling
- **See:** [Development Setup](#development-setup) section below

---

## Backend (`app/backend/`)

### Technology Stack

- **Language:** Go 1.24.7
- **Framework:** Connect-RPC (Protocol Buffers)
- **Database:** MongoDB
- **Infrastructure:** Pulumi CLI integration

### Key Features

- Cloud resource CRUD operations
- Asynchronous Pulumi deployments
- Credential management with automatic resolution
- Stack update tracking and streaming
- Server-side pagination

### Structure

```
backend/
├── cmd/server/        # Main entry point
├── internal/
│   ├── database/      # MongoDB repositories
│   ├── server/        # Connect-RPC server setup
│   └── service/       # Business logic (services)
├── apis/             # Proto definitions and generated code
└── pkg/models/       # Data models
```

### Development

```bash
cd app/backend

# Generate proto code
make generate

# Run locally
go run cmd/server/main.go

# Build
make build

# Run tests
go test ./...
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `50051` | Backend API port |
| `MONGODB_URI` | `mongodb://localhost:27017` | MongoDB connection string |
| `MONGODB_DATABASE` | `openmcf` | Database name |
| `PULUMI_HOME` | `/home/appuser/.pulumi` | Pulumi state directory |
| `PULUMI_CONFIG_PASSPHRASE` | `openmcf-default-passphrase` | Pulumi encryption key |

---

## Frontend (`app/frontend/`)

### Technology Stack

- **Framework:** Next.js 14 (App Router)
- **Language:** TypeScript
- **UI Library:** Material-UI (MUI)
- **State:** React Context API
- **RPC Client:** Connect-RPC
- **Package Manager:** Yarn 3

### Key Features

- Cloud resource management interface
- Dark/light theme system (200+ color definitions)
- Stack update tracking and detail pages
- Real-time deployment output display
- Server-side pagination
- YAML editor with syntax highlighting

### Structure

```
frontend/
├── src/
│   ├── app/              # Next.js pages (App Router)
│   │   ├── dashboard/
│   │   ├── cloud-resources/
│   │   └── stack-updates/
│   ├── components/       # Reusable UI components
│   │   ├── layout/       # Header, sidebar, theme switch
│   │   └── shared/       # Tables, drawers, dialogs
│   ├── contexts/         # React contexts (theme, app state)
│   ├── themes/           # Color schemes (dark/light)
│   ├── gen/proto/        # Generated proto TypeScript code
│   └── hooks/            # Custom React hooks
├── public/              # Static assets
└── pages/api/           # API routes (health check)
```

### Development

```bash
cd app/frontend

# Install dependencies
yarn install

# Generate proto code (from backend)
cd ../backend/apis && buf generate

# Run development server
yarn dev

# Build for production
yarn build

# Run production build
yarn start
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `NEXT_PUBLIC_API_URL` | `http://localhost:50051` | Backend API URL |
| `PORT` | `3000` | Frontend port |
| `NODE_ENV` | `production` | Node environment |

---

## Unified Container

### Architecture

The unified container runs three services using supervisord:

```
supervisord (PID 1)
├── MongoDB (priority 1)
│   └── Port 27017 (localhost only)
├── Backend (priority 2, depends on MongoDB)
│   └── Port 50051 (exposed)
└── Frontend (priority 3, depends on Backend)
    └── Port 3000 (exposed)
```

### Building the Unified Image

```bash
# From project root
docker build -f app/Dockerfile.unified -t openmcf:dev .

# Test locally
docker run -p 3000:3000 -p 50051:50051 \
  -v mongodb-data:/data/db \
  openmcf:dev
```

### Multi-stage Build

1. **Stage 1:** Build backend Go binary
2. **Stage 2:** Build frontend Next.js app
3. **Stage 3:** Combine all into Ubuntu base with MongoDB

### Image Size

- Base Ubuntu: ~77MB
- MongoDB: ~150MB
- Node.js: ~50MB
- Backend binary: ~50MB
- Frontend build: ~100MB
- **Total: ~500MB**

---

## Development Setup

### Prerequisites

- Docker & Docker Compose
- Go 1.24.7+
- Node.js 20+
- Yarn 3
- MongoDB (for local development)

### Quick Start

```bash
# 1. Start MongoDB (if not using Docker)
brew install mongodb-community
brew services start mongodb-community

# 2. Start backend
cd app/backend
make generate
go run cmd/server/main.go

# 3. Start frontend (in another terminal)
cd app/frontend
yarn install
yarn dev

# 4. Access http://localhost:3000
```

### Docker Compose (Unified Container)

```bash
# Start unified container (MongoDB + Backend + Frontend)
docker-compose up

# Rebuild after changes (or pull latest image)
docker-compose up --build

# View logs (all services in one container)
docker-compose logs -f planton

# Stop all services
docker-compose down

# Clean up volumes (removes Pulumi state and MongoDB data)
docker-compose down -v
```

---

## API Development

### Adding a New RPC Method

1. **Define proto** (`app/backend/apis/proto/your_service.proto`):
```protobuf
service YourService {
  rpc CreateResource(CreateResourceRequest) returns (CreateResourceResponse);
}
```

2. **Generate code**:
```bash
cd app/backend && make generate
cd ../frontend && # proto code auto-generated
```

3. **Implement backend** (`app/backend/internal/service/your_service.go`):
```go
func (s *YourService) CreateResource(ctx context.Context, req *pb.CreateResourceRequest) (*pb.CreateResourceResponse, error) {
    // Implementation
}
```

4. **Register service** (`app/backend/internal/server/server.go`):
```go
yourService := service.NewYourService(cfg.MongoDB)
path, handler := yourServiceConnect.NewYourServiceHandler(yourService)
mux.Handle(path, handler)
```

5. **Use in frontend** (`app/frontend/src/app/your-page/page.tsx`):
```typescript
const client = useConnectRpcClient();
const response = await client.yourService.createResource(request);
```

---

## Database Schema

### Collections

#### `cloud_resources`
```javascript
{
  _id: ObjectId,
  id: "uuid",
  name: "resource-name",
  kind: "GcpCloudSql",
  manifest: { /* full YAML as JSON */ },
  created_at: ISODate,
  updated_at: ISODate
}
```

#### `credentials`
```javascript
{
  _id: ObjectId,
  provider: "gcp",
  gcp_creds: { /* GCP credentials */ },
  created_at: ISODate
}
```

#### `stackupdates`
```javascript
{
  _id: ObjectId,
  id: "uuid",
  cloud_resource_id: "uuid",
  status: "running|completed|failed",
  created_at: ISODate,
  completed_at: ISODate
}
```

#### `stackupdate_streaming_responses`
```javascript
{
  _id: ObjectId,
  stack_update_id: "uuid",
  output: "deployment log line",
  timestamp: ISODate,
  sequence: 123
}
```

---

## Testing

### Backend Tests

```bash
cd app/backend
go test ./...

# With coverage
go test -cover ./...
```

### Frontend Tests

```bash
cd app/frontend

# Unit tests
yarn test

# E2E tests (if implemented)
yarn test:e2e
```

### Integration Testing

```bash
# Start all services
docker-compose up -d

# Run integration tests
./scripts/integration-test.sh

# Clean up
docker-compose down
```

---

## Contributing

### Code Style

- **Backend:** Follow Go standard formatting (`gofmt`, `golint`)
- **Frontend:** Prettier + ESLint configuration included
- **Commits:** Conventional commits format

### Pull Request Process

1. Create feature branch
2. Make changes with tests
3. Ensure all tests pass
4. Update documentation
5. Submit PR with description

### Protocol Buffer Changes

- Always update proto files first
- Run code generation
- Update both backend and frontend
- Test end-to-end

---

## Deployment

### Building Production Image

Images are automatically built and published to GitHub Container Registry (GHCR) via GitHub Actions.

**Automated build (recommended):**
1. Navigate to GitHub Actions
2. Run "Build and Push Docker Image to GHCR" workflow
3. Specify version tag (optional)

**Manual build (if needed):**
```bash
# Build locally
docker build -f app/Dockerfile.unified -t ghcr.io/plantonhq/openmcf:v1.0.0 .

# Tag as latest
docker tag ghcr.io/plantonhq/openmcf:v1.0.0 ghcr.io/plantonhq/openmcf:latest
```

### Release Process

1. Update version in project
2. Push git tag: `git tag v1.0.0 && git push origin v1.0.0` (future: automatic builds)
3. Or manually trigger GitHub Actions workflow (current: testing phase)
4. GitHub Actions builds multi-arch image and pushes to GHCR
5. Update Homebrew formula (if CLI needs update)
6. Create GitHub release with changelog

---

## Troubleshooting

### Backend Won't Start

- Check MongoDB connection: `mongosh`
- Verify environment variables
- Check logs: `docker-compose logs backend`

### Frontend Build Fails

- Clear node_modules: `rm -rf node_modules && yarn install`
- Clear Next.js cache: `rm -rf .next`
- Regenerate proto code

### Docker Build Fails

- Increase Docker memory/CPU limits
- Clear build cache: `docker builder prune`
- Check disk space

---

## Performance

### Backend
- Connection pooling for MongoDB
- Concurrent deployment execution
- Streaming responses for large outputs

### Frontend
- Server-side rendering (SSR)
- Code splitting by route
- Image optimization
- Bundle size monitoring

---

## Security

### Current (Development)
- ⚠️ No authentication
- ⚠️ MongoDB without password (localhost only)
- ⚠️ No HTTPS
- ⚠️ No rate limiting

### Future (Production)
- User authentication (OAuth2/JWT)
- Role-based access control
- Encrypted credentials at rest
- HTTPS/TLS
- API rate limiting
- Audit logging

---

## License

This project is part of OpenMCF. See main repository for license information.

---

## Getting Help

- Check documentation: `_projects/20251127-openmcf-web-app/docs/`
- View logs: `docker-compose logs -f`
- Open issue on GitHub
