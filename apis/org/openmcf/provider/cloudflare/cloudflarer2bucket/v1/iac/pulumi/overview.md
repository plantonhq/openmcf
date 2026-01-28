# Pulumi Module Overview: Cloudflare R2 Bucket

## Architecture

This document explains the architectural decisions, resource model, and design rationale for the Cloudflare R2 Bucket Pulumi module.

## Resource Model

The Cloudflare R2 Bucket module provisions the bucket and optionally a custom domain:

```
┌─────────────────────────────────────────────────────────────┐
│                    Cloudflare Account                        │
│                                                               │
│  ┌─────────────────────────┐    ┌─────────────────────────┐ │
│  │ R2 Bucket               │    │ R2 Custom Domain        │ │
│  │ - Name: string          │───▶│ - Hostname: string      │ │
│  │ - Location: enum        │    │ - ZoneId: string        │ │
│  │ - Public: bool          │    │ (optional)              │ │
│  └─────────────────────────┘    └─────────────────────────┘ │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

### Why So Simple?

R2 buckets are intentionally minimal compared to AWS S3:
- **No versioning** (R2 doesn't support it)
- **No bucket policies** (use R2 API tokens instead)
- **No lifecycle rules via API** (coming in future)
- **No website hosting** (use Cloudflare Pages)
- **Custom domains** (now supported via IaC for production access)

This simplicity is by design—R2 focuses on the 80% use case: store and serve objects efficiently, without egress fees.

## Implementation Details

### File Structure

```
module/
├── main.go           # Entry point - Resources() function
├── locals.go         # Local variables (metadata, spec shortcuts)
├── bucket.go         # Core logic (creates R2 bucket)
└── outputs.go        # Output constant definitions
```

### Function Flow

**1. `main.go:Resources()`**

- Entry point called by Pulumi CLI
- Initializes locals (shortcuts for metadata and spec)
- Creates Cloudflare provider
- Calls `bucket()` to create the R2 bucket

**2. `bucket.go:bucket()`**

This is the core implementation:

```go
func bucket(ctx *pulumi.Context, locals *Locals, provider *cloudflare.Provider) {
    // Step 1: Translate enum to string
    var bucketLocation string
    switch locals.CloudflareR2Bucket.Spec.Location {
    case cloudflarer2bucketv1.CloudflareR2Location_WNAM:
        bucketLocation = "WNAM"
    case cloudflarer2bucketv1.CloudflareR2Location_ENAM:
        bucketLocation = "ENAM"
    // ... other cases
    default:
        bucketLocation = "auto"
    }

    // Step 2: Create the bucket
    bucket, err := cloudflare.NewR2Bucket(ctx, "bucket", &cloudflare.R2BucketArgs{
        AccountId: pulumi.String(locals.CloudflareR2Bucket.Spec.AccountId),
        Name:      pulumi.String(locals.CloudflareR2Bucket.Spec.BucketName),
        Location:  pulumi.String(bucketLocation),
    })

    // Step 3: Handle limitations
    if locals.CloudflareR2Bucket.Spec.PublicAccess {
        ctx.Log.Warn("Public access must be enabled manually - not yet in provider", nil)
    }
    
    if locals.CloudflareR2Bucket.Spec.VersioningEnabled {
        ctx.Log.Warn("R2 does not support versioning - field ignored", nil)
    }

    // Step 4: Create custom domain (if enabled)
    if spec.CustomDomain != nil && spec.CustomDomain.Enabled {
        cloudflare.NewR2CustomDomain(ctx, "custom-domain", &cloudflare.R2CustomDomainArgs{
            AccountId:  pulumi.String(spec.AccountId),
            BucketName: bucket.Name,
            ZoneId:     pulumi.String(spec.CustomDomain.ZoneId.GetValue()),
            Domain:     pulumi.String(spec.CustomDomain.Domain),
            Enabled:    pulumi.Bool(true),
        })
        ctx.Export(OpCustomDomainUrl, pulumi.Sprintf("https://%s", spec.CustomDomain.Domain))
    }

    // Step 5: Export outputs
    ctx.Export(OpBucketName, bucket.Name)
}
```

### Enum to String Mapping

Protobuf enums are integers, but Cloudflare expects string values. The module translates:

```go
// Location enum values to Cloudflare location strings
enum (proto) → string (API)
0 (UNSPECIFIED) → "auto"
1 (WNAM) → "WNAM"
2 (ENAM) → "ENAM"
3 (WEUR) → "WEUR"
4 (EEUR) → "EEUR"
5 (APAC) → "APAC"
6 (OC) → "OC"
```

This abstraction allows the proto API to use strongly-typed enums while the Pulumi module handles Cloudflare's string-based API.

## Design Decisions

### Decision 1: Minimal Resource Scope

**Choice**: Only implement bucket creation with name, account, and location.

**Rationale**:
- R2 is intentionally simpler than S3
- Advanced features (CORS, lifecycle) require AWS S3-compatible API
- Adding S3 provider would complicate the module significantly
- 80/20 principle: Most users need basic bucket creation

**Trade-off**:
- ✅ Simple, maintainable module
- ❌ Advanced features require manual configuration

**Future Enhancement**: Add support for CORS/lifecycle when Cloudflare provider adds native support.

### Decision 2: Warning for Unsupported Fields

**Choice**: Log warnings for `public_access` and `versioning_enabled` instead of failing.

**Rationale**:
- `public_access`: Feature exists in R2 but not yet in Pulumi provider
- `versioning_enabled`: R2 doesn't support it (may never support it)
- Warnings inform users without blocking deployment

**Trade-off**:
- ✅ Graceful degradation
- ❌ User might not notice warnings

**Best Practice**: Documentation prominently explains these limitations.

### Decision 3: Location as Enum

**Choice**: Use protobuf enum for location, not freeform string.

**Rationale**:
- Type safety: Only valid locations can be specified
- Validation at proto level prevents invalid API calls
- Better IDE autocomplete and documentation

**Trade-off**:
- ✅ Compile-time validation
- ❌ Requires proto update if Cloudflare adds new locations

**Alternative Considered**: String field with validation regex. Rejected because enum is more maintainable.

### Decision 4: No CORS or Lifecycle

**Choice**: Don't implement CORS or lifecycle rules.

**Rationale**:
- Requires AWS S3 provider (complex setup)
- Not universally needed (80/20 split)
- Users can configure manually via S3 API if needed

**Trade-off**:
- ✅ Simple implementation
- ❌ Power users need manual steps

**Workaround**: Document how to configure CORS manually in research docs.

## OpenMCF Abstraction

### The Core Problem

When using the Cloudflare API or Pulumi directly, users must:
1. Understand Cloudflare's account structure
2. Know location strings exactly ("WNAM" not "wnam")
3. Handle provider configuration
4. Remember that R2 doesn't support versioning

### The Solution: Typed, Validated API

OpenMCF provides:

```yaml
# User's input (simple and validated)
spec:
  bucket_name: my-bucket
  account_id: "a1b2c3d4e5f6..."
  location: WEUR  # Enum, not string
  public_access: false
