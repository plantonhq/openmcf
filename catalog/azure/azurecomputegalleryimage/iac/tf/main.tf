# Create the gallery image definition. Almost the whole definition is
# create-only in the provider (only description, disk-type exclusions,
# recommended sizing, release notes, end-of-life, and tags update in
# place). The four security flags are a mutual-exclusion clique whose
# ConflictsWith fires on argument PRESENCE -- each is sent ONLY when
# true (an explicit false alongside another flag is provider-rejected).
# Unset architecture/hyper_v_generation ride the provider defaults
# (x64 / V1).
resource "azurerm_shared_image" "main" {
  name                = var.spec.name
  gallery_name        = var.spec.gallery_name
  resource_group_name = var.spec.resource_group
  location            = var.spec.region
  os_type             = var.spec.os_type

  identifier {
    publisher = var.spec.identifier.publisher
    offer     = var.spec.identifier.offer
    sku       = var.spec.identifier.sku
  }

  specialized        = var.spec.specialized
  architecture       = var.spec.architecture != "" ? var.spec.architecture : null
  hyper_v_generation = var.spec.hyper_v_generation != "" ? var.spec.hyper_v_generation : null

  trusted_launch_supported  = var.spec.trusted_launch_supported ? true : null
  trusted_launch_enabled    = var.spec.trusted_launch_enabled ? true : null
  confidential_vm_supported = var.spec.confidential_vm_supported ? true : null
  confidential_vm_enabled   = var.spec.confidential_vm_enabled ? true : null

  accelerated_network_support_enabled = var.spec.accelerated_network_support_enabled ? true : null
  hibernation_enabled                 = var.spec.hibernation_enabled ? true : null
  disk_controller_type_nvme_enabled   = var.spec.disk_controller_type_nvme_enabled ? true : null

  disk_types_not_allowed = length(var.spec.disk_types_not_allowed) > 0 ? var.spec.disk_types_not_allowed : null

  # Updatable, but CLEARING a previously set date forces replacement
  # (the provider's CustomizeDiff).
  end_of_life_date      = var.spec.end_of_life_date != "" ? var.spec.end_of_life_date : null
  eula                  = var.spec.eula != "" ? var.spec.eula : null
  privacy_statement_uri = var.spec.privacy_statement_uri != "" ? var.spec.privacy_statement_uri : null
  release_note_uri      = var.spec.release_note_uri != "" ? var.spec.release_note_uri : null
  description           = var.spec.description != "" ? var.spec.description : null

  dynamic "purchase_plan" {
    for_each = var.spec.purchase_plan != null ? [var.spec.purchase_plan] : []
    content {
      name      = purchase_plan.value.name
      publisher = purchase_plan.value.publisher != "" ? purchase_plan.value.publisher : null
      product   = purchase_plan.value.product != "" ? purchase_plan.value.product : null
    }
  }

  min_recommended_vcpu_count   = var.spec.min_recommended_vcpu_count
  max_recommended_vcpu_count   = var.spec.max_recommended_vcpu_count
  min_recommended_memory_in_gb = var.spec.min_recommended_memory_in_gb
  max_recommended_memory_in_gb = var.spec.max_recommended_memory_in_gb

  tags = local.final_tags
}

# Publish the image's versions -- one ARM child per entry, keyed by the
# version name (its address segment under the image). Each version has
# exactly one source (the spec's CEL enforces it); target regions and
# exclude-from-latest update in place, everything else is create-only.
# A target region's storage_account_type cannot be UPDATED by the API
# and the provider cannot force replacement for it (region-list
# membership changes in place) -- changing it on an existing region
# surfaces Azure's own error.
resource "azurerm_shared_image_version" "main" {
  for_each = local.versions_by_name

  name                = each.value.name
  gallery_name        = azurerm_shared_image.main.gallery_name
  image_name          = azurerm_shared_image.main.name
  resource_group_name = var.spec.resource_group
  location            = var.spec.region

  dynamic "target_region" {
    for_each = each.value.target_regions
    content {
      name                        = target_region.value.name
      regional_replica_count     = target_region.value.regional_replica_count
      disk_encryption_set_id     = target_region.value.disk_encryption_set_id != "" ? target_region.value.disk_encryption_set_id : null
      exclude_from_latest_enabled = target_region.value.exclude_from_latest_enabled
      storage_account_type       = target_region.value.storage_account_type != "" ? target_region.value.storage_account_type : null
    }
  }

  blob_uri            = each.value.blob_uri != "" ? each.value.blob_uri : null
  storage_account_id  = each.value.storage_account_id != "" ? each.value.storage_account_id : null
  os_disk_snapshot_id = each.value.os_disk_snapshot_id != "" ? each.value.os_disk_snapshot_id : null
  managed_image_id    = each.value.managed_image_id != "" ? each.value.managed_image_id : null

  replication_mode                         = each.value.replication_mode != "" ? each.value.replication_mode : null
  exclude_from_latest                      = each.value.exclude_from_latest
  deletion_of_replicated_locations_enabled = each.value.deletion_of_replicated_locations_enabled

  # Updatable, but CLEARING a previously set date forces replacement
  # (the provider's CustomizeDiff).
  end_of_life_date = each.value.end_of_life_date != "" ? each.value.end_of_life_date : null

  tags = merge(local.metadata_tags, each.value.tags)
}
