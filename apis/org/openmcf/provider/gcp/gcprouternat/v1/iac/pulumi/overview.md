# GCP Cloud Router NAT - Pulumi Module Architecture

This document provides an architectural overview of the Pulumi module implementation for GCP Cloud Router NAT, explaining design decisions, resource creation flow, and implementation patterns.

## Architecture Overview

The Pulumi module follows a modular, declarative approach to provision Cloud Router and NAT resources in Google Cloud Platform. The architecture separates concerns into distinct layers:

```
┌─────────────────────────────────────────────────────────────────┐
│                    OpenMCF CLI                          │
│                   (User Interface Layer)                        │
└─────────────────────┬───────────────────────────────────────────┘
                      │
                      │ GcpRouterNatStackInput (protobuf)
                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                  iac/pulumi/main.go                             │
│                  (Pulumi Entrypoint)                            │
│  • Deserializes stack input                                    │
│  • Invokes module.Resources()                                  │
└─────────────────────┬───────────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────────┐
│              iac/pulumi/module/main.go                          │
│               (Module Entrypoint)                               │
│  • Initializes locals (labels, config)                         │
│  • Calls routerNat() to provision resources                    │
└─────────────────────┬───────────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────────┐
│         iac/pulumi/module/router_nat.go                         │
│          (Resource Provisioning Logic)                          │
│  1. Creates Cloud Router                                        │
│  2. Conditionally creates static IPs (if manual allocation)    │
│  3. Determines NAT IP allocation strategy                       │
│  4. Handles subnet coverage (all vs. specific)                  │
│  5. Configures logging                                          │
│  6. Creates Cloud NAT gateway                                   │
│  7. Exports outputs                                             │
└─────────────────────┬───────────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                  Google Cloud Platform                          │
│  • Cloud Router (regional)                                      │
│  • Static IPs (regional, conditional)                           │
│  • Cloud NAT (regional)                                         │
└─────────────────────────────────────────────────────────────────┘
```

## Module Structure

### 1. Entrypoint (`main.go`)

**Purpose:** Initializes Pulumi program and invokes module.

**Responsibilities:**
- Deserialize `GcpRouterNatStackInput` from environment variable or file
- Initialize Pulumi context
- Call `module.Resources()` to provision infrastructure
- Handle errors and return exit code

**Key Code:**
```go
func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        stackInput := &gcprouternatv1.GcpRouterNatStackInput{}
        // Deserialize stack input
        
        _, err := module.Resources(ctx, stackInput)
        return err
    })
}
```

### 2. Module Entrypoint (`module/main.go`)

**Purpose:** Orchestrates resource provisioning.

**Responsibilities:**
- Initialize locals (labels, derived values)
- Create GCP provider with credentials
- Invoke `routerNat()` to create resources
- Return created resources for potential chaining

**Key Code:**
```go
func Resources(ctx *pulumi.Context, stackInput *GcpRouterNatStackInput) (*compute.RouterNat, error) {
    locals := initializeLocals(ctx, stackInput)
    
    gcpProvider := createGcpProvider(ctx, locals.GcpProviderConfig)
    
    return routerNat(ctx, locals, gcpProvider)
}
```

### 3. Locals (`module/locals.go`)

**Purpose:** Compute derived values and standardized labels.

**Responsibilities:**
- Transform stack input into easily accessible struct
- Generate GCP labels following OpenMCF conventions
- Store provider configuration

**Label Convention:**
```go
labels := map[string]string{
    "resource":      "true",
    "resource-name": metadata.name,
    "resource-kind": "gcprouternat",
    "organization":  metadata.org,    // if provided
    "environment":   metadata.env,    // if provided
    "resource-id":   metadata.id,     // if provided
}
```

### 4. Resource Provisioning (`module/router_nat.go`)

**Purpose:** Core resource creation logic.

**Responsibilities:**
- Create Cloud Router
- Conditionally create static external IPs
- Determine NAT IP allocation strategy
- Handle subnet coverage modes
- Configure logging
- Create Cloud NAT gateway
- Export outputs

## Resource Creation Flow

### Step 1: Cloud Router Creation

**When:** Always (required for Cloud NAT).

**Purpose:** Regional router resource that hosts the NAT gateway.

