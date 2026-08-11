# Overview

The **AzurePrivateDnsResolver** component deploys Azure DNS Private Resolver -- the managed DNS proxy that makes names resolve ACROSS the hybrid boundary without anyone running DNS server VMs. On-premises forwarders send queries to its inbound endpoint and get Azure's private answers (private zones, private endpoints, VM records); Azure workloads send queries for on-premises domains out through its outbound endpoint, steered by forwarding rulesets.

## Purpose

- **Retire the DNS forwarder VMs**: the two IaaS boxes most hybrid networks run purely to proxy DNS become one managed, zone-resilient service.
- **Both directions, one anchor**: inbound endpoints answer queries coming INTO Azure; outbound endpoints carry queries OUT -- deploy either or both.
- **The hub's DNS heart**: in hub-and-spoke, one resolver in the hub serves every spoke through ruleset links -- spokes need no peering to the resolver's network for DNS to work.

## Key Features

- Full azurerm v5 surface: the resolver plus its inbound and outbound endpoints (up to 5 each way), static or dynamic inbound IPs, composed as one component.
- The provider's Static/Dynamic address contract is validated in seconds -- an invalid combination never reaches Azure.
- Chart-ready: references the virtual network and endpoint subnets by typed outputs; publishes the primary inbound IP (the "point your DNS here" value) and the outbound endpoint id forwarding rulesets bind.

## Use Cases

- **On-premises resolves Azure names**: point datacenter DNS conditional forwarders at the inbound endpoint IP -- private endpoints and private zones resolve correctly from the office.
- **Azure resolves on-premises names**: an outbound endpoint plus a forwarding ruleset (AzurePrivateDnsResolverForwardingRuleset) sends corp-domain queries to the datacenter's DNS servers.
- **The full hybrid loop**: both endpoints on one resolver -- each side resolves the other's names through one managed service.

## Future Enhancements

- The endpoint subnets' delegation contract (`Microsoft.Network/dnsResolvers`, /28-/24, dedicated) stays documentation until references can be introspected offline.
