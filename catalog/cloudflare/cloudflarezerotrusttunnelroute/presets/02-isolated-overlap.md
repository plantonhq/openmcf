# Overlapping CIDR isolated in a virtual network

Advertise a CIDR that overlaps another already-connected network by scoping the route to
its own virtual network, so the two never collide.

## When to use

- Multiple sites/tenants each use the same private range (e.g. `10.0.0.0/8`) and must be
  reachable independently.

## Key choices

- `virtualNetworkId`: reference a dedicated `CloudflareZeroTrustTunnelVirtualNetwork`;
  each overlapping CIDR goes in its own virtual network through its own tunnel.

## Placeholders

| Placeholder | Description |
|---|---|
| `0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d` | 32-character Cloudflare account ID |
| `<tenant-a-tunnel>` | Tunnel serving tenant A's network |
| `<tenant-a-vnet>` | Virtual network isolating tenant A's CIDR |
