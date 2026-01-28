# KubernetesPerconaMysqlOperator API Resource Examples

Below are examples demonstrating how to configure and deploy the `KubernetesPerconaMysqlOperator` API resource using various specifications. Follow the instructions to create and apply each YAML configuration using the Planton CLI.

---

## Example 1: Basic Operator Deployment

### Description

This example demonstrates a basic deployment of the Percona Operator for MySQL within a Kubernetes cluster. It uses the default resource allocations suitable for most production environments.

### Create and Apply

1. **Create a YAML file** using the example below.
2. **Apply the configuration** using the following command:

    ```shell
    planton pulumi up --manifest <yaml-path>
    ```

### YAML Configuration

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesPerconaMysqlOperator
metadata:
  name: percona-mysql-operator-prod
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace:
    value: percona-mysql-operator
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

This example illustrates how to deploy the Percona Operator for MySQL with customized resource allocations. This configuration provides more resources to the operator, which may be beneficial in large clusters managing many MySQL instances.

### Create and Apply

1. **Create a YAML file** using the example below.
2. **Apply the configuration** using the following command:

    ```shell
    planton pulumi up --manifest <yaml-path>
    ```

### YAML Configuration

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesPerconaMysqlOperator
metadata:
  name: percona-mysql-operator-large
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: percona-mysql-operator
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
kind: KubernetesPerconaMysqlOperator
metadata:
  name: percona-mysql-operator-dev
spec:
  target_cluster:
    cluster_name: dev-gke-cluster
  namespace:
    value: percona-mysql-operator-dev
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

This example demonstrates deploying the Percona Operator for MySQL into a pre-existing namespace. This is useful when namespaces are managed separately by platform teams or when namespace-level policies, quotas, or network policies are pre-configured.

### Prerequisites

Before applying this configuration, ensure the namespace exists:

```shell
# Create the namespace if it doesn't exist
kubectl create namespace database-operators

# Verify namespace exists
kubectl get namespace database-operators
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
kind: KubernetesPerconaMysqlOperator
metadata:
  name: percona-mysql-operator-existing-ns
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: database-operators
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

**Important:** With `create_namespace: false`, the namespace must exist before deployment. If the namespace doesn't exist, the deployment will fail.

---

## Verifying the Deployment

After deploying the operator, you can verify its status using standard Kubernetes commands:

```shell
# Check operator pod status
kubectl get pods -n percona-mysql-operator

# Check operator logs
kubectl logs -n percona-mysql-operator -l app.kubernetes.io/name=kubernetes-percona-mysql-operator

# Verify CRDs are installed
kubectl get crds | grep percona

# Expected CRDs:
# - perconaservermysqls.ps.percona.com
# - perconaservermysqlbackups.ps.percona.com
# - perconaservermysqlrestores.ps.percona.com
```

## Next Steps

Once the operator is deployed, you can:

1. Deploy MySQL clusters using the `MySQLKubernetes` workload resource
2. Create custom PerconaServerMySQL resources directly
3. Configure automated backups and monitoring
4. Scale your MySQL deployments as needed

