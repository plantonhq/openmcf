# Multiple Examples for `KubernetesOpenFga` API-Resource

## Example with Minimal Configuration (Plain String Password)

### Create using CLI

Create a YAML file using the example shown below. After the YAML is created, use the command below to apply it.

```shell
planton apply -f <yaml-path>
```

### YAML Configuration

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
    port: 5432
    database: openfga
    username: user
    password:
      stringValue: my-password
```

---

## Example with Kubernetes Secret Reference (Production Recommended)

### Create using CLI

Create a YAML file using the example shown below. After the YAML is created, use the command below to apply it.

```shell
planton apply -f <yaml-path>
```

### YAML Configuration

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesOpenFga
metadata:
  name: openfga-prod
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: openfga-prod
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
    hostname: openfga.example.com
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

## Example with Ingress Enabled

### Create using CLI

Create a YAML file using the example shown below. After the YAML is created, use the command below to apply it.

```shell
planton apply -f <yaml-path>
```

### YAML Configuration

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
        cpu: 1000m
        memory: 1Gi
  ingress:
    enabled: true
    hostname: openfga.example.com
  datastore:
    engine: postgres
    host: db-host
    database: openfga
    username: user
    password:
      stringValue: password
```

---

## Example with MySQL Datastore

### Create using CLI

Create a YAML file using the example shown below. After the YAML is created, use the command below to apply it.

```shell
planton apply -f <yaml-path>
```

### YAML Configuration

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesOpenFga
metadata:
  name: open-fga-mysql
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: open-fga-mysql
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

## Example with High Availability Configuration

### Create using CLI

Create a YAML file using the example shown below. After the YAML is created, use the command below to apply it.

```shell
planton apply -f <yaml-path>
```

### YAML Configuration

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
    hostname: openfga-ha.example.com
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

## Example with Production-Grade Resources

### Create using CLI

Create a YAML file using the example shown below. After the YAML is created, use the command below to apply it.

```shell
planton apply -f <yaml-path>
```

### YAML Configuration

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesOpenFga
metadata:
  name: open-fga-production
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: open-fga-production
  createNamespace: true
  container:
    replicas: 3
    resources:
      requests:
        cpu: 1000m
        memory: 2Gi
      limits:
        cpu: 4000m
        memory: 8Gi
  ingress:
    enabled: true
    hostname: openfga-prod.company.com
  datastore:
    engine: postgres
    host: prod-db-host
    port: 5432
    database: openfga_prod
    username: openfga_user
    password:
      secretRef:
        name: openfga-prod-credentials
        key: password
    isSecure: true
```

---

## Example with Pre-existing Namespace

### Create using CLI

Create a YAML file using the example shown below. After the YAML is created, use the command below to apply it.

```shell
planton apply -f <yaml-path>
```

### YAML Configuration

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
  createNamespace: false  # Namespace already exists and is managed externally
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
