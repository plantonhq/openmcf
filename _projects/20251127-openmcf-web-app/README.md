# OpenMCF Web App

**Status:** 🚧 In Development
**Started:** November 27, 2025
**First Milestone:** November 29, 2025 (CNCF Hyderabad Meetup)
**Type:** Long-term Feature Development

---

## TL;DR

Building a comprehensive full-stack web application for OpenMCF to transform the CLI-only tool into a production-ready product with web interface, backend API layer, and persistent database storage. First working version targeted for CNCF Hyderabad meetup in 2 days, but this is a long-term project that will evolve with the product.

---

## Project Description

Building a comprehensive full-stack web application for OpenMCF with frontend UI, backend API layer, and database integration to transform the CLI-only tool into a production-ready product with persistence.

Currently, OpenMCF exists as a CLI tool without backend integration or persistence capabilities. All operations are ephemeral and CLI-driven. This project will add:

1. **Web frontend** - Intuitive UI exposing all CLI capabilities
2. **Backend services** - API layer with persistence
3. **Database integration** - Persistent storage for all operations
4. **CLI integration** - Connect existing CLI to backend

---

## Primary Goal

Transform Planton from a CLI-only tool into a self-contained, production-ready product with a complete web interface, backend API layer, and persistent database storage that users will find immediately useful and want to try out.

---

## Timeline

### Immediate Milestone (November 29, 2025)
- **2 days from now**: CNCF Hyderabad meetup demo
- **Target**: Working web interface that demonstrates value
- **Approach**: Frontend with mock data if needed, backend in progress

### Long-Term Vision
- This is NOT a one-off demo project
- Architecture designed for production use
- Will evolve with OpenMCF for its lifetime
- Meetup is first milestone, not final delivery

---

## Technology Stack

### Frontend
- **Framework**: Next.js
- **Language**: TypeScript/React
- **API Client**: Connect-RPC (TypeScript client)
- **State Management**: React hooks + Context
- **Styling**: TailwindCSS (following modern patterns)

### Backend
- **Language**: Golang
- **RPC Framework**: Connect-RPC (Buf Connect)
- **Why Connect-RPC**: Avoids gRPC-web complexity (no Envoy/Istio needed)
- **API Definitions**: Protocol Buffers

### Database
- **Primary**: MongoDB
- **Schema**: To be designed based on CLI operations
- **Driver**: Official MongoDB Go driver

### Infrastructure
- **Containerization**: Docker
- **Orchestration**: Kubernetes (for production)
- **Local Dev**: Docker Compose

---

## Project Type

**Feature Development** - Building major new capabilities across multiple components.

---

## Affected Components

1. **OpenMCF CLI** (existing Go codebase)
   - Integration with backend
   - Add persistence layer
   - Maintain CLI usability

2. **New Web Frontend** (Next.js application)
   - User interface for all CLI operations
   - List views and data visualization
   - Resource creation and management

3. **New Backend Services** (Golang + Connect-RPC)
   - RPC-based APIs
   - Business logic layer
   - Data persistence operations

4. **New Database Layer** (MongoDB)
   - Schema design
   - Data models
   - Query patterns

5. **Integration Layer**
   - CLI ↔ Backend
   - Frontend ↔ Backend
   - Shared data models (Protocol Buffers)

---

## Success Criteria

### For Meetup Milestone (November 29, 2025)
- ✅ Working web interface accessible in browser
- ✅ Clear demonstration of OpenMCF's value
- ✅ Something attendees can interact with and try
- ✅ Visual representation of deployment components
- ✅ Mock data acceptable if backend not fully ready

### For Long-Term Production Readiness
- ✅ Full backend persistence (all operations stored)
- ✅ CLI integration with backend
- ✅ Real-time data synchronization
- ✅ Authentication and authorization
- ✅ Multi-user support
- ✅ Production deployment capability
- ✅ Comprehensive documentation

---

## Technical Architecture

### Three-Phase Development Strategy

#### Phase 1: Frontend with Mock Data (Parallel Work)
**Goal**: Don't wait for backend - build UI with dummy data

**Deliverables**:
- Complete web interface with all views
- Mock data representing deployment components
- User interaction flows
- Component designs and layouts

**Benefit**: Frontend work proceeds independently, reveals exact backend requirements

#### Phase 2: Backend Requirements Identification
**Goal**: Define precise backend APIs needed

**Approach**:
- Document all RPC calls frontend needs
- Use CLI operations as blueprint
- Map CLI functionality to backend APIs
- Define data models and schemas

**Output**: Clear specification for backend implementation

#### Phase 3: Backend and Database Implementation
**Goal**: Build real backend with persistence

**Deliverables**:
- Backend RPC services (Connect-RPC)
- MongoDB schemas and collections
- CLI integration with backend
- Frontend connection to real APIs

**Result**: Fully integrated system with persistence

---

## Key Technical Decisions

### Why Connect-RPC Over gRPC?

**Problem with gRPC-web**:
- Requires Envoy proxy or Istio service mesh
- Additional infrastructure complexity
- Harder to deploy and maintain

