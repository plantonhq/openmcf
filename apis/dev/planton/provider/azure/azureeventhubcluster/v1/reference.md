# AzureEventHubCluster

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1`

**AzureEventHubClusterSpec** defines a dedicated Event Hubs cluster:
single-tenant capacity units (CUs) that namespaces are placed on via
their dedicated_cluster_id reference.

Dedicated clusters are the top of the Event Hubs capacity ladder --
above PREMIUM's shared-infrastructure processing units. A cluster buys
guaranteed, isolated throughput; up to 1024 partitions per hub; 90-day
retention; and namespace-level customer-managed-key encryption
(AzureEventHubNamespaceCustomerManagedKey requires the namespace to sit
on a cluster). Many namespaces share one cluster, which is why the
cluster is its own kind rather than a namespace property.

**Cost and lifecycle warnings**:
- Dedicated clusters bill per capacity unit per hour at enterprise
  rates -- this is the most expensive resource in the Event Hubs
  family. Provision one deliberately.
- **Azure forbids deleting a cluster for 4 HOURS after creation** (the
  deletion moratorium). Destroy operations inside that window retry
  until Azure permits the delete -- expect a destroy of a young
  cluster to take hours by the service's own rule.

**ForceNew fields**: `cluster_name`, `region`, `resource_group`.
Capacity (`capacity_units`) scales in place.

## Example

```yaml
# Offline-plan manifest: a dedicated cluster with an explicit capacity
# count and a user tag, exercising the composed Dedicated_{n} sku and
# the tag merge.
apiVersion: azure.planton.dev/v1
kind: AzureEventHubCluster
metadata:
  name: test-eventhub-cluster
spec:
  region: eastus
  resourceGroup:
    value: my-rg
  clusterName: hack-eventhub-cluster
  capacityUnits: 1
  tags:
    team: streaming
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.clusterName` | `string` | yes |  |  |
| `spec.capacityUnits` | `int32` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

The Azure region where the cluster is created. Examples: "eastus",
"westus2", "westeurope". Fixed at creation.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.resourceGroup

`string | valueFrom` · required

The Azure Resource Group the cluster lives in.
Can be a literal string or a reference to an AzureResourceGroup
output. Fixed at creation.

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.clusterName

`string` · required

The cluster's name -- unique within the resource group, 1-50
characters. Starts and ends with a letter or number; letters,
numbers, periods, hyphens, and underscores in between.

**ForceNew**: changing the name replaces the cluster (subject to the
4-hour deletion moratorium).

- rule: cluster_name must start and end with a letter or number and may contain letters, numbers, periods, hyphens, and underscores (max 50 characters)
- rule: {"required":true,"string":{"minLen":"1","maxLen":"50"}}

### spec.capacityUnits

`int32` · optional (explicit presence)

The cluster's size in capacity units (CUs) -- each CU is a slice of
guaranteed, single-tenant ingest/egress capacity. Azure sells the
Dedicated tier only (the modules compose the ARM sku "Dedicated_{n}"
from this count -- the tier name is a one-value constant, not
configuration). Scaling updates in place; 1 CU is the entry size.
Default: 1

- rule: {"int32":{"gte":1}}

### spec.tags

`map<string, string>`

Tags to apply to the cluster, merged over the Planton-derived
metadata tags (user values win on key conflicts). ARM tags are
Azure's first-class governance surface -- Azure Policy enforces them
and Microsoft Cost Management groups by them.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureEventHubCluster, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.cluster_id` | `string` | The Azure Resource Manager ID of the cluster -- what an AzureEventHubNamespace's dedicated_cluster_id references to place the namespace on this cluster. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.EventHub/clusters/{name} |
| `status.outputs.cluster_name` | `string` | The cluster's name. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureEventHubNamespace | `spec.dedicatedClusterId` | `status.outputs.cluster_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