```

Behind the scenes, the Pulumi module:
1. Validates bucket_name format (DNS-compatible, 3-63 chars)
2. Validates account_id length (exactly 32 hex chars)
3. Translates location enum to string
4. Creates Cloudflare provider
5. Provisions R2 bucket
6. Warns about unsupported features

**User sees**: Simple, typed API with validation  
**Module handles**: Provider setup, enum translation, error handling

## Secret Management

### API Token Handling

Cloudflare API tokens are highly sensitive. The module supports:

**1. Environment Variables (Recommended)**

```bash
export CLOUDFLARE_API_TOKEN="your-token"
pulumi up
```

The Pulumi Cloudflare provider automatically reads `CLOUDFLARE_API_TOKEN`.

**2. Pulumi Config Secrets**

```bash
pulumi config set cloudflare:apiToken --secret "your-token"
pulumi up
```

Pulumi encrypts the secret before storing it in the state file.

**Why environment variables are preferred**:
- Secrets never touch the state file
- Compatible with CI/CD systems (GitHub Actions secrets, etc.)
- Follows 12-factor app principles

## Error Handling

The module uses Go's standard error handling pattern:

```go
bucket, err := cloudflare.NewR2Bucket(...)
if err != nil {
    return errors.Wrap(err, "failed to create Cloudflare R2 bucket")
}
```

Pulumi captures these errors and displays them with context:

```
error: failed to create R2 bucket: API Error:
  bucket name 'invalid..name' contains invalid characters
```

## Outputs

The module exports the following outputs:

| Output | Type | Description |
|--------|------|-------------|
| `bucket_name` | string | The name of the created R2 bucket |
| `custom_domain_url` | string | The custom domain URL (only if custom domain is enabled) |

Access outputs:

```bash
# Get bucket name
BUCKET_NAME=$(pulumi stack output bucket_name)

# Get custom domain URL (if configured)
CUSTOM_DOMAIN_URL=$(pulumi stack output custom_domain_url)

