# Read-Only Policy on a Key Prefix

This preset defines the most common custom policy: read-only access
scoped to one application's key prefix -- finer than the built-in
"Data Reader", which reads every key on the cache.

## When to Use

- Multiple applications sharing one cache, each confined to its own
  prefix
- Dashboards, debuggers, and read replicas of application logic that
  must never write
- The least-privilege default for any consumer that only reads

## Key Configuration Choices

- **`+@read +@connection`** -- the read command category plus the
  connection-management commands every client needs (AUTH, PING, SELECT)
- **`~<keyPrefix>:*`** -- the key pattern confines every granted command
  to that prefix; add more `~pattern` terms for multiple prefixes
- **Permissions update in place** -- widen or narrow the grant without
  recreating assignments

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<cache-resource-name>` | The AzureRedisCache's Planton resource name | Your cache composition |
| `<policyName>` | The policy's name (what assignments reference) | Your naming convention, e.g. `orders-read-only` |
| `<keyPrefix>` | The application's key namespace | Your key-naming convention |
