# Pulumi Module Overview - DigitalOcean Spaces Bucket

## Architecture

This document explains the internal architecture and design decisions of the Pulumi module for deploying DigitalOcean Spaces buckets.

## Module Structure

```
iac/pulumi/
├── main.go           # Pulumi program entrypoint
├── Pulumi.yaml       # Pulumi project configuration
├── Makefile          # Build automation
├── debug.sh          # Debug helper script
└── module/
    ├── main.go       # Module entry point (Resources function)
    ├── locals.go     # Local variable initialization
    ├── outputs.go    # Output key constants
    └── bucket.go     # Bucket resource creation logic
```

## Design Principles

### 1. Single Responsibility

Each file has a focused purpose:
- `main.go` (module): Orchestrates resource creation flow
- `bucket.go`: Spaces bucket resource creation
- `locals.go`: Data transformation and metadata extraction
- `outputs.go`: Output key constants for type safety

### 2. S3 API Compatibility

The module leverages DigitalOcean Spaces' S3 compatibility:
- Bucket names follow S3 DNS naming rules
- ACL values map to S3 ACL strings (`private`, `public-read`)
- Versioning semantics match S3 behavior

### 3. Protobuf-Driven Configuration

The module is driven by protobuf schema validation:
- Required fields enforced at API level (bucket_name, region)
- Enum values validated before reaching Pulumi
- Field patterns (bucket naming) validated by buf.validate

## Resource Creation Flow

### Step 1: Initialize Locals

```go
func initializeLocals(ctx *pulumi.Context, stackInput *DigitalOceanBucketStackInput) *Locals
```

Prepares:
- Metadata (name, labels, tags)
- Target resource reference
- Configuration lookups

### Step 2: Create Provider

```go
digitalOceanProvider, err := pulumidigitaloceanprovider.Get(ctx, stackInput.ProviderConfig)
```

Initializes DigitalOcean provider with:
- API token from provider config
- Spaces access credentials (if provided)
- Region settings

### Step 3: Map Access Control

```go
var acl pulumi.StringPtrInput
if locals.DigitalOceanBucket.Spec.AccessControl == digitaloceanbucketv1.DigitalOceanBucketAccessControl_PUBLIC_READ {
    acl = pulumi.String("public-read")
} else {
    acl = pulumi.String("private")
}
```

**Why**: Protobuf uses enums (PRIVATE=0, PUBLIC_READ=1), but DigitalOcean API expects S3 ACL strings.

### Step 4: Build Bucket Arguments

```go
bucketArgs := &digitalocean.SpacesBucketArgs{
    Name:   pulumi.String(locals.DigitalOceanBucket.Spec.BucketName),
    Region: pulumi.String(locals.DigitalOceanBucket.Spec.Region.String()),
    Acl:    acl,
}
```

### Step 5: Configure Versioning (Optional)

```go
if locals.DigitalOceanBucket.Spec.VersioningEnabled {
    bucketArgs.Versioning = &digitalocean.SpacesBucketVersioningArgs{
        Enabled: pulumi.Bool(true),
    }
}
```

**Design decision**: Only configure versioning block if enabled. If false or unset, omit the block entirely (cleaner resource configuration).

### Step 6: Create Bucket

```go
createdBucket, err := digitalocean.NewSpacesBucket(
    ctx,
    "bucket",
    bucketArgs,
    pulumi.Provider(digitalOceanProvider),
)
```

### Step 7: Export Outputs

```go
ctx.Export(OpBucketId, createdBucket.ID())
ctx.Export(OpEndpoint, createdBucket.Endpoint)
```

Outputs match `stack_outputs.proto`:
- `bucket_id`: Composite ID (`region:bucket-name`)
- `endpoint`: Bucket FQDN endpoint

## Key Design Decisions

### Why Enum for Access Control?

**Protobuf enum**:
```protobuf
enum DigitalOceanBucketAccessControl {
  PRIVATE = 0;
  PUBLIC_READ = 1;
}
```

