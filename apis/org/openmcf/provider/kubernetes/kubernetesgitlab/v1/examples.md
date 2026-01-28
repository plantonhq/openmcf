# Kubernetes GitLab Examples

## Namespace Management

All examples below demonstrate namespace management using the `create_namespace` flag:
- When `create_namespace: true` (default), the module creates the namespace automatically
- When `create_namespace: false`, you must create the namespace beforehand

---

# Example 1: Basic GitLab Deployment

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGitlab
metadata:
  name: gitlab-instance
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace:
    value: gitlab-instance
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 50m
        memory: 256Mi
      limits:
        cpu: 1
        memory: 1Gi
```

---

# Example 2: GitLab with Ingress and Custom Hostname

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGitlab
metadata:
  name: gitlab-production
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: gitlab-production
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 200m
        memory: 512Mi
      limits:
        cpu: 2
        memory: 2Gi
  ingress:
    enabled: true
    hostname: gitlab.example.com
```

---

# Example 3: GitLab Deployment with Custom Resources

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGitlab
metadata:
  name: gitlab-custom
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace:
    value: gitlab-custom
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 500m
        memory: 1Gi
      limits:
        cpu: 3
        memory: 4Gi
```

---

# Example 4: GitLab Using Existing Namespace

This example demonstrates deploying GitLab into a pre-existing namespace.
The namespace must be created before applying this configuration.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGitlab
metadata:
  name: gitlab-shared
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace:
    value: shared-services
  create_namespace: false  # Namespace must already exist
  container:
    resources:
      requests:
        cpu: 250m
        memory: 512Mi
      limits:
        cpu: 1
        memory: 1Gi
```

---

# Example 5: Minimal GitLab Deployment

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesGitlab
metadata:
  name: minimal-gitlab
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace:
    value: minimal-gitlab
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 50m
        memory: 100Mi
      limits:
        cpu: 1
        memory: 1Gi
```
