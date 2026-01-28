# Pulumi Module Architecture - GCP GKE Node Pool

This document explains the internal architecture of the GcpGkeNodePool Pulumi module, including design patterns, resource relationships, and implementation details.

## Architecture Overview

### Module Structure

```
pulumi/
├── main.go              # Entrypoint: reads stack input, initializes module
├── module/
│   ├── main.go          # Orchestration: creates provider, looks up cluster, provisions node pool
│   ├── locals.go        # Data transformations and computed values
│   ├── outputs.go       # Output constants and mappings
│   └── node_pool.go     # Node pool resource provisioning logic
```

### Data Flow

```
stack_input.json (external)
    ↓
main.go: pulumi.Run()
    ↓
module.Resources(ctx, stackInput)
    ↓
initializeLocals() → Locals struct
    ↓
container.LookupCluster() → parent cluster info
    ↓
nodePool() → container.NewNodePool()
    ↓
ctx.Export() → stack outputs
```

## Core Components

### 1. Entrypoint (`main.go`)

**Responsibilities:**
- Read `stack_input` configuration from Pulumi config
- Unmarshal JSON/YAML into `GcpGkeNodePoolStackInput` protobuf message
- Invoke `module.Resources()` to create infrastructure

**Key code:**
```go
func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        var stackInput gcpgkenodepoolv1.GcpGkeNodePoolStackInput
        ctx.Config().GetObject("stack_input", &stackInput)
        return module.Resources(ctx, &stackInput)
    })
}
```

### 2. Module Orchestration (`module/main.go`)

**Responsibilities:**
- Initialize locals (labels, metadata, computed values)
- Create GCP provider with proper project configuration
- Lookup parent GKE cluster (validates cluster exists)
- Invoke node pool provisioning
- Handle errors

**Key functions:**
- `Resources(ctx, stackInput)`: Main orchestration entry point
- Creates `gcp.Provider` with project from `cluster_project_id`
- Calls `container.LookupCluster()` to fetch parent cluster details
- Delegates to `nodePool()` for resource creation

### 3. Locals (`module/locals.go`)

**Purpose:** Data transformations and computed values, mimicking Terraform's `locals` block.

**`Locals` struct:**
```go
type Locals struct {
    GcpGkeNodePool *gcpgkenodepoolv1.GcpGkeNodePool
    GcpLabels      map[string]string    // Standard GCP resource labels
    KubernetesLabels map[string]string  // Kubernetes-compatible labels
    NetworkTag     string               // Firewall tag: "gke-<cluster-name>"
    ClusterName    string               // Resolved cluster name
}
```

**Label Strategy:**
- **Base labels**: `resource`, `resource_kind`, `resource_name`
- **Conditional labels**: `organization`, `environment`, `resource_id` (added only if metadata fields are non-empty)
- **Network tag**: Follows GKE convention `gke-<cluster-name>` for firewall rule integration

**Initialization:**
```go
func initializeLocals(ctx *pulumi.Context, stackInput *gcpgkenodepoolv1.GcpGkeNodePoolStackInput) *Locals {
    // Extract target resource
    // Resolve foreign key references (cluster_name)
    // Build label maps
    // Compute network tag
}
```

### 4. Node Pool Resource (`module/node_pool.go`)

**Purpose:** Provisions the `google_container_node_pool` resource with all configuration.

**Key logic:**

#### Fixed Size vs. Autoscaling

```go
if spec.NodePoolSize != nil {
    switch x := spec.NodePoolSize.(type) {
    case *gcpgkenodepoolv1.GcpGkeNodePoolSpec_NodeCount:
        nodeCount = pulumi.IntPtr(int(x.NodeCount))
    case *gcpgkenodepoolv1.GcpGkeNodePoolSpec_Autoscaling:
        autoscaling = &container.NodePoolAutoscalingArgs{
            MinNodeCount: pulumi.Int(int(x.Autoscaling.MinNodes)),
            MaxNodeCount: pulumi.Int(int(x.Autoscaling.MaxNodes)),
            LocationPolicy: pulumi.StringPtr(x.Autoscaling.GetLocationPolicy()),
        }
    }
}
```

**Pattern:** `oneof` in protobuf maps to Go interface with type assertion. Exactly one must be set (enforced by proto validation).

#### Management Settings

```go
if spec.Management != nil {
    management = &container.NodePoolManagementArgs{
        AutoUpgrade: pulumi.Bool(!spec.Management.DisableAutoUpgrade),
        AutoRepair:  pulumi.Bool(!spec.Management.DisableAutoRepair),
    }
} else {
    management = &container.NodePoolManagementArgs{
        AutoUpgrade: pulumi.Bool(true),  // Default: enabled
        AutoRepair:  pulumi.Bool(true),
    }
}
```

**Pattern:** Invert `disable_*` flags (protobuf convention) to `enabled` flags (GCP API convention). Defaults to enabled if management block is omitted.

#### Label Merging

```go
mergedLabels := map[string]string{}
for k, v := range locals.GcpLabels {
    mergedLabels[k] = v
}
for k, v := range spec.NodeLabels {
    mergedLabels[k] = v
}
```

