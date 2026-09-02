# Allow staff

A simple `allow` policy that grants access to anyone with a corporate email domain,
connecting from an allowed country, with a 24-hour session.

## When to use

- A baseline "staff can access this" policy attached to one or more applications.

## Key choices

- `decision: allow` with `include` (corporate domain) and `require` (country).
- `sessionDuration`: how long before re-authentication.

## Placeholders

| Placeholder | Description |
|---|---|
| `0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d` | 32-character Cloudflare account ID |

## Referencing it from an application

```yaml
policies:
  - policy:
      valueFrom:
        kind: CloudflareZeroTrustAccessPolicy
        name: allow-staff
        fieldPath: status.outputs.policy_id
    precedence: 1
```
