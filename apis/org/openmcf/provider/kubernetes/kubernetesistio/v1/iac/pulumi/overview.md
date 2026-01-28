# Kubernetes Istio - Pulumi Module Overview

## Architecture

The Kubernetes Istio Pulumi module deploys a complete Istio service mesh using the official Helm charts. The module orchestrates the installation of three core components in the correct dependency order: base (CRDs), istiod (control plane), and ingress gateway.

## Architecture Diagram

```
┌────────────────────────────────────────────────────────────────┐
│                     Kubernetes Cluster                          │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │            Namespace: istio-system                        │ │
│  │                                                           │ │
│  │  ┌──────────────────────────────────────────────────┐    │ │
│  │  │           Istio Base (CRDs)                      │    │ │
│  │  │                                                  │    │ │
│  │  │  • VirtualService CRD                            │    │ │
│  │  │  • DestinationRule CRD                           │    │ │
│  │  │  • Gateway CRD                                   │    │ │
│  │  │  • ServiceEntry CRD                              │    │ │
│  │  └──────────────────────────────────────────────────┘    │ │
│  │                          │                               │ │
│  │  ┌──────────────────────▼───────────────────────────┐    │ │
│  │  │        Istiod (Control Plane)                    │    │ │
│  │  │                                                  │    │ │
│  │  │  ┌─────────┐  ┌─────────┐  ┌─────────┐         │    │ │
│  │  │  │ Pilot   │  │ Citadel │  │ Galley  │         │    │ │
│  │  │  │(Traffic)│  │(Security)│ │(Config) │         │    │ │
│  │  │  └─────────┘  └─────────┘  └─────────┘         │    │ │
│  │  │                                                  │    │ │
│  │  │  - Traffic Management                           │    │ │
│  │  │  - Certificate Management                       │    │ │
│  │  │  - Configuration Validation                     │    │ │
│  │  │  - Sidecar Injection                            │    │ │
│  │  └──────────────────────────────────────────────────┘    │ │
│  └───────────────────────────────────────────────────────────┘ │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │            Namespace: istio-ingress                       │ │
│  │                                                           │ │
│  │  ┌──────────────────────────────────────────────────┐    │ │
│  │  │         Istio Ingress Gateway                    │    │ │
│  │  │                                                  │    │ │
│  │  │  ┌──────────┐  ┌──────────┐                     │    │ │
│  │  │  │   Pod    │  │   Pod    │                     │    │ │
│  │  │  │ (Envoy)  │  │ (Envoy)  │                     │    │ │
│  │  │  └──────────┘  └──────────┘                     │    │ │
│  │  │                                                  │    │ │
│  │  │  - External Traffic Entry Point                 │    │ │
│  │  │  - TLS Termination                              │    │ │
│  │  │  - Load Balancing                               │    │ │
│  │  └──────────────────────────────────────────────────┘    │ │
│  └───────────────────────────────────────────────────────────┘ │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │         Application Namespaces (with sidecar injection)   │ │
│  │                                                           │ │
│  │  ┌─────────────────────────┐  ┌─────────────────────┐    │ │
│  │  │     Service A           │  │     Service B       │    │ │
│  │  │  ┌────────┬──────────┐  │  │  ┌────────┬──────┐  │    │ │
│  │  │  │  App   │  Envoy   │  │  │  │  App   │Envoy │  │    │ │
│  │  │  │  Pod   │  Sidecar │  │  │  │  Pod   │Sidecar│  │    │ │
│  │  │  └────────┴──────────┘  │  │  └────────┴──────┘  │    │ │
│  │  └─────────────────────────┘  └─────────────────────┘    │ │
│  └───────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

## Component Flow

### 1. Initialization

- Receive KubernetesIstioStackInput with target cluster and spec
- Set up Kubernetes provider from cluster credentials
- Extract container resource specifications

### 2. Namespace Creation

**istio-system namespace:**
- Hosts Istio control plane (istiod)
- Contains base resources and CRDs
- Apply resource labels for organization/environment tracking

**istio-ingress namespace:**
- Hosts ingress gateway pods
- Separate from control plane for security isolation
- Enables independent scaling and management

### 3. Component Installation Sequence

The module ensures proper dependency ordering:

**Step 1: Istio Base**
```
helm_release "istio-base"
  ├─ Install CRDs
  ├─ Create base resources
  └─ Wait for completion (timeout: 180s)
```

**Step 2: Istiod Control Plane**
```
helm_release "istiod"
  ├─ Deploy control plane
  ├─ Configure pilot resources from spec
  ├─ Enable atomic rollback
  └─ Depends on: istio-base
```

**Step 3: Ingress Gateway**
```
helm_release "istio-gateway"
  ├─ Deploy ingress gateway
  ├─ Configure as ClusterIP service
  ├─ Enable atomic rollback
  └─ Depends on: istiod
