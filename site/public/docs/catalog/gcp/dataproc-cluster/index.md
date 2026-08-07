---
title: "Dataproc Cluster"
description: "Dataproc Cluster deployment documentation"
icon: "package"
order: 100
componentName: "gcpdataproccluster"
---

# GCP Dataproc Cluster

Deploys a Dataproc cluster for Apache Spark, Hadoop, and related data processing frameworks — either the standard Compute Engine arm (configurable master and worker nodes, preemptible/spot secondary workers, software image versioning, optional components like Jupyter/Flink/Trino, in-cluster security, lifecycle management for cost control, and CMEK encryption) or the Dataproc-on-GKE virtual arm, where Spark workloads run as pods on an existing GKE cluster. The component integrates with Planton's Provider Connections for GCP credential management and supports ValueFromRef wiring to GCP projects, VPCs, subnets, service accounts, GCS buckets, KMS keys, autoscaling policies, and GKE clusters/node pools.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Dataproc Cluster** -- a managed cluster resource in the specified GCP project and region, configured with the chosen Dataproc image version, master/worker topology, and software components
- **Master Nodes** -- 1 instance for standard mode or 3 instances for high-availability mode, with configurable machine type, boot disk, local SSDs, and accelerators
- **Primary Worker Nodes** -- on-demand VMs for persistent compute capacity, with configurable machine type, disk, accelerators, and autoscaling minimum
- **Secondary Worker Nodes** -- created only when `clusterConfig.secondaryWorkerConfig` is set; preemptible or spot VMs providing cost-optimized burst capacity
- **Software Configuration** -- Dataproc image version, optional components (Jupyter, Flink, Presto, Trino), and framework property overrides for Spark, Hadoop, YARN, and HDFS
- **Initialization Actions** -- created only when `clusterConfig.initializationActions` is set; startup scripts that run on all nodes during cluster creation
- **Component Gateway** -- created only when `clusterConfig.endpointConfig.enableHttpPortAccess` is true; provides authenticated HTTPS access to Spark UI, YARN, HDFS, and optional component web interfaces
- **Lifecycle Configuration** -- created only when `clusterConfig.lifecycleConfig` is set; automatic cluster deletion after idle timeout or at a scheduled time
- **CMEK Encryption** -- created only when `clusterConfig.encryptionKmsKeyName` is provided; encrypts persistent disks with a customer-managed Cloud KMS key
- **In-Cluster Security** -- created only when `clusterConfig.securityConfig` is set; exactly one of Kerberos (Hadoop Secure Mode) or personal-cluster identity mapping
- **Autoscaling Attachment** -- created only when `clusterConfig.autoscalingPolicyUri` references a GcpDataprocAutoscalingPolicy; attaching, swapping, or detaching updates in place
- **Virtual Cluster (Dataproc-on-GKE)** -- created only when `virtualClusterConfig` is set instead of `clusterConfig`; registers Spark workloads as pods on an existing GKE cluster with node-pool role mapping and shared metastore/Spark History Server attachments
- **GCP Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking and governance; user labels from `spec.labels` propagate to cluster VMs (not supported by the API on virtual clusters)

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### GCP Project

- **A GCP project** where the cluster will be created. Provide the project ID directly or reference a GcpProject Cloud Resource via ValueFromRef.
- **Dataproc API** and **Compute Engine API** enabled in the target project.
- **A VPC network or subnetwork** for cluster node placement. Using a subnetwork is recommended for production clusters with controlled IP ranges. Provide directly or reference GcpVpcNetwork/GcpSubnetwork Cloud Resources via ValueFromRef.
- **A custom service account** (recommended for production) with minimal permissions for cluster VMs. The default Compute Engine service account works for development.
- **Cloud NAT or Private Google Access** (if using `internalIpOnly: true`) so private nodes can reach the internet for container image pulls and package installs.

## Deploy

### Console

