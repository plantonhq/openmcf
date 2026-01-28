# ClickHouseKubernetes - Example Configurations

This document provides examples demonstrating various configurations of the **ClickHouseKubernetes** API resource using the Altinity ClickHouse Operator. Each example shows a typical use case with corresponding YAML that can be applied via `planton pulumi up`.

---

## Prerequisites

The **Altinity ClickHouse Operator** must be installed on your cluster before deploying these examples. See the main [README.md](README.md) for installation instructions.

---

## Namespace Management

All examples below should include the `create_namespace` field to control namespace management:

- **`create_namespace: true`** - The component creates and manages the namespace with appropriate labels
- **`create_namespace: false`** - Deploy into an existing namespace (must exist before deployment)

**Example with managed namespace:**
```yaml
spec:
  namespace:
    value: "clickhouse-prod"
  create_namespace: true  # Component creates the namespace
```

**Example with existing namespace:**
```yaml
spec:
  namespace:
    value: "data-platform"
  create_namespace: false  # Must exist before deployment
```

---

## 1. Minimal Standalone Instance

A simple single-node ClickHouse deployment suitable for development and testing.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: ClickHouseKubernetes
metadata:
  name: dev-clickhouse
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: dev-clickhouse
  create_namespace: true
  clusterName: dev-cluster
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
```

**Key points**:
- Single replica with 20Gi persistent volume
- Suitable resource allocations for development
- Operator manages all Kubernetes resources
- No clustering - standalone mode

---

## 2. Production Standalone with Custom Cluster Name

Production-ready single-node deployment with larger resources and explicit cluster naming.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: ClickHouseKubernetes
metadata:
  name: prod-clickhouse
  org: my-org
  env: production
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: prod-clickhouse
  create_namespace: true
  clusterName: production-analytics
  version: "24.3"
  container:
    replicas: 1
    isPersistenceEnabled: true
    diskSize: 200Gi
    resources:
      requests:
        cpu: 2000m
        memory: 8Gi
      limits:
        cpu: 8000m
        memory: 32Gi
```

**Key points**:
- Custom cluster name: `production-analytics`
- Pinned ClickHouse version for stability
- Large persistent volume (200Gi)
- High resource allocation for production workloads
- Organization and environment metadata for tracking

---

## 3. Distributed Cluster with Auto-Managed ClickHouse Keeper

Horizontally scaled cluster with multiple shards for parallel query processing using auto-managed ClickHouse Keeper.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: ClickHouseKubernetes
metadata:
  name: distributed-analytics
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: distributed-analytics
  create_namespace: true
  clusterName: analytics-cluster
  container:
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
    shardCount: 4
    replicaCount: 1
  # coordination: not specified = auto-managed ClickHouse Keeper (default)
```

**Key points**:
- 4 shards for horizontal scaling
- 1 replica per shard (4 total nodes)
- ClickHouse Keeper automatically referenced (80% use case)
- Distributed query execution across shards
- Each shard processes subset of data
- Note: Deploy ClickHouseKeeperInstallation separately first

---

## 4. High Availability Cluster with Production Keeper Configuration

Production cluster with replication for high availability and 3-node ClickHouse Keeper ensemble.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: ClickHouseKubernetes
metadata:
  name: ha-clickhouse
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: ha-clickhouse
  create_namespace: true
  clusterName: ha-analytics
  version: "24.3"
  container:
    isPersistenceEnabled: true
    diskSize: 150Gi
    resources:
      requests:
        cpu: 2000m
        memory: 8Gi
      limits:
        cpu: 6000m
        memory: 24Gi
  cluster:
    isEnabled: true
    shardCount: 3
    replicaCount: 3
  coordination:
    type: keeper
    keeperConfig:
      replicas: 3  # Production HA: survives 1 node failure
      diskSize: "20Gi"
      resources:
        requests:
          cpu: "200m"
          memory: "512Mi"
        limits:
          cpu: "1000m"
          memory: "2Gi"
```

**Key points**:
- 3 shards with 3 replicas each (9 total ClickHouse nodes)
- High availability - survives 2 node failures per shard
- Data automatically replicated across replicas
- Query load balanced across replicas
- 3-node ClickHouse Keeper for production coordination (survives 1 Keeper failure)
- Custom Keeper resources for larger clusters

---

## 5. Cluster with External ClickHouse Keeper

