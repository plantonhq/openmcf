# GCP VPC Pulumi Module - Architecture Overview

## Purpose

This document explains the internal architecture of the Pulumi module for deploying GCP VPCs. It is intended for contributors, maintainers, and advanced users who want to understand how the module works under the hood.

## Design Philosophy

The Pulumi module follows OpenMCF's core principles:

1. **Declarative API**: Users specify the desired VPC configuration via protobuf messages, not imperative code
2. **Best-Practice Defaults**: Custom-mode VPCs and regional routing by default, preventing common pitfalls
3. **Dependency Management**: Automatically enables required GCP APIs before creating resources
4. **Idempotent Operations**: Safe to run multiple times—Pulumi handles state reconciliation
5. **Minimal Boilerplate**: Encapsulates complexity so users don't need to understand Pulumi internals

## Module Structure

```
iac/pulumi/
├── main.go           # Entry point: loads stack input and invokes module
├── Pulumi.yaml       # Project metadata (name, runtime, description)
├── Makefile          # Build automation
├── debug.sh          # Local debugging helper
└── module/
    ├── main.go       # Module coordinator: initializes locals and calls resource functions
    ├── vpc.go        # VPC resource creation with API enablement
    ├── locals.go     # Compute labels, project config, routing mode
    └── outputs.go    # Output constants
```

### Separation of Concerns

- **`main.go`** (entrypoint): Minimal logic—just loads input and delegates to module
- **`module/main.go`**: Orchestrates resource creation, ensures proper ordering
- **`module/locals.go`**: Pure data transformation (input → labels, config)
- **`module/vpc.go`**: Impure infrastructure logic (API calls, resource creation)
- **`module/outputs.go`**: Constants for output keys (avoids magic strings)

This separation makes testing easier and keeps each file focused on a single responsibility.

## Data Flow

### 1. Input Stage

```
User YAML manifest
    ↓
GcpVpcStackInput (protobuf)
    ↓
main.go (entrypoint)
```

The user provides a YAML manifest conforming to the `GcpVpc` API schema. Pulumi deserializes this into a `GcpVpcStackInput` protobuf message containing:
- `target`: The `GcpVpc` resource specification
- `providerConfig`: GCP credentials and configuration

### 2. Locals Initialization

```
GcpVpcStackInput
    ↓
module/locals.go: initializeLocals()
    ↓
Locals struct (labels, config, routing mode)
```

The `initializeLocals` function extracts:
- **GCP Labels**: From metadata (name, org, env, resource kind)
- **Provider Config**: GCP project, credentials
- **Routing Mode**: Maps proto enum (0=REGIONAL, 1=GLOBAL) to GCP API strings

**Why separate locals?**
- Reusable values across multiple resources (if the module grows)
- Clear transformation logic without side effects
- Easy to unit test

### 3. Provider Setup

```
Locals
    ↓
module/main.go: createGcpProvider()
    ↓
Pulumi GCP Provider
```

A GCP provider is created with credentials from `providerConfig`. This provider is used for all GCP API calls in the module, ensuring consistent authentication.

### 4. Resource Creation

```
Provider + Locals
    ↓
module/vpc.go: vpc()
    ↓
GCP Resources
```

The `vpc()` function:
1. **Enables Compute API**: Creates a `google_project_service` resource to enable `compute.googleapis.com`
2. **Creates VPC Network**: Provisions `google_compute_network` with dependencies on the API enablement
3. **Exports Outputs**: Saves `network_self_link` to Pulumi stack outputs

### 5. Output Stage

```
VPC Resource
    ↓
ctx.Export("network_self_link", ...)
    ↓
Pulumi Stack Outputs
```

The VPC's `self_link` attribute is exported so other resources (subnets, GKE clusters, etc.) can reference it.

## Key Implementation Details

### API Enablement Strategy

```go
createdComputeService, err := projects.NewService(ctx,
    "compute-api",
    &projects.ServiceArgs{
        Project: pulumi.String(locals.GcpVpc.Spec.ProjectId.GetValue()),
        Service: pulumi.String("compute.googleapis.com"),
        DisableDependentServices: pulumi.BoolPtr(true),
    }, pulumi.Provider(gcpProvider))
```

**Why enable APIs explicitly?**
- **Deterministic deployments**: Users don't need to manually enable APIs
- **Dependency ordering**: VPC creation depends on API enablement (via `pulumi.DependsOn`)
- **Error prevention**: Avoids cryptic "API not enabled" errors

