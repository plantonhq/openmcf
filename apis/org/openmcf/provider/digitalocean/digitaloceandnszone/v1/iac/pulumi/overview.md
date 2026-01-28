# Pulumi Implementation Overview

## Architecture Philosophy

The Pulumi implementation for DigitalOcean DNS Zone follows a **modular, resource-centric** design that separates concerns while maintaining simplicity. Unlike monolithic Pulumi programs, this implementation is structured as a reusable module that can be invoked programmatically by the OpenMCF CLI or used standalone.

## Design Principles

### 1. Protobuf-First API

The implementation is driven entirely by the protobuf-defined `DigitalOceanDnsZoneSpec`. This ensures:

- **Type safety**: Go types generated from protobuf prevent runtime errors
- **Validation**: Protobuf validation rules (`buf.validate`) enforce constraints before Pulumi runs
- **Consistency**: Same spec can drive both Pulumi and Terraform implementations
- **Versioning**: API evolution is managed through protobuf versions (v1, v2, etc.)

### 2. Module Pattern

The code is split into two layers:

**Entrypoint (`main.go`)**:
- Loads manifest from file system
- Unmarshals protobuf
- Initializes Pulumi context
- Invokes module

**Module (`module/`)**:
- Receives parsed protobuf spec
- Creates DigitalOcean provider
- Provisions infrastructure
- Returns outputs

This separation allows:
- **Testing**: Module can be tested independently with mocked specs
- **Reusability**: Module can be imported by other Pulumi programs
- **Clarity**: Entrypoint concerns (file I/O, CLI args) are separate from provisioning logic

### 3. Stateful Locals

The `Locals` struct acts as shared state across module functions:

```go
type Locals struct {
    DigitalOceanDnsZone *digitaloceandnszonev1.DigitalOceanDnsZone
    Labels              map[string]string
}
```

This pattern:
- Avoids passing the same parameters to every function
- Provides a single source of truth for derived values (e.g., labels)
- Enables future expansion (e.g., adding computed values)

## Code Structure

```
iac/pulumi/
│
├── main.go                  # Entrypoint
│   ├── Load manifest YAML
│   ├── Unmarshal protobuf
│   └── pulumi.Run(...)
│       └── calls module.Resources(...)
│
└── module/
    ├── main.go              # Module orchestration
    │   ├── Resources(...)   # Public entry point
    │   │   ├── createProvider()
    │   │   ├── initLocals()
    │   │   └── dnsZone()
    │   │
    │   └── createProvider() # Provider initialization
    │
    ├── locals.go            # Shared state
    │   └── initLocals()     # Initialize Locals struct
    │
    ├── dns_zone.go          # Core provisioning
    │   └── dnsZone()
    │       ├── Create digitalocean.Domain
    │       ├── Loop over records
    │       │   ├── Build DnsRecordArgs
    │       │   ├── Set type-specific fields
    │       │   └── Create digitalocean.DnsRecord
    │       └── Export outputs
    │
    └── outputs.go           # Output definitions
        ├── OpZoneName
        ├── OpZoneId
        └── OpNameServers
```

## Key Implementation Decisions

### Multi-Value Record Handling

**Problem**: A single DNS record in the spec can have multiple values (e.g., multiple MX servers).

**Solution**: Create one DigitalOcean DNS record resource per value.

```go
for recIdx, rec := range locals.DigitalOceanDnsZone.Spec.Records {
    for valIdx, val := range rec.Values {
        resourceName := fmt.Sprintf("%s-%d-%d", rec.Name, recIdx, valIdx)
        digitalocean.NewDnsRecord(ctx, resourceName, ...)
    }
}
```

**Why**: DigitalOcean's API treats each value as a separate record. This approach:
- Matches the DigitalOcean data model
- Allows independent priority assignment for MX records
- Simplifies Pulumi state management (each resource has a unique ID)

### Conditional Field Setting

**Problem**: Different record types require different fields (priority for MX, weight/port for SRV, flags/tag for CAA).

**Solution**: Conditionally populate `DnsRecordArgs` based on record type.

```go
recordArgs := &digitalocean.DnsRecordArgs{
    Domain: createdDomain.Name,
    Name:   pulumi.String(rec.Name),
    Type:   pulumi.String(rec.Type.String()),
    Value:  pulumi.String(val.GetValue()),
    Ttl:    pulumi.Int(ttl),
}

// Add type-specific fields
if rec.Type.String() == "MX" || rec.Type.String() == "SRV" {
    recordArgs.Priority = pulumi.Int(int(rec.Priority))
}

if rec.Type.String() == "SRV" {
    recordArgs.Weight = pulumi.Int(int(rec.Weight))
    recordArgs.Port = pulumi.Int(int(rec.Port))
}

if rec.Type.String() == "CAA" {
    recordArgs.Flags = pulumi.Int(int(rec.Flags))
    recordArgs.Tag = pulumi.String(rec.Tag)
}
```

