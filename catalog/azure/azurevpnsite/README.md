# Overview

The **Azure VPN Site API Resource** provides a consistent and standardized interface for deploying and managing VPN Sites -- the Virtual WAN address-book entries describing branch locations: their internet links, reachable address space, and terminating device. A site deploys nothing at the branch and costs nothing to keep; it is what a VPN Gateway Connection points at.

## Purpose

We developed this API resource so describing a branch -- and the links a connection tunnels to -- is one first-class, versioned object:

- **The branch, described once**: links (with per-link BGP), reachable prefixes, device metadata, O365 breakout
- **Connectable by name**: each link's ARM ID surfaces in the name-keyed `link_ids` output -- exactly what a connection's tunnels pin to
- **Chart-ready wiring**: the WAN is a typed reference; the site's ARM ID surfaces where connections consume it

## Key Features

- **Consistent Interface**: aligns with our existing APIs for deploying cloud infrastructure across multiple providers
- **Endpoint Contracts**: every link must carry a public endpoint (IP or FQDN) -- ARM's create-time rule enforced upfront
- **Active-Active Ready**: two links (two ISPs) model the dual-tunnel branch pattern the connection kind completes

## Use Cases

- **Describe a datacenter or branch office** ahead of connecting it to a hub's VPN gateway
- **Dual-ISP branches**: two links with distinct endpoints, connected active-active
- **BGP branches**: per-link BGP speakers instead of static prefix lists
- **SD-WAN O365 breakout**: declare which O365 categories exit the branch locally

## Future Enhancements

Future updates will include:

- **SD-WAN partner automation**: surfacing partner-managed site properties

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
