# KubernetesDaemonSet Pulumi Module Overview

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    KubernetesDaemonSet                       │
│                     Stack Input                              │
├─────────────────────────────────────────────────────────────┤
│  target: KubernetesDaemonSet                                │
│  provider_config: KubernetesProviderConfig                  │
│  docker_config_json: string (optional)                      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Pulumi Module                             │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌────────────────┐    ┌────────────────┐                   │
│  │   Namespace    │    │    Secret      │                   │
│  │  (optional)    │    │ (env-secrets)  │                   │
│  └────────────────┘    └────────────────┘                   │
│                                                              │
│  ┌────────────────┐    ┌────────────────┐                   │
│  │   DaemonSet    │    │    Secret      │                   │
│  │                │    │ (image-pull)   │                   │
│  │  ┌──────────┐  │    │   (optional)   │                   │
│  │  │   Pod    │  │    └────────────────┘                   │
│  │  │ Template │  │                                         │
│  │  └──────────┘  │                                         │
│  └────────────────┘                                         │
│                                                              │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Kubernetes Cluster                        │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│   Node 1          Node 2          Node 3          Node N    │
│  ┌──────┐        ┌──────┐        ┌──────┐        ┌──────┐  │
│  │ Pod  │        │ Pod  │        │ Pod  │   ...  │ Pod  │  │
│  └──────┘        └──────┘        └──────┘        └──────┘  │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

## Resource Flow

1. **initializeLocals**: Parse stack input and compute derived values
2. **Namespace Creation**: If `create_namespace: true`, create the namespace
3. **Secret Creation**: Create secret for environment secrets
4. **Image Pull Secret**: If docker config provided, create image pull secret
5. **DaemonSet Creation**: Create the DaemonSet with all configured options

## Key Features

### Node Selection
- `node_selector`: Key-value pairs to match node labels
- `tolerations`: Allow scheduling on tainted nodes

### Volume Mounts
- Host path volumes for accessing node-level files
- Support for different host path types (Directory, File, Socket, etc.)

### Security Context
- Privileged containers for node-level access
- Linux capabilities management
- User/Group ID configuration

### Update Strategy
- RollingUpdate: Progressive pod updates
- OnDelete: Manual pod updates

## Configuration Examples

### Minimal Configuration
```yaml
spec:
  namespace:
    value: default
  container:
    app:
      image:
        repo: busybox
        tag: latest
```

### Full Configuration
```yaml
spec:
  namespace:
    value: monitoring
  create_namespace: true
  container:
    app:
      image:
        repo: prom/node-exporter
        tag: v1.7.0
      resources:
        limits:
          cpu: "500m"
          memory: "256Mi"
        requests:
          cpu: "100m"
          memory: "128Mi"
      volume_mounts:
        - name: proc
          mount_path: /host/proc
          host_path: /proc
          read_only: true
      ports:
        - name: metrics
          container_port: 9100
          network_protocol: TCP
      security_context:
        privileged: true
  node_selector:
    kubernetes.io/os: linux
  tolerations:
    - operator: Exists
  update_strategy:
    type: RollingUpdate
    rolling_update:
      max_unavailable: "1"
```

