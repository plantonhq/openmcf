# Pulumi Module Overview

## Architecture

```
┌──────────────────────────────────────────────────────────────────┐
│                        Pulumi Module                             │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │  main.go (entrypoint)                                    │   │
│  │  - Loads stack input from environment                    │   │
│  │  - Calls module.Resources()                              │   │
│  └──────────────────────────────────────────────────────────┘   │
│                             │                                    │
│                             ▼                                    │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │  module/main.go (orchestrator)                           │   │
│  │  - Sets up Kubernetes provider                           │   │
│  │  - Creates namespace (if enabled)                        │   │
│  │  - Deploys Helm release                                  │   │
│  └──────────────────────────────────────────────────────────┘   │
│                             │                                    │
│          ┌──────────────────┼──────────────────┐                │
│          ▼                  ▼                  ▼                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │  locals.go   │  │  outputs.go  │  │   vars.go    │          │
│  │  - Computed  │  │  - Output    │  │  - Helm repo │          │
│  │    values    │  │    constants │  │  - Chart name│          │
│  │  - Helm vals │  │              │  │              │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└──────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│                    Kubernetes Cluster                            │
│                                                                  │
│  ┌─────────────┐  ┌─────────────────────────────────────────┐   │
│  │  Namespace  │  │           Helm Release                  │   │
│  │  (optional) │  │  - CephCluster                          │   │
│  └─────────────┘  │  - CephBlockPool(s) + StorageClass(es)  │   │
│                   │  - CephFilesystem(s) + StorageClass(es) │   │
│                   │  - CephObjectStore(s) + StorageClass(es)│   │
│                   └─────────────────────────────────────────────┘│
└──────────────────────────────────────────────────────────────────┘
```

## Key Design Decisions

### 1. Helm-Based Deployment

Uses the official `rook-ceph-cluster` Helm chart for:
- Validated configuration
- Upstream support
- Easy upgrades

### 2. Computed Helm Values

The `locals.go` transforms the protobuf spec into Helm values:
- Block pools → `cephBlockPools`
- Filesystems → `cephFileSystems`
- Object stores → `cephObjectStores`

### 3. StorageClass Integration

Each storage pool can optionally create a StorageClass:
- Block pools → RBD StorageClass
- Filesystems → CephFS StorageClass
- Object stores → Bucket StorageClass

### 4. Namespace Management

Namespace creation is optional (`create_namespace` flag):
- `true`: Module creates the namespace
- `false`: Assumes namespace exists

## Data Flow

```
manifest.yaml
     │
     ▼
KubernetesRookCephClusterStackInput (protobuf)
     │
     ▼
initializeLocals() → Locals struct
     │
     ├─→ Labels (common labels for all resources)
     ├─→ HelmValues (transformed configuration)
     └─→ Storage names (for outputs)
     │
     ▼
helm.NewRelease() → Kubernetes resources
```

## Module Files

| File | Purpose |
|------|---------|
| `main.go` | Entrypoint, loads stack input |
| `module/main.go` | Resource orchestration |
| `module/locals.go` | Value transformations |
| `module/outputs.go` | Output key constants |
| `module/vars.go` | Helm chart configuration |

## Extending the Module

To add new features:

1. Update `spec.proto` with new fields
2. Regenerate Go stubs
3. Update `locals.go` to handle new fields
4. Update `buildHelmValues()` to include in Helm values
5. Add tests in `spec_test.go`
