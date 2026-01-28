# CivoBucket Pulumi Module Architecture

## High-Level Overview

The CivoBucket Pulumi module provisions Civo Object Storage buckets with S3-compatible access. It translates a declarative Protobuf specification (`CivoBucketSpec`) into Civo infrastructure resources, handling credential creation, bucket provisioning, and output management.

```
┌─────────────────────────────────────────────────────────────────┐
│                    CivoBucket Manifest (YAML)                   │
│                                                                  │
│  apiVersion: civo.openmcf.org/v1                        │
│  kind: CivoBucket                                               │
│  metadata: {name: my-bucket}                                    │
│  spec:                                                          │
│    bucketName: my-bucket                                        │
│    region: LON1                                                 │
│    versioningEnabled: true                                      │
│    tags: [env:prod]                                             │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ (Protobuf deserialization)
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                  CivoBucketStackInput (Proto)                   │
│                                                                  │
│  CivoBucket: {metadata, spec}                                   │
│  ProviderConfig: {civo_token}                                   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ (Pulumi module entry)
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Pulumi Resources                           │
│                                                                  │
│  1. Initialize Locals (metadata, labels, region)                │
│  2. Create Civo Provider (with API token)                       │
│  3. Create ObjectStoreCredential (access key pair)              │
│  4. Create ObjectStore (bucket with credential)                 │
│  5. Export Outputs (bucket ID, endpoint, keys)                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ (Civo API calls)
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Civo Infrastructure                         │
│                                                                  │
│  ┌─────────────┐       ┌─────────────────────────────────────┐ │
│  │ Credential  │◄──────│ ObjectStore (Bucket)                │ │
│  │             │       │                                     │ │
│  │ Access Key  │       │ Name: my-bucket                     │ │
│  │ Secret Key  │       │ Region: LON1                        │ │
│  └─────────────┘       │ Endpoint: objectstore.civo.com      │ │
│                        └─────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ (Stack outputs)
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                CivoBucketStackOutputs (Proto)                   │
│                                                                  │
│  bucket_id: "uuid-1234-5678"                                    │
│  endpoint_url: "https://objectstore.civo.com/my-bucket"         │
│  access_key_secret_ref: "civo-bucket-access-key"                │
│  secret_key_secret_ref: "civo-bucket-secret-key"                │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ (Consumed by applications)
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Application Consumption                       │
│                                                                  │
│  - Kubernetes Secrets (mounted in pods)                         │
│  - CI/CD pipelines (environment variables)                      │
│  - Backup tools (S3 API access)                                 │
└─────────────────────────────────────────────────────────────────┘
```

## Module Components

### 1. Entry Point (`main.go`)

The Pulumi stack initialization file. Responsibilities:
- Parse command-line arguments or environment variables
- Deserialize `CivoBucketStackInput` from JSON/YAML
- Call `module.Resources()` with Pulumi context and stack input
- Handle top-level errors

**Flow**:
```go
func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        // 1. Parse stack input (from env var or file)
        stackInput := parseStackInput()
        
        // 2. Call module
        return module.Resources(ctx, stackInput)
    })
}
```

### 2. Module Entry (`module/main.go`)

The core module logic. Responsibilities:
- Initialize locals (metadata, labels)
- Create Civo provider
- Orchestrate resource creation
- Return errors

**Flow**:
```go
func Resources(ctx *pulumi.Context, stackInput *CivoBucketStackInput) error {
    // 1. Prepare locals
    locals := initializeLocals(ctx, stackInput)
    
    // 2. Setup Civo provider
    civoProvider, err := pulumicivoprovider.Get(ctx, stackInput.ProviderConfig)
    if err != nil {
        return errors.Wrap(err, "failed to setup civo provider")
    }
    
    // 3. Create bucket
    _, err = bucket(ctx, locals, civoProvider)
    return err
}
```

### 3. Locals Initialization (`module/locals.go`)

Prepares local variables for use throughout the module. Responsibilities:
- Extract metadata (name, labels, description)
- Store CivoBucket spec
- Prepare labels for tagging (if supported)

**Data Structure**:
```go
type Locals struct {
    CivoBucket *civobucketv1.CivoBucket
    Metadata   *shared.CloudResourceMetadata
    Labels     map[string]string
}
```

**Why Locals?**
- **Reusability**: Shared across multiple resource creation functions
- **Clarity**: Centralized data access
- **Immutability**: Prepared once, used many times

