# Civo Kubernetes Cluster - Pulumi Module Architecture

## Overview

This document provides an architectural overview of the Pulumi module for deploying managed Kubernetes clusters on Civo Cloud. It explains design decisions, implementation patterns, and how the module integrates with the broader OpenMCF ecosystem.

## Architecture Principles

### 1. Single Entry Point Pattern

The module exposes a single public function:

```go
func Resources(
    ctx *pulumi.Context,
    stackInput *civokubernetesclusterv1.CivoKubernetesClusterStackInput,
) error
```

**Rationale:**
- Simplifies invocation by OpenMCF CLI
- Ensures consistent interface across all deployment components
- Hides internal complexity from callers
- Enables version upgrades without breaking API contracts

### 2. Protobuf-First Configuration

Input is defined as a protobuf message (`CivoKubernetesClusterStackInput`) rather than JSON or YAML structures.

**Benefits:**
- **Type safety**: Compile-time validation of input structure
- **Schema evolution**: Protobuf supports backward-compatible changes
- **Cross-language**: Same schema used by CLI (Go), API server (Go), and potential future clients
- **Validation**: buf.validate annotations ensure data integrity before reaching Pulumi

### 3. Declarative Resource Model

The module implements Kubernetes-style declarative semantics:

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoKubernetesCluster
metadata:
  name: prod-cluster
spec:
  clusterName: prod-k8s
  region: fra1
  kubernetesVersion: "1.29.0+k3s1"
  network:
    value: "network-id"
  defaultNodePool:
    size: g4s.kube.large
    nodeCount: 5
