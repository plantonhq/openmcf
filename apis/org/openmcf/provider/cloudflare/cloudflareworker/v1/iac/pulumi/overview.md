# Pulumi Module Overview: Cloudflare Worker

## Architecture

This document explains the architectural decisions, resource dependency graph, and design rationale for the Cloudflare Worker Pulumi module.

## Resource Model

Cloudflare Workers require orchestration of multiple resources across different APIs:

```
┌─────────────────────────────────────────────────────────────┐
│                    R2 Storage (S3 API)                       │
│  ┌─────────────────────────┐                                │
│  │ Worker Bundle           │                                │
│  │ (worker-v1.0.0.js)      │                                │
│  └──────────┬──────────────┘                                │
│             │                                                 │
└─────────────┼─────────────────────────────────────────────────┘
              │ fetch during deployment
              ↓
┌─────────────────────────────────────────────────────────────┐
│                    Cloudflare Account                        │
│  ┌─────────────────────────┐                                │
│  │ Worker Script           │                                │
│  │ - Name: api-gateway     │                                │
│  │ - Content: [from R2]    │                                │
│  │ - Bindings: [KV, env]   │                                │
│  └──────────┬──────────────┘                                │
│             │                                                 │
└─────────────┼─────────────────────────────────────────────────┘
              │ referenced by
              ↓
┌─────────────────────────────────────────────────────────────┐
│                    Cloudflare Zone                           │
│  ┌─────────────────────────┐                                │
│  │ DNS Record              │                                │
│  │ api.example.com → 100:: │                                │
│  └──────────┬──────────────┘                                │
│             │                                                 │
│  ┌──────────┴──────────────┐                                │
│  │ Worker Route            │                                │
│  │ Pattern: api.example.com/*                               │
│  │ Script: api-gateway     │                                │
│  └─────────────────────────┘                                │
└─────────────────────────────────────────────────────────────┘
```

### Resource Hierarchy

**1. Worker Bundle (R2 Storage)**
- **Scope**: Stored in R2 bucket (S3-compatible)
- **Purpose**: Immutable artifact containing Worker JavaScript
- **Access**: Fetched via AWS S3 provider configured for R2

**2. Worker Script (Account-Level)**
- **Scope**: Cloudflare account
- **Purpose**: Deployable Worker with code, bindings, and configuration
- **Dependencies**: Fetches bundle from R2, references KV namespace IDs
- **API**: `POST /accounts/{account_id}/workers/scripts/{script_name}`

**3. DNS Record (Zone-Level) - Optional**
- **Scope**: Cloudflare DNS zone
- **Purpose**: Maps custom domain to Cloudflare proxy
- **Type**: AAAA record with dummy IPv6 (`100::`)
- **Requirement**: Must be proxied (orange cloud)

**4. Worker Route (Zone-Level) - Optional**
- **Scope**: Cloudflare DNS zone
- **Purpose**: Attaches Worker to URL pattern
- **Dependencies**: Requires DNS record and Worker script
- **API**: `POST /zones/{zone_id}/workers/routes`

## Dependency Graph

```
R2 Bundle → Worker Script → DNS Record → Worker Route
(external)  (creates)       (creates)    (creates)
```

### Explicit Dependencies

```go
// 1. Fetch bundle from R2
scriptObject := s3.GetObjectOutput(...)

// 2. Create Worker script using bundle content
workerScript, err := cloudflare.NewWorkersScript(...,
    Content: scriptObject.Body(),  // Dependency on R2
)

// 3. Create DNS record (no dependency on script)
dnsRecord, err := cloudflare.NewRecord(...)

// 4. Create route (depends on both)
route, err := cloudflare.NewWorkerRoute(...,
    ScriptName: workerScript.Name,  // Explicit dependency
    DependsOn: []pulumi.Resource{dnsRecord},  // Explicit dependency
)
```

Pulumi automatically orders resource creation based on these dependencies.

## OpenMCF Abstraction

### The Core Problem

When using Cloudflare API or Pulumi directly, users must:
1. Build and bundle Worker code
2. Upload bundle somewhere accessible
3. Create Worker script with inline content (size limits)
4. Manually manage KV binding IDs
5. Create DNS records and routes separately
6. Coordinate environment variables and secrets

This is complex, error-prone, and tightly couples code artifacts with infrastructure configuration.

### The Solution: R2-Based Artifacts + Declarative Config

OpenMCF separates concerns:

**Build Phase** (CI/CD):
```bash
wrangler build → dist/worker.js
aws s3 cp dist/worker.js s3://bucket/worker-v1.0.0.js
```

**Deploy Phase** (OpenMCF):
```yaml
spec:
  script:
    bundle:
      bucket: bucket
      path: worker-v1.0.0.js  # Reference, not inline content
```

