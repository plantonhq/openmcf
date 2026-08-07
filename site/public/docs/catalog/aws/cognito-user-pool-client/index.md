---
title: "Cognito User Pool Client"
description: "Cognito User Pool Client deployment documentation"
icon: "package"
order: 100
componentName: "awscognitouserpoolclient"
---

# AWS Cognito User Pool Client

Deploys a Cognito User Pool app client — the OAuth 2.0 / OIDC contract between ONE application and a user pool. Each application surface (web SPA, mobile app, M2M service) gets its own client with its own grant types, redirect URLs, token lifetimes, and client ID. Downstream systems reference that client ID (API Gateway JWT authorizer audiences, ALB authenticate-cognito actions, SDK configs).

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Cognito User Pool Client** -- registered on the target user pool with the configured OAuth flows, scopes, redirect URLs, auth flows, token lifetimes, attribute access, and optional Pinpoint analytics
- **Client Secret** -- minted only when `generateSecret` is true (confidential / M2M clients); public clients authenticate with PKCE instead
- **AWS Tags** -- resource metadata tags applied automatically for tracking and governance

## Before You Deploy

### Planton Setup

- **AWS Provider Connection** -- an active connection in the Connect module with credentials for the target AWS account.
- **Planton Runner** -- required when using Runner-based credential delivery.

### AWS Account

- **A Cognito User Pool** -- the directory this application authenticates against. Reference an AwsCognitoUserPool Cloud Resource or provide the pool ID directly.
- **Identity providers** (optional) -- list federated providers by name or reference AwsCognitoIdentityProvider resources in `supportedIdentityProviders`.
- **Resource server scopes** (optional) -- custom OAuth scopes from an AwsCognitoResourceServer use the form `identifier/scope-name`.

## Deploy

### Console

Open the deployment store, find **AWS Cognito User Pool Client**, and click **Deploy**. Start from the **Web SPA** preset for public clients, **M2M Client Credentials** for machine grants, or **Server Confidential** for server-side apps with a secret.

### CLI

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsCognitoUserPoolClient
metadata:
  name: web-spa-client
  org: acme-corp
  env: prod
spec:
  region: us-west-2
  userPoolId:
    valueFrom:
      kind: AwsCognitoUserPool
      name: app-auth
      fieldPath: status.outputs.user_pool_id
  allowedOauthFlowsUserPoolClient: true
  allowedOauthFlows:
    - code
  allowedOauthScopes:
    - openid
    - email
    - profile
  callbackUrls:
    - https://app.example.com/callback
  logoutUrls:
    - https://app.example.com/logout
  supportedIdentityProviders:
    - value: COGNITO
  explicitAuthFlows:
    - ALLOW_USER_SRP_AUTH
    - ALLOW_REFRESH_TOKEN_AUTH
  enableTokenRevocation: true
  preventUserExistenceErrors: ENABLED
```

```shell
planton apply -f cognito-user-pool-client.yaml
```

### InfraChart

Reference the pool and any federated providers with ValueFromRef so the deployment graph creates them first. Put custom resource-server scopes into `allowedOauthScopes` using the resource server's `scope_identifiers` output.

## Key Configuration

**Generate secret** -- ForceNew. On for confidential server-side and M2M clients; off for public SPAs and mobile apps that use PKCE.

**OAuth grants** -- Authorization Code is the modern default. Client Credentials is machine-to-machine and cannot share a client with user-facing grants. Callback URLs are required for browser redirects; the default redirect URI must be one of them.

**Refresh token rotation** -- When ENABLED, each refresh issues a new refresh token. Drop `ALLOW_REFRESH_TOKEN_AUTH` from explicit auth flows — AWS rejects the combination.

## After Deployment

Stack outputs include `client_id` (the hero identifier apps and authorizers embed), `client_secret` (sensitive; only when a secret was generated), and the resolved `user_pool_id`.

## Related Resources

- **AwsCognitoUserPool** -- the user directory this client authenticates against
- **AwsCognitoIdentityProvider** -- federated providers listed in `supportedIdentityProviders`
- **AwsCognitoResourceServer** -- custom OAuth scopes requested in `allowedOauthScopes`
