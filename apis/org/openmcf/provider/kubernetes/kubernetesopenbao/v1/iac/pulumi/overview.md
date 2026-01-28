# KubernetesOpenBao Pulumi Module - Architecture Overview

## Module Structure

```
iac/pulumi/
├── main.go           # Pulumi entry point
├── Pulumi.yaml       # Pulumi project configuration
├── Makefile          # Build and dependency management
├── debug.sh          # Debugging script
├── README.md         # Module documentation
├── overview.md       # This file
└── module/
    ├── main.go       # Resource orchestration
    ├── locals.go     # Local variables and computed values
    ├── vars.go       # Helm chart configuration
    ├── outputs.go    # Output constants
    ├── namespace.go  # Namespace creation
    └── helm_chart.go # Helm chart installation
```

## Resource Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                    KubernetesOpenBaoStackInput                   │
│  ┌─────────────────┐  ┌─────────────────────────────────────┐  │
│  │ ProviderConfig  │  │            Target                    │  │
│  │ (kubeconfig)    │  │  ┌─────────────────────────────┐    │  │
│  └────────┬────────┘  │  │    KubernetesOpenBaoSpec    │    │  │
│           │           │  │  - namespace                 │    │  │
│           │           │  │  - server_container          │    │  │
│           │           │  │  - high_availability         │    │  │
│           │           │  │  - ingress                   │    │  │
│           │           │  │  - injector                  │    │  │
│           │           │  │  - ui_enabled                │    │  │
│           │           │  └─────────────────────────────┘    │  │
│           │           └─────────────────────────────────────┘  │
└───────────┼────────────────────────────────────────────────────┘
            │
            ▼
┌───────────────────────────────────────────────────────────────────┐
│                      Module Orchestration                          │
│                        (module/main.go)                            │
│                                                                    │
│  1. initializeLocals() ─────────► Compute namespace, labels,       │
│                                   service names, endpoints         │
│                                                                    │
│  2. pulumikubernetesprovider ───► Create Kubernetes provider       │
│                                                                    │
│  3. namespace() ────────────────► Create namespace (if enabled)    │
│                                                                    │
│  4. helmChart() ────────────────► Deploy OpenBao Helm chart        │
│                                                                    │
└───────────────────────────────────────────────────────────────────┘
            │
            ▼
┌───────────────────────────────────────────────────────────────────┐
│                     Kubernetes Resources                           │
│                                                                    │
│  ┌─────────────┐  ┌─────────────────────────────────────────────┐ │
│  │  Namespace  │  │              Helm Release                    │ │
│  │ (optional)  │  │  ┌───────────────────────────────────────┐  │ │
│  └─────────────┘  │  │ OpenBao Server (StatefulSet)          │  │ │
│                   │  │ - Standalone OR HA mode                │  │ │
│                   │  │ - PersistentVolumeClaims               │  │ │
│                   │  │ - ConfigMaps                           │  │ │
│                   │  │ - Services                             │  │ │
│                   │  └───────────────────────────────────────┘  │ │
│                   │  ┌───────────────────────────────────────┐  │ │
│                   │  │ Agent Injector (Deployment, optional) │  │ │
│                   │  │ - MutatingWebhookConfiguration        │  │ │
│                   │  │ - ServiceAccount/RBAC                 │  │ │
│                   │  └───────────────────────────────────────┘  │ │
│                   │  ┌───────────────────────────────────────┐  │ │
│                   │  │ Ingress (optional)                    │  │ │
│                   │  │ - TLS termination                     │  │ │
│                   │  │ - External hostname                   │  │ │
│                   │  └───────────────────────────────────────┘  │ │
│                   │  ┌───────────────────────────────────────┐  │ │
│                   │  │ UI Service (optional)                 │  │ │
│                   │  └───────────────────────────────────────┘  │ │
│                   └─────────────────────────────────────────────┘ │
└───────────────────────────────────────────────────────────────────┘
            │
            ▼