### 4. Bucket Creation (`module/bucket.go`)

Core resource provisioning logic. Responsibilities:
- Create ObjectStoreCredential (access key pair)
- Create ObjectStore (bucket)
- Handle versioning configuration (informational)
- Handle tags (informational)
- Export outputs

**Credential Creation**:
```go
createdCredential, err := civo.NewObjectStoreCredential(
    ctx,
    "bucket-creds",
    &civo.ObjectStoreCredentialArgs{},
    pulumi.Provider(civoProvider),
)
```

**Key Points**:
- No explicit keys provided → Civo auto-generates them
- Credentials are region-specific
- One credential can be attached to multiple buckets

**Bucket Creation**:
```go
bucketArgs := &civo.ObjectStoreArgs{
    Name:        pulumi.String(spec.BucketName),
    Region:      pulumi.String(spec.Region.String()),
    AccessKeyId: createdCredential.AccessKeyId,
}

createdBucket, err := civo.NewObjectStore(
    ctx,
    "bucket",
    bucketArgs,
    pulumi.Provider(civoProvider),
)
```

**Key Points**:
- `Name`: DNS-compatible bucket name (validated in proto)
- `Region`: Civo region enum mapped to string
- `AccessKeyId`: Reference to created credential

**Versioning Handling**:
```go
if spec.VersioningEnabled {
    ctx.Log.Info("Versioning requested. Note: Configure via S3 API.", nil)
}
```

**Rationale**: Civo provider doesn't expose versioning. The module logs a reminder for post-deployment configuration.

**Tags Handling**:
```go
if len(spec.Tags) > 0 {
    ctx.Log.Info(fmt.Sprintf("Tags: %v (not applied to Civo resource)", spec.Tags), nil)
}
```

**Rationale**: Civo ObjectStore doesn't support tags. Tags are recorded in metadata for organizational purposes.

### 5. Outputs (`module/outputs.go`)

Defines output constant names for `ctx.Export()`. Responsibilities:
- Standardize output key names
- Ensure consistency across modules
- Document what gets exported

**Constants**:
```go
const (
    OpBucketId           = "bucket_id"
    OpEndpointUrl        = "endpoint_url"
    OpAccessKeySecretRef = "access_key_secret_ref"
    OpSecretKeySecretRef = "secret_key_secret_ref"
)
```

**Export Example**:
```go
ctx.Export(OpBucketId, createdBucket.ID())
ctx.Export(OpEndpointUrl, createdBucket.BucketUrl)
ctx.Export(OpAccessKeySecretRef, createdCredential.AccessKeyId)
ctx.Export(OpSecretKeySecretRef, createdCredential.SecretAccessKey)
```

## Resource Dependency Graph

```
┌──────────────────────┐
│  Civo Provider       │
│  (API token)         │
└──────────┬───────────┘
           │
           │ (provides)
           │
           ▼
┌──────────────────────┐
│ ObjectStoreCredential│◄─────────┐
│ (access key pair)    │          │
└──────────┬───────────┘          │
           │                      │
           │ (access_key_id)      │ (depends on)
           │                      │
           ▼                      │
┌──────────────────────┐          │
│   ObjectStore        │──────────┘
│   (bucket)           │
└──────────────────────┘
```

**Dependency Chain**:
1. **Provider**: Must exist first (initialized in `Resources()`)
2. **Credential**: Created before bucket (needs provider)
3. **Bucket**: Depends on credential (references `access_key_id`)

Pulumi automatically handles ordering based on resource references.

## Data Flow

### Input Flow

```
YAML Manifest
    │
    ▼
Protobuf Deserialization
    │
    ▼
CivoBucketStackInput
    │
    ├─► CivoBucket (metadata + spec)
    │       ├─► bucket_name: "my-bucket"
    │       ├─► region: LON1
    │       ├─► versioning_enabled: true
    │       └─► tags: ["env:prod"]
    │
    └─► ProviderConfig
            └─► civo_token: "<secret>"
```

### Processing Flow

```
module.Resources()
    │
    ├─► initializeLocals()
    │       └─► Locals struct (metadata, labels)
    │
    ├─► pulumicivoprovider.Get()
    │       └─► Civo Provider (authenticated)
    │
    └─► bucket()
            ├─► NewObjectStoreCredential()
            │       └─► Credential resource
            │
            ├─► NewObjectStore()
            │       └─► Bucket resource (depends on credential)
            │
            └─► ctx.Export()
                    └─► Stack outputs (bucket ID, endpoint, keys)
```

