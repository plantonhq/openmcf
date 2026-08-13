# Create the Bastion host -- the managed jump service that opens
# RDP/SSH sessions to virtual machines over their private addresses.
# Dedicated-infrastructure SKUs (Basic/Standard/Premium) deploy into
# the network's "AzureBastionSubnet" (/26+, the exact ARM name -- ARM
# validates it at deploy time) with a Standard static public IP the
# host binds EXCLUSIVELY; the Developer SKU instead attaches to a
# virtual network on Azure-shared infrastructure. The host bills
# hourly per SKU from the moment it provisions (~10 minutes); SKU
# upgrades are in-place, downgrades replace the host (provider
# CustomizeDiff).
resource "azurerm_bastion_host" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group
  sku                 = local.sku

  # Fixed at 2 on Developer/Basic (the provider default, always sent);
  # 2-50 on Standard/Premium (spec-validated).
  scale_units = coalesce(var.spec.scale_units, 2)

  # Available on every SKU; the provider defaults it true.
  copy_paste_enabled = coalesce(var.spec.copy_paste_enabled, true)

  # Standard/Premium feature knobs (spec-validated against the SKU).
  # kerberos_enabled is applied at CREATE only -- the provider has no
  # update path for it and silently ignores later changes.
  file_copy_enabled         = var.spec.file_copy_enabled
  ip_connect_enabled        = var.spec.ip_connect_enabled
  kerberos_enabled          = var.spec.kerberos_enabled
  shareable_link_enabled    = var.spec.shareable_link_enabled
  tunneling_enabled         = var.spec.tunneling_enabled
  session_recording_enabled = var.spec.session_recording_enabled

  # Developer SKU only (spec-validated): the shared-infrastructure host
  # attaches to a virtual network directly -- no subnet, no public IP.
  virtual_network_id = var.spec.virtual_network_id != "" ? var.spec.virtual_network_id : null

  # The dedicated-infrastructure binding. Premium may omit the public
  # IP to deploy private-only (surfaced in the private_only_enabled
  # output).
  dynamic "ip_configuration" {
    for_each = var.spec.ip_configuration != null ? [var.spec.ip_configuration] : []
    content {
      name                 = ip_configuration.value.name
      subnet_id            = ip_configuration.value.subnet_id
      public_ip_address_id = ip_configuration.value.public_ip_address_id != "" ? ip_configuration.value.public_ip_address_id : null
    }
  }

  # Availability zone pinning; empty deploys regionally. Fixed at
  # creation.
  zones = length(var.spec.zones) > 0 ? var.spec.zones : null

  tags = local.final_tags
}