Advanced configuration using external ClickHouse Keeper for shared coordination across multiple ClickHouse clusters.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: ClickHouseKubernetes
metadata:
  name: enterprise-clickhouse
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: enterprise-clickhouse
  create_namespace: true
  clusterName: enterprise-cluster
  container:
    isPersistenceEnabled: true
    diskSize: 200Gi
    resources:
      requests:
        cpu: 4000m
        memory: 16Gi
      limits:
        cpu: 12000m
        memory: 48Gi
  cluster:
    isEnabled: true
    shardCount: 6
    replicaCount: 2
  coordination:
    type: external_keeper
    externalConfig:
      nodes:
        - "keeper-shared:2181"
```

**Key points**:
- External ClickHouse Keeper for production-grade coordination
- 6 shards with 2 replicas (12 total nodes)
- Enterprise-scale resources
- Shared Keeper can coordinate multiple ClickHouse clusters
- Better isolation and control over coordination service

---

## 5b. Cluster with External ZooKeeper (Legacy)

For environments with existing ZooKeeper infrastructure (e.g., shared with Kafka).

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: ClickHouseKubernetes
metadata:
  name: enterprise-with-zk
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: enterprise-with-zk
  create_namespace: true
  clusterName: enterprise-cluster
  cluster:
    isEnabled: true
    shardCount: 6
    replicaCount: 2
  coordination:
    type: external_zookeeper
    externalConfig:
      nodes:
        - "zk-0.zookeeper.default.svc.cluster.local:2181"
        - "zk-1.zookeeper.default.svc.cluster.local:2181"
        - "zk-2.zookeeper.default.svc.cluster.local:2181"
```

**Key points**:
- External ZooKeeper for legacy integration
- Shared ZooKeeper infrastructure with Kafka, Solr, etc.
- Multiple ZK nodes for redundancy

---

## 6. Development Cluster with Ephemeral Storage

Lightweight cluster for testing and development with no persistence.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: ClickHouseKubernetes
metadata:
  name: test-clickhouse
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: test-clickhouse
  create_namespace: true
  clusterName: test-cluster
  container:
    replicas: 1
    isPersistenceEnabled: false
    resources:
      requests:
        cpu: 100m
        memory: 512Mi
      limits:
        cpu: 500m
        memory: 2Gi
```

**Key points**:
- No persistent storage - data lost on restart
- Minimal resources for quick testing
- Fast startup and teardown
- Cost-effective for temporary workloads

---

## 7. Ingress-Enabled Cluster with External Access

Cluster exposed to external clients via LoadBalancer with DNS.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: ClickHouseKubernetes
metadata:
  name: public-analytics
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: public-analytics
  create_namespace: true
  clusterName: public-cluster
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
  ingress:
    enabled: true
    hostname: clickhouse.example.com
```

**Key points**:
- LoadBalancer service for external access
- External DNS: `clickhouse.example.com`
- Both HTTP (8123) and native (9000) ports exposed
- Suitable for external analytics tools and dashboards
- Consider security implications (firewall, authentication)

---

## 8. SigNoz Backend Configuration

Optimized configuration for deploying ClickHouse as a SigNoz observability backend.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: ClickHouseKubernetes
metadata:
  name: signoz-backend
  org: my-org
  env: production
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: signoz-backend
  create_namespace: true
  clusterName: cluster  # SigNoz requires cluster name to be "cluster"
  version: "24.8"
  container:
    isPersistenceEnabled: true
    diskSize: "50Gi"  # Scale based on retention needs
    resources:
      requests:
        cpu: "500m"
        memory: "1Gi"
      limits:
        cpu: "2000m"
        memory: "4Gi"
  cluster:
    isEnabled: true
    shardCount: 2      # Start with 2, scale as data grows
    replicaCount: 2    # HA for production observability
  coordination:
    type: keeper
    keeperConfig:
      replicas: 3  # Production HA for observability infrastructure
  ingress:
    enabled: true
    hostname: signoz-clickhouse.example.com
