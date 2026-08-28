---
display_name: One-Time PIN
---

# One-time PIN

Cloudflare's own one-time PIN identity provider -- no IdP application, no client secret. Users receive a PIN at the Access login page. This is the smallest useful identity provider and the safest starting point.

## When to Use

- A fallback for contractors or guests who have no corporate IdP
- First Access setup on an account that does not yet have Google/Okta/Azure AD connected
- A second provider beside a corporate IdP, not a replacement for it

## Key Configuration Choices

- **type: onetimepin** -- omit `config`; the module sends the empty object Cloudflare expects. SCIM is forbidden on this type.
- **No client secret** -- there is nothing to leak and nothing to rotate.
- **OTP updates may 403 on API-token auth** -- create and destroy work; an in-place update of an OTP provider has historically been rejected. Prefer destroy-and-recreate over field edits.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `account_id` | The Cloudflare account that owns this provider | Cloudflare Dashboard -> Overview -> API section (right sidebar) |

## Related Presets

- **02-github-oauth** -- GitHub as the sign-in source
- **03-okta** -- Okta with optional SCIM