**Why `DisableDependentServices: true`?**
- When the VPC is destroyed, we want to disable the Compute API cleanly
- Prevents orphaned API enablements

### Routing Mode Mapping

```go
if locals.GcpVpc.Spec.GetRoutingMode() != gcpvpcv1.GcpVpcRoutingMode_REGIONAL {
    // GLOBAL is the only alternative at present.
    networkArgs.RoutingMode = pulumi.StringPtr("GLOBAL")
}
```

The protobuf enum (`GcpVpcRoutingMode`) uses integers:
- `REGIONAL = 0`
- `GLOBAL = 1`

The GCP API expects strings (`"REGIONAL"` or `"GLOBAL"`). The module maps the enum to the string, defaulting to `REGIONAL` if not specified.

**Design choice**: Default to REGIONAL (safest, simplest) unless user explicitly requests GLOBAL.

### Label Management

```go
locals.GcpLabels = map[string]string{
    gcplabelkeys.Resource:     strconv.FormatBool(true),
    gcplabelkeys.ResourceName: locals.GcpVpc.Metadata.Name,
    gcplabelkeys.ResourceKind: strings.ToLower(cloudresourcekind.CloudResourceKind_GcpVpc.String()),
}
```

Labels are generated from metadata for:
- **Cost tracking**: Identify resources by org/env
- **Resource management**: Filter resources by kind, name, etc.
- **Compliance**: Tag all resources with standard labels

**Note**: GCP label keys must be lowercase, so we use `gcplabelkeys` constants.

### Dependency Management

```go
createdNetwork, err := compute.NewNetwork(ctx,
    "vpc",
    networkArgs,
    pulumi.Provider(gcpProvider),
    pulumi.DependsOn([]pulumi.Resource{createdComputeService}))
```

Pulumi's `DependsOn` ensures the VPC is created only after the Compute API is enabled. This prevents race conditions where the API isn't ready yet.

## Error Handling

The module uses Go's standard error handling patterns:

```go
createdNetwork, err := compute.NewNetwork(ctx, "vpc", networkArgs, ...)
if err != nil {
    return nil, errors.Wrap(err, "failed to create vpc network")
}
```

Errors are wrapped with context (e.g., "failed to create vpc network") so users know which step failed. Pulumi propagates these errors to the CLI with stack traces.

## State Management

Pulumi automatically manages state:
- **State File**: Stores current state of all resources (e.g., VPC ID, self-link)
- **State Backend**: Configurable (local file, GCS, Pulumi Service)
- **State Locking**: Prevents concurrent modifications

When you run `pulumi up`:
1. Pulumi reads desired state (from code) and current state (from state file)
2. Computes the diff (what needs to be created/updated/deleted)
3. Shows preview to user
4. Executes changes if approved
5. Updates state file

## Idempotency

The module is idempotent:
- **First run**: Creates VPC
- **Subsequent runs**: No changes (Pulumi detects no diff)
- **Config change**: Pulumi updates VPC in place (e.g., changing routing mode)

**Why this matters**: Safe to run in CI/CD pipelines—won't create duplicate resources.

## Extensibility

### Adding New Configuration Options

To add a new field (e.g., `mtu`):

1. **Update `spec.proto`**:
   ```protobuf
   message GcpVpcSpec {
     // ...
     optional int32 mtu = 4;
   }
   ```

2. **Regenerate Go stubs**:
   ```bash
   make protos
   ```

3. **Update `vpc.go`**:
   ```go
   if locals.GcpVpc.Spec.Mtu != nil {
       networkArgs.Mtu = pulumi.IntPtr(int(locals.GcpVpc.Spec.GetMtu()))
   }
   ```

4. **Test and deploy**:
   ```bash
   pulumi preview
   ```

### Adding New Resources

To create additional resources (e.g., default firewall rules):

1. **Create new file**: `module/firewall.go`
2. **Define function**: `func firewall(ctx *pulumi.Context, locals *Locals, network *compute.Network) error`
3. **Call from `module/main.go`**:
   ```go
   if err := firewall(ctx, locals, createdNetwork); err != nil {
       return err
   }
   ```

## Testing Strategy

### Unit Tests

Test pure functions in isolation:
- `locals.go`: Test label generation, routing mode mapping
- Run: `go test ./module/...`

### Integration Tests

Use Pulumi's testing framework:

```go
func TestVpcCreation(t *testing.T) {
    pulumi.Run(func(ctx *pulumi.Context) error {
        network, err := compute.NewNetwork(ctx, "test-vpc", ...)
        assert.NoError(t, err)
        assert.Equal(t, "test-vpc", network.Name)
        return nil
    })
}
```