```

### 4. Resource Configuration

The module maps spec container resources to Helm values:

```go
pilot.resources.requests.cpu    = spec.container.resources.requests.cpu
pilot.resources.requests.memory = spec.container.resources.requests.memory
pilot.resources.limits.cpu      = spec.container.resources.limits.cpu
pilot.resources.limits.memory   = spec.container.resources.limits.memory
```

This allows users to control istiod resource consumption without writing Helm values.

### 5. Stack Outputs

Export deployment details:
- `namespace`: istio-system
- `service`: istiod
- `port_forward_command`: kubectl command for local access
- `kube_endpoint`: Internal cluster endpoint
- `ingress_endpoint`: Gateway service endpoint

## Module Design

### File Organization

- **`main.go`**: Orchestrates deployment sequence
  - Provider setup
  - Namespace creation (system and ingress)
  - Helm chart deployments with dependencies
  - Resource configuration mapping
  - Output exports

- **`vars.go`**: Deployment constants
  - Namespace names (istio-system, istio-ingress)
  - Helm repository URL
  - Chart names (base, istiod, gateway)
  - Pinned chart versions

- **`outputs.go`**: Defines output constant names
  - OpNamespace
  - OpService
  - OpPortForwardCommand
  - OpKubeEndpoint
  - OpIngressEndpoint

### Dependency Management

The module uses Pulumi's dependency tracking:

```go
helm.NewRelease(ctx, "istio-base", ...)
  // No dependencies

helm.NewRelease(ctx, "istiod", ..., 
  pulumi.DependsOn([]pulumi.Resource{baseRelease}))
  // Waits for base completion

helm.NewRelease(ctx, "istio-gateway", ..., 
  pulumi.DependsOn([]pulumi.Resource{istiodRelease}))
  // Waits for istiod completion
```

This ensures components are installed in the correct order.

## Data Flow

### Configuration Propagation Flow

1. **User** → Defines KubernetesIstio resource with container resources
2. **Pulumi Module** → Maps resources to Helm values
3. **Helm Charts** → Deploy components with configured resources
4. **Istiod** → Starts with specified CPU/memory
5. **Sidecars** → Connect to control plane for configuration

### Service Mesh Traffic Flow

1. **External Request** → Istio Ingress Gateway (istio-ingress namespace)
2. **Gateway** → Routes to VirtualService based on rules
3. **VirtualService** → Directs to DestinationRule
4. **DestinationRule** → Selects backend pods with Envoy sidecars
5. **Envoy Sidecar** → Enforces policies (retries, timeouts, circuit breakers)
6. **Application Container** → Processes request
7. **Response** → Flows back through sidecars (with mTLS encryption)

### Certificate Distribution Flow

1. **Istiod (Citadel)** → Generates root CA certificate
2. **Sidecars** → Request workload certificates via CSR
3. **Istiod** → Signs certificates and returns to sidecars
4. **Sidecars** → Use certificates for mTLS connections
5. **Istiod** → Rotates certificates automatically (default: every 24 hours)

## Design Decisions

### Why Three Separate Helm Charts?

**Modularity:**
- Base CRDs can be upgraded independently
- Control plane and gateways have different upgrade schedules
- Enables granular rollback if one component fails

**Flexibility:**
- Can deploy multiple gateways (ingress/egress) separately
- Different resource requirements for each component
- Separate lifecycle management

### Why Separate Namespaces?

**Security:**
- Control plane isolated from data plane
- Gateway pods run with minimal privileges
- Namespace-level RBAC enforcement

**Operational:**
- Independent resource quotas
- Separate network policies
- Easier troubleshooting and monitoring

### Why Atomic Releases?

**Reliability:**
- Automatic rollback on deployment failure
- No partial installations
- Consistent state after deployment

**Safety:**
- Production deployments remain stable
- Failed upgrades don't break existing mesh
- Easier disaster recovery

### Why ClusterIP for Gateway?

**Flexibility:**
- Users can choose exposure method (LoadBalancer, NodePort, Ingress)
- Doesn't force specific cloud load balancer configuration
- More portable across environments

**Cost:**
- Avoids automatic cloud load balancer provisioning
- Users opt-in to external access
- Suitable for internal-only meshes

## Resource Sizing Guidance

### Development Environment
```yaml
container:
  resources:
    requests: { cpu: 25m, memory: 64Mi }
    limits:   { cpu: 500m, memory: 256Mi }
```
- 5-10 services
- < 100 RPS
- Single developer

### Standard Production
```yaml
container:
  resources:
    requests: { cpu: 500m, memory: 512Mi }
    limits:   { cpu: 2000m, memory: 2Gi }
