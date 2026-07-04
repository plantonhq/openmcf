# Create the container registry -- the managed, private OCI registry the
# platform's workloads pull their images from.
#
# Lifecycle notes worth knowing before operating this resource:
# - Name and region are the registry's identity -- changing either
#   replaces it and its CONTENTS DO NOT MIGRATE; every image would need
#   re-pushing. Zone redundancy and CMK encryption are likewise fixed at
#   creation.
# - The SKU changes in place, but downgrading requires every Premium-only
#   feature (geo-replication, network rules, policies, CMK) to be unset
#   first -- the same ordering ARM enforces.
# - Geo-replications are managed inline by azurerm and can be added and
#   removed in place; each replica is its own tracked ARM resource with
#   its own tags.
resource "azurerm_container_registry" "main" {
  name                = var.spec.registry_name
  resource_group_name = var.spec.resource_group
  location            = var.spec.region

  # The STANDARD baseline is materialized here when the spec leaves the
  # SKU unspecified (azurerm requires an explicit value); mapped in locals.
  sku = local.sku

  admin_enabled = var.spec.admin_user_enabled

  public_network_access_enabled = var.spec.public_network_access_enabled

  # Premium-only; ForceNew. Spec-level validation enforces the SKU gate,
  # mirroring ARM's own CustomizeDiff.
  zone_redundancy_enabled = var.spec.zone_redundancy_enabled

  anonymous_pull_enabled    = var.spec.anonymous_pull_enabled
  data_endpoint_enabled     = var.spec.data_endpoint_enabled
  quarantine_policy_enabled = var.spec.quarantine_policy_enabled

  # Unset keeps untagged manifests forever (Azure's default); a value
  # (including 0 = purge immediately) turns the retention policy on.
  retention_policy_in_days = var.spec.retention_policy_in_days

  trust_policy_enabled  = var.spec.trust_policy_enabled
  export_policy_enabled = var.spec.export_policy_enabled

  # null lets Azure apply its default (AzureServices); mapped in locals.
  network_rule_bypass_option = local.network_rule_bypass_option

  # The public-registry allowlist (Premium). ARM only supports Allow
  # rules, so the action is constant and only the CIDR ranges vary.
  dynamic "network_rule_set" {
    for_each = var.spec.network_rule_set != null ? [var.spec.network_rule_set] : []
    content {
      default_action = local.network_rule_default_action
      dynamic "ip_rule" {
        for_each = network_rule_set.value.ip_rules
        content {
          action   = "Allow"
          ip_range = ip_rule.value.ip_range
        }
      }
    }
  }

  # Additional regions the registry replicates to (Premium). azurerm
  # expects the list in alphabetical location order; iterating a
  # location-keyed map yields exactly that, keeping manifests
  # order-insensitive instead of surfacing ARM's quirk to users.
  dynamic "georeplications" {
    for_each = { for g in var.spec.georeplications : g.location => g }
    content {
      location                  = georeplications.value.location
      zone_redundancy_enabled   = georeplications.value.zone_redundancy_enabled
      regional_endpoint_enabled = georeplications.value.regional_endpoint_enabled
      tags                      = georeplications.value.tags
    }
  }

  # The registry's managed identity; a USER_ASSIGNED identity is what
  # unwraps the CMK encryption key at boot.
  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type
      identity_ids = identity.value.identity_ids
    }
  }

  # Customer-managed-key encryption (Premium; fixed at creation). The
  # identity_client_id must belong to an identity listed in the identity
  # block that holds get/wrapKey/unwrapKey on the key's vault.
  dynamic "encryption" {
    for_each = var.spec.encryption != null ? [var.spec.encryption] : []
    content {
      identity_client_id = encryption.value.identity_client_id
      key_vault_key_id   = encryption.value.key_vault_key_id
    }
  }

  tags = local.final_tags
}
