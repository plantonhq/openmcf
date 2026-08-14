# CloudflareZeroTrustTunnel guide

Operational judgment for cloudflared tunnels. The README covers what each field is; this covers how the pieces interact.

## The tunnel object and the connector are two different things

Deploying this kind creates the tunnel's identity at Cloudflare's edge and exports the run token — it does not run cloudflared anywhere. Nothing flows until you start a connector (`cloudflared tunnel run --token <token>`) on a machine that can reach your origins. Expect `tunnel_status` to read `inactive` until then; that is correct, not broken.

## Remote config is the default for a reason

With `config_src: cloudflare` (the default), ingress rules live in this spec and editing them never recreates the tunnel — the config is a separate resource underneath. `local` mode hands ingress to a YAML file on the origin machine: choose it only when the machine's operator owns routing and the control plane should manage identity only. Switching modes later moves ownership of the ingress config, so pick deliberately.

## The catch-all rule is mandatory, and DNS is on you

Cloudflare requires the last ingress rule to answer for every unmatched hostname — `service: http_status:404` with no hostname is the canonical shape (validation enforces it). And each real hostname only resolves once a DNS record CNAMEs it to `tunnel_cname`; composing a CloudflareDnsRecord per hostname is the intended graph, not an extra.

## The run token is a credential

`tunnel_token` grants the ability to serve your tunnel's traffic. It lives in the outputs (marked sensitive) so machines can consume it; treat any place it lands — CI variables, VM bootstrap, dashboards — as secret storage. Rotating means recreating the tunnel today.

## Destroy is a soft delete

Cloudflare keeps a deleted tunnel's record around (`deleted_at` is set) rather than 404ing it. Dashboards and API lists may still show it filtered views apart; a deleted tunnel's ID cannot be re-adopted. Active connectors block deletion — stop cloudflared first.
