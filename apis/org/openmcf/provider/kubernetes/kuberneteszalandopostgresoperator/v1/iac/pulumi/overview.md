# Zalando Postgres Operator Pulumi Module - Architecture Overview

This document provides detailed architecture information, design decisions, and implementation patterns for the Zalando Postgres Operator Pulumi module.

## Module Purpose

Deploy and configure the Zalando Postgres Operator on Kubernetes clusters with:
- Automated namespace and Helm chart deployment
- Optional Cloudflare R2 backup configuration
- Label inheritance for consistent resource tagging
- Production-ready defaults with customizable resource limits

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│              Zalando Postgres Operator Pulumi Module                     │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                           │
│  ┌──────────────────┐                                                    │
│  │   Stack Input    │                                                    │
│  │  (Protobuf API)  │                                                    │
│  └────────┬─────────┘                                                    │
│           │                                                               │
│           ▼                                                               │
│  ┌──────────────────┐                                                    │
│  │  initializeLocals│──────► Labels (resource, org, env, kind, id)      │
│  └────────┬─────────┘                                                    │
│           │                                                               │
│           ▼                                                               │
│  ┌──────────────────────────────────────────────────────────┐            │
│  │            postgresOperator Function                      │            │
│  │                                                           │            │
│  │  1. Create or Use Namespace (conditional)                │            │
│  │     ├─► If create_namespace=true: Create namespace      │            │
│  │     └─► If create_namespace=false: Use existing         │            │
│  │                                                           │            │
│  │  2. createBackupResources (if backup_config provided)   │            │
│  │     ├─► Secret: r2-postgres-backup-credentials          │            │
│  │     │    └─► AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY   │            │
│  │     └─► ConfigMap: postgres-pod-backup-config           │            │
│  │          ├─► USE_WALG_BACKUP / RESTORE / CLONE          │            │
│  │          ├─► WALG_S3_PREFIX, AWS_ENDPOINT               │            │
│  │          ├─► AWS_REGION, AWS_FORCE_PATH_STYLE           │            │
│  │          └─► BACKUP_SCHEDULE                            │            │
│  │                                                           │            │
│  │  3. Deploy Helm Chart                                    │            │
│  │     ├─► Chart: postgres-operator/postgres-operator      │            │
│  │     ├─► Version: 1.12.2                                 │            │
│  │     ├─► Values:                                          │            │
│  │     │    ├─► inherited_labels (5 labels)                │            │
│  │     │    └─► pod_environment_configmap (if backup)      │            │
│  │     └─► Resources: CPU/Memory limits                    │            │
│  │                                                           │            │
│  │  4. Export Outputs                                       │            │
│  │     └─► namespace: postgres-operator                    │            │
│  └──────────────────────────────────────────────────────────┘            │
│                                                                           │
└─────────────────────────────────────────────────────────────────────────┘
                                   │
                                   ▼
                     ┌──────────────────────────────┐
                     │   Kubernetes Cluster         │
                     │                              │
                     │  Namespace: postgres-operator│
                     │  ├─► Deployment (operator)   │
                     │  ├─► Service (operator)      │
                     │  ├─► Secret (if backup)      │
                     │  ├─► ConfigMap (if backup)   │
                     │  └─► CRDs (postgresqls)      │
                     └──────────────────────────────┘
```

## File Organization

### Entry Points

| File | Purpose | Key Function |
|------|---------|--------------|
| `main.go` | CLI integration | Converts YAML manifest to Pulumi input |
| `module/main.go` | Module entry point | `Resources()` - orchestrates deployment |

### Core Module Files

| File | Responsibilities |
|------|------------------|
| `module/postgres_operator.go` | Namespace creation, backup resources, Helm deployment |
| `module/backup_config.go` | Secret and ConfigMap creation for R2 backups |
| `module/locals.go` | Label computation and input transformation |
| `module/outputs.go` | Output constant definitions |
| `module/vars.go` | Helm chart configuration constants |

## Data Flow

### Input → Locals → Resources

```
KubernetesZalandoPostgresOperatorStackInput (Protobuf)
    │
    ├─► Target (KubernetesZalandoPostgresOperator)
    │    ├─► Metadata (name, id, org, env, labels)
    │    └─► Spec (container, backup_config)
    │
    └─► ProviderConfig (Kubernetes cluster credentials)

                    ▼

              initializeLocals()
    ┌─────────────────────────────────┐
    │ Locals struct:                  │
    │  - KubernetesZalandoPostgresOp  │
    │  - KubernetesLabels (5 labels)  │
    └─────────────────────────────────┘

                    ▼

          postgresOperator() function
    ┌─────────────────────────────────┐
    │ 1. Namespace                    │
    │ 2. Backup Resources (optional)  │
    │ 3. Helm Chart                   │
    │ 4. Outputs                      │
    └─────────────────────────────────┘
