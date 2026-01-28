# Example 1: Basic Keycloak Kubernetes Setup

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KeycloakKubernetes
metadata:
  name: keycloak-instance-basic
spec:
  target_cluster:
    cluster_name: my-cluster
  namespace:
    value: keycloak
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 50m
        memory: 256Mi
      limits:
        cpu: 1
        memory: 1Gi
  ingress:
    enabled: false
```

# Example 2: Keycloak Kubernetes with Ingress Enabled

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KeycloakKubernetes
metadata:
  name: keycloak-instance-ingress
spec:
  target_cluster:
    cluster_name: my-cluster
  namespace:
    value: keycloak-ingress
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 100m
        memory: 512Mi
      limits:
        cpu: 2
        memory: 2Gi
  ingress:
    enabled: true
    hostname: keycloak.mydomain.com
```

# Example 3: Keycloak Kubernetes with Minimal Configuration

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KeycloakKubernetes
metadata:
  name: keycloak-minimal
spec:
  target_cluster:
    cluster_name: my-cluster
  namespace:
    value: keycloak-minimal
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 50m
        memory: 256Mi
      limits:
        cpu: 1
        memory: 1Gi
  ingress:
    enabled: false
```

# Example 4: Keycloak Kubernetes with High Resource Allocation

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KeycloakKubernetes
metadata:
  name: keycloak-high-resources
spec:
  target_cluster:
    cluster_name: prod-cluster
  namespace:
    value: keycloak-prod
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 500m
        memory: 2Gi
      limits:
        cpu: 4
        memory: 8Gi
  ingress:
    enabled: true
    hostname: keycloak-large.mydomain.com
```

# Example 5: Keycloak Kubernetes with Port Forwarding for Local Access

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KeycloakKubernetes
metadata:
  name: keycloak-port-forward
spec:
  target_cluster:
    cluster_name: my-cluster
  namespace:
    value: keycloak-dev
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 100m
        memory: 512Mi
      limits:
        cpu: 2
        memory: 2Gi
  ingress:
    enabled: false
```

# Example 6: Keycloak with Namespace Creation

```yaml
apiVersion: kubernetes.openmcf.cloud/v1
kind: KubernetesKeycloak
metadata:
  name: keycloak-with-namespace
spec:
  target_cluster:
    cluster_name: "prod-gke-cluster"
  namespace:
    value: "keycloak-dedicated"
  create_namespace: true  # Module creates the namespace
  container:
    resources:
      requests:
        cpu: "100m"
        memory: "256Mi"
      limits:
        cpu: "1000m"
        memory: "1Gi"
  ingress:
    enabled: true
    hostname: "auth.example.com"
```

# Example 7: Keycloak in Pre-existing Namespace

```yaml
apiVersion: kubernetes.openmcf.cloud/v1
kind: KubernetesKeycloak
metadata:
  name: keycloak-shared
spec:
  target_cluster:
    cluster_name: "prod-gke-cluster"
  namespace:
    value: "shared-auth"
  create_namespace: false  # Assumes namespace already exists
  container:
    resources:
      requests:
        cpu: "100m"
        memory: "256Mi"
      limits:
        cpu: "1000m"
        memory: "1Gi"
  ingress:
    enabled: false
```
