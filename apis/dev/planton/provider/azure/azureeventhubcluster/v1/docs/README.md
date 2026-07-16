# AzureEventHubCluster -- Design Research

## The Resource

A dedicated Event Hubs cluster (`Microsoft.EventHub/clusters`) is the
single-tenant top of the Event Hubs capacity ladder: guaranteed,
isolated throughput sold in capacity units. The component maps onto
`azurerm_eventhub_cluster` (azurerm v4.x,
`internal/services/eventhub/eventhub_cluster_resource.go`),
parity-verified against pulumi-azure v6 (`eventhub.Cluster`).

Many namespaces share one cluster (joined via the namespace's
`dedicated_cluster_id`, ForceNew on the namespace side), which is why
the cluster is its own kind rather than a namespace property. Placement
buys the namespaces: up to 1024 partitions per hub, 90-day retention,
and customer-managed-key eligibility
(AzureEventHubNamespaceCustomerManagedKey requires dedicated or
PREMIUM).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `location` | `region` | Required, ForceNew |
| `resource_group_name` | `resource_group` | FK → AzureResourceGroup; ForceNew |
| `name` | `cluster_name` | Required, ForceNew; 1-50 chars, alphanumeric ends, `[-._]` interior |
| `sku_name` | `capacity_units` | Composed (see below); scaling updates in place; unset deploys 1 CU |
| `tags` | `tags` | Merged over Planton identity tags; user values win |

Outputs: `cluster_id` (what namespaces reference), `cluster_name`.

## The Composed Sku

ARM's `sku_name` is the composite string `Dedicated_{n}` -- tier name
plus capacity count in one field. Dedicated is the ONLY sku family
Azure sells for clusters, so the tier name is a one-value constant, not
configuration: the spec models the actually-variable part
(`capacity_units`, min 1) and both modules compose the string from it.
Surfacing the raw sku string would invite invalid values for zero
expressive gain.

## The 4-Hour Deletion Moratorium

Azure FORBIDS deleting a cluster for 4 hours after creation. A destroy
inside that window retries against `ClusterMoratoriumInEffect` until
Azure permits the delete -- a destroy of a young cluster takes hours by
the service's own rule, billing dedicated capacity units the whole
time. The spec, both modules, and the presets all carry the warning.

## E2E Exclusion

Live E2E is excluded for this kind (see `e2e/profile.yaml`): the
moratorium makes an ephemeral create-verify-destroy cycle impossible
inside any reasonable window, at dedicated-tier cost throughout. The
component stands on its offline gate -- full parity audit, a complete
plan rendering the composed sku, spec tests covering every validation
rule, and the outputs conformance case. A live proof would need a
deliberate long-window run against a cluster older than 4 hours.

## Operational Behavior Worth Knowing

- **Per-CU-per-hour billing at dedicated-tier rates** -- the most
  expensive resource in the Event Hubs family; provision deliberately.
- **`capacity_units` scales in place** -- start at 1, grow with load.
- **Namespace placement is ForceNew on the namespace** -- moving a
  namespace onto or off a cluster replaces the namespace, so plan
  placement up front.

## Composition

- `resource_group` → `AzureResourceGroup.status.outputs.resource_group_name`
- `cluster_id` output ← `AzureEventHubNamespace.spec.dedicated_cluster_id`
  (placement) and, transitively, the namespaces' CMK eligibility