```

### Output Exports

| Export Key | Type | Source |
|------------|------|--------|
| `namespace` | string | Fixed value: `"postgres-operator"` |

## Design Decisions

### 1. Conditional Namespace Creation

**Decision**: Allow users to choose between creating namespace or using existing one via `create_namespace` flag

**Rationale**:
- Provides flexibility for different deployment scenarios
- Some environments have pre-existing namespaces managed by platform teams
- Allows integration with existing namespace management workflows
- Zalando operator is cluster-scoped (watches all namespaces)
- Single operator instance per cluster is the recommended pattern

**Implementation**:
```go
if locals.KubernetesZalandoPostgresOperator.Spec.CreateNamespace {
    // Create new namespace
    createdNamespace, err := corev1.NewNamespace(ctx, ...)
    namespaceOutput = createdNamespace.Metadata.Name()
} else {
    // Use existing namespace - just reference the name
    namespaceOutput = pulumi.String(namespace).ToStringOutput()
}
```

**Trade-offs**:
- More complexity in handling namespace references
- **Benefit**: Flexibility for different operational models

### 2. Conditional Backup Resources

**Decision**: Only create Secret and ConfigMap when `backup_config` is provided

**Rationale**:
- Backup is optional (development environments may not need it)
- Avoid creating unused resources
- Keep Secret management clean (no empty credentials)

**Implementation**:
```go
if backupConfig == nil {
    return pulumi.String("").ToStringOutput(), nil
}
```

Returns empty string when no backup, which is handled by Helm values logic.

### 3. Label Inheritance Configuration

**Decision**: Configure operator to inherit 5 specific labels

**Rationale**:
- Ensures all PostgreSQL databases have consistent metadata
- Enables multi-tenancy (org, env labels)
- Facilitates cost tracking and resource management
- Follows OpenMCF conventions

**Inherited Labels**:
1. `resource` - Marks as OpenMCF-managed
2. `organization` - Multi-tenancy support
3. `environment` - Environment segregation
4. `resource_kind` - Resource type tracking
5. `resource_id` - Unique identifier

### 4. Helm-Based Deployment

**Decision**: Use Zalando's official Helm chart instead of raw manifests

**Rationale**:
- Official support and updates from Zalando team
- Comprehensive configuration options via Helm values
- Easier upgrades (change version number)
- Community-tested and production-proven

**Trade-offs**:
- Dependency on Helm chart stability
- Limited customization beyond Helm values
- **Benefit**: Less maintenance burden, automatic updates

### 5. R2-Specific Backup Configuration

**Decision**: Only support Cloudflare R2 (not generic S3)

**Rationale**:
- R2 is S3-compatible with specific quirks (endpoint format, region "auto")
- Simplifies credential management (account ID + bucket + keys)
- Clear scope reduces complexity
- Future: Can add AWS S3, GCS support as separate config blocks

**R2 Specifics**:
```go
r2Endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId)
configMapData["AWS_REGION"] = "auto"
configMapData["AWS_FORCE_PATH_STYLE"] = "true"
```

### 6. WAL-G Environment Variables in ConfigMap

**Decision**: Store WAL-G configuration in ConfigMap (including credentials)

**Rationale**:
- Zalando operator expects `pod_environment_configmap` reference
- Operator mounts ConfigMap as environment variables in database pods
- Credentials are also in ConfigMap (duplicated from Secret for operator convenience)
- Future: Could use Secret-only approach with operator enhancement

**Security Note**: Credentials appear in both Secret and ConfigMap. The Secret is the source of truth; ConfigMap is for operator consumption.

### 7. Default WAL-G Settings

**Decision**: Default all WAL-G flags to `true` (backup, restore, clone)

**Rationale**:
- Production-ready defaults (encourage backup usage)
- Users must explicitly disable (opt-out pattern)
- Safer than opt-in (prevents accidental data loss)

**Implementation**:
```go
func boolToString(value bool, defaultWhenFalse bool) string {
    if value { return "true" }
    if defaultWhenFalse { return "true" }
    return "false"
}
```

## Resource Dependencies

### Explicit Dependencies

```
namespace
  │
  ├─► backup_secret (if backup configured)
  │    │
  │    └─► backup_configmap
  │         │
  │         └─► helm_release
  │
  └─► helm_release (if no backup)
