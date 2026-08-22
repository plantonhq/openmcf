# Cloudflare Zero Trust Access Identity Provider

Connects Cloudflare Access to an identity source -- Google, Okta, Azure AD, GitHub, a generic OIDC or SAML provider, or Cloudflare's own one-time PIN -- so users can sign in to Access-protected applications. The provider type is immutable: changing it replaces the provider and invalidates every policy rule that referenced the old ID.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Access identity provider** -- one `cloudflare_zero_trust_access_identity_provider` at the account (or zone) with the chosen type and config
- **SCIM endpoint + secret** -- created only when `scimConfig.enabled` is true; the secret is returned once in `scim_secret` and redacted on later reads

## Prerequisites

- **A Cloudflare account with Zero Trust enabled** -- the organization (team name) must already exist or every Access create fails at the API
- **A Cloudflare API token** with Account → Access: Organizations, Identity Providers, and Groups → Edit
- **IdP application credentials** for OAuth/OIDC/SAML types (client id/secret, or SAML metadata). One-time PIN needs none

## Quick Start

The smallest useful provider is Cloudflare's own one-time PIN -- no config object:

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZeroTrustAccessIdentityProvider
metadata:
  name: otp-fallback
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: one-time-pin
  type: onetimepin
```

```shell
planton apply -f identity-provider.yaml
```

Users see a one-time PIN option on the Access login page. No IdP application is required.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `accountId` XOR `zoneId` | string / StringValueOrRef | Exactly one. `zoneId` can reference a CloudflareDnsZone via `valueFrom` (defaults to `status.outputs.zone_id`). | Exactly one required. `accountId` is a 32-character hex string when set. |
| `name` | string | Shown on the Access login page. | Required, min length 1. |
| `type` | string | Identity provider type. Immutable at Cloudflare. | Required. One of `onetimepin`, `azureAD`, `saml`, `centrify`, `facebook`, `github`, `google-apps`, `google`, `linkedin`, `oidc`, `okta`, `onelogin`, `pingone`, `yandex`, `cloudflare`. |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `config` | object | omitted (empty object sent for `onetimepin`) | Connection parameters. Per-type fields are rejected on the wrong type. OAuth types need `clientId` + `clientSecret` (sensitive StringValueOrRef). |
| `config.clientId` | string | unset | OAuth / OIDC client id. |
| `config.clientSecret` | StringValueOrRef | unset | Sensitive. Prefer a managed-secret reference. |
| `config.claims` | string[] | unset | Extra OIDC claims to request. |
| `config.emailClaimName` | string | unset | Claim that carries the user's email. |
| `config.pkceEnabled` | bool | unset | Enable PKCE on the OAuth flow. |
| `samlCertificateSetId` | string | unset | Required when `config.enableEncryption` is true. |
| `scimConfig` | object | unset | SCIM provisioning. Forbidden when `type` is `onetimepin`. Enabling mints `scim_secret`. |
| `readOnly` | bool | unset | Cloudflare refuses updates and deletes while set. |

## Examples

### One-time PIN

No IdP application. Useful as a fallback for contractors.

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZeroTrustAccessIdentityProvider
metadata:
  name: otp-fallback
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: one-time-pin
  type: onetimepin
```

### GitHub OAuth

Dummy or real GitHub OAuth App credentials. Creation validates shape, not that users can log in.

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZeroTrustAccessIdentityProvider
metadata:
  name: github-login
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: github
  type: github
  config:
    clientId: your-github-oauth-app-id
    clientSecret:
      value: your-github-oauth-app-secret
```

### Okta with SCIM

Okta account plus SCIM so deprovisioned users lose Access seats automatically.

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZeroTrustAccessIdentityProvider
metadata:
  name: okta-corp
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: okta
  type: okta
  config:
    clientId: okta-client-id
    clientSecret:
      value: your-okta-client-secret
    oktaAccount: acme.okta.com
  scimConfig:
    enabled: true
    identityUpdateBehavior: automatic
    userDeprovision: true
    seatDeprovision: true
```

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `identity_provider_id` | string | UUID referenced by Access policy rules |
| `scim_base_url` | string | Cloudflare SCIM v2.0 endpoint (present when SCIM is enabled) |
| `scim_secret` | string | Sensitive. Minted once when SCIM is first enabled; never readable later and does not survive import |

## Related Components

- [Cloudflare Zero Trust Access Application](/docs/catalog/cloudflare/cloudflarezerotrustaccessapplication) -- the door users reach after signing in
- [Cloudflare Zero Trust Access Policy](/docs/catalog/cloudflare/cloudflarezerotrustaccesspolicy) -- the guards that match this provider's identities
- [Cloudflare Zero Trust Access Service Token](/docs/catalog/cloudflare/cloudflarezerotrustaccessservicetoken) -- machine credentials, not human sign-in
- [Cloudflare DNS Zone](/docs/catalog/cloudflare/cloudflarednszone) -- `zoneId` foreign key for a zone-scoped provider