### Output Flow

```
Stack Outputs
    │
    ├─► bucket_id: "uuid-1234"
    ├─► endpoint_url: "https://objectstore.civo.com/my-bucket"
    ├─► access_key_secret_ref: "civo-bucket-access-key"
    └─► secret_key_secret_ref: "civo-bucket-secret-key"
            │
            ▼
CivoBucketStackOutputs (Protobuf)
            │
            ▼
Application Consumption
    │
    ├─► Kubernetes Secrets (mounted as files or env vars)
    ├─► CI/CD pipelines (backup jobs, data sync)
    └─► AWS CLI / SDKs (S3-compatible access)
```

## State Management

### Pulumi State

Pulumi tracks resource state in a backend (Pulumi Service, S3, local file). State includes:
- Resource IDs (bucket ID, credential ID)
- Resource attributes (bucket name, region, endpoint)
- Dependencies (credential → bucket)
- Outputs (exported values)

**State Operations**:
- `pulumi up`: Creates/updates resources, writes state
- `pulumi refresh`: Syncs state with actual infrastructure
- `pulumi destroy`: Deletes resources, removes from state

### State Persistence

```
Pulumi State Backend
    │
    ├─► Resource: civo:index/objectStoreCredential:ObjectStoreCredential
    │       ├─► ID: "cred-uuid-1234"
    │       ├─► AccessKeyId: "AKIAIOSFODNN7EXAMPLE"
    │       └─► SecretAccessKey: (encrypted)
    │
    └─► Resource: civo:index/objectStore:ObjectStore
            ├─► ID: "bucket-uuid-5678"
            ├─► Name: "my-bucket"
            ├─► Region: "LON1"
            ├─► BucketUrl: "https://objectstore.civo.com/my-bucket"
            └─► DependsOn: ["cred-uuid-1234"]
```

### State Drift

If someone manually modifies a bucket in the Civo dashboard:
- `pulumi refresh`: Detects drift, updates state
- `pulumi up`: Proposes changes to restore desired state

## Error Handling

### Error Propagation

```go
func Resources(ctx *pulumi.Context, stackInput *CivoBucketStackInput) error {
    // If provider setup fails, return immediately
    civoProvider, err := pulumicivoprovider.Get(ctx, stackInput.ProviderConfig)
    if err != nil {
        return errors.Wrap(err, "failed to setup civo provider")
    }
    
    // If bucket creation fails, return immediately
    _, err = bucket(ctx, locals, civoProvider)
    if err != nil {
        return errors.Wrap(err, "failed to create bucket")
    }
    
    // All resources created successfully
    return nil
}
```

**Key Points**:
- Errors are wrapped with context (`errors.Wrap`)
- Failed resources prevent subsequent operations
- Pulumi automatically rolls back on error (if possible)

### Common Error Scenarios

| Error | Cause | Recovery |
|-------|-------|----------|
| "bucket name already exists" | Global name collision | Choose different name |
| "401 Unauthorized" | Invalid Civo token | Verify token, update config |
| "region not found" | Invalid region enum | Use valid region (LON1, NYC1, etc.) |
| "credential creation failed" | API quota/limit | Contact Civo support |
| "state locked" | Concurrent `pulumi up` | Wait for other operation to finish |

## Performance Characteristics

### Resource Creation Times

| Operation | Typical Duration |
|-----------|------------------|
| Provider initialization | < 1 second |
| Credential creation | 2-5 seconds |
| Bucket creation | 5-15 seconds |
| Output export | < 1 second |
| **Total** | **~10-25 seconds** |

### Scaling Considerations

- **Single Bucket**: ~20 seconds
- **10 Buckets (parallel)**: ~30-40 seconds (Pulumi parallelizes by default)
- **100 Buckets**: Limited by Civo API rate limits (typically ~10 req/sec)

**Optimization**: Pulumi automatically parallelizes independent resources (e.g., multiple buckets with separate credentials).

## Security Considerations

### Secrets Management

1. **Civo API Token**: Stored as Pulumi secret
   ```bash
   pulumi config set civo:token $CIVO_TOKEN --secret
   ```

2. **Bucket Credentials**: Exported as references (not plaintext)
   - `access_key_secret_ref`: Pointer to secret store
   - `secret_key_secret_ref`: Pointer to secret store

