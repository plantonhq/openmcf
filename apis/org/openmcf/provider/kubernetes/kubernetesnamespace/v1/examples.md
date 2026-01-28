# Kubernetes Namespace Examples

This document provides practical examples for deploying Kubernetes namespaces using OpenMCF.

## Understanding Target Cluster

All examples include the `target_cluster` field, which specifies which Kubernetes cluster the namespace should be created in. This field is required and consists of:

- `cluster_kind`: The type of Kubernetes cluster (e.g., `gcp_gke_cluster`, `aws_eks_cluster`, `azure_aks_cluster`, etc.)
- `cluster_name`: The name (or slug) of the Kubernetes cluster in your environment

**Supported cluster kinds:**
- `azure_aks_cluster` (400) - Azure AKS
- `aws_eks_cluster` (207) - AWS EKS
- `gcp_gke_cluster` (607) - GCP GKE
- `digital_ocean_kubernetes_cluster` (1208) - DigitalOcean Kubernetes
- `civo_kubernetes_cluster` (1507) - Civo Kubernetes

The cluster specified must exist in the same environment as the namespace resource.

## Example 1: Minimal Namespace

The simplest possible namespace configuration with just a name:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesNamespace
metadata:
  name: dev-team-alpha
spec:
  target_cluster:
    cluster_kind: gcp_gke_cluster
    cluster_name: prod-cluster
  name: dev-team-alpha
```

**CLI Commands:**

```bash
# Validate the manifest
openmcf validate --manifest minimal-namespace.yaml

# Deploy with Pulumi
openmcf pulumi up \
  --manifest minimal-namespace.yaml \
  --stack myorg/myproject/dev

# Deploy with Terraform
openmcf tofu apply \
  --manifest minimal-namespace.yaml \
  --auto-approve
```

## Example 2: Development Namespace with Small Resource Profile

Namespace for a development team with predefined resource limits:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesNamespace
metadata:
  name: dev-namespace
spec:
  target_cluster:
    cluster_kind: gcp_gke_cluster
    cluster_name: dev-cluster
  name: dev-namespace
  labels:
    team: platform-engineering
    environment: development
    cost-center: engineering
  resource_profile:
    preset: small
  pod_security_standard: baseline
```

**What this creates:**

- Namespace with 2-4 CPU cores and 4-8Gi memory quota
- Baseline pod security enforcement
- Labels for cost tracking

## Example 3: Production Namespace with Large Profile and Network Isolation

Production namespace with strict network controls:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesNamespace
metadata:
  name: prod-api-services
spec:
  target_cluster:
    cluster_kind: gcp_gke_cluster
    cluster_name: prod-cluster
  name: prod-api-services
  labels:
    team: backend-team
    environment: production
    cost-center: engineering
    criticality: high
  resource_profile:
    preset: large
  network_config:
    isolate_ingress: true
    restrict_egress: true
    allowed_ingress_namespaces:
      - istio-system
      - monitoring
    allowed_egress_cidrs:
      - "10.0.0.0/8"
    allowed_egress_domains:
      - "api.stripe.com"
      - "api.twilio.com"
  pod_security_standard: restricted
```

**What this creates:**

- Large resource quota (8-16 CPU, 16-32Gi memory)
- Default-deny network policy with explicit allow lists
- Restricted pod security (maximum hardening)
- External API access only to Stripe and Twilio

## Example 4: Staging Namespace with Istio Service Mesh

Namespace configured for Istio service mesh with revision tag for safe upgrades:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesNamespace
metadata:
  name: staging-services
spec:
  target_cluster:
    cluster_kind: gcp_gke_cluster
    cluster_name: staging-cluster
  name: staging-services
  labels:
    team: platform-team
    environment: staging
    cost-center: engineering
  annotations:
    description: "Staging environment for microservices"
  resource_profile:
    preset: medium
  service_mesh_config:
    enabled: true
    mesh_type: istio
    revision_tag: "prod-stable"
  pod_security_standard: baseline
```

**What this creates:**

- Medium resource quota (4-8 CPU, 8-16Gi memory)
- Automatic Istio sidecar injection enabled
- Uses "prod-stable" revision tag (allows mesh upgrades without changing manifest)
- Baseline pod security

## Example 5: Ephemeral PR Environment with TTL

Temporary namespace for pull request preview with automatic cleanup:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesNamespace
metadata:
  name: pr-1234-preview
spec:
  target_cluster:
    cluster_kind: gcp_gke_cluster
    cluster_name: preview-cluster
  name: pr-1234-preview
  labels:
    team: frontend-team
    environment: preview
    pr-number: "1234"
  annotations:
    janitor/ttl: "24h"
    description: "Temporary environment for PR #1234"
  resource_profile:
    preset: small
```

**What this creates:**

- Small resource quota for cost efficiency
- TTL annotation signals to cleanup operator to delete after 24 hours
- Labels for tracking preview environments

## Example 6: Custom Resource Quotas (Advanced)

Namespace with precise custom resource allocations:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesNamespace
metadata:
  name: custom-quota-namespace
spec:
  target_cluster:
    cluster_kind: gcp_gke_cluster
    cluster_name: prod-cluster
  name: custom-quota-namespace
  labels:
    team: data-science
    environment: production
  resource_profile:
    custom:
      cpu:
        requests: "10"
        limits: "20"
      memory:
        requests: "20Gi"
        limits: "40Gi"
      object_counts:
        pods: 50
        services: 20
        configmaps: 100
        secrets: 100
        persistent_volume_claims: 10
        load_balancers: 3
      default_limits:
        default_cpu_request: "250m"
        default_cpu_limit: "1000m"
        default_memory_request: "256Mi"
        default_memory_limit: "1Gi"
```

