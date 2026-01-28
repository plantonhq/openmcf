# DigitalOcean VPC - Pulumi Module Architecture

This document provides a technical overview of the Pulumi module implementation for DigitalOcean Virtual Private Cloud (VPC) networks.

## Architecture Overview

```
┌─────────────────────────────────────────────────┐
│          Pulumi Entrypoint (main.go)            │
│  - Loads stack input (manifest YAML)           │
│  - Initializes Pulumi context                   │
│  - Delegates to module                          │
└──────────────────┬──────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────┐
│          Module Layer (module/main.go)           │
│  - Processes DigitalOceanVpc spec               │
│  - Creates DigitalOcean provider                │
│  - Orchestrates VPC resource creation           │
└──────────────────┬──────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────┐
│       Resource Implementation                    │
│  ┌──────────┐  ┌─────────────┐                 │
│  │ vpc.go   │  │ locals.go   │                 │
│  │ - VPC    │  │ - Variables │                 │
│  │          │  │ - Labels    │                 │
│  └──────────┘  └─────────────┘                 │
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
├── main.go           # Pulumi program entrypoint
├── Pulumi.yaml       # Project configuration
├── Makefile          # Helper commands
├── debug.sh          # Debug script
├── module/           # Core implementation
│   ├── main.go       # Module orchestration
│   ├── vpc.go        # VPC resource logic
│   ├── locals.go     # Local variables
│   └── outputs.go    # Output constants
├── README.md         # Usage documentation
└── overview.md       # This file
```

## Implementation Details

### 1. Entrypoint (`main.go`)

**Responsibility:** Initialize Pulumi context and load stack input

```go
func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        stackInput := &digitaloceanvpcv1.DigitalOceanVpcStackInput{}
        
        // Load from Pulumi config or stack exports
        if err := loadStackInput(ctx, stackInput); err != nil {
            return err
        }
        
        // Delegate to module
        return module.Resources(ctx, stackInput)
    })
}
```

### 2. Module Orchestration (`module/main.go`)

**Responsibility:** Coordinate VPC provisioning

```go
func Resources(ctx *pulumi.Context, stackInput *StackInput) error {
    // 1. Initialize locals
    locals := initializeLocals(ctx, stackInput)
    
    // 2. Create DigitalOcean provider
    provider, err := createDigitalOceanProvider(ctx, locals)
    if err != nil {
        return err
    }
    
    // 3. Create VPC resource
    _, err = vpc(ctx, locals, provider)
    if err != nil {
        return err
    }
    
    return nil
}
```

### 3. VPC Resource (`module/vpc.go`)

**Responsibility:** Provision DigitalOcean VPC with conditional IP range

#### Core Logic

```go
func vpc(ctx *pulumi.Context, locals *Locals, provider *digitalocean.Provider) (*digitalocean.Vpc, error) {
    // 1. Build base arguments
    vpcArgs := &digitalocean.VpcArgs{
        Name:   pulumi.String(locals.DigitalOceanVpc.Metadata.Name),
        Region: pulumi.String(locals.DigitalOceanVpc.Spec.Region.String()),
    }
    
    // 2. Add optional description
    if locals.DigitalOceanVpc.Spec.Description != "" {
        vpcArgs.Description = pulumi.String(locals.DigitalOceanVpc.Spec.Description)
    }
    
    // 3. Add IP range if specified (80/20 principle)
    // When omitted, DigitalOcean auto-generates a /20 CIDR
    if locals.DigitalOceanVpc.Spec.IpRangeCidr != "" {
        vpcArgs.IpRange = pulumi.String(locals.DigitalOceanVpc.Spec.IpRangeCidr)
    }
    
    // 4. Create VPC
    createdVpc, err := digitalocean.NewVpc(ctx, "vpc", vpcArgs,
        pulumi.Provider(provider))
    if err != nil {
        return nil, errors.Wrap(err, "failed to create vpc")
    }
    
    // 5. Export outputs
    ctx.Export(OpVpcId, createdVpc.ID())
    
    return createdVpc, nil
}
```

#### Key Features

**1. Conditional IP Range (80/20 Principle)**

```go
// Only include IP range if user specified it
if spec.IpRangeCidr != "" {
    vpcArgs.IpRange = pulumi.String(spec.IpRangeCidr)
}
// Otherwise, DO auto-generates /20 block
```

**2. Optional Description**

```go
if spec.Description != "" {
    vpcArgs.Description = pulumi.String(spec.Description)
}
```

**3. Region Enum Conversion**

```go
Region: pulumi.String(spec.Region.String())
```

Converts proto enum (e.g., `NYC3`) to DigitalOcean region slug (`"nyc3"`).

### 4. Local Variables (`module/locals.go`)

**Responsibility:** Process metadata and generate labels

