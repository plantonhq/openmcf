# Auth0 Resource Server (API)

Deploys an Auth0 Resource Server (API) with configurable token settings, scope definitions, and optional RBAC policy enforcement. Resource Servers define the APIs that applications request access to via the OAuth 2.0 `audience` parameter and the permissions that can be granted. Integrates with Planton's Auth0 Provider Connection for credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Auth0 Resource Server** -- an API resource configured with the specified identifier, token lifetime, signing algorithm, and access control settings
- **Resource Server Scopes** -- created only when `scopes` is configured, a scopes resource defining the permissions available for this API

## Before You Deploy

### Planton Setup

- **Auth0 Provider Connection** -- an active connection in the Connect module with Auth0 domain, client ID, and client secret. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credential authentication.

### Auth0 Account

- **An Auth0 tenant** with sufficient API quota.
- **A unique API identifier** (audience URI) that has not already been registered in the tenant. Identifiers cannot be changed after creation, so choose a stable URI like `https://api.example.com/`.

## Deploy

### Console

Open the deployment store, find **Auth0 Resource Server (API)**, and click **Deploy**. The creation wizard walks you through environment and connection configuration and spec fields.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: auth0.planton.dev/v1
kind: Auth0ResourceServer
metadata:
  name: backend-api
  org: acme-corp
  env: prod
spec:
  identifier: https://api.example.com/
  scopes:
    - name: read:data
      description: Read access to data
    - name: write:data
      description: Write access to data
```

```shell
planton apply -f auth0-resource-server.yaml
```

This creates a Resource Server in Auth0 with the specified identifier, RS256 signing, default token lifetimes (86400 seconds), and two scopes. No RBAC enforcement is configured. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring an Auth0 Resource Server. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Identifier** -- The `identifier` field sets the API audience URI used in OAuth authorization requests. This value is permanent -- it cannot be changed after creation. Use a stable, descriptive URI like `https://api.example.com/` or a short identifier like `api.planton.live`.

**Token lifetime** -- The `tokenLifetime` field controls how long access tokens remain valid (default 86400 seconds / 24 hours). For sensitive APIs, reduce this to minutes (e.g., 900 for 15-minute tokens). The `tokenLifetimeForWeb` field sets a separate, typically shorter lifetime for implicit and hybrid flow tokens.

**RBAC enforcement** -- Set `enforcePolicies: true` to enable Auth0's built-in role-based access control. When enabled, role and permission assignments are evaluated during login. Pair with a `tokenDialect` ending in `_authz` (e.g., `access_token_authz`) to include permissions in access tokens.

**Token dialect** -- The `tokenDialect` field controls the access token format. Use `access_token` for standard Auth0 JWTs, or `rfc9068_profile` for IETF-compliant tokens. Append `_authz` to either variant to embed RBAC permissions in the token payload.

**Scopes** -- The `scopes` array defines the permissions applications can request for this API. Follow the `action:resource` naming pattern (e.g., `read:users`, `write:orders`). Scopes appear in access tokens and on consent screens.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `id` | Internal Auth0 identifier for the resource server | Auth0 API operations |
| `identifier` | API identifier (audience) for authorization requests | Auth0Client `apiGrants[].audience` references |
| `name` | Display name of the resource server | Dashboards, audit logs |
| `signing_alg` | Token signing algorithm (`RS256`, `HS256`, `PS256`) | Backend token validation configuration |
| `signing_secret` | Signing secret (HS256 only) | Backend token verification |
| `token_lifetime` | Token validity duration in seconds | Client-side token refresh logic |
| `token_lifetime_for_web` | Token validity for implicit/hybrid flows | SPA token handling |
| `allow_offline_access` | Whether refresh tokens can be issued | Client grant configuration |
| `skip_consent_for_verifiable_first_party_clients` | Whether consent is skipped for first-party apps | UX flow design |
| `enforce_policies` | Whether RBAC is enabled | Authorization middleware configuration |
| `token_dialect` | Access token format | Token parsing logic |
| `is_system` | Whether this is a system-managed resource server | Tenant management automation |
| `client_id` | Associated client ID if linked | Integration wiring |

## Common Patterns

No presets are available yet. Configure directly using the fields documented in the [API Explorer](#api-explorer) tab.

## Works With

This component operates independently and does not reference other deployment components.