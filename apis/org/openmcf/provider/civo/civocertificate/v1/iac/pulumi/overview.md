# CivoCertificate Pulumi Module Architecture

## High-Level Overview

The CivoCertificate Pulumi module provides validation and structure for TLS certificate specifications on Civo Cloud. Due to current provider limitations, the module cannot provision actual certificates but validates configurations and logs appropriate guidance.

**Provider Status**: The Civo Pulumi/Terraform provider does not expose certificate resources as of 2025. This module is prepared for future provider support.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                CivoCertificate Manifest (YAML)                  │
│                                                                  │
│  apiVersion: civo.openmcf.org/v1                        │
│  kind: CivoCertificate                                          │
│  spec:                                                          │
│    certificate_name: prod-cert                                  │
│    type: letsEncrypt                                            │
│    lets_encrypt:                                                │
│      domains: [example.com, www.example.com]                   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                  Protobuf Validation (buf.validate)             │
│                                                                  │
│  ✓ certificate_name: 1-64 chars                                │
│  ✓ type: letsEncrypt or custom                                 │
│  ✓ domains: valid FQDN/wildcard pattern                        │
│  ✓ oneof: exactly one of lets_encrypt or custom               │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Pulumi Module Entry                        │
│                                                                  │
│  1. Initialize Locals (metadata, labels, tags)                 │
│  2. Validate certificate configuration                          │
│  3. Log provider limitation warning                             │
│  4. Log certificate details (type, domains, auto-renew)        │
│  5. Export placeholder outputs                                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Current State (Logs)                        │
│                                                                  │
│  ⚠️ Provider does not support certificate resources              │
│  ℹ️ Let's Encrypt cert for: [example.com, www.example.com]      │
│  ℹ️ Auto-renewal: enabled                                        │
│  📤 Outputs: certificate_id="", expiry_rfc3339=""               │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Manual Provisioning Required                  │
│                                                                  │
│  Options:                                                       │
│  • Civo Dashboard: dashboard.civo.com/certificates             │
│  • Civo API: curl api.civo.com/v2/certificates                 │
│  • Civo CLI: civo certificate create                           │
└─────────────────────────────────────────────────────────────────┘
```

## Module Components

### 1. Entry Point (`main.go`)

Pulumi stack initialization. Deserializes `CivoCertificateStackInput` and calls `module.Resources()`.

**Flow**:
```go
pulumi.Run(func(ctx *pulumi.Context) error {
    stackInput := parseStackInput()  // From env or file
    return module.Resources(ctx, stackInput)
})
```

### 2. Module Entry (`module/main.go`)

Core logic:
1. Initialize locals (metadata, labels)
2. Call `certificate()` for validation
3. Return errors or success

### 3. Locals Initialization (`module/locals.go`)

Prepares:
- `CivoCertificate`: Spec from input
- `Metadata`: Name, labels, description
- `Labels`: Merged metadata labels + spec tags

### 4. Certificate Validation (`module/certificate.go`)

Current behavior:
- Logs **warning** about provider limitation
- Logs **info** about certificate type and configuration
- Exports **empty** placeholder outputs

Future behavior (when provider support is added):
- Create `civo.Certificate` resource
- Configure Let's Encrypt or custom params
- Export real certificate ID and expiry

### 5. Outputs (`module/outputs.go`)

Constants for output keys:
- `OpCertificateId` = "certificate_id"
- `OpExpiryRfc3339` = "expiry_rfc3339"

## Data Flow

### Input Flow

```
YAML Manifest
    │
    ▼
Protobuf Deserialization
    │
    ▼
CivoCertificateStackInput
    │
    ├─► CivoCertificate
    │       ├─► certificate_name
    │       ├─► type (letsEncrypt/custom)
    │       ├─► lets_encrypt/custom params
    │       ├─► description
    │       └─► tags
    │
    └─► ProviderConfig (future use)
            └─► civo_token
```

### Processing Flow

```
module.Resources()
    │
    ├─► initializeLocals()
    │       └─► Locals{metadata, labels}
    │
    └─► certificate()
            ├─► Log warning (provider limitation)
            ├─► Log certificate details
            └─► Export placeholder outputs
```

### Output Flow

```
Stack Outputs (Current)
    │
    ├─► certificate_id: ""
    └─► expiry_rfc3339: ""

Stack Outputs (Future)
    │
    ├─► certificate_id: "uuid-1234-5678"
    └─► expiry_rfc3339: "2025-04-15T10:30:00Z"
```

## Certificate Types

### Let's Encrypt

**Validation**:
- `domains`: Required, unique, valid FQDN/wildcard pattern
- `disable_auto_renew`: Optional boolean

**Example**:
```yaml
lets_encrypt:
  domains:
    - "*.example.com"
    - "example.com"
  disable_auto_renew: false  # Default
```

**Logged Output**:
```
ℹ️  Let's Encrypt certificate requested for domains: [*.example.com example.com] (auto-renew: true)
```

### Custom Certificate

**Validation**:
- `leaf_certificate`: Required, non-empty
- `private_key`: Required, non-empty
- `certificate_chain`: Optional but recommended

**Example**:
```yaml
custom:
  leaf_certificate: |
    -----BEGIN CERTIFICATE-----
    MIIFbjCCBFagAwIBAgIQDvO...
    -----END CERTIFICATE-----
  private_key: |
    -----BEGIN RSA PRIVATE KEY-----
    MIIEpAIBAAKCAQEAv0...
    -----END RSA PRIVATE KEY-----
  certificate_chain: |
    -----BEGIN CERTIFICATE-----
    (intermediate certs)
    -----END CERTIFICATE-----