**Benefits**:
- Type-safe at API level (can't pass invalid ACL strings)
- Clear semantics (PRIVATE vs PUBLIC_READ, not cryptic strings)
- Extensible (can add more ACL types without breaking existing specs)

**Trade-off**: Requires mapping to S3 ACL strings in implementation.

### Why Separate Versioning Boolean?

Instead of a versioning object with nested fields, we use a simple boolean:

**Rationale** (80/20 principle):
- 95% of users either enable versioning (true) or don't (false)
- Advanced versioning features (MFA delete, suspend state) are rarely used
- Simplifies API surface for common case

**Future**: Could extend with `DigitalOceanBucketVersioningConfig` message if advanced features needed.

### Why No CORS Configuration in Spec?

CORS is intentionally omitted from the protobuf spec:

**Reason**:
- CORS configuration is highly application-specific (allowed origins, methods, headers)
- Most users configure CORS via S3 API after bucket creation
- Adding to spec would significantly increase complexity for ~20% use case

**Workaround**: Users can configure CORS using AWS SDK, Terraform/Pulumi resources, or DigitalOcean control panel after bucket creation.

### Why Tags as Repeated String?

Tags are simple strings instead of key-value pairs:

**Rationale**:
- DigitalOcean Spaces tags are label-style (not AWS-style key=value)
- Simpler API surface (just a list of strings)
- Consistent with DigitalOcean's tagging model

Example: `["production", "team:backend", "purpose:media"]`

## Error Handling

The module uses Go error wrapping for context:

```go
if err != nil {
    return errors.Wrap(err, "failed to create digitalocean spaces bucket")
}
```

This provides:
- Stack traces for debugging
- Contextual error messages
- Pulumi error reporting integration

## State Management

Pulumi automatically tracks:
- Bucket ID (composite: `region:bucket-name`)
- Bucket endpoint
- Configuration state
- Resource dependencies

On subsequent `pulumi up`:
- Pulumi compares desired state (spec) with actual state (DigitalOcean API)
- Generates minimal diff
- Updates only changed configuration

**Note**: Some changes (like bucket name or region) require resource replacement (destroy + recreate).

## Provider Configuration

The module uses a shared provider helper:

```go
digitalOceanProvider, err := pulumidigitaloceanprovider.Get(ctx, stackInput.ProviderConfig)
```

This centralizes:
- Credential management
- Provider options
- Retry configuration
- Common settings

## Outputs Structure

Constants defined in `outputs.go`:

```go
const (
    OpBucketId = "bucket_id"
    OpEndpoint = "endpoint"
)
```

**Why constants**: Type safety and refactoring support. Changing output keys requires updating constant definitions.

## Future Enhancements

Potential additions (not in current 80/20 scope):
- CORS configuration in spec
- Lifecycle policy rules
- Object lock configuration
- Bucket logging settings
- CDN endpoint customization
- Transfer acceleration
- Bucket policies (beyond ACL)

## Testing Strategy

While this module lacks unit tests, recommended testing:
1. **Integration tests**: Deploy to test DO account
2. **Pulumi preview**: Review changes before applying
3. **Staging buckets**: Test on non-production stacks
4. **S3 compatibility**: Test with AWS CLI/SDK

## Performance Considerations

The module is optimized for:
- **Fast previews**: Minimal API calls during `pulumi preview`
- **Efficient updates**: Only modified fields trigger updates
- **Idempotent operations**: Safe to run multiple times

## Bucket Naming Conventions

The module enforces DNS-compatible bucket names:
- 3-63 characters
- Lowercase alphanumeric and hyphens
- Cannot start or end with hyphen
- No underscores or special characters

**Valid**: `my-app-media`, `prod-backups-2024`
**Invalid**: `My_App`, `-invalid-`, `a`, `TOO_MANY_UNDERSCORES`

## Region Enum Mapping

Protobuf region enums map to region slugs:
- `nyc3` → `"nyc3"`
- `sfo3` → `"sfo3"`
- `ams3` → `"ams3"`

The `.String()` method handles conversion.

## S3 Endpoint Construction

For each bucket, DigitalOcean provides:
- **Origin endpoint**: `{bucket-name}.{region}.digitaloceanspaces.com`
- **CDN endpoint**: `{bucket-name}.{region}.cdn.digitaloceanspaces.com` (for public buckets)

The module exports the origin endpoint. CDN endpoint can be constructed by clients.

## References

- [Pulumi Programming Model](https://www.pulumi.com/docs/intro/concepts/programming-model/)
- [DigitalOcean Provider SDK](https://github.com/pulumi/pulumi-digitalocean)
- [Spaces API Documentation](https://docs.digitalocean.com/reference/api/spaces-api/)
- [OpenMCF Architecture](../../../../../../architecture/deployment-component.md)

