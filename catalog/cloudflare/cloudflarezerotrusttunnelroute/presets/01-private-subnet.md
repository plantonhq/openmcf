# Private subnet via a tunnel

Advertise a private subnet through a tunnel so WARP clients can reach hosts in it — the
most common private-networking route.

## When to use

- You want to reach a private network (databases, internal apps) over Cloudflare Tunnel
  without exposing public hostnames.

## Key choices

- `network`: the private CIDR to advertise (use a `/32` or `/128` for a single host).
- `tunnelId`: reference the `CloudflareZeroTrustTunnel` so the graph deploys it first.
- Leave `virtualNetworkId` unset to use the account default; set it (referencing a
  `CloudflareZeroTrustTunnelVirtualNetwork`) when CIDRs overlap.

## Placeholders

| Placeholder | Description |
|---|---|
| `0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d` | 32-character Cloudflare account ID |
| `<tunnel-name>` | Name of the tunnel that serves this network |