```
- 20-100 services
- 1,000-5,000 RPS
- Multiple teams

### High Availability
```yaml
container:
  resources:
    requests: { cpu: 1000m, memory: 1Gi }
    limits:   { cpu: 4000m, memory: 8Gi }
```
- 100-500 services
- 10,000+ RPS
- Mission-critical workloads

### Enterprise Scale
```yaml
container:
  resources:
    requests: { cpu: 2000m, memory: 4Gi }
    limits:   { cpu: 8000m, memory: 16Gi }
```
- 500+ services
- 50,000+ RPS
- Multi-tenant, multi-cluster

## Helm Chart Configuration

### Base Chart
- **Purpose**: Install CRDs and foundational resources
- **Timeout**: 180 seconds
- **Atomic**: Yes (rollback on failure)
- **CleanupOnFail**: Yes
- **WaitForJobs**: Yes

### Istiod Chart
- **Purpose**: Deploy control plane
- **Timeout**: 180 seconds
- **Atomic**: Yes
- **Resource Mapping**: From spec.container.resources
- **Values Set**:
  - `pilot.resources.requests.cpu`
  - `pilot.resources.requests.memory`
  - `pilot.resources.limits.cpu`
  - `pilot.resources.limits.memory`

### Gateway Chart
- **Purpose**: Deploy ingress gateway
- **Timeout**: 180 seconds
- **Atomic**: Yes
- **Values Set**:
  - `service.type: ClusterIP`

## Error Handling

### Namespace Creation Failure
- Module returns error immediately
- No Helm charts are deployed
- User can fix cluster permissions and retry

### Base Chart Failure
- Atomic rollback removes partially installed resources
- Istiod and gateway are not attempted
- CRDs remain if manually created

### Istiod Failure
- Atomic rollback removes control plane
- Gateway deployment is not attempted
- Base CRDs remain (safe for retry)

### Gateway Failure
- Atomic rollback removes gateway
- Control plane remains functional
- Mesh works for east-west traffic (no ingress)

## Monitoring and Observability

### Control Plane Metrics

Istiod exposes metrics on port 15014:
- Configuration push latency
- Proxy connection count
- Certificate rotation status
- Memory and CPU usage

Access via port-forward:
```bash
kubectl port-forward -n istio-system svc/istiod 15014:15014
curl http://localhost:15014/metrics
```

### Gateway Metrics

Gateway pods expose Envoy metrics:
- Request rate
- Response codes
- Latency percentiles
- Active connections

### Stack Outputs for Debugging

The module exports utility commands:
```bash
# Get port-forward command
pulumi stack output port_forward_command

# Execute to access debug interface
kubectl port-forward -n istio-system svc/istiod 15014:15014
```

## Upgrade Strategy

### In-Place Upgrade
1. Update chart version in vars.go
2. Run `pulumi up`
3. Atomic releases handle rollback if needed

### Blue-Green Upgrade
1. Deploy new Istio instance to different namespaces
2. Migrate workloads gradually
3. Delete old instance when complete

### Canary Upgrade
1. Deploy new istiod alongside old
2. Configure canary sidecars to use new control plane
3. Monitor and rollout gradually

## Best Practices

### Resource Allocation
- Start with default values
- Monitor actual usage with `kubectl top`
- Scale up based on observed metrics
- Don't over-provision (wastes resources)

### Version Pinning
- Use specific chart versions in production
- Test upgrades in staging first
- Document version compatibility

### Namespace Strategy
- Keep control plane in istio-system
- Keep gateways in dedicated namespaces
- Enable sidecar injection per application namespace

### High Availability
- Run multiple istiod replicas (via Helm values)
- Spread gateway pods across nodes/zones
- Set appropriate resource requests for QoS

## Troubleshooting Workflows

### Control Plane Issues
1. Check pod status: `kubectl get pods -n istio-system`
2. View logs: `kubectl logs -n istio-system -l app=istiod`
3. Check resources: `kubectl top pods -n istio-system`
4. Verify Helm release: `helm status istiod -n istio-system`

### Gateway Issues
1. Check pod status: `kubectl get pods -n istio-ingress`
2. View logs: `kubectl logs -n istio-ingress -l app=istio-gateway`
3. Describe service: `kubectl describe svc -n istio-ingress`

### Sidecar Injection Issues
1. Verify namespace label: `kubectl get ns -L istio-injection`
2. Check webhook: `kubectl get mutatingwebhookconfiguration`
3. Describe pod: `kubectl describe pod <pod-name>`

## Additional Resources

- [Istio Installation Guide](https://istio.io/latest/docs/setup/install/helm/)
- [Istio Performance and Scalability](https://istio.io/latest/docs/ops/deployment/performance-and-scalability/)
- [Helm Chart Documentation](https://github.com/istio/istio/tree/master/manifests/charts)
- [Component README](README.md) - Usage examples and quick start

