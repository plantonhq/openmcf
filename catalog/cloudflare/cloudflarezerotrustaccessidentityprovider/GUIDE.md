# CloudflareZeroTrustAccessIdentityProvider guide

The judgment this guide protects you from: the provider type is a one-way door, and the SCIM secret is shown once. Pick the type deliberately, and capture `scim_secret` the moment SCIM is enabled -- Cloudflare will not show it again.

## Type is immutable: changing it replaces the provider

Cloudflare's `type` field RequiresReplace. A rename from `github` to `okta` is not an update -- it is a new identity provider with a new ID. Every Access policy rule that referenced the old `identity_provider_id` (azure_ad, github_organization, gsuite, okta, saml, oidc, login_method, auth_context) now points at a deleted object. The discipline is: create a second provider, attach policies to the new ID, then destroy the old one.

The spellings are Cloudflare's own. `azureAD` and `google-apps` are not typos. Lowercasing them fails validation here instead of failing at the API.

## Config is a union; the wrong field is a validation error

The provider's ConfigValidators are empty. Per-type gating (`directory_id` only on `azureAD`, `auth_url` only on `oidc`, SAML fields only on `saml`) is ours. A GitHub provider with `okta_account` set is rejected at manifest validation, not after a confusing API error.

`onetimepin` is the empty-config type. The module still sends Cloudflare the empty config object it expects -- omit `config` in the spec, do not invent a dummy client id.

`enable_encryption` on a SAML provider requires `saml_certificate_set_id`. The certificate set is created out-of-band via the Access SAML-certificate API; this kind does not mint it.

## The SCIM secret is create-only

Enabling `scim_config` mints a bearer token. Cloudflare returns it on that create (and when `enabled` flips from false/null to true) and redacts it on every later read. Import cannot recover it. Capture `status.outputs.scim_secret` into a secret store in the same change that turns SCIM on. If you lose it, refresh via the Access API's `refresh_scim_secret` endpoint -- do not destroy and recreate the provider just to see the secret again (that changes the provider ID).

`scim_config` is forbidden on `onetimepin`. Cloudflare rejects the combination; validation does too.

## OTP and API-token auth

The OTP Access endpoint has historically rejected API-token auth on **update**. The provider's own OTP tests unset `CLOUDFLARE_API_TOKEN` for that reason. Create and destroy with a token work; an update (including an idempotency re-plan that Cloudflare treats as an update) may 403. If that happens, the fix is not to switch credentials in the harness -- record the defect and keep OTP as a create/destroy-only path until Cloudflare accepts tokens on that endpoint.

## Read-only is a latch, not a suggestion

`read_only: true` tells Cloudflare to refuse API updates and deletes. Clearing it takes an explicit apply. Use it on the corporate IdP you cannot afford to delete from a bad pull request.

## Pairs well with

- [CloudflareZeroTrustAccessApplication](../cloudflarezerotrustaccessapplication/README.md) / [CloudflareZeroTrustAccessPolicy](../cloudflarezerotrustaccesspolicy/README.md) -- the doors and guards that consume `identity_provider_id`.
- [CloudflareZeroTrustAccessServiceToken](../cloudflarezerotrustaccessservicetoken/README.md) -- machine credentials, not human sign-in.
- [CloudflareDnsZone](../cloudflarednszone/README.md) -- `zone_id` via `valueFrom` for a zone-scoped provider.
