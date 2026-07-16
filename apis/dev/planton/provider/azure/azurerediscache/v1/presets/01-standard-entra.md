# Standard Cache with Entra Authentication

This preset creates the production default: a replicated Standard-tier
cache with Microsoft Entra token authentication enabled alongside the
access keys -- the on-ramp to a fully keyless cache.

## When to Use

- Application caching, session state, and rate limiting for most
  production workloads
- Teams moving toward token-based (secretless) Redis access
- Any cache that does not need VNet injection, clustering, or
  persistence (those are Premium capabilities)

## Key Configuration Choices

- **STANDARD tier** -- primary/replica pair, 99.9% SLA; upgrade to
  PREMIUM in place when clustering or persistence is needed (downgrades
  replace the cache)
- **Entra auth ON, keys still ON** -- grant identities through
  AzureRedisCacheAccessPolicyAssignment, migrate clients to tokens, then
  flip `accessKeysAuthenticationEnabled: false` for the keyless posture
- **`allkeys-lru` eviction** -- right for pure caches where every key is
  re-derivable; keep the default `volatile-lru` for mixed workloads

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<region>` | Azure region, e.g. `eastus` | Your region strategy |
| `<resource-group-resource-name>` | The AzureResourceGroup's Planton resource name | Your foundation composition |
| `myorg-app-cache` | Globally unique DNS name (1-63 letters/digits/hyphens) | Becomes `{cacheName}.redis.cache.windows.net` |