Open the deployment store, find **GCP Dataproc Cluster**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Dev Jupyter** preset in the [Presets](#presets) tab to pre-populate a development cluster with Jupyter notebooks.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpDataprocCluster
metadata:
  name: analytics-spark
  org: acme-corp
  env: prod
spec:
  projectId:
    value: "acme-prod-12345"
  region: us-central1
  clusterName: analytics-spark
```

```shell
planton apply -f dataproc-cluster.yaml
```

This creates a standard cluster with GCP defaults: 1 master and 2 workers on n2-standard-4, 500 GB pd-standard disks, latest Dataproc image, and no lifecycle management. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the cluster to VPC infrastructure deployed in the same InfraPipeline:

```yaml
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: production-project
      fieldPath: status.outputs.project_id
  clusterConfig:
    gceConfig:
      subnetwork:
        valueFrom:
          kind: GcpSubnetwork
          name: dataproc-subnet
          fieldPath: status.outputs.subnetwork_self_link
      serviceAccount:
        valueFrom:
          kind: GcpServiceAccount
          name: dataproc-sa
          fieldPath: status.outputs.email
    stagingBucket:
      valueFrom:
        kind: GcpGcsBucket
        name: spark-staging
        fieldPath: status.outputs.bucket_id
    encryptionKmsKeyName:
      valueFrom:
        kind: GcpKmsKey
        name: dataproc-key
        fieldPath: status.outputs.key_id
```

The InfraPipeline resolves the dependency graph, deploys the project, subnet, service account, bucket, and KMS key first, then provisions the Dataproc cluster with all resolved values.

## Key Configuration

These are the most important decisions when configuring a Dataproc cluster. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Cluster topology** -- Set `clusterConfig.masterConfig.numInstances` to 1 for standard mode or 3 for high-availability mode. HA mode tolerates a single master failure but triples master costs. Worker count is set via `clusterConfig.workerConfig.numInstances` (default 2).

**Secondary workers for cost optimization** -- Configure `clusterConfig.secondaryWorkerConfig` with `preemptibility: SPOT` for cost-optimized burst capacity. Spot VMs can be preempted at any time, so use them for fault-tolerant batch workloads with Spark's dynamic allocation enabled.

**Lifecycle management** -- Set `clusterConfig.lifecycleConfig.idleDeleteTtl` (e.g., `"1800s"` for 30 minutes) to auto-delete idle clusters. Critical for ephemeral batch clusters to avoid runaway costs. Use `autoDeleteTime` for time-boxed clusters with a known end date.

**Software and components** -- Set `clusterConfig.softwareConfig.imageVersion` (e.g., `"2.2-debian12"`) and add optional components like `JUPYTER`, `FLINK`, or `PRESTO` via `optionalComponents`. Override framework properties via the `properties` map (e.g., `"spark:spark.executor.memory": "12g"`).

**Network isolation** -- Set `clusterConfig.gceConfig.internalIpOnly: true` for private nodes with no external IPs. Requires Cloud NAT or Private Google Access for internet connectivity. Combine with `tags` for firewall rule targeting.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** | `projectId` | `status.outputs.project_id` |
| **GcpVpcNetwork** (optional) | `clusterConfig.gceConfig.network` | `status.outputs.network_self_link` |
| **GcpSubnetwork** (optional) | `clusterConfig.gceConfig.subnetwork` | `status.outputs.subnetwork_self_link` |
| **GcpServiceAccount** (optional) | `clusterConfig.gceConfig.serviceAccount` | `status.outputs.email` |
| **GcpGcsBucket** (optional) | `clusterConfig.stagingBucket` | `status.outputs.bucket_id` |
| **GcpGcsBucket** (optional) | `clusterConfig.tempBucket` | `status.outputs.bucket_id` |
| **GcpKmsKey** (optional) | `clusterConfig.encryptionKmsKeyName` | `status.outputs.key_id` |
| **GcpDataprocAutoscalingPolicy** (optional) | `clusterConfig.autoscalingPolicyUri` | `status.outputs.name` |
| **GcpKmsKey** (Kerberos only) | `clusterConfig.securityConfig.kerberosConfig.kmsKeyUri` | `status.outputs.key_id` |
| **GcpGkeCluster** (virtual arm) | `virtualClusterConfig.kubernetesClusterConfig.gkeClusterConfig.gkeClusterTarget` | `status.outputs.cluster_id` |
| **GcpGkeNodePool** (virtual arm) | `virtualClusterConfig...nodePoolTarget[].nodePool` | `status.outputs.node_pool_id` |
| **GcpGcsBucket** (virtual arm) | `virtualClusterConfig.stagingBucket` | `status.outputs.bucket_id` |
| **GcpDataprocCluster** (virtual arm) | `virtualClusterConfig.auxiliaryServicesConfig.sparkHistoryServerConfig.dataprocCluster` | `status.outputs.cluster_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `cluster_id` | Fully qualified cluster resource name (`projects/{p}/regions/{r}/clusters/{c}`) | Dataproc job submissions, workflow template references, another cluster's Spark History Server attachment |
| `cluster_name` | Short cluster name | Display, logging, job targeting |
| `staging_bucket` | GCS bucket used for staging job dependencies | Job dependency uploads, output inspection |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Dev with Jupyter** -- Single master, 2 workers on e2-standard-4, Jupyter notebooks enabled, Component Gateway for web UI access, and 30-minute idle auto-delete. Optimized for interactive data exploration and prototyping. Start from the **Dev Jupyter** preset.

**HA production** -- 3 masters for high availability, 5 workers with SSD disks and local SSDs, internal-only networking, CMEK encryption, custom service account, and graceful decommission timeout. Suitable for production ETL pipelines and long-running Spark applications. Start from the **HA Production** preset.

**Cost-optimized batch** -- Single master with 2 primary workers and 10 spot secondary workers, dynamic allocation enabled, and 15-minute idle auto-delete. Suitable for ephemeral batch ETL jobs where cost efficiency outweighs preemption risk. Start from the **Cost-Optimized Batch** preset.

**Spark on GKE** -- A Dataproc-on-GKE virtual cluster: Spark workloads run as pods on an existing GKE cluster referenced by ValueFromRef, sharing its capacity, autoscaling, and operational tooling. Start from the **Spark on GKE** preset.

## Works With

- [**GCP Project**](/cloud-catalog/gcp-project) -- provides the GCP project where the cluster is created
- [**GCP VPC Network**](/cloud-catalog/gcp-vpc-network) -- provides the VPC network for cluster node placement
- [**GCP Subnetwork**](/cloud-catalog/gcp-subnetwork) -- provides the subnet with controlled IP ranges for cluster nodes
- [**GCP Service Account**](/cloud-catalog/gcp-service-account) -- provides the identity for cluster VMs
- [**GCP GCS Bucket**](/cloud-catalog/gcp-gcs-bucket) -- provides staging and temp buckets for job dependencies and shuffle data
- [**GCP KMS Key**](/cloud-catalog/gcp-kms-key) -- provides the CMEK encryption key for cluster persistent disks
- [**GCP Dataproc Autoscaling Policy**](/cloud-catalog/gcp-dataproc-autoscaling-policy) -- provides the reusable worker-scaling contract attached via `autoscalingPolicyUri`
- [**GCP GKE Cluster**](/cloud-catalog/gcp-gke-cluster) -- hosts the virtual arm's Spark pods
- [**GCP GKE Node Pool**](/cloud-catalog/gcp-gke-node-pool) -- provides the node pools the virtual arm schedules Dataproc roles onto