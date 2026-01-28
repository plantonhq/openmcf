# Pulumi Examples for KubernetesOpenFga

This document provides example configurations for deploying OpenFGA (Fine-Grained Authorization) on Kubernetes using the Pulumi module.

---

## Example with Ingress Enabled and Plain String Password

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesOpenFga
metadata:
  name: open-fga-service
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: open-fga-service
  createNamespace: true
  container:
    replicas: 3
    resources:
      requests:
        cpu: 50m
        memory: 256Mi
      limits:
        cpu: 1
        memory: 1Gi
  ingress:
    enabled: true
    hostname: openfga.mycluster.example.com
  datastore:
    engine: postgres
    host: db-host
    port: 5432
    database: openfga
    username: user
    password:
      stringValue: password
```

---

## Example with Kubernetes Secret Reference (Production Recommended)

Using a Kubernetes Secret reference is the recommended approach for production deployments as it avoids storing sensitive credentials in plain text.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesOpenFga
metadata:
  name: open-fga-prod
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: open-fga-prod
  createNamespace: true
  container:
    replicas: 3
    resources:
      requests:
        cpu: 500m
        memory: 512Mi
      limits:
        cpu: 2000m
        memory: 2Gi
  ingress:
    enabled: true
    hostname: openfga-prod.example.com
  datastore:
    engine: postgres
    host: prod-db-host.example.com
    port: 5432
    database: openfga
    username: openfga_user
    password:
      secretRef:
        name: openfga-db-credentials
        key: password
    isSecure: true
```

---

## Example with Ingress Disabled and MySQL Datastore

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesOpenFga
metadata:
  name: open-fga-service
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: open-fga-service
  createNamespace: true
  container:
    replicas: 2
    resources:
      requests:
        cpu: 50m
        memory: 256Mi
      limits:
        cpu: 500m
        memory: 512Mi
  ingress:
    enabled: false
  datastore:
    engine: mysql
    host: mysql-db
    port: 3306
    database: openfga
    username: user
    password:
      stringValue: password
```

---

## Example with Minimum Required Fields

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesOpenFga
metadata:
  name: basic-openfga
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: basic-openfga
  createNamespace: true
  container:
    replicas: 1
    resources:
      requests:
        cpu: 50m
        memory: 128Mi
      limits:
        cpu: 200m
        memory: 512Mi
  datastore:
    engine: postgres
    host: db-host
    database: openfga
    username: user
    password:
      stringValue: password
```

---

## Example with High Availability Configuration

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesOpenFga
metadata:
  name: open-fga-high-availability
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: open-fga-high-availability
  createNamespace: true
  container:
    replicas: 5
    resources:
      requests:
        cpu: 500m
        memory: 512Mi
      limits:
        cpu: 2000m
        memory: 4Gi
  ingress:
    enabled: true
    hostname: open-fga-ha.example.com
  datastore:
    engine: postgres
    host: ha-db-host
    database: openfga
    username: user
    password:
      secretRef:
        name: openfga-ha-credentials
        key: password
    isSecure: true
```

---

## Example with Pre-existing Namespace

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesOpenFga
metadata:
  name: openfga-shared-namespace
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: shared-services
  createNamespace: false  # Namespace is managed externally
  container:
    replicas: 2
    resources:
      requests:
        cpu: 50m
        memory: 256Mi
      limits:
        cpu: 1000m
        memory: 1Gi
  datastore:
    engine: postgres
    host: db-host
    database: openfga
    username: user
    password:
      stringValue: password
```

---

## Password Configuration Options

The `password` field in the datastore configuration supports two modes:

### Plain String Value (Development/Testing)

For development and testing environments, you can provide the password directly:

```yaml
password:
  stringValue: my-password
```

### Kubernetes Secret Reference (Production)

For production deployments, reference an existing Kubernetes Secret:

```yaml
password:
  secretRef:
    name: openfga-db-credentials  # Name of the Kubernetes Secret
    key: password                  # Key within the Secret containing the password
    namespace: ""                  # Optional: defaults to deployment namespace
```

**Note:** The Kubernetes Secret must exist before deploying OpenFGA when using `secretRef`.
