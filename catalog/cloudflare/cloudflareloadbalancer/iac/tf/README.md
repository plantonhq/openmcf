# Terraform Module: Cloudflare Load Balancer

Provisions a single zone-scoped `cloudflare_load_balancer` that references
account-scoped pools. Pools and monitors are separate modules
(`cloudflareloadbalancerpool`, `cloudflareloadbalancermonitor`).

## Resources

- `cloudflare_load_balancer` (zone-scoped) — attaches the hostname to the given
  `default_pools` / `fallback_pool` and applies steering, session affinity,
  geo-pool maps, adaptive routing, location strategy, random steering, traffic
  rules, and private networks.

## Inputs

- `metadata` — name/labels.
- `spec` — see [variables.tf](./variables.tf). Required: `hostname`, `zone_id`,
  `default_pools`, `fallback_pool`. Pool references flatten from `StringValueOrRef`
  to plain strings (and lists to `list(string)`); enums flatten to their string
  names (top-level `none`/`off` are omitted so the provider applies its default).
- `spec.rules[]` — the typed rules list is rebuilt into the provider's `rules`
  shape in [locals.tf](./locals.tf). Rule OVERRIDES keep real presence: an
  unset override inherits the load balancer's setting, while an explicit value
  — including `none`/`off`/`0` for the presence-carrying `session_affinity`,
  `steering_policy`, and `priority` fields — is sent as a real override.
  `terminates`/`disabled` are sent only when true (a `fixed_response` rule is
  auto-marked terminating server-side).

## Outputs

| Output | Description |
|---|---|
| `load_balancer_id` | The load balancer ID |
| `load_balancer_dns_record_name` | The hostname |
| `load_balancer_cname_target` | The hostname clients point their DNS at |
| `zone_id` | The owning zone (the API identity is `zones/{zone_id}/load_balancers/{id}`) |

## Requirements

- **Load Balancing add-on** must be enabled on the account (paid add-on); otherwise
  the Load Balancing API returns `403`.
- The provider reads `CLOUDFLARE_API_TOKEN` from the environment. The token needs
  **Zone → Load Balancers → Edit** for the zone-scoped load balancer (distinct from
  the account-level "Load Balancers Account" permission), plus
  **Account → Load Balancing: Monitors and Pools → Edit** for the pools/monitors it
  references, and the zone in its Zone Resources scope.
