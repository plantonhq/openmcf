# Temporal Kubernetes Pulumi Module - Architecture Overview

This document describes the internal architecture, design decisions, and resource organization of the Temporal Kubernetes Pulumi module.

## Module Purpose

Deploy a production-ready Temporal cluster on Kubernetes with support for:
- Multiple database backends (Cassandra, PostgreSQL, MySQL)
- External or embedded database options
- Ingress for frontend (gRPC + HTTP) and Web UI
- External Elasticsearch integration
- Monitoring stack (Prometheus + Grafana)

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                     Temporal Kubernetes Module                       │
├─────────────────────────────────────────────────────────────────────┤
│                                                                       │
│  ┌──────────────┐    ┌──────────────────────────────────────────┐  │
│  │              │    │         Helm Chart Values                 │  │
│  │  Namespace   │    │  ┌────────────────────────────────────┐  │  │
│  │  Creation    │    │  │ Database Backend Selection          │  │  │
│  │              │    │  │  - Cassandra (embedded/external)    │  │  │
│  └──────┬───────┘    │  │  - PostgreSQL (embedded/external)   │  │  │
│         │            │  │  - MySQL (embedded/external)        │  │  │
│         │            │  └────────────────────────────────────┘  │  │
│         │            │                                           │  │
│  ┌──────▼───────┐    │  ┌────────────────────────────────────┐  │  │
│  │  DB Password │    │  │ Service Configuration               │  │  │
│  │  Secret      │◄───┼──│  - Frontend gRPC (7233)             │  │  │
│  │  (if ext DB) │    │  │  - Frontend HTTP (7243)             │  │  │
│  └──────┬───────┘    │  │  - Web UI (8080)                    │  │  │
│         │            │  └────────────────────────────────────┘  │  │
│         │            │                                           │  │
│  ┌──────▼───────┐    │  ┌────────────────────────────────────┐  │  │
│  │              │    │  │ Optional Features                   │  │  │
│  │  Helm Chart  │◄───┼──│  - Elasticsearch (embedded/ext)     │  │  │
│  │  Deployment  │    │  │  - Monitoring (Prom + Grafana)      │  │  │
│  │              │    │  │  - Auto Schema Setup                │  │  │
│  └──────┬───────┘    │  └────────────────────────────────────┘  │  │
│         │            └──────────────────────────────────────────┘  │
│         │                                                           │
│  ┌──────▼──────────────────────────────────────────────────────┐  │
│  │                 Ingress Resources                            │  │
│  │  ┌─────────────────┐  ┌──────────────┐  ┌────────────────┐ │  │
│  │  │ Frontend gRPC   │  │ Frontend HTTP│  │ Web UI HTTP    │ │  │
│  │  │ LoadBalancer    │  │ Gateway/Route│  │ Gateway/Route  │ │  │
│  │  │ (external-dns)  │  │ (Istio)      │  │ (Istio)        │ │  │
│  │  └─────────────────┘  └──────────────┘  └────────────────┘ │  │
│  └──────────────────────────────────────────────────────────────┘  │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## File Organization

The module is organized into focused, single-purpose files following Go conventions:

### Core Module Files

| File | Purpose | Key Responsibilities |
|------|---------|---------------------|
| `main.go` | Orchestration entry point | Wires all resources in dependency order |
| `locals.go` | Data transformations | Computes derived values, exports outputs |
| `variables.go` | Constants & defaults | Helm versions, ports, secret names |
| `outputs.go` | Output constants | Defines export key names |

### Resource-Specific Files

| File | Resource Type | Creates |
|------|---------------|---------|
| `namespace.go` | Namespace | Kubernetes namespace with labels (conditional) |
| `db_password_secret.go` | Secret | External database password (conditional) |
| `helm_chart.go` | Helm Chart | Temporal Helm chart with all values |
| `frontend_ingress.go` | LoadBalancer Service | gRPC frontend external access |
| `frontend_http_ingress.go` | Gateway/HTTPRoute | HTTP frontend via Istio |
| `web_ui_ingress.go` | Gateway/HTTPRoute | Web UI via Istio |

## Design Decisions

### 1. Conditional Database Password Secret

**Decision**: Only create the Kubernetes Secret when `external_database` is provided.

**Rationale**:
- Embedded databases (Cassandra, MySQL, PostgreSQL) manage their own credentials
- External databases require explicit password management
- Avoids creating unused resources

**Implementation**: `db_password_secret.go` checks for non-nil `ExternalDatabase`

### 2. Dual Ingress Strategy

**Decision**: Frontend uses LoadBalancer + external-dns; Web UI uses Gateway API (Istio)

**Rationale**:
- **Frontend gRPC**: Requires TCP LoadBalancer (port 7233) with DNS via external-dns annotation
- **Frontend HTTP**: Optional HTTP access via Gateway API HTTPRoute (port 7243)
- **Web UI HTTP**: Always HTTP-based, benefits from Istio's features (HTTPS redirect, cert management)

**Trade-offs**:
- Mixed ingress types add complexity but optimize for protocol requirements
- gRPC over HTTP/2 works best with direct TCP LoadBalancer
- HTTP services benefit from Gateway API features (path routing, TLS termination)

