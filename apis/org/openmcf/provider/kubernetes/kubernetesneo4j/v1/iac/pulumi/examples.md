# Neo4j Kubernetes API - Example Configurations

## Namespace Management

The `create_namespace` flag controls whether the Neo4j component creates the namespace or expects it to already exist:

- **`create_namespace: true`** (recommended for new deployments): The component creates and manages the namespace
- **`create_namespace: false`**: The component expects the namespace to already exist (useful when namespace is managed separately or by another component)

---

## Example w/ Basic Configuration

### Create Using CLI

Create a YAML file using the example shown below. After the YAML is created, use the following command to apply it:

```shell
planton apply -f <yaml-path>
```

### Basic Example

This is a basic configuration for deploying a Neo4j instance on Kubernetes. It specifies the necessary resources like CPU and memory, along with the Kubernetes cluster credential.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: Neo4jKubernetes
metadata:
  name: basic-neo4j
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: my-namespace
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 100m
        memory: 512Mi
      limits:
        cpu: 1000m
        memory: 2Gi
```

---

## Example w/ Ingress Enabled

In this example, ingress is enabled to allow external access to the Neo4j instance. This setup is particularly useful if the Neo4j database needs to be accessed from outside the Kubernetes cluster.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: Neo4jKubernetes
metadata:
  name: ingress-neo4j
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: my-namespace
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 200m
        memory: 1Gi
      limits:
        cpu: 1500m
        memory: 3Gi
  ingress:
    enabled: true
    hostname: neo4j.example.com
```

---

## Example w/ Custom Resource Limits

This example demonstrates how to configure custom resource limits for the Neo4j instance, tailoring the deployment to specific performance needs.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: Neo4jKubernetes
metadata:
  name: custom-resources-neo4j
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: my-namespace
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 250m
        memory: 1Gi
      limits:
        cpu: 2000m
        memory: 4Gi
```

---

## Example w/ Minimal Configuration

This minimal example specifies only the mandatory fields for deploying a Neo4j instance on Kubernetes. Although functional, additional configurations like ingress and persistence can be added as needed.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: Neo4jKubernetes
metadata:
  name: minimal-neo4j
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: my-namespace
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 50m
        memory: 256Mi
      limits:
        cpu: 500m
        memory: 1Gi
```