```

**Logged Output**:
```
ℹ️  Custom certificate provided for 'prod-cert'
```

## Validation Rules

All validation is enforced at the Protobuf level via `buf.validate`:

| Field | Rule | Error |
|-------|------|-------|
| certificate_name | required, 1-64 chars | "certificate_name is required" |
| type | required, defined_only | "type must be letsEncrypt or custom" |
| certificate_source | oneof required | "must specify lets_encrypt or custom" |
| domains (LE) | required, unique, pattern | "invalid domain format" |
| leaf_certificate (custom) | required, min_len=1 | "leaf_certificate is required" |
| private_key (custom) | required, min_len=1 | "private_key is required" |
| description | max_len=128 | "description too long" |
| tags | unique | "duplicate tags" |

## Error Handling

### Validation Errors

Caught by `protovalidate.Validate()` before reaching Pulumi:

```go
err := protovalidate.Validate(input)
if err != nil {
    // Error contains field-level details
    return fmt.Errorf("validation failed: %w", err)
}
```

### Runtime Errors

Currently minimal (no network calls). Future errors:
- Provider authentication failure
- Rate limiting (Let's Encrypt)
- Domain validation failure
- Invalid certificate format

## Provider Limitation Details

### Current State

**What works**:
- Specification validation ✅
- Configuration logging ✅
- Module structure ✅
- Test coverage ✅

**What doesn't work**:
- Actual certificate provisioning ❌
- Certificate renewal ❌
- Certificate deletion ❌
- Real outputs (ID, expiry) ❌

### Why This Limitation Exists

The Civo provider (both Terraform and Pulumi) does not expose certificate resources. Checked:
- https://registry.terraform.io/providers/civo/civo/latest/docs (no `civo_certificate`)
- https://www.pulumi.com/registry/packages/civo/api-docs/ (no `Certificate` resource)

This is a **provider gap**, not a Civo API limitation. The Civo REST API supports certificates:
- https://www.civo.com/api/certificates

### When Support Will Be Added

Monitor:
- **Terraform Provider Changelog**: https://github.com/civo/terraform-provider-civo/releases
- **Pulumi Provider Updates**: https://github.com/pulumi/pulumi-civo/releases

Once Terraform provider adds `civo_certificate`, Pulumi will automatically bridge it.

## Future Enhancements

### When Provider Support is Added

**Module Updates Required**:
1. Import Civo provider package
2. Implement `civo.NewCertificate()` in `certificate.go`
3. Handle Let's Encrypt vs custom conditional logic
4. Export real certificate ID and expiry
5. Update tests with actual provisioning
6. Add integration tests against Civo staging

**Example Future Implementation**:
```go
cert, err := civo.NewCertificate(ctx, "cert", &civo.CertificateArgs{
    Name:    pulumi.String(spec.CertificateName),
    Domains: pulumi.ToStringArray(spec.GetLetsEncrypt().Domains),
}, pulumi.Provider(civoProvider))

ctx.Export(OpCertificateId, cert.ID())
ctx.Export(OpExpiryRfc3339, cert.Expiry)
```

### Potential Features

- **Auto-renewal monitoring**: Poll certificate expiry, alert before expiration
- **Certificate rotation**: Automate custom cert updates
- **Multi-region support**: Deploy same cert across regions
- **Integration with cert-manager**: Sync with Kubernetes cert-manager

## Security Considerations

### Current (Validation Only)

- **Secrets in specs**: Private keys stored in YAML (use Pulumi secrets)
- **Git commits**: Never commit PEM files (use `.gitignore`)
- **Validation only**: No network exposure, no security risk

### Future (When Provisioning)

- **Pulumi secrets**: Encrypt private keys in state
- **Least privilege**: Separate Civo API tokens per environment
- **Certificate rotation**: Automate renewal before expiry
- **Audit logging**: Track certificate creation/deletion

## Testing Strategy

### Current Tests

- **spec_test.go**: Validates all buf.validate rules
  - Certificate name length
  - Type enumeration
  - Domain patterns (FQDN, wildcard)
  - Custom cert requirements
  - Tags uniqueness
  - Description length

### Future Tests

When provider support is added:
- **Integration tests**: Provision real certificates
- **E2E tests**: Attach to load balancers, verify HTTPS
- **Renewal tests**: Verify auto-renewal workflow
- **Cleanup tests**: Ensure certificate deletion works

## Performance Characteristics

### Current

- **Validation**: < 1ms (no network)
- **Module execution**: < 100ms (logging only)

### Future

- **Let's Encrypt provisioning**: 30-90 seconds (ACME challenge)
- **Custom cert upload**: 5-10 seconds
- **Certificate renewal**: ~60 seconds

## Related Resources

- **Civo API Docs**: https://www.civo.com/api/certificates
- **Let's Encrypt**: https://letsencrypt.org/docs/
- **Pulumi Civo Provider**: https://www.pulumi.com/registry/packages/civo/
- **ACME Protocol**: https://tools.ietf.org/html/rfc8555

## Conclusion

The CivoCertificate Pulumi module provides a production-ready structure for certificate management on Civo. While current provider limitations prevent actual provisioning, the module validates configurations, logs actionable guidance, and is prepared for seamless certificate deployment once provider support is added.

For now, use the module for validation and leverage manual provisioning via Civo dashboard or API.
