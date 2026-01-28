Here are a few examples for the `PostgresKubernetes` API resource, following the same format as provided:

---

# Basic Example

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: PostgresKubernetes
metadata:
  name: postgres-db
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: "postgres-db"
  create_namespace: true
  container:
    replicas: 1
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 1
        memory: 1Gi
    diskSize: 10Gi
  ingress:
    enabled: true
    hostname: postgres-db.example.com
```

---

# Example with Environment Variables

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: PostgresKubernetes
metadata:
  name: postgres-db
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: "postgres-db"
  create_namespace: true
  container:
    replicas: 2
    resources:
      requests:
        cpu: 200m
        memory: 512Mi
      limits:
        cpu: 2000m
        memory: 2Gi
    diskSize: 20Gi
  ingress:
    enabled: true
    hostname: postgres-prod.example.com
  env:
    variables:
      DATABASE_USER: admin
      DATABASE_NAME: prod-db
```

---

# Example with Secrets for Sensitive Information

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: PostgresKubernetes
metadata:
  name: secure-postgres-db
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: "secure-postgres-db"
  create_namespace: true
  container:
    replicas: 3
    resources:
      requests:
        cpu: 300m
        memory: 1Gi
      limits:
        cpu: 4000m
        memory: 4Gi
    diskSize: 50Gi
  ingress:
    enabled: false
  env:
    secrets:
      DATABASE_PASSWORD: ${gcpsm-my-org-prod-gcp-secrets.db-password}
    variables:
      DATABASE_NAME: secure-db
```

---

# Example with Minimum Required Fields

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: PostgresKubernetes
metadata:
  name: minimal-postgres-db
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: "minimal-postgres-db"
  create_namespace: true
  container:
    replicas: 1
    resources:
      requests:
        cpu: 50m
        memory: 128Mi
      limits:
        cpu: 500m
        memory: 512Mi
    diskSize: 5Gi
```

---

# Example with Existing Namespace

This example shows how to deploy PostgreSQL into an existing namespace that is managed separately.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: PostgresKubernetes
metadata:
  name: shared-postgres-db
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: "shared-database-namespace"
  create_namespace: false  # Use existing namespace
  container:
    replicas: 2
    resources:
      requests:
        cpu: 200m
        memory: 512Mi
      limits:
        cpu: 1000m
        memory: 2Gi
    diskSize: 20Gi
  ingress:
    enabled: false
```

---

These examples demonstrate varying configurations for deploying PostgreSQL on Kubernetes using the `PostgresKubernetes` API resource. Each example includes different use cases such as handling environment variables, secrets, basic deployment, minimum configuration setups, and namespace management options for flexibility.