---
title: "Cognito User Pool Client"
description: "Cognito User Pool Client deployment documentation"
icon: "package"
order: 100
componentName: "awscognitouserpoolclient"
---

# AWS Cognito User Pool Client

Deploys a Cognito user pool app client -- the OAuth 2.0 / OIDC contract between one application and a user pool. Each client carries its own grant types, redirect URLs, token lifetimes, authentication flows, and client ID; that client ID is what API Gateway JWT authorizers validate as the token audience and what ALB authenticate-cognito actions reference.

## What Gets Created

When you deploy an AwsCognitoUserPoolClient resource, Planton provisions:

- **App Client** -- an `aws_cognito_user_pool_client` on the referenced pool, with the configured OAuth contract, token lifetimes, and security posture

## Prerequisites

- **An AwsCognitoUserPool** -- the pool this client authenticates against (referenced by `userPoolId`)
- (Optional) **AwsCognitoIdentityProvider** resources if this client offers federated sign-in
- (Optional) **An AwsCognitoResourceServer** if this client requests custom OAuth scopes (required for `client_credentials`)

## Quick Start

Create a file `client.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsCognitoUserPoolClient
metadata:
  name: web-app
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsCognitoUserPoolClient.web-app
spec:
  region: us-east-1
  userPoolId:
    valueFrom:
      kind: AwsCognitoUserPool
      name: my-auth
      fieldPath: status.outputs.user_pool_id
  explicitAuthFlows:
    - ALLOW_USER_SRP_AUTH
    - ALLOW_REFRESH_TOKEN_AUTH
  preventUserExistenceErrors: ENABLED
```

Deploy:

