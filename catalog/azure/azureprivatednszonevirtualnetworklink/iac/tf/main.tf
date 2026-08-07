# Create the virtual network link -- the attachment that makes the parent
# private DNS zone resolvable from the referenced virtual network.
#
# Lifecycle notes worth knowing before operating this resource:
# - registration_enabled, resolution_policy, and tags update IN PLACE;
#   name, zone, and network are the link's ARM identity, so changing any
#   of them replaces the link (a brief resolution gap for the affected
#   network, nothing else).
# - Azure allows only ONE registration-enabled link per virtual network;
#   additional links to the same network must keep VM auto-registration
#   off.
# - The zone's name and resource group are derived from the referenced
#   zone's ARM ID (see locals), so the module never asks the user to
#   restate derivable state that could then disagree with the zone.
resource "azurerm_private_dns_zone_virtual_network_link" "main" {
  name                  = var.spec.name
  resource_group_name   = local.zone_resource_group_name
  private_dns_zone_name = local.zone_name
  virtual_network_id    = var.spec.virtual_network_id

  # Azure defaults VM auto-registration to off; it is only meaningful for
  # custom internal zones (privatelink zones are populated by private
  # endpoints, not VM registration).
  registration_enabled = var.spec.registration_enabled

  # null lets Azure choose its per-zone-type default (privatelink zones get
  # their platform-managed policy); only an explicit policy is ever sent.
  resolution_policy = local.resolution_policy

  tags = local.final_tags
}
