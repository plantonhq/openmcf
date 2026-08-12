# Create the Azure Cosmos DB for MongoDB vCore cluster. The provider
# owns the mode machinery this module deliberately does not
# reimplement: Default mode stages the Data API in a separate
# post-create update, upgrades away from Free/M25 stage a tier-first
# update, and a create_mode change is forced to a replacement (Azure
# never returns the mode on reads).
resource "azurerm_mongo_cluster" "main" {
  name                = var.spec.name
  resource_group_name = var.spec.resource_group
  location            = var.spec.region

  # Platform default "Default" -- always sent, so the rendered plan
  # states the mode.
  create_mode = coalesce(var.spec.create_mode, "Default")

  # The native administrator pair travels together (spec CEL mirrors
  # the provider's RequiredWith); replicas and restores inherit the
  # source's administrator, so both stay unsent when unset.
  administrator_username = var.spec.administrator_username != "" ? var.spec.administrator_username : null
  administrator_password = var.spec.administrator_username != "" ? var.spec.administrator_password : null

  # Sizing fields are sent only when set: replica and restore modes
  # inherit them from the source, and the provider sends each to ARM
  # only when it carries a value.
  version                = var.spec.version
  compute_tier           = var.spec.compute_tier
  storage_size_in_gb     = var.spec.storage_size_in_gb
  shard_count            = var.spec.shard_count
  high_availability_mode = var.spec.high_availability_mode

  # The storage type rides the size (the provider's RequiredWith):
  # platform default "PremiumSSD" when the size is set at all.
  storage_type = var.spec.storage_size_in_gb != null ? coalesce(var.spec.storage_type, "PremiumSSD") : null

  # Sent only when set: Azure defaults an unset list to ["NativeAuth"]
  # server-side (the provider models the argument Optional+Computed for
  # exactly that reason).
  authentication_methods = length(var.spec.authentication_methods) > 0 ? var.spec.authentication_methods : null

  # User-assigned is the only identity flavor the service supports on
  # this resource. Adding the first identity or removing the last one
  # is a REPLACEMENT (the provider forces it -- Azure rejects the
  # in-place transition); documented on the spec field.
  dynamic "identity" {
    for_each = length(var.spec.user_assigned_identity_ids) > 0 ? [var.spec.user_assigned_identity_ids] : []
    content {
      type         = "UserAssigned"
      identity_ids = identity.value
    }
  }

  dynamic "customer_managed_key" {
    for_each = var.spec.customer_managed_key != null ? [var.spec.customer_managed_key] : []
    content {
      key_vault_key_id          = customer_managed_key.value.key_vault_key_id
      user_assigned_identity_id = customer_managed_key.value.user_assigned_identity_id
    }
  }

  preview_features = length(var.spec.preview_features) > 0 ? var.spec.preview_features : null

  # GeoReplica coordinates -- the spec CELs mirror the provider's
  # create-time contract (both required in GeoReplica mode; location
  # requires the server id).
  source_server_id = var.spec.source_server_id != "" ? var.spec.source_server_id : null
  source_location  = var.spec.source_location != "" ? var.spec.source_location : null

  dynamic "restore" {
    for_each = var.spec.restore != null ? [var.spec.restore] : []
    content {
      point_in_time_utc = restore.value.point_in_time_utc
      source_id         = restore.value.source_id
    }
  }

  # Sent only when the manifest carries it: the provider ERRORS when
  # the raw config sets this on a non-Default-mode cluster (even
  # false), and the spec CEL front-loads that contract.
  data_api_mode_enabled = var.spec.data_api_mode_enabled

  # Platform default true, mapped to the provider's Enabled/Disabled
  # tokens -- always sent (mirrors Azure's own default).
  public_network_access = coalesce(var.spec.public_network_access_enabled, true) ? "Enabled" : "Disabled"

  tags = local.final_tags
}

# Composed firewall rules -- one provider resource per named rule,
# keyed by the rule's name (renames replace only that rule, sibling
# rules stay untouched), in lockstep with the Pulumi module's per-name
# resources. Start/end addresses update in place.
resource "azurerm_mongo_cluster_firewall_rule" "main" {
  for_each = local.firewall_rules_by_name

  name             = each.value.name
  mongo_cluster_id = azurerm_mongo_cluster.main.id
  start_ip_address = each.value.start_ip_address
  end_ip_address   = each.value.end_ip_address
}