```

**Key points**:
- Cluster name "cluster" is required by SigNoz (hardcoded in migrations)
- 2 shards × 2 replicas = 4 ClickHouse pods for distributed data
- 3 ClickHouse Keeper nodes for coordination reliability
- Ingress enabled for external access
- Sized for production observability workloads
- ClickHouse Keeper (not ZooKeeper) for resource efficiency

---

## 9. Multi-Environment with Custom Versions

Example showing version pinning for different environments.

```yaml
# Staging Environment
apiVersion: kubernetes.openmcf.org/v1
kind: ClickHouseKubernetes
metadata:
  name: staging-clickhouse
  env: staging
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: staging-clickhouse
  create_namespace: true
  clusterName: staging-analytics
  version: "24.4"  # Testing newer version
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

---

# Production Environment
apiVersion: kubernetes.openmcf.org/v1
kind: ClickHouseKubernetes
metadata:
  name: production-clickhouse
  env: production
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: production-clickhouse
  create_namespace: true
  clusterName: production-analytics
  version: "24.3"  # Stable version
  container:
    isPersistenceEnabled: true
    diskSize: 200Gi
    resources:
      requests:
        cpu: 2000m
        memory: 8Gi
      limits:
        cpu: 8000m
        memory: 32Gi
  cluster:
    isEnabled: true
    shardCount: 4
    replicaCount: 2
```

**Key points**:
- Different versions for testing vs production
- Staging validates new versions before production
- Environment-specific resource allocations
- Consistent configuration structure across environments

---

## Deployment Instructions

1. **Ensure Operator is Installed**:
   ```bash
   kubectl get pods -n clickhouse-operator
   ```

2. **Apply Configuration**:
   ```bash
   planton pulumi up --stack-input clickhouse.yaml \
     --module-dir apis/project/planton/provider/kubernetes/clickhousekubernetes/v1/iac/pulumi
   ```

3. **Verify Deployment**:
   ```bash
   # Check ClickHouseInstallation
   kubectl get clickhouseinstallations -n <namespace>
   
   # Check pods
   kubectl get pods -n <namespace>
   
   # Check services
   kubectl get svc -n <namespace>
   ```

4. **Access ClickHouse**:
   ```bash
   # Use port-forward command from stack outputs
   kubectl port-forward -n <namespace> service/<cluster-name> 8123:8123
   
   # Connect with clickhouse-client
   clickhouse-client --host localhost --port 8123
   ```

---

## Best Practices

### Resource Sizing
- **Development**: 500m CPU, 1-2Gi memory
- **Production standalone**: 2-4 CPU cores, 8-16Gi memory  
- **Production clustered**: 4-8 CPU cores, 16-48Gi memory per node

### Storage Planning
- **Small datasets** (<100GB): 20-50Gi per node
- **Medium datasets** (100GB-1TB): 100-500Gi per node
- **Large datasets** (>1TB): 500Gi-2Ti per node, consider sharding

### Clustering Strategy
- **Development/Testing**: Standalone (1 node)
- **Small production**: 2-3 shards, 2 replicas (4-6 nodes)
- **Large production**: 4-8 shards, 2-3 replicas (8-24 nodes)

### Version Management
- Pin specific versions in production
- Test new versions in staging first
- Follow Altinity/ClickHouse release notes for breaking changes

---

## Troubleshooting

### Check Operator Logs
```bash
kubectl logs -n clickhouse-operator deployment/clickhouse-operator
```

### Check ClickHouseInstallation Status
```bash
kubectl describe clickhouseinstallation <cluster-name> -n <namespace>
```

### Check ClickHouse Logs
```bash
kubectl logs -n <namespace> <pod-name>
```

### Verify Coordination Connectivity (Clustered Mode)
```bash
# Verify ClickHouse can connect to coordination service (Keeper or ZooKeeper)
kubectl exec -n <namespace> <pod-name> -- clickhouse-client --query "SELECT * FROM system.zookeeper WHERE path='/'"

# Check cluster topology
kubectl exec -n <namespace> <pod-name> -- clickhouse-client --query "SELECT * FROM system.clusters"
```

### Check ClickHouse Keeper Status (if using auto-managed Keeper)
```bash
# Check ClickHouseKeeperInstallation resource
kubectl get clickhousekeeperinstallation -n <namespace>

# Check Keeper pods
kubectl get pods -n <namespace> -l "app.kubernetes.io/name=clickhouse-keeper"
```

---

For additional details and advanced configurations, see the [ClickHouseKubernetes API documentation](../README.md) and the [Altinity Operator documentation](https://github.com/Altinity/clickhouse-operator).
