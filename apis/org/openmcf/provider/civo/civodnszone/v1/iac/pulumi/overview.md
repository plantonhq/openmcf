# Civo DNS Zone - Pulumi Module Architecture

## Overview

This document provides an architectural overview of the Pulumi module for managing Civo DNS zones. It explains design decisions, implementation patterns, and how the module integrates with the broader OpenMCF ecosystem.

## Architecture Principles

### 1. Single Entry Point Pattern

The module exposes a single public function:

```go
func Resources(
    ctx *pulumi.Context,
    stackInput *civodnszonev1.CivoDnsZoneStackInput,
) error
```

**Rationale:**
- Simplifies invocation by OpenMCF CLI
- Ensures consistent interface across all deployment components
- Hides internal complexity from callers
- Enables version upgrades without breaking API contracts

### 2. Protobuf-First Configuration

Input is defined as a protobuf message (`CivoDnsZoneStackInput`) rather than JSON or YAML structures.

**Benefits:**
- **Type safety**: Compile-time validation of input structure
- **Schema evolution**: Protobuf supports backward-compatible changes
- **Cross-language**: Same schema used by CLI (Go), API server (Go), and potential future clients
- **Validation**: buf.validate annotations ensure data integrity before reaching Pulumi

### 3. Declarative Resource Model

The module implements Kubernetes-style declarative semantics:

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoDnsZone
metadata:
  name: example-zone
spec:
  domainName: example.com
  records: [...]
