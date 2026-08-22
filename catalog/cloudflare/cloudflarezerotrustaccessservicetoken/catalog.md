# Cloudflare Zero Trust Access Service Token

A machine credential -- a client-ID / client-secret pair that non-human clients present in the `CF-Access-Client-ID` / `CF-Access-Client-Secret` headers to pass through Access-protected applications without an identity-provider login. The secret is returned only at creation and at rotation; capture it at deploy or you will have to rotate.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Access service token** -- one `cloudflare_zero_trust_access_service_token` at the account (or zone)
- **Client ID + client secret** -- exported as stack outputs; the secret is marked sensitive and is never readable again except on rotation

## Prerequisites

- **A Cloudflare account with Zero Trust enabled** -- the organization (team name) must already exist or every Access create fails at the API
- **A Cloudflare API token** with Account → Access: Service Tokens → Edit
- **A secret store** to capture `client_secret` at deploy -- the value cannot be recovered from Cloudflare later

## Quick Start

A named token with Cloudflare's one-year default duration:

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZeroTrustAccessServiceToken
metadata:
  name: ci-deployer
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: ci-deployer
```

```shell
planton apply -f service-token.yaml
```

Capture `status.outputs.client_id` and `status.outputs.client_secret` into your CI secret store before anything else.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `accountId` XOR `zoneId` | string / StringValueOrRef | Exactly one. `zoneId` can reference a CloudflareDnsZone via `valueFrom`. | Exactly one required. `accountId` is a 32-character hex string when set. |
| `name` | string | Display name. | Required, min length 1. |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `duration` | string | Cloudflare's 8760h (one year) | Go-style duration (`2h45m`, `8760h`) or `forever`. |
| `clientSecretVersion` | optional int32 | unset (Cloudflare treats the initial secret as version 1) | Increment to rotate. Must be set together with `previousClientSecretExpiresAt`. ≥ 1. |
| `previousClientSecretExpiresAt` | string | unset | RFC3339. When the previous secret stops being accepted after a rotation. Required whenever `clientSecretVersion` is set. |

## Examples

### Minimal

Name only. Duration defaults to one year. Capture the secret at deploy.

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZeroTrustAccessServiceToken
metadata:
  name: ci-deployer
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: ci-deployer
```

### Non-expiring token

For a long-lived automation identity. Still rotate if the secret leaks.

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZeroTrustAccessServiceToken
metadata:
  name: monitoring-scraper
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: monitoring-scraper
  duration: forever
```

### Rotation

Mint version 2 and keep the old secret working for 30 days so services can migrate.

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZeroTrustAccessServiceToken
metadata:
  name: ci-deployer
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: ci-deployer
  clientSecretVersion: 2
  previousClientSecretExpiresAt: "2026-09-15T00:00:00Z"
```

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `service_token_id` | string | UUID used for import and policy `service_token` rules |
| `client_id` | string | Presented in `CF-Access-Client-ID` |
| `client_secret` | string | Sensitive. Returned only at create and rotation; does not survive import |
| `expires_at` | string | RFC3339 expiry |

## Related Components

- [Cloudflare Zero Trust Access Application](/docs/catalog/cloudflare/cloudflarezerotrustaccessapplication) -- the apps the token unlocks
- [Cloudflare Zero Trust Access Policy](/docs/catalog/cloudflare/cloudflarezerotrustaccesspolicy) -- attach the token via a `service_token` rule
- [Cloudflare Zero Trust Access Identity Provider](/docs/catalog/cloudflare/cloudflarezerotrustaccessidentityprovider) -- human sign-in, not machine credentials
- [Cloudflare DNS Zone](/docs/catalog/cloudflare/cloudflarednszone) -- `zoneId` foreign key for a zone-scoped token
