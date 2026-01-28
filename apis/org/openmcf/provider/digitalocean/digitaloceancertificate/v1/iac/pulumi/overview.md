# DigitalOcean Certificate - Pulumi Module Architecture

## Purpose

This document explains the internal architecture, design decisions, and implementation patterns of the DigitalOcean Certificate Pulumi module.

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
    │ locals.go│      │certificate.go│     │outputs.go│
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

4. **`module/certificate.go`**
   - Core resource creation logic
   - Implements discriminated union (Let's Encrypt vs Custom)
   - Exports stack outputs

5. **`module/outputs.go`**
   - Defines output constant names
   - Centralizes output keys for consistency

---

## Data Flow

### 1. Input: DigitalOceanCertificateStackInput

```protobuf
message DigitalOceanCertificateStackInput {
  ProviderConfig provider_config = 1;          // DigitalOcean credentials
  DigitalOceanCertificate target = 2;          // Certificate spec
}

message DigitalOceanCertificate {
  CloudResourceMetadata metadata = 1;          // Name, ID, env, org
  DigitalOceanCertificateSpec spec = 2;        // Certificate configuration
}
```

### 2. Processing: Locals Initialization

```go
// In locals.go
type Locals struct {
    DigitalOceanCertificate *digitaloceancertificatev1.DigitalOceanCertificate
    Labels                   map[string]string  // Derived from metadata
}

func initializeLocals(ctx *pulumi.Context, stackInput *stackInput) *Locals {
    // Extract spec, metadata
    // Transform into Pulumi-usable format
}
```

### 3. Resource Creation: Certificate

```go
// In certificate.go
func certificate(
    ctx *pulumi.Context,
    locals *Locals,
    provider *digitalocean.Provider,
) (*digitalocean.Certificate, error) {
    // 1. Build CertificateArgs based on type (Let's Encrypt vs Custom)
    // 2. Call digitalocean.NewCertificate()
    // 3. Export outputs (certificate_id, expiry_rfc3339)
}
```

### 4. Output: Stack Exports

```go
ctx.Export(OpCertificateId, createdCertificate.ID())
ctx.Export(OpExpiryRfc3339, createdCertificate.NotAfter)
```

Users retrieve these via `pulumi stack output certificate_id`.

---

## Discriminated Union Implementation

The spec uses protobuf `oneof` to enforce mutual exclusivity:

```protobuf
message DigitalOceanCertificateSpec {
  string certificate_name = 1;
  DigitalOceanCertificateType type = 2;

  oneof certificate_source {
    DigitalOceanCertificateLetsEncryptParams lets_encrypt = 3;
    DigitalOceanCertificateCustomParams custom = 4;
  }
}
```

### Go Implementation

The module uses conditional logic to handle both paths:

```go
certArgs := &digitalocean.CertificateArgs{
    Name: pulumi.String(locals.DigitalOceanCertificate.Spec.CertificateName),
    Type: pulumi.String(locals.DigitalOceanCertificate.Spec.Type.String()),
}

// Let's Encrypt path
if locals.DigitalOceanCertificate.Spec.Type == digitaloceancertificatev1.DigitalOceanCertificateType_letsEncrypt {
    certArgs.Domains = domains
}

// Custom path
if locals.DigitalOceanCertificate.Spec.Type == digitaloceancertificatev1.DigitalOceanCertificateType_custom {
    certArgs.LeafCertificate = pulumi.String(locals.DigitalOceanCertificate.Spec.GetCustom().LeafCertificate)
    certArgs.PrivateKey = pulumi.String(locals.DigitalOceanCertificate.Spec.GetCustom().PrivateKey)
    if locals.DigitalOceanCertificate.Spec.GetCustom().CertificateChain != "" {
        certArgs.CertificateChain = pulumi.StringPtr(locals.DigitalOceanCertificate.Spec.GetCustom().CertificateChain)
    }
}
```

**Why this works:**
- Protobuf guarantees only one of `lets_encrypt` or `custom` can be set
- Type safety prevents invalid configurations at compile-time
- Runtime checks are defensive (already validated by protobuf)

---

## Design Decisions

### 1. Why Separate `certificate.go`?

**Rationale**: Following OpenMCF's convention, resource-specific logic is isolated in dedicated files (`certificate.go`, not inline in `main.go`). This:
- Improves readability (main.go is orchestration-only)
- Simplifies testing (can test certificate logic independently)
- Enables future extensions (e.g., adding validation pre-checks)

### 2. Why Not Use Pulumi's Native `oneof` Support?

**Context**: Pulumi's Go SDK doesn't have first-class `oneof` support (unlike TypeScript's union types).

**Solution**: Use Go's type assertions and conditional logic:
```go
spec.GetLetsEncrypt()  // Returns nil if Custom is set
spec.GetCustom()        // Returns nil if LetsEncrypt is set
```

This pattern is idiomatic for protobuf in Go and matches other OpenMCF components.

### 3. Why Ignore `tags` and `disable_auto_renew`?

**Context**: The DigitalOcean Pulumi provider (version 4.x) does not expose:
- A `tags` field for certificates (tags are supported for Droplets, but not certs)
- A `disable_auto_renew` field (Let's Encrypt renewal is always automatic)

**Solution**: Document these limitations in code comments and module README. Future provider updates may expose these fields.

```go
// NOTE: The DigitalOcean Pulumi provider currently lacks fields for tags
// and automatic‑renew configuration, so spec.tags and disable_auto_renew are ignored.
```

### 4. Why Export Both `certificate_id` and `expiry_rfc3339`?

**Rationale**:
- **`certificate_id`**: Required for integration with Load Balancers, Spaces CDN. This is the primary output.
- **`expiry_rfc3339`**: Enables expiration monitoring (e.g., parse in Prometheus, alert 30 days before expiry).

Both outputs are defined in `stack_outputs.proto`, ensuring consistency with the protobuf contract.

---

## Provider Integration

### DigitalOcean Provider Setup

```go
digitalOceanProvider, err := pulumidigitaloceanprovider.Get(
    ctx,
    stackInput.ProviderConfig,
)
```

**What this does:**
- Reads DigitalOcean API token from `stackInput.ProviderConfig.credential_id`
- Creates a configured `digitalocean.Provider` instance
- Handles token resolution from environment variables or Pulumi secrets

**Token Sources** (in order of precedence):
1. `DIGITALOCEAN_TOKEN` environment variable
2. Pulumi config: `pulumi config set digitalocean:token <token>`
3. Provider config in stack input (protobuf)

---

## Resource Lifecycle

### Let's Encrypt Certificate Lifecycle

```
┌──────────┐   DNS-01 Challenge   ┌──────────┐   Issued   ┌──────────┐
│ Pending  │ ──────────────────▶ │ Verified │ ────────▶ │  Active  │
└──────────┘                       └──────────┘            └──────────┘
     │                                                           │
     │ Timeout / DNS Failure                        Auto-renew (30 days before expiry)
     ▼                                                           │
┌──────────┐                                                     ▼
│  Error   │                                             ┌──────────┐
└──────────┘                                             │  Active  │
                                                          └──────────┘
```

**Key Points:**
- `pending` → `verified` typically takes 10-30 seconds
- Pulumi exports outputs immediately (certificate may not be usable until `verified`)
- Auto-renewal happens silently (certificate UUID doesn't change)

### Custom Certificate Lifecycle

```
┌──────────┐   Upload   ┌──────────┐   Manual Renewal   ┌──────────┐
│ Created  │ ─────────▶ │  Active  │ ─────────────────▶ │ Replaced │
└──────────┘             └──────────┘                     └──────────┘
                              │
                              │ Expiration (no auto-renewal)
                              ▼
                         ┌──────────┐
                         │ Expired  │
                         └──────────┘
```

**Key Points:**
- Custom certs have no lifecycle states (immediately active)
- **No auto-renewal**: You must rotate before expiration
- Pulumi's dependency engine handles zero-downtime rotation (create new → update refs → delete old)

---

## Error Handling

### Common Errors and Handling

1. **"certificate name already exists"**
   - **Cause**: Duplicate certificate name in the same DigitalOcean account
   - **Handling**: Module returns error, user must choose unique name
   - **Prevention**: Use Pulumi stack name as part of certificate name

2. **"domain validation failed" (Let's Encrypt)**
   - **Cause**: DNS not managed by DigitalOcean or DNS propagation delay
   - **Handling**: DigitalOcean retries DNS-01 challenge for ~5 minutes, then fails
   - **Prevention**: Pre-validate DNS zone exists in DigitalOcean

3. **"invalid PEM format" (Custom)**
   - **Cause**: Malformed PEM blocks (missing headers, invalid base64)
   - **Handling**: DigitalOcean API rejects immediately
   - **Prevention**: Validate PEM with `openssl` before uploading

### Module Error Propagation

```go
createdCertificate, err := digitalocean.NewCertificate(ctx, "certificate", certArgs, ...)
if err != nil {
    return nil, errors.Wrap(err, "failed to create digitalocean certificate")
}
```

Errors are wrapped with context for easier debugging.

---

## Testing Strategy

### Unit Tests

```go
// Example: Test discriminated union logic
func TestCertificateArgsBuilding(t *testing.T) {
    // Test 1: Let's Encrypt path
    spec := &DigitalOceanCertificateSpec{Type: letsEncrypt, ...}
    args := buildCertArgs(spec)
    assert.NotNil(t, args.Domains)
    assert.Nil(t, args.PrivateKey)

    // Test 2: Custom path
    spec = &DigitalOceanCertificateSpec{Type: custom, ...}
    args = buildCertArgs(spec)
    assert.Nil(t, args.Domains)
    assert.NotNil(t, args.PrivateKey)
}
```

### Integration Tests

```bash
# Test full deployment with test credentials
pulumi stack init test
pulumi config set digitalocean:token $TEST_DO_TOKEN
pulumi up --stack-input test-input.yaml
pulumi destroy --yes
```

---

## Performance Considerations

### Certificate Creation Timing

- **Let's Encrypt**: 10-30 seconds (DNS-01 challenge)
- **Custom**: < 1 second (no validation)

### State File Size

- **Minimal**: Certificate resources are lightweight (~1 KB in state)
- **Secrets**: Pulumi encrypts secrets (private keys) in state file

### DigitalOcean API Rate Limits

- **Certificate API**: 5,000 requests/hour (per account)
- **Module Impact**: Each deployment = 1 CREATE + 1 READ (~2 requests)
- **Best Practice**: Use Pulumi state locking to prevent concurrent deployments

---

## Future Enhancements

### 1. Support for Tags (when provider adds support)

```go
// Proposed change in certificate.go
if len(locals.DigitalOceanCertificate.Spec.Tags) > 0 {
    certArgs.Tags = convertToStringArray(locals.DigitalOceanCertificate.Spec.Tags)
}
```

### 2. Pre-Validation Checks

Add preflight checks before creating resources:

```go
// Check DNS zone exists (for Let's Encrypt)
if spec.Type == letsEncrypt {
    zone, err := digitalocean.GetDomain(ctx, extractApexDomain(spec.Domains[0]))
    if err != nil {
        return errors.Wrap(err, "DNS zone not found in DigitalOcean")
    }
}
```

### 3. Expiration Monitoring Integration

Export additional outputs for monitoring:

```go
ctx.Export("days_until_expiry", computeDaysUntilExpiry(cert.NotAfter))
ctx.Export("needs_renewal", pulumi.Bool(daysUntilExpiry < 30))
```

---

## Debugging Tips

### Enable Debug Logging

```bash
export PULUMI_DEBUG_COMMANDS=true
export PULUMI_DEBUG_GRPC=debug.log
pulumi up --logtostderr -v=9
```

### Inspect Pulumi State

```bash
# View full state (including secrets)
pulumi stack export > state.json

# View specific resource
pulumi stack --show-secrets | jq '.deployment.resources[] | select(.type=="digitalocean:index/certificate:Certificate")'
```

### Validate Provider Configuration

```bash
# Check DigitalOcean API connectivity
curl -H "Authorization: Bearer $DIGITALOCEAN_TOKEN" https://api.digitalocean.com/v2/account
```

---

## Comparison to Terraform Module

| Aspect | Pulumi Module | Terraform Module |
|--------|---------------|------------------|
| **Language** | Go (with protobuf) | HCL (declarative) |
| **State Management** | Pulumi state backend | Terraform state backend |
| **Secret Encryption** | Built-in (encrypted by default) | Manual (via backend config) |
| **Dependency Handling** | Implicit (DAG-based) | Explicit (`depends_on`) |
| **Lifecycle Hooks** | Go code (pre/post resource creation) | `lifecycle` blocks |
| **Testing** | Go unit tests + integration tests | `terraform plan` + Terratest |

**When to use Pulumi:**
- Need real programming logic (loops, conditionals, functions)
- Prefer type-safe languages (Go, TypeScript, Python)
- Want built-in secret encryption

**When to use Terraform:**
- Team prefers declarative HCL
- Mature ecosystem (community modules)
- Broader provider support

---

## Contributing

### Adding New Features

1. Update `spec.proto` (if changing API contract)
2. Regenerate Go stubs: `make protos`
3. Update `certificate.go` with new logic
4. Add tests in `certificate_test.go`
5. Update `README.md` and this `overview.md`

### Code Review Checklist

- [ ] Protobuf spec changes are backward-compatible
- [ ] Discriminated union logic is type-safe
- [ ] Errors are wrapped with context (`errors.Wrap`)
- [ ] Outputs are exported with correct keys (from `outputs.go`)
- [ ] Documentation updated (README.md, overview.md)

---

## References

- **DigitalOcean API Docs**: [Certificates API](https://docs.digitalocean.com/reference/api/api-reference/#tag/Certificates)
- **Pulumi DigitalOcean Provider**: [pulumi-digitalocean](https://www.pulumi.com/registry/packages/digitalocean/)
- **OpenMCF Conventions**: [architecture/deployment-component.md](../../../../../../architecture/deployment-component.md)
- **Protobuf Guide**: [buf.build](https://buf.build/)

---

**Maintainers**: For questions or clarifications, see the main [README.md](../../README.md) or file an issue.