```

**Rationale:**
- Familiar to Kubernetes users
- Enables GitOps workflows
- Supports multi-cloud consistency (same pattern for AWS Route 53, GCP Cloud DNS, etc.)
- Allows future CRD (Custom Resource Definition) support for Kubernetes operators

### 4. Explicit Provider Management

The module receives Civo provider configuration via `stackInput.ProviderConfig` and instantiates a Pulumi provider explicitly:

```go
civoProvider, err := pulumicivoprovider.Get(ctx, stackInput.ProviderConfig)
```

**Why not use default provider?**
- Enables multi-account deployments (different zones in different Civo accounts)
- Supports credential rotation without stack recreation
- Allows testing with mock providers
- Aligns with OpenMCF's multi-tenancy model

## Component Breakdown

### `main.go` (Entrypoint)

**Purpose**: Bridge between Pulumi runtime and module logic.

```go
func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        var stackInput civodnszonev1.CivoDnsZoneStackInput
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
func Resources(ctx *pulumi.Context, stackInput *CivoDnsZoneStackInput) error {
    locals := initializeLocals(ctx, stackInput)
    civoProvider, err := pulumicivoprovider.Get(ctx, stackInput.ProviderConfig)
    if err != nil {
        return errors.Wrap(err, "failed to set up Civo provider")
    }
    
    if _, err := dnsZone(ctx, locals, civoProvider); err != nil {
        return errors.Wrap(err, "failed to create Civo DNS zone")
    }
    
    return nil
}
```

**Flow:**
1. Initialize locals (consolidate frequently-used values)
2. Set up Civo provider
3. Create DNS zone and records
4. Return any errors

**Error handling**: Uses `github.com/pkg/errors` for error wrapping, providing rich context for debugging.

### `module/locals.go` (Context Initialization)

**Purpose**: Extract and organize input data for easy access.

```go
type Locals struct {
    CivoProviderConfig *civoprovider.CivoProviderConfig
    CivoDnsZone        *civodnszonev1.CivoDnsZone
    CivoLabels         map[string]string
}
```

**CivoLabels**: Standard Planton labels attached to all resources:
- `resource: "true"` - Marks resource as managed by Planton
- `resource_name: <metadata.name>` - User-provided name
- `resource_kind: "CivoDnsZone"` - Component type
- `resource_id: <metadata.id>` - Unique identifier
- `organization: <metadata.org>` - Org context (if provided)
- `environment: <metadata.env>` - Environment (dev/staging/prod)

**Design rationale**: Labels enable:
- Cost tracking and reporting
- Resource discovery and querying
- Multi-tenancy isolation
- Audit trails

### `module/dns_zone.go` (Resource Provisioning)

**Purpose**: Create Civo DNS domain and records.

#### Zone Creation

```go
createdDomain, err := civo.NewDnsDomainName(
    ctx,
    "dns_zone",
    &civo.DnsDomainNameArgs{
        Name: pulumi.String(locals.CivoDnsZone.Spec.DomainName),
    },
    pulumi.Provider(civoProvider),
)
```

**Key points:**
- Resource name is hardcoded as `"dns_zone"` for stability
- Domain name comes from spec
- Explicit provider ensures correct credentials

#### Record Creation Strategy

**Challenge**: Protobuf defines records as:

```protobuf
message CivoDnsZoneRecord {
    string name = 1;
    repeated StringValueOrRef values = 2;
    uint32 ttl_seconds = 3;
    DnsRecordType type = 4;
}
```

Multiple values per record (for round-robin), but Civo API expects one value per record resource.

**Solution**: Flatten records to one Pulumi resource per value:

```go
for recIdx, rec := range locals.CivoDnsZone.Spec.Records {
    ttl := int(rec.TtlSeconds)
    if ttl == 0 {
        ttl = 3600  // Default TTL
    }
    
    for valIdx, val := range rec.Values {
        resourceName := fmt.Sprintf("%s-%d-%d", rec.Name, recIdx, valIdx)
        
        _, err := civo.NewDnsDomainRecord(
            ctx,
            resourceName,
            &civo.DnsDomainRecordArgs{
                DomainId: createdDomain.ID(),
                Name:     pulumi.String(rec.Name),
                Type:     pulumi.String(rec.Type.String()),
                Value:    pulumi.String(val.GetValue()),
                Ttl:      pulumi.Int(ttl),
            },
            pulumi.Provider(civoProvider),
        )
    }
}
```

**Implications:**
- A record with 3 values creates 3 Pulumi resources
- Resource names are deterministic: `<record-name>-<record-index>-<value-index>`
- Adding/removing values causes resource creation/deletion (not updates)
- TTL defaults to 3600 if not specified

**Alternative considered**: Create one Pulumi resource per record, with arrays of values. **Rejected** because Civo API doesn't support multi-value records natively, and this approach simplifies state tracking.

### `module/outputs.go` (Stack Outputs)

**Purpose**: Define constants for output keys.

```go
const (
    OpZoneName    = "zone_name"
    OpZoneId      = "zone_id"
    OpNameServers = "name_servers"
)
```

**Exported in `dns_zone.go`:**

```go
ctx.Export(OpZoneName, pulumi.String(locals.CivoDnsZone.Spec.DomainName))
ctx.Export(OpZoneId, createdDomain.ID())
ctx.Export(OpNameServers, pulumi.StringArray{
    pulumi.String("ns0.civo.com"),
    pulumi.String("ns1.civo.com"),
    pulumi.String("ns2.civo.com"),
})
```

**Design note**: Nameservers are **hardcoded** because Civo uses a fixed set of nameservers for all zones. This could theoretically change in the future, but Civo's API currently doesn't return zone-specific nameservers.

## Design Decisions

### 1. Why One Resource Per Record Value?

**Decision**: Create separate Pulumi resources for each value in a multi-value record.

**Alternatives:**
- **Option A**: Single resource with array of values
- **Option B**: One resource per unique (name, type) pair

**Rationale for chosen approach:**
- Civo API models records as individual entities
- Pulumi state tracking is cleaner with 1:1 mapping
- Enables fine-grained change detection
- Simplifies error handling (one failed value doesn't block others)

**Trade-off**: Slightly more resources in Pulumi state, but improved clarity and debuggability.

### 2. TTL Default Handling

**Decision**: Apply default TTL of 3600 seconds in Pulumi code if `ttl_seconds` is 0 or unspecified.

**Alternatives:**
- Rely on Civo's server-side default
- Make TTL required in protobuf (no default)

**Rationale:**
- Explicit defaults in code are self-documenting
- Protobuf's default value for `uint32` is 0, which is semantically "unspecified"
- 3600 (1 hour) is a sensible default per industry best practices
- Allows future customization of defaults per record type

### 3. StringValueOrRef for Record Values

**Decision**: Use `org.openmcf.shared.foreignkey.v1.StringValueOrRef` for record values.

**Structure:**
```protobuf
message StringValueOrRef {
    string value = 1;
    CloudResourceValueRef ref = 2;
}
```

**Rationale:**
- Supports literal values (common case): `{value: "192.0.2.1"}`
- Supports references to other resources: `{ref: {resource_id: "...", output_key: "ip"}}`
- Enables cross-resource dependencies (e.g., A record pointing to output of a compute instance)
- Future-proof for advanced orchestration

**Current limitation**: Reference resolution is not yet implemented. Only literal `value` is used. References will fail silently (empty string).

**Future work**: Implement reference resolution in a shared library that all modules can use.

### 4. Error Wrapping Strategy

**Decision**: Use `github.com/pkg/errors.Wrap` for all errors.

**Example:**
```go
if err != nil {
    return errors.Wrap(err, "failed to create DNS record %s", resourceName)
}
```

**Benefits:**
- Stack traces show full error propagation path
- Context added at each layer
- Debugging production issues becomes much easier

**Alternative considered**: Standard Go errors with `fmt.Errorf`. **Rejected** because stack traces are invaluable in complex Pulumi programs.

### 5. Hardcoded Nameservers

**Decision**: Export hardcoded nameservers instead of querying Civo API.

**Rationale:**
- Civo uses the same nameservers for all zones
- Querying API adds latency and potential failure point
- Nameservers are stable and unlikely to change
- If they do change, it's a breaking change requiring Civo migration guide

**Monitoring**: If Civo changes nameservers, users will report DNS resolution failures, triggering manual update of constants.

### 6. No DNS Record Validation

**Decision**: Don't validate record values (e.g., IP address format, domain name syntax) in Pulumi code.

**Rationale:**
- Validation belongs in protobuf schema (buf.validate)
- Pulumi runs after validation, so input is assumed valid
- Civo API provides server-side validation as last line of defense
- Avoiding duplication of validation logic

**Exception**: We validate that `values` array is not empty (required by buf.validate, but defensive check in code).

## Integration Points

### OpenMCF CLI

**Invocation flow:**
1. User runs `planton apply -f dns-zone.yaml`
2. CLI parses YAML → protobuf `CivoDnsZone`
3. CLI validates via buf.validate
4. CLI constructs `CivoDnsZoneStackInput` (adds provider config)
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
- Consistency across all Civo modules (compute, DNS, Kubernetes, etc.)
- Single place to add features (e.g., proxy support, retries)
- Easier credential rotation

## State Management

### Pulumi State Structure

For a zone with 2 records:

```
Zone: cidns-example-123
├── civo:index/dnsDomainName:DnsDomainName
│   └── dns_zone
└── Records
    ├── civo:index/dnsDomainRecord:DnsDomainRecord
    │   └── @-0-0 (A record, value 1)
    ├── civo:index/dnsDomainRecord:DnsDomainRecord
    │   └── www-1-0 (CNAME record, value 1)