**Why**: This approach:
- Avoids setting null/zero values for irrelevant fields
- Makes code self-documenting (clear which fields apply to which types)
- Prevents API errors from invalid field combinations

### Default TTL Handling

**Problem**: TTL is optional in the spec with a recommended default of 3600.

**Solution**: Apply default in Go code.

```go
ttl := int(rec.TtlSeconds)
if ttl == 0 {
    ttl = 3600
}
```

**Why**: Protobuf's `optional` fields with defaults use field presence detection, but zero values can be ambiguous. Explicit handling in code:
- Provides clear fallback behavior
- Avoids DigitalOcean API rejecting 0 TTL
- Allows users to omit TTL in manifests for cleaner YAML

### Static Nameserver Exports

**Problem**: DigitalOcean doesn't return nameservers in the domain API response.

**Solution**: Hardcode DigitalOcean's nameservers in outputs.

```go
ctx.Export(
    OpNameServers,
    pulumi.StringArray{
        pulumi.String("ns1.digitalocean.com"),
        pulumi.String("ns2.digitalocean.com"),
        pulumi.String("ns3.digitalocean.com"),
    },
)
```

**Why**: DigitalOcean's nameservers are static and documented. Hardcoding:
- Provides immediate value to users (no need to look up nameservers)
- Simplifies delegation instructions
- Matches behavior of other DNS providers (Route 53, Cloudflare)

## Resource Naming Strategy

Pulumi resource names follow the pattern: `{record.name}-{recordIndex}-{valueIndex}`

**Example**:
```yaml
records:
  - name: "@"
    type: dns_record_type_mx
    values:
      - value: "mx1.example.com."
      - value: "mx2.example.com."
```

Generates Pulumi resources:
- `@-0-0` → MX record pointing to mx1.example.com
- `@-0-1` → MX record pointing to mx2.example.com

**Rationale**:
- **Uniqueness**: Guaranteed unique names even with duplicate record names
- **Stability**: Names don't change when record values are modified (only when order changes)
- **Traceability**: Easy to map Pulumi state to spec records

## Error Handling

The implementation uses structured error wrapping for clear failure messages:

```go
createdDomain, err := digitalocean.NewDomain(...)
if err != nil {
    return nil, errors.Wrap(err, "failed to create digitalocean domain")
}

createdDnsRecord, err := digitalocean.NewDnsRecord(...)
if err != nil {
    return nil, errors.Wrapf(err, "failed to create dns record %s", resourceName)
}
```

**Benefits**:
- Error context includes the failing operation
- Stack traces preserved through `errors.Wrap`
- Pulumi displays rich error messages with resource context

## Output Definitions

Outputs are defined as constants in `outputs.go`:

```go
const (
    OpZoneName    = "zone_name"
    OpZoneId      = "zone_id"
    OpNameServers = "name_servers"
)
```

**Why use constants**:
- **Compile-time safety**: Typos are caught by the compiler
- **Refactoring**: Renaming an output updates all references
- **Documentation**: Single source of truth for available outputs

## Dependencies and Versioning

### Pulumi Provider Version

```go
// go.mod
github.com/pulumi/pulumi-digitalocean/sdk/v4 v4.x.x
```

**Rationale**: v4 is the latest stable version of the DigitalOcean Pulumi provider, offering:
- Full API coverage (domains, records, all field types)
- Active maintenance and bug fixes
- Compatibility with Pulumi 3.x

### Pulumi SDK Version

```go
github.com/pulumi/pulumi/sdk/v3 v3.x.x
```

**Rationale**: v3 is the current major version, providing:
- Stable API for Go programs
- Modern resource lifecycle management
- Support for all Pulumi features (secrets, transformations, etc.)

## Testing Strategy

### Unit Testing (Future)

Planned approach:
```go
func TestDnsZone_BasicRecords(t *testing.T) {
    spec := &digitaloceandnszonev1.DigitalOceanDnsZone{
        Spec: &digitaloceandnszonev1.DigitalOceanDnsZoneSpec{
            DomainName: "test.com",
            Records: []*digitaloceandnszonev1.DigitalOceanDnsZoneRecord{
                {Name: "@", Type: dnsrecordtype.DnsRecordType_dns_record_type_a, ...},
            },
        },
    }
    
    // Use Pulumi's testing framework
    err := pulumi.RunErr(func(ctx *pulumi.Context) error {
        return module.Resources(ctx, spec, "fake-token")
    }, pulumi.WithMocks("project", "stack", &MockDigitalOcean{}))
    
    assert.NoError(t, err)
}
```

