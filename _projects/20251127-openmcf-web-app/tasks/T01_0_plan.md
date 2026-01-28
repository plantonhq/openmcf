# T01: Initial Task Plan - OpenMCF Web App

**Status:** 📋 PENDING REVIEW  
**Created:** November 27, 2025  
**Plan Type:** Initial breakdown for full-stack web application

---

## Overview

This plan breaks down the development of OpenMCF Web App into manageable phases, balancing the immediate need for a working demo (2 days) with the long-term goal of a production-ready full-stack application.

**Strategy**: Frontend-first with mock data, backend requirements emerge naturally, full integration follows.

---

## High-Level Task Breakdown

### Phase 1: Foundation & Mock Frontend (Priority: URGENT)
**Timeline:** Day 1 (Today - November 27)  
**Goal:** Working web interface with mock data for meetup demo

#### Task 1.1: Project Setup & Infrastructure
**Estimated Time:** 2-3 hours

**Frontend Setup:**
- [ ] Initialize Next.js project with TypeScript
- [ ] Configure TailwindCSS
- [ ] Set up project structure (pages, components, lib)
- [ ] Configure Protocol Buffers code generation
- [ ] Set up Connect-RPC client stub generation

**Backend Setup:**
- [ ] Initialize Go module for backend service
- [ ] Configure Protocol Buffers + Connect-RPC (Golang)
- [ ] Set up project structure (cmd, internal, api)
- [ ] Configure buf.yaml for proto generation

**Deliverables:**
- Working dev environments (frontend + backend)
- Shared proto definitions
- Build scripts and Makefiles

---

#### Task 1.2: Define Core Data Models (Protocol Buffers)
**Estimated Time:** 2-3 hours

**API Definitions:**
- [ ] Define `DeploymentComponent` message
- [ ] Define `CloudResource` message
- [ ] Define `Provider` enum (AWS, GCP, Azure, Kubernetes, etc.)
- [ ] Define `ResourceStatus` enum and messages
- [ ] Define query/list/filter messages

**Service Definitions (Connect-RPC):**
- [ ] `DeploymentComponentService` (list, get, search)
- [ ] `CloudResourceService` (create, list, get, update, delete)
- [ ] `StackUpdateService` (get status, list jobs)

**Generate Code:**
- [ ] Run buf generate for TypeScript (frontend)
- [ ] Run buf generate for Go (backend)

**Deliverables:**
- Complete proto definitions
- Generated client stubs (TypeScript)
- Generated server stubs (Go)

---

#### Task 1.3: Create Mock Data Layer
**Estimated Time:** 1-2 hours

**Mock Data Files:**
- [ ] Create mock deployment components (10-15 examples)
  - PostgresKubernetes
  - AwsVpc
  - GcpGkeCluster
  - RedisKubernetes
  - AwsEksCluster
  - etc.

- [ ] Create mock cloud resources (5-10 examples)
  - Sample deployed resources
  - Various states (pending, running, failed)
  - Realistic outputs

- [ ] Create mock stack-updates (3-5 examples)
  - In-progress deployments
  - Completed deployments
  - Failed deployments with errors

**Mock Service Implementation (TypeScript):**
- [ ] `mockDeploymentComponentService.ts`
- [ ] `mockCloudResourceService.ts`
- [ ] `mockStackUpdateService.ts`

**Deliverables:**
- Rich, realistic mock data
- Frontend can run completely standalone

---

#### Task 1.4: Build Frontend UI Components
**Estimated Time:** 4-6 hours

**Core Layout:**
- [ ] Main layout with navigation
- [ ] Sidebar with sections (Deployment Store, Resources, Jobs)
- [ ] Header with search and actions

**Deployment Component Store:**
- [ ] Grid/list view of deployment components
- [ ] Filter by provider (AWS, GCP, Azure, K8s)
- [ ] Search by name/description
- [ ] Component detail modal/page
- [ ] "Deploy" button (opens creation form)

**Cloud Resource Management:**
- [ ] List view of deployed resources
- [ ] Filter by kind, provider, status
- [ ] Resource detail view
- [ ] Status indicators (running, pending, failed)
- [ ] Actions (update, delete, view logs)

**Stack Job Monitoring:**
- [ ] List of recent deployments
- [ ] Real-time status indicators
- [ ] Log viewer (mock streaming logs)
- [ ] Job history

**Resource Creation Flow:**
- [ ] Multi-step form for creating resources
- [ ] Provider selection
- [ ] Component selection
- [ ] Configuration form (based on schema)
- [ ] Validation preview
- [ ] Confirmation and deploy

