# Cloudflare Zero Trust MCP Server

An MCP server registration behind Access AI Controls: the upstream tool server Cloudflare proxies, audits, and gates. Portals publish registered servers to users.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **MCP server registration** -- one `cloudflare_zero_trust_access_ai_controls_mcp_server`

## Prerequisites

- **Zero Trust enabled on the account** (the team-name onboarding step)
- **A Cloudflare API token** with Account → Zero Trust → Edit
- **The upstream MCP server's URL** (and its bearer token or OAuth client secret, as managed secrets, unless unauthenticated)

## Quick Start

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZeroTrustMcpServer
metadata:
  name: docs-search
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  serverId: docs-search
  name: Docs Search
  hostname: https://mcp.example.com
  authType: bearer
  authCredentials:
    value: $secret/mcp-docs-search-token
```

```shell
planton apply -f mcp-server.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `accountId` | string | The Cloudflare account. | Required, 32-hex; replaces on change. |
| `serverId` | string | Your identifier for the registration. | Required; replaces on change. |
| `name` | string | Display name. | Required; updates in place. |
| `hostname` | string | The upstream MCP server URL. | Required; replaces on change. |
| `authType` | string | oauth, bearer, or unauthenticated. | Required; replaces on change. |

### Optional Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `authCredentials` | secret ref | Bearer token for auth_type bearer. | Write-only; never read back. |
| `clientSecret` | secret ref | OAuth client secret (manual OAuth mode). | Write-only; never read back. |
| `description` | string | Free-form description. | |
| `updatedPrompts` / `updatedTools` | list | Per-item overrides (alias, description, enabled). | `name` required per row. |
| `isSharedOauthCallbackEnabled` | bool | Use Cloudflare's shared OAuth callback URL. | |
| `secureWebGateway` | bool | Filter upstream traffic through Gateway. | |

## Destroy Semantics

Real delete. Portals referencing the server must drop their rows first.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `server_id` | string | The server's identifier -- what portals reference |

## Related Components

- [Cloudflare Zero Trust MCP Portal](/docs/catalog/cloudflare/cloudflarezerotrustmcpportal) -- publishes this server
- [Cloudflare Zero Trust Gateway Settings](/docs/catalog/cloudflare/cloudflarezerotrustgatewaysettings) -- the Gateway posture
