# DigitalOcean Kubernetes Node Pool - Pulumi Module Architecture

This document provides a technical overview of the Pulumi module implementation for DigitalOcean Kubernetes node pools.

## Architecture Overview

```
┌─────────────────────────────────────────────────┐
│          Pulumi Entrypoint (main.go)            │
│  - Loads stack input                            │
│  - Initializes Pulumi context                   │
└──────────────────┬──────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────┐
│          Module Layer (module/main.go)           │
│  - Processes node pool spec                     │
│  - Creates DigitalOcean provider                │
│  - Orchestrates resource creation               │
└──────────────────┬──────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────┐
│       Resource Implementation                    │
│  ┌──────────────┐  ┌─────────────┐             │
│  │ node_pool.go │  │ locals.go   │             │
│  │ - Pool       │  │ - Variables │             │
│  │ - Labels     │  │ - Merging   │             │
│  │ - Taints     │  │             │             │
│  └──────────────┘  └─────────────┘             │
│                                                  │
│  ┌──────────────┐                               │
│  │ outputs.go   │                               │
│  │ - Exports    │                               │
│  └──────────────┘                               │
└─────────────────────────────────────────────────┘
```

## Module Structure

### File Organization

```
iac/pulumi/
├── main.go           # Pulumi entrypoint
├── module/
│   ├── main.go       # Module orchestration
│   ├── node_pool.go  # Node pool resource
│   ├── locals.go     # Local variables
│   └── outputs.go    # Output constants
```

## Implementation Details

### 1. Locals (locals.go)

**Responsibility:** Process metadata and generate standard labels

```go
type Locals struct {
    DigitalOceanKubernetesNodePool *NodePool
    DigitalOceanLabels             map[string]string
}
```

**Generated Labels:**
- `planton-resource`: "true"
- `planton-resource-name`: metadata.name
- `planton-resource-kind`: "DigitalOceanKubernetesNodePool"
- `planton-organization`: metadata.org (if set)
- `planton-environment`: metadata.env (if set)
- `planton-resource-id`: metadata.id (if set)

### 2. Node Pool Resource (node_pool.go)

**Responsibility:** Provision DOKS node pool with all configurations

#### Label Merging

```go
// Merge metadata labels with spec labels
labels := pulumi.StringMap{}
for k, v := range locals.DigitalOceanLabels {
    labels[k] = pulumi.String(v)
}
for k, v := range spec.Labels {
    labels[k] = pulumi.String(v)
}
```

**Merge Strategy:**
1. Add metadata-derived labels first
2. Add user-specified labels (can override metadata)

#### Taint Processing

```go
var taints KubernetesNodePoolTaintArray
for _, taint := range spec.Taints {
    taints = append(taints, &KubernetesNodePoolTaintArgs{
        Key:    pulumi.String(taint.Key),
        Value:  pulumi.String(taint.Value),
        Effect: pulumi.String(taint.Effect),
    })
}
```

#### Autoscaling Configuration

```go
if spec.AutoScale {
    args.AutoScale = pulumi.BoolPtr(true)
    args.MinNodes = pulumi.IntPtr(int(spec.MinNodes))
    args.MaxNodes = pulumi.IntPtr(int(spec.MaxNodes))
}
```

### 3. Outputs (outputs.go)

**Exports:**
- `node_pool_id`: Node pool UUID for cross-stack references

## Resource Provisioning Flow

```
User
  │
  ▼
pulumi up
  │
  ▼
main.go
  │
  ├─ Load manifest YAML
  ├─ Initialize Pulumi context
  │
  ▼
module/main.go
  │
  ├─ Initialize locals
  ├─ Create DO provider
  │
  ▼
module/node_pool.go
  │
  ├─ Merge labels
  ├─ Build taints array
  ├─ Build tags array
  ├─ Create node pool resource
  ├─ Export outputs
  │
  ▼
DigitalOcean API
  │
  ├─ Create node pool
  ├─ Provision nodes
  │
  ▼
Complete
```

## Configuration Processing

### Spec Validation

Validation happens at protobuf level before Pulumi execution:

