# Overview

The **Azure VPN Server Configuration API Resource** provides a consistent and standardized interface for deploying and managing VPN Server Configurations -- the reusable point-to-site authentication policies (Entra ID, certificate, RADIUS) that point-to-site VPN gateways attach to. The configuration is free, deploys in seconds, and one policy can serve many gateways.

## Purpose

We developed this API resource so "who may connect and how" is one first-class, versioned object instead of settings scattered across gateways:

- **Authentication, declared once**: Entra ID, certificate, and RADIUS parameters with their trust anchors, reusable across every hub's point-to-site gateway
- **Contracts enforced upfront**: each enabled authentication type requires its block (ARM's create rule, validated in seconds instead of at deploy)
- **Policy groups by name**: named member-matching rules deploy as ARM children and surface in the name-keyed `policy_group_ids` output

## Key Features

- **Consistent Interface**: aligns with our existing APIs for deploying cloud infrastructure across multiple providers
- **Three authentication families**: Entra ID (managed sign-in), certificate (offline-capable), RADIUS (existing enterprise auth) -- combinable on one configuration
- **Pinned IPsec proposals**: optional client-side algorithm pinning with the provider's full vocabularies validated upfront
- **Chart-ready wiring**: the configuration's ARM ID surfaces where point-to-site gateways consume it

## Use Cases

- **Remote workforce over Entra ID**: managed sign-in with conditional access, no certificate distribution
- **Lab and contractor access over certificates**: offline-capable auth anchored to your own root
- **Enterprise RADIUS integration**: forward authentication to existing NPS/RADIUS infrastructure
- **User segmentation**: policy groups matching Entra ID groups or certificate common names, mapped to different address pools by the gateway

## Future Enhancements

Future updates will include:

- **Per-group experience surfacing**: richer policy-group diagnostics as Azure exposes them

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
