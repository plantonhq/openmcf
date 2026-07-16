# Development Cache

This preset creates the smallest, cheapest cache Azure offers -- a
single Basic-tier C0 node with an IP allow-list on the public endpoint.

## When to Use

- Local development and integration testing against a real Azure cache
- CI environments that need a disposable Redis endpoint
- Never for production: Basic has no replica and NO SLA -- a node
  restart loses the cache and takes the endpoint down while it rebuilds

## Key Configuration Choices

- **BASIC / C0** -- 250 MB for the lowest cost; memory dials
  (`maxmemoryReserved` and friends) are not configurable on this tier
- **Firewall allow-list** -- even a dev cache holds real data; scope the
  public endpoint to your office/VPN range (rule names allow letters,
  digits, and underscores only)
- **Upgrade path** -- moving to STANDARD/PREMIUM applies in place;
  the reverse replaces the cache

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<region>` | Azure region, e.g. `eastus` | Your region strategy |
| `<resource-group-resource-name>` | The AzureResourceGroup's Planton resource name | Your foundation composition |
| `myorg-dev-cache` | Globally unique DNS name (1-63 letters/digits/hyphens) | Becomes `{cacheName}.redis.cache.windows.net` |
| `203.0.113.0` / `203.0.113.255` | The allowed IPv4 range (inclusive) | Your office/VPN egress range |
