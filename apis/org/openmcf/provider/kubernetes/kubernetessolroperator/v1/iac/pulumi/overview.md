# Pulumi Module Architecture Overview

## Purpose

This document explains the design decisions, architecture, and implementation details of the KubernetesSolrOperator Pulumi module.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│  Stack Input                                                │
│  - Provider Config (Kubernetes credential)                  │
│  - Metadata (name, org, env)                                │
│  - Spec (container resources)                               │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ↓
┌─────────────────────────────────────────────────────────────┐
│  Pulumi Module (module/main.go)                             │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  1. Setup Kubernetes Provider                        │   │
│  │     - Authenticate using credential                  │   │
│  └──────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  2. Create Namespace                                 │   │
│  │     - namespace: "solr-operator"                     │   │
│  └──────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  3. Apply CRDs                                       │   │
│  │     - SolrCloud CRD                                  │   │
│  │     - SolrBackup CRD                                 │   │
│  │     - SolrPrometheusExporter CRD                     │   │
│  └──────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  4. Deploy Helm Chart                                │   │
│  │     - Chart: solr-operator                           │   │
│  │     - Version: 0.7.0 (configurable)                  │   │
│  │     - Wait for deployment to complete                │   │
│  └──────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  5. Export Outputs                                   │   │
│  │     - namespace                                      │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                     │
                     ↓
┌─────────────────────────────────────────────────────────────┐
│  Kubernetes Cluster                                         │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  Namespace: solr-operator                             │  │
│  │  ┌─────────────────────────────────────────────────┐  │  │
│  │  │  Deployment: solr-operator                      │  │  │
│  │  │  - Watches SolrCloud CRDs                       │  │  │
│  │  │  - Manages SolrCloud deployments                │  │  │
│  │  │  - Handles backups and exports                  │  │  │
│  │  └─────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## Key Design Decisions

### 1. CRDs Before Helm Chart

**Decision**: Apply CRDs explicitly before deploying the Helm chart.

**Rationale**:
- Helm's CRD handling has limitations (doesn't update existing CRDs)
- Explicit CRD application ensures latest versions are installed
- Allows verification that CRDs are present before operator starts
- Enables operator to function immediately on deployment

**Implementation**:
```go
crds, err := pulumiyaml.NewConfigFile(ctx, "solr-operator-crds",
    &pulumiyaml.ConfigFileArgs{
        File: vars.CrdManifestDownloadURL,
    },
    pulumi.Provider(kubeProvider),
    pulumi.Parent(ns))

// Helm release depends on CRDs
pulumi.DependsOn([]pulumi.Resource{crds})
```

### 2. Namespace Management

**Decision**: Support both namespace creation and usage of existing namespaces via `create_namespace` flag.

**Rationale**:
- **Flexibility**: Users can choose their namespace management strategy
- **GitOps Compatibility**: Allows separate namespace lifecycle management
- **Multi-tenant Support**: Platform teams can pre-provision namespaces with policies
- **Backward Compatibility**: Default behavior (false) prevents accidental namespace creation

**Implementation**:
- When `create_namespace: true`: Create namespace resource, set as parent for CRDs and Helm
- When `create_namespace: false`: Reference namespace by name, no parent relationship

**Tradeoff**: Requires users to understand namespace lifecycle implications.

### 3. Atomic Helm Deployment

**Decision**: Enable `Atomic: true` and `CleanupOnFail: true` on Helm release.

**Rationale**:
- Atomic deployment ensures all-or-nothing semantics
- Failed deployments automatically roll back
- Prevents partial operator states that could break SolrCloud management
- Production-safe deployment strategy

**Tradeoff**: Slower deployments (waits for all pods to be ready), but safer.

### 4. Hardcoded Stable Version

**Decision**: Default to a known-stable chart version (`0.7.0`), not `latest`.

**Rationale**:
- Reproducible deployments across environments
- Avoids unexpected breaking changes
- Allows controlled upgrades by explicitly changing version
- Production best practice for infrastructure-as-code

**Flexibility**: Version is defined in `vars.go` and can be easily updated for upgrades.

### 5. Minimal Values Customization

**Decision**: Deploy with default Helm values (empty `Values` map).

**Rationale**:
- Apache Solr Operator Helm chart has sensible defaults
- Reduces complexity and configuration drift
- Operator's own configuration should be in SolrCloud CRDs, not Helm values
- Users can fork/extend module if deep customization needed

**Future Enhancement**: Could expose common customizations (resource limits, affinity) through StackInput spec.

## Module Structure Explained

### main.go

Entry point for the Pulumi program. Minimal logic—just instantiates the module with stack input.

### module/main.go

Core resource creation logic. Orchestrates:
1. Provider setup
2. Namespace creation
3. CRD application
4. Helm chart deployment

**Error Handling**: Wraps all errors with context for clear debugging.

**Dependencies**: Explicit dependency chain (namespace → CRDs → Helm) ensures correct ordering.

### module/locals.go

Computes derived values from stack input. Currently:
- Extracts metadata for labels
- Builds common label set for resources
- Prepares operator name from metadata

