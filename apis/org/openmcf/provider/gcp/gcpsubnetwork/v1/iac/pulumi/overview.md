# Pulumi Implementation Overview: GCP Subnetwork

## Purpose

This document provides an architectural overview of the Pulumi implementation for `GcpSubnetwork`, explaining design decisions, module structure, and key implementation patterns.

## Why Pulumi?

OpenMCF uses **Pulumi** as the default IaC engine for several reasons:

### 1. **Multi-Cloud Consistency**

Pulumi's programming model works uniformly across AWS, Azure, GCP, and Kubernetes. The same patterns for resource creation, dependency management, and output handling apply everywhere.

### 2. **Real Programming Languages**

Using Go (instead of HCL or YAML) enables:
- **Type safety**: Compile-time checks for invalid configurations
- **Code reuse**: Share common logic via Go packages and modules
- **Complex logic**: Express conditionals, loops, and transformations naturally
- **IDE support**: IntelliSense, refactoring, and debugging tools

### 3. **Protobuf Integration**

OpenMCF's API definitions are protobuf-based. Pulumi's Go SDK integrates seamlessly:
- Load `GcpSubnetworkStackInput` from Pulumi config
- Unmarshal protobuf JSON into Go structs
- Pass structs directly to Pulumi resource constructors

### 4. **Open Source**

Pulumi's core engine is Apache 2.0 licensed, avoiding vendor lock-in. State can be self-hosted in GCS, S3, or Azure Blob Storage.

## Architecture

### Module Structure

```
iac/pulumi/
├── main.go                # Entry point (loads stack input, calls module)
├── module/
│   ├── main.go            # Module entry (orchestrates resources)
│   ├── locals.go          # Locals bundle (spec, metadata, labels)
│   ├── outputs.go         # Output constants
│   └── subnetwork.go      # Core implementation (google_compute_subnetwork)
├── Pulumi.yaml            # Pulumi project config
├── Makefile               # Deployment helpers
└── debug.sh               # Local debugging script
```

### Data Flow

```
manifest.yaml
    ↓ (openmcf CLI)
GcpSubnetworkStackInput (JSON)
    ↓ (Pulumi config)
main.go
    ↓ (unmarshal protobuf)
module.Resources()
    ↓ (create GCP provider)
subnetwork()
    ↓ (Pulumi resources)
google_compute_subnetwork
    ↓ (GCP API)
Created Subnet
    ↓ (Pulumi exports)
Stack Outputs
```

## Key Components

### 1. Entry Point: `main.go`

**Responsibilities**:
- Load `GcpSubnetworkStackInput` from Pulumi config (`stack-input` key)
- Unmarshal protobuf JSON into Go struct
- Invoke `module.Resources()` to create infrastructure
- Return error if unmarshalling or resource creation fails

**Code Pattern**:

```go
func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        // Load stack input from Pulumi config
        stackInput := &gcpsubnetworkv1.GcpSubnetworkStackInput{}
        if err := loadStackInput(ctx, stackInput); err != nil {
            return err
        }
        
        // Create resources
        return module.Resources(ctx, stackInput)
    })
}
```

**Why Separate?**: Keeps module code testable and reusable across different entry points (CLI, server, tests).

### 2. Module Orchestrator: `module/main.go`

**Responsibilities**:
- Create `Locals` bundle (aggregates spec, metadata, labels)
- Configure GCP provider with credentials from `GcpProviderConfig`
- Call resource creation functions (`subnetwork()`)
- Propagate errors up to entry point

**Code Pattern**:

```go
func Resources(ctx *pulumi.Context, stackInput *gcpsubnetworkv1.GcpSubnetworkStackInput) error {
    locals := &Locals{GcpSubnetwork: stackInput.Target}
    
    provider, err := createGcpProvider(ctx, stackInput.ProviderConfig)
    if err != nil {
        return err
    }
    
    _, err = subnetwork(ctx, locals, provider)
    return err
}
```

**Why Orchestrate?**: Centralizes provider setup and error handling, making it easy to add new resources later.

### 3. Locals Bundle: `module/locals.go`

**Responsibilities**:
- Aggregate `GcpSubnetwork` spec and metadata
- Compute resource labels (resource_id, resource_kind, organization, environment)
- Provide utility functions for label merging and defaults

**Code Pattern**:

```go
type Locals struct {
    GcpSubnetwork *gcpsubnetworkv1.GcpSubnetwork
}

func (l *Locals) ResourceLabels() pulumi.StringMap {
    return pulumi.StringMap{
        "resource":      pulumi.String("true"),
        "resource_id":   pulumi.String(l.resourceId()),
        "resource_kind": pulumi.String("gcp_subnetwork"),
        "organization":  pulumi.String(l.GcpSubnetwork.Metadata.Org),
        "environment":   pulumi.String(l.GcpSubnetwork.Metadata.Env),
    }
}
```

