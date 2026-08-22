# CloudflareZeroTrustAccessIdentityProvider Terraform Module

Terraform IaC module for connecting Cloudflare Access to an identity source (Google, Okta, Azure AD, GitHub, OIDC, SAML, or one-time PIN).

## Architecture

```
provider.tf   — Cloudflare provider configuration (~> 5.23)
variables.tf  — Input variables mirroring CloudflareZeroTrustAccessIdentityProviderSpec
locals.tf     — Scope + config / scim_config shaping
main.tf       — cloudflare_zero_trust_access_identity_provider
outputs.tf    — identity_provider_id, scim_base_url, scim_secret (sensitive)
```

## Behavior

Set exactly one of `account_id` or `zone_id`. `type` is immutable at Cloudflare (RequiresReplace). For `onetimepin` the module sends the empty config object Cloudflare expects. `scim_secret` is returned only when SCIM is first enabled and is marked sensitive.

## Outputs

| Name | Description |
|------|-------------|
| `identity_provider_id` | UUID referenced by Access policy rules |
| `scim_base_url` | SCIM v2.0 endpoint (present when SCIM is enabled) |
| `scim_secret` | Sensitive. Minted once; does not survive import |

## Provider Version

Uses `cloudflare/cloudflare ~> 5.23`.
