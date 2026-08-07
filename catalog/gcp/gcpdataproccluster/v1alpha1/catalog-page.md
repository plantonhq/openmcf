# GCP Dataproc Cluster

Deploys a Google Cloud Dataproc cluster for Apache Spark, Hadoop, Hive, and related frameworks — on dedicated Compute Engine VMs (the standard arm) or as Kubernetes pods on an existing GKE cluster (the Dataproc-on-GKE virtual arm). Supports HA masters, Spot secondary workers with machine-type flexibility, autoscaling by policy reference, Kerberos or personal-cluster authentication, CMEK, metastore attachment, OSS metric collection, and lifecycle-driven auto-delete.

## What Gets Created

When you deploy a GcpDataprocCluster resource, Planton provisions:

- **Dataproc Cluster** — a `google_dataproc_cluster` resource on one of the two arms
- **GCE arm** — master nodes, primary workers, optional secondary (Spot/preemptible) workers, and optional dedicated DRIVER node groups
- **Virtual arm** — Dataproc control-plane and Spark pods scheduled onto the referenced GKE cluster's node pools
- **GCS staging/temp buckets** — auto-created by GCP when not referenced
- **Component Gateway endpoints** — authenticated HTTPS URLs for Spark UI, YARN, Jupyter, and other web UIs (when enabled)

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **VPC network or subnetwork** if specifying custom networking (otherwise GCP uses the default network)
- **A service account** with the Dataproc Worker role if using a custom node identity
- **A Cloud KMS key** if enabling CMEK or Kerberos (Kerberos secret files are KMS-encrypted GCS objects)
- **A GKE cluster and node pools** (`GcpGkeCluster`, `GcpGkeNodePool`) if deploying the virtual arm

## Quick Start

Create a file `dataproc.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpDataprocCluster
metadata:
  name: my-spark-cluster
spec:
  projectId:
    value: my-gcp-project
  region: us-central1
  clusterName: my-spark-cluster
  clusterConfig:
    masterConfig:
      machineType: n2-standard-4
    workerConfig:
      numInstances: 2
      machineType: n2-standard-4
    softwareConfig:
      imageVersion: "2.2-debian12"
    endpointConfig:
      enableHttpPortAccess: true
    lifecycleConfig:
      idleDeleteTtl: "1800s"
```

Deploy:

```shell
planton apply -f dataproc.yaml
```

This creates a cluster with 1 master, 2 workers, Spark 3.5, Component Gateway enabled, and auto-delete after 30 minutes idle.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | GCP region for the cluster (e.g., `us-central1`, `europe-west12`). Immutable. | Required, `^[a-z]+-[a-z]+[0-9]+$` |
| `clusterName` | `string` | Cluster name. Lowercase letters, numbers, hyphens. Immutable. | 2-55 chars, `^[a-z][a-z0-9-]{0,53}[a-z0-9]$` |

### Spec-Level Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `projectId` | `StringValueOrRef` | provider default | GCP project. Can reference GcpProject via `valueFrom`. |
| `gracefulDecommissionTimeout` | `string` | `"0s"` | YARN drain window on worker scale-down (e.g., `"3600s"`). GCE arm only. |
| `labels` | `map<string,string>` | `{}` | User labels merged beneath platform labels. Rejected on the virtual arm (API limitation). |
| `clusterConfig` | `object` | — | The Compute Engine arm. Mutually exclusive with `virtualClusterConfig`. |
| `virtualClusterConfig` | `object` | — | The Dataproc-on-GKE arm. Mutually exclusive with `clusterConfig`. Omitting both arms yields a default GCE cluster. |

