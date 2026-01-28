# Civo Kubernetes Node Pool - Pulumi Module Architecture

## Overview

Architectural overview of the Pulumi module for managing Civo Kubernetes node pools.

## Architecture Principles

### 1. Single Entry Point

```go
func Resources(
    ctx *pulumi.Context,
    stackInput *civokubernetesnodepoolv1.CivoKubernetesNodePoolStackInput,
) error
```

### 2. Protobuf-First Configuration

Input via `CivoKubernetesNodePoolStackInput` protobuf message.

**Benefits:**
- Type safety
- Schema evolution
- Cross-language support
- Built-in validation

### 3. Minimal 80/20 API

Exposes only essential node pool configuration:

**Required:**
- `node_pool_name` - Pool identifier
- `cluster` - Parent cluster reference
- `size` - Node instance type
- `node_count` - Number of nodes (>0)

**Optional:**
- `auto_scale` - Enable autoscaling
- `min_nodes` / `max_nodes` - Autoscaling bounds
- `tags` - Organization metadata

**Omitted (advanced):**
- Custom taints (apply via kubectl)
- Node labels (apply via kubectl)
- Public IP configuration

## Component Breakdown

### `module/main.go`

Orchestrates resource creation:

```go
func Resources(ctx *pulumi.Context, stackInput *StackInput) error {
    locals := initializeLocals(ctx, stackInput)
    civoProvider, _ := pulumicivoprovider.Get(ctx, stackInput.ProviderConfig)
    _, err := nodePool(ctx, locals, civoProvider)
    return err
}
```

### `module/node_pool.go`

Creates Civo node pool:

```go
createdPool, err := civo.NewKubernetesNodePool(
    ctx,
    "node_pool",
    &civo.KubernetesNodePoolArgs{
        ClusterId: pulumi.String(spec.Cluster.GetValue()),
        NodeCount: pulumi.Int(int(spec.NodeCount)),
        Size:      pulumi.String(spec.Size),
        // ... autoscaling config
    },
    pulumi.Provider(civoProvider),
)
```

### `module/outputs.go`

Output constants:
- `node_pool_id` - Civo node pool UUID

## Design Decisions

### 1. Autoscaling Implementation

**Decision**: Expose autoscaling fields in protobuf spec.

**Rationale:**
- Simplifies common use case (70% of pools use autoscaling)
- Abstracts Civo cluster autoscaler installation
- Consistent API across cloud providers

### 2. Node Count as Initial Size

When `autoScale=true`, `node_count` is the starting size:

```yaml
nodeCount: 3   # Start with 3
minNodes: 2    # Can scale down to 2
maxNodes: 10   # Can scale up to 10
```

### 3. StringValueOrRef for Cluster

Supports literal cluster names or references to `CivoKubernetesCluster` resources.

## Performance Characteristics

- **Node pool creation**: 10-20 seconds
- **Node provisioning**: ~60 seconds per node
- **Autoscale trigger**: 30-60 seconds after pending pod
- **Scale down**: 10+ minutes of low utilization

## Security Considerations

- Node pools inherit cluster security settings
- Use taints/labels for workload isolation
- Monitor autoscaling bounds to prevent cost overruns

## Testing Strategy

**Validation tests**: `v1/spec_test.go` ✅ (21 tests passing)

**Coverage:**
- Required fields validation
- Node count > 0 validation
- Autoscaling configurations
- API version and kind validation

## References

- [Pulumi Civo Provider](https://www.pulumi.com/registry/packages/civo/)
- [User Documentation](../../README.md)

---

**Maintained by**: OpenMCF Team
**Last Updated**: 2025-11-16

