# The Event Hubs namespace: the container and billing boundary for
# high-throughput event streaming. Hubs, consumer groups, SAS rules,
# schema groups, the geo-DR pairing, and CMK encryption are all
# first-class kinds that reference this namespace's ARM id -- nothing is
# bundled here, so stream teams own their entities independently of the
# namespace's owner.
resource "azurerm_eventhub_namespace" "main" {
  # ForceNew: the name is the public DNS identity
  # ({name}.servicebus.windows.net -- Event Hubs shares the Service Bus
  # DNS zone) and the Kafka bootstrap host; changing it replaces the
  # namespace and every entity in it.
  name                = var.spec.namespace_name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # BASIC <-> STANDARD updates in place; moving into or out of PREMIUM
  # replaces the namespace (the provider encodes the boundary as a
  # ForceNew diff -- Azure cannot convert across the reserved/
  # multi-tenant boundary).
  sku = local.sku

  # Throughput units (BASIC/STANDARD) or processing units (PREMIUM).
  # Sent only when present so Azure's default (1) applies otherwise.
  capacity = var.spec.capacity

  # STANDARD's elastic scaling: Azure grows TUs up to the ceiling under
  # load but never shrinks them back -- scale-down is a manual capacity
  # edit. Azure validates the ceiling/enable pairing at apply time; the
  # provider's only guard is zeroing the ceiling on a Basic downgrade,
  # which it performs itself.
  auto_inflate_enabled     = var.spec.auto_inflate_enabled
  maximum_throughput_units = var.spec.maximum_throughput_units

  # ForceNew: a namespace cannot move on or off a dedicated cluster in
  # place. Placement buys single-tenant capacity, 1024-partition hubs,
  # 90-day retention, and CMK eligibility.
  dedicated_cluster_id = var.spec.dedicated_cluster_id

  # Managed identity -- required for identity-based capture auth and for
  # CMK (the unwrapping identity must be attached here), usable anywhere
  # the namespace itself authenticates to other services.
  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_map[identity.value.type]
      identity_ids = length(identity.value.user_assigned_identity_ids) > 0 ? identity.value.user_assigned_identity_ids : null
    }
  }

  # False = keyless posture: every SAS rule's keys (including the root
  # rule surfaced in this module's outputs) stop being usable
  # credentials, and clients authenticate with Entra identities.
  local_authentication_enabled = var.spec.local_authentication_enabled

  public_network_access_enabled = var.spec.public_network_access_enabled

  # The namespace firewall (not available on BASIC -- front-loaded as a
  # spec CEL). The provider models it as an inline block riding a
  # SEPARATE ARM operation after the namespace create. The block-level
  # public_network_access_enabled must agree with the namespace-level
  # one (Azure validates the pair server-side; the spec CEL front-loads
  # it). Azure requires an explicit default_action when the block is
  # declared.
  dynamic "network_rulesets" {
    for_each = var.spec.network_rule_sets != null ? [var.spec.network_rule_sets] : []
    content {
      default_action                 = local.network_default_action_map[network_rulesets.value.default_action]
      public_network_access_enabled  = network_rulesets.value.public_network_access_enabled
      trusted_service_access_enabled = network_rulesets.value.trusted_service_access_enabled

      # Each entry is an allow rule: Azure's per-rule action accepts
      # exactly one value (Allow), so the spec models just the mask.
      dynamic "ip_rule" {
        for_each = network_rulesets.value.ip_rules
        content {
          ip_mask = ip_rule.value
          action  = "Allow"
        }
      }

      dynamic "virtual_network_rule" {
        for_each = network_rulesets.value.virtual_network_rules
        content {
          subnet_id                                       = virtual_network_rule.value.subnet_id
          ignore_missing_virtual_network_service_endpoint = virtual_network_rule.value.ignore_missing_virtual_network_service_endpoint
        }
      }
    }
  }

  tags = local.final_tags
}
