# Civo Firewall - Pulumi Module Architecture

## Overview

This document provides an architectural overview of the Pulumi module for managing Civo firewalls. It explains design decisions, implementation patterns, and how the module integrates with the broader OpenMCF ecosystem.

## Architecture Principles

### 1. Single Entry Point Pattern

The module exposes a single public function:

```go
func Resources(
    ctx *pulumi.Context,
    stackInput *civofirewallv1.CivoFirewallStackInput,
) error
```

**Rationale:**
- Simplifies invocation by OpenMCF CLI
- Ensures consistent interface across all deployment components
- Hides internal complexity from callers
- Enables version upgrades without breaking API contracts

### 2. Protobuf-First Configuration

Input is defined as a protobuf message (`CivoFirewallStackInput`) rather than JSON or YAML structures.

**Benefits:**
- **Type safety**: Compile-time validation of input structure
- **Schema evolution**: Protobuf supports backward-compatible changes
- **Cross-language**: Same schema used by CLI (Go), API server (Go), and potential future clients
- **Validation**: buf.validate annotations ensure data integrity before reaching Pulumi

### 3. Declarative Resource Model

The module implements Kubernetes-style declarative semantics:

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoFirewall
metadata:
  name: web-server-firewall
spec:
  name: web-server-fw
  networkId:
    value: "network-123"
  inboundRules: [...]
