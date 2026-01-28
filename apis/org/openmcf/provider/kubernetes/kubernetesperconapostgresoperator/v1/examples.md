# KubernetesPerconaPostgresOperator API Resource Examples

Below are examples demonstrating how to configure and deploy the `KubernetesPerconaPostgresOperator` API resource using various specifications. Follow the instructions to create and apply each YAML configuration using the Planton CLI.

## Namespace Management

The `create_namespace` field controls whether the deployment creates the namespace or uses an existing one:

- **`create_namespace: true`** - The deployment creates and manages the namespace. Recommended for new deployments.
- **`create_namespace: false`** - The deployment uses an existing namespace. Use this when namespaces are managed separately.

**Important**: When using `create_namespace: false`, ensure the namespace exists before deployment.

---

## Example 1: Basic Operator Deployment

### Description

This example demonstrates a basic deployment of the Percona Operator for PostgreSQL within a Kubernetes cluster. It uses the default resource allocations suitable for most production environments.

### Create and Apply

1. **Create a YAML file** using the example below.
2. **Apply the configuration** using the following command:

    ```shell
    planton pulumi up --manifest <yaml-path>
    ```

### YAML Configuration

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesPerconaPostgresOperator
metadata:
  name: percona-pg-operator-prod
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: "kubernetes-percona-postgres-operator"
  create_namespace: true
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

This example illustrates how to deploy the Percona Operator for PostgreSQL with customized resource allocations. This configuration provides more resources to the operator, which may be beneficial in large clusters managing many PostgreSQL instances.

### Create and Apply

1. **Create a YAML file** using the example below.
2. **Apply the configuration** using the following command:

    ```shell
    planton pulumi up --manifest <yaml-path>
    ```

### YAML Configuration

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesPerconaPostgresOperator
metadata:
  name: percona-pg-operator-large
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: "kubernetes-percona-postgres-operator"
  create_namespace: true
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
kind: KubernetesPerconaPostgresOperator
metadata:
  name: percona-pg-operator-dev
spec:
  target_cluster:
    cluster_name: "my-local-gke-cluster"
  namespace:
    value: "kubernetes-percona-postgres-operator-dev"
  create_namespace: true
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

## Example 4: Using Existing Namespace

### Description

This example demonstrates deploying the operator into a pre-existing namespace managed by your platform team or GitOps system.

### Prerequisites

Ensure the namespace exists before deployment:

```shell
kubectl create namespace percona-operators
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
kind: KubernetesPerconaPostgresOperator
metadata:
  name: percona-pg-operator-shared
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: "percona-operators"
  create_namespace: false
  container:
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 1000m
        memory: 1Gi
```

### Notes

- The namespace `percona-operators` must exist before deployment
- This approach is useful when namespace management is centralized
- On deletion, the namespace will not be removed, only the operator resources

---

## Verifying the Deployment

After deploying the operator, you can verify its status using standard Kubernetes commands:

```shell
# Check operator pod status
kubectl get pods -n kubernetes-percona-postgres-operator

# Check operator logs
kubectl logs -n kubernetes-percona-postgres-operator -l app.kubernetes.io/name=kubernetes-percona-postgres-operator

# Verify CRDs are installed
kubectl get crds | grep percona

# Expected CRDs:
# - perconapgclusters.pgv2.percona.com
# - perconapgbackups.pgv2.percona.com
# - perconapgrestores.pgv2.percona.com
```

## Next Steps

Once the operator is deployed, you can:

1. Deploy PostgreSQL clusters using the `PostgreSQLKubernetes` workload resource
2. Create custom PerconaPGCluster resources directly
3. Configure automated backups and monitoring
4. Scale your PostgreSQL deployments as needed

