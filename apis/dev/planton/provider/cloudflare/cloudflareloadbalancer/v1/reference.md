# CloudflareLoadBalancer

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `cloudflare.planton.dev/v1`

CloudflareLoadBalancerSpec defines a zone-scoped Cloudflare Load Balancer. The
load balancer attaches a DNS hostname to a set of account-scoped pools
(CloudflareLoadBalancerPool) and steers traffic across them with the chosen
steering policy, session affinity, and optional geo-routing. Pools and their
monitors are independent, reusable resources referenced here by ID or reference.

## Example

```yaml
apiVersion: cloudflare.planton.dev/v1
kind: CloudflareLoadBalancer
metadata:
  name: lb-hack
spec:
  hostname: lb.planton-example.com
  zoneId:
    value: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  defaultPools:
    - value: "f1e2d3c4b5a6978869584a3b2c1d0e0f"
  fallbackPool:
    value: "f1e2d3c4b5a6978869584a3b2c1d0e0f"
  proxied: true
  steeringPolicy: random
  sessionAffinity: cookie
  sessionAffinityTtl: 1800
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.hostname` | `string` | yes |  |  |
| `spec.zoneId` | `string \| valueFrom` | yes |  | CloudflareDnsZone (`status.outputs.zone_id`) |
| `spec.proxied` | `bool` |  | `true` |  |
| `spec.sessionAffinity` | `enum` |  |  |  |
| `spec.steeringPolicy` | `enum` |  |  |  |
| `spec.defaultPools` | `[]string \| valueFrom` | yes |  | CloudflareLoadBalancerPool (`status.outputs.pool_id`) |
| `spec.fallbackPool` | `string \| valueFrom` | yes |  | CloudflareLoadBalancerPool (`status.outputs.pool_id`) |
| `spec.description` | `string` |  |  |  |
| `spec.enabled` | `bool` |  | `true` |  |
| `spec.ttl` | `double` |  |  |  |
| `spec.sessionAffinityTtl` | `double` |  |  |  |
| `spec.sessionAffinityAttributes` | `CloudflareLoadBalancerSessionAffinityAttributes` |  |  |  |
| `spec.sessionAffinityAttributes.drainDuration` | `double` |  |  |  |
| `spec.sessionAffinityAttributes.headers` | `[]string` |  |  |  |
| `spec.sessionAffinityAttributes.requireAllHeaders` | `bool` |  |  |  |
| `spec.sessionAffinityAttributes.samesite` | `string` |  |  |  |
| `spec.sessionAffinityAttributes.secure` | `string` |  |  |  |
| `spec.sessionAffinityAttributes.zeroDowntimeFailover` | `string` |  |  |  |
| `spec.regionPools` | `[]CloudflareLoadBalancerGeoPools` |  |  |  |
| `spec.regionPools[].code` | `string` | yes |  |  |
| `spec.regionPools[].poolIds` | `[]string \| valueFrom` | yes |  | CloudflareLoadBalancerPool (`status.outputs.pool_id`) |
| `spec.countryPools` | `[]CloudflareLoadBalancerGeoPools` |  |  |  |
| `spec.countryPools[].code` | `string` | yes |  |  |
| `spec.countryPools[].poolIds` | `[]string \| valueFrom` | yes |  | CloudflareLoadBalancerPool (`status.outputs.pool_id`) |
| `spec.popPools` | `[]CloudflareLoadBalancerGeoPools` |  |  |  |
| `spec.popPools[].code` | `string` | yes |  |  |
| `spec.popPools[].poolIds` | `[]string \| valueFrom` | yes |  | CloudflareLoadBalancerPool (`status.outputs.pool_id`) |
| `spec.adaptiveRouting` | `CloudflareLoadBalancerAdaptiveRouting` |  |  |  |
| `spec.adaptiveRouting.failoverAcrossPools` | `bool` |  |  |  |
| `spec.locationStrategy` | `CloudflareLoadBalancerLocationStrategy` |  |  |  |
| `spec.locationStrategy.mode` | `string` |  |  |  |
| `spec.locationStrategy.preferEcs` | `string` |  |  |  |
| `spec.randomSteering` | `CloudflareLoadBalancerRandomSteering` |  |  |  |
| `spec.randomSteering.defaultWeight` | `double` |  |  |  |
| `spec.randomSteering.poolWeights` | `map<string, double>` |  |  |  |

## Field Details

### spec.hostname

`string` · required

The DNS hostname to associate with this load balancer (e.g.
"app.example.com"). If a DNS record with this name already exists, the load
balancer takes precedence.

- rule: {"required":true}

### spec.zoneId

`string | valueFrom` · required

The Cloudflare zone that owns the hostname, as a literal zone ID or a
reference to a CloudflareDnsZone.

