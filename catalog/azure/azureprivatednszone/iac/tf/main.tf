# Create the private DNS zone -- name resolution inside virtual networks
# without running a DNS server.
#
# Lifecycle notes worth knowing before operating this resource:
# - Tags update IN PLACE; the zone's name is its ARM identity, so renaming
#   replaces the zone AND every record in it. The SOA record is written at
#   creation and cannot be customized afterwards.
# - Private DNS zones are global resources -- no location/region argument.
# - The zone is deliberately just the zone: which networks can resolve it
#   is declared through AzurePrivateDnsZoneVirtualNetworkLink resources
#   referencing this zone's zone_id output, one per network. A zone with
#   no links answers nobody.
resource "azurerm_private_dns_zone" "main" {
  name                = var.spec.name
  resource_group_name = var.spec.resource_group

  # Absent means Azure creates its standard SOA record; the block is only
  # sent when the spec customizes it. Unset timers fall back to Azure's
  # defaults so a partially-specified block and Azure's own values deploy
  # identically on both engines.
  dynamic "soa_record" {
    for_each = var.spec.soa_record != null ? [var.spec.soa_record] : []
    content {
      email        = soa_record.value.email
      expire_time  = soa_record.value.expire_time
      minimum_ttl  = soa_record.value.minimum_ttl
      refresh_time = soa_record.value.refresh_time
      retry_time   = soa_record.value.retry_time
      ttl          = soa_record.value.ttl
      tags         = soa_record.value.tags
    }
  }

  tags = local.final_tags
}
