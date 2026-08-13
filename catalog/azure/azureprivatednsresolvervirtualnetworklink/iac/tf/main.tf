# Create the virtual network link -- the attachment that makes a DNS
# forwarding ruleset take effect in ONE virtual network. The linked
# network must be in the ruleset's region but does NOT need to be
# peered with the resolver's network (hub-and-spoke: spokes link to
# the hub's ruleset and their queries egress through the hub's
# outbound endpoint). One link per ruleset-network pair, up to 500 per
# ruleset; links are free at rest. Everything except metadata is
# create-only.
resource "azurerm_private_dns_resolver_virtual_network_link" "main" {
  name                      = var.spec.name
  dns_forwarding_ruleset_id = var.spec.dns_forwarding_ruleset_id
  virtual_network_id        = var.spec.virtual_network_id

  # ARM's free-form annotation map on the link itself (links carry no
  # tags) -- the only surface updatable in place.
  metadata = length(var.spec.metadata) > 0 ? var.spec.metadata : null
}