```go
type Locals struct {
    DigitalOceanVpc            *DigitalOceanVpc
    DigitalOceanProviderConfig *ProviderConfig
    DigitalOceanLabels         map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *StackInput) *Locals {
    locals := &Locals{
        DigitalOceanVpc: stackInput.Target,
    }
    
    // Generate standard OpenMCF labels
    locals.DigitalOceanLabels = map[string]string{
        "planton-resource":      "true",
        "planton-resource-name": locals.DigitalOceanVpc.Metadata.Name,
        "planton-resource-kind": "DigitalOceanVpc",
    }
    
    // Add optional metadata labels
    if locals.DigitalOceanVpc.Metadata.Org != "" {
        locals.DigitalOceanLabels["planton-organization"] = locals.Metadata.Org
    }
    
    return locals
}
```

### 5. Outputs (`module/outputs.go`)

**Responsibility:** Define output constants

```go
const (
    OpVpcId = "vpc_id"
)
```

**Exported Values:**
- `vpc_id`: VPC UUID for cross-stack references

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
module/vpc.go
  │
  ├─ Build VPC args
  ├─ Conditionally add IP range
  ├─ Conditionally add description
  ├─ Create VPC resource
  ├─ Export VPC ID
  │
  ▼
DigitalOcean API
  │
  ├─ Validate CIDR (if provided)
  ├─ Auto-generate CIDR (if omitted)
  ├─ Create VPC
  │
  ▼
Complete (~30 seconds)
```

## Configuration Processing

### Spec Validation

Validation occurs at protobuf level before Pulumi execution:

```proto
message DigitalOceanVpcSpec {
  DigitalOceanRegion region = 2 [(buf.validate.field).required = true];
  
  string ip_range_cidr = 3 [
    // Optional - supports 80/20 principle
    (buf.validate.field).string.pattern = "^([0-9]{1,3}\\.){3}[0-9]{1,3}/(16|20|24)$"
  ];
}
```

**Required:**
- `region` - Must be valid DigitalOcean region

**Optional but Validated:**
- `ip_range_cidr` - If provided, must match /16, /20, or /24 pattern

### Field Transformations

**Region Enum → String:**
```go
pulumi.String(spec.Region.String())
```

**Empty String → nil (Auto-Generation):**
```go
if spec.IpRangeCidr != "" {
    vpcArgs.IpRange = pulumi.String(spec.IpRangeCidr)
}
// Otherwise, field is omitted → DigitalOcean auto-generates
```

## Design Decisions

### 1. Optional IP Range (80/20 Principle)

**Decision:** Only include `IpRange` in API call if user specified it

**Rationale:**
- 80% of users want auto-generated CIDR
- DigitalOcean API interprets missing `ip_range` as "auto-generate /20"
- Explicit nil or empty string triggers auto-generation
- Matches the platform's native behavior

**Implementation:**
```go
if spec.IpRangeCidr != "" {
    vpcArgs.IpRange = pulumi.String(spec.IpRangeCidr)
}
```

**Tradeoff:** Requires conditional logic, but enables the 80% use case.

### 2. Immutability Awareness

**Decision:** VPC IP ranges are immutable

**Rationale:**
- DigitalOcean doesn't support resizing VPCs
- Terraform lifecycle ignore_changes prevents drift
- Pulumi doesn't need explicit ignore (updates fail naturally)

**Impact:** Users must plan CIDR carefully before deployment.

### 3. Name from Metadata

**Decision:** Use `metadata.name` for VPC name, not `spec.name`

**Rationale:**
- Consistent with OpenMCF patterns
- Enables resource naming conventions
- Supports multi-environment naming (e.g., `{env}-vpc`)

### 4. No Default VPC Handling

**Decision:** Don't implement `is_default_for_region` as input

**Rationale:**
- DigitalOcean automatically sets first VPC in region as default
- The Ansible `default: true` pattern is a non-atomic abstraction
- OpenMCF avoids abstractions that don't match platform behavior

**Future:** Could add as a separate update operation if needed.

## Error Handling

```go
createdVpc, err := digitalocean.NewVpc(ctx, "vpc", vpcArgs)
if err != nil {
    return nil, errors.Wrap(err, "failed to create digitalocean vpc")
}
```

**Pattern:**
1. Wrap errors with context
2. Return immediately on failure
3. Propagate to Pulumi runtime

**Common Errors:**
- `overlapping CIDR`: IP range conflicts with existing VPC
- `invalid region`: Region doesn't exist
- `invalid CIDR format`: Not /16, /20, or /24

## Production Patterns

### Auto-Generated CIDR (Dev/Test)

```yaml
spec:
  region: nyc3
  # No ip_range_cidr → auto-generate
```

**Benefits:**
- Zero planning overhead
- DigitalOcean ensures no conflicts
- Perfect for ephemeral environments

### Explicit CIDR (Production)

```yaml
spec:
  region: nyc1
  ipRangeCidr: "10.101.0.0/16"
  description: "Production VPC"
```

**Benefits:**
- Full control over IP allocation
- Documented in IPAM
- Non-overlapping for VPC peering

### Multi-Environment Strategy

```
Dev:     10.100.0.0/20   (auto-generated or explicit)
Staging: 10.100.16.0/20  (explicit)
Prod:    10.101.0.0/16   (explicit /16 for scale)
```

## Customization Points

### Adding Additional Outputs

To export more VPC properties:

```go
ctx.Export("vpc_urn", createdVpc.Urn)
ctx.Export("ip_range", createdVpc.IpRange)
ctx.Export("is_default", createdVpc.Default)
ctx.Export("created_at", createdVpc.CreatedAt)
```

### Cross-Stack References

Use VPC ID in other stacks:

```go
// In another Pulumi program
vpcRef := pulumi.NewStackReference(ctx, "org/prod-vpc/prod", nil)
vpcId := vpcRef.GetOutput(pulumi.String("vpc_id"))

