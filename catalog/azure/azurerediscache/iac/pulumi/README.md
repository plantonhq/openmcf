# AzureRedisCache - Pulumi Module

Pulumi implementation for the AzureRedisCache component.

## Architecture

```
redis.Cache (the cache)
└── redis.FirewallRule (xN, public-endpoint IP allow-list)
```

## Key Design Decisions

- **The size family is derived, never spelled** -- Azure sizes a cache as
  `{family}{capacity}` and the family letter ("C" vs "P") is fully
  determined by the tier, so the module derives it and the spec carries
  only tier + capacity.
- **The tier default materializes in the module** -- stack inputs never
  carry proto defaults, so an unspecified sku deploys STANDARD
  explicitly (matching the Terraform module's coalesce).
- **`redis_configuration` is emitted only when the spec carries it** --
  an omitted block deploys Azure's engine defaults identically on both
  engines; unset memory dials are never sent so Azure sizes them from
  total memory.
- **Presence-guarded proto defaults** -- redis_version "6",
  public_network_access true, access_keys_authentication true,
  maxmemory_policy volatile-lru, authentication_enabled true, patch
  windows PT5H starting at hour 0.
- **Only `replicas_per_primary` is modeled** -- ARM's legacy
  `replicas_per_master` alias mirrors it server-side; modeling both
  would be contradictable redundant state.
- **Firewall rules fold in; grants and links do not** -- rules are pure
  IP filters with no life outside the cache, while Entra grants
  (AzureRedisCacheAccessPolicyAssignment) and geo-replication links
  (AzureRedisLinkedServer) are first-class kinds with their own
  lifecycles.
- **Keys and connection strings are exported as secret-bearing outputs**
  -- both primary AND secondary faces, so clients can rotate with zero
  downtime; keyless (Entra-only) caches leave them empty.
- **PARITY-EXCEPTION (tag shape)**: `resource_kind` here is the lowered
  CloudResourceKind enum string and `resource_id` is omitted when
  metadata.id is empty, while the Terraform module emits the family-wide
  snake-case literal and falls back to metadata.name. Output-neutral;
  aligning the shapes is a family-wide convention change.

## Provider

The Azure provider is built by the shared
`pulumiazureprovider.Get(ctx, stackInput.ProviderConfig)` builder, which
dispatches static client-secret, keyless web-identity (OIDC), and
ambient credential chains. Never construct a provider inline.