### 3. Database Backend Selection Logic

**Decision**: Helm chart values mutually exclude database backends

**Rationale**:
- Temporal Helm chart's embedded databases conflict if multiple are enabled
- Clear backend selection prevents misconfiguration
- External database requires disabling all embedded options

**Implementation** (`helm_chart.go`):
```
if ExternalDatabase != nil:
  - Disable all embedded (cassandra, mysql, postgresql)
  - Configure SQL driver based on backend enum
  - Set connection details and TLS settings
else:
  - Enable selected embedded database
  - Disable other embedded databases
```

### 4. Monitoring Stack Auto-Enable

**Decision**: Enable monitoring when `EnableMonitoringStack=true` OR `ExternalElasticsearch` is configured

**Rationale**:
- External Elasticsearch implies production use → needs monitoring
- Explicit flag allows monitoring without Elasticsearch
- Avoids deploying heavy monitoring stack in development

### 5. Schema Management

**Decision**: Auto schema setup is enabled by default, can be disabled via `DisableAutoSchemaSetup`

**Rationale**:
- Temporal requires schema initialization for new databases
- Production environments may want manual schema control
- Default behavior favors ease-of-use for development

### 6. Version Pinning

**Decision**: Default Helm chart version defined in `variables.go`, overridable via `spec.version`

**Rationale**:
- Stability: Known working version as default
- Flexibility: Allow upgrades via spec field
- Visibility: Version documented in one place

**Default**: `0.62.0` (tested and verified)

### 7. Conditional Namespace Creation

**Decision**: Namespace creation is controlled by `spec.create_namespace` boolean flag

**Rationale**:
- Production environments often have pre-existing namespaces with policies, quotas, RBAC
- Some organizations require namespace creation through separate governance processes
- Provides flexibility for both greenfield and brownfield deployments

**Implementation**: 
- When `create_namespace = true`: Module creates namespace with proper labels
- When `create_namespace = false`: Module assumes namespace exists, skips creation
- All child resources work identically in both modes
- Default is `true` for backward compatibility and ease of use

**Trade-offs**:
- Users must ensure namespace exists when `create_namespace = false`
- No validation that existing namespace has proper labels or configuration

## Data Flow

### Input → Locals → Resources

1. **Input**: `KubernetesTemporalStackInput` contains:
   - `Target` (KubernetesTemporal API resource)
   - `ProviderConfig` (Kubernetes connection details)
   - `KubernetesNamespace` (optional override)

2. **Locals Initialization** (`initializeLocals`):
   - Computes namespace (priority: stackInput > label > metadata.name)
   - Generates labels (resource, org, env)
   - Derives service names and endpoints
   - Extracts ingress hostnames
   - Exports all outputs immediately

3. **Resource Creation** (dependency order):
   ```
   namespace → db_password_secret → helm_chart → ingress_resources
   ```

### Output Exports

Outputs are exported **eagerly** during `initializeLocals` to ensure they're available even if resource creation fails:

| Export Key | Source | Type |
|------------|--------|------|
| `namespace` | Computed from metadata/stackInput | String |
| `frontend_service` / `ui_service` | Derived from metadata.name | String |
| `frontend_endpoint` / `ui_endpoint` | FQDN constructed from service + namespace | String |
| `port_forward_*_command` | Template command string | String |
| `external_frontend_hostname` | From ingress.frontend.grpcHostname (if enabled) | String |
| `external_ui_hostname` | From ingress.webUi.hostname (if enabled) | String |

## Resource Dependencies

### Explicit Dependencies

```
namespace (if create_namespace = true)
  ├─> db_password_secret (if external DB)
  │     └─> helm_chart
  └─> helm_chart
        ├─> frontend_ingress
        ├─> frontend_http_ingress
        └─> web_ui_ingress

If create_namespace = false, all resources deploy to existing namespace
```

### Implicit Dependencies (via values)

- Helm chart references `db_password_secret.name` in values (when external DB)
- Ingress resources select Helm-created pods via labels
- All resources use `locals.Namespace` for placement

## Helm Chart Value Construction

The Temporal Helm chart is configured via a `pulumi.Map` with nested structure:

### Key Value Groups

1. **Identity**: `fullnameOverride` → metadata.name
2. **Database**: Backend selection, connection details, TLS config
3. **Services**: Port configuration for frontend gRPC/HTTP
4. **Persistence**: Default and visibility database settings
5. **Schema**: Auto-setup flags (createDatabase, setup, update)
6. **Web UI**: Enable/disable flag
7. **Monitoring**: Prometheus, Grafana, KubePrometheusStack
8. **Elasticsearch**: External connection or embedded enable

### Example Value Mapping

| Spec Field | Helm Value Path | Notes |
|------------|-----------------|-------|
| `database.backend = postgresql` | `server.config.persistence.default.sql.driver = "postgres12"` | SQL driver string |
| `database.externalDatabase.host` | `server.config.persistence.default.sql.host` | Direct mapping |
| `database.externalDatabase.password` | `server.config.persistence.default.sql.existingSecret` | Via Secret reference |
| `disableWebUi = true` | `web.enabled = false` | Boolean inversion |
| `enableMonitoringStack = true` | `prometheus.enabled`, `grafana.enabled`, etc. | Fan-out to multiple values |

