# Publish a private app on a public hostname

Expose a single private web app at a public hostname through the tunnel — the most common
Cloudflare Tunnel setup.

## When to use

- You have an internal HTTP service and want it reachable at `app.example.com` without
  opening inbound ports.

## Key choices

- `ingress[].service`: the local address of your app (e.g. `http://localhost:8080`).
- The trailing `service: http_status:404` rule is the required catch-all.
- After applying, CNAME `app.example.com` to `status.outputs.tunnel_cname` with a
  `CloudflareDnsRecord`, and run the connector with `status.outputs.tunnel_token`.

## Placeholders

| Placeholder | Description |
|---|---|
| `0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d` | 32-character Cloudflare account ID |