**Why Bundle?**: Avoids passing multiple arguments to every resource function; centralizes label computation.

### 4. Core Implementation: `module/subnetwork.go`

**Responsibilities**:
- Enable required GCP APIs (`compute.googleapis.com`)
- Prepare secondary IP ranges from spec
- Create `google_compute_subnetwork` resource
- Export outputs (self-link, region, CIDR, secondary ranges)

**Implementation Breakdown**:

#### Step 1: Enable APIs

```go
var subnetworkProjectApis = []string{"compute.googleapis.com"}

for _, api := range subnetworkProjectApis {
    createdProjectService, err := projects.NewService(ctx, "subnetwork-"+api, 
        &projects.ServiceArgs{
            Project: pulumi.String(locals.GcpSubnetwork.Spec.ProjectId.GetValue()),
            Service: pulumi.String(api),
            DisableDependentServices: pulumi.BoolPtr(true),
        }, pulumi.Provider(provider))
    if err != nil {
        return nil, errors.Wrapf(err, "failed to enable %s api", api)
    }
    createdGoogleApiResources = append(createdGoogleApiResources, createdProjectService)
}
```

**Why?**: GCP requires compute API to be explicitly enabled. Automation ensures subnets can be created even in new projects.

**Performance**: APIs are enabled in parallel (no sequential blocking).

#### Step 2: Prepare Secondary Ranges

```go
var secondaryRanges compute.SubnetworkSecondaryIpRangeArray
for _, r := range locals.GcpSubnetwork.Spec.SecondaryIpRanges {
    secondaryRanges = append(secondaryRanges, &compute.SubnetworkSecondaryIpRangeArgs{
        RangeName:   pulumi.String(r.RangeName),
        IpCidrRange: pulumi.String(r.IpCidrRange),
    })
}
```

**Why?**: Pulumi expects `SubnetworkSecondaryIpRangeArray`, so we convert protobuf repeated fields.

**Handling Optionals**: If spec has zero secondary ranges, `secondaryRanges` is empty (valid).

#### Step 3: Create Subnet

```go
createdSubnetwork, err := compute.NewSubnetwork(ctx, "subnetwork",
    &compute.SubnetworkArgs{
        Name:                  pulumi.String(locals.GcpSubnetwork.Metadata.Name),
        Project:               pulumi.StringPtr(locals.GcpSubnetwork.Spec.ProjectId.GetValue()),
        Region:                pulumi.String(locals.GcpSubnetwork.Spec.Region),
        Network:               pulumi.String(locals.GcpSubnetwork.Spec.VpcSelfLink.GetValue()),
        IpCidrRange:           pulumi.String(locals.GcpSubnetwork.Spec.IpCidrRange),
        PrivateIpGoogleAccess: pulumi.BoolPtr(locals.GcpSubnetwork.Spec.PrivateIpGoogleAccess),
        SecondaryIpRanges:     secondaryRanges,
    },
    pulumi.Provider(provider),
    pulumi.DependsOn(createdGoogleApiResources))
```

**Key Design Choices**:
- **Resource Name**: Pulumi logical name is `"subnetwork"` (stable across deployments)
- **GCP Name**: Uses `metadata.name` from spec (user-controlled, matches K8s patterns)
- **Dependencies**: Explicit `DependsOn` ensures APIs are enabled first
- **Provider**: Uses custom provider configured with user credentials

#### Step 4: Export Outputs

```go
ctx.Export(OpSubnetworkSelfLink, createdSubnetwork.SelfLink)
ctx.Export(OpRegion, createdSubnetwork.Region)
ctx.Export(OpIpCidrRange, createdSubnetwork.IpCidrRange)
ctx.Export(OpSecondaryRanges, createdSubnetwork.SecondaryIpRanges)
```

**Why?**: Downstream resources (GKE clusters, VMs) need these outputs to reference the subnet.

**Output Constants**: Defined in `module/outputs.go` for consistency:

```go
const (
    OpSubnetworkSelfLink = "subnetwork_self_link"
    OpRegion             = "region"
    OpIpCidrRange        = "ip_cidr_range"
    OpSecondaryRanges    = "secondary_ranges"
)
```

## Design Patterns

### Pattern 1: API Enablement as Resources

**Problem**: GCP APIs must be enabled before using them, but checking enablement is slow and racy.

