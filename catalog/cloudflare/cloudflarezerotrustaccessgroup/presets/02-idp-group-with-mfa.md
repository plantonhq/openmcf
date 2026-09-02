---
display_name: IdP Group With MFA
---

# IdP group with MFA login method

An account-scoped group that matches an Okta group (federated through a configured
identity provider) and additionally requires a hardware-key authentication method.

## When to use

- Membership is driven by an IdP group and you want an extra authentication-method
  requirement layered on.

## Key choices

- `include.okta`: the Okta group name plus the Cloudflare identity-provider ID.
- `require.authMethod`: an AMR value such as `hwk` (hardware key) or `mfa`.

## Placeholders

| Placeholder | Description |
|---|---|
| `0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d` | 32-character Cloudflare account ID (replace with yours) |
| `<okta-identity-provider-id>` | The Cloudflare identity-provider ID for the Okta connection (a value-or-reference field: set `value`, or reference a `CloudflareZeroTrustAccessIdentityProvider`) |
