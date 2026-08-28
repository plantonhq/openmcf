---
display_name: GitHub OAuth
---

# GitHub OAuth

Connects Access to a GitHub OAuth App so users sign in with their GitHub identity. Creation validates the credential shape, not that users can complete a login -- a dummy client id/secret is enough to create the provider.

## When to Use

- Engineering teams already living in GitHub
- A second sign-in option beside a corporate IdP
- An Access policy that matches `github_organization` rules

## Key Configuration Choices

- **type: github** -- immutable. Changing it later replaces the provider and breaks policy references.
- **client_secret as StringValueOrRef** -- use a managed-secret `value_from` in production; a literal value is only for first-run experiments.
- **No SCIM** -- GitHub is not a SCIM source in this preset. Add `scim_config` only for IdPs that push user lifecycle events (Okta, Azure AD).

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `account_id` | The Cloudflare account that owns this provider | Cloudflare Dashboard -> Overview -> API section |
| `config.client_id` | GitHub OAuth App client id | GitHub -> Settings -> Developer settings -> OAuth Apps |
| `config.client_secret.value` | GitHub OAuth App client secret | The same OAuth App page; prefer a secret reference |

## Related Presets

- **01-onetimepin** -- no IdP application required
- **03-okta** -- corporate IdP with SCIM
