# Cloudflare Account API Token

An account-owned API token: a scoped Cloudflare credential with permission policies, client-IP conditions, and an optional validity window. Real create, update, delete.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Account API token** -- one `cloudflare_account_token`

## Prerequisites

- **A Cloudflare API token** with Account API Tokens → Write (a token that can mint tokens)
- **Permission-group UUIDs** from `GET /accounts/{account_id}/tokens/permission_groups`

## Quick Start

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareAccountApiToken
metadata:
  name: ci-dns-editor
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: ci-dns-editor
  policies:
    - effect: allow
      permissionGroupIds:
        - "4755a26eedb94da69e1066f98e79d058"
      resources:
        "com.cloudflare.api.account.0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d":
          permission: "*"
```

```shell
planton apply -f account-api-token.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `accountId` | string | The owning account. | Required, 32-hex; replaces on change. |
| `name` | string | Token name. | Required. |
| `policies` | list | What the token may do. | At least one; each needs an effect, at least one permission group, and at least one resource entry. |

### Optional Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `expiresOn` | string | Expiry timestamp. | RFC 3339; must be after `notBefore`. |
| `notBefore` | string | Activation timestamp. | RFC 3339. |
| `condition.requestIp.inCidrs` | list | Allowed client CIDRs. | IPv4/IPv6 CIDR strings. |
| `condition.requestIp.notInCidrs` | list | Barred client CIDRs. | Evaluated before the allow list. |
| `status` | string | Administrative state. | `active` or `disabled` (`expired` and `revoked (exposed)` are server-reported only). |

### Policy Resource Entries

Each entry under `policies[].resources` grants either the whole resource (`permission`, normally `*`) or a nested scoping (`subresources`) -- exactly one of the two. Cloudflare takes this as one raw JSON object; the modules serialize it for you.

## Destroy Semantics

Real delete. The credential stops working immediately -- retire consumers before destroying the token.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `tokenId` | string | The token's management id (not the credential) |
| `value` | string | The secret token value -- returned once at create; secret-marked |

## Related Components

- [Cloudflare Zero Trust Access Service Token](/docs/catalog/cloudflare/cloudflarezerotrustaccessservicetoken) -- machine credentials for Access-protected apps
- [Cloudflare Secrets Store Secret](/docs/catalog/cloudflare/cloudflaresecretsstoresecret) -- where a minted value can live for Workers to consume
