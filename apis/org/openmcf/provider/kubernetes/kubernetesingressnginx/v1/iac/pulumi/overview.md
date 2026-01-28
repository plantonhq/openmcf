# Kubernetes Ingress NGINX - Pulumi Module Overview

## Architecture

The Kubernetes Ingress NGINX Pulumi module deploys the official NGINX Ingress Controller with cloud-native integrations. The module automatically configures load balancers, networking, and cloud-specific features based on the target platform.

## Architecture Diagram

```
┌──────────────────────────────────────────────────────────────┐
│                     Kubernetes Cluster                        │
│                                                               │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │        Namespace: kubernetes-ingress-nginx              │ │
│  │                                                         │ │
│  │  ┌───────────────────────────────────────────────────┐ │ │
│  │  │         NGINX Ingress Controller                  │ │ │
│  │  │                                                   │ │ │
│  │  │  ┌──────────┐  ┌──────────┐  ┌──────────┐       │ │ │
│  │  │  │   Pod    │  │   Pod    │  │   Pod    │       │ │ │
│  │  │  │    1     │  │    2     │  │    3     │       │ │ │
│  │  │  └──────────┘  └──────────┘  └──────────┘       │ │ │
│  │  │                                                   │ │ │
│  │  │  - Watches Ingress resources                     │ │ │
│  │  │  - Configures NGINX dynamically                  │ │ │
│  │  │  - Routes traffic to services                    │ │ │
│  │  └───────────────────────────────────────────────────┘ │ │
│  │                            │                           │ │
│  │                            │                           │ │
│  │                   ┌────────▼────────┐                  │ │
│  │                   │  LoadBalancer   │                  │ │
│  │                   │    Service      │                  │ │
│  │                   └────────┬────────┘                  │ │
│  └────────────────────────────┼──────────────────────────┘ │
│                                │                            │
└────────────────────────────────┼────────────────────────────┘
                                 │
                     ┌───────────▼──────────┐
                     │  Cloud Load Balancer │
                     │                      │
                     │  • GCP: GCLB         │
                     │  • AWS: NLB          │
                     │  • Azure: Azure LB   │
                     └──────────────────────┘
```

## Component Flow

### 1. Initialization

- Receive KubernetesIngressNginxStackInput with target cluster and spec
- Initialize locals with resource labels and computed values
- Set up Kubernetes provider from cluster credentials

### 2. Namespace Creation

- Create dedicated namespace: `kubernetes-ingress-nginx`
- Apply resource labels for organization/environment tracking
- Use as parent for all subsequent resources

### 3. Load Balancer Annotation Logic

The module determines load balancer configuration based on:

**GKE:**
- External: `cloud.google.com/load-balancer-type: external`
- Internal: `cloud.google.com/load-balancer-type: internal`

**EKS:**
- External: `service.beta.kubernetes.io/aws-load-balancer-scheme: internet-facing`
- Internal: `service.beta.kubernetes.io/aws-load-balancer-scheme: internal`

**AKS:**
- Internal: `service.beta.kubernetes.io/azure-load-balancer-internal: true`
- External: No annotation needed (default)

### 4. Helm Chart Deployment

- Deploy official `ingress-nginx` Helm chart
- Configure controller service type as LoadBalancer
- Set ingress class as default
- Enable watching ingresses without explicit class
- Apply cloud-specific service annotations

### 5. Stack Outputs

Export deployment details:
- Namespace name
- Helm release name
- Service name (controller service)
- Service type (LoadBalancer)

## Module Design

### File Organization

- **`main.go`**: Orchestrates deployment sequence
  - Provider setup
  - Namespace creation
  - Service annotation logic
  - Helm chart deployment

- **`locals.go`**: Initializes local variables
  - Resource labels
  - Chart version selection
  - Service name calculation
  - Stack output exports

- **`outputs.go`**: Defines output constant names
  - OpNamespace
  - OpReleaseName
  - OpServiceName
  - OpServiceType

- **`vars.go`**: Deployment constants
  - Default namespace
  - Helm chart repository URL
  - Default chart version

### Cloud Provider Detection

The module uses Go's type switch on the `provider_config` oneof:

```go
if gke := spec.GetGke(); gke != nil {
    // Apply GKE-specific configuration
} else if eks := spec.GetEks(); eks != nil {
    // Apply EKS-specific configuration
} else if aks := spec.GetAks(); aks != nil {
    // Apply AKS-specific configuration
}
```

This ensures only one provider configuration is active at a time.

## Data Flow

### Ingress Request Flow

