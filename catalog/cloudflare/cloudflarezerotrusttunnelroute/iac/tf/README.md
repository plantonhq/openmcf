# Terraform Module: Cloudflare Zero Trust Tunnel Route

Provisions a single `cloudflare_zero_trust_tunnel_cloudflared_route` — a
private CIDR made reachable to WARP clients through an existing cloudflared
tunnel, optionally inside a specific virtual network (omit for the account's
default network).

## Layout

```
iac/tf/
├── provider.tf    # cloudflare provider ~> 5.23
├── variables.tf   # metadata + spec
├── locals.tf      # ref resolution for tunnel_id / virtual_network_id
├── main.tf        # the cloudflare_zero_trust_tunnel_cloudflared_route resource
└── outputs.tf     # route_id, network
```

## Inputs

A `spec` matching `CloudflareZeroTrustTunnelRouteSpec`. Required: `account_id`,
`network` (the CIDR — unique per virtual network), and `tunnel_id` (literal or
a reference to a CloudflareZeroTrustTunnel's output). Optional:
`virtual_network_id` (literal or a reference to a
CloudflareZeroTrustTunnelVirtualNetwork's output) and `comment`.

## Outputs

- `route_id` — the route's UUID.
- `network` — the routed CIDR as recorded by Cloudflare.

## Requirements

- API token with **Account → Cloudflare Tunnel → Edit**.