**Implementation:**
```go
createdRouter, err := compute.NewRouter(ctx,
    "router",
    &compute.RouterArgs{
        Name:    pulumi.String(metadata.name),
        Region:  pulumi.String(spec.region),
        Network: pulumi.String(spec.vpc_self_link.GetValue()),
    },
    pulumi.Provider(gcpProvider))
```

**Outputs:**
- `router.SelfLink` → exported as `router_self_link`

### Step 2: Static IP Creation (Conditional)

**When:** If `spec.nat_ip_names` is non-empty (manual allocation mode).

**Purpose:** Reserve static external IP addresses for deterministic egress.

**Implementation:**
```go
if len(spec.nat_ip_names) > 0 {
    for idx, natIpName := range spec.nat_ip_names {
        createdNatIp, err := compute.NewAddress(ctx,
            fmt.Sprintf("nat-ip-%d", idx),
            &compute.AddressArgs{
                Name:        pulumi.String(natIpName.GetValue()),
                Region:      pulumi.String(spec.region),
                AddressType: pulumi.String("EXTERNAL"),
                Labels:      pulumi.ToStringMap(locals.GcpLabels),
            },
            pulumi.Provider(gcpProvider),
            pulumi.Parent(createdRouter))
        
        natIps = append(natIps, createdNatIp.SelfLink)
    }
}
```

**Decision Logic:**
- **Empty list** → Auto-allocation mode (no IPs created)
- **Non-empty list** → Manual allocation mode (create static IPs)

### Step 3: NAT IP Allocation Strategy

**Purpose:** Determine how Cloud NAT allocates external IPs.

**Implementation:**
```go
natIpAllocateOption := pulumi.String("AUTO_ONLY")
var natIps pulumi.StringArray

if len(spec.nat_ip_names) > 0 {
    natIpAllocateOption = pulumi.String("MANUAL_ONLY")
    natIps = [...] // Self-links of created static IPs
}
```

**Modes:**
- **AUTO_ONLY**: Google automatically allocates and scales IPs (default, recommended).
- **MANUAL_ONLY**: Use specified static IPs (for partner allowlisting).

### Step 4: Subnet Coverage Strategy

**Purpose:** Determine which subnets use NAT for egress.

**Implementation:**
```go
subnetworks := compute.RouterNatSubnetworkArray{}
sourceRangeSetting := pulumi.String("ALL_SUBNETWORKS_ALL_IP_RANGES")

if len(spec.subnetwork_self_links) > 0 {
    sourceRangeSetting = pulumi.String("LIST_OF_SUBNETWORKS")
    
    for _, subnet := range spec.subnetwork_self_links {
        subnetworks = append(subnetworks, &compute.RouterNatSubnetworkArgs{
            Name:                  pulumi.String(subnet.GetValue()),
            SourceIpRangesToNats:  pulumi.StringArray{pulumi.String("ALL_IP_RANGES")},
            SecondaryIpRangeNames: pulumi.StringArray{},
        })
    }
}
```

**Modes:**
- **ALL_SUBNETWORKS_ALL_IP_RANGES**: All subnets in the region (default, future-proof).
- **LIST_OF_SUBNETWORKS**: Only specified subnets (fine-grained control).

### Step 5: Logging Configuration

**Purpose:** Configure NAT translation logging for troubleshooting and auditing.

**Implementation:**
```go
var logConfig *compute.RouterNatLogConfigArgs

switch spec.log_filter {
case gcprouternatv1.GcpRouterNatLogFilter_ERRORS_ONLY:
    logConfig = &compute.RouterNatLogConfigArgs{
        Enable: pulumi.Bool(true),
        Filter: pulumi.String("ERRORS_ONLY"),
    }
case gcprouternatv1.GcpRouterNatLogFilter_ALL:
    logConfig = &compute.RouterNatLogConfigArgs{
        Enable: pulumi.Bool(true),
        Filter: pulumi.String("ALL"),
    }
case gcprouternatv1.GcpRouterNatLogFilter_DISABLED:
    logConfig = &compute.RouterNatLogConfigArgs{
        Enable: pulumi.Bool(false),
    }
default:
    // Default to ERRORS_ONLY
    logConfig = &compute.RouterNatLogConfigArgs{
        Enable: pulumi.Bool(true),
        Filter: pulumi.String("ERRORS_ONLY"),
    }
}
```