**Module handles**:
1. Configure AWS S3 provider for R2
2. Fetch bundle from R2 dynamically
3. Create Worker script with fetched content
4. Wire KV bindings by ID
5. Create DNS and routes atomically
6. Upload secrets via Cloudflare API

**Benefits**:
- **Artifact immutability**: Bundles versioned in R2, enable rollbacks
- **Zero egress fees**: R2 charges nothing for bundle downloads
- **Clean separation**: Build once, deploy many times
- **Declarative management**: IaC tracks state, enables drift detection

## Implementation Details

### File Structure

```
module/
├── main.go           # Entry point - sets up providers, orchestrates resources
├── locals.go         # Local variables (shortcuts, transformations)
├── worker_script.go  # Creates Worker script with R2 bundle
├── route.go          # Creates DNS record and Worker route
├── secrets.go        # Uploads secrets via Cloudflare API
└── outputs.go        # Output constant definitions
```

### Function Flow

**1. `main.go:Resources()`**

```go
func Resources(ctx *pulumi.Context, stackInput *CloudflareWorkerStackInput) error {
    // Step 1: Initialize locals
    locals := initializeLocals(ctx, stackInput)
    
    // Step 2: Create Cloudflare provider
    cloudflareProvider, err := pulumicloudflareprovider.Get(ctx, stackInput.ProviderConfig)
    
    // Step 3: Create AWS provider for R2 (S3-compatible)
    r2Provider, err := aws.NewProvider(ctx, "r2-provider", &aws.ProviderArgs{
        Endpoints: aws.ProviderEndpointArray{
            aws.ProviderEndpointArgs{
                S3: pulumi.String(r2Endpoint),
            },
        },
        // ... S3 compatibility flags
    })
    
    // Step 4: Create Worker script (fetches from R2)
    workerScript, err := createWorkerScript(ctx, locals, cloudflareProvider, r2Provider)
    
    // Step 5: Create route (if DNS enabled)
    route, err := route(ctx, locals, cloudflareProvider, workerScript)
    
    return nil
}
```

**2. `worker_script.go:createWorkerScript()`**

```go
func createWorkerScript(...) (*cloudflare.WorkersScript, error) {
    // Fetch bundle from R2
    scriptObject := s3.GetObjectOutput(ctx, s3.GetObjectOutputArgs{
        Bucket: pulumi.String(bundle.Bucket),
        Key:    pulumi.String(bundle.Path),
    }, pulumi.Provider(r2Provider))
    
    scriptContent := scriptObject.Body()
    
    // Build bindings (env vars + KV)
    var bindings []cloudflare.WorkersScriptBindingArgs
    for k, v := range locals.CloudflareWorker.Spec.Env.Variables {
        bindings = append(bindings, cloudflare.WorkersScriptBindingArgs{
            Name: pulumi.String(k),
            Type: pulumi.String("plain_text"),
            Text: pulumi.String(v),
        })
    }
    
    // Create Worker script
    workerScript, err := cloudflare.NewWorkersScript(ctx, "workers-script", &args)
    
    return workerScript, nil
}
```

**3. `route.go:route()`**

```go
func route(...) (*cloudflare.WorkerRoute, error) {
    if !locals.CloudflareWorker.Spec.Dns.Enabled {
        return nil, nil  // DNS disabled, skip
    }
    
    // Create DNS record
    dnsRecord, err := cloudflare.NewRecord(ctx, "worker-dns", &cloudflare.RecordArgs{
        ZoneId:  pulumi.String(locals.CloudflareWorker.Spec.Dns.ZoneId),
        Name:    pulumi.String(locals.CloudflareWorker.Spec.Dns.Hostname),
        Type:    pulumi.String("AAAA"),
        Value:   pulumi.String("100::"),  // Dummy IPv6 for Workers
        Proxied: pulumi.Bool(true),
    })
    
    // Create Worker route
    workerRoute, err := cloudflare.NewWorkerRoute(ctx, "worker-route", &cloudflare.WorkerRouteArgs{
        ZoneId:     pulumi.String(locals.CloudflareWorker.Spec.Dns.ZoneId),
        Pattern:    pulumi.String(routePattern),
        ScriptName: workerScript.ScriptName,
    }, pulumi.DependsOn([]pulumi.Resource{dnsRecord}))
    
    return workerRoute, nil
}
```

## Design Decisions

### Decision 1: R2-Based Bundle Storage

**Choice**: Store Worker bundles in R2, fetch during deployment

**Rationale**:
- **Artifact immutability**: Versioned bundles enable rollbacks
- **Build once, deploy many**: Reuse same bundle across environments
- **Zero egress**: Free to download bundles during deployment
- **S3 compatibility**: Standard tools work (AWS CLI, SDKs)