**Pattern:** Merge OpenMCF standard labels with user-specified node labels. User labels take precedence on conflicts.

#### Node Config

```go
nodeConfig := &container.NodePoolNodeConfigArgs{
    MachineType: pulumi.String(spec.GetMachineType()),  // Uses proto default
    Preemptible: pulumi.Bool(spec.Spot),
    Labels:      pulumi.ToStringMap(mergedLabels),
    Tags:        pulumi.StringArray{pulumi.String(locals.NetworkTag)},
    Metadata: pulumi.StringMap{
        "disable-legacy-endpoints": pulumi.String("true"),  // Security best practice
    },
    OauthScopes: pulumi.StringArray{
        pulumi.String("https://www.googleapis.com/auth/monitoring"),
        pulumi.String("https://www.googleapis.com/auth/logging.write"),
        pulumi.String("https://www.googleapis.com/auth/devstorage.read_only"),
    },
    ImageType: pulumi.StringPtr(spec.GetImageType()),  // Uses proto default
}
```

**Key details:**
- `GetMachineType()` and `GetImageType()` return proto defaults (`e2-medium`, `COS_CONTAINERD`)
- `disable-legacy-endpoints: true` is a GKE security best practice (prevents GCE metadata API v1 access)
- OAuth scopes are minimal: monitoring, logging, read-only GCS

#### Conditional Fields

```go
if spec.DiskSizeGb > 0 {
    nodeConfig.DiskSizeGb = pulumi.IntPtr(int(spec.DiskSizeGb))
}
if spec.DiskType != nil && *spec.DiskType != "" {
    nodeConfig.DiskType = pulumi.StringPtr(*spec.DiskType)
}
if spec.ServiceAccount != "" {
    nodeConfig.ServiceAccount = pulumi.StringPtr(spec.ServiceAccount)
}
```

**Pattern:** Conditionally set fields only if non-zero/non-empty to allow GKE defaults to apply.

#### Resource Options

```go
container.NewNodePool(ctx, name, args,
    pulumi.Provider(gcpProvider),           // Use configured provider
    pulumi.IgnoreChanges([]string{"nodeCount"}),  // Ignore autoscaler changes
    pulumi.DeleteBeforeReplace(true))       // Avoid duplicate pools
```

**Key options:**
- `IgnoreChanges(["nodeCount"])`: When autoscaling is enabled, the autoscaler changes `node_count` dynamically. Pulumi should not revert these changes.
- `DeleteBeforeReplace(true)`: For node pool replacement (e.g., machine type change), delete old pool before creating new one to avoid name conflicts.

### 5. Outputs (`module/outputs.go`)

**Purpose:** Define output key constants and export stack outputs.

**Output constants:**
```go
const (
    OpNodePoolName      = "node_pool_name"
    OpInstanceGroupUrls = "instance_group_urls"
    OpCurrentNodeCount  = "current_node_count"
    OpMinNodes          = "min_nodes"
    OpMaxNodes          = "max_nodes"
)
```

**Export logic:**
```go
ctx.Export(OpNodePoolName, createdNodePool.Name)
ctx.Export(OpInstanceGroupUrls, createdNodePool.InstanceGroupUrls)
ctx.Export(OpCurrentNodeCount, createdNodePool.NodeCount)

if nodeCount != nil {
    ctx.Export(OpMinNodes, nodeCount)
    ctx.Export(OpMaxNodes, nodeCount)
} else if autoscaling != nil {
    ctx.Export(OpMinNodes, autoscaling.MinNodeCount)
    ctx.Export(OpMaxNodes, autoscaling.MaxNodeCount)
}
```

**Pattern:** For fixed-size pools, min/max both equal `node_count`. For autoscaling pools, export the configured min/max values.

## Resource Relationships

### Parent-Child Hierarchy

```
GcpProject (project_id)
    ↓
GcpGkeCluster (control plane)
    ↓
GcpGkeNodePool (worker nodes)
    ↓
Compute Engine VMs (managed by GKE)
```

**Foreign key resolution:**
- `cluster_project_id`, `cluster_name`, and `cluster_location` are foreign key references to the parent `GcpGkeCluster`
- All required cluster information is provided in the spec (no additional API lookups needed)
- Node pool is created in the same location (region or zone) as the parent cluster

### Implicit Dependencies

- **VPC/Subnetwork**: Inherited from parent cluster (not specified in node pool spec)
- **Firewall rules**: Node pool VMs get network tag `gke-<cluster-name>`, enabling firewall rule targeting
- **IAM**: Nodes use either default GKE service account or custom `service_account` specified in spec

## Design Patterns

### 1. Separate Resource Pattern

**Rationale:** Node pools are provisioned as **separate, independent resources** from the cluster, not as inline configuration blocks within the cluster definition.

**Benefits:**
- **Lifecycle independence**: Update node pool machine type without risking cluster control plane disruption
- **Scalability**: Manage multiple node pools per cluster declaratively
- **Industry alignment**: Matches Terraform, GKE API, and production best practices

### 2. Foreign Key References

