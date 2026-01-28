# KubernetesGatewayApiCrds Examples

Complete YAML manifests for installing Kubernetes Gateway API CRDs on different cluster types.

---

## Example 1: Standard CRDs on GKE

Install stable Gateway API CRDs on a Google Kubernetes Engine cluster.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGatewayApiCrds
metadata:
  name: gateway-api-crds-gke
spec:
  target_cluster:
    cluster_kind: GcpGkeCluster
    cluster_name: prod-gke-cluster
  version: v1.2.1
  install_channel:
    channel: standard
```

**What gets installed:**
- `gatewayclasses.gateway.networking.k8s.io`
- `gateways.gateway.networking.k8s.io`
- `httproutes.gateway.networking.k8s.io`
- `referencegrants.gateway.networking.k8s.io`

**Use case:** Production clusters using Istio, GKE Gateway Controller, or other standard-conformant implementations.

---

## Example 2: Standard CRDs on EKS

Install Gateway API CRDs on Amazon Elastic Kubernetes Service.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGatewayApiCrds
metadata:
  name: gateway-api-crds-eks
spec:
  target_cluster:
    cluster_kind: AwsEksCluster
    cluster_name: prod-eks-cluster
  version: v1.2.1
  install_channel:
    channel: standard
```

**Use case:** EKS clusters using AWS Gateway API Controller or Envoy Gateway.

---

## Example 3: Standard CRDs on AKS

Install Gateway API CRDs on Azure Kubernetes Service.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGatewayApiCrds
metadata:
  name: gateway-api-crds-aks
spec:
  target_cluster:
    cluster_kind: AzureAksCluster
    cluster_name: prod-aks-cluster
  version: v1.2.1
  install_channel:
    channel: standard
```

**Use case:** AKS clusters with Application Gateway for Containers or other Gateway API implementations.

---

## Example 4: Experimental CRDs for Advanced Routing

Install experimental CRDs to enable TCP, UDP, TLS, and gRPC routing.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGatewayApiCrds
metadata:
  name: gateway-api-crds-experimental
spec:
  target_cluster:
    cluster_kind: GcpGkeCluster
    cluster_name: dev-gke-cluster
  version: v1.2.1
  install_channel:
    channel: experimental
```

**Additional CRDs installed:**
- `tcproutes.gateway.networking.k8s.io` - TCP traffic routing
- `udproutes.gateway.networking.k8s.io` - UDP traffic routing
- `tlsroutes.gateway.networking.k8s.io` - TLS passthrough routing
- `grpcroutes.gateway.networking.k8s.io` - Native gRPC routing

**Use case:** Development environments testing advanced routing features, or production clusters requiring TCP/UDP routing.

---

## Example 5: Latest Version

Install the latest Gateway API version.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGatewayApiCrds
metadata:
  name: gateway-api-crds-latest
spec:
  target_cluster:
    cluster_kind: GcpGkeCluster
    cluster_name: staging-cluster
  version: v1.3.0
  install_channel:
    channel: standard
```

**Use case:** Staging environments testing the latest Gateway API features before production rollout.

---

## Example 6: Minimal Configuration (Defaults)

Let the module use default version and channel.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGatewayApiCrds
metadata:
  name: gateway-api-crds-minimal
spec:
  target_cluster:
    cluster_kind: GcpGkeCluster
    cluster_name: dev-cluster
```

**Defaults applied:**
- `version`: v1.2.1
- `install_channel.channel`: standard (when unspecified, treated as standard)

**Use case:** Quick setup for development or testing environments.

---

## Example 7: DigitalOcean Kubernetes

Install Gateway API CRDs on DigitalOcean Kubernetes.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGatewayApiCrds
metadata:
  name: gateway-api-crds-doks
spec:
  target_cluster:
    cluster_kind: DigitalOceanKubernetesCluster
    cluster_name: prod-doks-cluster
  version: v1.2.1
  install_channel:
    channel: standard
