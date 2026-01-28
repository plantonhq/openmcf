# Overview

The **MongodbKubernetes** API resource provides an intuitive, production-grade way to deploy and manage MongoDB database clusters on Kubernetes using the **Percona Server for MongoDB Operator**. This Pulumi module interprets a `MongodbKubernetesStackInput`, which includes Kubernetes credentials and your MongoDB cluster specification, and generates a `PerconaServerMongoDB` custom resource that the Percona operator reconciles into a fully functional MongoDB deployment.

## Namespace Management

This module provides flexible namespace management through the `create_namespace` field in the spec:

- **`create_namespace: true`**: Creates a new namespace with resource labels for tracking
  - Ideal for new deployments and isolated environments
  - Namespace is automatically created before MongoDB resources
  
- **`create_namespace: false`**: Uses an existing namespace
  - Namespace must exist before applying this module
  - Suitable for environments with pre-configured policies, quotas, or RBAC
  - Allows sharing the namespace with other resources

## Architecture

The deployment follows a modern, operator-based architecture that separates concerns and leverages Kubernetes-native patterns:

### Component Overview

```
┌─────────────────────────────────────────────────────────────┐
│  Kubernetes Cluster                                         │
│                                                             │
│  ┌───────────────────────────────────────────────────────┐ │
│  │  mongodb-operator namespace                           │ │
│  │  ┌─────────────────────────────────────────────────┐  │ │
│  │  │  Percona MongoDB Operator                       │  │ │
│  │  │  (watches PerconaServerMongoDB CRDs)            │  │ │
│  │  └─────────────────────────────────────────────────┘  │ │
│  └───────────────────────────────────────────────────────┘ │
│                                                             │
│  ┌───────────────────────────────────────────────────────┐ │
│  │  Application Namespace (e.g., my-mongodb)             │ │
│  │                                                         │ │
│  │  1. Pulumi Module Creates:                             │ │
│  │  ┌─────────────────────────────────────────────────┐  │ │
│  │  │  PerconaServerMongoDB CRD                        │  │ │
│  │  │  (replica set config, resources, persistence)    │  │ │
│  │  └─────────────────────────────────────────────────┘  │ │
│  │  ┌─────────────────────────────────────────────────┐  │ │
│  │  │  Secret (auto-generated password)                │  │ │
│  │  └─────────────────────────────────────────────────┘  │ │
│  │                                                         │ │
│  │  2. Operator Reconciles to Create:                     │ │
│  │  ┌─────────────────────────────────────────────────┐  │ │
│  │  │  StatefulSets (MongoDB pods)                     │  │ │
│  │  │  ├─ Pod with persistent volumes                  │  │ │
│  │  │  ├─ Pod with persistent volumes                  │  │ │
│  │  │  └─ ...                                          │  │ │
│  │  └─────────────────────────────────────────────────┘  │ │
│  │  ┌─────────────────────────────────────────────────┐  │ │
│  │  │  Services (cluster communication)                │  │ │
│  │  │  - Cluster service (load balances)               │  │ │
│  │  │  - Headless service (replica set)                │  │ │
│  │  └─────────────────────────────────────────────────┘  │ │
│  │  ┌─────────────────────────────────────────────────┐  │ │
│  │  │  ConfigMaps (MongoDB configuration)              │  │ │
│  │  └─────────────────────────────────────────────────┘  │ │
│  │  ┌─────────────────────────────────────────────────┐  │ │
│  │  │  LoadBalancer Service (optional ingress)         │  │ │
│  │  └─────────────────────────────────────────────────┘  │ │
│  └───────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### Deployment Flow

1. **User defines** MongoDB cluster spec (replicas, resources, persistence)
2. **Pulumi module generates** PerconaServerMongoDB CRD and Secret
3. **Percona operator watches** for PerconaServerMongoDB resources
4. **Operator creates** StatefulSets, Services, ConfigMaps
5. **MongoDB replica set** becomes operational with automatic failover
6. **Operator continuously reconciles** to maintain desired state

### Key Benefits of Operator-Based Architecture

- **Declarative Management**: Define desired state, operator ensures actual state matches
- **Self-Healing**: Operator automatically recovers from failures
- **Rolling Updates**: Zero-downtime upgrades and configuration changes
- **Automated Failover**: Replica sets automatically elect new primaries on failure
- **Best Practices**: Operator encodes years of MongoDB operational expertise

## Replica Set Architecture

MongoDB replica sets provide high availability and data redundancy. The Percona operator automatically configures and manages replica set members.

### Single Replica (Development)

```
┌───────────────────────────────────────────────────────┐
│  Application Namespace                                │
│                                                       │
│  ┌─────────────────────────────────────────────────┐ │
│  │  PerconaServerMongoDB                            │ │
│  │  replicas: 1                                     │ │
│  └─────────────────────────────────────────────────┘ │
│                                                       │
│  ┌─────────────────────────────────────────────────┐ │
│  │  StatefulSet (1 pod)                             │ │
│  │  ┌───────────────────────────────────────────┐  │ │
│  │  │  MongoDB Pod (rs0-0)                       │  │ │
│  │  │  - Primary member                          │  │ │
│  │  │  - CPU: 1 core, Memory: 1Gi               │  │ │
│  │  │  - Persistent Volume: 10Gi                 │  │ │
│  │  │  - Port 27017                              │  │ │
│  │  └───────────────────────────────────────────┘  │ │
│  └─────────────────────────────────────────────────┘ │
│                                                       │
│  ┌─────────────────────────────────────────────────┐ │
│  │  Services                                        │ │
│  │  - Cluster service (mongodb.namespace)          │ │
│  │  - Headless service (for replica set)           │ │
│  └─────────────────────────────────────────────────┘ │
└───────────────────────────────────────────────────────┘
```

**Characteristics:**
- Single point of failure (acceptable for non-critical workloads)
- Simple configuration and operation
- Lower resource usage
- Faster startup time

## High Availability Architecture (3+ Replicas)

For production workloads requiring high availability:

```
┌──────────────────────────────────────────────────────────────────────┐
│  Application Namespace                                               │
│                                                                      │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │  PerconaServerMongoDB                                           │ │
│  │  replicas: 3                                                    │ │
│  └────────────────────────────────────────────────────────────────┘ │
│                                                                      │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │  StatefulSet (3 pods)                                           │ │
│  │                                                                  │ │
│  │  ┌──────────────────────────┐                                   │ │
│  │  │ MongoDB Pod (rs0-0)      │  ← PRIMARY                         │ │
│  │  │ + Persistent Volume      │     Handles writes                │ │
│  │  └──────────────────────────┘                                   │ │
│  │                                                                  │ │
│  │  ┌──────────────────────────┐                                   │ │
│  │  │ MongoDB Pod (rs0-1)      │  ← SECONDARY                       │ │
│  │  │ + Persistent Volume      │     Replicates data from primary  │ │
│  │  └──────────────────────────┘                                   │ │
│  │                                                                  │ │
│  │  ┌──────────────────────────┐                                   │ │
│  │  │ MongoDB Pod (rs0-2)      │  ← SECONDARY                       │ │
│  │  │ + Persistent Volume      │     Replicates data from primary  │ │
│  │  └──────────────────────────┘                                   │ │
│  └────────────────────────────────────────────────────────────────┘ │
│                                                                      │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │  Services                                                       │ │
│  │  - Cluster service (load balances reads across secondaries)    │ │
│  │  - Headless service (replica set member discovery)             │ │
│  └────────────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────────────┘
```

**Characteristics:**
- **High Availability**: Can tolerate 1 node failure (majority = 2/3)
- **Automatic Failover**: Secondaries automatically elect new primary if primary fails
- **Data Replication**: All data replicated to all members
- **Read Distribution**: Reads can be distributed across secondaries (configurable)
- **Zero Downtime**: Rolling updates don't impact availability

### Replica Set Behavior

- **Write Operations**: Always go to PRIMARY
- **Read Operations**: Can be configured to read from SECONDARIES (reduces primary load)
- **Automatic Elections**: If primary fails, secondaries vote to elect new primary (typically 10-30 seconds)
- **Data Consistency**: All members eventually have identical data

## Operator Capabilities

The Percona Server for MongoDB Operator provides enterprise-grade automation:

### Lifecycle Management
- **Rolling Upgrades**: Update MongoDB versions without downtime
- **Scaling**: Add/remove replica set members dynamically
- **Configuration Changes**: Update settings with automatic rollout

### High Availability
- **Automatic Failover**: Replicas take over if primary fails
- **Self-Healing**: Operator restarts failed pods automatically
- **Health Monitoring**: Continuous health checks and auto-recovery

### Storage Management
- **Persistent Volumes**: Automatically provisions PVCs for each pod
- **Volume Expansion**: Supports dynamic volume resizing (if storage class allows)
- **Backup Integration**: Built-in support for automated backups

### Security
- **Secret Management**: Integrates with Kubernetes Secrets for credentials
- **TLS Support**: Can enable TLS for client and inter-replica communication
- **RBAC Integration**: Works with Kubernetes RBAC

## Use Cases

### Application Databases
Deploy MongoDB as the primary database for web applications, mobile apps, and microservices requiring flexible schema and document storage.

### Content Management
Store and retrieve unstructured or semi-structured content like articles, comments, user profiles, and media metadata.

### Real-time Analytics
Handle real-time data ingestion and queries for dashboards, monitoring systems, and operational analytics.

### IoT Data Storage
Store high-volume time-series data from IoT devices with flexible schema requirements.

### Catalog and Inventory
Manage product catalogs, inventory systems, and e-commerce platforms with complex, hierarchical data structures.

### User Profiles and Session Storage
Store user profiles, preferences, and session data for high-traffic applications requiring low-latency access.

## Why Percona Operator?

The Percona Server for MongoDB Operator is a production-grade solution for running MongoDB on Kubernetes:

- **Battle-Tested**: Used by thousands of production deployments worldwide
- **Active Development**: Regular updates and security patches from Percona team
- **Community Support**: Large community, extensive documentation, and professional support available
- **Feature-Rich**: Comprehensive feature set including backups, monitoring, and TLS
- **MongoDB Native**: Built by database experts, optimized for MongoDB-specific patterns
- **Cloud-Agnostic**: Works on any Kubernetes cluster (EKS, GKE, AKS, self-hosted)

## Deployment-Agnostic Design

A key feature of the `MongodbKubernetes` specification is its deployment-agnostic design. The same YAML manifest works regardless of the underlying deployment mechanism:

- Originally implemented with Bitnami Helm charts
- Migrated to Percona operator CRDs **without changing the specification**
- Could be implemented with other tools (Zalando operator, raw Kubernetes manifests) in the future

This design ensures that:
- User configurations remain stable over time
- Migrations between deployment tools are seamless
- The API focuses on user intent, not implementation details

## Summary

The **MongodbKubernetes** module provides a clean, intuitive API that abstracts the complexity of operator-based MongoDB deployments. You define your desired cluster topology and resource allocations; the Percona operator handles all the operational complexity. This results in production-grade MongoDB clusters that are reliable, scalable, and easy to manage.

By leveraging the Percona operator, you benefit from years of operational expertise and best practices encoded into the operator's reconciliation logic, allowing you to focus on your application rather than database operations.