- references: CloudflareDnsZone (`status.outputs.zone_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: CloudflareDnsZone, name: <that resource's name>, fieldPath: status.outputs.zone_id}} -- a bare string does not parse

### spec.proxied

`bool` · optional (explicit presence)

Whether the hostname is proxied through Cloudflare (orange cloud). Defaults
to false (gray cloud) when unset; most HTTP load balancers set this true.

- default: `true`

### spec.sessionAffinity

`enum`

Session-affinity mode. Defaults to none.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `none` -- No session affinity (the default).
- `cookie` -- Cookie-based affinity.
- `ip_cookie` -- Cookie-based affinity with stable initial selection by client IP.
- `header` -- Header-based affinity (see session_affinity_attributes.headers).

### spec.steeringPolicy

`enum`

Traffic-steering policy. Defaults to off (static failover over default_pools).

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `off` -- Use default_pools as a static failover priority list (also the default).
- `geo` -- Use region_pools / country_pools / pop_pools (geo-routing).
- `random` -- Select a pool at random, weighted by random_steering.
- `dynamic_latency` -- Use round-trip time to pick the closest healthy pool in default_pools.
- `proximity` -- Use pool latitude/longitude to pick the closest pool.
- `least_outstanding_requests` -- Pick the pool with the fewest outstanding requests (weighted).
- `least_connections` -- Pick the pool with the fewest open connections (weighted).

### spec.defaultPools

`[]string | valueFrom` · required

Ordered list of pools used by default (and when no geo mapping matches),
by failover priority. Each is a literal pool ID or a reference to a
CloudflareLoadBalancerPool. At least one is required.

- references: CloudflareLoadBalancerPool (`status.outputs.pool_id`)
- rule: {"repeated":{"minItems":"1"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: CloudflareLoadBalancerPool, name: <that resource's name>, fieldPath: status.outputs.pool_id}} -- a bare string does not parse

### spec.fallbackPool

`string | valueFrom` · required

The pool of last resort, used when every other pool is unhealthy. A literal
pool ID or a reference to a CloudflareLoadBalancerPool.

- references: CloudflareLoadBalancerPool (`status.outputs.pool_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: CloudflareLoadBalancerPool, name: <that resource's name>, fieldPath: status.outputs.pool_id}} -- a bare string does not parse

### spec.description

`string`

Human-readable description shown in the Cloudflare dashboard.

### spec.enabled

`bool` · optional (explicit presence)

Whether the load balancer is enabled. Defaults to true when unset.

- default: `true`

### spec.ttl

`double`

TTL (seconds) of the DNS entry returned by the load balancer. Applies only to
gray-clouded (unproxied) load balancers. Leave 0 to use Cloudflare's default.

- rule: ttl must be 0 (default) or a positive number of seconds

### spec.sessionAffinityTtl

`double`

Seconds until a client's affinity session expires. Leave 0 for the
Cloudflare default for the chosen session_affinity mode.

- rule: session_affinity_ttl must be 0 (default) or a positive number of seconds

### spec.sessionAffinityAttributes

`CloudflareLoadBalancerSessionAffinityAttributes`

Fine-grained session-affinity attributes (drain, headers, cookie flags,
zero-downtime failover). Omit to use defaults.

### spec.sessionAffinityAttributes.drainDuration

`double`

Drain duration (seconds) applied when affinity is enabled.

- rule: drain_duration must be 0 or a positive number of seconds

### spec.sessionAffinityAttributes.headers

`[]string`

HTTP header names that header-based affinity is keyed on (required when
session_affinity is header). Use "cookie:<name>,<name>" to scope to cookies.

### spec.sessionAffinityAttributes.requireAllHeaders

`bool`

When header affinity is enabled, require ALL listed headers (true) or AT
LEAST ONE (false) for a session to be created.

### spec.sessionAffinityAttributes.samesite

`string`

SameSite attribute on the affinity cookie: one of "Auto", "Lax", "None",
"Strict". Empty uses Cloudflare's default ("Auto").

- rule: samesite must be one of "Auto", "Lax", "None", "Strict"

### spec.sessionAffinityAttributes.secure

`string`

Secure attribute on the affinity cookie: one of "Auto", "Always", "Never".
Empty uses Cloudflare's default ("Auto").

- rule: secure must be one of "Auto", "Always", "Never"

### spec.sessionAffinityAttributes.zeroDowntimeFailover

`string`

Zero-downtime failover behavior for pinned sessions: one of "none",
"temporary", "sticky". Empty uses Cloudflare's default ("none").

- rule: zero_downtime_failover must be one of "none", "temporary", "sticky"

### spec.regionPools

`[]CloudflareLoadBalancerGeoPools`

Region-code -> ordered pool list for geo steering. Used when
steering_policy is geo (or auto-geo). Regions not listed fall back to
default_pools.

### spec.regionPools[].code

`string` · required

The geo code this mapping applies to: a region code (e.g. "WNAM"), an ISO
country code (e.g. "US"), or a Cloudflare PoP code (e.g. "LAX"), depending on
which map this entry belongs to.

- rule: {"required":true}

### spec.regionPools[].poolIds

`[]string | valueFrom` · required

Ordered pools (by failover priority) for this code. Each is a literal pool ID
or a reference to a CloudflareLoadBalancerPool.

- references: CloudflareLoadBalancerPool (`status.outputs.pool_id`)
- rule: {"repeated":{"minItems":"1"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: CloudflareLoadBalancerPool, name: <that resource's name>, fieldPath: status.outputs.pool_id}} -- a bare string does not parse

### spec.countryPools

`[]CloudflareLoadBalancerGeoPools`

Country-code (ISO 3166-1 alpha-2) -> ordered pool list. Countries not listed
fall back to the matching region_pools entry, else default_pools.

### spec.countryPools[].code

`string` · required

The geo code this mapping applies to: a region code (e.g. "WNAM"), an ISO
country code (e.g. "US"), or a Cloudflare PoP code (e.g. "LAX"), depending on
which map this entry belongs to.

- rule: {"required":true}

### spec.countryPools[].poolIds

`[]string | valueFrom` · required

Ordered pools (by failover priority) for this code. Each is a literal pool ID
or a reference to a CloudflareLoadBalancerPool.

- references: CloudflareLoadBalancerPool (`status.outputs.pool_id`)
- rule: {"repeated":{"minItems":"1"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: CloudflareLoadBalancerPool, name: <that resource's name>, fieldPath: status.outputs.pool_id}} -- a bare string does not parse

### spec.popPools

`[]CloudflareLoadBalancerGeoPools`

Cloudflare PoP code -> ordered pool list (Enterprise only). PoPs not listed
fall back to country_pools, then region_pools, then default_pools.

### spec.popPools[].code

`string` · required

The geo code this mapping applies to: a region code (e.g. "WNAM"), an ISO
country code (e.g. "US"), or a Cloudflare PoP code (e.g. "LAX"), depending on
which map this entry belongs to.

- rule: {"required":true}

### spec.popPools[].poolIds

`[]string | valueFrom` · required

Ordered pools (by failover priority) for this code. Each is a literal pool ID
or a reference to a CloudflareLoadBalancerPool.

- references: CloudflareLoadBalancerPool (`status.outputs.pool_id`)
- rule: {"repeated":{"minItems":"1"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: CloudflareLoadBalancerPool, name: <that resource's name>, fieldPath: status.outputs.pool_id}} -- a bare string does not parse

### spec.adaptiveRouting

`CloudflareLoadBalancerAdaptiveRouting`

Adaptive-routing controls (zero-downtime failover across pools).

### spec.adaptiveRouting.failoverAcrossPools

`bool`

Extend zero-downtime failover to healthy origins in alternate pools when no
healthy origin remains in the same pool. Defaults to false.

### spec.locationStrategy

`CloudflareLoadBalancerLocationStrategy`

Location-steering controls for non-proxied requests.

### spec.locationStrategy.mode

`string`

Authoritative location when ECS is not used: "pop" (Cloudflare PoP) or
"resolver_ip" (DNS resolver GeoIP). Empty uses Cloudflare's default ("pop").

- rule: mode must be "pop" or "resolver_ip"

### spec.locationStrategy.preferEcs

`string`

Whether to prefer the ECS GeoIP location: one of "always", "never",
"proximity", "geo". Empty uses Cloudflare's default ("proximity").

- rule: prefer_ecs must be one of "always", "never", "proximity", "geo"

### spec.randomSteering

`CloudflareLoadBalancerRandomSteering`

Pool weighting for random / least-* steering policies.

### spec.randomSteering.defaultWeight

`double`

Default weight for pools not present in pool_weights (0.0–1.0). Leave 0 to
use the Cloudflare default (1).

- rule: default_weight must be between 0 and 1

### spec.randomSteering.poolWeights

`map<string, double>`

Per-pool weight overrides, keyed by pool ID. Each weight is relative to the
other pools in this load balancer.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: CloudflareLoadBalancer, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.load_balancer_id` | `string` | Unique identifier of the Cloudflare load balancer. |
| `status.outputs.load_balancer_dns_record_name` | `string` | The hostname DNS record associated with the load balancer. |
| `status.outputs.load_balancer_cname_target` | `string` | The canonical CNAME target that the hostname resolves to (Cloudflare endpoint). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.zoneId` | CloudflareDnsZone | `status.outputs.zone_id` |
| `spec.defaultPools` | CloudflareLoadBalancerPool | `status.outputs.pool_id` |
| `spec.fallbackPool` | CloudflareLoadBalancerPool | `status.outputs.pool_id` |
| `spec.regionPools[].poolIds` | CloudflareLoadBalancerPool | `status.outputs.pool_id` |
| `spec.countryPools[].poolIds` | CloudflareLoadBalancerPool | `status.outputs.pool_id` |
| `spec.popPools[].poolIds` | CloudflareLoadBalancerPool | `status.outputs.pool_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
