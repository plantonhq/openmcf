# Pulumi Module Overview - DigitalOcean App Platform Service

## Architecture

This document explains the internal architecture and design decisions of the Pulumi module for deploying DigitalOcean App Platform services.

## Module Structure

```
iac/pulumi/
├── main.go              # Pulumi program entrypoint
├── Pulumi.yaml          # Pulumi project configuration
├── Makefile             # Build and deployment automation
├── debug.sh             # Debug helper script
└── module/
    ├── main.go          # Module entry point (Resources function)
    ├── locals.go        # Local variable initialization
    ├── outputs.go       # Output constants
    └── app_platform_service.go  # Core resource creation logic
```

## Design Principles

### 1. Single Responsibility

Each file has a clear, focused purpose:
- `main.go` (module): Entry point that orchestrates resource creation
- `app_platform_service.go`: App Platform resource creation and configuration
- `locals.go`: Data transformation and metadata preparation
- `outputs.go`: Output key constants

### 2. Type Safety

The module leverages Go's type system and protobuf-generated types to:
- Catch configuration errors at compile time
- Provide IDE autocomplete and type hints
- Enforce required fields through protobuf validation

### 3. Abstraction Levels

The module operates at three levels:
- **API Level**: Protobuf spec (80/20 simplified configuration)
- **Module Level**: Go functions that transform spec to Pulumi resources
- **Provider Level**: DigitalOcean provider resources (complete App Platform API)

## Resource Creation Flow

### Step 1: Initialize Locals

```go
func initializeLocals(ctx *pulumi.Context, stackInput *DigitalOceanAppPlatformServiceStackInput) *Locals
```

Prepares:
- Metadata (name, labels)
- Target resource reference
- Configuration lookups

### Step 2: Create Provider

```go
digitalOceanProvider, err := pulumidigitaloceanprovider.Get(ctx, stackInput.ProviderConfig)
```

Initializes DigitalOcean provider with:
- API token from provider config
- Region configuration
- Retry settings

### Step 3: Build Service Configuration

The module branches based on `service_type`:

#### Web Service
```go
func buildWebService(spec, instanceCount, instanceSizeSlug) *digitalocean.AppSpecServiceArgs
```
- Creates HTTP-accessible service
- Configures autoscaling (if enabled)
- Sets up load balancing
- Exports live URL

#### Worker
```go
func buildWorker(spec, instanceCount, instanceSizeSlug) *digitalocean.AppSpecWorkerArgs
```
- Creates background processing service
- No HTTP ingress
- Fixed instance count (no autoscaling)

#### Job
```go
func buildJob(spec, instanceSizeSlug) *digitalocean.AppSpecJobArgs
```
- Creates pre-deployment job
- Runs before each deployment
- Single instance execution

### Step 4: Configure Source

#### Git Source
```go
func configureSource(spec, serviceArgs)
```
For git-based deployments:
- Sets `repo_clone_url` from spec
- Configures branch
- Optionally overrides build/run commands

#### Image Source
For container image deployments:
- Sets registry URL (DOCR by default)
- Configures repository and tag
- Handles registry authentication

### Step 5: Apply Environment Variables

```go
func buildEnvVars(env map[string]string) digitalocean.AppSpecServiceEnvArray
```

Transforms env var map to Pulumi array:
- Sets scope (`RUN_AND_BUILD_TIME` for services, `RUN_TIME` for workers/jobs)
- Marks type as `GENERAL` (DigitalOcean encrypts on server-side)

### Step 6: Configure Autoscaling

For web services with `enable_autoscale: true`:
```go
Autoscaling: &digitalocean.AppSpecServiceAutoscalingArgs{
    MinInstanceCount: pulumi.Int(spec.MinInstanceCount),
    MaxInstanceCount: pulumi.Int(spec.MaxInstanceCount),
    Metrics: &digitalocean.AppSpecServiceAutoscalingMetricsArgs{
        Cpu: &digitalocean.AppSpecServiceAutoscalingMetricsCpuArgs{
            Percent: pulumi.Int(80), // Default CPU threshold
        },
    },
}
```