**Deliverables:**
- Complete, functional UI
- Works entirely with mock data
- Beautiful, modern design
- Demonstrates all key workflows

---

### Phase 2: Backend API Implementation (Priority: HIGH)
**Timeline:** Day 2 (November 28)  
**Goal:** Real backend services with database persistence

#### Task 2.1: MongoDB Setup & Schema Design
**Estimated Time:** 2-3 hours

**Database Design:**
- [ ] Install/configure MongoDB (local + Docker)
- [ ] Design collections:
  - `deployment_components` (catalog)
  - `cloud_resources` (user resources)
  - `stack_jobs` (deployment history)
  - `users` (future: authentication)

**Schema Considerations:**
- [ ] Use Protocol Buffer messages as schema guide
- [ ] Index strategy (queries from frontend)
- [ ] Relationships (resources → jobs, components → resources)

**Go Database Layer:**
- [ ] MongoDB connection setup
- [ ] Repository pattern implementation
- [ ] CRUD operations for each collection
- [ ] Query builders for filters/search

**Deliverables:**
- Running MongoDB instance
- Database schemas
- Go repository layer

---

#### Task 2.2: Implement Backend Services (Connect-RPC)
**Estimated Time:** 4-6 hours

**Service Implementation:**
- [ ] `DeploymentComponentService` implementation
  - List with filters
  - Get by ID
  - Search by keyword
  - Pagination support

- [ ] `CloudResourceService` implementation
  - Create new resource
  - List user's resources
  - Get resource details
  - Update resource
  - Delete resource

- [ ] `StackUpdateService` implementation
  - Get job status
  - List jobs for resource
  - Stream logs (future: real-time)

**Business Logic:**
- [ ] Validation (use proto-validate)
- [ ] Status state machines
- [ ] Error handling and responses
- [ ] Logging

**Deliverables:**
- Fully functional backend APIs
- Database persistence
- Connect-RPC handlers

---

#### Task 2.3: CLI Integration with Backend
**Estimated Time:** 3-4 hours

**CLI Modifications:**
- [ ] Add backend connection configuration
- [ ] Implement API clients for backend services
- [ ] Modify existing commands to use backend:
  - `openmcf list` → query backend
  - `openmcf create` → call backend API
  - `openmcf get` → fetch from backend

**Backward Compatibility:**
- [ ] CLI works standalone (without backend)
- [ ] Flag to enable/disable backend mode
- [ ] Graceful fallback if backend unavailable

**Data Synchronization:**
- [ ] CLI operations write to backend
- [ ] Backend operations accessible via CLI
- [ ] Shared data model (Protocol Buffers)

**Deliverables:**
- CLI integrated with backend
- Maintains standalone capability
- All operations persisted

---

### Phase 3: Integration & Polish (Priority: MEDIUM)
**Timeline:** Day 2 Evening + Ongoing  
**Goal:** Connect frontend to real backend, polish for demo

#### Task 3.1: Connect Frontend to Real Backend
**Estimated Time:** 2-3 hours

**Replace Mock Services:**
- [ ] Update frontend to use real Connect-RPC clients
- [ ] Point to backend URL (localhost:8080 for dev)
- [ ] Remove/toggle mock data layer

**Connection Configuration:**
- [ ] Environment variable for backend URL
- [ ] CORS configuration
- [ ] Error handling for backend failures

**Real-Time Updates:**
- [ ] Polling for status updates (short-term)
- [ ] Future: WebSocket/SSE for real-time

**Deliverables:**
- Frontend talks to real backend
- Full end-to-end flow working
- Mock mode still available

---

#### Task 3.2: Deployment & Docker Setup
**Estimated Time:** 2-3 hours

**Containerization:**
- [ ] Dockerfile for backend service
- [ ] Dockerfile for frontend (Next.js)
- [ ] Docker Compose setup:
  - Backend service
  - Frontend service
  - MongoDB
  - Network configuration

**Local Development:**
- [ ] docker-compose.dev.yml for development
- [ ] Volume mounts for hot-reload
- [ ] Environment variable management

**Production Preparation:**
- [ ] docker-compose.prod.yml
- [ ] Environment-specific configs
- [ ] Health checks

**Deliverables:**
- One-command local setup
- Production-ready containers
- Easy deployment

---

#### Task 3.3: Testing & Documentation
**Estimated Time:** 2-3 hours

