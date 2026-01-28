# Overview

The **ClickHouseKubernetes** API resource provides an intuitive, production-grade way to deploy and manage ClickHouse database clusters on Kubernetes using the **Altinity ClickHouse Operator**. This Pulumi module interprets a `ClickHouseKubernetesStackInput`, which includes Kubernetes credentials and your ClickHouse cluster specification, and generates a ClickHouseInstallation custom resource that the Altinity operator reconciles into a fully functional ClickHouse deployment.

## Architecture

The deployment follows a modern, operator-based architecture that separates concerns and leverages Kubernetes-native patterns:

### Component Overview

```
┌─────────────────────────────────────────────────────────────┐
│  Kubernetes Cluster                                         │
│                                                             │
│  ┌───────────────────────────────────────────────────────┐ │
│  │  clickhouse-operator namespace                        │ │
│  │  ┌─────────────────────────────────────────────────┐  │ │
│  │  │  Altinity ClickHouse Operator                   │  │ │
│  │  │  (watches ClickHouseInstallation CRDs)          │  │ │
│  │  └─────────────────────────────────────────────────┘  │ │
│  └───────────────────────────────────────────────────────┘ │
│                                                             │
│  ┌───────────────────────────────────────────────────────┐ │
│  │  Application Namespace (e.g., my-clickhouse)          │ │
│  │                                                         │ │
│  │  1. Pulumi Module Creates:                             │ │
│  │  ┌─────────────────────────────────────────────────┐  │ │
│  │  │  ClickHouseInstallation CRD                      │  │ │
│  │  │  (cluster topology, resources, persistence)      │  │ │
│  │  └─────────────────────────────────────────────────┘  │ │
│  │  ┌─────────────────────────────────────────────────┐  │ │
│  │  │  Secret (auto-generated password)                │  │ │
│  │  └─────────────────────────────────────────────────┘  │ │
│  │                                                         │ │
│  │  2. Operator Reconciles to Create:                     │ │
│  │  ┌─────────────────────────────────────────────────┐  │ │
│  │  │  StatefulSets (ClickHouse pods)                  │  │ │
│  │  │  ├─ Pod with persistent volumes                  │  │ │
│  │  │  ├─ Pod with persistent volumes                  │  │ │
│  │  │  └─ ...                                          │  │ │
│  │  └─────────────────────────────────────────────────┘  │ │
│  │  ┌─────────────────────────────────────────────────┐  │ │
│  │  │  Services (cluster communication)                │  │ │
│  │  │  ├─ cluster service                              │  │ │
│  │  │  └─ shard services                               │  │ │
│  │  └─────────────────────────────────────────────────┘  │ │
│  │  ┌─────────────────────────────────────────────────┐  │ │
│  │  │  ConfigMaps (ClickHouse configuration)           │  │ │
│  │  └─────────────────────────────────────────────────┘  │ │
│  │  ┌─────────────────────────────────────────────────┐  │ │
│  │  │  ZooKeeper StatefulSet (for clustering)         │  │ │
│  │  │  (optional, auto-managed by operator)            │  │ │
│  │  └─────────────────────────────────────────────────┘  │ │
│  │  ┌─────────────────────────────────────────────────┐  │ │
│  │  │  LoadBalancer Service (optional ingress)         │  │ │
│  │  └─────────────────────────────────────────────────┘  │ │
│  └───────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### Deployment Flow

1. **User defines** ClickHouse cluster spec (cluster name, resources, topology)
2. **Pulumi module generates** ClickHouseInstallation CRD and Secret
3. **Altinity operator watches** for ClickHouseInstallation resources
4. **Operator creates** StatefulSets, Services, ConfigMaps, ZooKeeper (if needed)
5. **ClickHouse cluster** becomes operational and self-healing
6. **Operator continuously reconciles** to maintain desired state

### Key Benefits of Operator-Based Architecture

- **Declarative Management**: Define desired state, operator ensures actual state matches
- **Self-Healing**: Operator automatically recovers from failures
- **Rolling Updates**: Zero-downtime upgrades and configuration changes
- **Scaling**: Add/remove shards and replicas without manual intervention
- **Best Practices**: Operator encodes years of ClickHouse operational expertise

## Standalone Mode Architecture

For development, testing, or small production workloads, deploy a single ClickHouse instance:

```
┌───────────────────────────────────────────────────────┐
│  Application Namespace                                │
│                                                       │
│  ┌─────────────────────────────────────────────────┐ │
│  │  ClickHouseInstallation                          │ │
│  │  clusterName: my-clickhouse                      │ │
│  │  shards: 1, replicas: 1                          │ │
│  └─────────────────────────────────────────────────┘ │
│                                                       │
│  ┌─────────────────────────────────────────────────┐ │
│  │  StatefulSet (1 pod)                             │ │
│  │  ┌───────────────────────────────────────────┐  │ │
│  │  │  ClickHouse Pod                            │  │ │
│  │  │  - CPU: 2 cores, Memory: 4Gi              │  │ │
│  │  │  - Persistent Volume: 50Gi                 │  │ │
│  │  │  - Port 8123 (HTTP), 9000 (Native)        │  │ │
│  │  └───────────────────────────────────────────┘  │ │
│  └─────────────────────────────────────────────────┘ │
│                                                       │
│  ┌─────────────────────────────────────────────────┐ │
│  │  Service (ClusterIP)                             │ │
│  │  my-clickhouse.namespace.svc.cluster.local       │ │
│  └─────────────────────────────────────────────────┘ │
│                                                       │
│  ┌─────────────────────────────────────────────────┐ │
│  │  LoadBalancer (if ingress enabled)               │ │
│  │  External: my-clickhouse.example.com             │ │
│  └─────────────────────────────────────────────────┘ │
└───────────────────────────────────────────────────────┘
```

**Characteristics:**
- Single point of failure (acceptable for non-critical workloads)
- Simple configuration and operation
- Lower resource usage
- Faster startup time

## Clustered Mode Architecture

For production workloads requiring high availability and horizontal scaling:

```
┌──────────────────────────────────────────────────────────────────────┐
│  Application Namespace                                               │
│                                                                      │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │  ClickHouseInstallation                                         │ │
│  │  clusterName: production-cluster                                │ │
│  │  shards: 3, replicas: 2                                         │ │
│  └────────────────────────────────────────────────────────────────┘ │
│                                                                      │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │  ZooKeeper StatefulSet (auto-managed)                           │ │
│  │  - Coordinates distributed operations                           │ │
│  │  - Manages replica synchronization                              │ │
│  └────────────────────────────────────────────────────────────────┘ │
│                                                                      │
│  ┌─────────────────────────────────┐                                │
│  │  Shard 1                        │                                │
│  │  ┌──────────────────────────┐   │                                │
│  │  │ Replica 1 Pod            │   │  ← Data Shard A                │
│  │  │ + Persistent Volume      │   │                                │
│  │  └──────────────────────────┘   │                                │
│  │  ┌──────────────────────────┐   │                                │
│  │  │ Replica 2 Pod            │   │  ← Data Shard A (copy)         │
│  │  │ + Persistent Volume      │   │                                │
│  │  └──────────────────────────┘   │                                │
│  └─────────────────────────────────┘                                │
│                                                                      │
│  ┌─────────────────────────────────┐                                │
│  │  Shard 2                        │                                │
│  │  ┌──────────────────────────┐   │                                │
│  │  │ Replica 1 Pod            │   │  ← Data Shard B                │
│  │  │ + Persistent Volume      │   │                                │
│  │  └──────────────────────────┘   │                                │
│  │  ┌──────────────────────────┐   │                                │
│  │  │ Replica 2 Pod            │   │  ← Data Shard B (copy)         │
│  │  │ + Persistent Volume      │   │                                │
│  │  └──────────────────────────┘   │                                │
│  └─────────────────────────────────┘                                │
│                                                                      │
│  ┌─────────────────────────────────┐                                │
│  │  Shard 3                        │                                │
│  │  ┌──────────────────────────┐   │                                │
│  │  │ Replica 1 Pod            │   │  ← Data Shard C                │
│  │  │ + Persistent Volume      │   │                                │
│  │  └──────────────────────────┘   │                                │
│  │  ┌──────────────────────────┐   │                                │
│  │  │ Replica 2 Pod            │   │  ← Data Shard C (copy)         │
│  │  │ + Persistent Volume      │   │                                │
│  │  └──────────────────────────┘   │                                │
│  └─────────────────────────────────┘                                │
│                                                                      │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │  Services                                                       │ │
│  │  - Cluster service (load balances across all replicas)         │ │
│  │  - Shard services (direct access to specific shards)           │ │
│  └────────────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────────────┘
```

**Characteristics:**
- **Horizontal Scaling**: Data distributed across 3 shards (A, B, C)
- **High Availability**: 2 replicas per shard (survives 1 node failure per shard)
- **Query Parallelization**: Each shard processes queries in parallel
- **Fault Tolerance**: ZooKeeper ensures consistency during failures
- **Load Balancing**: Queries distributed across replicas

### Data Distribution Example

With 3 shards and 2 replicas:
- **Total Nodes**: 6 ClickHouse pods
- **Data Distribution**: Each shard holds 1/3 of the data
- **Replication**: Each shard's data is duplicated on 2 nodes
- **Query Performance**: ~3x faster for distributed queries (parallel processing)
- **Fault Tolerance**: Can lose 1 node per shard without data loss

## Operator Capabilities

The Altinity ClickHouse Operator provides enterprise-grade automation:

### Lifecycle Management
- **Rolling Upgrades**: Update ClickHouse versions without downtime
- **Scaling**: Add/remove shards and replicas dynamically
- **Configuration Changes**: Update settings with automatic rollout

### High Availability
- **Automatic Failover**: Replicas take over if primary fails
- **ZooKeeper Management**: Operator provisions and manages ZooKeeper automatically
- **Health Monitoring**: Continuous health checks and auto-recovery

### Storage Management
- **Persistent Volumes**: Automatically provisions PVCs for each pod
- **Volume Expansion**: Supports dynamic volume resizing (if storage class allows)
- **Backup Integration**: Hooks for backup and recovery workflows

### Security
- **Secret Management**: Integrates with Kubernetes Secrets for credentials
- **Network Policies**: Supports pod-to-pod encryption
- **RBAC Integration**: Works with Kubernetes RBAC

## Use Cases

### Real-time Analytics
Deploy ClickHouse for sub-second OLAP queries on streaming data from Kafka, Kinesis, or other sources.

### Data Warehousing
Build scalable data warehouses leveraging ClickHouse's columnar storage and compression (90%+ compression ratios).

### Log Analytics
Process and analyze billions of log events per day with millisecond query latency.

### Time-Series Data
Handle time-series workloads like metrics, IoT sensor data, and financial tick data with optimized storage engines.

### Business Intelligence
Power real-time dashboards and BI tools with instant query responses on massive datasets.

### Event Analytics
Track and analyze user behavioral events, clickstreams, and product analytics at scale.

### Machine Learning Features
Generate real-time ML features and aggregations for training and inference pipelines.

## Why Altinity Operator?

The Altinity ClickHouse Operator is the industry standard for running ClickHouse on Kubernetes:

- **Battle-Tested**: Used by thousands of production deployments worldwide
- **Active Development**: Regular updates and security patches from Altinity team
- **Community Support**: Large community, extensive documentation, and professional support available
- **Feature-Rich**: Comprehensive feature set beyond basic deployment
- **ClickHouse Native**: Built by ClickHouse experts, optimized for ClickHouse-specific patterns
- **Cloud-Agnostic**: Works on any Kubernetes cluster (EKS, GKE, AKS, self-hosted)

## Summary

The **ClickHouseKubernetes** module provides a clean, intuitive API that abstracts the complexity of operator-based ClickHouse deployments. You define your desired cluster topology and resource allocations; the Altinity operator handles all the operational complexity. This results in production-grade ClickHouse clusters that are reliable, scalable, and easy to manage.

By leveraging the Altinity operator, you benefit from years of operational expertise and best practices encoded into the operator's reconciliation logic, allowing you to focus on your data and analytics rather than Kubernetes complexity.