```

**Resource naming:**
- Zone: `"dns_zone"` (constant)
- Records: `<name>-<recIdx>-<valIdx>`

**Why stable names matter**: Pulumi uses resource names (URNs) to track state. Changing a resource name causes Pulumi to delete the old resource and create a new one, even if the configuration is identical.

### State Drift Scenarios

#### Scenario 1: Manual change in Civo dashboard

User manually adds a record via Civo dashboard.

**Result**: Pulumi doesn't know about it. On next `pulumi refresh`, Pulumi detects drift but doesn't import the new record.

**Resolution**: Run `pulumi up` to converge to desired state (no-op if spec unchanged) or manually remove record from dashboard.

#### Scenario 2: Record deleted externally

User deletes a record via Civo CLI.

**Result**: Pulumi shows resource in state but no longer exists in Civo. On next `pulumi up`, Pulumi will attempt to recreate it.

**Resolution**: Automatic recovery. Pulumi recreates the missing record.

#### Scenario 3: Concurrent modifications

Two users run `pulumi up` simultaneously on the same zone.

**Result**: Pulumi's state locking prevents corruption, but last write wins. Changes from first user may be overwritten.

**Resolution**: Use remote state backend with locking (S3, Pulumi Cloud). Educate users to coordinate deployments or use OpenMCF CLI (which serializes operations).

## Performance Characteristics

### Zone Creation Time

- **DNS domain creation**: ~2-3 seconds
- **Per record**: ~1-2 seconds
- **Total for 50 records**: ~1-2 minutes

**Bottleneck**: Civo API rate limits. Each record is a separate API call.

**Optimization**: Records are created in sequence (not parallel) to respect rate limits and avoid errors. Pulumi's dependency graph automatically parallelizes where possible.

### Update Operations

- **TTL change**: In-place update (~1-2 seconds per record)
- **Value change**: Replace record (~2-4 seconds per record)
- **Type change**: Replace record (delete + create)

**Note**: Changing record name causes replacement because Civo treats name as part of record identity.

### Destroy Operations

- **Delete all records**: ~1-2 seconds per record
- **Delete zone**: ~2-3 seconds

**Total for 50 records**: ~1-2 minutes

**Order**: Pulumi automatically deletes records before zone due to dependency graph.

## Security Considerations

### Credential Handling

**API Token Storage:**
- Passed via `ProviderConfig.Credential.ApiToken` (string)
- **Never logged** by Pulumi or module code
- Should be encrypted at rest in OpenMCF's credential store
- Recommended: Use Pulumi secrets for standalone usage

**Permissions**: Token needs `DNS: Write` permission in Civo. **Do not** use tokens with broader permissions (e.g., compute, Kubernetes).

### Zone Takeover Prevention

**Risk**: If a zone is deleted, someone else could recreate it with the same domain name in a different Civo account.

**Mitigation**: 
- Pulumi state tracks zone ID (UUID), not just name
- If zone is recreated externally, Pulumi detects ID mismatch and errors
- User must manually resolve conflict

**Best practice**: Use `pulumi destroy` to cleanly delete zones. Avoid manual deletion in Civo dashboard.

### DNS Hijacking

**Risk**: Attacker gains access to Civo API token and modifies DNS records.

**Mitigation**:
- Rotate API tokens regularly
- Use time-limited tokens (if Civo adds support)
- Monitor DNS records for unexpected changes (external monitoring)
- Enable Civo account 2FA

**Future enhancement**: Implement output hash in stack state. On refresh, compare hash to detect tampering.

## Testing Strategy

### Unit Tests

**Location**: `module/*_test.go` (not yet implemented)

**Coverage:**
- `initializeLocals`: Verify label construction
- TTL defaulting logic
- Resource name generation

**Mocking**: Use `pulumi-go-provider` test utilities to mock Pulumi context.

### Integration Tests

**Location**: `iac/pulumi/integration_test.go` (not yet implemented)

**Approach:**
1. Create test Civo account/API token
2. Run `module.Resources` with test input
3. Query Civo API to verify zone and records created
4. Run `pulumi destroy`
5. Verify cleanup

**Challenges**: Requires live Civo credentials. Consider using Civo's staging environment if available.

### Validation Tests

**Location**: `v1/spec_test.go` ✅ (implemented)

**Coverage:**
- Valid domain patterns
- Invalid domain patterns
- Required fields
- Record validation (types, values, TTL)

**Benefit**: Catches invalid input before reaching Pulumi.

## Troubleshooting Guide

### Issue: "Domain already exists"

**Symptom**: Pulumi error during `civo.NewDnsDomainName`.

**Cause**: Domain is already registered in Civo (possibly in another account or region).

**Resolution:**
1. Check `civo domain list` for existing zone
2. If zone exists, import it: `pulumi import civo:index/dnsDomainName:DnsDomainName dns_zone <domain-name>`
3. If zone is in wrong account, delete it first

### Issue: Records not resolving

**Symptom**: `dig example.com` returns NXDOMAIN.

**Cause**: Nameservers not updated at registrar.

**Resolution:**
1. Get nameservers: `pulumi stack output name_servers`
2. Update at registrar (GoDaddy, Namecheap, etc.)
3. Wait up to 48 hours for propagation
4. Test directly: `dig @ns0.civo.com example.com`

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
2. Test token: `curl -H "Authorization: Bearer $TOKEN" https://api.civo.com/v2/dns`
3. Update `ProviderConfig.Credential.ApiToken`
4. Retry `pulumi up`

## Future Enhancements

### Planned

1. **Reference resolution**: Implement `StringValueOrRef.ref` to support cross-resource dependencies
2. **Import command**: Add helper script to import existing zones
3. **Bulk operations**: Batch record creation API calls for performance
4. **Health checks**: Optional external DNS monitoring integration

### Under Consideration

1. **DNSSEC support**: When Civo adds API support
2. **Geo-routing**: If Civo introduces advanced DNS features
3. **Record policies**: Validation rules per record type (e.g., MX must have priority)
4. **Change previews**: Show TTL impact, propagation estimates

## References

- [Pulumi Civo Provider Docs](https://www.pulumi.com/registry/packages/civo/)
- [Civo DNS API](https://www.civo.com/api/dns)
- [OpenMCF Architecture](../../../../../architecture/)
- [User Documentation](../../README.md)

## Changelog

- **2025-11-16**: Initial architecture documentation
- **2025-11-14**: Implementation completed
- **2025-11-10**: Module scaffolding created

---

**Maintained by**: OpenMCF Team  
**Last Updated**: 2025-11-16

