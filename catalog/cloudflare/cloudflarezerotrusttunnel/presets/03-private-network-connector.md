# Private-network connector (for WARP access)

A tunnel with no public hostnames, used purely to make private IP ranges reachable to
WARP clients via routes.

## When to use

- You want enrolled devices to reach internal subnets (databases, internal apps) over the
  private network, with no public ingress.

## Key choices

- No `ingress` rules — reachability comes from `CloudflareZeroTrustTunnelRoute` resources
  that reference this tunnel's `status.outputs.tunnel_id`.
- Run the connector with `status.outputs.tunnel_token`.

## Composition

Pair this with one or more `CloudflareZeroTrustTunnelRoute` resources (optionally scoped
to a `CloudflareZeroTrustTunnelVirtualNetwork` when CIDRs overlap).

## Placeholders

| Placeholder | Description |
|---|---|
| `0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d` | 32-character Cloudflare account ID |
