# Create the Azure AI Search service. Capacity is sku x partitions x
# replicas; the spec's CEL contracts already enforce the provider's
# per-SKU caps and pairing rules (high-density on standard3 only,
# failure mode only with local auth, semantic not on free), so by the
# time this module runs the shape is legal.
#
# The SKU changes in place ONLY along basic -> standard -> standard2 ->
# standard3 (the provider's update contract) -- every other SKU change
# replaces the service. The service name is GLOBALLY unique (it forms
# the endpoint DNS name).
resource "azurerm_search_service" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group
  sku                 = var.spec.sku

  # Null (unset) lets the provider apply its default, 1 -- the same
  # value the spec's proto default carries.
  replica_count   = var.spec.replica_count
  partition_count = var.spec.partition_count

  # Enum name -> wire value; unspecified ("") omits the property so the
  # provider applies its default, "default". ForceNew.
  hosting_mode = lookup(local.hosting_mode_wire, var.spec.hosting_mode, null)

  # Null (unset) lets the provider apply its default, true. Setting
  # false is the RBAC-only posture -- admin/query keys stop working.
  local_authentication_enabled = var.spec.local_authentication_enabled

  # Setting a failure mode is what enables RBAC alongside API keys;
  # omitted, the service stays in API-keys-only mode.
  authentication_failure_mode = var.spec.authentication_failure_mode != "" ? var.spec.authentication_failure_mode : null

  # Provider default false -- an explicit false is the same wire.
  customer_managed_key_enforcement_enabled = var.spec.customer_managed_key_enforcement_enabled

  public_network_access_enabled = var.spec.public_network_access_enabled

  # Omitted, the provider sends "disabled" -- semantic ranking off.
  semantic_search_sku = var.spec.semantic_search_sku != "" ? var.spec.semantic_search_sku : null

  allowed_ips = length(var.spec.allowed_ips) > 0 ? var.spec.allowed_ips : null

  # Unspecified ("") omits the property so the provider applies its
  # default, "None".
  network_rule_bypass_option = var.spec.network_rule_bypass_option != "" ? var.spec.network_rule_bypass_option : null

  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_wire[identity.value.type]
      identity_ids = length(identity.value.identity_ids) > 0 ? identity.value.identity_ids : null
    }
  }

  tags = local.final_tags
}

# The composed shared private links: standalone ARM children
# (.../sharedPrivateLinkResources/{name}), one per spec entry, keyed by
# name (uniqueness is spec CEL). Each link sits "Pending" until the
# target resource's owner approves it -- creating the link never
# requires the target side's consent.
resource "azurerm_search_shared_private_link_service" "links" {
  for_each = { for link in var.spec.shared_private_link_services : link.name => link }

  name               = each.value.name
  search_service_id  = azurerm_search_service.main.id
  subresource_name   = each.value.subresource_name
  target_resource_id = each.value.target_resource_id
  request_message    = each.value.request_message != "" ? each.value.request_message : null
}