```

**Rationale:**
- Familiar to Kubernetes users
- Enables GitOps workflows
- Supports multi-cloud consistency
- Allows future CRD support for Kubernetes operators

### 4. Minimal 80/20 API Surface

The protobuf spec captures only essential cluster configuration:

**Required fields:**
- `cluster_name` - Cluster identifier
- `region` - Deployment location
- `kubernetes_version` - K8s version (explicit, no "latest")
- `network` - VPC for cluster networking
- `default_node_pool.size` - Node instance type
- `default_node_pool.node_count` - Number of nodes

**Optional fields:**
- `highly_available` - HA control plane
- `auto_upgrade` - Automatic patch upgrades
- `disable_surge_upgrade` - Upgrade strategy
- `tags` - Resource organization

**Intentionally omitted (advanced use cases):**
- Multiple node pools (use direct Civo API)
- CNI selection (requires platform decision)
- Marketplace applications (complex dependencies)
- Custom firewall rules (use `CivoFirewall` resource)

**Rationale**: 80% of users need simple cluster provisioning. Advanced features can be accessed via direct IaC or Civo CLI.

## Component Breakdown

### `module/main.go` (Core Logic)

**Purpose**: Orchestrate resource creation.

```go
func Resources(ctx *pulumi.Context, stackInput *CivoKubernetesClusterStackInput) error {
    locals := initializeLocals(ctx, stackInput)
    civoProvider, err := pulumicivoprovider.Get(ctx, stackInput.ProviderConfig)
    if err != nil {
        return errors.Wrap(err, "failed to setup Civo provider")
    }
    
    if _, err := cluster(ctx, locals, civoProvider); err != nil {
        return errors.Wrap(err, "failed to create Kubernetes cluster")
    }
    
    return nil
}
```

**Flow:**
1. Initialize locals
2. Set up Civo provider
3. Create Kubernetes cluster
4. Return any errors

### `module/locals.go` (Context Initialization)

**Purpose**: Extract and organize input data.

```go
type Locals struct {
    CivoProviderConfig    *civoprovider.CivoProviderConfig
    CivoKubernetesCluster *civokubernetesclusterv1.CivoKubernetesCluster
}
```

**Design note**: Unlike other modules, cluster labels are not explicitly set here. Civo manages cluster metadata through its own tagging system.

### `module/cluster.go` (Resource Provisioning)

**Purpose**: Create Civo Kubernetes cluster.

```go
createdCluster, err := civo.NewKubernetesCluster(
    ctx,
    "cluster",
    &civo.KubernetesClusterArgs{
        Name:              pulumi.String(spec.ClusterName),
        Region:            pulumi.String(spec.Region.String()),
        KubernetesVersion: pulumi.String(spec.KubernetesVersion),
        NetworkId:         pulumi.String(spec.Network.GetValue()),
        // ... default node pool configuration
    },
    pulumi.Provider(civoProvider),
)
```

**Key points:**
- Resource name is hardcoded as `"cluster"` for stability
- Region converted from enum to string
- Network ID extracted from StringValueOrRef
- Explicit provider ensures correct credentials

### `module/outputs.go` (Stack Outputs)

**Purpose**: Define constants for output keys.

```go
const (
    OpClusterId          = "cluster_id"
    OpKubeconfig         = "kubeconfig_b64"
    OpApiServerEndpoint  = "api_server_endpoint"
)
```

**Exported outputs:**
- `cluster_id` - Civo cluster UUID
- `kubeconfig_b64` - Base64-encoded kubeconfig
- `api_server_endpoint` - Kubernetes API URL

## Design Decisions

### 1. Why No Multiple Node Pools?

**Decision**: Support only default node pool in initial API version.

**Rationale:**
- 80% of clusters use single homogeneous node pool
- Multiple pools add significant API complexity
- Users needing mixed node types can use direct Civo API
- Keeps protobuf spec simple and approachable

**Future consideration**: Add optional `additional_node_pools` array in v2.

### 2. CNI Not Exposed in API

**Decision**: CNI plugin (Flannel vs Cilium) not configurable via protobuf spec.

**Rationale:**
- CNI is platform-level decision, not per-cluster
- Cannot be changed after cluster creation (breaking change)
- Requires understanding of NetworkPolicy implications
- Most users don't need to think about CNI

**Current approach**: CNI configured directly in Pulumi code or via Civo CLI flag during creation.

**Future consideration**: Add `cni: enum(flannel, cilium)` with strong warnings in documentation.

### 3. Marketplace Applications Omitted

**Decision**: Don't expose Civo marketplace applications in protobuf spec.

**Rationale:**
- 80+ marketplace apps create large API surface
- App installation can fail, complicating cluster provisioning
- Apps better installed post-deployment via Helm/kubectl
- Reduces blast radius (cluster creation independent of app issues)

**Alternative**: Users install apps after cluster is ready:

```bash
kubectl apply -f cert-manager.yaml
helm install prometheus prometheus-community/kube-prometheus-stack
```

### 4. Version Pinning Required

**Decision**: `kubernetes_version` is required field with no "latest" default.

**Rationale:**
- Explicit versions prevent surprise upgrades
- Enables reproducible deployments
- Forces users to consciously choose K8s version
- Aligns with production best practice (no implicit defaults)

**Trade-off**: Users must check available versions (`civo kubernetes versions`) before deploying.

### 5. StringValueOrRef for Network

**Decision**: Use `org.openmcf.shared.foreignkey.v1.StringValueOrRef` for network field.

**Structure:**
```protobuf
message StringValueOrRef {
    oneof literal_or_ref {
        string value = 1;
        ValueFromRef value_from = 2;
    }
}
```

**Rationale:**
- Supports literal values: `{value: "network-uuid"}`
- Supports references: `{value_from: {kind: "CivoVpc", name: "..."}}`
- Enables cross-resource dependencies
- Future-proof for orchestration

**Current limitation**: Reference resolution not yet implemented. Only literal `value` works.

### 6. Kubeconfig in Outputs

**Decision**: Export base64-encoded kubeconfig as stack output.

**Security considerations:**
- Kubeconfig contains cluster credentials
- Must be stored securely in Pulumi backend
- Use `--show-secrets` flag to view
- Never commit kubeconfig to Git

**Alternative considered**: Don't export kubeconfig, require users to fetch via Civo API. **Rejected** because:
- Less convenient workflow
- OpenMCF manages credentials centrally
- Pulumi backends are designed for secrets

## Integration Points

### OpenMCF CLI

**Invocation flow:**
1. User runs `planton apply -f cluster.yaml`
2. CLI parses YAML → protobuf `CivoKubernetesCluster`
3. CLI validates via buf.validate
4. CLI constructs `CivoKubernetesClusterStackInput`
5. CLI invokes `module.Resources(ctx, stackInput)`
6. CLI captures outputs and stores kubeconfig securely

**Kubeconfig access:**

```bash
planton kubeconfig civokubernetesclusters/dev-cluster > ~/.kube/config
```

### Standalone Pulumi Usage

**Invocation flow:**
1. User creates `Pulumi.yaml` and `stack-input.json`
2. User runs `pulumi up`
3. `main.go` deserializes stack input
4. Calls `module.Resources(ctx, stackInput)`
5. Pulumi manages state in configured backend

**Kubeconfig access:**

```bash
pulumi stack output kubeconfig_b64 --show-secrets | base64 -d > ~/.kube/config
```

### Civo Provider Setup

Uses shared helper:

```go
import "github.com/plantonhq/openmcf/pkg/iac/pulumi/pulumimodule/provider/civo/pulumicivoprovider"

