# SigNoz Kubernetes API - Example Configurations

## Example w/ Basic Configuration (Evaluation)

### Create Using CLI

Create a YAML file using the example shown below. After the YAML is created, use the following command to apply it:

```shell
planton apply -f <yaml-path>
```

### Basic Example

This basic example demonstrates a minimal configuration for deploying SigNoz Kubernetes for evaluation or development purposes, using default settings with a single-node ClickHouse instance.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SignozKubernetes
metadata:
  name: basic-signoz
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: "signoz"
  create_namespace: true
  signozContainer:
    replicas: 1
    resources:
      requests:
        cpu: 200m
        memory: 512Mi
      limits:
        cpu: 1000m
        memory: 2Gi
  otelCollectorContainer:
    replicas: 2
    resources:
      requests:
        cpu: 500m
        memory: 1Gi
      limits:
        cpu: 2000m
        memory: 4Gi
  database:
    isExternal: false
    managedDatabase:
      container:
        replicas: 1
        isPersistenceEnabled: true
        diskSize: 20Gi
        resources:
          requests:
            cpu: 500m
            memory: 1Gi
          limits:
            cpu: 2000m
            memory: 4Gi
      cluster:
        isEnabled: false
      zookeeper:
        isEnabled: false
```

---

## Example w/ Production Configuration (High Availability)

This example demonstrates a production-ready deployment with high availability, featuring a distributed ClickHouse cluster with sharding and replication, along with Zookeeper coordination.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SignozKubernetes
metadata:
  name: production-signoz
spec:
  target_cluster:
    cluster_name: "prod-gke-cluster"
  namespace:
    value: "signoz-prod"
  create_namespace: true
  signozContainer:
    replicas: 2
    resources:
      requests:
        cpu: 500m
        memory: 1Gi
      limits:
        cpu: 2000m
        memory: 4Gi
  otelCollectorContainer:
    replicas: 3
    resources:
      requests:
        cpu: 1000m
        memory: 2Gi
      limits:
        cpu: 4000m
        memory: 8Gi
  database:
    isExternal: false
    managedDatabase:
      container:
        replicas: 1
        isPersistenceEnabled: true
        diskSize: 100Gi
        resources:
          requests:
            cpu: 1000m
            memory: 4Gi
          limits:
            cpu: 4000m
            memory: 16Gi
      cluster:
        isEnabled: true
        shardCount: 2
        replicaCount: 2
      zookeeper:
        isEnabled: true
        container:
          replicas: 3
          diskSize: 10Gi
          resources:
            requests:
              cpu: 200m
              memory: 512Mi
            limits:
              cpu: 1000m
              memory: 1Gi
```

---

## Example w/ External ClickHouse (Password from Kubernetes Secret - Recommended)

This example demonstrates how to connect SigNoz to an external ClickHouse instance using a password stored in a Kubernetes Secret. This is the recommended approach for production deployments.

First, create a Kubernetes Secret with the ClickHouse password:
```bash
kubectl create secret generic clickhouse-credentials \
  --from-literal=password=your-secure-password \
  --namespace=signoz
```

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SignozKubernetes
metadata:
  name: signoz-external-clickhouse
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: "signoz"
  create_namespace: true
  signozContainer:
    replicas: 2
    resources:
      requests:
        cpu: 500m
        memory: 1Gi
      limits:
        cpu: 2000m
        memory: 4Gi
  otelCollectorContainer:
    replicas: 3
    resources:
      requests:
        cpu: 1000m
        memory: 2Gi
      limits:
        cpu: 4000m
        memory: 8Gi
  database:
    isExternal: true
    externalDatabase:
      host: clickhouse.database.svc.cluster.local
      httpPort: 8123
      tcpPort: 9000
      clusterName: cluster
      isSecure: false
      username: signoz
      # Reference an existing Kubernetes Secret for the password
      password:
        secretRef:
          name: clickhouse-credentials
          key: password
