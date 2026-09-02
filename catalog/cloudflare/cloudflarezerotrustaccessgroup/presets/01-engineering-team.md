# Engineering team group

A reusable account-scoped group that matches your engineering staff by email domain,
requires an allowed country, and excludes a known contractor account.

## When to use

- You repeatedly grant the same team access across multiple applications.

## Key choices

- `include`: any matching rule adds the user (here, corporate email domains).
- `require`: every rule must also hold (here, an allowed country).
- `exclude`: any matching rule removes the user.

## Placeholders

| Placeholder | Description |
|---|---|
| `0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d` | 32-character Cloudflare account ID |

## Referencing it from a policy

```yaml
include:
  - group:
      id:
        valueFrom:
          kind: CloudflareZeroTrustAccessGroup
          name: engineering-team
          fieldPath: status.outputs.group_id
```
