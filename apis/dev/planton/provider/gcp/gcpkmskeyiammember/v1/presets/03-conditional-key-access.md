# Conditional Key Access (Time-Bound)

This preset grants key access that expires on its own: an IAM condition gates the grant with a CEL expression, so a human operator gets encrypt/decrypt on one key until a fixed timestamp — no cleanup task, no forgotten standing access to key material.

## When to Use

- A migration or incident requires temporary human access to encrypted data
- Break-glass procedures where key access must provably end
- Any grant on key material that should not become permanent by inertia

## Key Configuration Choices

- **Expiry via `request.time`** — the safest condition pattern; access ends even if everyone forgets the grant exists (the expired grant remains in the policy as a no-op until removed)
- **Condition is part of the grant's identity** — the same role with and without a condition are two independent grants; "editing" the expiry replaces the grant atomically
- **Human member** — conditions commonly gate `user:`/`group:` access; service-to-service access is usually better served by the unconditioned presets

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<kms-key-resource-name>` | The Planton resource name of the GcpKmsKey | Your GcpKmsKey manifest's `metadata.name` |
| `<user-email>` | The human receiving temporary access | Your identity provider |
| `<expiry-timestamp>` | RFC 3339 expiry, e.g. `2027-01-01T00:00:00Z` | Your migration plan |

## Related Presets

- **01-storage-cmek-grant** — Authorize a Google service agent for CMEK
- **02-workload-key-user** — Grant a workload's own service account use of a key