```

### Pulumi Dependency Management

- `pulumi.DependsOn([]pulumi.Resource{createdSecret})` for ConfigMap
- `pulumi.Parent(createdNamespace)` for Helm release
- Automatic dependency inference via resource references

## Helm Chart Configuration

### Base Values

Always set:

```yaml
configKubernetes:
  inherited_labels:
    - resource
    - organization
    - environment
    - resource_kind
    - resource_id
```

### Conditional Values

Set only when `backup_config` is provided:

```yaml
configKubernetes:
  pod_environment_configmap: "postgres-operator/postgres-pod-backup-config"
```

### Dynamic Values Construction

```go
baseValues := pulumi.Map{
    "configKubernetes": pulumi.Map{
        "inherited_labels": pulumi.ToStringArray([...]string),
    },
}

if cmName != "" {
    baseValues["configKubernetes"].(pulumi.Map)["pod_environment_configmap"] = pulumi.String(cmName)
}
```

## Backup Architecture

### R2 Endpoint Construction

```go
// Input: cloudflare_account_id = "abc123"
// Output: https://abc123.r2.cloudflarestorage.com
r2Endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", r2Config.CloudflareAccountId)
```

### S3 Prefix Template

Default template uses Zalando operator variables:

```
backups/$(SCOPE)/$(PGVERSION)
```

- `$(SCOPE)` = PostgreSQL cluster name
- `$(PGVERSION)` = PostgreSQL version (e.g., "15")

Example expansion:
```
s3://postgres-backups-prod/backups/my-app-db/15/
```

Custom template example:
```
prod/pg/$(SCOPE)/v$(PGVERSION)/$(CLUSTER)
```

### WAL-G S3 Prefix

Full path with bucket:

```go
walgS3Prefix := fmt.Sprintf("s3://%s/%s", bucketName, s3PrefixTemplate)
// Result: s3://postgres-backups-prod/backups/$(SCOPE)/$(PGVERSION)
```

### ConfigMap Data Structure

```yaml
data:
  # WAL-G Feature Flags
  USE_WALG_BACKUP: "true"
  USE_WALG_RESTORE: "true"
  CLONE_USE_WALG_RESTORE: "true"
  
  # S3/R2 Configuration
  WALG_S3_PREFIX: "s3://bucket-name/backups/$(SCOPE)/$(PGVERSION)"
  AWS_ENDPOINT: "https://account-id.r2.cloudflarestorage.com"
  AWS_REGION: "auto"
  AWS_FORCE_PATH_STYLE: "true"
  
  # Backup Schedule
  BACKUP_SCHEDULE: "0 2 * * *"
  
  # Credentials (from Secret)
  AWS_ACCESS_KEY_ID: "..."
  AWS_SECRET_ACCESS_KEY: "..."