```

**Rationale:**
- Familiar to Kubernetes users
- Enables GitOps workflows
- Supports multi-cloud consistency (same pattern for AWS security groups, GCP firewall rules, etc.)
- Allows future CRD (Custom Resource Definition) support for Kubernetes operators

### 4. Explicit Provider Management

The module receives Civo provider configuration via `stackInput.ProviderConfig` and instantiates a Pulumi provider explicitly:

```go
civoProvider, err := pulumicivoprovider.Get(ctx, stackInput.ProviderConfig)
```

**Why not use default provider?**
- Enables multi-account deployments (different firewalls in different Civo accounts)
- Supports credential rotation without stack recreation
- Allows testing with mock providers
- Aligns with OpenMCF's multi-tenancy model

## Component Breakdown

### `main.go` (Entrypoint)

**Purpose**: Bridge between Pulumi runtime and module logic.

```go
func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        var stackInput civofirewallv1.CivoFirewallStackInput
        // ... deserialize from config ...
        return module.Resources(ctx, &stackInput)
    })
}
```

**Key responsibilities:**
- Parse stack input from Pulumi config
- Invoke `module.Resources`
- Handle top-level errors

**Design note**: When invoked by OpenMCF CLI, this `main` function is **not** used. The CLI directly calls `module.Resources` as a library. The `main` is provided for standalone Pulumi usage.

### `module/main.go` (Core Logic)

**Purpose**: Orchestrate resource creation.

```go
func Resources(ctx *pulumi.Context, stackInput *CivoFirewallStackInput) error {
    locals := initializeLocals(ctx, stackInput)
    civoProvider, err := pulumicivoprovider.Get(ctx, stackInput.ProviderConfig)
    if err != nil {
        return errors.Wrap(err, "failed to set up Civo provider")
    }
    
    if _, err := firewall(ctx, locals, civoProvider); err != nil {
        return errors.Wrap(err, "failed to create firewall")
    }
    
    return nil
}
```

**Flow:**
1. Initialize locals (consolidate frequently-used values)
2. Set up Civo provider
3. Create firewall with rules
4. Return any errors

**Error handling**: Uses `github.com/pkg/errors` for error wrapping, providing rich context for debugging.

### `module/locals.go` (Context Initialization)

**Purpose**: Extract and organize input data for easy access.

```go
type Locals struct {
    CivoProviderConfig *civoprovider.CivoProviderConfig
    CivoFirewall       *civofirewallv1.CivoFirewall
    CivoLabels         map[string]string
}
```

**CivoLabels**: Standard Planton labels attached to all resources:
- `resource: "true"` - Marks resource as managed by Planton
- `resource_name: <metadata.name>` - User-provided name
- `resource_kind: "CivoFirewall"` - Component type
- `resource_id: <metadata.id>` - Unique identifier
- `organization: <metadata.org>` - Org context (if provided)
- `environment: <metadata.env>` - Environment (dev/staging/prod)

**Design rationale**: Labels enable:
- Cost tracking and reporting
- Resource discovery and querying
- Multi-tenancy isolation
- Audit trails

### `module/firewall.go` (Resource Provisioning)

**Purpose**: Create Civo firewall with inbound and outbound rules.

#### Firewall Creation

```go
createdFirewall, err := civo.NewFirewall(
    ctx,
    "firewall",
    &civo.FirewallArgs{
        Name:         pulumi.String(locals.CivoFirewall.Spec.Name),
        NetworkId:    pulumi.String(locals.CivoFirewall.Spec.NetworkId.GetValue()),
        IngressRules: inboundRules,
        EgressRules:  outboundRules,
    },
    pulumi.Provider(civoProvider),
)
```

**Key points:**
- Resource name is hardcoded as `"firewall"` for stability
- Network ID extracted from StringValueOrRef (supports literal values and references)
- Explicit provider ensures correct credentials
- Rules translated from protobuf to Civo provider format

#### Rule Translation Strategy

**Challenge**: Protobuf defines rules as:

```protobuf
message CivoFirewallInboundRule {
    string protocol = 1;
    string port_range = 2;
    repeated string cidrs = 3;
    string action = 4;
    string label = 5;
}
```

Civo's Pulumi provider expects:

```go
type FirewallIngressRuleArgs struct {
    Protocol  pulumi.StringInput
    PortRange pulumi.StringInput
    Cidrs     pulumi.StringArrayInput
    Action    pulumi.StringInput
    Label     pulumi.StringInput
}
```

**Solution**: Map protobuf fields to Pulumi provider format:

```go
inboundRules := make(civo.FirewallIngressRuleArray, 0, len(locals.CivoFirewall.Spec.InboundRules))
for _, rule := range locals.CivoFirewall.Spec.InboundRules {
    inboundRules = append(inboundRules, civo.FirewallIngressRuleArgs{
        Protocol:  pulumi.String(rule.Protocol),
        PortRange: pulumi.String(rule.PortRange),
        Cidrs:     pulumi.ToStringArray(rule.Cidrs),
        Action:    pulumi.String(rule.Action),
        Label:     pulumi.String(rule.Label),
    })
}
```

**Implications:**
- All rules created in a single Pulumi resource (atomic operation)
- Adding/removing rules updates the firewall in-place
- Rule order is preserved from protobuf spec
- Empty action defaults to "allow" (Civo default)

### `module/outputs.go` (Stack Outputs)

**Purpose**: Define constants for output keys.

```go
const (
    OpFirewallId = "firewall_id"
)
```

**Exported in `firewall.go`:**

```go
ctx.Export(OpFirewallId, createdFirewall.ID())
```

**Design note**: Using constants instead of magic strings prevents typos and enables IDE autocomplete when consuming outputs.

## Design Decisions

### 1. Why Single Firewall Resource (Not Separate Rule Resources)?

**Decision**: Create all rules within a single `civo.Firewall` resource, not separate `civo.FirewallRule` resources.

**Alternatives:**
- **Option A**: Separate Pulumi resource per rule (like Terraform's `civo_firewall_rule`)
- **Option B**: All rules in single firewall resource (chosen)

**Rationale for chosen approach:**
- Civo's Pulumi provider supports inline rules (IngressRules, EgressRules arrays)
- Atomic updates: All rules change together or not at all
- Simpler state management (one resource vs many)
- Better alignment with protobuf spec structure
- Reduces Pulumi state size

**Trade-off**: Large firewalls (50+ rules) result in large single resource, but this is acceptable for typical use cases.

### 2. StringValueOrRef for Network ID

**Decision**: Use `org.openmcf.shared.foreignkey.v1.StringValueOrRef` for network_id.

**Structure:**
```protobuf
message StringValueOrRef {
    oneof literal_or_ref {
        string value = 1;
        ValueFromRef value_from = 2;
    }
}
```

**Rationale:**
- Supports literal values (common case): `{value: "network-uuid"}`
- Supports references to other resources: `{value_from: {kind: "CivoVpc", name: "...", field_path: "status.outputs.network_id"}}`
- Enables cross-resource dependencies (e.g., firewall referencing VPC)
- Future-proof for advanced orchestration

**Current limitation**: Reference resolution is not yet implemented. Only literal `value` is used. References will fail silently (empty string).

**Future work**: Implement reference resolution in a shared library that all modules can use.

### 3. Protocol Validation in Protobuf

**Decision**: Validate protocol at protobuf level with regex pattern.

```protobuf
string protocol = 1 [(buf.validate.field).string.pattern = "^(tcp|udp|icmp)$"];
```

**Rationale:**
- Validation happens before Pulumi runs
- Fails fast with clear error message
- Prevents costly Pulumi updates with invalid data
- Protobuf validation is portable across languages

**Alternative considered**: Validate in Pulumi code. **Rejected** because:
- Duplicate validation logic
- Later error detection (after stack input deserialization)
- Less user-friendly error messages

### 4. Empty CIDRs Default Behavior

**Decision**: Allow empty `cidrs` array, defaulting to `["0.0.0.0/0"]` in Civo.

**Protobuf:**
```protobuf
repeated string cidrs = 3;  // No required validation
```

**Rationale:**
- Civo's default behavior is "from anywhere" if no CIDRs specified
- Matches user expectation for public services (HTTP/HTTPS)
- Reduces verbosity for common cases

**Security note**: Users should be aware that empty CIDRs = public access. Documentation emphasizes explicit CIDR specification.

### 5. Error Wrapping Strategy

**Decision**: Use `github.com/pkg/errors.Wrap` for all errors.

**Example:**
```go
if err != nil {
    return errors.Wrap(err, "failed to create Civo firewall")
}
```

**Benefits:**
- Stack traces show full error propagation path
- Context added at each layer
- Debugging production issues becomes much easier

**Alternative considered**: Standard Go errors with `fmt.Errorf`. **Rejected** because stack traces are invaluable in complex Pulumi programs.

### 6. No Rule Validation in Module

**Decision**: Don't validate rule specifics (IP format, port numbers) in Pulumi code.

**Rationale:**
- Validation belongs in protobuf schema (buf.validate)
- Pulumi runs after validation, so input is assumed valid
- Civo API provides server-side validation as last line of defense
- Avoiding duplication of validation logic

**Exception**: We do check that network_id is not empty (defensive programming).

## Integration Points

### OpenMCF CLI

**Invocation flow:**
1. User runs `planton apply -f firewall.yaml`
2. CLI parses YAML → protobuf `CivoFirewall`
3. CLI validates via buf.validate
4. CLI constructs `CivoFirewallStackInput` (adds provider config)
5. CLI invokes `module.Resources(ctx, stackInput)` directly (not via `pulumi up`)
6. CLI captures outputs and displays to user

**Key point**: No `Pulumi.yaml` or stack files involved. CLI manages Pulumi state in its own storage.

### Standalone Pulumi Usage

**Invocation flow:**
1. User creates `Pulumi.yaml` and `stack-input.json`
2. User runs `pulumi up`
3. `main.go` reads config, deserializes to protobuf
4. Calls `module.Resources(ctx, stackInput)`
5. Pulumi manages state in configured backend (S3, local, Pulumi Cloud)

**Key point**: Same `module.Resources` function, different entry points.

### Civo Provider Setup

The module uses a shared helper to set up the Civo Pulumi provider:

```go
import "github.com/plantonhq/openmcf/pkg/iac/pulumi/pulumimodule/provider/civo/pulumicivoprovider"