**Pattern**: Follows OpenMCF convention of separating input transformation from resource creation.

### module/vars.go

Constants and configuration shared across module. Includes:
- Namespace name
- Helm chart details
- Version defaults
- CRD manifest URLs

**Pattern**: Single source of truth for deployment parameters. Easy to update for new versions.

### module/outputs.go

Output constant definitions. Keeps output keys consistent and typed.

**Current Outputs**:
- `namespace`: Where operator is deployed

**Future Outputs**: Could add operator version, CRD status, etc.

## Resource Ordering

The deployment follows strict ordering:

```
1. Kubernetes Provider
   ↓
2. Namespace
   ↓
3. CRDs (parallel with namespace as parent)
   ↓
4. Helm Release (depends on CRDs, namespace as parent)
```

This ordering ensures:
- CRDs exist before operator starts
- Operator has valid namespace to deploy into
- All resources are properly parented for cascade deletion

## Idempotency

The module is fully idempotent:

- **CRDs**: ConfigFile resource detects existing CRDs and updates if needed
- **Helm**: Helm provider handles in-place updates of existing releases
- **Namespace**: Kubernetes provider won't recreate existing namespace

Running `pulumi up` multiple times safely converges to desired state.

## Upgrade Strategy

### Operator Upgrades

1. Update `DefaultStableVersion` in `vars.go`
2. Run `pulumi preview` to see changes
3. Run `pulumi up` to apply

Helm handles rolling update of operator pod.

### CRD Upgrades

1. Update `CrdManifestDownloadURL` in `vars.go`
2. Run `pulumi preview`
3. **Important**: Check for API-breaking changes in release notes
4. Run `pulumi up`

**Warning**: CRD updates may require SolrCloud resource recreation if API version changes.

## State Management

Pulumi state tracks:
- Namespace UID
- Helm release revision
- CRD ConfigFile checksum

This allows Pulumi to detect drift and plan corrective actions.

## Error Scenarios

### CRD Download Fails

If `vars.CrdManifestDownloadURL` is unreachable:
- Pulumi fails immediately (before Helm deployment)
- No partial state created
- User can fix URL and retry

### Helm Chart Not Found

If chart repository is down or chart doesn't exist:
- Namespace and CRDs are created
- Helm deployment fails
- User can manually install chart or wait for repository recovery
- Re-running `pulumi up` resumes from checkpoint

### Operator Pod CrashLoops

If operator pod fails to start:
- Helm considers deployment failed (due to `WaitForJobs: true`)
- Rollback occurs (due to `Atomic: true`)
- Pulumi reports failure
- Previous operator version remains running (if upgrade) or nothing deployed (if new install)

## Comparison to Terraform

| Aspect | Pulumi (This Module) | Terraform |
|--------|----------------------|-----------|
| CRD Handling | Explicit YAML ConfigFile | Must use kubectl provider or null_resource |
| Type Safety | Strong (Go types) | Weak (HCL maps) |
| Helm Integration | Native, first-class | Via provider, less integrated |
| Dependency Management | Explicit, compile-time checked | Implicit, runtime only |
| State | JSON-based, flexible backends | JSON/HCL, limited backends |

## Future Enhancements

1. **Custom Values Support**: Expose Helm values through StackInput spec
2. **Multi-Version Support**: Allow users to specify operator version in spec
3. **Resource Customization**: Support custom resource limits for operator pod
4. **Health Checks**: Add readiness probes and status checks
5. **Metrics**: Export operator deployment metrics as outputs

## Testing Approach

Manual testing checklist:

- [ ] Deploy to fresh cluster
- [ ] Verify CRDs are installed (`kubectl get crds`)
- [ ] Verify operator pod is running
- [ ] Create test SolrCloud resource
- [ ] Verify SolrCloud reconciles
- [ ] Upgrade operator version
- [ ] Verify existing SolrCloud still works
- [ ] Destroy stack, verify cleanup

## Debugging

Enable verbose logging:

```bash
pulumi up --logtostderr -v=9
```

Check Kubernetes resources:

```bash
kubectl get all -n solr-operator
kubectl describe deployment solr-operator -n solr-operator
kubectl logs -n solr-operator -l app.kubernetes.io/name=solr-operator
```

## Security Considerations

- **RBAC**: Operator requires cluster-admin permissions (defined in Helm chart)
- **Secrets**: Module doesn't create or manage secrets
- **Network**: No network policies created (cluster defaults apply)
- **Image Pull**: Uses public images (no registry credentials needed)

## Performance

Typical deployment times:
- **Fresh install**: 60-90 seconds (CRD download + Helm deploy)
- **Update**: 30-45 seconds (Helm rolling update)
- **Destroy**: 20-30 seconds (cascade deletion)

## Conclusion

This module prioritizes **safety**, **simplicity**, and **production-readiness**:

- Safety through atomic deployments and explicit dependencies
- Simplicity through minimal configuration and convention-over-configuration
- Production-readiness through stable versions and comprehensive error handling

The architecture balances OpenMCF's abstraction goals with Apache Solr Operator's operational requirements, providing a reliable deployment foundation for SolrCloud management.