```proto
message DigitalOceanKubernetesNodePoolSpec {
  string node_pool_name = 1 [(buf.validate.field).required = true];
  // ...
  uint32 node_count = 4 [
    (buf.validate.field).required = true,
    (buf.validate.field).uint32.gt = 0
  ];
}
```

### Field Transformations

#### Cluster Reference Resolution

```go
ClusterId: pulumi.String(spec.Cluster.GetValue())
```

Extracts cluster ID from `StringValueOrRef` (supports both direct values and cross-stack references).

#### Conditional Taints

```go
if len(taints) > 0 {
    args.Taints = taints
}
```

Only includes taints if specified (avoids empty array in API call).

## Design Decisions

### 1. Label Merging Pattern

**Decision:** Merge metadata labels with spec labels, allowing overrides

**Rationale:**
- Metadata provides consistent labeling across resources
- Spec labels enable user customization
- Override capability provides flexibility

**Tradeoff:** Users can override standard labels (intentional for flexibility)

### 2. Separate Taint Processing

**Decision:** Build taints array separately before node pool args

**Rationale:**
- Clearer code structure
- Easier to test taint logic
- Conditional inclusion logic is explicit

### 3. Autoscaling Conditional

**Decision:** Only set min/max nodes if autoscaling enabled

**Rationale:**
- API doesn't accept min/max for fixed-size pools
- Explicit conditional prevents API errors
- Clearer intent in code

## Error Handling

```go
createdNodePool, err := digitalocean.NewKubernetesNodePool(...)
if err != nil {
    return nil, errors.Wrap(err, "failed to create node pool")
}
```

**Pattern:**
1. Wrap errors with context
2. Return immediately on failure
3. Propagate to Pulumi runtime

## Production Patterns

### Multi-Pool Architecture

Create multiple node pools for workload separation:

```go
// System pool (minimal, never touched)
systemPool := createNodePool(ctx, systemSpec)

// Application pool (autoscaling)
appPool := createNodePool(ctx, appSpec)

// Batch pool (scale-to-zero)
batchPool := createNodePool(ctx, batchSpec)
```

### Label Conventions

**Recommended labels:**
- `workload`: application|batch|system
- `tier`: frontend|backend|data
- `env`: production|staging|dev
- `team`: team-name

### Taint Use Cases

**GPU Isolation:**
```go
{
    Key: "nvidia.com/gpu",
    Value: "true",
    Effect: "NoSchedule",
}
```

**Dedicated Pools:**
```go
{
    Key: "dedicated",
    Value: "database",
    Effect: "NoSchedule",
}
```

## Customization Points

### Adding Node Pool Metadata

To expose additional metadata fields:

1. Update `spec.proto` with new field
2. Update `locals.go` to process new field
3. Update `node_pool.go` to use new field

### Supporting Multiple Pools

Current implementation creates single pool. For multiple pools:

```go
for _, poolSpec := range specs {
    pool, err := nodePool(ctx, poolSpec, provider)
    // ...
}
```

## Performance Considerations

- **Deployment time:** ~2-3 minutes per node pool
- **State size:** ~2 KB per node pool
- **API calls:** ~3-5 per deployment

## Testing

### Unit Tests

```go
func TestNodePoolCreation(t *testing.T) {
    // Mock Pulumi context
    // Test node pool creation
    // Verify label merging
    // Verify taint processing
}
```

### Integration Tests

```bash
# Deploy test pool
pulumi stack init test
pulumi up --yes

# Verify
kubectl get nodes -l planton-resource-name=test-pool

# Cleanup
pulumi destroy --yes
```

## Reference

- [Pulumi DigitalOcean Provider Source](https://github.com/pulumi/pulumi-digitalocean)
- [DigitalOcean Kubernetes API](https://docs.digitalocean.com/reference/api/api-reference/#tag/Kubernetes)
- [OpenMCF Architecture](../../../../../../architecture/README.md)

---

**Version:** 1.0.0  
**Last Updated:** 2025-11-16  
**Authors:** OpenMCF Engineering Team