```

## Error Handling

### Pre-Flight Validation

None currently implemented. Validation is handled by Protobuf buf.validate rules at API level.

### Resource Creation Errors

All resource creation functions return errors wrapped with context:

```go
if err != nil {
    return errors.Wrap(err, "failed to create backup credentials secret")
}
```

### Backup Config Validation

```go
if backupConfig.R2Config == nil {
    return pulumi.String("").ToStringOutput(), 
           errors.New("backup_config.r2_config is required when backup_config is specified")
}
```

## Testing Strategy

### Unit Tests

Currently not implemented. Future: Add tests for:
- `boolToString()` helper function
- R2 endpoint construction
- S3 prefix template expansion

### Manual Testing

Use `debug.sh` script with `hack/manifest.yaml`:

```bash
cd iac/pulumi
./debug.sh
```

Verifies:
1. Namespace creation
2. Secret creation (if backup configured)
3. ConfigMap creation (if backup configured)
4. Helm release deployment
5. Operator pod startup

## Extension Points

### Adding New Backup Providers

To add AWS S3 or Google Cloud Storage:

1. Add new message to `spec.proto`:
   ```protobuf
   message KubernetesZalandoPostgresOperatorBackupS3Config { ... }
   ```

2. Update `backup_config.go`:
   ```go
   func createBackupResources(...) {
       if backupConfig.S3Config != nil {
           // Handle AWS S3
       } else if backupConfig.R2Config != nil {
           // Handle Cloudflare R2
       }
   }
   ```

3. Add endpoint construction logic for each provider

### Customizing Helm Values

To expose more Helm chart options:

1. Add fields to `spec.proto`
2. Pass to `postgresOperator()` function
3. Add to `helmValues` map:
   ```go
   helmValues["newSetting"] = pulumi.String(spec.NewSetting)
   ```

### Additional Outputs

To add new exports:

1. Define constant in `outputs.go`:
   ```go
   const OpService = "service"
   ```

2. Export in `postgresOperator()`:
   ```go
   ctx.Export(OpService, pulumi.String(serviceName))
   ```

## Known Limitations

1. **Single Operator Per Cluster**: Module assumes one operator instance per Kubernetes cluster
2. **R2 Only**: Backup currently only supports Cloudflare R2 (not AWS S3, GCS, Azure Blob)
3. **No Operator HA**: Module doesn't configure operator high availability (multiple replicas)
4. **Limited Helm Customization**: Only exposes a subset of Helm chart values
5. **Namespace Existence**: When `create_namespace: false`, the namespace must already exist (no validation at plan time)

## Performance Considerations

### Deployment Time

- Namespace creation: <1 second
- Secret/ConfigMap creation: <1 second each
- Helm chart deployment: 30-60 seconds (includes CRD installation)
- Total deployment: ~1 minute

### Resource Usage

**Operator Pod** (default):
- Requests: 50m CPU, 100Mi memory
- Limits: 1000m CPU, 1Gi memory

**Recommended Production**:
- Requests: 100m CPU, 256Mi memory
- Limits: 2000m CPU, 2Gi memory

## Security Considerations

1. **Secret Management**: R2 credentials stored in Kubernetes Secret (base64-encoded, not encrypted at rest by default)
2. **ConfigMap Credentials**: Credentials duplicated in ConfigMap for operator consumption
3. **RBAC**: Operator requires cluster-admin permissions (installed by Helm chart)
4. **Network Policies**: Not configured by module (operator needs cluster-wide access)

**Recommendations**:
- Use external secret management (Vault, AWS Secrets Manager) for production
- Enable Kubernetes encryption at rest
- Review and restrict operator RBAC if possible

## Comparison to Terraform Module

| Aspect | Pulumi (Go) | Terraform (HCL) |
|--------|-------------|-----------------|
| Language | Go | HCL |
| Conditionals | Native `if` | `count`, `dynamic` blocks |
| String Manipulation | Native Go | `format()`, `join()` functions |
| Error Handling | Explicit errors | Plan-time validation |
| Type Safety | Compile-time | Plan-time |
| Dependency Management | Explicit `DependsOn` | Implicit + `depends_on` |
| Resource Organization | Multi-file Go package | Multi-file HCL |

**Shared Concepts**:
- Helm provider for chart deployment
- Kubernetes provider for raw resources
- Conditional resource creation
- Label management and propagation

## Maintenance Guidelines

1. **Keep Files Focused**: Each file should manage one logical concern
2. **Follow "Thin Locals" Pattern**: Locals should be simple transformations, not complex logic
3. **Wrap Errors**: Always wrap errors with context for debugging
4. **Test Backup Configurations**: All backup settings must be tested with real R2 buckets
5. **Version Pin Carefully**: Only update Helm chart version after testing

## References

- [Spec Definition](../../spec.proto)
- [Stack Outputs](../../stack_outputs.proto)
- [Zalando Operator Configuration](https://postgres-operator.readthedocs.io/en/latest/reference/operator_parameters/)
- [Zalando Helm Chart](https://github.com/zalando/postgres-operator/tree/master/charts/postgres-operator)
- [WAL-G Environment Variables](https://github.com/wal-g/wal-g/blob/master/docs/STORAGES.md)
- [Cloudflare R2 Documentation](https://developers.cloudflare.com/r2/)

