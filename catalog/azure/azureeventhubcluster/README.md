# AzureEventHubCluster

A dedicated Event Hubs cluster: single-tenant capacity units (CUs) of
guaranteed, isolated throughput at the top of the Event Hubs capacity
ladder, above PREMIUM's shared-infrastructure processing units.
Namespaces are placed on the cluster via their `dedicated_cluster_id`
reference -- many namespaces share one cluster, which is why the
cluster is its own kind rather than a namespace property.

## When to Use

Use AzureEventHubCluster when you need:

- **Guaranteed single-tenant throughput** -- sustained high-volume
  estates that outgrow PREMIUM's processing units
- **1024-partition hubs and 90-day retention** -- limits only
  dedicated placement unlocks for the namespaces on the cluster
- **Customer-managed-key encryption** -- namespace-level CMK
  (AzureEventHubNamespaceCustomerManagedKey) requires the namespace to
  sit on a dedicated cluster or be PREMIUM

Provision one deliberately: dedicated clusters bill per capacity unit
per hour at dedicated-tier rates -- the most expensive resource in the
Event Hubs family.

## Key Configuration

- `region` / `resource_group` -- placement, fixed at creation
- `cluster_name` -- unique within the resource group, 1-50 characters;
  ForceNew (renaming replaces the cluster, subject to the deletion
  moratorium below)
- `capacity_units` -- the cluster's size; each CU is a slice of
  guaranteed single-tenant ingest/egress capacity. Scales in place;
  unset deploys 1 CU, Azure's entry size. The modules compose the ARM
  sku `Dedicated_{n}` from this count -- Dedicated is the only sku
  family Azure sells for clusters, so the tier name is a constant, not
  configuration
- `tags` -- merged over the Planton-derived identity tags (user values
  win)

Know the moratorium: **Azure forbids deleting a cluster for 4 hours
after creation.** A destroy inside that window retries until Azure
permits the delete -- expect a destroy of a young cluster to take
hours by the service's own rule.

## Composition

Namespaces join the cluster at creation (ForceNew on the namespace):

```yaml
# On an AzureEventHubNamespace
dedicatedClusterId:
  valueFrom:
    kind: AzureEventHubCluster
    name: streaming-cluster
    fieldPath: status.outputs.cluster_id
```

## Documentation

- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
