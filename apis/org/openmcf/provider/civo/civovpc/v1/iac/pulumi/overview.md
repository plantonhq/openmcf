# CivoVpc Pulumi Module - Architecture Overview

## Purpose

This Pulumi module provides Infrastructure as Code (IaC) for provisioning isolated private networks (VPCs) on the Civo cloud platform. It translates OpenMCF's protobuf-defined `CivoVpc` specification into actual Civo network resources.

## Design Philosophy

### 1. Protobuf-First API

The module is driven entirely by a protobuf specification (`CivoVpcSpec`), ensuring:

- **Type Safety**: Compile-time validation of configuration
- **Versioning**: Clear API versioning (v1, v2, etc.)
- **Language Agnostic**: Protobuf stubs can be generated for any language
- **Documentation as Code**: Proto comments serve as API documentation

### 2. Declarative Infrastructure

The module is fully declarative:

- **Input**: `CivoVpcStackInput` protobuf message
- **Process**: Idempotent resource creation/update
- **Output**: `CivoVpcStackOutputs` with network ID and CIDR

Running the same input multiple times produces the same result (Pulumi's state management ensures idempotency).

### 3. Minimal Configuration, Maximum Safety

Following the **80/20 principle**, the module exposes only essential fields:

- **Required**: `network_name`, `region`, `civo_credential_id`
- **Optional**: `ip_range_cidr` (auto-allocated if omitted)
- **Advanced**: `is_default_for_region`, `description`

Fields not applicable to public cloud VPCs (e.g., VLAN configuration for private cloud) are intentionally omitted to avoid confusion.

## Architecture

### Component Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                        OpenMCF                          │
│                      Orchestration Layer                        │
└────────────────────────────┬────────────────────────────────────┘
                             │ CivoVpcStackInput (protobuf)
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Pulumi Module (This Code)                     │
│                                                                 │
│  ┌──────────────┐      ┌──────────────┐      ┌──────────────┐ │
│  │   main.go    │─────▶│ module/      │─────▶│  vpc.go      │ │
│  │ (entrypoint) │      │ main.go      │      │ (resource)   │ │
│  └──────────────┘      └──────────────┘      └──────┬───────┘ │
│                             │                        │         │
│                             ▼                        │         │
│                      ┌──────────────┐                │         │
│                      │  locals.go   │                │         │
│                      │  (metadata)  │                │         │
│                      └──────────────┘                │         │
│                                                       │         │
└───────────────────────────────────────────────────────┼─────────┘
                                                        │
                             ┌──────────────────────────┘
                             │ Pulumi Civo Provider
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                          Civo Cloud                             │
│                        civo.Network                             │
│                    (Private Network / VPC)                      │
└─────────────────────────────────────────────────────────────────┘
```

### Data Flow

1. **Input Loading**: `main.go` loads `CivoVpcStackInput` from Pulumi context
2. **Locals Initialization**: `module/locals.go` extracts metadata and builds label map
3. **Provider Setup**: `module/main.go` creates Civo provider from credentials
4. **Resource Creation**: `module/vpc.go` provisions `civo.Network` resource
5. **Output Export**: Exports `network_id` and `cidr_block` to Pulumi state
6. **Status Update**: OpenMCF captures outputs and stores in resource status

## Key Decisions

### Why Pulumi (Go) Instead of Terraform?

**Rationale:**

1. **Language Consistency**: OpenMCF's core orchestration is written in Go
2. **Type Safety**: Strong typing at compile-time catches errors before deployment
3. **Programming Model**: Complex logic (conditionals, loops) is more natural in Go
4. **Testability**: Go's testing framework makes it easy to unit test IaC logic

**Trade-off:**
- Terraform has a larger ecosystem and broader adoption
- Pulumi's Civo provider is bridged from Terraform (not native)

**Result:**
- Both tools work equally well for standard VPC provisioning
- The choice is an implementation detail; the protobuf API remains the same

### Why Bridge from Terraform Provider?

The Pulumi Civo provider is **bridged** from the official Terraform provider:

**Benefits:**
- **1:1 Feature Parity**: Guaranteed coverage of all Civo resources
- **Automatic Updates**: New Terraform features flow to Pulumi
- **Reduced Maintenance**: No need to maintain separate provider code

**Limitations:**
- **Slight Lag**: Updates appear in Pulumi after Terraform release
- **Non-Native Feel**: Some APIs feel like Terraform HCL translated to Go

### CIDR Auto-Allocation Strategy

**Design Choice:** Support both explicit CIDR and auto-allocation

```go
// User specifies CIDR
networkArgs := &civo.NetworkArgs{
    Label:   pulumi.String("prod-network"),
    Region:  pulumi.String("LON1"),
    CidrV4:  pulumi.String("10.10.1.0/24"),  // Explicit
}

// User omits CIDR (Civo auto-allocates)
networkArgs := &civo.NetworkArgs{
    Label:  pulumi.String("dev-network"),
    Region: pulumi.String("LON1"),
    // CidrV4 omitted → Civo picks from available pool
}
```

**Rationale:**
- **Dev/Test**: Auto-allocation simplifies quick experiments
- **Production**: Explicit CIDRs enable predictable network planning and future VPN connectivity

**Implementation:**
```go
if locals.CivoVpc.Spec.IpRangeCidr != "" {
    networkArgs.CidrV4 = pulumi.String(locals.CivoVpc.Spec.IpRangeCidr)
}
```

### Handling Provider Limitations

The Pulumi Civo provider has some limitations compared to the full Civo API:

#### 1. `is_default_for_region` Not Supported

**Problem:** The `NetworkArgs` struct doesn't have a `Default` field in v2.4.8

**Solution:**
```go
if locals.CivoVpc.Spec.IsDefaultForRegion {
    ctx.Log.Warn("is_default_for_region is not supported by provider", nil)
}
```

**Workaround:** Users can set default via Civo CLI after creation

#### 2. `description` Field Not Exposed

**Problem:** Civo Network resource doesn't expose a description attribute

**Solution:**
```go
if locals.CivoVpc.Spec.Description != "" {
    ctx.Log.Info(fmt.Sprintf("Description: %s (not applied to resource)", desc), nil)
}
```

**Impact:** Description is stored in OpenMCF metadata but not on the Civo resource

#### 3. `created_at` Timestamp Not Available

**Problem:** Provider doesn't expose network creation timestamp

**Solution:**
- Define `created_at_rfc3339` in `CivoVpcStackOutputs` for API consistency
- Leave empty in actual output (document limitation)

**Alternative:** Could use `time.Now()` at creation time, but this would be Pulumi's creation time, not Civo's

## Resource Lifecycle

### Creation

```go
createdNetwork, err := civo.NewNetwork(
    ctx,
    "network",              // Pulumi resource name
    networkArgs,            // Configuration
    pulumi.Provider(civoProvider),  // Provider with credentials
)
```

**Behavior:**
- Creates new Civo network via API
- Returns `network.id` and `network.cidr_v4`
- Idempotent: Running again with same input updates (if possible) or no-ops

### Update

**Immutable Fields:**
- `region`: Cannot be changed (Civo constraint - networks are regional)
- `cidr_v4`: Cannot be changed after creation (Civo limitation)

**Mutable Fields:**
- `label`: Can be renamed

**Implication:**
- Changing `region` or `cidr_v4` forces **replacement** (destroy + recreate)
- Pulumi's plan will show this as a "replacement" operation

### Deletion

```go
pulumi destroy
```

**Behavior:**
- Deletes Civo network via API
- Resources (clusters, instances) must be deleted first (dependency constraint)

**Safety:**
- Pulumi tracks dependencies and will error if dependent resources exist
- Always destroy clusters/instances before destroying the network

## State Management

### Pulumi State

Pulumi tracks all created resources in state:

```json
{
  "resources": [
    {
      "type": "civo:index/network:Network",
      "urn": "urn:pulumi:stack::project::civo:index/network:Network::network",
      "id": "civo-network-uuid",
      "outputs": {
        "id": "civo-network-uuid",
        "label": "prod-network",
        "cidr_v4": "10.10.1.0/24",
        "region": "LON1"
      }
    }
  ]
}
```

**State Storage:**
- **Local Development**: `~/.pulumi/stacks/` (not recommended for teams)
- **Production**: Pulumi Cloud, S3, or **Civo Object Store** (self-contained)

### Recommended: Civo Object Store Backend

Use Civo's S3-compatible Object Store to create a fully self-contained stack:

```bash
pulumi login s3://my-civo-bucket?endpoint=https://objectstore.lon1.civo.com
```

**Benefits:**
- All Civo resources and state are in Civo
- No external dependencies (Pulumi Cloud, AWS S3)
- Cost-effective and performant

## Output Handling

### Exported Outputs

```go
ctx.Export(OpNetworkId, createdNetwork.ID())     // network_id
ctx.Export(OpCidrBlock, createdNetwork.CidrV4)   // cidr_block
```

### Consumption by Dependent Resources

Other resources (e.g., `CivoKubernetesCluster`) reference the `network_id`:

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoKubernetesCluster
spec:
  network_id: ${civovpc.outputs.network_id}  # Reference from stack output
```

Pulumi automatically builds a dependency graph:
1. Create network first
2. Create cluster second (after network exists)

## Labels and Metadata

### OpenMCF Labels

The module builds a standard label map for resource tracking:

```go
locals.CivoLabels = map[string]string{
    civolabelkeys.Resource:     "true",
    civolabelkeys.ResourceName: "prod-network",
    civolabelkeys.ResourceKind: "CivoVpc",
    civolabelkeys.Organization: "myorg",
    civolabelkeys.Environment:  "production",
    civolabelkeys.ResourceId:   "civpc-abc123",
}
```

**Current Limitation:**
- Civo Network resource doesn't support labels/tags via the provider
- These labels are stored in OpenMCF metadata only
- Future: If Civo adds label support, these can be applied to resources

## Error Handling

### Error Wrapping

The module uses `github.com/pkg/errors` for context-rich errors:

```go
if err != nil {
    return errors.Wrap(err, "failed to create civo network")
}
```

**Error Chain Example:**
```
failed to create vpc:
  failed to create civo network:
    API error: CIDR block 10.10.1.0/24 overlaps with existing network
```

### Common Errors

| Error | Cause | Solution |
|-------|-------|----------|
| "failed to setup civo provider" | Invalid credentials | Verify `civo_credential_id` |
| "CIDR block already in use" | Overlapping CIDR | Choose different CIDR or use auto-allocation |
| "region not found" | Invalid region code | Use valid region (LON1, NYC1, FRA1) |
| "replacement required" | Changed `region` or `cidr_v4` | Accept replacement or revert change |

## Testing Strategy

### Unit Testing

**Not yet implemented** (future work)

Potential approach:
```go
func TestLocalsInitialization(t *testing.T) {
    stackInput := &civovpcv1.CivoVpcStackInput{
        Target: &civovpcv1.CivoVpc{
            Metadata: &shared.CloudResourceMetadata{Name: "test"},
            Spec: &civovpcv1.CivoVpcSpec{NetworkName: "test"},
        },
    }
    locals := initializeLocals(nil, stackInput)
    assert.Equal(t, "test", locals.CivoVpc.Metadata.Name)
}
```

### Integration Testing

**Current approach:**
1. Create test stack input YAML
2. Run `pulumi preview` to verify plan
3. Run `pulumi up` to create resources
4. Verify outputs via Civo CLI: `civo network show <id>`
5. Run `pulumi destroy` to clean up

## Performance Considerations

### Resource Creation Time

- **Network creation**: ~2-5 seconds (Civo API is fast)
- **Full stack deployment**: Depends on dependent resources (clusters take longer)

### State Size

- Minimal: Single network resource produces ~1KB of state
- Scales linearly with number of resources

### API Rate Limits

- Civo API has rate limits (exact limits not documented)
- Pulumi respects rate limits via retries with exponential backoff
- For bulk operations, consider batching

## Security Considerations

### Credential Handling

- **Never hardcode API tokens** in code
- Credentials are managed by OpenMCF's credential store
- Pulumi provider receives credentials at runtime via environment variables

### Network Isolation

- Each Civo network is Layer 3 isolated (OpenStack-based)
- True tenant isolation (not pseudo-private)
- No cross-network communication without explicit routing

### CIDR Planning

- Plan CIDRs to avoid overlaps with on-premise networks
- Use RFC1918 private address space (10.0.0.0/8 recommended)
- Document allocation in version control

## Future Enhancements

### 1. Provider Feature Parity

**When Pulumi Civo provider adds support:**
- Apply `is_default_for_region` directly (currently logs warning)
- Apply `description` to resource (currently metadata-only)
- Expose `created_at` timestamp (currently unavailable)

### 2. Advanced Networking Features

**If Civo adds platform support:**
- VPC Peering (connect networks within same region)
- Inter-region connectivity (VPN or transit gateway)
- Network ACLs (beyond firewall rules)

### 3. Observability

**Potential additions:**
- Metrics: Network creation time, API latency
- Logging: Structured logs for all operations
- Tracing: Distributed tracing for multi-resource deployments

## Related Modules

- **CivoFirewall**: Creates firewall rules attached to VPCs
- **CivoKubernetesCluster**: Deploys K8s clusters within VPCs
- **CivoInstance**: Provisions compute instances within VPCs

## References

- **Pulumi Civo Provider**: [GitHub](https://github.com/pulumi/pulumi-civo)
- **Civo API Documentation**: [civo.com/api](https://www.civo.com/api)
- **Pulumi Go SDK**: [Pulumi Docs](https://www.pulumi.com/docs/languages-sdks/go/)

---

**Maintained by:** OpenMCF  
**Last Updated:** 2025-11-21

