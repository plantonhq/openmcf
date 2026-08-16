# Rotation

Increments the client-secret version so Cloudflare mints a new secret, and keeps the previous secret working until the expiry. The token ID and client ID stay the same, so Access policy `service_token` rules keep matching.

## When to Use

- A leaked or aging secret
- Scheduled rotation without deleting the token
- Killing a compromised secret immediately (set the expiry in the past)

## Key Configuration Choices

- **Both-or-neither** -- `client_secret_version` and `previous_client_secret_expires_at` must be set together. Validation rejects exactly one.
- **Version 2** -- the first rotation from the implicit version 1. Subsequent rotations increment again.
- **Capture the new secret** -- the `client_secret` output changes on this apply. The old value is gone from Cloudflare after the expiry.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `account_id` | The Cloudflare account that owns this token | Cloudflare Dashboard -> Overview -> API section |
| `previous_client_secret_expires_at` | RFC3339 instant the old secret dies | Pick a migration window; past = revoke now |

## Related Presets

- **01-minimal** -- first create, no rotation
