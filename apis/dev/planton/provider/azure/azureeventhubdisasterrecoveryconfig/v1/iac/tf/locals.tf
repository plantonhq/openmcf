locals {
  # Geo-DR configs carry no Azure tags: ARM does not support tags on
  # disasterRecoveryConfigs, so the platform's identity tags live on the
  # paired namespaces.

  # azurerm addresses the PRIMARY side by discrete names (namespace name
  # + resource group) rather than by ARM ID, so parse the resolved
  # primary_namespace_id into those parts. The anchored regex fails the
  # plan loudly on a malformed id instead of sending garbage to the API.
  primary_namespace_parts = regex(
    "/resourceGroups/(?P<rg>[^/]+)/providers/Microsoft.EventHub/namespaces/(?P<ns>[^/]+)$",
    var.spec.primary_namespace_id
  )
  primary_resource_group_name = local.primary_namespace_parts.rg
  primary_namespace_name      = local.primary_namespace_parts.ns
}