civoProvider, err := pulumicivoprovider.Get(ctx, stackInput.ProviderConfig)
```

**Benefits:**
- Consistency across all Civo modules
- Centralized credential handling
- Single place for provider enhancements

## State Management

### Pulumi State Structure

For a 3-node cluster:

```
Cluster: cikc-prod-123
└── civo:index/kubernetesCluster:KubernetesCluster
    └── cluster
        ├── Name: "prod-k8s"
        ├── Region: "fra1"
        ├── NumTargetNodes: 3
        └── Pools: [default-pool]
```

**Resource naming:**
- Cluster: `"cluster"` (constant)

**Why stable names matter**: Pulumi tracks resources by URN. Changing resource name causes replacement (delete + create), losing all cluster data.

### State Drift Scenarios

#### Scenario 1: Manual node scaling

User scales nodes via Civo dashboard.

**Result**: Pulumi shows `node_count: 3` in spec but cluster has 5 nodes.

**Resolution**: Run `pulumi refresh` to sync state. On next `pulumi up`, Pulumi will scale back to 3 (spec is source of truth).

#### Scenario 2: Kubernetes version upgrade externally

User upgrades K8s version via Civo CLI.

**Result**: State drift. Pulumi shows old version in state.

**Resolution**: Update `kubernetes_version` in spec to match actual cluster version, then run `pulumi up` (no-op).

#### Scenario 3: Cluster deleted manually

User deletes cluster via Civo dashboard.

**Result**: Pulumi state references non-existent cluster. On next `pulumi up`, Pulumi will recreate it.

**Warning**: Recreating cluster creates new UUID. Any external references (DNS, LoadBalancers) will break.

## Performance Characteristics

### Cluster Creation Time

- **API call to Civo**: ~2-3 seconds
- **Cluster provisioning**: 90-120 seconds
- **Total**: ~2 minutes for ready cluster

**Breakdown:**
1. Network setup: 5-10 seconds
2. Control plane creation: 30-40 seconds
3. Worker nodes provisioning: 40-60 seconds
4. K3s initialization: 10-20 seconds

### Scaling Operations

- **Add node**: ~60 seconds per node
- **Remove node**: ~30 seconds (drain + delete)
- **Change node size**: Replace operation (delete + create)

### Upgrade Operations

- **Patch upgrade** (1.29.0 → 1.29.1): ~5-10 minutes
- **Minor upgrade** (1.29 → 1.30): ~10-15 minutes

**Note**: Civo handles rolling upgrades. Downtime depends on workload configuration (PodDisruptionBudgets).

### Destroy Operations

- **Delete cluster**: ~60-90 seconds
- **Network cleanup**: Automatic (if using OpenMCF)

**Warning**: Destroying cluster also deletes:
- All pods, deployments, services
- Persistent volumes (if using Civo volumes)
- Load balancers
- Associated IPs (if not reserved)

## Security Considerations

### Kubeconfig Protection

**Risk**: Kubeconfig grants full cluster admin access.

**Mitigation**:
- Stored encrypted in Pulumi state backend
- Access controlled via Pulumi/AWS IAM
- Never logged or exposed in plain text
- Rotate credentials regularly

**Best practice**: Use RBAC to create limited-access service accounts instead of admin kubeconfig.

### API Server Exposure

**Default**: API server (port 6443) publicly accessible.

**Recommendation**: Restrict via firewall:

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoFirewall
spec:
  inboundRules:
    - protocol: tcp
      portRange: "6443"
      cidrs: ["203.0.113.0/24"]  # Office only
      label: K8s API restricted
```

### Cluster Updates

**Risk**: Kubernetes upgrades can introduce breaking changes.

**Mitigation**:
- Test upgrades in dev/staging first
- Review K8s release notes
- Use `pulumi preview` to see changes
- Take cluster backups before upgrades

## Testing Strategy

### Unit Tests

**Location**: `module/*_test.go` (not yet implemented)

**Coverage:**
- `initializeLocals`: Verify field extraction
- Resource name generation
- Tag conversion logic

### Integration Tests

**Location**: `iac/pulumi/integration_test.go` (not yet implemented)

**Approach:**
1. Create test Civo account
2. Deploy minimal cluster
3. Verify cluster is ready
4. Test kubectl access
5. Destroy cluster
6. Verify cleanup