### GCE Arm (`clusterConfig`)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `stagingBucket` / `tempBucket` | `StringValueOrRef` | auto-created | GCS buckets for dependencies and shuffle/spill. Can reference GcpGcsBucket. |
| `clusterTier` | `string` | `CLUSTER_TIER_STANDARD` | `CLUSTER_TIER_STANDARD` or `CLUSTER_TIER_PREMIUM`. Immutable. |
| `gceConfig.network` | `StringValueOrRef` | default VPC | VPC network. Mutually exclusive with `subnetwork`. Can reference GcpVpcNetwork. |
| `gceConfig.subnetwork` | `StringValueOrRef` | — | VPC subnetwork. Mutually exclusive with `network`. Can reference GcpSubnetwork. |
| `gceConfig.serviceAccount` | `StringValueOrRef` | default CE SA | Node identity. Can reference GcpServiceAccount. |
| `gceConfig.serviceAccountScopes` | `string[]` | GCP default | OAuth scopes for the node SA. |
| `gceConfig.zone` | `string` | auto-selected | Zone within the region for node placement. |
| `gceConfig.internalIpOnly` | `bool` | `false` | Nodes get no external IPs; pair with Cloud NAT or Private Google Access. |
| `gceConfig.tags` / `gceConfig.metadata` | — | — | GCE network tags and instance metadata. |
| `gceConfig.shieldedInstanceConfig` | `object` | off | Secure boot, vTPM, integrity monitoring. |
| `gceConfig.confidentialInstanceConfig` | `object` | off | Confidential VMs (requires N2D machine types). |
| `gceConfig.reservationAffinity` | `object` | ANY | `NO_RESERVATION` / `ANY_RESERVATION` / `SPECIFIC_RESERVATION` (requires `key` + `values`). |
| `gceConfig.nodeGroupAffinity.nodeGroupUri` | `string` | — | Sole-tenant node group placement. |
| `masterConfig` | `object` | 1 master | `numInstances` (1 or 3 for HA), `machineType`, `diskConfig`, `accelerators`, `minCpuPlatform`, `imageUri`. |
| `workerConfig` | `object` | 2 workers | Same shape plus `minNumInstances` (autoscaler floor). `numInstances` and `minNumInstances` update in place. |
| `secondaryWorkerConfig` | `object` | none | `numInstances` (in place), `preemptibility` (`SPOT`/`PREEMPTIBLE`/`NON_PREEMPTIBLE`), `diskConfig`, `instanceFlexibilityPolicy`. |
| `secondaryWorkerConfig.instanceFlexibilityPolicy.instanceSelectionList` | `object[]` | — | Ranked machine-type preferences (`machineTypes` min 1, `rank` >= 0; lower rank preferred). |
| `secondaryWorkerConfig.instanceFlexibilityPolicy.provisioningModelMix` | `object` | — | `standardCapacityBase` + `standardCapacityPercentAboveBase` (0-100) — the on-demand/spot blend. |
| `*.diskConfig` | `object` | GCP defaults | `bootDiskSizeGb` (>=10), `bootDiskType` (`pd-standard`/`pd-ssd`/`pd-balanced`), `numLocalSsds`, `localSsdInterface` (`scsi`/`nvme`). |
| `softwareConfig` | `object` | latest stable | `imageVersion`, `optionalComponents` (JUPYTER, DOCKER, TRINO, ZEPPELIN, FLINK, …), `properties` (`"prefix:property"` overrides). |
| `initializationActions` | `object[]` | `[]` | Startup scripts (`script` GCS URI, optional `timeoutSec`). |
| `autoscalingPolicyUri` | `StringValueOrRef` | — | Autoscaling policy attachment. Can reference GcpDataprocAutoscalingPolicy (resolves to `status.outputs.name`). Attach/swap/detach updates in place. |
| `encryptionKmsKeyName` | `StringValueOrRef` | Google-managed | CMEK for persistent disks. Can reference GcpKmsKey. |
| `securityConfig` | `object` | — | Exactly one of `kerberosConfig` XOR `identityConfig`. |
| `securityConfig.kerberosConfig` | `object` | — | Hadoop Secure Mode. `rootPrincipalPasswordUri` + `kmsKeyUri` required; every secret field is a GCS URI of a KMS-encrypted file. |
| `securityConfig.identityConfig.userServiceAccountMapping` | `map` | — | Personal-cluster auth: user → service account (min 1 pair). |
| `endpointConfig.enableHttpPortAccess` | `bool` | `false` | Component Gateway web UIs. |
| `lifecycleConfig.idleDeleteTtl` / `autoDeleteTime` | `string` | — | Idle-delete duration (`"1800s"`) / absolute RFC3339 deadline. Both update in place. |
| `metastoreConfig.dataprocMetastoreService` | `StringValueOrRef` | — | Persistent Hive metastore resource name (`projects/{p}/locations/{l}/services/{s}`). |
| `dataprocMetricConfig.metrics` | `object[]` | — | Min 1. `metricSource` in `MONITORING_AGENT_DEFAULTS`/`HDFS`/`SPARK`/`YARN`/`SPARK_HISTORY_SERVER`/`HIVESERVER2`; optional `metricOverrides`. |
| `auxiliaryNodeGroups` | `object[]` | `[]` | Dedicated `DRIVER` groups: `roles` (DRIVER only), optional `nodeGroupConfig` sizing, optional `nodeGroupId` (3-33 chars). |

