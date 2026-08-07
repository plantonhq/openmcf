# AWS Cognito Identity Provider

Deploys an external identity provider federated into an Amazon Cognito User Pool. Configures social (Google, Facebook, Amazon, Apple), enterprise OIDC (Okta, Azure AD, Auth0), or SAML 2.0 providers so users can sign in through the IdP and receive Cognito tokens. The `provider_name` output is referenced by User Pool Clients in `supportedIdentityProviders` to enable federated sign-in.

## What Gets Created

When you deploy an AwsCognitoIdentityProvider resource, Planton provisions:

- **Cognito Identity Provider** — an `aws_cognito_identity_provider` resource attached to the specified User Pool, with provider-specific configuration (OAuth credentials, OIDC issuer, or SAML metadata) and attribute mapping

## Prerequisites

- **AwsCognitoUserPool** (or equivalent) must exist; `userPoolId` references its `status.outputs.user_pool_id`
- **Provider credentials** — OAuth client ID/secret from the IdP, or SAML metadata URL/file
- **AWS credentials** configured via environment variables or Planton provider config

## Quick Start

Create a file `google-idp.yaml`:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsCognitoIdentityProvider
metadata:
  name: google-idp
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsCognitoIdentityProvider.google-idp
spec:
  region: us-west-2
  userPoolId:
    valueFrom:
      kind: AwsCognitoUserPool
      name: my-auth
      fieldPath: status.outputs.user_pool_id
  providerName: Google
  providerType: Google
  google:
    clientId: "${GOOGLE_CLIENT_ID}"
    clientSecret: "${GOOGLE_CLIENT_SECRET}"
    authorizeScopes: "email profile openid"
  attributeMapping:
    email: email
    username: sub
```

Deploy:

```shell
planton apply -f google-idp.yaml
```

Then add `Google` to the User Pool Client's `supportedIdentityProviders` list.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | AWS region where the resource will be created (e.g., `us-west-2`). | Required |
| `userPoolId` | `StringValueOrRef` | User Pool ID (e.g., `us-east-1_Ab1Cd2EfG`). Can reference AwsCognitoUserPool via `valueFrom`. | Required. ForceNew. |
| `providerName` | `string` | Display name for this IdP. Referenced in User Pool Client `supportedIdentityProviders`. | 1-32 UTF-8 chars. ForceNew. |
| `providerType` | `string` | Provider type, exactly as the AWS API spells it: `Google`, `Facebook`, `LoginWithAmazon`, `SignInWithApple`, `OIDC`, `SAML`. | Required. ForceNew. |
| `google` / `facebook` / `loginWithAmazon` / `signInWithApple` / `oidc` / `saml` | oneof | Provider-specific config. Must match `providerType`. | Exactly one required. |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `attributeMapping` | `map<string, string>` | Provider defaults | Maps IdP attributes to Cognito user pool attributes. Keys: Cognito attrs; values: IdP claim names. |
| `idpIdentifiers` | `string[]` | `[]` | Alternative identifiers for `idp_identifier` query param. Max 50, each 1-40 chars. |
| `facebook.apiVersion` | `string` | — | Graph API version (e.g., `v17.0`). Omit for Cognito default. |
| `oidc.authorizeScopes` | `string` | — | Space-separated scopes. Defaults to `openid email profile` when omitted. |
| `oidc.clientSecret` | `string` | — | OIDC client secret. Omit for public clients using PKCE. |
| `oidc.attributesRequestMethod` | `string` | `GET` | `GET` or `POST` for userinfo endpoint. |
| `oidc.authorizeUrl` | `string` | — | Override auto-discovered authorization endpoint. |
| `oidc.tokenUrl` | `string` | — | Override auto-discovered token endpoint. |
| `oidc.attributesUrl` | `string` | — | Override auto-discovered userinfo endpoint. |
| `oidc.jwksUri` | `string` | — | Override auto-discovered JWKS endpoint. |
| `oidc.attributesUrlAddAttributes` | `bool` | `false` | Append requested attributes as query parameters to the userinfo request (only some providers need this). |
| `saml.metadataFile` | `string` | — | Inline SAML metadata XML. For SAML, set one of `metadataFile` or `metadataUrl`. |
| `saml.metadataUrl` | `string` | — | URL to IdP SAML metadata. For SAML, set one of `metadataFile` or `metadataUrl`. |
| `saml.idpSignOut` | `bool` | `false` | Enable single logout (SLO). |
| `saml.idpInit` | `bool` | `false` | Enable IdP-initiated SSO. |
| `saml.encryptedResponses` | `bool` | `false` | Require encrypted SAML assertions. |
| `saml.requestSigningAlgorithm` | `string` | — | Algorithm for signing AuthnRequest (e.g., `rsa-sha256`). |

## Examples

### Google OAuth

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsCognitoIdentityProvider
metadata:
  name: google-idp
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: acme
    pulumi.planton.dev/project: auth
    pulumi.planton.dev/stack.name: prod.AwsCognitoIdentityProvider.google-idp
spec:
  region: us-west-2
  userPoolId:
    valueFrom:
      kind: AwsCognitoUserPool
      name: prod-auth
      fieldPath: status.outputs.user_pool_id
  providerName: Google
  providerType: Google
  google:
    clientId: "${GOOGLE_CLIENT_ID}"
    clientSecret: "${GOOGLE_CLIENT_SECRET}"
    authorizeScopes: "email profile openid"
  attributeMapping:
    email: email
    username: sub
    given_name: given_name
    family_name: family_name
```

