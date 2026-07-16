# Custom Policy Grant

This preset completes the three-kind composition: the cache defines the
boundary, an AzureRedisCacheAccessPolicy defines WHAT is allowed, and
this grant defines WHO gets it.

## When to Use

- Granting the prefix-scoped or command-scoped policies the built-ins
  cannot express (see the AzureRedisCacheAccessPolicy presets)
- Multi-app caches where every identity gets exactly its own namespace

## Key Configuration Choices

- **The policy is a reference** -- deploy ordering resolves
  automatically: cache, then policy, then grant, all in one manifest set
- **One policy, many grants** -- reuse the same policy resource across
  every identity that needs the same shape of access
- **Everything is ForceNew** -- changing the grant replaces it
  (momentary revoke-and-regrant, safe for the grant class)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<cache-resource-name>` | The AzureRedisCache's Planton resource name | Your cache composition |
| `<assignmentName>` | The grant's name, unique within the cache | Convention: `<identity>-<policy>` |
| `<policy-resource-name>` | The AzureRedisCacheAccessPolicy's Planton resource name | Your policy composition |
| `<identity-resource-name>` | The AzureUserAssignedIdentity's Planton resource name | Your identity composition |