## Ingress Architecture Details

### Frontend gRPC Ingress (`frontend_ingress.go`)

**Resource**: Kubernetes LoadBalancer Service

**When Created**: `ingress.frontend.enabled == true` AND `grpc_hostname != ""`

**Configuration**:
- Service type: LoadBalancer
- Port: 7233 (gRPC)
- Annotation: `external-dns.alpha.kubernetes.io/hostname`
- Selector: Matches Temporal frontend pods created by Helm

**DNS Management**: external-dns controller watches annotation and creates DNS records

### Frontend HTTP Ingress (`frontend_http_ingress.go`)

**Resources**: Gateway + HTTPRoute (Kubernetes Gateway API)

**When Created**: `ingress.frontend.enabled == true` AND `http_hostname != ""`

**Configuration**:
- Gateway class: `istio`
- Protocol: HTTP/HTTPS (with redirect)
- Backend: Points to frontend service (port 7243)
- Hostname: `ingress.frontend.httpHostname`

**Purpose**: Allows HTTP-based access to Temporal's REST API

### Web UI Ingress (`web_ui_ingress.go`)

**Resources**: Gateway + HTTPRoute (Kubernetes Gateway API)

**When Created**: `ingress.webUi.enabled == true` AND `hostname != ""`

**Configuration**:
- Gateway class: `istio`
- Protocol: HTTP/HTTPS (with redirect)
- Backend: Points to Web UI service (port 8080)
- Hostname: `ingress.webUi.hostname`

**Purpose**: Exposes Temporal Web UI for workflow monitoring

## Error Handling

### Pre-Flight Validation

**Check**: External database required for PostgreSQL/MySQL backends

```go
if backend != cassandra && externalDatabase == nil {
    return errors.New("external_database must be provided when backend is not cassandra")
}
```

**Rationale**: Embedded PostgreSQL/MySQL are not recommended for production

### Resource Creation Errors

All resource creation functions return errors wrapped with context:

```go
if err := helmChart(ctx, locals, createdNamespace); err != nil {
    return errors.Wrap(err, "failed to install Temporal Helm chart")
}
```

This provides clear error traces when debugging deployment failures.

## Testing Strategy

### Validation Tests

Located in `api_test.go` (should be `spec_test.go` per convention):

1. **Database Backend Tests**: Validate each backend configuration
2. **Ingress CEL Rules**: Test conditional hostname requirements
3. **External Database**: Verify connection parameters

### Manual Testing

Use `hack/manifest.yaml` for local testing:

```bash
cd iac/pulumi
make debug  # Uses debug.sh script
```

## Extension Points

### Adding New Features

1. **New Helm Value**: Add to `helm_chart.go` values map
2. **New Resource**: Create dedicated file (e.g., `backup_config.go`)
3. **New Output**: Add to `locals.go` initialization and `outputs.go` constants
4. **New Validation**: Update `main.go` pre-flight checks

### Version Updates

Update `variables.go`:
```go
HelmChartVersion: "0.XX.Y"
```

Test with `spec.version` field before changing default.

## Known Limitations

1. **Cassandra Only Embedded DB**: PostgreSQL/MySQL embedded options exist in Helm chart but are not production-ready
2. **Single Namespace**: All Temporal components must reside in same namespace
3. **Istio Required for HTTP Ingress**: Gateway API resources assume Istio gateway controller
4. **external-dns Required for gRPC**: LoadBalancer DNS automation needs external-dns

## Comparison to Terraform Module

Both modules achieve the same outcome with different tooling:

| Aspect | Pulumi (Go) | Terraform (HCL) |
|--------|-------------|-----------------|
| Structure | Multi-file Go package | HCL files (variables, locals, main, outputs) |
| Conditionals | Native Go `if` statements | `dynamic` blocks, `count` |
| Loops | Go loops | `for_each`, `dynamic` |
| Values | Go maps/structs | HCL maps/objects |
| Dependencies | Explicit resource passing | `depends_on` |
| Type Safety | Compile-time | Runtime (plan-time) |

**Shared Concepts**:
- Locals for derived values
- Outputs for stack exports
- Helm provider for chart deployment
- Kubernetes provider for raw resources

## Maintenance Guidelines

1. **Keep Files Focused**: Each resource file should manage one logical resource group
2. **Mirror Pulumi README**: Keep user-facing docs in sync with implementation
3. **Test Backend Changes**: All three database backends must work
4. **Version Pin Carefully**: Only update default Helm version after testing
5. **Preserve Output Keys**: Changing output key names breaks downstream consumers

## References

- Temporal Helm Chart: https://github.com/temporalio/helm-charts
- OpenMCF Architecture: `/architecture/deployment-component.md`
- Spec Definition: `../spec.proto`
- Stack Outputs: `../stack_outputs.proto`