**Solution**: Treat API enablement as Pulumi resources:

```go
projects.NewService(ctx, "subnetwork-compute.googleapis.com", ...)
```

**Benefits**:
- **Idempotent**: Pulumi tracks whether API is enabled in state
- **Dependency-aware**: Other resources wait for API enablement via `DependsOn`
- **Parallel**: Multiple APIs enabled concurrently

### Pattern 2: Optional Fields with Protobuf

**Problem**: Protobuf optional fields need special handling in Go.

**Solution**: Use `GetValue()` for `StringValueOrRef` types:

```go
// spec.proto:
// org.openmcf.shared.foreignkey.v1.StringValueOrRef project_id = 1;

// Go code:
Project: pulumi.String(locals.GcpSubnetwork.Spec.ProjectId.GetValue())
```

**Why?**: `StringValueOrRef` supports both direct values (`"my-project"`) and references to other resources.

### Pattern 3: Locals Bundle for DRY Code

**Problem**: Resource functions need access to spec, metadata, labels, etc.

**Solution**: Aggregate everything in a `Locals` struct:

```go
type Locals struct {
    GcpSubnetwork *gcpsubnetworkv1.GcpSubnetwork
}

func subnetwork(ctx *pulumi.Context, locals *Locals, provider *gcp.Provider) { ... }
```

**Benefits**:
- Single argument instead of 5+
- Centralizes label computation
- Easy to extend with new fields

### Pattern 4: Error Wrapping for Context

**Problem**: Generic errors like "failed to create subnetwork" lack context.

**Solution**: Wrap errors with `github.com/pkg/errors`:

```go
if err != nil {
    return nil, errors.Wrap(err, "failed to create subnetwork")
}
```

**Benefits**:
- Stack traces show exactly where failures occur
- Error messages include operation context

### Pattern 5: Explicit Dependencies

**Problem**: Pulumi infers dependencies from property references, but sometimes relationships are implicit.

**Solution**: Use explicit `DependsOn`:

```go
pulumi.DependsOn(createdGoogleApiResources)
```

**Why?**: API enablement doesn't produce a property that subnet needs, but subnet *must* wait for it.

## State Management

### Pulumi State

Pulumi stores:
- **Resources**: Mapping of logical names to cloud resource IDs
- **Outputs**: Exported values (self-link, region, etc.)
- **Inputs**: Configuration and stack input (encrypted secrets)

### State Backends

#### 1. Pulumi Service (Default)

```shell
pulumi login
```

**Pros**:
- Managed service, no infrastructure to maintain
- Built-in secret encryption (AES-256)
- Team collaboration (RBAC, audit logs)

**Cons**:
- Requires Pulumi account
- SaaS dependency

#### 2. Self-Hosted (GCS)

```shell
pulumi login gs://my-state-bucket
```

**Pros**:
- Full control over state storage
- Use existing GCS bucket infrastructure
- No external dependencies

**Cons**:
- Must manage encryption and access control
- No built-in team features

#### 3. Local (Development Only)

```shell
pulumi login --local
```

**Pros**: No setup, instant start

**Cons**: State in `~/.pulumi/`, not shareable, easily lost

## Resource Naming

### Pulumi Logical Names

Logical names are **stable across deployments** and used for state tracking:

```go
compute.NewSubnetwork(ctx, "subnetwork", ...)
```

**Rule**: Use fixed strings like `"subnetwork"`, not dynamic values like `metadata.name`.

**Why?**: Changing logical names causes Pulumi to recreate resources.

### GCP Physical Names

Physical names are **user-visible** in GCP Console:

```go
Name: pulumi.String(locals.GcpSubnetwork.Metadata.Name)
```

**Rule**: Use `metadata.name` from spec.

**Why?**: Users expect resource names to match their YAML manifests.

## Error Handling Strategy

### Fail Fast

If any resource creation fails, Pulumi halts and rolls back:

```go
if err != nil {
    return nil, errors.Wrap(err, "failed to create subnetwork")
}
```

**Benefits**:
- No partial deployments
- Clear error messages
- Automatic rollback on failure

### Retry Behavior

Pulumi retries transient failures automatically:
- API rate limits (429)
- Network timeouts
- "Operation in progress" conflicts

**Configuration**: Controlled by Pulumi SDK (no custom retry logic needed).

## Testing Strategy

### Unit Tests

Test individual functions in isolation:

