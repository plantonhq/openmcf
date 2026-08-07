# Auth0 Application (Client)

Deploys an Auth0 Application with configurable OAuth flows, token settings, and optional API access grants. Supports native, SPA, regular web, and machine-to-machine application types. Integrates with Planton's Auth0 Provider Connection for credential management and ValueFromRef for wiring connection and resource server dependencies.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Auth0 Client (Application)** -- an application registered in Auth0 with the specified type, OAuth settings, callback URLs, JWT configuration, and refresh token behavior
- **Client Grants** -- created only when `apiGrants` is configured, one grant per entry authorizing this client to call the specified API with the listed scopes

## Before You Deploy

### Planton Setup

- **Auth0 Provider Connection** -- an active connection in the Connect module with Auth0 domain, client ID, and client secret. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credential authentication.

### Auth0 Account

- **An Auth0 tenant** with sufficient application quota.
- **An Auth0 Resource Server** if configuring `apiGrants` to authorize this client for API access. Provide the audience directly or reference an Auth0ResourceServer Cloud Resource via ValueFromRef.
- **An Auth0 Connection** if restricting the client to specific identity providers via `enabledConnections`. Provide the connection name directly or reference an Auth0Connection Cloud Resource via ValueFromRef.

## Deploy

### Console

Open the deployment store, find **Auth0 Application (Client)**, and click **Deploy**. The creation wizard walks you through environment and connection configuration and spec fields.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: auth0.planton.dev/v1
kind: Auth0Client
metadata:
  name: my-spa
  org: acme-corp
  env: prod
spec:
  applicationType: spa
  callbacks:
    - https://app.example.com/callback
  allowedLogoutUrls:
    - https://app.example.com
  webOrigins:
    - https://app.example.com
```

```shell
planton apply -f auth0-client.yaml
```

This creates a Single Page Application in Auth0 with OIDC-conformant defaults. No API grants or connection restrictions are configured. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the client to connections and resource servers deployed in the same InfraPipeline:

```yaml
spec:
  enabledConnections:
    - valueFrom:
        kind: Auth0Connection
        name: user-db
        fieldPath: status.outputs.name
  apiGrants:
    - audience:
        valueFrom:
          kind: Auth0ResourceServer
          name: backend-api
          fieldPath: status.outputs.identifier
      scopes:
        - read:data
```

The InfraPipeline resolves the dependency graph, deploys the connection and resource server first, then provisions the client with the resolved values.

## Key Configuration

These are the most important decisions when configuring an Auth0 Application. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Application type** -- The `applicationType` field determines which OAuth flows and security settings Auth0 applies. Use `spa` for browser-based JavaScript apps, `regular_web` for server-rendered apps that can store secrets, `native` for mobile or desktop apps, and `non_interactive` for backend M2M communication using client credentials.

**Grant types** -- The `grantTypes` field controls which OAuth flows the application can use. If omitted, Auth0 assigns defaults based on `applicationType`. Override explicitly when you need refresh tokens (`refresh_token`) or device authorization flows.

**JWT signing algorithm** -- The `jwtConfiguration.alg` field sets the token signing method. RS256 (asymmetric) is recommended for most applications. Use HS256 only when the consumer can securely store the client secret.

**Refresh token rotation** -- The `refreshToken.rotationType` field controls whether a new refresh token is issued on each use. Set to `rotating` for production to limit the impact of token theft. Pair with `expirationType: expiring` and explicit lifetimes for time-bounded sessions.

**API grants** -- The `apiGrants` array authorizes this client to call specific APIs. Each entry pairs an audience (resource server identifier) with a list of scopes. Required for M2M applications that need API access beyond authentication.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **Auth0Connection** (optional) | `enabledConnections` | `status.outputs.name` |
| **Auth0ResourceServer** (optional) | `apiGrants[].audience` | `status.outputs.identifier` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `id` | Internal Auth0 identifier for the client | Auth0 API operations |
| `client_id` | OAuth 2.0 client identifier (public) | Auth0Connection `enabledClients`, application config |
| `client_secret` | OAuth 2.0 client secret (confidential clients only) | Backend service authentication |
| `name` | Application name derived from metadata | Audit logs, monitoring dashboards |
| `application_type` | Configured application type | Downstream validation logic |
| `signing_keys` | RS256 signing keys with cert and thumbprint | JWT signature verification in backend services |
| `callback_url_template` | Whether callback URL templating is enabled | Application URL configuration |
| `allowed_clients` | Clients allowed to perform delegation | Legacy delegation flows |
| `global` | Whether this is the tenant's default client | Tenant identification |
| `token_endpoint_auth_method` | Authentication method for the token endpoint | OAuth integration configuration |

## Common Patterns

No presets are available yet. Configure directly using the fields documented in the [API Explorer](#api-explorer) tab.

## Works With

- [**Auth0 Connection (Identity Provider)**](/cloud-catalog/auth0-connection) -- provides identity provider connections that this client can use for authentication via `enabledConnections`
- [**Auth0 Resource Server (API)**](/cloud-catalog/auth0-resource-server) -- defines APIs that this client can be authorized to access via `apiGrants`