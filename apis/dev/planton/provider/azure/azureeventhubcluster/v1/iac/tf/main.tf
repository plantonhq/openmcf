# The dedicated Event Hubs cluster: single-tenant capacity units that
# namespaces are placed on via their dedicated_cluster_id reference.
# Many namespaces share one cluster, which is why the cluster is its own
# resource rather than a namespace property.
#
# Cost: dedicated clusters bill per capacity unit per hour at
# dedicated-tier rates -- the most expensive resource in the Event Hubs
# family. Provision one deliberately.
#
# Lifecycle: Azure FORBIDS deleting a cluster for 4 HOURS after creation
# (the deletion moratorium). A destroy inside that window retries until
# Azure permits the delete -- expect a destroy of a young cluster to
# take hours by the service's own rule.
resource "azurerm_eventhub_cluster" "main" {
  # ForceNew: renaming replaces the cluster (subject to the 4-hour
  # deletion moratorium above).
  name                = var.spec.cluster_name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # The ARM sku is the composite string "Dedicated_{CUs}". The tier name
  # is a one-value constant -- Dedicated is the ONLY sku family Azure
  # sells for clusters -- so the module composes the string from the
  # capacity count instead of surfacing it as configuration. Unset
  # deploys 1 CU, Azure's entry size; scaling updates in place.
  sku_name = "Dedicated_${coalesce(var.spec.capacity_units, 1)}"

  tags = local.final_tags
}