```go
func TestSecondaryRangeParsing(t *testing.T) {
    spec := &gcpsubnetworkv1.GcpSubnetworkSpec{
        SecondaryIpRanges: []*gcpsubnetworkv1.GcpSubnetworkSecondaryRange{
            {RangeName: "pods", IpCidrRange: "10.0.0.0/18"},
        },
    }
    
    ranges := prepareSecondaryRanges(spec)
    assert.Equal(t, 1, len(ranges))
    assert.Equal(t, "pods", ranges[0].RangeName)
}
```

### Integration Tests

Deploy to real GCP project and verify outputs:

```shell
make deploy MANIFEST=../../hack/manifest.yaml
pulumi stack output subnetwork_self_link
# Expected: projects/my-project/regions/us-central1/subnetworks/test-subnet
```

### E2E Tests

Create subnet → Create GKE cluster using subnet → Verify cluster uses correct IP ranges.

## Performance Characteristics

### Deployment Time

Typical timings for fresh deployment:

| Operation | Duration | Notes |
|-----------|----------|-------|
| API Enablement | 10-30s | Only first time per project |
| Subnet Creation | 5-10s | GCP resource creation |
| Total | 15-40s | First deployment |
| Subsequent Deploys | 5-10s | APIs already enabled |

### Resource Limits

- **Subnets per VPC**: 500 (GCP quota)
- **Secondary ranges per subnet**: 170 (GCP limit)
- **Pulumi state size**: <100KB per subnet (minimal)

## Security Considerations

### Credential Handling

- **Never hardcode** service account keys in code
- **Use environment variables** or Pulumi secrets for credentials
- **Prefer Workload Identity** (GKE) or Application Default Credentials (local)

### Secret Management

Encrypt sensitive outputs:

```go
ctx.Export("private_key", createdKey.PrivateKey, pulumi.Secret(true))
```

Pulumi encrypts secrets in state using AES-256.

### IAM Requirements

Minimum IAM roles for deployment:

- `roles/compute.networkAdmin` - Create subnetworks
- `roles/serviceusage.serviceUsageAdmin` - Enable APIs

## Extending the Module

### Adding Flow Logs

1. Update `spec.proto`:

```protobuf
message GcpSubnetworkSpec {
    // ... existing fields
    optional bool enable_flow_logs = 7;
    optional GcpSubnetworkFlowLogsConfig flow_logs_config = 8;
}
```

2. Update `module/subnetwork.go`:

```go
LogConfig: pulumi.SubnetworkLogConfigPtrInput(&compute.SubnetworkLogConfigArgs{
    Enable: pulumi.Bool(locals.GcpSubnetwork.Spec.EnableFlowLogs),
}),
```

3. Regenerate protos: `make protos`

4. Test with updated manifest

### Adding Purpose/Role Fields

For special subnet types (e.g., `INTERNAL_HTTPS_LOAD_BALANCER`):

```protobuf
optional string purpose = 9; // e.g., "INTERNAL_HTTPS_LOAD_BALANCER"
optional string role = 10;    // e.g., "ACTIVE"
```

Then map in Pulumi:

```go
Purpose: pulumi.StringPtr(locals.GcpSubnetwork.Spec.Purpose),
Role:    pulumi.StringPtr(locals.GcpSubnetwork.Spec.Role),
```

## Comparison to Terraform

| Aspect | Pulumi (This Implementation) | Terraform |
|--------|------------------------------|-----------|
| Language | Go (strongly typed) | HCL (domain-specific) |
| State | Pulumi Service or GCS | GCS, S3, or Terraform Cloud |
| Secrets | Native encryption | External secret providers |
| Dependencies | Inferred + explicit `DependsOn` | Inferred from references |
| Output Handling | Strongly typed (Go structs) | Map[string]interface{} |
| API Enablement | Explicit resources | Manual or module-based |
| Module Reuse | Go packages | Terraform modules (HCL) |

**Verdict**: Pulumi offers better type safety and code reuse, while Terraform has a larger ecosystem and more mature tooling.

## Further Reading

- [Pulumi Go SDK Documentation](https://www.pulumi.com/docs/reference/pkg/go/)
- [GCP Compute Subnetwork API](https://cloud.google.com/compute/docs/reference/rest/v1/subnetworks)
- [Pulumi GCP Provider](https://www.pulumi.com/registry/packages/gcp/)
- [OpenMCF Architecture](../../../../../architecture/deployment-component.md)

## Related Documentation

- **User Guide**: `../../README.md` - How to use `GcpSubnetwork`
- **Examples**: `../../examples.md` - Practical usage scenarios
- **Research**: `../../docs/README.md` - Deep dive into GCP networking
- **Pulumi Deployment**: `README.md` - Local development workflow

---

**Questions?** Consult the [module README](README.md) or [component documentation](../../README.md).