**Pattern:** Use `StringValueOrRef` for cluster references, allowing either:
- **Literal values**: `{ value: "my-cluster" }`
- **Resource references**: `{ resource: { kind: "GcpGkeCluster", name: "my-cluster", field_path: "metadata.name" } }`

**Implementation:** Module reads `.GetValue()` which resolves either path. CLI handles reference resolution and dependency ordering.

### 3. Proto Defaults

**Pattern:** Use `(org.openmcf.shared.options.default)` in protobuf for opinionated defaults:
- `machine_type`: `e2-medium`
- `disk_type`: `pd-standard`
- `image_type`: `COS_CONTAINERD`
- `location_policy`: `BALANCED`

**Implementation:** `spec.GetMachineType()` returns default if field is unset. No conditional logic needed in module code.

### 4. Inversion of Boolean Flags

**Pattern:** Proto uses `disable_auto_upgrade` and `disable_auto_repair` (explicit opt-out), GCP API uses `auto_upgrade: true/false` (enabled by default).

**Implementation:** Module inverts: `AutoUpgrade: pulumi.Bool(!spec.Management.DisableAutoUpgrade)`

**Rationale:** Makes enabled state explicit in API while defaulting to secure, production-ready configuration.

### 5. Label Standardization

**Pattern:** All OpenMCF resources apply standard labels:
- `resource: "true"`
- `resource_kind: "gcp-gke-node-pool"`
- `resource_name: <metadata.name>`
- `resource_id: <metadata.id>` (if set)
- `organization: <metadata.org>` (if set)
- `environment: <metadata.env>` (if set)

**Benefits:**
- Unified cost tracking across all OpenMCF resources
- Consistent filtering and querying in GCP Console
- Integration with billing dashboards and reporting tools

## Error Handling

### Cluster Lookup Failure

If `container.LookupCluster()` fails (cluster doesn't exist or wrong project/name):

```go
clusterInfo, err := container.LookupCluster(ctx, &container.LookupClusterArgs{
    Name:     locals.ClusterName,
    Project:  locals.GcpGkeNodePool.Spec.ClusterProjectId.GetValue(),
    Location: "*",
})
if err != nil {
    return errors.Wrap(err, "failed to lookup parent GKE cluster")
}
```

**Error message:** Clear indication that parent cluster is missing or inaccessible.

### Node Pool Creation Failure

Common causes and error messages:
- **Quota exceeded**: "Quota 'CPUS' exceeded. Limit: X in region Y."
- **Invalid machine type**: "Invalid value for field 'machineType': 'invalid-type'."
- **Spot VM unavailable**: Partial success (some nodes created, others pending capacity).

**Handling:** Pulumi propagates GCP API error messages directly. Retry logic is handled by Pulumi's state management.

## Testing and Validation

### Unit Tests (`module/*_test.go`)

Future improvement: Test local value transformations, label merging, and conditional logic in isolation.

### Integration Tests (`hack/manifest.yaml`)

The `debug.sh` script reads `hack/manifest.yaml` and provisions a real node pool for end-to-end testing.

### Validation Points

1. **Proto validation** (spec_test.go): Validates buf.validate rules before Pulumi runs
2. **Cluster lookup**: Fails fast if parent cluster doesn't exist
3. **GCP API validation**: GCP validates machine types, disk types, quotas, etc.

## Performance Considerations

### Resource Creation Time

- **Fixed-size pool**: ~3-5 minutes (depends on `node_count` and VM startup time)
- **Autoscaling pool**: ~3-5 minutes for initial min_nodes, additional nodes scale as needed
- **Spot VM pool**: May take longer if hunting for capacity across zones

### State Size

Pulumi state for a single node pool is ~2-3 KB, negligible for most backends.

### API Rate Limits

GKE API has rate limits (quota `READ_REQUESTS` and `WRITE_REQUESTS`). For bulk operations, consider:
- Batching node pool creation across time
- Using separate GCP projects to distribute quota

## Future Enhancements

### Roadmap Items

1. **Taints and tolerations**: Allow workload isolation via node taints
2. **GPU accelerators**: Support `guest_accelerator` configuration for ML workloads
3. **Advanced kernel tuning**: Expose sysctls and kubelet config
4. **Surge upgrade settings**: Make `max_surge` and `max_unavailable` configurable

### Extensibility

The module is designed for easy extension:
- Add new fields to `GcpGkeNodePoolSpec` proto
- Update `node_pool.go` to map new fields to Pulumi args
- Proto validation and defaults handle most configuration logic

## Related Documentation

- **[Pulumi README](./README.md)**: Usage guide and troubleshooting
- **[Component Overview](../../README.md)**: High-level user documentation
- **[Research Document](../../docs/README.md)**: Deep dive into GKE node pools
- **[Pulumi GCP Provider Docs](https://www.pulumi.com/registry/packages/gcp/)**: API reference

## Contributing

When modifying this module:
1. Maintain the separate resource pattern (don't inline node pools in cluster)
2. Follow label standardization conventions
3. Use proto defaults instead of hardcoded values
4. Test with both fixed-size and autoscaling configurations
5. Validate with `hack/manifest.yaml` before committing

