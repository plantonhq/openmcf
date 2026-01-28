# MongoDB Kubernetes API - Example Configurations

## Example w/ Basic Configuration

### Create Using CLI

Create a YAML file using the example shown below. After the YAML is created, use the following command to apply it:

```shell
planton apply -f <yaml-path>
```

### Basic Example

This basic example demonstrates a minimal configuration for deploying a MongoDB instance on Kubernetes using the Percona Server for MongoDB Operator with default settings.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesMongodb
metadata:
  name: basic-mongodb
spec:
  namespace: mongodb-namespace
  createNamespace: true
  container:
    replicas: 1
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 1000m
        memory: 1Gi
    persistenceEnabled: false
```

---

## Example w/ Persistence Enabled

In this example, MongoDB persistence is enabled, and a persistent volume is created for each MongoDB pod to ensure data durability. The `diskSize` field defines the storage size allocated to the MongoDB pods.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesMongodb
metadata:
  name: persistent-mongodb
spec:
  namespace: mongodb-namespace
  createNamespace: true
  container:
    replicas: 3
    persistenceEnabled: true
    diskSize: 10Gi
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 1000m
        memory: 1Gi
```

---

## Example w/ Custom Helm Values

This example demonstrates how to customize the MongoDB deployment using Helm chart values. These values allow for advanced configuration options available in the Bitnami MongoDB Helm chart.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesMongodb
metadata:
  name: custom-mongodb
spec:
  namespace: mongodb-namespace
  createNamespace: true
  container:
    replicas: 2
    persistenceEnabled: true
    diskSize: 20Gi
    resources:
      requests:
        cpu: 200m
        memory: 512Mi
      limits:
        cpu: 2000m
        memory: 2Gi
  helmValues:
    mongodbUsername: myuser
    mongodbDatabase: mydatabase
```

---

## Example w/ Ingress Enabled

In this example, ingress is enabled to allow external access to the MongoDB service. A LoadBalancer service is created with external-dns annotations for automatic DNS configuration.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesMongodb
metadata:
  name: ingress-mongodb
spec:
  namespace: mongodb-namespace
  createNamespace: true
  container:
    replicas: 1
    persistenceEnabled: true
    diskSize: 5Gi
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 1000m
        memory: 1Gi
  ingress:
    enabled: true
    hostname: mongodb.example.com
```

---

## Example w/ Production Configuration

This example demonstrates a production-ready configuration with multiple replicas, persistence enabled, and appropriate resource allocations for a production MongoDB deployment.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesMongodb
metadata:
  name: production-mongodb
spec:
  namespace: mongodb-prod
  createNamespace: true
  container:
    replicas: 3
    persistenceEnabled: true
    diskSize: 50Gi
    resources:
      requests:
        cpu: 500m
        memory: 1Gi
      limits:
        cpu: 2000m
        memory: 4Gi
  ingress:
    enabled: true
    hostname: mongodb.prod.example.com
```

---

## Example w/ Existing Namespace

This example demonstrates using an existing namespace instead of creating a new one. This is useful when you have pre-configured policies, quotas, or RBAC rules in the namespace.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesMongodb
metadata:
  name: existing-ns-mongodb
spec:
  namespace: existing-mongodb-namespace
  createNamespace: false
  container:
    replicas: 3
    persistenceEnabled: true
    diskSize: 20Gi
    resources:
      requests:
        cpu: 200m
        memory: 512Mi
      limits:
        cpu: 1000m
        memory: 2Gi
```
