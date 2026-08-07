---
title: "Connection (Identity Provider)"
description: "Connection (Identity Provider) deployment documentation"
icon: "package"
order: 100
componentName: "auth0connection"
---

# Auth0 Connection (Identity Provider)

Deploys an Auth0 Connection that bridges Auth0 with an identity source -- databases, social providers (Google, Facebook, GitHub), or enterprise identity providers (SAML, OIDC, Azure AD/Entra ID). Each connection is configured with a single strategy and its corresponding provider-specific options, then linked to one or more Auth0 applications. Integrates with Planton's Auth0 Provider Connection for credential management and ValueFromRef for wiring client dependencies.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Auth0 Connection** -- a connection resource configured with the specified strategy, display name, and provider-specific options (database, social, SAML, OIDC, or Azure AD)
- **Connection Clients** -- created only when `enabledClients` is configured, a resource linking this connection to the specified Auth0 applications so they can offer it as a login option

## Before You Deploy

### Planton Setup

- **Auth0 Provider Connection** -- an active connection in the Connect module with Auth0 domain, client ID, and client secret. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credential authentication.

### Auth0 Account

- **An Auth0 tenant** with connection quota available.
- **OAuth credentials from the social provider** if using a social strategy (Google, Facebook, GitHub). Obtain a client ID and secret from the provider's developer console.
- **Identity Provider metadata** if using an enterprise strategy -- SAML sign-in endpoint and X.509 certificate, OIDC issuer URL, or Azure AD app registration credentials.
- **Auth0 Client application IDs** if linking the connection to specific applications via `enabledClients`. Provide client IDs directly or reference Auth0Client Cloud Resources via ValueFromRef.

## Deploy

### Console

Open the deployment store, find **Auth0 Connection (Identity Provider)**, and click **Deploy**. The creation wizard walks you through environment and connection configuration and spec fields.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: auth0.planton.dev/v1
kind: Auth0Connection
metadata:
  name: user-db
  org: acme-corp
  env: prod
spec:
  strategy: auth0
  databaseOptions:
    passwordPolicy: good
    bruteForceProtection: true
```

```shell
planton apply -f auth0-connection.yaml
```

This creates an Auth0 hosted database connection with a "good" password policy and brute-force protection enabled. No client restrictions are configured, so all applications in the tenant can use it. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the connection to applications deployed in the same InfraPipeline:

```yaml
spec:
  enabledClients:
    - valueFrom:
        kind: Auth0Client
        name: my-spa
        fieldPath: status.outputs.client_id
```

The InfraPipeline resolves the dependency graph, deploys the client first, then provisions the connection with the resolved client ID.

## Key Configuration

These are the most important decisions when configuring an Auth0 Connection. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Strategy** -- The `strategy` field determines the identity source type and which options block applies. Use `auth0` for a hosted database, social strategy names like `google-oauth2` or `github` for social login, or enterprise strategies like `samlp`, `oidc`, or `waad` for corporate SSO. This field cannot be changed after creation.

**Password policy** -- The `databaseOptions.passwordPolicy` field sets complexity requirements for database connections. Options range from `none` to `excellent` (10+ characters with mixed case, numeric, and special). Use `good` or `excellent` for production workloads with `bruteForceProtection: true`.

**Enterprise provider configuration** -- Each enterprise strategy requires its own options block: `samlOptions` for SAML with sign-in endpoint and signing certificate, `oidcOptions` for OIDC with issuer and client credentials, and `azureAdOptions` for Azure AD with tenant domain and app registration details.

**Enabled clients** -- The `enabledClients` field restricts which Auth0 applications can use this connection. If left empty, no applications can authenticate through it. Provide client IDs directly or reference Auth0Client Cloud Resources via ValueFromRef.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **Auth0Client** (optional) | `enabledClients` | `status.outputs.client_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `id` | Unique Auth0 connection identifier (e.g., `con_...`) | Auth0 API operations |
| `name` | Unique connection name within the tenant | Auth0Client `enabledConnections` references |
| `strategy` | Identity provider strategy type | Downstream validation logic |
| `is_enabled` | Whether the connection is currently enabled | Health checks, monitoring |
| `provisioning_ticket_url` | Self-service setup URL for enterprise connections | IdP onboarding workflows |
| `callback_url` | Auth0 callback URL to register with the identity provider | Social and enterprise IdP configuration |
| `metadata_url` | SAML metadata URL (SAML connections only) | SAML IdP configuration |
| `entity_id` | SAML Service Provider Entity ID (SAML connections only) | SAML IdP trust configuration |
| `enabled_client_ids` | List of application client IDs linked to this connection | Audit, downstream wiring |
| `realms` | Realms/domains for identifier-first authentication | Authentication routing |

## Common Patterns

No presets are available yet. Configure directly using the fields documented in the [API Explorer](#api-explorer) tab.

## Works With

- [**Auth0 Application (Client)**](/cloud-catalog/auth0-client) -- applications that use this connection for authentication, linked via `enabledClients`