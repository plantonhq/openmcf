# AWS Cognito Resource Server

Deploys a Cognito resource server -- the OAuth 2.0 resource (an API) a user pool mints custom access-token scopes for. Each scope it defines becomes requestable by app clients as `{identifier}/{scope_name}`, and machine-to-machine clients using the `client_credentials` grant can request only these custom scopes.

## What Gets Created

When you deploy an AwsCognitoResourceServer resource, Planton provisions:

- **Resource Server** -- an `aws_cognito_resource_server` on the referenced pool with the configured identifier and scope vocabulary

## Prerequisites

- **An AwsCognitoUserPool** -- the pool this resource server belongs to (referenced by `userPoolId`)

## Quick Start

Create a file `resource-server.yaml`:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsCognitoResourceServer
metadata:
  name: orders-api
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsCognitoResourceServer.orders-api
spec:
  region: us-east-1
  userPoolId:
    valueFrom:
      kind: AwsCognitoUserPool
      name: my-auth
      fieldPath: status.outputs.user_pool_id
  identifier: https://api.example.com
  name: Orders API
  scopes:
    - scopeName: read
      scopeDescription: Read orders
    - scopeName: write
      scopeDescription: Create and update orders
```

Deploy:

```shell
planton apply -f resource-server.yaml
```

A machine-to-machine client can now list `https://api.example.com/read` in its `allowedOauthScopes`.

## Configuration Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `region` | `string` | Yes | AWS region. |
| `userPoolId` | `StringValueOrRef` | Yes | The pool this resource server belongs to. ForceNew. |
| `identifier` | `string` | Yes | Unique identifier within the pool -- the scope prefix access tokens carry (conventionally the API's audience URI, e.g. `https://api.example.com`). 1-256 characters. ForceNew. |
| `name` | `string` | Yes | Human-readable display name shown in the Cognito console. 1-256 characters. |
| `scopes` | `object[]` | No | Up to 100 custom scopes, each with `scopeName` (no spaces, `/`, `"`, `\`) and `scopeDescription`. Scope names must be unique. Scopes update in place. |

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `resource_server_identifier` | `string` | The identifier -- the scope prefix access tokens carry. |
| `scope_identifiers` | `string[]` | Fully-qualified scope strings (`{identifier}/{scope_name}`) -- the exact values app clients list in `allowedOauthScopes`. |
| `user_pool_id` | `string` | The pool this resource server belongs to, resolved from the spec reference. |

## Related Components

- [AWS Cognito User Pool](/docs/catalog/aws/cognito-user-pool) -- the pool this resource server belongs to
- [AWS Cognito User Pool Client](/docs/catalog/aws/cognito-user-pool-client) -- clients that request this resource server's scopes, most importantly machine-to-machine `client_credentials` clients
