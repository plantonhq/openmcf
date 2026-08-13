# Azure Machine Learning Compute Cluster

Creates an auto-scaling compute cluster on an Azure Machine Learning workspace -- the pool of VMs that training jobs and pipelines run on, growing and shrinking between its configured node bounds. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Compute Cluster** -- an ARM child of the workspace (`.../workspaces/{ws}/computes/{name}`) with its VM size, priority, scale bounds, identity, and networking

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **An AzureMachineLearningWorkspace** -- the cluster is created on it.

### Azure Subscription

- **Regional vCPU quota for the VM family** -- the cluster provisions nodes only up to your quota; GPU families usually need a quota request first.
- **Update surface is narrow** -- only identity, scale settings, and tags update in place; everything else replaces the cluster (running jobs on it fail).

## Deploy

### Console

Open the deployment store, find **Azure Machine Learning Compute Cluster**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **CPU Training Cluster** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureMachineLearningComputeCluster
metadata:
  name: cpu-cluster
  org: acme-corp
  env: prod
spec:
  workspaceId:
    valueFrom:
      name: ml-prod
  name: cpu-cluster
  region: eastus
  vmSize: STANDARD_DS3_V2
  vmPriority: DEDICATED
  scaleSettings:
    minNodeCount: 0
    maxNodeCount: 4
    scaleDownNodesAfterIdleDuration: PT30M
  identity:
    type: SYSTEM_ASSIGNED
```

```shell
planton apply -f azure-machine-learning-compute-cluster.yaml
```

The cluster object creates in a few minutes; nodes provision on demand as jobs arrive.

### InfraChart

In an ML-platform chart the order is: workspace → **compute cluster**, wiring the workspace by reference; jobs then reference the cluster by its name as their compute target.

## Key Configuration

These are the most important decisions when configuring the cluster. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Scale bounds** -- `minNodeCount: 0` is the economical posture (free while idle, node spin-up latency when work arrives); a non-zero minimum keeps warm nodes for latency-sensitive teams. The idle duration (`PT30M` style) decides how long a node survives between jobs.

**VM priority** -- `DEDICATED` nodes stay up until scaled down; `LOW_PRIORITY` (spot-class) nodes cost a fraction but can be evicted mid-job -- suit checkpointed, fault-tolerant training only.

**Identity** -- give the cluster a managed identity and grant it data access (Storage Blob Data Reader on training data, AcrPull on the registry); jobs then run credential-free.

**Region** -- uniquely in the ML family, the cluster's nodes may run in a different region from the workspace (useful when GPU quota lives elsewhere); ARM still reports the cluster envelope at the workspace's region.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureMachineLearningWorkspace** | `workspaceId` | `status.outputs.machine_learning_workspace_id` |
| **AzureSubnet** | `subnetId` | `status.outputs.subnet_id` |
| **AzureUserAssignedIdentity** | `identity.identityIds[]` | `status.outputs.identity_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `machine_learning_compute_cluster_id` | ARM ID of the cluster | Operational tooling |
| `machine_learning_compute_cluster_name` | The cluster's name | What jobs and pipelines reference as their compute target |
| `system_assigned_identity_principal_id` | The system identity's principal ID | Storage / Key Vault / ACR role assignments |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**CPU training pool** -- a scale-to-zero dedicated cluster for everyday jobs. Start from the **CPU Training Cluster** preset.

**GPU training** -- a GPU-family cluster with a system identity. Start from the **GPU Training Cluster** preset.

**Spot batch** -- a low-priority cluster for checkpointed workloads. Start from the **Low-Priority Batch Cluster** preset.

## Works With

- [**Azure Machine Learning Workspace**](/cloud-catalog/azure-machine-learning-workspace) -- the parent workspace
- [**Azure Machine Learning Datastore**](/cloud-catalog/azure-machine-learning-datastore) -- where the jobs' data lives
- [**Azure Subnet**](/cloud-catalog/azure-subnet) -- VNet placement for the nodes
- [**Azure User Assigned Identity**](/cloud-catalog/azure-user-assigned-identity) -- bring-your-own identity for data and registry grants
