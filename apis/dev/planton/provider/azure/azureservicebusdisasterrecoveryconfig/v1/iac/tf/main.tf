# The geo-DR pairing: metadata (entities, rules, SAS rules -- not
# message data) continuously replicates primary -> partner, and the
# alias DNS name fronts whichever namespace is currently primary.
#
# Provider-managed lifecycle choreography worth knowing (all inside
# azurerm, no module-side steps): create waits for the pairing to reach
# Succeeded; changing partner_namespace_id breaks the existing pairing,
# waits, then re-pairs; destroy breaks the pairing, deletes the config,
# and polls until the alias NAME is released so an immediate re-create
# with the same alias does not collide. Failover itself is an
# operational action taken on the SECONDARY during an incident -- never
# a config change here.
resource "azurerm_servicebus_namespace_disaster_recovery_config" "main" {
  # ForceNew pair: the alias identity and the primary side are fixed at
  # creation.
  name                 = var.spec.alias_name
  primary_namespace_id = var.spec.primary_namespace_id

  partner_namespace_id = var.spec.partner_namespace_id

  # Unset defaults the alias connection strings to the namespace's
  # built-in root rule; a scoped rule gives least-privilege alias
  # credentials.
  alias_authorization_rule_id = local.alias_authorization_rule_id
}