```

---

## Example w/ External ClickHouse (Plain String Password)

This example demonstrates how to connect SigNoz to an external ClickHouse instance using a plain string password. This approach is suitable for development/testing but not recommended for production.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SignozKubernetes
metadata:
  name: signoz-external-clickhouse-dev
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: "signoz-dev"
  create_namespace: true
  signozContainer:
    replicas: 2
    resources:
      requests:
        cpu: 500m
        memory: 1Gi
      limits:
        cpu: 2000m
        memory: 4Gi
  otelCollectorContainer:
    replicas: 3
    resources:
      requests:
        cpu: 1000m
        memory: 2Gi
      limits:
        cpu: 4000m
        memory: 8Gi
  database:
    isExternal: true
    externalDatabase:
      host: clickhouse.database.svc.cluster.local
      httpPort: 8123
      tcpPort: 9000
      clusterName: cluster
      isSecure: false
      username: signoz
      # Plain string password (not recommended for production)
      password:
        stringValue: my-secure-password
```

---

## Example w/ Ingress Configuration

This example demonstrates how to configure ingress for both the SigNoz UI and the OpenTelemetry Collector ingestion endpoints.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SignozKubernetes
metadata:
  name: signoz-with-ingress
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: "signoz"
  create_namespace: true
  signozContainer:
    replicas: 2
    resources:
      requests:
        cpu: 500m
        memory: 1Gi
      limits:
        cpu: 2000m
        memory: 4Gi
  otelCollectorContainer:
    replicas: 3
    resources:
      requests:
        cpu: 1000m
        memory: 2Gi
      limits:
        cpu: 4000m
        memory: 8Gi
  database:
    isExternal: false
    managedDatabase:
      container:
        replicas: 1
        isPersistenceEnabled: true
        diskSize: 50Gi
        resources:
          requests:
            cpu: 500m
            memory: 2Gi
          limits:
            cpu: 2000m
            memory: 8Gi
      cluster:
        isEnabled: false
      zookeeper:
        isEnabled: false
  ingress:
    ui:
      enabled: true
      hostname: signoz.example.com
    otelCollector:
      enabled: true
      hostname: signoz-ingest.example.com
```

---

## Example w/ Custom Container Images

This example shows how to specify custom container images for different SigNoz components, useful for using specific versions or custom-built images.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SignozKubernetes
metadata:
  name: signoz-custom-images
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: "signoz"
  create_namespace: true
  signozContainer:
    replicas: 2
    image:
      repo: signoz/signoz
      tag: 0.40.0
    resources:
      requests:
        cpu: 500m
        memory: 1Gi
      limits:
        cpu: 2000m
        memory: 4Gi
  otelCollectorContainer:
    replicas: 2
    image:
      repo: signoz/signoz-otel-collector
      tag: 0.88.11
    resources:
      requests:
        cpu: 500m
        memory: 1Gi
      limits:
        cpu: 2000m
        memory: 4Gi
  database:
    isExternal: false
    managedDatabase:
      container:
        replicas: 1
        image:
          repo: clickhouse/clickhouse-server
          tag: 23.11-alpine
        isPersistenceEnabled: true
        diskSize: 50Gi
        resources:
          requests:
            cpu: 500m
            memory: 2Gi
          limits:
            cpu: 2000m
            memory: 8Gi
      cluster:
        isEnabled: false
      zookeeper:
        isEnabled: false
```

---

## Example w/ Custom Helm Values