**Connect-RPC Benefits**:
- Works directly over HTTP/1.1 and HTTP/2
- No proxy required
- Browser-compatible out of the box
- Same Protocol Buffers definitions
- Simpler deployment

### Why MongoDB?

**Alignment with Planton Cloud**:
- Same database as Planton Cloud InfraHub
- Proven patterns available
- Flexible schema for cloud resources
- Good query performance

**Development Benefits**:
- Schema flexibility during early iterations
- Easy to add fields
- Natural fit for nested resource structures

### Why Mock-First Frontend?

**Speed**:
- Frontend development can start immediately
- No waiting for backend to be ready

**Requirements Discovery**:
- Building UI reveals what backend needs to do
- Prevents over-engineering backend
- Frontend team defines API contract

**Demo Capability**:
- Working UI can demo even without backend
- Visual progress for stakeholders
- Easier to iterate on UX

---

## Dependencies and Blockers

### External Dependencies
- Protocol Buffers (buf CLI)
- Connect-RPC code generation
- MongoDB instance (local + production)
- Next.js and Go toolchains

### Internal Dependencies
- Existing OpenMCF CLI codebase
- CloudResourceKind enum and definitions
- Deployment component schemas
- Validation rules from proto definitions

### Timeline Constraint
- **2 days to meetup**: Aggressive timeline
- **Mitigation**: Mock-first approach ensures something to show
- **Realistic**: Full integration may take longer, but demo will work

### Integration Complexity
- Three major components must work together
- Data models must align across CLI, backend, frontend
- Protocol Buffers provide type safety across boundaries

---

## Risks and Mitigation

### Risk: Cannot finish backend in 2 days

**Mitigation**:
- Frontend with mock data is still valuable demo
- Shows vision and capability
- Backend can be finished post-meetup

### Risk: Three components difficult to integrate

**Mitigation**:
- Protocol Buffers provide shared type system
- Connect-RPC simplifies RPC layer
- Start with small end-to-end slice

### Risk: CLI integration changes break existing functionality

**Mitigation**:
- CLI continues to work standalone
- Backend is additive, not replacement
- Comprehensive testing before merging

### Risk: Database schema changes during development

**Mitigation**:
- MongoDB schema flexibility
- Start simple, iterate
- Version migrations if needed

---

## Related Work and Context

### OpenMCF Architecture
- Open-source multi-cloud deployment framework
- 100+ deployment components
- Protocol Buffers-based APIs
- Dual IaC support (Pulumi + Terraform)

### Planton Cloud Reference
- Commercial platform built on OpenMCF
- InfraHub architecture patterns available
- MongoDB + Neo4j proven in production
- gRPC-based backend (switching to Connect-RPC)

### CNCF Hyderabad Meetup Context
- Speaking opportunity November 29, 2025
- Platform engineering audience
- Demo of OpenMCF OSS
- First public showcase of web interface

---

## Special Requirements

### Code Quality
- Follow OpenMCF Go coding guidelines
- Use Protocol Buffers for all API definitions
- Maintain feature parity with CLI
- Comprehensive error handling

### Documentation
- API documentation from proto files
- Frontend component documentation
- Backend service documentation
- Deployment guides

### Testing
- Unit tests for backend services
- Integration tests for API endpoints
- Frontend component tests
- End-to-end workflows

### Deployment
- Dockerized services
- Docker Compose for local development
- Kubernetes manifests for production
- CI/CD pipeline (future)

---

## Next Steps

See `next-task.md` for immediate next actions and current task status.

To resume work on this project in any session, simply drag `next-task.md` into the chat.

---

## Project Structure

```
_projects/20251127-openmcf-web-app/
├── README.md                      # This file - project overview
├── next-task.md                   # Quick resume file (drag into chat)
├── tasks/                         # Task management
│   ├── T01_0_plan.md             # Initial task plan (archived)
│   ├── T02_4_completion.md       # ✅ Cloud Resource Web UI (Dec 1)
│   ├── T03_4_completion.md       # ✅ Dashboard Simplification (Dec 1)
│   ├── T04_4_completion.md       # ✅ UI Enhancements & Pagination (Dec 3)
│   ├── T05_4_completion.md       # ✅ Pulumi Stack Job API (Dec 3)
│   ├── T06_4_completion.md       # ✅ KubernetesRedis Fix (Dec 4)
│   ├── T07_4_completion.md       # ✅ Stack Jobs UI Integration (Dec 4)
│   └── T08_4_completion.md       # ✅ Credential Management (Dec 8-9)
├── checkpoints/                   # Milestone documentation
├── design-decisions/              # Architecture decisions
├── coding-guidelines/             # Project-specific guidelines
├── wrong-assumptions/             # Learnings and corrections
└── dont-dos/                      # Anti-patterns and gotchas
```

---

**Status**: ✅ **OPERATIONAL & PRODUCTION READY** - 7 major features completed (Dec 1-9, 2025)

**Last Updated**: December 9, 2025