// Use in cluster
cluster := digitalocean.NewKubernetesCluster(ctx, "cluster", &Args{
    VpcUuid: vpcId,
})
```

### Post-Creation Actions

Execute actions after VPC creation:

```go
createdVpc.ID().ApplyT(func(id string) error {
    // Configure firewall rules
    // Set up VPC peering
    // Tag VPC resources
    return nil
})
```

## Performance Considerations

- **Deployment time:** ~30 seconds (very fast)
- **State size:** Minimal (~1 KB per VPC)
- **API calls:** 1-2 per deployment

## VPC Constraints

### Immutability

**IP ranges cannot be changed:**
- Once created with CIDR, it's permanent
- Over-provision to avoid future migrations
- Plan multi-environment ranges carefully

### VPC-First Requirement

**These resources MUST be created in VPC from day one:**
- DOKS clusters
- Load balancers

**These can be migrated (with downtime):**
- Droplets (via snapshots)
- Managed databases (via control panel)

### Regional Limits

- One VPC is regional (cannot span regions)
- First VPC in region becomes default automatically
- 10,000 resources per VPC maximum

## Testing

### Integration Tests

```bash
# Deploy test VPC
pulumi stack init test
pulumi up --yes

# Verify
VPC_ID=$(pulumi stack output vpc_id)
doctl vpcs get $VPC_ID

# Verify auto-generated CIDR
doctl vpcs get $VPC_ID --format IPRange

# Cleanup
pulumi destroy --yes
```

### Validation Tests

```bash
# Test with explicit CIDR
pulumi config set vpc-cidr "10.100.0.0/20"
pulumi up

# Test with auto-generated CIDR
pulumi config set vpc-cidr ""
pulumi up
```

## 80/20 Principle Implementation

### The 80% Use Case: Auto-Generated CIDR

**What it means:**
- Most users don't want to manage IP planning
- DigitalOcean can auto-generate non-conflicting /20 blocks
- This is the recommended approach for dev/test

**How we support it:**
```go
// Don't include ip_range in args if not specified
if spec.IpRangeCidr != "" {
    vpcArgs.IpRange = pulumi.String(spec.IpRangeCidr)
}
// Field omission triggers auto-generation
```

### The 20% Use Case: Explicit CIDR

**What it means:**
- Production environments need IP planning
- Corporate IPAM requirements
- Multi-region architectures need non-overlapping ranges
- VPC peering requires careful planning

**How we support it:**
```go
// Include explicit IP range when user specifies it
vpcArgs.IpRange = pulumi.String(spec.IpRangeCidr)
```

## Best Practices

### 1. Use Stacks for Environments

```bash
pulumi stack init dev      # Auto-generated CIDR
pulumi stack init staging  # Explicit /20
pulumi stack init prod     # Explicit /16
```

### 2. Stack Configuration

```bash
# Dev: Let it auto-generate
pulumi config set region nyc3

# Prod: Explicit CIDR
pulumi config set region nyc1
pulumi config set cidr "10.101.0.0/16"
```

### 3. Export VPC ID for Reuse

```bash
# Export from VPC stack
pulumi stack output vpc_id > ../cluster/vpc-id.txt

# Or use stack references in code
```

### 4. Document Auto-Generated Ranges

```bash
# After deployment, document what was auto-generated
VPC_CIDR=$(pulumi stack output ip_range)
echo "Dev VPC auto-generated: $VPC_CIDR" >> IPAM.md
```

## Common Patterns

### Pattern 1: VPC + Cluster (Same Stack)

```go
// Create VPC first
vpc, err := createVpc(ctx, vpcSpec)

// Use VPC ID in cluster
cluster, err := createCluster(ctx, &ClusterArgs{
    VpcUuid: vpc.ID(),
})
```

### Pattern 2: VPC + Cluster (Separate Stacks)

**VPC Stack:**
```go
vpc, err := createVpc(ctx, vpcSpec)
ctx.Export("vpc_id", vpc.ID())
```

**Cluster Stack:**
```go
vpcStack := pulumi.NewStackReference(ctx, "org/prod-vpc/prod")
vpcId := vpcStack.GetOutput(pulumi.String("vpc_id"))

cluster, err := createCluster(ctx, &ClusterArgs{
    VpcUuid: vpcId,
})
```

## Reference

- [Pulumi DigitalOcean Provider Source](https://github.com/pulumi/pulumi-digitalocean)
- [DigitalOcean VPC API](https://docs.digitalocean.com/reference/api/api-reference/#tag/VPCs)
- [OpenMCF Architecture](../../../../../../architecture/README.md)

---

**Version:** 1.0.0  
**Last Updated:** 2025-11-16  
**Authors:** OpenMCF Engineering Team