```

**Use case:** DigitalOcean clusters using Traefik, NGINX Gateway Fabric, or other Gateway API implementations.

---

## Example 8: Civo Kubernetes

Install Gateway API CRDs on Civo Kubernetes.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGatewayApiCrds
metadata:
  name: gateway-api-crds-civo
spec:
  target_cluster:
    cluster_kind: CivoKubernetesCluster
    cluster_name: prod-civo-cluster
  version: v1.2.1
  install_channel:
    channel: standard
```

**Use case:** Civo Kubernetes clusters with any Gateway API implementation.

---

## Using Gateway API After Installation

Once CRDs are installed, you can create Gateway API resources. Here's a typical flow:

### Step 1: Create a GatewayClass

The GatewayClass is typically created by the Gateway API implementation (Istio, Envoy Gateway, etc.). Example for reference:

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: GatewayClass
metadata:
  name: istio
spec:
  controllerName: istio.io/gateway-controller
```

### Step 2: Create a Gateway

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: my-gateway
  namespace: default
spec:
  gatewayClassName: istio
  listeners:
  - name: http
    port: 80
    protocol: HTTP
    hostname: "*.example.com"
```

### Step 3: Create an HTTPRoute

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: my-route
  namespace: default
spec:
  parentRefs:
  - name: my-gateway
  hostnames:
  - "app.example.com"
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /
    backendRefs:
    - name: my-service
      port: 8080
```

---

## Verification

After deployment, verify the CRDs are installed:

```bash
# List Gateway API CRDs
kubectl get crds | grep gateway.networking.k8s.io

# Expected output for standard channel:
# gatewayclasses.gateway.networking.k8s.io
# gateways.gateway.networking.k8s.io
# httproutes.gateway.networking.k8s.io
# referencegrants.gateway.networking.k8s.io

# For experimental channel, also expect:
# tcproutes.gateway.networking.k8s.io
# udproutes.gateway.networking.k8s.io
# tlsroutes.gateway.networking.k8s.io
# grpcroutes.gateway.networking.k8s.io
```

---

## Version Compatibility

| Gateway API Version | Kubernetes Minimum | Recommended For |
|---------------------|-------------------|-----------------|
| v1.2.1 | 1.25+ | Production |
| v1.3.0 | 1.26+ | Staging/Testing |
| v1.1.0 | 1.24+ | Legacy clusters |

**Note:** Check your Gateway API implementation's documentation for specific version requirements.

---

## Common Gateway API Implementations

After installing CRDs, you'll need a Gateway API implementation:

| Implementation | Best For |
|----------------|----------|
| [Istio](https://istio.io/latest/docs/tasks/traffic-management/ingress/gateway-api/) | Service mesh users, advanced features |
| [Envoy Gateway](https://gateway.envoyproxy.io/) | Standalone Envoy, cloud-native |
| [NGINX Gateway Fabric](https://github.com/nginxinc/nginx-gateway-fabric) | NGINX users |
| [Traefik](https://doc.traefik.io/traefik/routing/providers/kubernetes-gateway/) | Simple setup, automatic discovery |
| [GKE Gateway Controller](https://cloud.google.com/kubernetes-engine/docs/concepts/gateway-api) | GKE native, Cloud Load Balancer |
| [AWS Gateway API Controller](https://www.gateway-api-controller.eks.aws.dev/) | EKS native, AWS ALB/NLB |

---

## Next Steps

1. Choose the example matching your cluster type
2. Deploy the KubernetesGatewayApiCrds manifest
3. Verify CRDs are installed
4. Install a Gateway API implementation
5. Create Gateway and Route resources

For more details, see:
- [README](README.md) - Overview and feature details
- [Research Documentation](docs/README.md) - Deep dive into Gateway API landscape
- [Gateway API Official Docs](https://gateway-api.sigs.k8s.io/)
