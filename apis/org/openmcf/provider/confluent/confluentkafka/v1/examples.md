# Confluent Kafka Examples

This document provides comprehensive examples of how to use the ConfluentKafka resource to deploy Kafka clusters on Confluent Cloud across different environments and configurations.

## Table of Contents

- [Basic Development Cluster](#basic-development-cluster)
- [Standard Production Cluster](#standard-production-cluster)
- [Enterprise Cluster with Private Networking](#enterprise-cluster-with-private-networking)
- [Dedicated Cluster with Provisioned Capacity](#dedicated-cluster-with-provisioned-capacity)
- [Multi-Region Disaster Recovery Setup](#multi-region-disaster-recovery-setup)

---

## Basic Development Cluster

A minimal configuration suitable for development and testing environments. Uses the BASIC cluster type with single-zone deployment for cost efficiency.

```yaml
apiVersion: confluent.openmcf.org/v1
kind: ConfluentKafka
metadata:
  name: dev-kafka
spec:
  cloud: AWS
  region: us-east-2
  availability: SINGLE_ZONE
  environmentId: env-dev-abc123
  clusterType: BASIC
```

**Key characteristics:**
- **Cost:** Lowest cost option (pay-as-you-go pricing)
- **Performance:** Suitable for development and testing
- **Availability:** Single-zone only, no SLA
- **Networking:** Public internet only

---

## Standard Production Cluster

A production-ready cluster with multi-zone high availability and elastic scaling. Recommended for most production workloads.

```yaml
apiVersion: confluent.openmcf.org/v1
kind: ConfluentKafka
metadata:
  name: prod-orders-kafka
  labels:
    environment: production
    team: platform
    cost-center: engineering
spec:
  cloud: GCP
  region: us-central1
  availability: MULTI_ZONE
  environmentId: env-prod-xyz789
  clusterType: STANDARD
  displayName: Production Orders Kafka Cluster
```

**Key characteristics:**
- **Cost:** Elastic pricing (pay for what you use)
- **Performance:** Auto-scales based on throughput
- **Availability:** Multi-zone with 99.99% SLA
- **Networking:** Public internet only

---

## Enterprise Cluster with Private Networking

An enterprise-grade cluster with private networking using AWS PrivateLink for secure, private connectivity from your VPC.

```yaml
apiVersion: confluent.openmcf.org/v1
kind: ConfluentKafka
metadata:
  name: enterprise-kafka
  labels:
    environment: production
    security-zone: restricted
    compliance: pci-dss
spec:
  cloud: AWS
  region: us-west-2
  availability: MULTI_ZONE
  environmentId: env-prod-enterprise
  clusterType: ENTERPRISE
  displayName: Enterprise Secure Kafka Cluster
  networkConfig:
    networkId: n-abc123
```

**Prerequisites:**
- A Confluent Cloud network resource (`n-abc123`) must be pre-created in the same environment
- The network must be configured for AWS PrivateLink
- Your VPC must have a PrivateLink endpoint configured

**Key characteristics:**
- **Cost:** Elastic pricing (similar to STANDARD)
- **Performance:** Auto-scales based on throughput
- **Availability:** Multi-zone with 99.99% SLA
- **Networking:** Private connectivity via PrivateLink
- **Security:** Traffic never traverses the public internet

---

## Dedicated Cluster with Provisioned Capacity

A dedicated, single-tenant cluster with provisioned capacity (CKU) for predictable, high-throughput workloads.

```yaml
apiVersion: confluent.openmcf.org/v1
kind: ConfluentKafka
metadata:
  name: dedicated-high-throughput
  labels:
    environment: production
    workload: high-throughput
    criticality: tier-1
spec:
  cloud: AZURE
  region: eastus
  availability: MULTI_ZONE
  environmentId: env-prod-critical
  clusterType: DEDICATED
  displayName: High-Throughput Dedicated Kafka
  dedicatedConfig:
    cku: 4
  networkConfig:
    networkId: n-azure-private-123
```

**Prerequisites:**
- A Confluent Cloud network resource must be pre-created for private networking (optional but recommended)

**Key characteristics:**
- **Cost:** Provisioned pricing (pay for CKU capacity)
- **Performance:** Predictable, high-throughput (scales with CKU count)
- **Availability:** Multi-zone with 99.99% SLA
- **Networking:** Supports private connectivity via Azure Private Link
- **Tenancy:** Single-tenant isolation
- **Scaling:** Manual CKU adjustments (can scale up/down)

**CKU Sizing Guidelines:**
- **1-2 CKU:** Small to medium workloads (up to 50 MB/s ingress)
- **4 CKU:** High-throughput workloads (up to 150 MB/s ingress)
- **8+ CKU:** Very high-throughput or mission-critical workloads

---

## Multi-Region Disaster Recovery Setup

Deploy clusters in multiple regions for disaster recovery and business continuity.

### Primary Cluster (us-east-1)

```yaml
apiVersion: confluent.openmcf.org/v1
kind: ConfluentKafka
metadata:
  name: prod-primary-kafka
  labels:
    environment: production
    region-role: primary
    dr-pair: prod-kafka-dr
spec:
  cloud: AWS
  region: us-east-1
  availability: MULTI_ZONE
  environmentId: env-prod-primary
  clusterType: ENTERPRISE
  displayName: Production Primary Kafka (us-east-1)
  networkConfig:
    networkId: n-us-east-1
```

### Secondary Cluster (us-west-2)

```yaml
apiVersion: confluent.openmcf.org/v1
kind: ConfluentKafka
metadata:
  name: prod-secondary-kafka
  labels:
    environment: production
    region-role: secondary
    dr-pair: prod-kafka-dr
spec:
  cloud: AWS
  region: us-west-2
  availability: MULTI_ZONE
  environmentId: env-prod-secondary
  clusterType: ENTERPRISE
  displayName: Production Secondary Kafka (us-west-2)
  networkConfig:
    networkId: n-us-west-2
```

**Disaster Recovery Strategy:**
1. Use [Cluster Linking](https://docs.confluent.io/cloud/current/multi-cloud/cluster-linking/index.html) to replicate data from primary to secondary cluster
2. Configure applications to fail over to secondary cluster in case of disaster
3. Both clusters should be in separate regions for true regional fault tolerance

**Key characteristics:**
- **Availability:** Regional fault tolerance
- **RPO:** Near-zero (with Cluster Linking)
- **RTO:** Depends on failover automation
- **Cost:** Runs two full clusters (higher cost for critical workloads)

---

## Additional Configuration Examples

### Using Legacy Availability Values (Basic Clusters)

For backward compatibility with older BASIC cluster configurations:

```yaml
apiVersion: confluent.openmcf.org/v1
kind: ConfluentKafka
metadata:
  name: legacy-basic-kafka
spec:
  cloud: GCP
  region: us-central1
  availability: LOW  # Legacy value for BASIC clusters
  environmentId: env-dev
  clusterType: BASIC
```

**Note:** `LOW` and `HIGH` are legacy availability values supported for BASIC clusters. For new deployments, use `SINGLE_ZONE` or `MULTI_ZONE`.

---

### Custom Display Name

Specify a human-readable display name different from the metadata name:

```yaml
apiVersion: confluent.openmcf.org/v1
kind: ConfluentKafka
metadata:
  name: kafka-prod-01  # Used for infrastructure naming
spec:
  cloud: AWS
  region: us-east-2
  availability: MULTI_ZONE
  environmentId: env-prod
  clusterType: STANDARD
  displayName: "Production Kafka Cluster - Orders & Payments"  # Shown in UI
```

---

## Best Practices

### 1. **Environment Isolation**
- Use separate Confluent Cloud environments for dev, staging, and production
- Provision clusters in the appropriate environment via `environment_id`

### 2. **Naming Conventions**
- Use consistent naming patterns: `{env}-{purpose}-kafka` (e.g., `prod-orders-kafka`)
- Add descriptive labels for cost allocation and filtering

### 3. **High Availability**
- Always use `MULTI_ZONE` availability for production workloads
- This enables 99.99% SLA and automatic fault tolerance

### 4. **Private Networking**
- Use ENTERPRISE or DEDICATED cluster types with `network_config` for production
- Ensure PrivateLink/Private Link endpoints are configured in your VPCs

### 5. **Cluster Sizing**
- Start with STANDARD for most production workloads (elastic scaling)
- Use DEDICATED with provisioned CKU only for predictable, high-throughput workloads
- Monitor usage and adjust cluster type or CKU as needed

### 6. **Cost Optimization**
- Use BASIC clusters for dev/test environments
- Use STANDARD for variable workloads (auto-scales and optimizes cost)
- Use DEDICATED for stable, high-volume workloads where provisioned capacity is more economical

---

## Next Steps

- Review the [README](./README.md) for detailed field descriptions
- Check the [Pulumi implementation](./iac/pulumi/README.md) for advanced customization
- Explore the [Terraform implementation](./iac/tf/README.md) for Terraform-based deployments
- Learn about [Confluent Cloud networking](https://docs.confluent.io/cloud/current/networking/overview.html)
- Set up [Schema Registry](https://docs.confluent.io/cloud/current/sr/index.html) for data governance