civoProvider, err := pulumicivoprovider.Get(ctx, stackInput.ProviderConfig)
```

**`pulumicivoprovider.Get` does:**
- Extract API token from `ProviderConfig`
- Handle token source (direct value, secret ref, env var)
- Create Pulumi explicit provider
- Set region if specified
- Apply common options

**Benefits of centralized provider setup:**
- Consistency across all Civo modules (firewall, compute, DNS, Kubernetes, etc.)
- Single place to add features (e.g., proxy support, retries)
- Easier credential rotation

## State Management

### Pulumi State Structure

For a firewall with 3 inbound rules:

```
Firewall: cifw-web-123
└── civo:index/firewall:Firewall
    └── firewall
        ├── IngressRules (3 rules)
        └── EgressRules (0 rules)
```

**Resource naming:**
- Firewall: `"firewall"` (constant)

**Why stable names matter**: Pulumi uses resource names (URNs) to track state. Changing a resource name causes Pulumi to delete the old resource and create a new one, even if the configuration is identical.

### State Drift Scenarios

#### Scenario 1: Manual change in Civo dashboard

User manually adds a rule via Civo dashboard.

**Result**: Pulumi doesn't know about it. On next `pulumi refresh`, Pulumi detects drift but doesn't import the new rule.

**Resolution**: Run `pulumi up` to converge to desired state (removes the manually-added rule) or update spec to match.

#### Scenario 2: Rule deleted externally

User deletes a rule via Civo CLI.

**Result**: Pulumi shows resource in state but no longer exists in Civo. On next `pulumi up`, Pulumi will attempt to recreate it.

**Resolution**: Automatic recovery. Pulumi recreates the missing rule.

#### Scenario 3: Concurrent modifications

Two users run `pulumi up` simultaneously on the same firewall.

**Result**: Pulumi's state locking prevents corruption, but last write wins. Changes from first user may be overwritten.

**Resolution**: Use remote state backend with locking (S3, Pulumi Cloud). Educate users to coordinate deployments or use OpenMCF CLI (which serializes operations).

## Performance Characteristics

### Firewall Creation Time

- **Firewall creation**: ~2-3 seconds
- **Per rule**: Negligible (rules created inline)
- **Total for 20 rules**: ~3-5 seconds

**Bottleneck**: Civo API latency. Rules are submitted as a batch in the firewall creation call.

### Update Operations

- **Add rule**: In-place update (~2-3 seconds)
- **Remove rule**: In-place update (~2-3 seconds)
- **Change rule**: In-place update (~2-3 seconds)
- **Change network_id**: Firewall replacement (delete + create, ~5-10 seconds)

**Note**: Changing firewall name or network_id causes replacement because these are identity fields.

### Destroy Operations

- **Delete firewall**: ~2-3 seconds

**Order**: Pulumi automatically handles dependencies (firewall deleted after instances detached).

## Security Considerations

### Credential Handling

**API Token Storage:**
- Passed via `ProviderConfig.Credential.ApiToken` (string)
- **Never logged** by Pulumi or module code
- Should be encrypted at rest in OpenMCF's credential store
- Recommended: Use Pulumi secrets for standalone usage

**Permissions**: Token needs `Network: Write` permission in Civo. **Do not** use tokens with broader permissions (e.g., compute, Kubernetes).

### Default-Deny Security Model

**Risk**: If inbound rules array is empty, Civo creates a default-deny firewall (no traffic allowed).

**Mitigation**: 
- Documentation clearly states default-deny behavior
- Users should explicitly define rules for all required traffic
- Testing workflow ensures connectivity before production deployment

**Best practice**: Start with restrictive rules, test, then expand as needed.

### Firewall Takeover Prevention

**Risk**: If a firewall is deleted, someone else could recreate it with the same name in a different account.

**Mitigation**: 
- Pulumi state tracks firewall ID (UUID), not just name
- If firewall is recreated externally, Pulumi detects ID mismatch and errors
- User must manually resolve conflict

**Best practice**: Use `pulumi destroy` to cleanly delete firewalls. Avoid manual deletion in Civo dashboard.

## Testing Strategy

### Unit Tests

**Location**: `module/*_test.go` (not yet implemented)

**Coverage:**
- `initializeLocals`: Verify label construction
- Rule translation: Ensure protobuf → Civo provider mapping is correct
- Resource name generation

**Mocking**: Use `pulumi-go-provider` test utilities to mock Pulumi context.

### Integration Tests

**Location**: `iac/pulumi/integration_test.go` (not yet implemented)

**Approach:**
1. Create test Civo account/API token
2. Run `module.Resources` with test input
3. Query Civo API to verify firewall and rules created
4. Run `pulumi destroy`
5. Verify cleanup

**Challenges**: Requires live Civo credentials. Consider using Civo's staging environment if available.

### Validation Tests

**Location**: `v1/spec_test.go` ✅ (implemented - 23 tests)

**Coverage:**
- Valid protocol patterns (tcp, udp, icmp)
- Invalid protocol patterns (uppercase, empty, invalid values)
- Required fields (name, network_id)
- Port range formats
- CIDR notation

**Benefit**: Catches invalid input before reaching Pulumi.

## Troubleshooting Guide

### Issue: "Firewall name already exists"

**Symptom**: Pulumi error during `civo.NewFirewall`.

**Cause**: Firewall with same name exists in the network (possibly in another account or region).

**Resolution:**
1. Check `civo firewall list` for existing firewall
2. If firewall exists, import it: `pulumi import civo:index/firewall:Firewall firewall <firewall-id>`
3. If firewall is in wrong account, delete it first

### Issue: Rules not blocking traffic

**Symptom**: Traffic still reaching instance despite deny rule.

**Cause**: 
- Civo's stateful firewall automatically allows return traffic
- Rule order might be incorrect (allow before deny)

**Resolution:**
1. Verify rule syntax: `planton get civofirewalls/your-firewall -o yaml`
2. Check rule order (deny rules should come first)
3. Test from expected source/destination
4. Verify firewall is attached to correct instances

### Issue: Pulumi state desync

**Symptom**: `pulumi preview` shows unexpected changes.

**Cause**: Manual changes in Civo dashboard.

**Resolution:**
1. Run `pulumi refresh` to sync state
2. Run `pulumi preview` to see actual diff
3. Update spec to match desired state
4. Run `pulumi up` to converge

### Issue: Provider authentication failed

**Symptom**: Error: "401 Unauthorized" or "Invalid API token".

**Cause**: Wrong token or expired credentials.

**Resolution:**
1. Verify token: `civo apikey list`
2. Test token: `curl -H "Authorization: Bearer $TOKEN" https://api.civo.com/v2/firewalls`
3. Update `ProviderConfig.Credential.ApiToken`
4. Retry `pulumi up`

## Future Enhancements

### Planned

1. **Reference resolution**: Implement `StringValueOrRef.value_from` to support cross-resource dependencies
2. **Import command**: Add helper script to import existing firewalls
3. **Rule templates**: Pre-defined rule sets for common scenarios (web server, database, K8s)
4. **Validation**: Additional runtime validation of CIDR formats and port ranges

### Under Consideration

1. **Firewall groups**: Manage multiple related firewalls as a unit
2. **Change previews**: Show impact of rule changes (e.g., which instances affected)
3. **Audit logging**: Track all firewall modifications with timestamps and users
4. **Cost tracking**: Estimate costs based on network egress (Civo doesn't charge for firewalls, but network usage matters)

## References

- [Pulumi Civo Provider Docs](https://www.pulumi.com/registry/packages/civo/)
- [Civo Firewall API](https://www.civo.com/api/firewalls)
- [OpenMCF Architecture](../../../../../architecture/)
- [User Documentation](../../README.md)

## Changelog

- **2025-11-16**: Initial architecture documentation
- **2025-11-14**: Implementation completed
- **2025-11-10**: Module scaffolding created

---

**Maintained by**: OpenMCF Team  
**Last Updated**: 2025-11-16

