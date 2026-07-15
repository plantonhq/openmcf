# Azure Event Hub Cluster

Creates a dedicated Event Hubs cluster -- single-tenant capacity units of guaranteed throughput that namespaces join for 1024-partition hubs, 90-day retention, and customer-managed-key eligibility.

## What Gets Created

When you deploy an AzureEventHubCluster resource, Planton provisions:

- **Event Hubs Cluster** -- an `azurerm_eventhub_cluster` with the composed `Dedicated_{n}` sku, sized by your capacity-unit count

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An Azure Resource Group** to create the cluster in (referenced through `resourceGroup`)

## Quick Start

Create a file `cluster.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureEventHubCluster
metadata:
  name: streaming-cluster
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureEventHubCluster.streaming-cluster
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: streaming-rg
      fieldPath: status.outputs.resource_group_name
  clusterName: myorg-streaming-dedicated
  capacityUnits: 1
```

Deploy:

```shell
planton apply -f cluster.yaml
```

Provision deliberately: clusters bill per capacity unit per hour at dedicated-tier rates -- the most expensive resource in the Event Hubs family. `capacityUnits` scales in place (unset deploys 1 CU, the entry size). Azure forbids deleting a cluster for 4 hours after creation; a destroy inside that window retries until Azure permits the delete.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `cluster_id` | What an AzureEventHubNamespace's `dedicatedClusterId` references to place the namespace on this cluster |
| `cluster_name` | The cluster's name |

## Related Resources

- [Azure Resource Group](/docs/catalog/azure/azureresourcegroup) -- the container
- [Azure Event Hub Namespace](/docs/catalog/azure/azureeventhubnamespace) -- namespaces placed on the cluster via `dedicatedClusterId`
- [Azure Event Hub Namespace Customer Managed Key](/docs/catalog/azure/azureeventhubnamespacecustomermanagedkey) -- BYOK encryption, which dedicated placement unlocks