This example demonstrates how to provide additional customization through Helm chart values, such as configuring data retention policies and alerting integrations.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SignozKubernetes
metadata:
  name: signoz-custom-helm
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: "signoz"
  create_namespace: true
  signozContainer:
    replicas: 2
    resources:
      requests:
        cpu: 500m
        memory: 1Gi
      limits:
        cpu: 2000m
        memory: 4Gi
  otelCollectorContainer:
    replicas: 3
    resources:
      requests:
        cpu: 1000m
        memory: 2Gi
      limits:
        cpu: 4000m
        memory: 8Gi
  database:
    isExternal: false
    managedDatabase:
      container:
        replicas: 1
        isPersistenceEnabled: true
        diskSize: 100Gi
        resources:
          requests:
            cpu: 1000m
            memory: 4Gi
          limits:
            cpu: 4000m
            memory: 16Gi
      cluster:
        isEnabled: false
      zookeeper:
        isEnabled: false
  helmValues:
    # Configure data retention for traces (in hours)
    clickhouse.ttl.traces: "720"  # 30 days
    # Configure data retention for metrics (in hours)
    clickhouse.ttl.metrics: "2160"  # 90 days
    # Configure data retention for logs (in hours)
    clickhouse.ttl.logs: "168"  # 7 days
    # Enable email alerting
    signoz.env.SMTP_HOST: smtp.gmail.com
    signoz.env.SMTP_PORT: "587"
    signoz.env.SMTP_USER: alerts@example.com
    signoz.env.SMTP_PASSWORD: my-smtp-password
```

---

## Example w/ High-Volume Ingestion

This example is optimized for environments with high telemetry data volumes, featuring scaled-up OpenTelemetry Collectors and ClickHouse resources.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SignozKubernetes
metadata:
  name: signoz-high-volume
spec:
  target_cluster:
    cluster_name: "prod-gke-cluster"
  namespace:
    value: "signoz-prod"
  create_namespace: true
  signozContainer:
    replicas: 3
    resources:
      requests:
        cpu: 1000m
        memory: 2Gi
      limits:
        cpu: 4000m
        memory: 8Gi
  otelCollectorContainer:
    replicas: 6
    resources:
      requests:
        cpu: 2000m
        memory: 4Gi
      limits:
        cpu: 8000m
        memory: 16Gi
  database:
    isExternal: false
    managedDatabase:
      container:
        replicas: 1
        isPersistenceEnabled: true
        diskSize: 500Gi
        resources:
          requests:
            cpu: 4000m
            memory: 16Gi
          limits:
            cpu: 16000m
            memory: 64Gi
      cluster:
        isEnabled: true
        shardCount: 3
        replicaCount: 2
      zookeeper:
        isEnabled: true
        container:
          replicas: 5
          diskSize: 20Gi
          resources:
            requests:
              cpu: 500m
              memory: 1Gi
            limits:
              cpu: 2000m
              memory: 4Gi
```

---

## Example w/ Minimal Resources (Development)

This example uses minimal resource allocations suitable for local development or resource-constrained environments.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SignozKubernetes
metadata:
  name: signoz-dev
spec:
  target_cluster:
    cluster_name: "dev-cluster"
  namespace:
    value: "signoz-dev"
  create_namespace: true
  signozContainer:
    replicas: 1
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 500m
        memory: 1Gi
  otelCollectorContainer:
    replicas: 1
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 500m
        memory: 1Gi
  database:
    isExternal: false
    managedDatabase:
      container:
        replicas: 1
        isPersistenceEnabled: false
        resources:
          requests:
            cpu: 100m
            memory: 512Mi
          limits:
            cpu: 1000m
            memory: 2Gi
      cluster:
        isEnabled: false
      zookeeper:
        isEnabled: false
```

---

## Example w/ S3 Archiving (via Helm Values)

This example demonstrates configuring the OpenTelemetry Collector to archive telemetry data to S3 for long-term retention using Helm values.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SignozKubernetes
metadata:
  name: signoz-s3-archiving
spec:
  target_cluster:
    cluster_name: "prod-gke-cluster"
  namespace:
    value: "signoz-prod"
  create_namespace: true
  signozContainer:
    replicas: 2
    resources:
      requests:
        cpu: 500m
        memory: 1Gi
      limits:
        cpu: 2000m
        memory: 4Gi
  otelCollectorContainer:
    replicas: 3
    resources:
      requests:
        cpu: 1000m
        memory: 2Gi
      limits:
        cpu: 4000m
        memory: 8Gi
  database:
    isExternal: false
    managedDatabase:
      container:
        replicas: 1
        isPersistenceEnabled: true
        diskSize: 100Gi
        resources:
          requests:
            cpu: 1000m
            memory: 4Gi
          limits:
            cpu: 4000m
            memory: 16Gi
      cluster:
        isEnabled: false
      zookeeper:
        isEnabled: false
  helmValues:
    # Configure S3 archiving for long-term retention
    otelCollector.exporters.awss3.region: us-west-2
    otelCollector.exporters.awss3.s3_bucket: my-signoz-archive-bucket
    otelCollector.exporters.awss3.s3_prefix: signoz-telemetry
    otelCollector.exporters.awss3.compression: gzip
```

