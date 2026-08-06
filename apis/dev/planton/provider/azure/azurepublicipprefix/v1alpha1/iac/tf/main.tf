# Create the public IP prefix -- a reserved, contiguous range of public IP
# addresses that public IPs allocate from and NAT gateways associate for
# outbound SNAT.
#
# Lifecycle notes worth knowing before operating this resource:
# - The prefix is essentially immutable: everything except tags is fixed
#   at creation, and replacing it changes the ACTUAL reserved range --
#   everything partners have allowlisted. Treat replacement as a
#   coordinated migration, never a casual update.
# - The prefix cannot be deleted while any of its addresses are in use by
#   public IPs or NAT gateway associations.
resource "azurerm_public_ip_prefix" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # Azure's default length is 28 (16 addresses); null lets the provider
  # apply it, so an unspecified spec and Azure's default deploy identically
  # on both engines. The same only-send-explicit rule applies to the
  # version/SKU/tier enums (mapped in locals).
  prefix_length = var.spec.prefix_length
  ip_version    = local.ip_version
  sku           = local.sku
  sku_tier      = local.sku_tier

  zones = length(var.spec.zones) > 0 ? var.spec.zones : null

  custom_ip_prefix_id = var.spec.custom_ip_prefix_id

  tags = local.final_tags
}
