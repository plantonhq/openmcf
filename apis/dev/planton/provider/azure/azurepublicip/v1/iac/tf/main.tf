# Create the public IP address -- a static, internet-routable address that
# load balancers, application gateways, NAT gateways, firewalls, and VMs
# attach for inbound or outbound connectivity.
#
# Lifecycle notes worth knowing before operating this resource:
# - reverse_fqdn, ddos settings, idle timeout, and tags update IN PLACE.
#   Name, SKU/tier, version, zones, prefix membership, ip_tags, and edge
#   zone are fixed at creation -- changing any of them replaces the
#   resource and with it the ACTUAL ADDRESS, so treat replacement as a
#   coordinated migration (DNS, allowlists).
# - Allocation is always Static: dynamic allocation existed only for the
#   Basic SKU, whose creation Azure discontinued in 2025, and every
#   current SKU requires static.
resource "azurerm_public_ip" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group
  allocation_method   = "Static"

  # Enum choices mapped in locals; null lets Azure apply its defaults so an
  # unspecified spec deploys identically on both engines.
  sku        = local.sku
  sku_tier   = local.sku_tier
  ip_version = local.ip_version

  zones = length(var.spec.zones) > 0 ? var.spec.zones : null

  public_ip_prefix_id = var.spec.public_ip_prefix_id

  domain_name_label       = var.spec.domain_name_label
  domain_name_label_scope = local.domain_name_label_scope
  reverse_fqdn            = var.spec.reverse_fqdn

  idle_timeout_in_minutes = var.spec.idle_timeout_in_minutes

  ip_tags = length(var.spec.ip_tags) > 0 ? var.spec.ip_tags : null

  # The plan is only valid alongside the ENABLED mode (spec-level
  # validation enforces the pairing).
  ddos_protection_mode    = local.ddos_protection_mode
  ddos_protection_plan_id = var.spec.ddos_protection_plan_id

  edge_zone = var.spec.edge_zone

  tags = local.final_tags
}
