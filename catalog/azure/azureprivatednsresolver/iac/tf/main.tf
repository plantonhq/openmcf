# Create the DNS Private Resolver -- the managed DNS proxy that
# resolves names across the hybrid boundary -- and its inbound/outbound
# endpoints as child resources. The resolver anchors to ONE virtual
# network (Azure allows at most one resolver per network) and each
# endpoint occupies its own dedicated subnet delegated to
# "Microsoft.Network/dnsResolvers" (/28-/24, nothing else in it) -- ARM
# validates the delegation and network membership at deploy time.
# Everything except tags is create-only; endpoints bill hourly, the
# resolver object is free.
resource "azurerm_private_dns_resolver" "main" {
  name                = var.spec.name
  resource_group_name = var.spec.resource_group
  location            = var.spec.region
  virtual_network_id  = var.spec.virtual_network_id
  tags                = local.final_tags
}

# Inbound endpoints -- the private IPs on-premises DNS forwarders send
# queries to. Keyed by endpoint name (spec-validated unique) so adding
# or removing one endpoint never replaces its siblings. Endpoint
# deletes outlive the API's first "deleted" answer -- the provider
# polls until the endpoint is verifiably gone.
resource "azurerm_private_dns_resolver_inbound_endpoint" "main" {
  for_each = { for endpoint in var.spec.inbound_endpoints : endpoint.name => endpoint }

  name                    = each.value.name
  private_dns_resolver_id = azurerm_private_dns_resolver.main.id
  # Endpoints deploy into the resolver's region (their subnets belong
  # to the resolver's own network).
  location = var.spec.region

  ip_configurations {
    subnet_id                    = each.value.subnet_id
    private_ip_allocation_method = lookup(local.allocation_method_wire, coalesce(each.value.private_ip_allocation_method, "DYNAMIC"), "Dynamic")
    # Only STATIC allocation carries an address (spec-validated); the
    # provider rejects an address on Dynamic before touching ARM.
    private_ip_address = each.value.private_ip_address != "" ? each.value.private_ip_address : null
  }

  tags = local.final_tags
}

# Outbound endpoints -- the egress points queries leave Azure through,
# steered by the forwarding rulesets that bind them. Same keying and
# delete-poller reality as the inbound endpoints.
resource "azurerm_private_dns_resolver_outbound_endpoint" "main" {
  for_each = { for endpoint in var.spec.outbound_endpoints : endpoint.name => endpoint }

  name                    = each.value.name
  private_dns_resolver_id = azurerm_private_dns_resolver.main.id
  location                = var.spec.region
  subnet_id               = each.value.subnet_id
  tags                    = local.final_tags
}