**Logging Modes:**
- **DISABLED**: No logging (use for non-production to reduce costs).
- **ERRORS_ONLY**: Log translation errors only (recommended for production).
- **ALL**: Log all translations (security auditing, high log volume).

### Step 6: Cloud NAT Creation

**Purpose:** Create the NAT gateway with all computed configuration.

**Implementation:**
```go
createdRouterNat, err := compute.NewRouterNat(ctx,
    "router-nat",
    &compute.RouterNatArgs{
        Name:                          pulumi.String(metadata.name),
        Router:                        createdRouter.Name,
        Region:                        pulumi.String(spec.region),
        NatIpAllocateOption:           natIpAllocateOption,
        NatIps:                        natIps,
        SourceSubnetworkIpRangesToNat: sourceRangeSetting,
        Subnetworks:                   subnetworks,
        LogConfig:                     logConfig,
    },
    pulumi.Provider(gcpProvider),
    pulumi.Parent(createdRouter))
```

**Key Parameters:**
- `Router`: References created Cloud Router (dependency).
- `NatIpAllocateOption`: AUTO_ONLY or MANUAL_ONLY.
- `NatIps`: Static IP self-links (only if MANUAL_ONLY).
- `SourceSubnetworkIpRangesToNat`: ALL or LIST_OF_SUBNETWORKS.
- `Subnetworks`: Subnet configuration (only if LIST_OF_SUBNETWORKS).
- `LogConfig`: Logging configuration.

### Step 7: Output Export

**Purpose:** Export resource outputs for OpenMCF status tracking.

**Implementation:**
```go
ctx.Export("name", createdRouterNat.Name)
ctx.Export("router_self_link", createdRouter.SelfLink)
ctx.Export("nat_ip_addresses", createdRouterNat.NatIps)
```

**Outputs:**
- `name`: Cloud NAT gateway name (for status display).
- `router_self_link`: Cloud Router self-link (for resource references).
- `nat_ip_addresses`: List of NAT IPs (auto-allocated or static).

## Design Decisions

### 1. Why Separate `router_nat.go` from `main.go`?

**Reason:** Separation of concerns.

- `main.go`: Orchestration (locals, provider initialization, coordination).
- `router_nat.go`: Resource provisioning logic (focused on Cloud Router and NAT).

**Benefit:** Easier to test, maintain, and extend. If additional resources are needed (firewall rules, monitoring), they can be added as separate files.

### 2. Why Compute Locals Instead of Passing Raw Spec?

**Reason:** Simplifies resource logic and provides reusable derived values.

**Example:**
- Raw: `stackInput.Target.Metadata.Name` (verbose, error-prone).
- Locals: `locals.GcpRouterNat.Metadata.Name` (clear, type-safe).

**Benefit:** Reduces duplication and makes code more readable.

### 3. Why Use `pulumi.Parent()` for Static IPs?

**Reason:** Establishes dependency hierarchy.

**Behavior:**
- Static IPs are logically children of the Cloud Router.
- If Cloud Router is deleted, Pulumi automatically deletes child IPs.

**Benefit:** Prevents orphaned resources and simplifies cleanup.

### 4. Why Default to ERRORS_ONLY Logging?

**Reason:** Production-ready balance between visibility and cost.

**Rationale:**
- **DISABLED**: Saves costs but hides critical errors (port exhaustion, connection failures).
- **ERRORS_ONLY**: Detects issues without excessive log volume (recommended).
- **ALL**: Security auditing but generates massive log volume (use sparingly).

**Benefit:** Users get sensible defaults without manual configuration.

### 5. Why Conditional Resource Creation?

**Reason:** Avoid creating resources that aren't needed.

**Example:**
- If `nat_ip_names` is empty, don't create static IPs.
- If `subnetwork_self_links` is empty, don't build subnetwork configuration.

**Benefit:** Cleaner state, lower costs, faster provisioning.

## Error Handling

The module uses wrapped errors for context:

```go
if err != nil {
    return nil, errors.Wrap(err, "failed to create router")
}
```

**Error Messages:**
- `"failed to create router"`: Cloud Router creation failed.
- `"failed to create static nat ip"`: Static IP creation failed.
- `"failed to create router nat"`: Cloud NAT creation failed.

