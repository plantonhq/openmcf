# The Front Door profile -- deliberately just the container. Endpoints,
# origin groups, origins, and routes are their own components referencing
# this profile's outputs, mirroring Azure's child-resource model. Front
# Door is a global service: the provider forces location to "global",
# so there is no region to configure.
resource "azurerm_cdn_frontdoor_profile" "main" {
  name                = var.spec.profile_name
  resource_group_name = var.spec.resource_group

  # ForceNew, and Azure additionally refuses a Premium -> Standard
  # downgrade outright -- choose Premium deliberately.
  sku_name = local.sku_name

  # Sent only when set: Azure applies its own 120 s default when the
  # property is omitted, matching the spec's documented default.
  response_timeout_seconds = var.spec.response_timeout_seconds

  # The managed identity is how Front Door reads customer-managed TLS
  # certificates from Key Vault without an access-policy secret. The
  # spec's CEL guarantees identity ids are present exactly when the type
  # includes UserAssigned.
  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_map[identity.value.type]
      identity_ids = identity.value.user_assigned_identity_ids
    }
  }

  # Log scrubbing is enabled by the PRESENCE of rules (no rules ==
  # scrubbing disabled). The service supports only the match-everything
  # operator on profile scrubbing rules, so each entry is just the
  # request part to mask.
  dynamic "log_scrubbing_rule" {
    for_each = local.log_scrubbing_variables
    content {
      match_variable = local.log_scrubbing_variable_map[log_scrubbing_rule.value]
    }
  }

  tags = local.final_tags
}
