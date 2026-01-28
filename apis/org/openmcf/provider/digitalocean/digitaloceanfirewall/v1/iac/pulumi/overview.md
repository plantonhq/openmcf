# DigitalOcean Firewall - Pulumi Module Architecture

## Purpose

This document explains the internal architecture, design decisions, and implementation patterns of the DigitalOcean Firewall Pulumi module.

**Target Audience**: Contributors, maintainers, and engineers extending or debugging the module.

For **usage instructions**, see [README.md](./README.md).

---

## Architecture Overview

The module follows OpenMCF's standard Pulumi module pattern:

```
┌─────────────────────────────────────────────────────────────┐
│                         main.go                              │
│  (Pulumi entrypoint - calls module.Resources)               │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    module/main.go                            │
│  • Resources() - entry point                                 │
│  • Orchestrates: locals → provider → resources → outputs    │
└─────────────────────────────────────────────────────────────┘
                              │
          ┌───────────────────┼───────────────────┐
          ▼                   ▼                   ▼
    ┌──────────┐      ┌──────────────┐     ┌──────────┐
    │ locals.go│      │ firewall.go  │     │outputs.go│
    │ (context)│      │ (resources)  │     │(exports) │
    └──────────┘      └──────────────┘     └──────────┘
```

### Key Files

1. **`main.go`** (root)
   - Pulumi program entrypoint
   - Parses stack input (protobuf)
   - Calls `module.Resources(ctx, stackInput)`

2. **`module/main.go`**
   - Module orchestration
   - Initializes locals, provider, resources
   - Follows OpenMCF's standard pattern

3. **`module/locals.go`**
   - Initializes locals (metadata, labels, spec)
   - Transforms protobuf spec into module-usable format

4. **`module/firewall.go`**
   - Core resource creation logic
   - Translates protobuf rules to DigitalOcean Pulumi format
   - Exports stack outputs

5. **`module/outputs.go`**
   - Defines output constant names
   - Centralizes output keys for consistency

---

## Data Flow

### 1. Input: DigitalOceanFirewallStackInput

```protobuf
message DigitalOceanFirewallStackInput {
  ProviderConfig provider_config = 1;          // DigitalOcean credentials
  DigitalOceanFirewall target = 2;             // Firewall spec
}

message DigitalOceanFirewall {
  CloudResourceMetadata metadata = 1;          // Name, ID, env, org
  DigitalOceanFirewallSpec spec = 2;           // Firewall configuration
}
```

### 2. Processing: Locals Initialization

```go
// In locals.go
type Locals struct {
    DigitalOceanFirewall *digitaloceanfirewallv1.DigitalOceanFirewall
    Labels                map[string]string  // Derived from metadata
}

func initializeLocals(ctx *pulumi.Context, stackInput *stackInput) *Locals {
    // Extract spec, metadata
    // Transform into Pulumi-usable format
}
```

### 3. Resource Creation: Firewall

```go
// In firewall.go
func firewall(
    ctx *pulumi.Context,
    locals *Locals,
    digitalOceanProvider *digitalocean.Provider,
) (*digitalocean.Firewall, error) {
    // 1. Translate inbound rules
    inboundRules := make(digitalocean.FirewallInboundRuleArray, 0)
    for _, rule := range locals.DigitalOceanFirewall.Spec.InboundRules {
        inboundRules = append(inboundRules, digitalocean.FirewallInboundRuleArgs{
            Protocol:               pulumi.String(rule.Protocol),
            PortRange:              pulumi.String(rule.PortRange),
            SourceAddresses:        pulumi.ToStringArray(rule.SourceAddresses),
            SourceTags:             pulumi.ToStringArray(rule.SourceTags),
            SourceLoadBalancerUids: pulumi.ToStringArray(rule.SourceLoadBalancerUids),
            // ... other source types
        })
    }

    // 2. Translate outbound rules (similar pattern)

    // 3. Create firewall resource
    createdFirewall, err := digitalocean.NewFirewall(ctx, "firewall", &digitalocean.FirewallArgs{
        Name:          pulumi.String(locals.DigitalOceanFirewall.Spec.Name),
        InboundRules:  inboundRules,
        OutboundRules: outboundRules,
        DropletIds:    int64SliceToIntArray(locals.DigitalOceanFirewall.Spec.DropletIds),
        Tags:          pulumi.ToStringArray(locals.DigitalOceanFirewall.Spec.Tags),
    }, pulumi.Provider(digitalOceanProvider))

    // 4. Export output
    ctx.Export(OpFirewallId, createdFirewall.ID())

    return createdFirewall, nil
}
```

### 4. Output: Stack Outputs

```go
// In outputs.go
const (
    OpFirewallId = "firewall_id"  // DigitalOcean firewall UUID
)
```

---

## Design Decisions

### 1. Rule Translation Pattern

**Challenge**: Protobuf `InboundRule` and `OutboundRule` messages have multiple optional source/destination fields. Pulumi DigitalOcean provider expects `FirewallInboundRuleArgs`.

**Solution**: Iterate over protobuf rules and translate each to Pulumi args:

