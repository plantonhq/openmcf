# SignozKubernetes - Example Configurations

This document provides a series of examples demonstrating various configurations of the **SignozKubernetes** API resource. Each example shows a typical use case, with corresponding YAML that can be applied via `planton apply -f <filename>`.

---

## 1. Minimal Configuration (Evaluation)

A simple example deploying SigNoz with default settings for evaluation or development purposes.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SignozKubernetes
metadata:
  name: signoz-dev
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: signoz-dev
  createNamespace: true
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

**Key points**:
- Single-node ClickHouse with 20Gi storage
- 1 SigNoz pod and 2 OpenTelemetry Collector pods
- No external ingress configured
- Suitable for development and testing

---

## 2. Production Configuration with High Availability

Demonstrates a production-ready SigNoz deployment with distributed ClickHouse cluster, sharding, replication, and Zookeeper coordination.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SignozKubernetes
metadata:
  name: signoz-production
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: signoz-production
  createNamespace: true
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

**Key points**:
- 2 SigNoz pods for high availability
- 3 OpenTelemetry Collector pods for distributed ingestion
- Distributed ClickHouse with 2 shards and 2 replicas
- 3 Zookeeper nodes for quorum
- 100Gi storage for telemetry data

---

## 3. External ClickHouse with Password from Kubernetes Secret (Recommended)

Example connecting SigNoz to an existing external ClickHouse instance using a password stored in a Kubernetes Secret. This is the recommended approach for production deployments.

First, create a Kubernetes Secret with the ClickHouse password:
```bash
kubectl create secret generic clickhouse-credentials \
  --from-literal=password=your-secure-password \
  --namespace=signoz-external-ch
```

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SignozKubernetes
metadata:
  name: signoz-external-ch
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: signoz-external-ch
  createNamespace: true
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

**Key points**:
- No ClickHouse or Zookeeper deployed by SigNoz
- Connects to external ClickHouse instance
- Password securely stored in Kubernetes Secret
- Reduced operational overhead
- Ideal for centralized database management

---

## 3b. External ClickHouse with Plain String Password

Example connecting SigNoz to an external ClickHouse instance using a plain string password. This approach is suitable for development/testing but not recommended for production.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SignozKubernetes
metadata:
  name: signoz-external-ch-dev
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: signoz-external-ch
  createNamespace: true
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

**Key points**:
- Plain text password in configuration (avoid in production)
- Suitable for development/testing environments

---

## 4. Ingress-Enabled Deployment for External Access

In this example, ingress is enabled to allow external access to both the SigNoz UI and OpenTelemetry Collector endpoints.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SignozKubernetes
metadata:
  name: signoz-ingress
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: signoz-ingress
  createNamespace: true
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

**Key points**:
- Gateway API resources for external access
- External hostname for UI: `signoz.example.com`
- External hostname for OTel Collector HTTP endpoint: `signoz-ingest.example.com`
- Both UI and data ingestion accessible externally via HTTPS

---

## 5. Custom Container Images

This example shows how to specify custom container images for different SigNoz components.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SignozKubernetes
metadata:
  name: signoz-custom-images
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: signoz
  createNamespace: true
  signozContainer:
    replicas: 2
    image:
      repo: signoz/signoz
      tag: "0.40.0"
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
      tag: "0.88.11"
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
          tag: "23.11-alpine"
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

**Key points**:
- Custom SigNoz version specified
- Custom OpenTelemetry Collector version
- Official ClickHouse image instead of Bitnami
- Version control for all components

---

## 6. Custom Helm Values for Advanced Configuration

Demonstrates how to customize the SigNoz deployment using Helm chart values for advanced configurations like data retention and alerting.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SignozKubernetes
metadata:
  name: signoz-custom-helm
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: signoz
  createNamespace: true
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
    clickhouse.ttl.traces: "720"    # 30 days retention for traces
    clickhouse.ttl.metrics: "2160"  # 90 days retention for metrics
    clickhouse.ttl.logs: "168"      # 7 days retention for logs
    signoz.env.SMTP_HOST: smtp.gmail.com
    signoz.env.SMTP_PORT: "587"
    signoz.env.SMTP_USER: alerts@example.com
    signoz.env.SMTP_PASSWORD: my-smtp-password
```

**Key points**:
- Custom data retention policies for each telemetry type
- SMTP configuration for email alerting
- Leverage full power of SigNoz Helm chart
- Fine-tune deployment parameters

---

## 7. High-Volume Ingestion Deployment

Example optimized for environments with high telemetry data volumes, featuring scaled-up components and distributed ClickHouse.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SignozKubernetes
metadata:
  name: signoz-high-volume
spec:
  targetCluster:
    clusterName: prod-gke-cluster
  namespace:
    value: signoz-prod
  createNamespace: true
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

**Key points**:
- 3 SigNoz pods for query load distribution
- 6 OpenTelemetry Collector pods for high ingestion capacity
- Distributed ClickHouse with 3 shards and 2 replicas
- 5 Zookeeper nodes for quorum
- 500Gi storage for high data volume

---

## 8. Using Alternative Docker Registry (Avoiding Bitnami Legacy)

Override the default bitnamilegacy registry to use official ClickHouse images or other sources.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: SignozKubernetes
metadata:
  name: signoz-official-images
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: signoz
  createNamespace: true
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
  helmValues:
    global.imageRegistry: ""  # Clear default registry
    clickhouse.image.registry: docker.io
    clickhouse.image.repository: clickhouse/clickhouse-server
    clickhouse.image.tag: "24.8"
    zookeeper.image.registry: docker.io
    zookeeper.image.repository: zookeeper
    zookeeper.image.tag: "3.8"
```

**Key points**:
- Override default bitnamilegacy registry
- Use official vendor images directly
- Provides long-term stability
- Custom image tags for version control

---

## Conclusion

These examples illustrate the breadth of **SignozKubernetes** features, from basic single-node deployments to advanced distributed clusters with high availability. By consolidating Kubernetes manifests behind a concise API resource definition, you can maintain consistency, reduce error-prone manual config, and accelerate delivery cycles.

> **Getting Started**
> 1. Create a YAML file for your SigNoz cluster (e.g., `signoz.yaml`).
> 2. Run:
>    ```shell
>    planton apply -f signoz.yaml
>    ```
> 3. Verify the deployment:
>    ```shell
>    kubectl get pods -n <namespace>
>    kubectl logs -n <namespace> <pod-name>
>    ```

For additional details, see the [SignozKubernetes API documentation](../README.md), or reach out to our support team.