```shell
planton apply -f client.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | AWS region. | Required |
| `userPoolId` | `StringValueOrRef` | The pool this client authenticates against. ForceNew. | Required |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `generateSecret` | `bool` | `false` | Generate a client secret (confidential/server-side clients). ForceNew. Public clients (SPA, mobile) use PKCE instead. |
| `allowedOauthFlowsUserPoolClient` | `bool` | `false` | Enable the OAuth 2.0 fields below. |
| `allowedOauthFlows` | `string[]` | `[]` | `"code"` (recommended), `"implicit"` (legacy), `"client_credentials"` (M2M -- cannot mix with the user-facing grants). |
| `allowedOauthScopes` | `string[]` | `[]` | Built-in scopes (`"openid"`, `"email"`, `"profile"`, `"phone"`, `"aws.cognito.signin.user.admin"`) and resource-server scopes (`"{identifier}/{scope}"`). |
| `callbackUrls` | `string[]` | `[]` | OAuth redirect URLs after sign-in (max 100). Required for `code`/`implicit`. |
| `logoutUrls` | `string[]` | `[]` | Redirect URLs after sign-out (max 100). |
| `defaultRedirectUri` | `string` | — | Default callback; must be one of `callbackUrls`. |
| `supportedIdentityProviders` | `StringValueOrRef[]` | all pool providers | `"COGNITO"` for the pool's own directory plus federated provider names -- literals or references to AwsCognitoIdentityProvider resources. |
| `explicitAuthFlows` | `string[]` | AWS defaults | `"ALLOW_USER_SRP_AUTH"`, `"ALLOW_REFRESH_TOKEN_AUTH"`, `"ALLOW_USER_PASSWORD_AUTH"`, `"ALLOW_ADMIN_USER_PASSWORD_AUTH"`, `"ALLOW_CUSTOM_AUTH"`, `"ALLOW_USER_AUTH"` (choice-based/passwordless). |
| `authSessionValidity` | `int` | 3 (AWS) | Challenge-session validity in minutes (3-15). |
| `accessTokenValidity` / `idTokenValidity` | `int` | 1 hour (AWS) | Token lifetimes in `tokenValidityUnits` units. AWS bounds: 5 minutes - 24 hours. |
| `refreshTokenValidity` | `int` | 30 days (AWS) | Refresh lifetime. AWS bounds: 60 minutes - 10 years. |
| `tokenValidityUnits.accessToken` / `idToken` / `refreshToken` | `string` | `hours`/`hours`/`days` | `"seconds"`, `"minutes"`, `"hours"`, `"days"`. |
| `refreshTokenRotation.feature` | `string` | — | `"ENABLED"` / `"DISABLED"` -- each refresh issues a new refresh token and retires the old one. When ENABLED, do not also list `"ALLOW_REFRESH_TOKEN_AUTH"` in `explicitAuthFlows` (AWS rejects the combination). |
| `refreshTokenRotation.retryGracePeriodSeconds` | `int` | 0 | 0-60 seconds the retired token stays usable after rotation. |
| `enableTokenRevocation` | `bool` | `true` (AWS) | Sign-out revokes the refresh token and tokens minted from it. |
| `enablePropagateAdditionalUserContextData` | `bool` | `false` | Forward client IP/user-agent to threat protection in server-side flows. Requires a client secret. |
| `preventUserExistenceErrors` | `string` | — | `"ENABLED"` (recommended -- prevents user enumeration) or `"LEGACY"`. |
| `readAttributes` / `writeAttributes` | `string[]` | all | Attribute-level access for this client. |
| `analyticsConfiguration.applicationArn` | `string` | — | Pinpoint project ARN (mutually exclusive with `applicationId`). |
| `analyticsConfiguration.applicationId` + `externalId` + `roleArn` | mixed | — | The cross-account Pinpoint arm; `roleArn` accepts an AwsIamRole reference. |
| `analyticsConfiguration.userDataShared` | `bool` | `false` | Include endpoint attributes in published events. |

## Examples

### Public SPA Client (PKCE, Hosted UI)

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsCognitoUserPoolClient
metadata:
  name: web-spa
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: acme
    pulumi.planton.dev/project: platform
    pulumi.planton.dev/stack.name: prod.AwsCognitoUserPoolClient.web-spa
spec:
  region: us-east-1
  userPoolId:
    valueFrom:
      kind: AwsCognitoUserPool
      name: prod-auth
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

### Machine-to-Machine Client (client_credentials)

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsCognitoUserPoolClient
metadata:
  name: billing-service
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: acme
    pulumi.planton.dev/project: platform
    pulumi.planton.dev/stack.name: prod.AwsCognitoUserPoolClient.billing-service
spec:
  region: us-east-1
  userPoolId:
    valueFrom:
      kind: AwsCognitoUserPool
      name: prod-auth
      fieldPath: status.outputs.user_pool_id
  generateSecret: true
  allowedOauthFlowsUserPoolClient: true
  allowedOauthFlows:
    - client_credentials
  allowedOauthScopes:
    - https://api.example.com/read
    - https://api.example.com/write
  accessTokenValidity: 30
  tokenValidityUnits:
    accessToken: minutes
```

The custom scopes come from an AwsCognitoResourceServer with identifier `https://api.example.com`.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `client_id` | `string` | The app client ID -- the JWT `aud` claim JWT authorizers validate, the ALB `user_pool_client_id`, and the SDK client identifier. |
| `client_secret` | `string` | The client secret; populated only when `generateSecret: true`. Sensitive -- treat as a credential. |
| `user_pool_id` | `string` | The pool this client belongs to, resolved from the spec reference -- application configs typically need the (pool id, client id) pair together. |

## Related Components

- [AWS Cognito User Pool](/docs/catalog/aws/cognito-user-pool) -- the pool this client authenticates against
- [AWS Cognito Identity Provider](/docs/catalog/aws/cognito-identity-provider) -- federated providers this client can offer at sign-in
- [AWS Cognito Resource Server](/docs/catalog/aws/cognito-resource-server) -- custom scopes for machine-to-machine clients
- [AWS HTTP API Gateway](/docs/catalog/aws/http-api-gateway) -- JWT authorizers listing this client's `client_id` as an audience