### E2E Tests

Deploy to a real GCP project and verify:
1. VPC exists: `gcloud compute networks describe <vpc-name>`
2. Routing mode is correct
3. Labels are applied
4. Clean up: `pulumi destroy`

## Performance Considerations

### Parallel Resource Creation

Pulumi automatically parallelizes independent resources. In this module:
- **Sequential**: Compute API → VPC (explicit dependency)
- **Parallel**: If we added multiple firewall rules, they'd be created in parallel

### State File Size

Each resource adds metadata to the state file. For this module:
- 2 resources (API enablement + VPC)
- Minimal state overhead (~1-2 KB)

As the module grows, consider splitting into separate Pulumi projects if state file exceeds 10 MB.

## Security Considerations

### Credential Management

Credentials are passed via `providerConfig.gcpCredential`:
- **Option 1**: Service account key (JSON) base64-encoded
- **Option 2**: Workload Identity (for GKE environments)
- **Option 3**: Application Default Credentials (local development)

**Best practice**: Never commit credentials to version control. Use secret management tools (Pulumi secrets, Google Secret Manager, etc.).

### IAM Permissions Required

The module requires these IAM roles:
- `roles/compute.networkAdmin` (to create/manage VPCs)
- `roles/serviceusage.serviceUsageAdmin` (to enable APIs)

For production, use a dedicated service account with minimal permissions.

### Least Privilege Principle

The module only requests permissions it needs:
- No `roles/owner` required
- No project-wide modifications (only network resources)

## Comparison to Terraform Module

| Aspect | Pulumi Module | Terraform Module |
|--------|--------------|------------------|
| Language | Go | HCL |
| State | Pulumi state backend | Terraform state backend |
| Type Safety | Strong (Go compiler checks) | Weak (validated at runtime) |
| Testing | Go unit/integration tests | `terraform validate` + `terratest` |
| Abstraction | Full programming language | Declarative DSL |

Both modules generate the same GCP resources. Choose based on team preference and existing workflows.

## Common Patterns

### Conditional Resource Creation

Pulumi uses standard Go conditionals:

```go
if locals.GcpVpc.Spec.EnableFlowLogs {
    // Create flow logs resource
}
```

Compare to Terraform's `count` or `for_each` meta-arguments.

### Resource Loops

Use Go loops to create multiple resources:

```go
for _, subnet := range locals.GcpVpc.Spec.Subnets {
    compute.NewSubnetwork(ctx, subnet.Name, ...)
}
```

### Cross-Stack References

Export outputs from one stack and import in another:

```go
// Stack A (VPC)
ctx.Export("network_self_link", createdNetwork.SelfLink)

// Stack B (Subnet)
vpcStack := pulumi.StackReference(...)
vpcSelfLink := vpcStack.GetOutput(pulumi.String("network_self_link"))
```

## Troubleshooting Tips for Developers

### Debugging Pulumi Programs

1. **Add logging**:
   ```go
   pulumi.Log.Info(ctx, fmt.Sprintf("Creating VPC: %s", locals.GcpVpc.Metadata.Name), nil)
   ```

2. **Use `pulumi preview --debug`**:
   ```bash
   pulumi preview --debug --logtostderr -v=9 2>&1 | tee debug.log
   ```

3. **Inspect state**:
   ```bash
   pulumi stack export > state.json
   jq '.deployment.resources' state.json
   ```

### Common Issues

**Issue**: "Resource already exists"
**Solution**: Import existing resource or delete it

**Issue**: "Credentials not found"
**Solution**: Run `gcloud auth application-default login`

**Issue**: "State locked"
**Solution**: Wait or run `pulumi cancel`

## Future Enhancements

Potential improvements for future versions:

1. **VPC Peering Support**: Add ability to peer VPCs via configuration
2. **Default Firewall Rules**: Optionally create allow-internal, deny-all rules
3. **VPC Flow Logs**: Enable flow logging for troubleshooting
4. **IPv6 Support**: Add dual-stack configuration
5. **Service Projects**: Automatically attach service projects to Shared VPC

## Conclusion

The Pulumi module for GCP VPC is designed to be simple, safe, and extensible. It abstracts away boilerplate while exposing the configuration that matters, following OpenMCF's philosophy of **guard rails, not handcuffs**.

For implementation details, see the source code. For usage instructions, see the [README](README.md).

