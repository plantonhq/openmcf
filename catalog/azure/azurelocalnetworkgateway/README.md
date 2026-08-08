# Overview

The **Azure Local Network Gateway API Resource** provides a consistent and standardized interface for deploying and managing local network gateways -- Azure's description of the ON-PREMISES side of a site-to-site VPN: the VPN device's public endpoint and the address space reachable behind it. Nothing deploys on-premises; the resource is the address-book entry a gateway connection points at.

## Purpose

We developed this API resource so every on-premises site is a first-class, named object in the resource graph. One site, one description, referenced by every tunnel that reaches it:

- **Site descriptions**: the device's public IP (or re-resolved FQDN) plus the CIDR ranges behind it
- **Static or dynamic routing**: route the declared address spaces, or let the site's BGP speaker advertise them
- **Multi-site topologies**: a hub gateway terminates one connection per site, each pointing at its own description

## Key Features

- **Consistent Interface**: aligns with our existing APIs for deploying cloud infrastructure across multiple providers
- **Endpoint Flexibility**: exactly one of a static public IPv4 address or an FQDN (for sites whose address changes -- Azure re-resolves periodically), enforced at validation time
- **Routing-Source Contract**: ARM's requirement that a site carries static prefixes, BGP settings, or both is spec-enforced upfront
- **Zero Cost**: the object is pure ARM metadata -- it provisions in seconds and costs nothing to keep

## Use Cases

- **Datacenter description**: the HQ firewall's public IP and the RFC 1918 ranges behind it
- **Branch fleet**: one description per branch office, each consumed by its own tunnel to a shared hub gateway
- **Dynamic sites**: FQDN-addressed sites on consumer connections with changing IPs
- **BGP-routed sites**: descriptions carrying the on-premises AS number and tunnel-interior peering address instead of static prefixes

## Future Enhancements

Future updates will include:

- **Site inventory views**: console surfacing of every described site and the tunnels consuming it
- **Reachability insights**: pairing descriptions with connection health for per-site status

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