1. **External Request** → Cloud Load Balancer
2. **Load Balancer** → NGINX Controller Service (LoadBalancer type)
3. **Service** → NGINX Controller Pods (via label selector)
4. **Controller** → Routes to backend service based on Ingress rules
5. **Backend Service** → Application Pods

### Configuration Flow

1. **User** → Defines KubernetesIngressNginx resource
2. **Pulumi Module** → Reads spec and determines cloud provider
3. **Module** → Applies cloud-specific annotations
4. **Helm Chart** → Deploys controller with annotations
5. **Cloud Controller** → Provisions load balancer with specified configuration

## Cloud-Specific Implementations

### GKE Implementation

**Features:**
- Global or regional load balancers
- Static IP assignment via `static_ip_name`
- Internal LB requires subnetwork specification
- Automatic Cloud Armor integration available

**Annotations Applied:**
- `cloud.google.com/load-balancer-type`
- Additional annotations if static IP configured

### EKS Implementation

**Features:**
- Network Load Balancer (NLB) for performance
- Security group attachment
- Subnet placement control
- IAM Roles for Service Accounts (IRSA) support

**Annotations Applied:**
- `service.beta.kubernetes.io/aws-load-balancer-scheme`
- Security groups via controller configuration

### AKS Implementation

**Features:**
- Azure Load Balancer integration
- Workload Identity support
- Public IP resource reuse
- Internal LB for private access

**Annotations Applied:**
- `service.beta.kubernetes.io/azure-load-balancer-internal`
- Additional annotations for managed identity

## Design Decisions

### Why Helm Chart?

- **Community Standard**: Official deployment method for ingress-nginx
- **Battle-Tested**: Used by thousands of production clusters
- **Feature Complete**: Includes all controller capabilities
- **Actively Maintained**: Regular updates from Kubernetes community

### Why LoadBalancer Service?

- **Cloud Native**: Leverages cloud provider load balancers
- **Production Ready**: Built for high-traffic scenarios
- **HA Support**: Load balances across multiple controller pods
- **Health Checks**: Automatic endpoint health monitoring

### Why Fixed Namespace?

- **Standardization**: Consistent namespace across deployments
- **RBAC Simplification**: Well-known namespace for policy definitions
- **Operational Clarity**: Easy to find and monitor
- **Convention**: Matches upstream Helm chart defaults

### Why vars.go Instead of Hardcoding?

- **Maintainability**: Single source of truth for constants
- **Upgradability**: Easy chart version updates
- **Testability**: Can override values in tests
- **Documentation**: Self-documenting constant values

## Monitoring and Observability

### Prometheus Metrics

The controller exposes metrics on `/metrics` endpoint:

- Request rate and latency
- Error rates by ingress
- SSL certificate expiration
- Configuration reload statistics

### Logs

Controller logs are available via:

```bash
kubectl logs -n kubernetes-ingress-nginx -l app.kubernetes.io/component=controller
```

### Stack Outputs

Use outputs for automation:

```go
ctx.Export("namespace", pulumi.String(locals.Namespace))
ctx.Export("service-name", pulumi.String(locals.ServiceName))
```

## High Availability

For production HA deployments:

1. **Multiple Replicas**: Configure via Helm values
2. **Pod Anti-Affinity**: Spread across nodes/zones
3. **Pod Disruption Budget**: Maintain minimum available replicas
4. **Resource Limits**: Set appropriate CPU/memory limits
5. **Health Checks**: Configure readiness/liveness probes

## Security Considerations

1. **Network Policies**: Restrict controller network access
2. **RBAC**: Controller has minimal necessary permissions
3. **TLS Termination**: Controller handles SSL/TLS
4. **Rate Limiting**: Configure via annotations on ingresses
5. **ModSecurity WAF**: Can be enabled via Helm values

## Troubleshooting

### Common Issues

1. **Load Balancer Not Provisioned**
   - Check cloud provider quotas
   - Verify subnet configuration
   - Review security groups
   - Check service annotations

2. **Ingress Not Working**
   - Verify ingress class is set correctly
   - Check controller logs for errors
   - Ensure backend services exist
   - Validate ingress resource spec

3. **SSL Certificate Issues**
   - Check cert-manager integration
   - Verify certificate secrets exist
   - Review TLS configuration in ingress

## References

- [NGINX Ingress Controller Documentation](https://kubernetes.github.io/ingress-nginx/)
- [Helm Chart Values](https://github.com/kubernetes/ingress-nginx/blob/main/charts/ingress-nginx/values.yaml)
- [Pulumi Kubernetes Provider](https://www.pulumi.com/registry/packages/kubernetes/)

