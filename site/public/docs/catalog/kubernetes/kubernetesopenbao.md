---
title: "Openbao"
description: "Openbao deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesopenbao"
---

# OpenBao on Kubernetes - Technical Research

## Overview

OpenBao is an open-source secrets management solution that is a community-driven fork of HashiCorp Vault, created after Vault's license change to BSL (Business Source License). OpenBao is developed under the OpenSSF (Open Source Security Foundation) umbrella and maintains API compatibility with Vault while being released under the MPL-2.0 license.

## Architecture

### Core Components

1. **OpenBao Server**: The main secrets management engine
2. **Storage Backend**: Persistent storage for encrypted secrets
3. **Agent Injector**: Kubernetes admission webhook for automatic secret injection
4. **CSI Provider**: Secrets Store CSI driver integration (optional)

### Deployment Modes

#### Standalone Mode
- Single server instance
- File-based storage backend
- Suitable for development and small deployments
- No built-in redundancy

#### High Availability (HA) Mode
- Multiple server instances (3+ recommended)
- Raft integrated storage for consensus
- Automatic leader election
- Data replication across all nodes

## Helm Chart Details

### Repository Information
- **Chart Repository**: https://openbao.github.io/openbao-helm
- **OCI Registry**: oci://ghcr.io/openbao/charts/openbao
- **Chart Version**: 0.23.3 (as of implementation)
- **App Version**: v2.4.4

### Key Helm Values

#### Server Configuration
```yaml
server:
  enabled: true
  image:
    registry: "quay.io"
    repository: "openbao/openbao"
    tag: ""  # defaults to appVersion
  standalone:
    enabled: "-"  # enabled when HA is disabled
    config: |
      ui = true
      listener "tcp" {
        tls_disable = 1
        address = "[::]:8200"
        cluster_address = "[::]:8201"
      }
      storage "file" {
        path = "/openbao/data"
      }
  ha:
    enabled: false
    replicas: 3
    raft:
      enabled: false
      config: |
        ui = true
        listener "tcp" {
          tls_disable = 1
          address = "[::]:8200"
          cluster_address = "[::]:8201"
        }
        storage "raft" {
          path = "/openbao/data"
        }
        service_registration "kubernetes" {}
  dataStorage:
    enabled: true
    size: 10Gi
    storageClass: null
    accessMode: ReadWriteOnce
```

#### Agent Injector Configuration
```yaml
injector:
  enabled: "-"  # follows global.enabled
  replicas: 1
  image:
    registry: "docker.io"
    repository: "hashicorp/vault-k8s"
    tag: "1.7.2"
  agentImage:
    registry: "quay.io"
    repository: "openbao/openbao"
```

#### UI Configuration
```yaml
ui:
  enabled: false
  serviceType: "ClusterIP"
  externalPort: 8200
```

## Network Ports

| Port | Protocol | Purpose |
|------|----------|---------|
| 8200 | TCP | API and UI access |
| 8201 | TCP | Cluster communication (HA mode) |

## Storage Backends

### File Backend (Standalone)
- Simple file-based storage
- Single node only
- Good for development

### Raft Backend (HA)
- Built-in consensus protocol
- Leader election
- Data replication
- Recommended for production

### External Backends (via config override)
- Consul
- PostgreSQL
- MySQL
- DynamoDB
- And more

## Security Considerations

### Seal/Unseal Process
OpenBao starts in a sealed state and requires unseal keys to decrypt the master key:
- **Shamir's Secret Sharing**: Default method, splits master key into shares
- **Auto-Unseal**: Uses external KMS (GCP, AWS, Azure) to automatically unseal

### TLS Configuration
- TLS disabled by default (`global.tlsDisable: true`)
- Can be enabled for production deployments
- Supports custom certificates

### Authentication Methods
- Kubernetes Service Account
- Token-based
- LDAP
- OIDC
- And more

## Initialization Process

After deployment, OpenBao must be initialized:

```bash
# Initialize with 5 key shares, 3 required to unseal
bao operator init -key-shares=5 -key-threshold=3

# Unseal (repeat 3 times with different keys)
bao operator unseal <key>

# Login with root token
bao login <root-token>
```

## Kubernetes Integration

### Service Account Auth
```bash
# Enable Kubernetes auth
bao auth enable kubernetes

# Configure auth
bao write auth/kubernetes/config \
    kubernetes_host="https://$KUBERNETES_PORT_443_TCP_ADDR:443"

# Create role
bao write auth/kubernetes/role/app \
    bound_service_account_names=app-sa \
    bound_service_account_namespaces=default \
    policies=app-policy \
    ttl=24h
```

### Agent Injection Annotations
```yaml
annotations:
  vault.hashicorp.com/agent-inject: "true"
  vault.hashicorp.com/role: "app"
  vault.hashicorp.com/agent-inject-secret-config: "secret/data/app/config"
```

## Monitoring and Telemetry

### Prometheus Integration
```yaml
serverTelemetry:
  serviceMonitor:
    enabled: true
    interval: 30s
    scrapeTimeout: 10s
  prometheusRules:
    enabled: true
```

### Grafana Dashboard
The Helm chart includes a pre-built Grafana dashboard for monitoring OpenBao metrics.

## Best Practices

### Production Deployment
1. **Use HA Mode**: Deploy at least 3 replicas for fault tolerance
2. **Enable TLS**: Secure all communications with TLS
3. **Auto-Unseal**: Configure auto-unseal for operational simplicity
4. **Audit Logging**: Enable audit devices for compliance
5. **Backup Strategy**: Regular snapshots of Raft storage

### Resource Sizing
| Deployment Size | CPU Request | Memory Request | Storage |
|-----------------|-------------|----------------|---------|
| Development | 100m | 128Mi | 1Gi |
| Small | 250m | 256Mi | 10Gi |
| Medium | 500m | 512Mi | 50Gi |
| Large | 1000m | 1Gi | 100Gi |

## References

- [OpenBao Documentation](https://openbao.org/docs/)
- [OpenBao Helm Chart](https://github.com/openbao/openbao-helm)
- [OpenBao GitHub Repository](https://github.com/openbao/openbao)
- [Kubernetes Auth Method](https://openbao.org/docs/auth/kubernetes)
- [Agent Injector](https://openbao.org/docs/platform/k8s/injector)