# Use in other Pulumi stacks via StackReference
bucket_name = pulumi.StackReference("org/stack-name").get_output("bucket_name")
```

## Provider Configuration

The module uses the `pulumicloudflareprovider.Get()` helper to create a Cloudflare provider:

```go
cloudflareProvider, err := pulumicloudflareprovider.Get(ctx, stackInput.ProviderConfig)
```

This helper:
1. Extracts API token from `ProviderConfig` or environment variables
2. Creates a Cloudflare provider instance
3. Handles authentication errors gracefully

All resources are created using this provider:

```go
bucket, err := cloudflare.NewR2Bucket(ctx, "bucket", &args,
    pulumi.Provider(cloudflareProvider),  // Use explicit provider
)
```

**Why explicit providers?** Allows multiple Cloudflare accounts in the same Pulumi program (e.g., dev and prod accounts).

## Testing Strategy

### Unit Tests

The module includes unit tests for:
- Enum-to-string conversion logic
- Local variable initialization
- Error handling paths

Run tests:

```bash
cd module
go test -v ./...
```

### Integration Tests

Integration tests deploy real R2 buckets (requires API token):

```bash
export CLOUDFLARE_API_TOKEN="..."
export CLOUDFLARE_ACCOUNT_ID="..."
go test -v -tags=integration ./...
```

**Cost**: Each integration test run costs ~$0 (R2 free tier covers testing).

## Performance Considerations

### Resource Creation Time

Typical deployment timeline:
- Bucket creation: ~2-3 seconds
- **Total**: 2-3 seconds

R2 bucket creation is fast because:
- No complex dependencies
- Single API call
- No health checks required

### State File Size

Each R2 bucket adds ~500 bytes to the Pulumi state file:
- Bucket: ~500 bytes (minimal metadata)
- **Total**: ~500 bytes per bucket

## Limitations and Workarounds

### Limitation 1: No Public Access Toggle

**Issue**: Pulumi provider doesn't expose r2.dev public URL toggle.

**Workaround**:
1. Deploy bucket with Pulumi
2. Enable public access manually via Cloudflare Dashboard
3. Or use Cloudflare API directly

**Future**: Wait for Pulumi provider to add support.

### Limitation 2: No CORS Configuration

**Issue**: CORS requires AWS S3 provider setup.

**Workaround**:
1. Deploy bucket with Pulumi
2. Configure CORS manually using AWS CLI with R2 endpoint:
   ```bash
   aws s3api put-bucket-cors --bucket my-bucket \
     --cors-configuration file://cors.json \
     --endpoint-url https://<account-id>.r2.cloudflarestorage.com
   ```

**Future**: Add CORS support when widely requested.

### Limitation 3: No Versioning

**Issue**: R2 doesn't support object versioning.

**Workaround**: Implement application-level versioning (e.g., append timestamps to object keys).

**Future**: R2 may add versioning in the future.

## Comparison: Pulumi vs. Terraform

| Aspect | Pulumi (This Module) | Terraform |
|--------|----------------------|-----------|
| **Language** | Go (strongly typed) | HCL (declarative) |
| **Secret Management** | Built-in encryption | Requires external vault |
| **Type Safety** | Compile-time checks | Runtime errors |
| **Conditionals** | Go if/else | HCL count/for_each |
| **Reusability** | Go functions | Terraform modules |
| **Learning Curve** | Moderate (requires Go knowledge) | Low (HCL is simple) |

**Recommendation**:
- Use Pulumi if you prefer type-safe, general-purpose languages
- Use Terraform if you prefer declarative DSLs and have an existing Terraform codebase

## Future Enhancements

### 1. Public Access Toggle

When Pulumi provider adds support:

```go
bucket, err := cloudflare.NewR2Bucket(ctx, "bucket", &cloudflare.R2BucketArgs{
    // ... existing fields
    PublicAccess: pulumi.Bool(locals.CloudflareR2Bucket.Spec.PublicAccess),
})
```

### 2. CORS Configuration

Add optional CORS rules in spec.proto:

```protobuf
message CloudflareR2BucketSpec {
  // ... existing fields
  repeated CorsRule cors_rules = 6;
}

message CorsRule {
  repeated string allowed_origins = 1;
  repeated string allowed_methods = 2;
  repeated string allowed_headers = 3;
  int32 max_age_seconds = 4;
}
```

### 3. Lifecycle Policies

Add optional lifecycle rules:

```protobuf
message CloudflareR2BucketSpec {
  // ... existing fields
  repeated LifecycleRule lifecycle_rules = 7;
}

message LifecycleRule {
  int32 expiration_days = 1;
  string prefix = 2;
}
```

### 4. Custom Domain Attachment

Automatically configure custom domains:

```protobuf
message CloudflareR2BucketSpec {
  // ... existing fields
  repeated string custom_domains = 8;
}
```

## References

- [Cloudflare R2 API Docs](https://developers.cloudflare.com/api/operations/r2-create-bucket)
- [Pulumi Cloudflare Provider](https://www.pulumi.com/registry/packages/cloudflare/)
- [Component README](../../README.md)
- [Research Documentation](../../docs/README.md)

---

**Questions?** Review the [README.md](./README.md) or consult Cloudflare's API documentation.

