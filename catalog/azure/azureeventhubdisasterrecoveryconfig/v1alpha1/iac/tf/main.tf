# The geo-DR pairing: metadata (hubs, consumer groups, authorization
# rules -- not event data) continuously replicates primary -> partner,
# and the alias DNS name fronts whichever namespace is currently primary.
#
# Provider-managed lifecycle choreography worth knowing (all inside
# azurerm, no module-side steps): create waits for the pairing to reach
# the Succeeded provisioning state (polling Accepted -> Succeeded);
# changing partner_namespace_id BREAKS the existing pairing first, waits,
# then re-pairs to the new partner; destroy breaks the pairing, deletes
# the config, then waits BOTH for the config to 404 AND for the alias
# NAME to be released by Azure's name-availability check -- the alias
# name stays reserved briefly after deletion, so destroys take minutes by
# the service's own design. Failover itself is an operational action
# performed from the SECONDARY side (portal/CLI/SDK), never a config
# change on this resource.
resource "azurerm_eventhub_namespace_disaster_recovery_config" "main" {
  # ForceNew trio: the alias identity and the primary side (addressed by
  # discrete names parsed from the primary's ARM ID) are fixed at
  # creation.
  name                = var.spec.alias_name
  namespace_name      = local.primary_namespace_name
  resource_group_name = local.primary_resource_group_name

  partner_namespace_id = var.spec.partner_namespace_id
}
