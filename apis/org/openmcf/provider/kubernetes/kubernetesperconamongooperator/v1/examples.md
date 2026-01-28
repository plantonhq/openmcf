# KubernetesPerconaMongoOperator API Resource Examples

Below are examples demonstrating how to configure and deploy the `KubernetesPerconaMongoOperator` API resource using various specifications. Follow the instructions to create and apply each YAML configuration using the Planton CLI.

---

## Example 1: Basic Operator Deployment

### Description

This example demonstrates a basic deployment of the Percona Operator for MongoDB within a Kubernetes cluster. It uses the default resource allocations suitable for most production environments.

### Create and Apply

1. **Create a YAML file** using the example below.
2. **Apply the configuration** using the following command:

    ```shell
    planton pulumi up --manifest <yaml-path>
    ```

### YAML Configuration

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesPerconaMongoOperator
metadata:
  name: percona-operator-prod
spec:
  targetCluster:
    clusterName: "my-gke-cluster"
  namespace:
    value: "percona-operator"
  createNamespace: true
  container:
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 1000m
        memory: 1Gi
```

---

## Example 2: Operator Deployment with Custom Resources

### Description

This example illustrates how to deploy the Percona Operator for MongoDB with customized resource allocations. This configuration provides more resources to the operator, which may be beneficial in large clusters managing many MongoDB instances.

### Create and Apply

1. **Create a YAML file** using the example below.
2. **Apply the configuration** using the following command:

    ```shell
    planton pulumi up --manifest <yaml-path>
    ```

### YAML Configuration

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesPerconaMongoOperator
metadata:
  name: percona-operator-large
spec:
  targetCluster:
    clusterName: "production-gke-cluster"
  namespace:
    value: "percona-operator"
  createNamespace: true
  container:
    resources:
      requests:
        cpu: 200m
        memory: 512Mi
      limits:
        cpu: 2000m
        memory: 2Gi
```

---

## Example 3: Development Environment Deployment

### Description

This example shows a minimal resource configuration suitable for development or testing environments where resource constraints are tighter.

### Create and Apply

1. **Create a YAML file** using the example below.
2. **Apply the configuration** using the following command:

    ```shell
    planton pulumi up --manifest <yaml-path>
    ```

### YAML Configuration

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesPerconaMongoOperator
metadata:
  name: percona-operator-dev
spec:
  targetCluster:
    clusterName: "my-local-k8s-cluster"
  namespace:
    value: "percona-operator-dev"
  createNamespace: true
  container:
    resources:
      requests:
        cpu: 50m
        memory: 128Mi
      limits:
        cpu: 500m
        memory: 512Mi
```

---

## Example 4: Using Pre-existing Namespace

### Description

This example demonstrates deploying the operator into an existing namespace. This is useful when the namespace is managed separately (e.g., via GitOps) or when you have pre-configured namespace settings.

### Prerequisites

Ensure the namespace exists before applying:

```shell
kubectl create namespace shared-operators
```

### Create and Apply

1. **Create a YAML file** using the example below.
2. **Apply the configuration** using the following command:

    ```shell
    planton pulumi up --manifest <yaml-path>
    ```

### YAML Configuration

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesPerconaMongoOperator
metadata:
  name: percona-operator-existing-ns
spec:
  targetCluster:
    clusterName: "my-gke-cluster"
  namespace:
    value: "shared-operators"
  createNamespace: false
  container:
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 1000m
        memory: 1Gi
```

**Note**: The namespace `shared-operators` must exist before deployment. If it doesn't exist, the deployment will fail.

---

## Verifying the Deployment

After deploying the operator, you can verify its status using standard Kubernetes commands:

```shell
# Check operator pod status
kubectl get pods -n percona-operator

# Check operator logs
kubectl logs -n percona-operator -l app.kubernetes.io/name=kubernetes-percona-mongo-operator

# Verify CRDs are installed
kubectl get crds | grep percona

# Expected CRDs:
# - perconaservermongodbs.psmdb.percona.com
# - perconaservermongodbbackups.psmdb.percona.com
# - perconaservermongodbrestores.psmdb.percona.com
```

## Next Steps

Once the operator is deployed, you can:

1. Deploy MongoDB clusters using the `MongoDBKubernetes` workload resource
2. Create custom PerconaServerMongoDB resources directly
3. Configure automated backups and monitoring
4. Scale your MongoDB deployments as needed

