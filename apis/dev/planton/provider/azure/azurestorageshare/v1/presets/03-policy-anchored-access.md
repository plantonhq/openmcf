# Policy-Anchored Access Share

This preset creates a share whose external access rides stored access
policies (signed identifiers). SAS tokens issued against a policy
inherit its window and permissions -- and revoking the policy instantly
revokes every token anchored to it, the operational lever ad-hoc SAS
tokens lack.

## When to Use

- Partner or vendor file exchange without identity federation
- Time-boxed external access that must be revocable in one action
- Any SAS-based integration that outlives a single token's lifetime

## Key Configuration Choices

- **Two policies, two blast radii** -- `partner-readers` (`rl`) for
  consumers, `partner-writers` (`rwdl`) for producers; revoke either
  without touching the other
- **Permission letters follow Azure's strict order** -- r, w, d, l;
  out-of-order strings are rejected by the spec before they reach Azure
- **Policy-level expiry is the revocation lever** -- set real windows;
  at most five policies per share (Azure's limit)
- **`accessTier: COOL`** -- drop zones are read rarely; pay less at rest

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<storage-account-resource-name>` | The AzureStorageAccount's Planton resource name | Your storage composition |
| `<share-name>` | 3-63 lowercase letters/digits/hyphens | Your naming convention |
| `<window-start-rfc3339>` | Policy validity start, e.g. `2026-07-01T00:00:00Z` | Your exchange agreement |
| `<window-expiry-rfc3339>` | Policy validity end -- the revocation backstop | Your exchange agreement |
