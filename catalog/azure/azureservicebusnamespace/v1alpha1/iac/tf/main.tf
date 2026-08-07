# The Service Bus namespace: the container and billing boundary for
# enterprise messaging. Queues, topics, subscriptions, SAS rules, and the
# geo-DR pairing are all first-class kinds that reference this
# namespace's ARM id -- nothing is bundled here, so entity teams own
# their entities independently of the namespace's owner.
resource "azurerm_servicebus_namespace" "main" {
  # ForceNew: the name is the public DNS identity
  # ({name}.servicebus.windows.net); changing it replaces the namespace
  # and every entity in it.
  name                = var.spec.namespace_name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  sku = local.sku

  # PREMIUM pairings: the spec's CELs mirror Azure's create-time
  # contract (capacity {1,2,4,8,16} and partitions {1,2,4} required on
  # PREMIUM, forbidden otherwise), so the module sends them only when
  # present. Partitions are ForceNew -- the layout is fixed at creation.
  capacity                     = var.spec.capacity
  premium_messaging_partitions = var.spec.premium_messaging_partitions

  # Managed identity -- required for CMK (the unwrapping identity must
  # be attached here), usable anywhere the namespace itself
  # authenticates to other services.
  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_map[identity.value.type]
      identity_ids = length(identity.value.user_assigned_identity_ids) > 0 ? identity.value.user_assigned_identity_ids : null
    }
  }

  # Customer-managed-key encryption (PREMIUM only). Azure cannot remove
  # CMK once set -- dropping this block forces namespace replacement
  # (the provider encodes that as a ForceNew diff).
  dynamic "customer_managed_key" {
    for_each = var.spec.customer_managed_key != null ? [var.spec.customer_managed_key] : []
    content {
      key_vault_key_id                  = customer_managed_key.value.key_vault_key_id
      identity_id                       = customer_managed_key.value.user_assigned_identity_id
      infrastructure_encryption_enabled = customer_managed_key.value.infrastructure_encryption_enabled
    }
  }

  # False = keyless posture: every SAS rule's keys (including the root
  # rule surfaced in this module's outputs) stop being usable
  # credentials, and clients authenticate with Entra identities.
  local_auth_enabled = var.spec.local_auth_enabled

  public_network_access_enabled = var.spec.public_network_access_enabled

  # The namespace firewall (PREMIUM only). The provider models it as an
  # inline block riding a SEPARATE ARM operation after the namespace
  # create; Azure rejects DENY with no admitted sources (front-loaded as
  # a spec CEL). The block-level public_network_access_enabled must
  # agree with the namespace-level one -- the module wires the block's
  # dial independently because Azure validates the pair server-side.
  dynamic "network_rule_set" {
    for_each = var.spec.network_rule_set != null ? [var.spec.network_rule_set] : []
    content {
      default_action = (
        network_rule_set.value.default_action != null && network_rule_set.value.default_action != ""
        ? local.network_default_action_map[network_rule_set.value.default_action]
        : "Allow"
      )
      public_network_access_enabled = network_rule_set.value.public_network_access_enabled
      trusted_services_allowed      = network_rule_set.value.trusted_services_allowed
      ip_rules                      = network_rule_set.value.ip_rules

      dynamic "network_rules" {
        for_each = network_rule_set.value.network_rules
        content {
          subnet_id                            = network_rules.value.subnet_id
          ignore_missing_vnet_service_endpoint = network_rules.value.ignore_missing_vnet_service_endpoint
        }
      }
    }
  }

  tags = local.final_tags
}