```go
inboundRules := make(digitalocean.FirewallInboundRuleArray, 0, len(spec.InboundRules))
for _, rule := range spec.InboundRules {
    inboundRules = append(inboundRules, digitalocean.FirewallInboundRuleArgs{
        Protocol:               pulumi.String(rule.Protocol),
        PortRange:              pulumi.String(rule.PortRange),
        SourceAddresses:        pulumi.ToStringArray(rule.SourceAddresses),
        SourceDropletIds:       int64SliceToIntArray(rule.SourceDropletIds),
        SourceTags:             pulumi.ToStringArray(rule.SourceTags),
        SourceKubernetesIds:    pulumi.ToStringArray(rule.SourceKubernetesIds),
        SourceLoadBalancerUids: pulumi.ToStringArray(rule.SourceLoadBalancerUids),
    })
}
```

**Why this pattern**:
- ✅ Type-safe: Pulumi validates types at compile-time
- ✅ Explicit: No magic conversions, clear mapping
- ✅ Maintainable: Easy to add new rule types

### 2. int64 to Int Conversion

**Challenge**: Protobuf uses `int64` for Droplet IDs. Pulumi DigitalOcean provider expects `pulumi.Int`.

**Solution**: Helper function `int64SliceToIntArray()`:

```go
func int64SliceToIntArray(values []int64) pulumi.IntArray {
    intInputs := make(pulumi.IntArray, 0, len(values))
    for _, v := range values {
        intInputs = append(intInputs, pulumi.Int(int(v)))
    }
    return intInputs
}
```

**Why**: Pulumi type system requires explicit conversion.

### 3. Provider Instantiation

**Pattern**: Provider is instantiated in `module/main.go` and passed to resource functions:

```go
digitalOceanProvider, err := pulumidigitaloceanprovider.Get(ctx, stackInput.ProviderConfig)
if err != nil {
    return errors.Wrap(err, "failed to setup digitalocean provider")
}

createdFirewall, err := firewall(ctx, locals, digitalOceanProvider)
```

**Why**:
- ✅ Centralized provider management
- ✅ Supports credential switching (multi-account)
- ✅ Testable (mock provider in tests)

### 4. Output Naming

**Pattern**: Output keys are defined as constants in `outputs.go`:

```go
const OpFirewallId = "firewall_id"
```

**Why**:
- ✅ Prevents typos (compile-time check)
- ✅ Consistent across components
- ✅ Easy to reference in tests

---

## Resource Graph

For a production multi-tier architecture, the Pulumi resource graph looks like:

```
┌──────────────────────────────────────────────────────────┐
│              DigitalOcean Provider                        │
│  (Configured with API token from ProviderConfig)         │
└──────────────────────────────────────────────────────────┘
                        │
        ┌───────────────┼───────────────┬──────────────────┐
        ▼               ▼               ▼                  ▼
  ┌───────────┐   ┌───────────┐   ┌───────────┐   ┌───────────┐
  │   Web     │   │   DB      │   │  Cache    │   │Management │
  │ Firewall  │   │ Firewall  │   │ Firewall  │   │ Firewall  │
  └───────────┘   └───────────┘   └───────────┘   └───────────┘
        │               │               │                  │
        │               │               │                  │
        ▼               ▼               ▼                  ▼
  ┌───────────────────────────────────────────────────────────┐
  │                  Droplets (tagged)                         │
  │  • web-tier Droplets get web-tier-firewall rules          │
  │  • db-tier Droplets get db-tier-firewall rules            │
  │  • All Droplets get management-firewall rules             │
  └───────────────────────────────────────────────────────────┘
```

**Note**: Firewalls are applied via tags. A Droplet with tags `["web-tier", "all-instances"]` receives rules from both `web-tier-firewall` and `management-firewall`.

---

## Tag-Based Targeting (Production Pattern)

The module encourages tag-based targeting over static Droplet IDs:

### Protobuf Spec
```protobuf
message DigitalOceanFirewallSpec {
  string name = 1;
  repeated DigitalOceanFirewallInboundRule inbound_rules = 2;
  repeated DigitalOceanFirewallOutboundRule outbound_rules = 3;
  repeated int64 droplet_ids = 4;  // Dev/testing only (max 10)
  repeated string tags = 5;        // Production standard (unlimited)
}
```

### Pulumi Translation
```go
firewallArgs := &digitalocean.FirewallArgs{
    Name:          pulumi.String(locals.DigitalOceanFirewall.Spec.Name),
    DropletIds:    int64SliceToIntArray(locals.DigitalOceanFirewall.Spec.DropletIds),
    Tags:          pulumi.ToStringArray(locals.DigitalOceanFirewall.Spec.Tags),
    InboundRules:  inboundRules,
    OutboundRules: outboundRules,
}
```

**Production**: `tags` is populated, `droplet_ids` is empty  
**Dev**: `droplet_ids` is populated (≤10), `tags` is empty

---

## Error Handling

### Provider Setup Errors
```go
digitalOceanProvider, err := pulumidigitaloceanprovider.Get(ctx, stackInput.ProviderConfig)
if err != nil {
    return errors.Wrap(err, "failed to setup digitalocean provider")
}
```

