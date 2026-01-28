# Kubernetes Helm Release - Pulumi Examples

This document provides examples for deploying Helm charts using the KubernetesHelmRelease API resource with Pulumi.

---

## Example 1: Basic Helm Release Deployment

A minimal example deploying NGINX from the Bitnami Helm repository.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesHelmRelease
metadata:
  name: basic-nginx
spec:
  target_cluster:
    cluster_name: my-cluster
  namespace:
    value: default
  create_namespace: true
  repo: https://charts.bitnami.com/bitnami
  name: nginx
  version: 15.14.0
  values: {}
```

**Use Case:** Quick deployment for testing or development environments.

---

## Example 2: Helm Release with Custom Values

Deploy NGINX with custom configuration values.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesHelmRelease
metadata:
  name: nginx-ha
spec:
  target_cluster:
    cluster_name: my-cluster
  namespace:
    value: default
  create_namespace: true
  repo: https://charts.bitnami.com/bitnami
  name: nginx
  version: 15.14.0
  values:
    replicaCount: "3"
    service:
      type: "LoadBalancer"
    resources:
      requests:
        memory: "256Mi"
        cpu: "100m"
      limits:
        memory: "512Mi"
        cpu: "500m"
```

**Use Case:** Production deployment with resource limits and multiple replicas.

---

## Example 3: Using Existing Namespace

Deploy Redis to an existing namespace without creating it.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesHelmRelease
metadata:
  name: redis-existing-ns
spec:
  target_cluster:
    cluster_name: my-cluster
  namespace:
    value: existing-redis-namespace
  create_namespace: false
  repo: https://charts.bitnami.com/bitnami
  name: redis
  version: 18.19.0
  values: {}
```

**Use Case:** Deploy to a namespace that already exists and is managed separately.

---

## Example 4: Helm Release with Multiple Custom Values

Deploy Redis with Sentinel for high availability.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesHelmRelease
metadata:
  name: redis-ha
spec:
  target_cluster:
    cluster_name: my-cluster
  namespace:
    value: redis
  create_namespace: true
  repo: https://charts.bitnami.com/bitnami
  name: redis
  version: 18.19.0
  values:
    architecture: "replication"
    auth:
      enabled: "true"
      password: "super-secret-password"
    sentinel:
      enabled: "true"
      quorum: "2"
    metrics:
      enabled: "true"
    replica:
      replicaCount: "3"
```

**Use Case:** Production Redis cluster with high availability.

---

## Example 5: WordPress with Ingress Configuration

Deploy WordPress with Ingress enabled for external access.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesHelmRelease
metadata:
  name: wordpress-blog
spec:
  target_cluster:
    cluster_name: my-cluster
  namespace:
    value: wordpress
  create_namespace: true
  repo: https://charts.bitnami.com/bitnami
  name: wordpress
  version: 19.0.0
  values:
    wordpressUsername: "admin"
    wordpressEmail: "admin@example.com"
    ingress:
      enabled: "true"
      hostname: "blog.example.com"
      pathType: "Prefix"
      ingressClassName: "nginx"
    service:
      type: "ClusterIP"
    persistence:
      enabled: "true"
      size: "10Gi"
```

**Use Case:** Public-facing blog or website with persistent storage.

---

## Deployment Commands

### Deploy a Helm Release

```bash
planton pulumi up --stack-input <example-file>.yaml
```

### Check Deployment Status

```bash
# View stack outputs
planton pulumi stack output --stack-input <example-file>.yaml

# Check Helm releases
helm list -A

# Check pods in namespace
kubectl get pods -n <namespace>
```

### Delete Deployment

```bash
planton pulumi destroy --stack-input <example-file>.yaml
```

## Notes

- Always specify the `target_cluster` to identify which Kubernetes cluster to deploy to
- The `namespace` field supports both literal values and references to KubernetesNamespace resources
- Set `create_namespace: true` to create the namespace automatically (default behavior)
- Set `create_namespace: false` to use an existing namespace
- All Helm chart values should be provided as string key-value pairs in the `values` map