**Trade-off**:
- ✅ Clean separation of build and deploy
- ❌ Requires R2 credentials

**Alternative Considered**: Inline script content in proto. Rejected due to size limits and versioning challenges.

### Decision 2: Separate AWS Provider for R2

**Choice**: Create dedicated AWS provider configured for R2

**Rationale**:
- R2 uses S3-compatible API
- Pulumi AWS provider handles authentication and S3 operations
- Cleaner than custom HTTP client

**Trade-off**:
- ✅ Reuse battle-tested S3 client
- ❌ Requires two providers (Cloudflare + AWS)

### Decision 3: Dummy IPv6 for DNS Record

**Choice**: Use `100::` as dummy IPv6 address for Worker DNS records

**Rationale**:
- Workers don't have real IPs (they're edge functions)
- Cloudflare requires AAAA record for proxied Workers
- `100::` is a reserved prefix that won't conflict

**Trade-off**:
- ✅ Works reliably
- ❌ Slightly confusing (why dummy IP?)

**Cloudflare Recommendation**: This is the official pattern.

### Decision 4: Unified Bindings Array

**Choice**: Combine env vars and KV bindings into single array

**Rationale**:
- Cloudflare Workers API v6 uses unified bindings model
- Pulumi provider expects single bindings array
- Simplifies code

**Trade-off**:
- ✅ Matches API design
- ❌ Slightly more complex than separate fields

## Secret Management

### Secrets vs Variables

**Variables** (plain text):
- Stored in Worker metadata
- Visible in Cloudflare Dashboard
- Good for non-sensitive config

**Secrets** (encrypted):
- Encrypted at rest by Cloudflare
- Never visible in Dashboard or logs
- Uploaded via separate Secrets API

### Implementation

```go
// Plain text variables → bindings
for k, v := range locals.CloudflareWorker.Spec.Env.Variables {
    bindings = append(bindings, WorkersScriptBindingArgs{
        Name: pulumi.String(k),
        Type: pulumi.String("plain_text"),
        Text: pulumi.String(v),
    })
}

// Secrets uploaded separately (future enhancement)
// Currently logged as warning - manual upload required
```

**Note**: Pulumi Cloudflare provider doesn't yet support Worker secrets upload. This requires direct API call or manual configuration.

## Performance Considerations

### Deployment Time

- Bundle fetch from R2: ~1-2 seconds
- Worker script creation: ~2-3 seconds
- DNS record creation: ~1-2 seconds
- Route attachment: ~1-2 seconds
- **Total**: 5-10 seconds

### State File Size

- Worker script: ~2-3 KB (includes bindings)
- DNS record: ~500 bytes
- Route: ~500 bytes
- **Total**: ~3-4 KB per Worker

## Error Handling

The module uses Go's error wrapping:

```go
if err != nil {
    return errors.Wrap(err, "failed to create worker script")
}
```

Pulumi displays full error chain:

```
error: failed to create cloudflare worker route:
  failed to create cloudflare workers script:
    failed to fetch bundle from R2:
      NoSuchKey: The specified key does not exist
```

## Future Enhancements

### 1. D1 Database Bindings

```protobuf
repeated ValueFromRef d1_bindings = 8;
```

### 2. R2 Bucket Bindings

```protobuf
repeated ValueFromRef r2_bindings = 9;
```

### 3. Durable Object Bindings

```protobuf
repeated DurableObjectBinding do_bindings = 10;
```

### 4. Cron Triggers

```protobuf
repeated string cron_triggers = 11;  // e.g., "0 0 * * *"
```

### 5. Service Bindings

```protobuf
repeated ServiceBinding service_bindings = 12;
```

## Comparison: Pulumi vs. Terraform

| Aspect | Pulumi (This Module) | Terraform |
|--------|----------------------|-----------|
| **R2 Bundle Fetch** | Native via AWS provider | Native via AWS provider |
| **Type Safety** | Compile-time Go checks | Runtime HCL validation |
| **Secret Management** | Built-in encryption | Requires external vault |
| **Conditionals** | Go if/else | HCL dynamic blocks |
| **Multi-Provider** | Natural (Go imports) | Requires careful configuration |

**Recommendation**: Both are production-ready. Choose based on team preference.

## References

- [Cloudflare Workers API](https://developers.cloudflare.com/api/operations/worker-script-upload-worker-module)
- [Pulumi Cloudflare Provider](https://www.pulumi.com/registry/packages/cloudflare/)
- [Component README](../../README.md)
- [Research Documentation](../../docs/README.md)

---

**Questions?** Review the [README.md](./README.md) or consult Cloudflare's API documentation.