### Enterprise OIDC

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsCognitoIdentityProvider
metadata:
  name: corp-oidc-idp
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: acme
    pulumi.planton.dev/project: auth
    pulumi.planton.dev/stack.name: prod.AwsCognitoIdentityProvider.corp-oidc-idp
spec:
  region: us-west-2
  userPoolId:
    valueFrom:
      kind: AwsCognitoUserPool
      name: prod-auth
      fieldPath: status.outputs.user_pool_id
  providerName: CorpSSO
  providerType: OIDC
  oidc:
    clientId: "${OIDC_CLIENT_ID}"
    clientSecret: "${OIDC_CLIENT_SECRET}"
    oidcIssuer: "https://login.microsoftonline.com/${TENANT_ID}/v2.0"
    authorizeScopes: "openid email profile"
  attributeMapping:
    email: email
    username: sub
    given_name: given_name
    family_name: family_name
```

### SAML Federation

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsCognitoIdentityProvider
metadata:
  name: corp-saml-idp
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: acme
    pulumi.planton.dev/project: auth
    pulumi.planton.dev/stack.name: prod.AwsCognitoIdentityProvider.corp-saml-idp
spec:
  region: us-west-2
  userPoolId:
    valueFrom:
      kind: AwsCognitoUserPool
      name: prod-auth
      fieldPath: status.outputs.user_pool_id
  providerName: CorpAD
  providerType: SAML
  saml:
    metadataUrl: "https://idp.example.com/saml/metadata"
    idpSignOut: true
    idpInit: true
  attributeMapping:
    email: "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress"
    given_name: "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/givenname"
    family_name: "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/surname"
```

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `provider_name` | `string` | Name of the identity provider. Add this value to the User Pool Client's `supportedIdentityProviders` to enable federated sign-in. |
| `provider_type` | `string` | Provider type (e.g., `Google`, `OIDC`, `SAML`). Informational. |
| `user_pool_id` | `string` | The pool this provider is attached to, resolved from the spec reference. |

## Related Components

- [AWS Cognito User Pool](/docs/catalog/aws/cognito-user-pool) — parent resource; provides the `user_pool_id` this IdP attaches to
- [AWS Cognito User Pool Client](/docs/catalog/aws/cognito-user-pool-client) — lists this IdP (by reference to its `provider_name` output) in `supportedIdentityProviders` to enable federated sign-in
