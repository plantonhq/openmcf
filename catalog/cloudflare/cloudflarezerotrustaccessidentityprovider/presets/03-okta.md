# Okta with SCIM

Connects Access to an Okta tenant and turns on SCIM so user create/update/deprovision events land in Zero Trust without waiting for re-authentication. Enabling SCIM mints `scim_secret` once -- capture it in the same change.

## When to Use

- Okta is the corporate identity source
- You want deprovisioned users to lose Access seats automatically
- Access policies will match Okta groups

## Key Configuration Choices

- **type: okta** -- `okta_account` is gated to this type; setting it on any other type fails validation.
- **scim_config.enabled: true** -- mints `status.outputs.scim_secret` and `scim_base_url`. Configure those at Okta. The secret is redacted on later reads and does not survive import.
- **identity_update_behavior: automatic** -- Cloudflare applies IdP-side profile changes without forcing a re-login. `reauth` and `no_action` are the other values.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `account_id` | The Cloudflare account that owns this provider | Cloudflare Dashboard -> Overview -> API section |
| `config.client_id` / `config.client_secret` | Okta application credentials | Okta Admin -> Applications -> the Access application |
| `config.okta_account` | Okta tenant hostname | Your Okta URL, e.g. `acme.okta.com` |

## Related Presets

- **01-onetimepin** -- no IdP application required
- **02-github-oauth** -- GitHub as the sign-in source
