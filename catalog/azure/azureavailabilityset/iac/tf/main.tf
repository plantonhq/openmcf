# Create the availability set -- the classic pre-zones placement
# grouping that spreads VMs across separate fault and update domains.
# The whole configuration is fixed at creation (only tags update in
# place). Unset optional fields are passed as null so the provider's
# own defaults apply (5 update domains, 3 fault domains, managed=true
# -- managed aligns fault domains with the VMs' managed-disk storage).
resource "azurerm_availability_set" "main" {
  name                = var.spec.name
  resource_group_name = var.spec.resource_group
  location            = var.spec.region

  platform_update_domain_count = var.spec.platform_update_domain_count
  platform_fault_domain_count  = var.spec.platform_fault_domain_count
  managed                      = var.spec.managed

  proximity_placement_group_id = (
    var.spec.proximity_placement_group_id != null && var.spec.proximity_placement_group_id != ""
    ? var.spec.proximity_placement_group_id
    : null
  )

  tags = local.final_tags
}
