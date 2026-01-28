Here are a few examples for the `KubernetesArgocd` API resource. These examples demonstrate how to configure and deploy Argo CD on a Kubernetes cluster using OpenMCF's unified API structure.

---

# Example 1: Basic Argo CD Deployment with Namespace Creation

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesArgocd
metadata:
  name: my-argocd
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace:
    value: argocd
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 50m
        memory: 100Mi
      limits:
        cpu: 1000m
        memory: 1Gi
```

**Description:** A basic Argo CD deployment with default resource allocations. The namespace will be automatically created by the module.

---

# Example 2: Argo CD with Ingress Enabled

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesArgocd
metadata:
  name: argocd-prod
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: argocd-prod
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 100m
        memory: 512Mi
      limits:
        cpu: 2000m
        memory: 2Gi
  ingress:
    enabled: true
    hostname: argocd.example.com
```

**Description:** Production Argo CD deployment with ingress enabled for external access at `argocd.example.com`.

---

# Example 3: Argo CD Deployment with Custom Resources

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesArgocd
metadata:
  name: argocd-large
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace:
    value: argocd-large
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 200m
        memory: 1Gi
      limits:
        cpu: 3000m
        memory: 4Gi
```

**Description:** Argo CD deployment with higher resource allocations for managing large-scale GitOps workflows with many applications.

---

# Example 4: Minimal Argo CD Deployment

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesArgocd
metadata:
  name: minimal-argocd
spec:
  target_cluster:
    cluster_name: dev-cluster
  namespace:
    value: argocd-dev
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 50m
        memory: 100Mi
      limits:
        cpu: 1000m
        memory: 1Gi
```

**Description:** Minimal Argo CD deployment using default container resources. Suitable for development and testing environments.

---

# Example 5: Argo CD with Pre-existing Namespace

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesArgocd
metadata:
  name: argocd-existing-ns
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace:
    value: platform-tools  # Must exist before deployment
  create_namespace: false  # Use existing namespace
  container:
    resources:
      requests:
        cpu: 100m
        memory: 512Mi
      limits:
        cpu: 2000m
        memory: 2Gi
  ingress:
    enabled: true
    hostname: argocd.platform.example.com
```

**Description:** Argo CD deployment using a pre-existing namespace. This is useful in production environments where namespace creation is managed separately or requires elevated privileges. Ensure the namespace exists before deploying.
