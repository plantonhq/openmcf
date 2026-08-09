# Overview

The **Azure VPN Gateway Connection API Resource** provides a consistent and standardized interface for deploying and managing VPN Gateway Connections -- the tunnel bundles joining branch sites to a Virtual WAN hub's VPN gateway. One connection per branch; one tunnel (`vpnLinks` entry) per site link being connected.

## Purpose

We developed this API resource so the branch tunnel -- with its per-link IPsec, BGP, and NAT choices, and the hub routing that shapes what the branch reaches -- is one first-class, versioned object:

- **One tunnel per site link**: each `vpnLinks` entry pins to a site link by its ARM ID (the site's name-keyed `link_ids` output) -- the active-active branch pattern, typed
- **Secrets referenced, never embedded**: per-link `sharedKey` is a sensitive reference; omit it and Azure generates one
- **Topology as configuration**: the routing block reuses the hub's route-table and route-map vocabulary

## Key Features

- **Consistent Interface**: aligns with our existing APIs for deploying cloud infrastructure across multiple providers
- **Provider Contracts Upfront**: a configured routing block must name its association; policy-based selectors require a pinned IPsec proposal
- **Provisioned-vs-Connected Honesty**: deployment success and tunnel establishment are different events -- documented where operators look

## Use Cases

- **Connect a branch** described by a VPN Site to a hub's VPN gateway
- **Dual-ISP active-active branches**: two tunnels with per-link parameters
- **Compliance-pinned IPsec**: exact cipher suites per tunnel
- **Overlapping branch address spaces**: NAT rule opt-ins per tunnel direction

## Future Enhancements

Future updates will include:

- **Tunnel health surfaces**: connection state and traffic telemetry integration

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
