# Cloudflare Zero Trust MCP Portal

An MCP portal: the Access-protected endpoint publishing a curated collection of registered MCP servers to users and their AI clients.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **MCP portal** -- one `cloudflare_zero_trust_access_ai_controls_mcp_portal`

## Prerequisites

- **Zero Trust enabled on the account** (the team-name onboarding step)
- **A Cloudflare API token** with Account → Zero Trust → Edit
- **Registered MCP servers** (`CloudflareZeroTrustMcpServer`) to publish, unless starting empty

## Quick Start

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZeroTrustMcpPortal
metadata:
  name: eng-tools
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  portalId: eng-tools
  hostname: mcp.example.com
  name: Engineering Tools
```

```shell
planton apply -f portal.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `accountId` | string | The Cloudflare account. | Required, 32-hex; replaces on change. |
| `portalId` | string | Your slug for the portal. | Required; replaces on change. |
| `hostname` | string | The serving hostname. | Required; updates in place. |
| `name` | string | Display name. | Required. |

### Optional Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `description` | string | Free-form description. | |
| `allowCodeMode` | bool | Client code execution against portal tools. | Cloudflare default: enabled. |
| `secureWebGateway` | bool | Filter portal traffic through Gateway. | |
| `servers` | list | Published servers with per-portal overrides. | Unordered set; `serverId` required per row. |

## Destroy Semantics

Real delete of the portal endpoint; the server registrations it published are untouched.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `portal_id` | string | The portal's identifier (the user-chosen slug) |
| `hostname` | string | The hostname users and AI clients connect to |

## Related Components

- [Cloudflare Zero Trust MCP Server](/docs/catalog/cloudflare/cloudflarezerotrustmcpserver) -- the registrations published here
- [Cloudflare Zero Trust Access Application](/docs/catalog/cloudflare/cloudflarezerotrustaccessapplication) -- the broader Access family
