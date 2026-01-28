# RedisKubernetes API-Resource Examples

Below are two examples demonstrating how to configure and deploy the `RedisKubernetes` API resource using various specifications. Follow the instructions to create and apply each YAML configuration using the Planton CLI.

---

## Example 1: Basic Redis Deployment

### Description

This example demonstrates a basic deployment of Redis within a Kubernetes cluster without enabling data persistence. It sets up a single Redis pod with default resource allocations and ingress configurations.

### Create and Apply

1. **Create a YAML file** using the example below.
2. **Apply the configuration** using the following command:

    ```shell
    planton apply -f <yaml-path>
    ```

### YAML Configuration

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: RedisKubernetes
metadata:
  name: basic-redis
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: my-namespace
  create_namespace: true
  container:
    replicas: 1
    resources:
      requests:
        cpu: 50m
        memory: 256Mi
      limits:
        cpu: 1
        memory: 1Gi
    persistenceEnabled: false
```

---

## Example 2: Redis Deployment with Persistence

### Description

This example illustrates how to deploy Redis with data persistence enabled. It configures multiple Redis replicas and attaches persistent storage to ensure data durability across pod restarts. The disk size for the persistent volume is customized to meet specific storage requirements.

### Create and Apply

1. **Create a YAML file** using the example below.
2. **Apply the configuration** using the following command:

    ```shell
    planton apply -f <yaml-path>
    ```

### YAML Configuration

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: RedisKubernetes
metadata:
  name: persistent-redis
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: my-namespace
  create_namespace: true
  container:
    replicas: 3
    resources:
      requests:
        cpu: 100m
        memory: 512Mi
      limits:
        cpu: 2
        memory: 2Gi
    persistenceEnabled: true
    diskSize: 20Gi
  ingress:
    enabled: true
    hostname: redis.example.com
```

---

## Example 3: Redis with Existing Namespace

### Description

This example demonstrates deploying Redis into an existing namespace that has been pre-configured with ResourceQuotas, NetworkPolicies, or RBAC rules. By setting `create_namespace: false`, the module will use the existing namespace without attempting to create it.

### Create and Apply

1. **Ensure the namespace exists**:
    ```shell
    kubectl create namespace production-cache
    ```

2. **Create a YAML file** using the example below.
3. **Apply the configuration** using the following command:

    ```shell
    planton apply -f <yaml-path>
    ```

### YAML Configuration

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: RedisKubernetes
metadata:
  name: production-redis
spec:
  target_cluster:
    cluster_name: "production-gke-cluster"
  namespace:
    value: production-cache
  create_namespace: false
  container:
    replicas: 3
    resources:
      requests:
        cpu: 200m
        memory: 1Gi
      limits:
        cpu: 2
        memory: 4Gi
    persistenceEnabled: true
    diskSize: 50Gi
  ingress:
    enabled: true
    hostname: redis-prod.example.com
```
