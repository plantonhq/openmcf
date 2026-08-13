# Create the Virtual WAN -- the free, lightweight umbrella object of
# Azure's managed hub-and-spoke networking. Hubs (and the gateways on
# them) are separate resources that reference this WAN's ID; ARM refuses
# to delete a WAN that still has hubs.
resource "azurerm_virtual_wan" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  disable_vpn_encryption = var.spec.disable_vpn_encryption

  # ARM defaults branch-to-branch transit ON -- most of the point of a
  # Virtual WAN; unset (null) mirrors it, matching the Pulumi module's
  # nil handling.
  allow_branch_to_branch_traffic = (
    var.spec.allow_branch_to_branch_traffic == null
    ? true
    : var.spec.allow_branch_to_branch_traffic
  )

  office365_local_breakout_category = local.office365_local_breakout_category

  # "Standard" is the full-mesh tier and ARM's default; "Basic" is the
  # constrained legacy tier (upgradeable in place, never downgradeable).
  type = (
    var.spec.type == null || var.spec.type == ""
    ? "Standard"
    : var.spec.type
  )

  tags = local.final_tags
}