### Integration Testing

Current approach:
- Use `debug.sh` with real DigitalOcean credentials
- Deploy to a test domain (e.g., `test-planton-dns.com`)
- Verify records via `dig`
- Destroy resources

## Performance Characteristics

### Resource Creation

- **Domain creation**: ~1-2 seconds
- **DNS record creation**: ~0.5 seconds per record
- **Total deployment**: ~(1 + 0.5 * num_records) seconds

Example: Zone with 20 records ≈ 11 seconds

### API Rate Limiting

DigitalOcean API limit: **250 requests/minute**

Pulumi's default parallelism: **10 concurrent operations**

**Calculation**:
- 10 concurrent operations = 10 req/sec max
- 250/60 = 4.16 req/sec limit
- Pulumi stays well below limit with default settings

For very large zones (100+ records), Pulumi may hit rate limits. The DigitalOcean provider automatically retries with exponential backoff.

## Comparison to Terraform Implementation

| Aspect | Pulumi (Go) | Terraform |
|--------|-------------|-----------|
| **Language** | Go (strongly typed) | HCL (declarative) |
| **State** | Pulumi Service or file | State file |
| **Loops** | Native Go `for` loops | `for_each` meta-argument |
| **Conditionals** | Native Go `if` | Ternary operators |
| **Error Handling** | Go errors with stack traces | Error messages only |
| **Testing** | Go testing framework | Terratest (Go) |
| **Provider Version** | pulumi-digitalocean v4 | digitalocean provider v2 |

**When to choose Pulumi**:
- Team prefers general-purpose languages over DSLs
- Need complex conditional logic or custom functions
- Want to share code with application layer (e.g., validation helpers)
- Prefer integrated testing (Pulumi mocks vs. Terratest)

**When to choose Terraform**:
- Team is already proficient in HCL
- Want wider community ecosystem (modules, examples)
- Need to interoperate with existing Terraform infrastructure
- Prefer vendor-neutral tooling (Terraform is more widely adopted)

## Future Enhancements

### Planned Features

1. **Import Support**: Automate importing existing DigitalOcean domains into Pulumi state
2. **DNSSEC Warning**: Detect if DNSSEC is enabled at registrar and warn users
3. **Zone Transfer**: Utility to export DigitalOcean zone to BIND format for migration
4. **Record Validation**: Pre-flight validation for common mistakes (CNAME at apex, missing trailing dots)
5. **Diff Optimization**: Smart diffing to minimize record churn during updates

### Architectural Improvements

1. **Separate Record Types**: Create dedicated functions for each record type (e.g., `createMxRecord()`, `createCaaRecord()`) for better maintainability
2. **Resource Options**: Support Pulumi resource options (protect, delete before replace, etc.)
3. **Custom Timeouts**: Allow users to configure timeouts for slow DNS propagation scenarios
4. **Parallel Control**: Expose parallelism settings in spec for large zones

## Debugging Tips

### Enable Pulumi Debug Logging

```bash
export PULUMI_DEBUG_COMMANDS=true
pulumi up --logtostderr -v=9 2>&1 | tee pulumi-debug.log
```

### Inspect Provider Requests

```bash
export TF_LOG=DEBUG  # DigitalOcean provider is Terraform-based
pulumi up
```

### View Generated Resources

```bash
pulumi stack export | jq '.deployment.resources[] | select(.type == "digitalocean:index/dnsRecord:DnsRecord")'
```

### Test Record Creation Independently

```go
// In dns_zone.go, add logging
fmt.Printf("Creating record: name=%s type=%s value=%s\n", rec.Name, rec.Type, val.GetValue())
```

## References

- **Pulumi Go SDK**: https://www.pulumi.com/docs/languages-sdks/go/
- **DigitalOcean Provider Docs**: https://www.pulumi.com/registry/packages/digitalocean/api-docs/
- **Protobuf Spec**: ../../spec.proto
- **Research Document**: ../../docs/README.md

## Contributing

When modifying this implementation:

1. Maintain the module pattern (don't merge module back into main.go)
2. Update this overview.md if architectural changes are made
3. Add comments for non-obvious logic
4. Test with both simple (1 record) and complex (50+ records) zones
5. Verify outputs match expected format

---

**Last Updated**: 2025-11-16  
**Pulumi Version**: 3.x  
**DigitalOcean Provider Version**: 4.x

