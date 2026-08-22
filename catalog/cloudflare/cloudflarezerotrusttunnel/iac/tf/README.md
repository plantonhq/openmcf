# Terraform Module: Cloudflare Zero Trust Tunnel

Provisions a `cloudflare_zero_trust_tunnel_cloudflared` — a cloudflared tunnel
whose connector makes an outbound-only connection to Cloudflare's edge — plus,
when the tunnel is remotely managed and declares ingress rules, its companion
`cloudflare_zero_trust_tunnel_cloudflared_config` (the remote ingress
configuration, a separate provider resource so editing ingress never recreates
the tunnel).

## Layout

```
iac/tf/
├── provider.tf    # cloudflare provider ~> 5.23
├── variables.tf   # metadata + spec
├── locals.tf      # ingress + origin_request shaping (unset values become null)
├── main.tf        # the tunnel + conditional config resource
└── outputs.tf     # tunnel_id, tunnel_cname, tunnel_token (sensitive), ...
```

## Inputs

A `spec` matching `CloudflareZeroTrustTunnelSpec`. Required: `account_id` and
`name`. `config_src` selects remote (`cloudflare`, the default) vs `local`
management; `ingress` rules (with a mandatory trailing catch-all) and
`origin_request` defaults apply only to a remotely-managed tunnel — enforced
by the spec's CEL before the module ever runs. `tunnel_secret` is optional
(omit to let Cloudflare generate one).

## Outputs

- `tunnel_id` — the tunnel's UUID (routes reference it).
- `tunnel_cname` — the `<id>.cfargotunnel.com` target DNS records CNAME to.
- `tunnel_token` — the connector run token (**sensitive**).
- `tunnel_status`, `account_tag`, `created_on`.

## Requirements

- API token with **Account → Cloudflare Tunnel → Edit**.
