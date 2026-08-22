# Terraform Module: Cloudflare Zero Trust Tunnel Virtual Network

Provisions a single `cloudflare_zero_trust_tunnel_cloudflared_virtual_network`
— a named routing namespace that lets overlapping private CIDRs coexist in one
account (each tunnel route may target a specific virtual network).

## Layout

```
iac/tf/
├── provider.tf    # cloudflare provider ~> 5.23
├── variables.tf   # metadata + spec
├── locals.tf      # labels
├── main.tf        # the ..._virtual_network resource
└── outputs.tf     # virtual_network_id, virtual_network_name
```

## Inputs

A `spec` matching `CloudflareZeroTrustTunnelVirtualNetworkSpec`. Required:
`account_id` and `name` (unique within the account). Optional: `comment` and
`is_default_network` (promoting a network to the account default demotes the
current one — flip it deliberately).

## Outputs

- `virtual_network_id` — the network's UUID (routes reference it).
- `virtual_network_name`.

## Requirements

- API token with **Account → Cloudflare Tunnel → Edit**.