**Propagation:** Errors bubble up to the Pulumi entrypoint, which logs them and returns a non-zero exit code.

## Dependency Management

Pulumi automatically handles resource dependencies through implicit references:

```go
&compute.RouterNatArgs{
    Router: createdRouter.Name,  // Implicit dependency on createdRouter
    NatIps: natIps,               // Implicit dependency on created IPs
}
```

**Execution Order:**
1. Cloud Router (no dependencies)
2. Static IPs (parent: Router)
3. Cloud NAT (depends on Router and IPs)

**Benefit:** No manual dependency ordering required—Pulumi's DAG handles it.

## Testing Strategy

### Unit Testing

**Scope:** Protobuf validation rules in `spec_test.go`.

**Example:**
```go
It("should not return a validation error for minimal valid fields", func() {
    spec := &GcpRouterNatSpec{
        VpcSelfLink: &StringValueOrRef{Value: "https://..."},
        Region:      "us-central1",
    }
    err := protovalidate.Validate(spec)
    Expect(err).To(BeNil())
})
```

### Integration Testing

**Scope:** End-to-end deployment via OpenMCF CLI.

**Process:**
1. Create test manifest (`hack/manifest.yaml`)
2. Deploy: `planton apply -f hack/manifest.yaml`
3. Verify outputs: `planton status gcprouternat test-nat`
4. Cleanup: `planton destroy gcprouternat test-nat`

## Performance Considerations

### Parallelism

Pulumi provisions resources in parallel when dependencies allow:

- Cloud Router (independent) → starts immediately
- Static IPs (depend on Router) → start after Router
- Cloud NAT (depends on Router and IPs) → starts after all dependencies

**Benefit:** Faster provisioning compared to sequential execution.

### State Management

Pulumi maintains state in backend (e.g., S3, GCS, Pulumi Service):

- Tracks created resources
- Detects configuration drift
- Enables `pulumi refresh` to sync state

**Benefit:** Reliable infrastructure management and change detection.

## Best Practices Applied

1. **Sensible Defaults**: Auto-allocation, all-subnets coverage, ERRORS_ONLY logging.
2. **Conditional Resources**: Only create what's needed (static IPs, subnetwork config).
3. **Dependency Hierarchy**: Use `pulumi.Parent()` for logical resource relationships.
4. **Error Context**: Wrap errors with descriptive messages.
5. **Standardized Labels**: Apply consistent labels following OpenMCF conventions.
6. **Protobuf-Driven**: Configuration driven by strongly-typed protobuf schema.

## Future Enhancements

Potential improvements (not currently implemented):

1. **Port Allocation Configuration**: Expose `minPortsPerVm`, `maxPortsPerVm`, `dynamicPortAllocation`.
2. **Custom Timeouts**: Allow configuring TCP/UDP/ICMP timeouts.
3. **NAT Rules**: Support advanced NAT rules for destination-specific routing.
4. **Endpoint-Independent Mapping**: Make configurable (currently always enabled).
5. **Monitoring**: Auto-create Cloud Monitoring dashboards and alerts.

## Troubleshooting

### Issue: "Network not found"
**Cause:** VPC self-link is incorrect.
**Debug:** Check `spec.vpc_self_link` value, verify VPC exists.

### Issue: "IP address already exists"
**Cause:** Static IP name collision.
**Debug:** Check existing IPs: `gcloud compute addresses list --region=<region>`.

### Issue: "Region mismatch"
**Cause:** Static IPs and NAT in different regions.
**Debug:** Verify all resources use the same `spec.region`.

### Issue: "Pulumi stack drift detected"
**Cause:** Manual changes via Console or gcloud CLI.
**Solution:** Run `pulumi refresh` to sync state, then `pulumi up` to reconcile.

## Related Documentation

- **API Specification**: [../../spec.proto](../../spec.proto)
- **Examples**: [../../examples.md](../../examples.md)
- **User Guide**: [../../README.md](../../README.md)
- **Pulumi Usage**: [README.md](README.md)
- **Terraform Module**: [../../iac/tf/README.md](../../iac/tf/README.md)

---

This architecture document provides the foundation for understanding, maintaining, and extending the GCP Cloud Router NAT Pulumi module. For usage examples, see [README.md](README.md).