**Testing:**
- [ ] Backend unit tests (key services)
- [ ] Integration tests (API endpoints)
- [ ] Frontend component tests (critical paths)
- [ ] End-to-end smoke tests

**Documentation:**
- [ ] API documentation (from proto files)
- [ ] Setup guide (README)
- [ ] Development workflow guide
- [ ] Deployment instructions

**Demo Preparation:**
- [ ] Seed data for demo
- [ ] Demo script/scenarios
- [ ] Screenshots/recordings
- [ ] Talking points

**Deliverables:**
- Test coverage for critical paths
- Comprehensive documentation
- Demo-ready environment

---

## Minimum Viable Demo (for Meetup)

If time is tight, **absolute minimum for meetup demo**:

✅ **Must Have:**
1. Frontend with mock data showing:
   - Deployment component catalog
   - Resource list view
   - Creation flow (even if mocked)
   - Visual polish

2. Clear vision demonstrated:
   - What OpenMCF does
   - How web UI makes it accessible
   - Value proposition obvious

❌ **Nice to Have (can come after meetup):**
- Real backend integration
- Database persistence
- CLI integration
- Full deployment capability

**Fallback Position**: Frontend-only demo with mock data is still valuable and impressive.

---

## Success Metrics

### For Meetup (November 29)
- [ ] Working web interface accessible in browser
- [ ] Attendees can interact with UI
- [ ] Demonstrates OpenMCF's capabilities
- [ ] Visual appeal and polish
- [ ] Clear value proposition

### For Production (Long-term)
- [ ] Full backend persistence
- [ ] CLI integration complete
- [ ] Multi-user support
- [ ] Authentication/authorization
- [ ] Real deployment capability
- [ ] Production deployment

---

## Resource Allocation

**Frontend Work:** ~12-16 hours
- UI components and layouts
- Mock data integration
- Visual design and polish
- Connect-RPC client integration

**Backend Work:** ~10-14 hours
- Database setup and schema
- Service implementation
- CLI integration
- Testing

**Integration Work:** ~6-8 hours
- Frontend ↔ Backend connection
- Docker setup
- Testing and debugging
- Documentation

**Total Estimated:** ~30-40 hours of work
**Available:** ~48 hours until meetup (with parallel work)

**Realistic Assessment**: Can deliver frontend + partial backend by meetup. Full integration shortly after.

---

## Risks & Mitigations

### Risk: Backend not ready for meetup
**Probability:** Medium  
**Impact:** Low (frontend with mocks still works)  
**Mitigation:** Frontend-first approach ensures demo works regardless

### Risk: Integration issues at the last minute
**Probability:** Medium  
**Impact:** Medium  
**Mitigation:** Frontend can run standalone with mocks, integrate post-meetup

### Risk: Database schema changes during development
**Probability:** High  
**Impact:** Low  
**Mitigation:** MongoDB flexibility, start simple

### Risk: Connect-RPC learning curve
**Probability:** Medium  
**Impact:** Medium  
**Mitigation:** Excellent docs, simpler than gRPC-web setup

---

## Next Steps

**Immediate Actions (once plan approved):**

1. **Start Frontend Setup** (parallel workstream)
   - Initialize Next.js project
   - Set up TailwindCSS
   - Create project structure

2. **Start Backend Setup** (parallel workstream)
   - Initialize Go module
   - Set up Protocol Buffers
   - Configure Connect-RPC

3. **Define Proto APIs** (required by both)
   - Core message types
   - Service definitions
   - Generate code for both sides

4. **Create Mock Data** (unblocks frontend)
   - Realistic deployment components
   - Sample cloud resources
   - Mock stack-updates

5. **Build UI Components** (main frontend work)
   - Layout and navigation
   - Component catalog
   - Resource management views

---

## Dependencies Between Tasks

```
Proto Definitions (1.2)
    ↓
    ├─→ Frontend Setup (1.1) → Mock Data (1.3) → UI Components (1.4)
    └─→ Backend Setup (1.1) → Database (2.1) → Services (2.2) → CLI Integration (2.3)
                                                      ↓
                                                Integration (3.1) → Docker (3.2) → Testing (3.3)
```

**Critical Path**: Proto Definitions → Mock Data → UI Components (for demo)

**Parallel Tracks**: Frontend and backend can work independently after proto definitions.

---

**Plan Status:** 📋 PENDING REVIEW

**Please review this plan and provide feedback. Once approved, execution will begin with T01_3_execution.md**