┌───────────────────────────────────────────────────────────────────┐
│                        Stack Outputs                               │
│                                                                    │
│  - namespace              - kube_endpoint                         │
│  - service                - external_hostname                     │
│  - port_forward_command   - api_address                          │
│  - cluster_address        - ha_enabled                           │
└───────────────────────────────────────────────────────────────────┘
```

## Deployment Modes

### Standalone Mode
```
┌────────────────────────────────────────┐
│          OpenBao Server Pod            │
│  ┌─────────────────────────────────┐  │
│  │         openbao container       │  │
│  │  - API: 8200                    │  │
│  │  - Cluster: 8201                │  │
│  └────────────┬────────────────────┘  │
│               │                        │
│  ┌────────────▼────────────────────┐  │
│  │    PersistentVolumeClaim        │  │
│  │    (file storage backend)       │  │
│  └─────────────────────────────────┘  │
└────────────────────────────────────────┘
```

### High Availability Mode (Raft)
```
┌────────────────────────────────────────────────────────────────────┐
│                       Raft Cluster (3+ replicas)                   │
│                                                                    │
│  ┌─────────────┐     ┌─────────────┐     ┌─────────────┐         │
│  │ openbao-0   │◄───►│ openbao-1   │◄───►│ openbao-2   │         │
│  │  (Leader)   │     │ (Standby)   │     │ (Standby)   │         │
│  └──────┬──────┘     └──────┬──────┘     └──────┬──────┘         │
│         │                   │                   │                 │
│  ┌──────▼──────┐     ┌──────▼──────┐     ┌──────▼──────┐         │
│  │    PVC-0    │     │    PVC-1    │     │    PVC-2    │         │
│  │ (Raft data) │     │ (Raft data) │     │ (Raft data) │         │
│  └─────────────┘     └─────────────┘     └─────────────┘         │
│                                                                    │
│  Services:                                                         │
│  - openbao (ClusterIP): Routes to active leader                   │
│  - openbao-active: Selects active node only                       │
│  - openbao-standby: Selects standby nodes                         │
│  - openbao-internal (Headless): Pod-to-pod communication          │
└────────────────────────────────────────────────────────────────────┘
```

## Helm Values Mapping

| Spec Field | Helm Value | Description |
|------------|------------|-------------|
| `server_container.replicas` | `server.ha.replicas` | Number of server replicas (HA mode) |
| `server_container.resources` | `server.resources` | CPU/memory allocation |
| `server_container.data_storage_size` | `server.dataStorage.size` | PVC size |
| `high_availability.enabled` | `server.ha.enabled` | Enable HA mode |
| `high_availability.replicas` | `server.ha.replicas` | HA replica count |
| `ingress.enabled` | `server.ingress.enabled` | Enable ingress |
| `ingress.hostname` | `server.ingress.hosts[].host` | Ingress hostname |
| `ui_enabled` | `ui.enabled` | Enable UI service |
| `injector.enabled` | `injector.enabled` | Enable agent injector |
| `tls_enabled` | `global.tlsDisable` (inverted) | TLS configuration |

## Key Implementation Details

### Labels
All resources are tagged with standard OpenMCF labels:
- `planton-cloud-resource: "true"`
- `planton-cloud-resource-name: <name>`
- `planton-cloud-resource-kind: KubernetesOpenBao`
- `planton-cloud-organization: <org>` (if set)
- `planton-cloud-environment: <env>` (if set)

### Service Discovery
- Internal: `<name>.<namespace>.svc.cluster.local:8200`
- Active (HA): `<name>-active.<namespace>.svc.cluster.local:8200`
- Standby (HA): `<name>-standby.<namespace>.svc.cluster.local:8200`

### Port Forwarding
For local development access:
```bash
kubectl port-forward -n <namespace> service/<name> 8200:8200
```

## Post-Deployment Requirements

After deployment, OpenBao requires initialization:
1. Initialize: `bao operator init`
2. Unseal: `bao operator unseal` (repeat with threshold keys)
3. Login: `bao login <root-token>`
4. Configure auth methods and policies
