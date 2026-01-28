# Kubernetes Istio - Example Configurations

This document provides comprehensive examples for deploying the Istio service mesh using the KubernetesIstio API resource.

## Table of Contents

1. [Minimal Development Setup](#example-1-minimal-development-setup)
2. [Standard Production Deployment](#example-2-standard-production-deployment)
3. [High-Availability Production](#example-3-high-availability-production)
4. [Resource-Constrained Environment](#example-4-resource-constrained-environment)
5. [Large-Scale Production Cluster](#example-5-large-scale-production-cluster)

---

## Example 1: Minimal Development Setup

The simplest deployment for development and testing environments with minimal resource consumption.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesIstio
metadata:
  name: dev-istio
  labels:
    environment: development
spec:
  target_cluster:
    cluster_name: dev-gke-cluster
  namespace:
    value: istio-system
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 25m
        memory: 64Mi
      limits:
        cpu: 500m
        memory: 256Mi
```

**Use Case:** Local development clusters, testing environments, proof-of-concept deployments.

**Resource Footprint:**
- Minimal CPU and memory for small clusters
- Suitable for local Kind/k3s clusters
- Handles 5-10 services in the mesh

**Expected Behavior:**
- Istio control plane starts with minimal resources
- Suitable for basic traffic management testing
- May not handle high traffic loads

**Deployment:**

Using OpenMCF CLI:
```bash
openmcf apply -f dev-istio.yaml
```

**Verification:**

Check Istio control plane status:
```bash
kubectl get pods -n istio-system
kubectl get pods -n istio-ingress
```

Expected output:
```
NAME                      READY   STATUS    RESTARTS   AGE
istiod-xxxxx             1/1     Running   0          2m
```

---

## Example 2: Standard Production Deployment

Recommended configuration for production clusters with moderate traffic.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesIstio
metadata:
  name: prod-istio
  labels:
    environment: production
    cluster: primary
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: istio-system
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 500m
        memory: 512Mi
      limits:
        cpu: 2000m
        memory: 2Gi
```

**Use Case:** Production Kubernetes clusters with 20-100 services, moderate traffic volume.

**Resource Footprint:**
- Balanced resource allocation
- Handles typical production workloads
- Supports 20-100 services in the mesh
- Processes 1000-5000 requests per second

**Expected Behavior:**
- Stable control plane under normal load
- Quick configuration propagation to data plane
- Suitable for most production use cases

**Deployment:**

```bash
openmcf apply -f prod-istio.yaml
```

**Post-Deployment Configuration:**

Enable automatic sidecar injection for your application namespace:
```bash
kubectl label namespace my-app istio-injection=enabled
```

Deploy a sample application to verify:
```bash
kubectl create deployment echo --image=gcr.io/google_containers/echoserver:1.10 -n my-app
kubectl expose deployment echo --port=8080 -n my-app
```

Verify sidecar injection:
```bash
kubectl get pods -n my-app
# Should show 2/2 containers (app + istio-proxy)
```

---

## Example 2.5: Using Existing Namespaces

When namespaces are managed externally (e.g., via KubernetesNamespace resource or pre-provisioned).

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesIstio
metadata:
  name: prod-istio-existing-ns
  labels:
    environment: production
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: istio-system
  create_namespace: false  # Use existing namespaces
  container:
    resources:
      requests:
        cpu: 500m
        memory: 512Mi
      limits:
        cpu: 2000m
        memory: 2Gi
```

**Use Case:** Centralized namespace management, multi-tenant clusters with namespace-level policies, GitOps workflows where namespaces are managed separately.

**Resource Footprint:**
- Same as standard production deployment
- Assumes namespaces are pre-configured with appropriate labels and RBAC

**Prerequisites:**
- Both `istio-system` and `istio-ingress` namespaces must already exist in the cluster
- Namespaces should have appropriate labels and RBAC configured
- Sufficient resource quotas allocated to both namespaces

**Deployment:**

Before deploying Istio, ensure namespaces exist:
```bash
kubectl get namespace istio-system istio-ingress
```

If namespaces don't exist, create them:
```bash
kubectl create namespace istio-system
kubectl create namespace istio-ingress
```

Then deploy Istio:
```bash
openmcf apply -f prod-istio-existing-ns.yaml
```

**Verification:**

Confirm Istio is using the existing namespaces:
```bash
kubectl get pods -n istio-system
kubectl get pods -n istio-ingress
```

---

## Example 3: High-Availability Production

Configuration for high-traffic production environments requiring maximum reliability.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesIstio
metadata:
  name: ha-prod-istio
  labels:
    environment: production
    tier: critical
    availability: high
spec:
  target_cluster:
    cluster_name: prod-ha-gke-cluster
  namespace:
    value: istio-system
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 1000m
        memory: 1Gi
      limits:
        cpu: 4000m
        memory: 8Gi
```

**Use Case:** Large production clusters with 100+ services, high traffic volume, strict SLAs.

**Resource Footprint:**
- High resource allocation for control plane
- Handles 100-500 services in the mesh
- Processes 10,000+ requests per second
- Low latency configuration propagation

**Capacity Planning:**
- Suitable for clusters with 100+ nodes
- Handles complex traffic routing rules
- Supports advanced features (external authorization, WebAssembly filters)
- Enables detailed telemetry without performance degradation

**Deployment:**

```bash
openmcf apply -f ha-prod-istio.yaml
```

**Advanced Configuration:**

After deployment, configure high-availability features:

Create PeerAuthentication for strict mTLS:
```yaml
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: default
  namespace: istio-system
spec:
  mtls:
    mode: STRICT
```

Apply the policy:
```bash
kubectl apply -f peer-authentication.yaml
```

Monitor control plane performance:
```bash
kubectl top pods -n istio-system
```

---

## Example 4: Resource-Constrained Environment

Optimized for edge clusters or cost-sensitive environments.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesIstio
metadata:
  name: edge-istio
  labels:
    environment: edge
    tier: lightweight
spec:
  target_cluster:
    cluster_name: edge-k3s-cluster
  namespace:
    value: istio-system
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 10m
        memory: 32Mi
      limits:
        cpu: 250m
        memory: 128Mi
```

**Use Case:** Edge computing, IoT gateways, cost-optimized environments, ARM-based clusters.

**Resource Footprint:**
- Absolute minimum resource requirements
- Suitable for single-node clusters
- Handles 2-5 services
- Limited traffic volume (< 100 RPS)

**Trade-offs:**
- May experience higher latency during configuration changes
- Not suitable for production traffic
- Limited observability features
- Reduced telemetry collection

**Deployment:**

```bash
openmcf apply -f edge-istio.yaml
```

**Optimization Tips:**

Disable unnecessary features to save resources:
```yaml
# Create minimal Gateway configuration
apiVersion: networking.istio.io/v1beta1
kind: Gateway
metadata:
  name: minimal-gateway
  namespace: istio-system
spec:
  selector:
    istio: gateway
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "*"
```

---

## Example 5: Large-Scale Production Cluster

Enterprise-grade configuration for massive production deployments.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesIstio
metadata:
  name: enterprise-istio
  labels:
    environment: production
    scale: enterprise
    region: us-east-1
spec:
  target_cluster:
    cluster_name: enterprise-gke-cluster
  namespace:
    value: istio-system
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 2000m
        memory: 4Gi
      limits:
        cpu: 8000m
        memory: 16Gi
```

**Use Case:** Enterprise deployments with 500+ services, multi-tenant clusters, very high traffic.

**Resource Footprint:**
- Maximum resource allocation
- Handles 500+ services in the mesh
- Processes 50,000+ requests per second
- Supports complex routing and security policies
- Full telemetry and observability enabled

**Cluster Requirements:**
- Minimum cluster size: 50+ nodes
- Dedicated node pool for Istio control plane (recommended)
- High-performance storage for control plane persistence
- Multi-zone deployment for HA

**Deployment:**

```bash
openmcf apply -f enterprise-istio.yaml
```

**Enterprise Features Setup:**

Configure namespace-level traffic policies:
```yaml
apiVersion: networking.istio.io/v1beta1
kind: Sidecar
metadata:
  name: default
  namespace: production
spec:
  egress:
  - hosts:
    - "./*"
    - "istio-system/*"
```

Set up authorization policies:
```yaml
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: frontend-policy
  namespace: production
spec:
  selector:
    matchLabels:
      app: frontend
  rules:
  - from:
    - source:
        principals: ["cluster.local/ns/production/sa/api-gateway"]
    to:
    - operation:
        methods: ["GET", "POST"]
```

Apply the policies:
```bash
kubectl apply -f sidecar-config.yaml
kubectl apply -f authz-policy.yaml
```

Monitor at scale:
```bash
# Check control plane metrics
kubectl top pods -n istio-system

# View proxy configuration sync status
istioctl proxy-status
```

---

## Common Post-Deployment Tasks

### Enable Automatic Sidecar Injection

Label namespaces for automatic injection:
```bash
kubectl label namespace production istio-injection=enabled
kubectl label namespace staging istio-injection=enabled
```

### Verify Installation

Check all Istio components:
```bash
kubectl get pods -n istio-system
kubectl get pods -n istio-ingress
kubectl get svc -n istio-system
kubectl get svc -n istio-ingress
```

Verify with istioctl (if installed):
```bash
istioctl version
istioctl verify-install
```

### Access Control Plane

Port-forward to istiod for debugging:
```bash
kubectl port-forward -n istio-system svc/istiod 15014:15014
# Access debug interface at http://localhost:15014/debug
```

### Deploy Sample Application

Deploy the Istio BookInfo sample:
```bash
kubectl create namespace bookinfo
kubectl label namespace bookinfo istio-injection=enabled
kubectl apply -n bookinfo -f https://raw.githubusercontent.com/istio/istio/release-1.22/samples/bookinfo/platform/kube/bookinfo.yaml
```

Create Gateway and VirtualService:
```bash
kubectl apply -n bookinfo -f https://raw.githubusercontent.com/istio/istio/release-1.22/samples/bookinfo/networking/bookinfo-gateway.yaml
```

---

## Troubleshooting

### Check Control Plane Logs

View istiod logs:
```bash
kubectl logs -n istio-system -l app=istiod --tail=100
```

### Verify Gateway Status

Check ingress gateway:
```bash
kubectl get pods -n istio-ingress
kubectl logs -n istio-ingress -l app=istio-gateway
```

### Debug Configuration Issues

Use istioctl to analyze configuration:
```bash
istioctl analyze --all-namespaces
```

### Resource Issues

If control plane is under-resourced:
```bash
# Check resource usage
kubectl top pods -n istio-system

# Describe pod for events
kubectl describe pod -n istio-system <istiod-pod-name>
```

Update the KubernetesIstio resource with higher limits and reapply.

---

## Upgrading Istio

To upgrade to a newer version:

1. Update the resource with new chart version (when available)
2. Apply the updated configuration
3. Restart workloads to get new sidecar versions

```bash
# After updating the YAML with new version
openmcf apply -f updated-istio.yaml

# Restart workloads to update sidecars
kubectl rollout restart deployment -n production
```

---

## Best Practices

1. **Start Conservative**: Begin with standard production resources and scale up based on metrics
2. **Monitor Resource Usage**: Use `kubectl top` and cluster monitoring to track control plane resource consumption
3. **Enable Automatic Sidecar Injection**: Label namespaces rather than manually injecting sidecars
4. **Test Upgrades**: Always test Istio upgrades in non-production environments first
5. **Use PeerAuthentication**: Enable strict mTLS for production environments
6. **Implement Gradual Rollouts**: Use VirtualServices for canary deployments and gradual traffic shifting
7. **Monitor Control Plane**: Set up alerts for istiod pod restarts and resource exhaustion
8. **Segregate Gateways**: Use separate ingress and egress gateways for better security and scaling

---

## Additional Resources

- [Official Istio Documentation](https://istio.io/latest/docs/)
- [Istio Performance Tuning](https://istio.io/latest/docs/ops/deployment/performance-and-scalability/)
- [Istio Security Best Practices](https://istio.io/latest/docs/ops/best-practices/security/)
- [OpenMCF Istio Research](docs/README.md)