**Challenges**: Requires live Civo credentials and ~2 minutes per test.

### Validation Tests

**Location**: `v1/spec_test.go` ✅ (implemented - 25 tests)

**Coverage:**
- Required fields validation
- Node count validation (must be > 0)
- Region enum validation
- API version and kind validation

**Benefit**: Catches invalid input before expensive cluster provisioning.

## Troubleshooting Guide

### Issue: "Cluster name already exists"

**Symptom**: Pulumi error during cluster creation.

**Cause**: Cluster with same name exists (possibly in another account/region).

**Resolution:**
1. Check: `civo kubernetes list`
2. If cluster exists, import it: `pulumi import civo:index/kubernetesCluster:KubernetesCluster cluster <cluster-id>`
3. Or use different name

### Issue: Cluster stuck in "Building"

**Symptom**: `pulumi up` shows "still creating" for > 5 minutes.

**Cause**: Civo provisioning delay or capacity issue.

**Resolution:**
1. Check Civo dashboard: https://dashboard.civo.com/kubernetes
2. Verify network has available IPs
3. Check region capacity
4. Wait up to 10 minutes (Civo SLA)

### Issue: kubectl connection refused

**Symptom**: `kubectl get nodes` fails after deployment.

**Cause**: API server not reachable or firewall blocking.

**Resolution:**
1. Test API endpoint: `curl -k $(pulumi stack output api_server_endpoint)`
2. Check firewall rules allow port 6443
3. Verify kubeconfig is correct: `kubectl config view`

### Issue: Nodes not ready

**Symptom**: `kubectl get nodes` shows NotReady status.

**Cause**: K3s agent startup failure or network issues.

**Resolution:**
1. Describe node: `kubectl describe node <name>`
2. SSH to node: `civo kubernetes show <cluster>` → get node IP → `ssh civo@<ip>`
3. Check K3s logs: `sudo journalctl -u k3s-agent -f`

## Future Enhancements

### Planned

1. **Multiple node pools** - Support heterogeneous node configurations
2. **CNI selection** - Expose CNI choice in protobuf (Flannel vs Cilium)
3. **Marketplace apps** - Optional application installation during cluster creation
4. **Autoscaler integration** - Built-in cluster autoscaler configuration
5. **Custom firewall** - Reference `CivoFirewall` directly in spec

### Under Consideration

1. **GPU node support** - When Civo adds GPU instances
2. **Spot instances** - Cost optimization for non-critical workloads
3. **Private clusters** - API server not publicly accessible
4. **Cluster templates** - Pre-defined configurations (dev, staging, prod)

## Comparison: Civo vs Other Providers

### K3s vs Standard Kubernetes

Civo uses K3s (lightweight Kubernetes):

**Benefits:**
- 50% less memory footprint
- Faster startup (30-40 seconds vs 2-3 minutes)
- Embedded SQLite (no separate etcd for single-master)
- Perfect for edge and resource-constrained scenarios

**Trade-offs:**
- containerd only (no Docker runtime)
- Some legacy features removed
- Smaller ecosystem (fewer K3s-specific tools)

**Verdict**: For 95% of workloads, K3s is indistinguishable from standard K8s. It's CNCF-certified and API-compliant.

### Civo vs EKS/GKE/AKS

| Feature | Civo K3s | AWS EKS | GCP GKE | Azure AKS |
|---------|----------|---------|---------|-----------|
| **Startup time** | 90-120s | 10-15 min | 5-10 min | 5-10 min |
| **Min cost** | ~$11/mo | ~$73/mo | ~$75/mo | ~$0 (free control plane) |
| **Control plane HA** | Optional | Built-in | Built-in | Built-in |
| **Node autoscaling** | Add-on | Native | Native | Native |
| **Best for** | Cost-conscious, fast iteration | Enterprise AWS | GCP-native | Azure-native |

**Civo's niche**: Development clusters, startups, cost-sensitive workloads, rapid experimentation.

## References

- [Pulumi Civo Provider](https://www.pulumi.com/registry/packages/civo/)
- [Civo Kubernetes API](https://www.civo.com/api/kubernetes)
- [K3s Documentation](https://docs.k3s.io/)
- [User Documentation](../../README.md)

## Changelog

- **2025-11-16**: Initial architecture documentation
- **2025-11-14**: Implementation completed
- **2025-11-10**: Module scaffolding created

---

**Maintained by**: OpenMCF Team  
**Last Updated**: 2025-11-16