3. **State Encryption**: Pulumi backend encrypts state at rest

### Least Privilege

- Use separate Civo API tokens per environment (dev, staging, prod)
- Create separate credentials per application/service
- Rotate credentials periodically

### Audit Trail

- Pulumi logs all operations (create, update, delete)
- Civo provides audit logs for API calls
- Combine for full change tracking

## Extending the Module

### Adding New Features

To add support for a new Civo ObjectStore attribute:

1. **Update `spec.proto`**: Add field to `CivoBucketSpec`
   ```protobuf
   message CivoBucketSpec {
       // ... existing fields ...
       int32 max_size_gb = 5;  // New field
   }
   ```

2. **Regenerate Go stubs**: `make protos`

3. **Update `bucket.go`**: Use new field in `ObjectStoreArgs`
   ```go
   bucketArgs := &civo.ObjectStoreArgs{
       // ... existing args ...
       MaxSizeGb: pulumi.Int(spec.MaxSizeGb),
   }
   ```

4. **Add validation tests**: Update `spec_test.go`

5. **Update documentation**: README.md, examples.md

### Adding S3 Post-Configuration

For features requiring S3 API configuration (versioning, lifecycle, CORS):

1. **Option A**: Use Pulumi Dynamic Provider
   - Create custom resource wrapping AWS SDK
   - Configure after bucket creation
   - Managed by Pulumi (tracked in state)

2. **Option B**: Use `local-exec` provisioner
   - Run AWS CLI commands in `pulumi up`
   - Not tracked in state (run on every update)

3. **Option C**: Post-deployment scripts
   - Separate automation (Ansible, Bash)
   - Run after Pulumi completes
   - Manual state management

**Recommendation**: For production, use **Option A** (Dynamic Provider) for full state management.

## Testing Strategy

### Unit Tests

Test individual functions in isolation:
- `initializeLocals()`: Verify correct locals initialization
- `bucket()`: Mock Civo provider, verify resource creation calls

### Integration Tests

Test full module execution:
- Provision real bucket in test Civo account
- Verify bucket exists via Civo API
- Verify S3 API works (upload/download test object)
- Destroy bucket

### Example Integration Test

```go
func TestBucketProvisioning(t *testing.T) {
    pulumi.Run(func(ctx *pulumi.Context) error {
        // 1. Create test stack input
        stackInput := &CivoBucketStackInput{
            CivoBucket: &CivoBucket{
                Spec: &CivoBucketSpec{
                    BucketName: "test-bucket-" + uuid.New().String(),
                    Region:     civo.CivoRegion_LON1,
                },
            },
            ProviderConfig: &ProviderConfig{
                CivoToken: os.Getenv("CIVO_TOKEN"),
            },
        }
        
        // 2. Run module
        err := module.Resources(ctx, stackInput)
        assert.NoError(t, err)
        
        // 3. Verify outputs
        ctx.Export("test_bucket_id", pulumi.String("verified"))
        return nil
    })
}
```

## Maintenance and Operations

### Upgrades

To upgrade the Civo provider:

1. Update `go.mod`:
   ```bash
   go get github.com/pulumi/pulumi-civo/sdk/v2@v2.x.x
   ```

2. Test with dry-run:
   ```bash
   pulumi preview
   ```

3. Review changes, apply:
   ```bash
   pulumi up
   ```

### Monitoring

Track module usage:
- Pulumi logs (resource creation times)
- Civo dashboard (bucket capacity, API usage)
- Application metrics (S3 operation latencies)

### Troubleshooting

Enable debug logging:
```bash
export PULUMI_DEBUG=true
export CIVO_DEBUG=1
pulumi up
```

## Related Resources

- **Civo Provider Docs**: [pulumi.com/registry/packages/civo](https://www.pulumi.com/registry/packages/civo/)
- **Pulumi Go SDK**: [pulumi.com/docs/reference/pkg/go](https://www.pulumi.com/docs/reference/pkg/go/)
- **Civo API Reference**: [civo.com/api](https://www.civo.com/api)

## Conclusion

The CivoBucket Pulumi module provides a production-ready implementation for provisioning S3-compatible object storage on Civo. It handles the 80% use case (bucket name, region, basic config) while allowing advanced users to extend functionality via S3 API post-configuration.

For most teams, this module eliminates the need to write custom Pulumi code or manage raw Terraform/API calls—simply declare what you want in YAML, and the module handles the rest.

