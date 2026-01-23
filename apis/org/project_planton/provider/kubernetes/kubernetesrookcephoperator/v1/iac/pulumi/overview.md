# KubernetesRookCephOperator Pulumi Module Architecture

## Module Structure

```
pulumi/
├── main.go              # Entry point - loads stack input and calls module
├── Pulumi.yaml          # Pulumi project configuration
├── Makefile             # Build and test automation
├── README.md            # Module usage documentation
├── overview.md          # This file - architecture overview
└── module/
    ├── main.go          # Core resource creation logic
    ├── locals.go        # Computed values and transformations
    ├── outputs.go       # Output constant definitions
    └── vars.go          # Static configuration (chart names, URLs)
```

## Data Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                      Stack Input                                 │
│  (KubernetesRookCephOperatorStackInput from environment)        │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     main.go (entrypoint)                         │
│  - Loads stack input from STACK_INPUT environment variable      │
│  - Deserializes JSON into protobuf message                      │
│  - Calls module.Resources()                                      │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                   module/main.go                                 │
│  - Initializes locals (computed values)                         │
│  - Sets up Kubernetes provider                                  │
│  - Creates namespace (if requested)                             │
│  - Deploys Helm release                                         │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                 Kubernetes Resources                             │
│  - Namespace (optional)                                         │
│  - Helm Release (rook-ceph operator)                            │
│    └── Operator Deployment                                      │
│    └── CRDs (CephCluster, CephBlockPool, etc.)                 │
│    └── RBAC (ClusterRoles, ServiceAccounts)                    │
│    └── CSI Drivers (RBD, CephFS)                               │
└─────────────────────────────────────────────────────────────────┘
```

## Key Components

### Entrypoint (main.go)

The entrypoint:
1. Creates Pulumi program context
2. Loads `KubernetesRookCephOperatorStackInput` from environment
3. Delegates to `module.Resources()` for actual resource creation

```go
func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        stackInput := &kubernetesrookcephoperatorv1.KubernetesRookCephOperatorStackInput{}
        if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
            return err
        }
        return module.Resources(ctx, stackInput)
    })
}
```

### Locals (module/locals.go)

Computes derived values from the input:
- **Namespace**: Extracted from spec
- **Labels**: Standard Project Planton labels
- **HelmReleaseName**: Based on metadata.name
- **ChartVersion**: Strips 'v' prefix from operator_version
- **HelmValues**: Transforms spec into Helm values map

### Resource Creation (module/main.go)

Creates resources in order:

1. **Kubernetes Provider**: Sets up authentication from `provider_config`
2. **Namespace** (conditional): Created if `create_namespace: true`
3. **Helm Release**: Deploys rook-ceph chart with computed values

### Outputs

The module exports:
- `namespace`: Where operator is deployed
- `helm_release_name`: Helm release identifier
- `webhook_service`: Webhook service name

## Helm Values Mapping

The module translates spec fields to Helm values:

| Spec Field | Helm Value |
|------------|------------|
| `crds_enabled` | `crds.enabled` |
| `container.resources` | `resources` |
| `csi.enable_rbd_driver` | `csi.enableRbdDriver` |
| `csi.enable_cephfs_driver` | `csi.enableCephfsDriver` |
| `csi.disable_csi_driver` | `csi.disableCsiDriver` |
| `csi.enable_csi_host_network` | `csi.enableCSIHostNetwork` |
| `csi.provisioner_replicas` | `csi.provisionerReplicas` |
| `csi.enable_csi_addons` | `csi.csiAddons.enabled` |
| `csi.enable_nfs_driver` | `csi.nfs.enabled` |

## Design Decisions

### Why Helm Instead of Direct Resources?

The Rook Ceph Operator includes:
- Complex RBAC configurations
- Multiple CRDs
- CSI driver deployments
- Webhook configurations

Using Helm:
- Leverages official, tested templates
- Handles version-specific differences
- Simplifies upgrades

### Namespace Deletion Policy

The namespace uses `background` deletion policy to prevent timeout issues during `pulumi destroy`. The Helm release and CRDs can have finalizers that block foreground deletion.

### Atomic Helm Release

`Atomic: true` ensures the release rolls back on failure, maintaining cluster consistency.

## Extension Points

### Adding New CSI Options

1. Add field to `KubernetesRookCephOperatorCsiSpec` in `spec.proto`
2. Regenerate proto stubs
3. Map field in `buildHelmValues()` in `locals.go`

### Custom Resource Additions

To add resources beyond Helm:
1. Create new `.go` file in `module/`
2. Import and call from `main.go` after Helm release

## Testing

Local testing without a cluster:

```bash
# Build check
make build

# Preview (requires stack setup)
make test
```

Integration testing requires a Kubernetes cluster with raw block devices.
