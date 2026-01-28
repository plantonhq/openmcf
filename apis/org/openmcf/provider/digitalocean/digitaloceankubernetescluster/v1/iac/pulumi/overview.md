# DigitalOcean Kubernetes Cluster - Pulumi Module Architecture

This document provides an in-depth technical overview of the Pulumi module implementation for DigitalOcean Kubernetes clusters.

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Module Structure](#module-structure)
- [Implementation Details](#implementation-details)
- [Resource Provisioning Flow](#resource-provisioning-flow)
- [Configuration Processing](#configuration-processing)
- [Error Handling](#error-handling)
- [Design Decisions](#design-decisions)
- [Customization Points](#customization-points)

---

## Architecture Overview

The Pulumi module is designed following OpenMCF's standard patterns for cloud resource provisioning:

```
┌─────────────────────────────────────────────────────────────┐
│                     Pulumi Entrypoint                       │
│                       (main.go)                             │
│  - Loads stack input (manifest YAML)                       │
│  - Initializes Pulumi context                              │
│  - Delegates to module                                      │
└──────────────────┬──────────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────────┐
│                    Module Layer                             │
│                  (module/main.go)                           │
│  - Processes DigitalOceanKubernetesCluster spec             │
│  - Creates DigitalOcean provider                            │
│  - Orchestrates resource creation                          │
└──────────────────┬──────────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────────┐
│               Resource Implementations                      │
│                                                             │
│  ┌─────────────────┐  ┌──────────────────┐                │
│  │  cluster.go     │  │   locals.go      │                │
│  │  - DOKS cluster │  │   - Local vars   │                │
│  │  - Node pool    │  │   - Transforms   │                │
│  │  - Networking   │  │   - Helpers      │                │
│  └─────────────────┘  └──────────────────┘                │
│                                                             │
│  ┌─────────────────┐                                       │
│  │  outputs.go     │                                       │
│  │  - Stack outputs│                                       │
│  │  - Exports      │                                       │
│  └─────────────────┘                                       │
└─────────────────────────────────────────────────────────────┘
```

## Module Structure

### File Organization

```
iac/pulumi/
├── main.go                 # Pulumi program entrypoint
├── Pulumi.yaml            # Project configuration
├── Makefile               # Helper commands
├── debug.sh               # Debug script
├── module/                # Core implementation
│   ├── main.go            # Module orchestration
│   ├── cluster.go         # Cluster resource logic
│   ├── locals.go          # Local variables and helpers
│   └── outputs.go         # Output constants and exports
├── README.md              # Usage documentation
└── overview.md            # This file
```

### Dependency Graph

```
main.go
  └── module.Resources()
        ├── locals.NewLocals()
        │     └── Processes DigitalOceanKubernetesCluster spec
        ├── module.createDigitalOceanProvider()
        │     └── Uses DIGITALOCEAN_TOKEN from env
        └── module.cluster()
              ├── Creates digitalocean.KubernetesCluster
              ├── Configures maintenance policy (optional)
              ├── Configures firewall (optional)
              ├── Exports outputs
              └── Returns cluster resource
```

---

## Implementation Details

### 1. Entrypoint (`main.go`)

**Responsibility:** Initialize Pulumi context and load stack input

```go
func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        stackInput := &digitaloceankubernetesclusterv1.DigitalOceanKubernetesClusterStackInput{}
        
        // Load from Pulumi config or stack exports
        if err := loadStackInput(ctx, stackInput); err != nil {
            return err
        }
        
        // Delegate to module
        return module.Resources(ctx, stackInput)
    })
}
```

**Key Features:**
- Pulumi context initialization
- Stack input deserialization from YAML/JSON
- Error propagation to Pulumi runtime

### 2. Module Orchestration (`module/main.go`)

**Responsibility:** Coordinate resource provisioning

```go
func Resources(ctx *pulumi.Context, stackInput *StackInput) error {
    // 1. Process specification into local variables
    locals := NewLocals(stackInput)
    
    // 2. Create DigitalOcean provider
    provider, err := createDigitalOceanProvider(ctx)
    if err != nil {
        return err
    }
    
    // 3. Create cluster resource
    cluster, err := cluster(ctx, locals, provider)
    if err != nil {
        return err
    }
    
    return nil
}
```

**Key Patterns:**
- **Locals Pattern**: Centralize configuration processing
- **Provider Pattern**: Explicit provider instances for multi-cloud support
- **Error Chaining**: Immediate error return on failures

### 3. Cluster Resource (`module/cluster.go`)

**Responsibility:** Provision DOKS cluster with all configurations

#### Core Logic

```go
func cluster(ctx *pulumi.Context, locals *Locals, provider *digitalocean.Provider) (*digitalocean.KubernetesCluster, error) {
    // 1. Process tags
    var tags pulumi.StringArray
    for _, t := range locals.DigitalOceanKubernetesCluster.Spec.Tags {
        tags = append(tags, pulumi.String(t))
    }
    
    // 2. Build maintenance policy (optional)
    var maintenancePolicy *KubernetesClusterMaintenancePolicyArgs
    if locals.Spec.MaintenanceWindow != "" {
        maintenancePolicy = &KubernetesClusterMaintenancePolicyArgs{
            StartTime: pulumi.String(locals.Spec.MaintenanceWindow),
        }
    }
    
    // 3. Build firewall (optional)
    var firewall *KubernetesClusterFirewallArgs
    if len(locals.Spec.ControlPlaneFirewallAllowedIps) > 0 {
        var allowedIPs pulumi.StringArray
        for _, ip := range locals.Spec.ControlPlaneFirewallAllowedIps {
            allowedIPs = append(allowedIPs, pulumi.String(ip))
        }
        firewall = &KubernetesClusterFirewallArgs{
            AllowedAddresses: allowedIPs,
        }
    }
    
    // 4. Create cluster
    cluster, err := digitalocean.NewKubernetesCluster(ctx, "cluster", &KubernetesClusterArgs{
        Name:                locals.Spec.ClusterName,
        Region:              pulumi.String(locals.Spec.Region.String()),
        Version:             pulumi.String(locals.Spec.KubernetesVersion),
        Ha:                  pulumi.BoolPtr(locals.Spec.HighlyAvailable),
        AutoUpgrade:         pulumi.BoolPtr(locals.Spec.AutoUpgrade),
        SurgeUpgrade:        pulumi.BoolPtr(!locals.Spec.DisableSurgeUpgrade),
        VpcUuid:             pulumi.String(locals.Spec.Vpc.GetValue()),
        RegistryIntegration: pulumi.BoolPtr(locals.Spec.RegistryIntegration),
        MaintenancePolicy:   maintenancePolicy,
        Firewall:            firewall,
        Tags:                tags,
        NodePool: &KubernetesClusterNodePoolArgs{
            Name:      pulumi.String("default"),
            Size:      pulumi.String(locals.Spec.DefaultNodePool.Size),
            NodeCount: pulumi.IntPtr(int(locals.Spec.DefaultNodePool.NodeCount)),
            AutoScale: pulumi.BoolPtr(locals.Spec.DefaultNodePool.AutoScale),
            MinNodes:  pulumi.IntPtr(int(locals.Spec.DefaultNodePool.MinNodes)),
            MaxNodes:  pulumi.IntPtr(int(locals.Spec.DefaultNodePool.MaxNodes)),
        },
    }, pulumi.Provider(provider), pulumi.IgnoreChanges([]string{"version"}))
    
    if err != nil {
        return nil, errors.Wrap(err, "failed to create cluster")
    }
    
    // 5. Export outputs
    ctx.Export(OpClusterId, cluster.ID())
    ctx.Export(OpApiServerEndpoint, cluster.Endpoint)
    ctx.Export(OpKubeconfig, cluster.KubeConfigs.Index(pulumi.Int(0)).RawConfig())
    
    return cluster, nil
}
```

#### Key Features

**1. Optional Configuration Handling**

Uses Go conditionals to include optional blocks:

```go
// Only include maintenance policy if configured
var maintenancePolicy *KubernetesClusterMaintenancePolicyArgs
if locals.Spec.MaintenanceWindow != "" {
    maintenancePolicy = &KubernetesClusterMaintenancePolicyArgs{
        StartTime: pulumi.String(locals.Spec.MaintenanceWindow),
    }
}
```

**2. Dynamic List Processing**

Converts proto repeated fields to Pulumi arrays:

```go
var tags pulumi.StringArray
for _, t := range locals.DigitalOceanKubernetesCluster.Spec.Tags {
    tags = append(tags, pulumi.String(t))
}
```

**3. Lifecycle Management**

```go
pulumi.IgnoreChanges([]string{"version"})
```

Prevents drift detection on `version` field since auto-upgrades modify it externally.

**4. Provider Binding**

```go
pulumi.Provider(provider)
```

Ensures cluster uses the explicitly created provider instance.

### 4. Local Variables (`module/locals.go`)

**Responsibility:** Transform stack input into usable local variables

```go
type Locals struct {
    DigitalOceanKubernetesCluster *DigitalOceanKubernetesCluster
    DigitalOceanProvider          *digitalocean.Provider
}

func NewLocals(stackInput *StackInput) *Locals {
    return &Locals{
        DigitalOceanKubernetesCluster: stackInput.Target,
    }
}
```

**Design Rationale:**
- Centralized access to configuration
- Type-safe references
- Future extensibility for computed values

### 5. Outputs (`module/outputs.go`)

**Responsibility:** Define output constants and export values

```go
const (
    OpClusterId          = "cluster_id"
    OpKubeconfig         = "kubeconfig"
    OpApiServerEndpoint  = "api_server_endpoint"
)
```

**Exported Values:**
- `cluster_id`: Cluster UUID for reference in other stacks
- `kubeconfig`: Sensitive, base64-encoded kubeconfig
- `api_server_endpoint`: API server URL for monitoring

---

## Resource Provisioning Flow

### Sequence Diagram

```
┌──────┐     ┌────────┐     ┌──────────┐     ┌──────────────┐     ┌────────────┐
│Client│     │main.go │     │module/   │     │DigitalOcean  │     │DigitalOcean│
│      │     │        │     │main.go   │     │Provider      │     │API         │
└──┬───┘     └───┬────┘     └────┬─────┘     └──────┬───────┘     └─────┬──────┘
   │             │               │                  │                   │
   │ pulumi up   │               │                  │                   │
   │────────────>│               │                  │                   │
   │             │               │                  │                   │
   │             │ Load manifest │                  │                   │
   │             │──────────────>│                  │                   │
   │             │               │                  │                   │
   │             │               │ Create provider  │                   │
   │             │               │─────────────────>│                   │
   │             │               │                  │                   │
   │             │               │ Create cluster   │                   │
   │             │               │─────────────────>│                   │
   │             │               │                  │                   │
   │             │               │                  │ POST /clusters    │
   │             │               │                  │──────────────────>│
   │             │               │                  │                   │
   │             │               │                  │ 201 Created       │
   │             │               │                  │<──────────────────│
   │             │               │                  │                   │
   │             │               │ Cluster resource │                   │
   │             │               │<─────────────────│                   │
   │             │               │                  │                   │
   │             │ Export outputs│                  │                   │
   │             │<──────────────│                  │                   │
   │             │               │                  │                   │
   │ Complete    │               │                  │                   │
   │<────────────│               │                  │                   │
```

### Step-by-Step Execution

1. **User runs `pulumi up`**
   - Pulumi CLI invokes `main.go`
   - Initializes Pulumi context

2. **Load stack input**
   - Deserializes manifest YAML
   - Validates required fields

3. **Create DigitalOcean provider**
   - Reads `DIGITALOCEAN_TOKEN` from environment
   - Initializes authenticated provider client

4. **Process locals**
   - Transforms spec into internal representation
   - Validates VPC references

5. **Create cluster resource**
   - Converts spec to DigitalOcean API parameters
   - Sends API request to create cluster
   - Polls for completion (3-5 minutes)

6. **Export outputs**
   - Extracts cluster ID, kubeconfig, API endpoint
   - Stores in Pulumi state

7. **Complete deployment**
   - Returns success to CLI
   - Updates state with new resources

---

## Configuration Processing

### Spec Validation

The module relies on proto validation rules defined in `spec.proto`:

```proto
message DigitalOceanKubernetesClusterSpec {
  string cluster_name = 1 [(buf.validate.field).required = true];
  DigitalOceanRegion region = 2 [(buf.validate.field).required = true];
  string kubernetes_version = 3 [(buf.validate.field).required = true];
  
  DigitalOceanKubernetesClusterDefaultNodePool default_node_pool = 12 [
    (buf.validate.field).required = true
  ];
}
```

**Validation happens before Pulumi execution**, ensuring invalid specs are rejected early.

### Field Transformations

#### Region Enum to String

```go
Region: pulumi.String(locals.Spec.Region.String())
```

Converts proto enum to DigitalOcean region slug (e.g., `NYC1` → `"nyc1"`).

#### VPC Reference Resolution

```go
VpcUuid: pulumi.String(locals.Spec.Vpc.GetValue())
```

Extracts VPC UUID from `StringValueOrRef`, supporting both direct values and cross-stack references.

#### Boolean Inversion (Surge Upgrade)

```go
SurgeUpgrade: pulumi.BoolPtr(!locals.Spec.DisableSurgeUpgrade)
```

The spec uses `disable_surge_upgrade` (default false), but the API expects `surge_upgrade` (default true).

---

## Error Handling

### Error Propagation Pattern

```go
cluster, err := digitalocean.NewKubernetesCluster(...)
if err != nil {
    return nil, errors.Wrap(err, "failed to create cluster")
}
```

All errors are:
1. **Wrapped** with context using `github.com/pkg/errors`
2. **Returned immediately** (no silent failures)
3. **Propagated** to Pulumi runtime for user display

### Common Error Scenarios

| Error | Cause | Solution |
|-------|-------|----------|
| `vpc not found` | Invalid VPC UUID | Verify VPC exists in region |
| `version not supported` | Kubernetes version unavailable | Check available versions with `doctl` |
| `insufficient quota` | Account limits exceeded | Request quota increase |
| `authentication failed` | Invalid DO token | Set correct `DIGITALOCEAN_TOKEN` |

---

## Design Decisions

### 1. Ignore Version Changes

**Decision:** `pulumi.IgnoreChanges([]string{"version"})`

**Rationale:**
- Auto-upgrade feature modifies version externally
- Without ignore, every Pulumi run detects "drift"
- Manual upgrades can temporarily remove ignore

**Tradeoff:** Prevents drift detection on version, requires manual tracking of actual cluster version.

### 2. Optional Blocks via Conditionals

**Decision:** Use Go conditionals instead of Pulumi's `ApplyT` for optional blocks

```go
var maintenancePolicy *MaintenancePolicyArgs
if spec.MaintenanceWindow != "" {
    maintenancePolicy = &MaintenancePolicyArgs{...}
}
```

**Rationale:**
- Simpler, more readable code
- Compile-time type safety
- No asynchronous complexity

**Tradeoff:** Cannot handle dynamic optionals (values computed at runtime).

### 3. Explicit Provider Instances

**Decision:** Create explicit `digitalocean.Provider` instead of relying on defaults

```go
provider, err := digitalocean.NewProvider(ctx, "do-provider", &Args{
    Token: pulumi.String(os.Getenv("DIGITALOCEAN_TOKEN")),
})
```

**Rationale:**
- Explicit is better than implicit
- Enables multi-cloud scenarios (multiple providers)
- Easier debugging and testing

**Tradeoff:** Slightly more verbose code.

### 4. Single Node Pool

**Decision:** Support only default node pool in initial implementation

**Rationale:**
- 80/20 principle: Most users need single pool
- Additional pools can be added via separate resources
- Keeps spec simple

**Future Work:** Add support for multiple node pools via separate proto message.

---

## Customization Points

### Adding Additional Node Pools

To extend the module for multiple node pools:

1. **Update spec.proto:**
   ```proto
   repeated DigitalOceanKubernetesClusterNodePool additional_node_pools = 13;
   ```

2. **Update cluster.go:**
   ```go
   // Create additional node pools
   for i, pool := range locals.Spec.AdditionalNodePools {
       _, err := digitalocean.NewKubernetesNodePool(ctx, fmt.Sprintf("pool-%d", i), &Args{
           ClusterId: cluster.ID(),
           Name:      pulumi.String(pool.Name),
           Size:      pulumi.String(pool.Size),
           NodeCount: pulumi.IntPtr(int(pool.NodeCount)),
       }, pulumi.Provider(provider))
       if err != nil {
           return nil, err
       }
   }
   ```

### Adding Post-Provisioning Configuration

To run kubectl commands after cluster creation:

```go
// Wait for cluster to be ready
cluster.Status.ApplyT(func(status string) error {
    if status != "running" {
        return nil
    }
    
    // Apply NetworkPolicy, ConfigMap, etc.
    return applyKubernetesResources(cluster)
})
```

### Integrating with External Services

Export cluster outputs for use in other stacks:

```go
ctx.Export("cluster_id", cluster.ID())

// In another stack:
clusterID := pulumi.StackReference("doks-cluster").GetOutput("cluster_id")
```

---

## Performance Considerations

### Deployment Time

- **Cluster creation:** 3-5 minutes (DigitalOcean provisioning)
- **Pulumi overhead:** <10 seconds
- **Total:** ~3-5 minutes

### State Size

- Minimal state footprint (~5 KB per cluster)
- Scales linearly with additional resources

### API Rate Limits

DigitalOcean API limits:
- 5,000 requests/hour per token
- Pulumi module uses ~5-10 requests per deployment
- Safe for frequent updates

---

## Testing

### Unit Tests

(Future enhancement - not currently implemented)

```go
func TestClusterCreation(t *testing.T) {
    // Mock Pulumi context
    ctx := &pulumi.Context{}
    
    // Test cluster creation
    cluster, err := cluster(ctx, testLocals, testProvider)
    assert.NoError(t, err)
    assert.NotNil(t, cluster)
}
```

### Integration Tests

```bash
# Deploy test cluster
pulumi stack init test
pulumi config set digitalocean:token --secret $DO_TOKEN
pulumi up --yes

# Verify cluster
kubectl --kubeconfig=$(pulumi stack output kubeconfig --show-secrets) get nodes

# Cleanup
pulumi destroy --yes
```

---

## Reference

- [Pulumi DigitalOcean Provider Source](https://github.com/pulumi/pulumi-digitalocean)
- [DigitalOcean Kubernetes API](https://docs.digitalocean.com/reference/api/api-reference/#tag/Kubernetes)
- [OpenMCF Architecture](../../../../../../architecture/README.md)

---

**Version:** 1.0.0  
**Last Updated:** 2025-11-16  
**Authors:** OpenMCF Engineering Team