### Virtual Arm (`virtualClusterConfig`)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `stagingBucket` | `StringValueOrRef` | auto-created | GCS staging bucket. |
| `kubernetesClusterConfig` | `object` | — (required) | The Kubernetes-side configuration. |
| `kubernetesClusterConfig.kubernetesNamespace` | `StringValueOrRef` | derived | Namespace for Spark pods. Can reference KubernetesNamespace. |
| `kubernetesClusterConfig.gkeClusterConfig.gkeClusterTarget` | `StringValueOrRef` | — (required) | Target GKE cluster. Can reference GcpGkeCluster (resolves to `status.outputs.cluster_id`). |
| `kubernetesClusterConfig.gkeClusterConfig.nodePoolTarget[]` | `object[]` | Dataproc-managed pool | `nodePool` (required; reference GcpGkeNodePool) + `roles` (min 1 of `DEFAULT`/`CONTROLLER`/`SPARK_DRIVER`/`SPARK_EXECUTOR`) + optional `nodePoolConfig`. |
| `...nodePoolConfig` | `object` | — | For pools Dataproc creates: `locations` (min 1), `autoscaling` (max >= min), `machineType`, `localSsdCount`, `minCpuPlatform`, `preemptible` XOR `spot`. Ignored for pre-existing pools. |
| `kubernetesClusterConfig.kubernetesSoftwareConfig` | `object` | — (required) | `componentVersion` (min 1 pair; SPARK required, e.g. `"3.5-dataproc-17"`) + `properties`. |
| `auxiliaryServicesConfig.metastoreConfig` | `object` | — | Persistent Hive metastore for the virtual cluster's jobs. |
| `auxiliaryServicesConfig.sparkHistoryServerConfig.dataprocCluster` | `StringValueOrRef` | — | The Dataproc cluster hosting a persistent Spark History Server (another cluster's `cluster_id` output). |

## Examples

### HA Production Cluster (GCE arm)

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpDataprocCluster
metadata:
  name: prod-spark
spec:
  projectId:
    value: my-gcp-project
  region: us-central1
  clusterName: prod-spark
  gracefulDecommissionTimeout: "3600s"
  labels:
    team: data-platform
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
      internalIpOnly: true
      shieldedInstanceConfig:
        enableSecureBoot: true
        enableVtpm: true
        enableIntegrityMonitoring: true
    masterConfig:
      numInstances: 3
      machineType: n2-standard-8
      diskConfig:
        bootDiskSizeGb: 200
        bootDiskType: pd-ssd
    workerConfig:
      numInstances: 5
      minNumInstances: 3
      machineType: n2-standard-8
      diskConfig:
        bootDiskSizeGb: 500
        bootDiskType: pd-ssd
        numLocalSsds: 2
        localSsdInterface: nvme
    softwareConfig:
      imageVersion: "2.2-debian12"
    encryptionKmsKeyName:
      valueFrom:
        kind: GcpKmsKey
        name: dataproc-cmek
        fieldPath: status.outputs.key_id
    endpointConfig:
      enableHttpPortAccess: true
```

### Autoscaled Batch Cluster with Spot Secondaries

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpDataprocCluster
metadata:
  name: batch-spark
spec:
  projectId:
    value: my-gcp-project
  region: us-central1
  clusterName: batch-spark
  clusterConfig:
    workerConfig:
      numInstances: 2
      machineType: n2-standard-4
    secondaryWorkerConfig:
      numInstances: 10
      preemptibility: SPOT
      instanceFlexibilityPolicy:
        instanceSelectionList:
          - machineTypes:
              - n2-standard-4
              - n2d-standard-4
            rank: 0
        provisioningModelMix:
          standardCapacityBase: 2
          standardCapacityPercentAboveBase: 0
    autoscalingPolicyUri:
      valueFrom:
        kind: GcpDataprocAutoscalingPolicy
        name: batch-autoscaling
        fieldPath: status.outputs.name
    lifecycleConfig:
      idleDeleteTtl: "900s"
```

### Spark on GKE (virtual arm)

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpDataprocCluster
metadata:
  name: spark-on-gke
spec:
  projectId:
    value: my-gcp-project
  region: us-central1
  clusterName: spark-on-gke
  virtualClusterConfig:
    stagingBucket:
      valueFrom:
        kind: GcpGcsBucket
        name: dataproc-staging
        fieldPath: status.outputs.bucket_id
    kubernetesClusterConfig:
      gkeClusterConfig:
        gkeClusterTarget:
          valueFrom:
            kind: GcpGkeCluster
            name: platform-gke
            fieldPath: status.outputs.cluster_id
        nodePoolTarget:
          - nodePool:
              valueFrom:
                kind: GcpGkeNodePool
                name: dataproc-pool
                fieldPath: status.outputs.node_pool_id
            roles:
              - DEFAULT
      kubernetesSoftwareConfig:
        componentVersion:
          SPARK: "3.5-dataproc-17"
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `cluster_id` | `string` | Fully qualified cluster resource name (`projects/{project}/regions/{region}/clusters/{cluster}`) — the composition handle, including another cluster's `sparkHistoryServerConfig` |
| `cluster_name` | `string` | Short cluster name (same as `spec.clusterName`) |
| `staging_bucket` | `string` | GCS bucket used for staging (user-supplied or auto-created), on either arm |

## Related Components

- [GcpDataprocAutoscalingPolicy](/docs/catalog/gcp/gcpdataprocautoscalingpolicy) — the shared autoscaling policy referenced by `autoscalingPolicyUri`
- [GcpGcsBucket](/docs/catalog/gcp/gcpgcsbucket) — staging and temp buckets for job artifacts
- [GcpVpcNetwork](/docs/catalog/gcp/gcpvpcnetwork) — VPC network for cluster node placement
- [GcpSubnetwork](/docs/catalog/gcp/gcpsubnetwork) — subnetwork for controlled IP range allocation
- [GcpServiceAccount](/docs/catalog/gcp/gcpserviceaccount) — custom IAM identity for cluster VMs
- [GcpKmsKey](/docs/catalog/gcp/gcpkmskey) — CMEK and Kerberos file encryption keys
- [GcpGkeCluster](/docs/catalog/gcp/gcpgkecluster) — the virtual arm's compute substrate
- [GcpGkeNodePool](/docs/catalog/gcp/gcpgkenodepool) — node pools Dataproc roles schedule onto