**What this creates:**

- Custom CPU quota: 10 cores requested, 20 cores limit
- Custom memory quota: 20Gi requested, 40Gi limit
- Limits on Kubernetes object counts
- Default resource requests/limits for containers that don't specify their own

## Example 7: Linkerd Service Mesh Integration

Namespace with Linkerd instead of Istio:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesNamespace
metadata:
  name: linkerd-services
spec:
  target_cluster:
    cluster_kind: gcp_gke_cluster
    cluster_name: prod-cluster
  name: linkerd-services
  labels:
    team: backend-team
    environment: production
  resource_profile:
    preset: medium
  service_mesh_config:
    enabled: true
    mesh_type: linkerd
  pod_security_standard: baseline
```

**What this creates:**

- Linkerd sidecar injection enabled
- Appropriate labels/annotations for Linkerd

## Example 8: Maximum Security Namespace

Namespace with maximum isolation and security for sensitive workloads:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesNamespace
metadata:
  name: secure-banking-services
spec:
  target_cluster:
    cluster_kind: gcp_gke_cluster
    cluster_name: prod-cluster
  name: secure-banking-services
  labels:
    team: security-team
    environment: production
    compliance: "pci-dss"
    criticality: critical
  resource_profile:
    preset: large
  network_config:
    isolate_ingress: true
    restrict_egress: true
    allowed_ingress_namespaces:
      - istio-system
    allowed_egress_domains:
      - "api.plaid.com"
  service_mesh_config:
    enabled: true
    mesh_type: istio
    revision_tag: "prod-stable"
  pod_security_standard: restricted
```

**What this creates:**

- Strict network isolation (only Istio ingress allowed)
- Minimal egress (only to Plaid API)
- Istio service mesh for mTLS and observability
- Restricted pod security (prevents privilege escalation)
- Large resource quota for production workloads

## Example 9: Multi-Environment Setup

Create namespaces for dev, staging, and prod with appropriate configurations:

**dev-environment.yaml:**
```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesNamespace
metadata:
  name: myapp-dev
spec:
  target_cluster:
    cluster_kind: gcp_gke_cluster
    cluster_name: dev-cluster
  name: myapp-dev
  labels:
    app: myapp
    environment: dev
    team: product-team
  resource_profile:
    preset: small
  pod_security_standard: baseline
```

**staging-environment.yaml:**
```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesNamespace
metadata:
  name: myapp-staging
spec:
  target_cluster:
    cluster_kind: gcp_gke_cluster
    cluster_name: staging-cluster
  name: myapp-staging
  labels:
    app: myapp
    environment: staging
    team: product-team
  resource_profile:
    preset: medium
  service_mesh_config:
    enabled: true
    mesh_type: istio
    revision_tag: "prod-stable"
  pod_security_standard: baseline
```

**prod-environment.yaml:**
```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesNamespace
metadata:
  name: myapp-prod
spec:
  target_cluster:
    cluster_kind: gcp_gke_cluster
    cluster_name: prod-cluster
  name: myapp-prod
  labels:
    app: myapp
    environment: prod
    team: product-team
    criticality: high
  resource_profile:
    preset: large
  network_config:
    isolate_ingress: true
    restrict_egress: true
    allowed_ingress_namespaces:
      - istio-system
  service_mesh_config:
    enabled: true
    mesh_type: istio
    revision_tag: "prod-stable"
  pod_security_standard: restricted
```

## Verification

After deployment, verify the namespace configuration:

```bash
# Check namespace exists
kubectl get namespace <namespace-name>

# Check resource quotas
kubectl get resourcequota -n <namespace-name>

# Check limit ranges
kubectl get limitrange -n <namespace-name>

# Check network policies
kubectl get networkpolicy -n <namespace-name>

# Check labels and annotations
kubectl get namespace <namespace-name> -o yaml

# Check stack outputs (Pulumi)
openmcf pulumi stack output --manifest <file> --stack <stack>
```

## Common Patterns

### Pattern 1: Team-Based Namespaces

Create a namespace per team with appropriate quotas:

```yaml
metadata:
  name: team-{team-name}
spec:
  labels:
    team: {team-name}
    cost-center: engineering
  resource_profile:
    preset: medium
```

### Pattern 2: Environment-Based Namespaces

Create namespaces per environment (dev/staging/prod):

```yaml
metadata:
  name: {app-name}-{environment}
spec:
  labels:
    app: {app-name}
    environment: {environment}
```

### Pattern 3: Microservices Namespace

Pre-configure namespace for microservices with mesh:

```yaml
spec:
  service_mesh_config:
    enabled: true
    mesh_type: istio
  network_config:
    isolate_ingress: true
  pod_security_standard: baseline
```

## Troubleshooting

**Issue: Namespace stuck in Terminating state**

```bash
# Check for finalizers
kubectl get namespace <name> -o yaml | grep finalizers

# Force removal (use with caution)
kubectl patch namespace <name> -p '{"metadata":{"finalizers":[]}}' --type=merge
```

**Issue: Pods can't be scheduled (quota exceeded)**

```bash
# Check quota usage
kubectl describe resourcequota -n <namespace>

# Increase quota by updating the manifest with higher preset or custom values
```

**Issue: Network policies blocking expected traffic**

```bash
# Check network policies
kubectl get networkpolicy -n <namespace>
kubectl describe networkpolicy -n <namespace>

# Add allowed namespaces/CIDRs to network_config in manifest
```

## Next Steps

1. Review the [README.md](README.md) for detailed component documentation
2. Check the [research documentation](docs/README.md) for architecture deep-dive
3. Deploy your first namespace using one of the examples above
4. Customize based on your organization's requirements