### Step 7: Add Custom Domain

If `custom_domain` is specified:
```go
appSpecArgs.Domains = digitalocean.AppSpecDomainArray{
    &digitalocean.AppSpecDomainArgs{
        Name: pulumi.String(spec.CustomDomain.GetValue()),
        Type: pulumi.String("PRIMARY"),
    },
}
```

### Step 8: Create App Platform Resource

```go
createdApp, err := digitalocean.NewApp(ctx, "app", &digitalocean.AppArgs{Spec: appSpecArgs}, pulumi.Provider(digitalOceanProvider))
```

### Step 9: Export Outputs

```go
ctx.Export(OpAppId, createdApp.ID())
ctx.Export(OpLiveUrl, createdApp.LiveUrl)
```

## Key Design Decisions

### Why Separate Functions for Service Types?

Each service type (web/worker/job) has different Pulumi types:
- `digitalocean.AppSpecServiceArgs`
- `digitalocean.AppSpecWorkerArgs`
- `digitalocean.AppSpecJobArgs`

Separate functions provide:
- Type safety (Go compiler enforces correct fields)
- Clear separation of concerns
- Easier testing and maintenance

### Why Default CPU Threshold to 80%?

Autoscaling at 80% CPU provides:
- Headroom for traffic spikes
- Cost efficiency (not over-provisioning)
- Industry-standard threshold

Users cannot override this currently (80/20 principle - most users want the default).

### Why DOCR as Default Registry Type?

DigitalOcean Container Registry is:
- Tightly integrated with App Platform
- No authentication configuration needed
- Automatic credential rotation
- Cost-effective for DO customers

For non-DOCR registries, users can still provide custom registry URLs.

### Instance Size Slug Transformation

The module converts underscores to hyphens:
```go
instanceSizeSlug := strings.ReplaceAll(spec.InstanceSizeSlug.String(), "_", "-")
```

**Reason**: Protobuf enum names use underscores (`basic_xxs`), but DigitalOcean API expects hyphens (`basic-xxs`).

## Error Handling

The module uses Go's error wrapping for context:

```go
if err != nil {
    return errors.Wrap(err, "failed to create digitalocean app")
}
```

This provides:
- Stack traces for debugging
- Contextual error messages
- Pulumi error reporting integration

## State Management

Pulumi automatically tracks:
- App ID
- Live URL
- Deployment status
- Resource dependencies

On subsequent `pulumi up`:
- Pulumi compares desired state (spec) with actual state (DigitalOcean API)
- Generates minimal diff
- Updates only changed resources

## Foreign Key Resolution

The module uses OpenMCF's foreign key pattern for `custom_domain` and `registry`:

```go
spec.CustomDomain.GetValue()
```

This resolves:
- Direct values (`"api.myapp.com"`)
- References to other resources (via foreign key lookup)

## Instance Count Defaults

```go
instanceCount := spec.InstanceCount
if instanceCount == 0 {
    instanceCount = 1 // Default to 1
}
```

**Reason**: Protobuf `uint32` defaults to `0`, but DigitalOcean requires at least 1 instance.

## Future Enhancements

Potential additions (not currently in 80/20 scope):
- Custom health check intervals
- Advanced routing rules
- Blue-green deployment strategies
- Custom buildpack configurations
- Alert policies
- Database connection pooling parameters

## Testing Strategy

While this module doesn't include unit tests, recommended testing:
1. **Integration tests**: Deploy to a test DO account
2. **Pulumi preview**: Review changes before applying
3. **Staging environments**: Test on non-production stacks
4. **Rollback capability**: Use Pulumi's stack history

## Performance Considerations

The module is optimized for:
- **Fast preview**: Minimal API calls during `pulumi preview`
- **Efficient updates**: Only modified resources are updated
- **Parallel resource creation**: Pulumi handles dependencies

## References

- [Pulumi Programming Model](https://www.pulumi.com/docs/intro/concepts/programming-model/)
- [DigitalOcean Provider SDK](https://github.com/pulumi/pulumi-digitalocean)
- [OpenMCF Architecture](../../../../../../architecture/deployment-component.md)

