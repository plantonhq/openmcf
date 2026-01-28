# Overview

The **KubernetesOpenBao** API resource enables deployment and management of OpenBao on Kubernetes clusters. OpenBao is an open-source secrets management solution forked from HashiCorp Vault, providing secure secret storage, dynamic secrets generation, data encryption, leasing/renewal, and revocation capabilities.

## Why We Created This API Resource

Managing secrets securely in cloud-native environments is critical but complex. OpenBao provides enterprise-grade secrets management, but deploying and configuring it on Kubernetes requires significant expertise. The KubernetesOpenBao API resource simplifies this process by providing:

- **Simplified Deployment**: Deploy OpenBao on Kubernetes with minimal configuration using sensible defaults
- **Flexible Architecture**: Support for both standalone (development/testing) and high-availability (production) deployment modes
- **Integrated Storage**: Automatic configuration of persistent storage using file backend (standalone) or Raft consensus (HA)
- **Secret Injection**: Optional Agent Injector deployment for automatic secret injection into application pods
- **External Access**: Built-in ingress configuration for secure external access

## Key Features

### Deployment Modes

#### Standalone Mode
- **Single Replica**: Ideal for development, testing, and small-scale deployments
- **File Storage Backend**: Persistent data storage using Kubernetes PersistentVolumeClaims
- **Quick Setup**: Minimal configuration required to get started

#### High Availability Mode
- **Multi-Replica Clustering**: Deploy 3 or more replicas for fault tolerance
- **Raft Integrated Storage**: Built-in consensus protocol for leader election and data replication
- **Automatic Failover**: Seamless leader election when the active node fails

### Container Configuration

- **Resource Management**: Fine-grained control over CPU and memory allocation
  - Default CPU: 100m requests, 500m limits
  - Default Memory: 128Mi requests, 256Mi limits
- **Data Storage**: Configurable persistent volume size (default: 10Gi)
- **Replica Count**: Adjustable based on deployment mode (1 for standalone, 3+ for HA)

### Agent Injector

The OpenBao Agent Injector provides automatic secret injection into Kubernetes pods:

- **Mutating Webhook**: Automatically injects OpenBao Agent sidecars into annotated pods
- **Secret Templating**: Transform secrets into application-specific formats
- **Token Management**: Automatic authentication and token renewal

### Ingress Configuration

- **External Access**: Configure ingress for secure external access to the OpenBao UI and API
- **Custom Hostname**: Full control over the ingress hostname (e.g., `openbao.example.com`)
- **TLS Support**: Optional TLS termination at the ingress controller
- **Ingress Class**: Support for different ingress controllers (nginx, traefik, etc.)

### Namespace Management

- **Flexible Namespace Control**: The `create_namespace` flag provides control over namespace creation:
  - **When `true`**: Creates a dedicated namespace with appropriate resource labels
  - **When `false`**: Uses an existing namespace (must exist before deployment)

## Security Features

- **TLS Encryption**: Optional end-to-end TLS encryption for all OpenBao traffic
- **Auto-Unseal**: Support for various auto-unseal mechanisms through Helm value overrides
- **Audit Logging**: Configurable audit storage for compliance requirements

## Outputs

After deployment, the following outputs are available:

- **namespace**: The Kubernetes namespace where OpenBao is deployed
- **service**: The Kubernetes service name for accessing OpenBao
- **kube_endpoint**: Internal cluster endpoint (FQDN)
- **external_hostname**: External hostname when ingress is enabled
- **port_forward_command**: kubectl command for local access
- **root_token_secret**: Reference to the root token Kubernetes secret
- **unseal_keys_secret**: Reference to the unseal keys Kubernetes secret

## Benefits

- **Open Source**: Community-driven development under the OpenSSF umbrella
- **Vault Compatible**: API-compatible with HashiCorp Vault, enabling migration from existing deployments
- **Production Ready**: Enterprise-grade features including HA, audit logging, and encryption
- **Kubernetes Native**: Deep integration with Kubernetes authentication and service accounts