**What it catches**: Invalid API token, network issues

### Resource Creation Errors
```go
createdFirewall, err := digitalocean.NewFirewall(ctx, "firewall", firewallArgs, pulumi.Provider(digitalOceanProvider))
if err != nil {
    return nil, errors.Wrap(err, "failed to create digitalocean firewall")
}
```

**What it catches**: API errors (duplicate name, invalid rule syntax, Droplet not found)

---

## Testing Strategy

### Unit Tests (Future)

Test rule translation logic in isolation:

```go
func TestInboundRuleTranslation(t *testing.T) {
    input := &digitaloceanfirewallv1.DigitalOceanFirewallInboundRule{
        Protocol:        "tcp",
        PortRange:       "443",
        SourceAddresses: []string{"0.0.0.0/0"},
    }

    result := translateInboundRule(input)

    assert.Equal(t, "tcp", result.Protocol)
    assert.Equal(t, "443", result.PortRange)
    assert.Equal(t, []string{"0.0.0.0/0"}, result.SourceAddresses)
}
```

### Integration Tests (Manual)

1. **Test tag-based targeting**:
   - Create firewall with `tags: ["test-web"]`
   - Create Droplet with tag `test-web`
   - Verify Droplet appears in firewall's Droplets list

2. **Test multi-tier architecture**:
   - Create web-tier-firewall (outbound to `db-tier` tag)
   - Create db-tier-firewall (inbound from `web-tier` tag)
   - Verify connectivity (web → db allowed, db → internet blocked)

3. **Test Load Balancer UID targeting**:
   - Create Load Balancer
   - Create firewall with `source_load_balancer_uids: [lb.id]`
   - Verify HTTPS traffic flows only through LB

---

## Common Patterns

### Pattern 1: Management Firewall (Applied to All Droplets)

```yaml
spec:
  name: management-firewall
  tags:
    - all-instances  # Apply to all Droplets
  inbound_rules:
    - protocol: tcp
      port_range: "22"
      source_addresses:
        - "203.0.113.10/32"  # Office bastion
```

**Usage**: Apply `all-instances` tag to every Droplet for centralized SSH access.

### Pattern 2: Web Tier with Load Balancer

```yaml
spec:
  name: web-tier-firewall
  tags:
    - web-tier
  inbound_rules:
    - protocol: tcp
      port_range: "443"
      source_load_balancer_uids:
        - "${load_balancer.id}"  # Pulumi interpolation
```

**Usage**: Web tier only accepts HTTPS from Load Balancer, never directly from internet.

### Pattern 3: Database Tier (Maximum Security)

```yaml
spec:
  name: db-tier-firewall
  tags:
    - db-tier
  inbound_rules:
    - protocol: tcp
      port_range: "5432"
      source_tags:
        - web-tier  # Only web tier can connect
  outbound_rules:
    - protocol: tcp
      port_range: "443"
      destination_addresses:
        - "91.189.88.0/21"  # Specific Ubuntu repos
```

**Usage**: Database is completely isolated from internet, only internal communication allowed.

---

## Debugging

### Enable Pulumi Debug Logging

```bash
export PULUMI_DEBUG_COMMANDS=true
pulumi up --logtostderr -v=9
```

### Check DigitalOcean API Connectivity

```bash
curl -H "Authorization: Bearer $DIGITALOCEAN_TOKEN" https://api.digitalocean.com/v2/account
```

### Verify Firewall Applied to Droplets

```bash
doctl compute firewall list
doctl compute firewall get <firewall-id>
```

### Check Droplet Tags

```bash
doctl compute droplet get <droplet-id> --format Name,Tags
```

---

## Performance Considerations

### Rule Count Limit

DigitalOcean enforces a **maximum of 50 rules per firewall** (inbound + outbound combined).

**Mitigation**: Split into multiple firewalls applied via tags:
- `management-firewall` (SSH, monitoring) → 5 rules
- `web-tier-firewall` (HTTPS, DB access) → 10 rules
- Total: 15 rules per Droplet (well under limit)

### Droplet ID Limit

Static Droplet IDs are limited to **10 per firewall**.

**Mitigation**: Use tag-based targeting (unlimited Droplets).

---

## Further Reading

- **Pulumi DigitalOcean Provider**: [Registry Docs](https://www.pulumi.com/registry/packages/digitalocean/)
- **DigitalOcean Firewall API**: [API Reference](https://docs.digitalocean.com/reference/api/api-reference/#tag/Firewalls)
- **Component Documentation**: See [../../docs/README.md](../../docs/README.md)

---

## Contributing

When extending this module:

1. **Preserve type safety**: Use `pulumi.String()`, `pulumi.Int()`, etc.
2. **Error wrapping**: Wrap errors with context using `errors.Wrap(err, "context")`
3. **Output naming**: Add new outputs to `outputs.go` with constants
4. **Testing**: Add unit tests for new translation logic
5. **Documentation**: Update this overview.md with design decisions

---

**Last Updated**: 2025-11-16

