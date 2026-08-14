---
title: "Pinned IPsec Connection"
description: "This preset pins the tunnel to an exact proposal (AES-256/SHA-256, DH group 14, PFS 2048 -- a widely supported compliance baseline) and carries a pre-agreed key. With a pinned proposal there is NO..."
type: "preset"
rank: "02"
presetSlug: "02-pinned-ipsec-connection"
componentSlug: "vpn-gateway-connection"
componentTitle: "VPN Gateway Connection"
provider: "azure"
icon: "package"
order: 2
---

# Pinned IPsec Connection

This preset pins the tunnel to an exact proposal (AES-256/SHA-256, DH group 14, PFS 2048 -- a widely supported compliance baseline) and carries a pre-agreed key. With a pinned proposal there is NO fallback: the branch device must offer exactly this suite.

## When to Use

- Compliance regimes that mandate specific cipher suites
- Branch devices whose vendors publish an exact supported matrix

## Key Configuration Choices

- **Pin from the device's matrix** -- take the values from the branch device vendor's documentation, not a generic hardening guide; a mismatch leaves the tunnel provisioned-but-never-Connected
- **Reference a secret for the key** -- the field is sensitive and reference-capable; replace the literal placeholder with a `valueFrom` secret reference in real use
- **Explicit DPD** -- 45s is ARM's default made visible; tune it for flaky branch uplinks

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-vpn-gateway-arm-id>` | ARM ID of the hub's VPN gateway | `AzureVpnGateway` status outputs (`vpn_gateway_id`), or reference it with valueFrom |
| `<your-vpn-site-arm-id>` | ARM ID of the branch's site | `AzureVpnSite` status outputs (`vpn_site_id`), or reference it with valueFrom |
| `<your-vpn-site-link-arm-id>` | ARM ID of the site link being connected | `AzureVpnSite` status outputs (`link_ids.<link-name>`), or reference it with valueFrom |
| `<your-pre-shared-key>` | The agreed tunnel key | Your secret store -- reference it, never embed it |