---

---

## Example w/ Existing Namespace (create_namespace: false)

This example demonstrates using an existing namespace that is managed separately. The namespace must already exist on the cluster before deploying SigNoz.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SignozKubernetes
metadata:
  name: signoz-shared-namespace
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: "observability"  # This namespace must already exist
  create_namespace: false
  signozContainer:
    replicas: 2
    resources:
      requests:
        cpu: 500m
        memory: 1Gi
      limits:
        cpu: 2000m
        memory: 4Gi
  otelCollectorContainer:
    replicas: 3
    resources:
      requests:
        cpu: 1000m
        memory: 2Gi
      limits:
        cpu: 4000m
        memory: 8Gi
  database:
    isExternal: false
    managedDatabase:
      container:
        replicas: 1
        isPersistenceEnabled: true
        diskSize: 50Gi
        resources:
          requests:
            cpu: 500m
            memory: 2Gi
          limits:
            cpu: 2000m
            memory: 8Gi
      cluster:
        isEnabled: false
      zookeeper:
        isEnabled: false
```

**Important**: When `create_namespace` is `false`, you must ensure the namespace exists before deploying. This is useful when:
- The namespace is shared with other applications
- Namespace creation is controlled by a separate team or process
- You're following organizational policies that require centralized namespace management

---

## Important Notes

### Namespace Management

All examples above use `create_namespace: true`, which means the component will create the namespace automatically. If you prefer to manage the namespace separately:

- Set `create_namespace: false`
- Ensure the namespace exists on the cluster before deployment
- The deployment will fail if the namespace doesn't exist

See the "Example w/ Existing Namespace" above for a complete example.

### Ingress Configuration

When configuring ingress for the OpenTelemetry Collector, you may need to create separate ingress resources for gRPC and HTTP endpoints due to annotation requirements:

- **gRPC Endpoint**: Requires `nginx.ingress.kubernetes.io/backend-protocol: "GRPC"` annotation
- **HTTP Endpoint**: Uses standard HTTP routing without special annotations

Consult the SigNoz Helm chart documentation for detailed ingress configuration examples.

### Storage Class

Ensure your Kubernetes cluster has a default StorageClass configured for dynamic persistent volume provisioning. If not, you may need to specify a storage class through Helm values:

```yaml
helmValues:
  global.storageClass: gp2-csi  # or your cluster's storage class name
```

### Resource Planning

The resource allocations in these examples are starting points. Monitor your actual usage using SigNoz's built-in Kubernetes dashboards and adjust based on:
- Number of monitored applications
- Telemetry data volume (traces/second, metrics/second, logs/second)
- Query complexity and concurrency
- Data retention policies

### External ClickHouse Prerequisites

When using external ClickHouse, ensure:
- The external ClickHouse instance is accessible from the Kubernetes cluster
- A Zookeeper instance is available for distributed cluster support
- The ClickHouse cluster is properly configured with a distributed cluster named "cluster" (or update `clusterName`)
- The provided credentials have necessary permissions to create databases and tables

### High Availability Recommendations

For production deployments with high availability:
- Use 2+ replicas for SigNoz pods
- Scale OpenTelemetry Collector based on ingestion volume (start with 3+)
- Enable ClickHouse clustering with 2 shards and 2 replicas per shard
- Deploy Zookeeper with 3 or 5 nodes (odd number for quorum)
- Configure pod anti-affinity to spread replicas across nodes
- Monitor persistent volume usage to prevent outages

