# AWS Cognito Resource Server

Deploys a Cognito Resource Server — an OAuth 2.0 scope namespace attached to a user pool. Each scope becomes an `identifier/scope-name` string that app clients request and your APIs enforce. The resource server does not host users; it declares the permissions tokens may carry.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Cognito Resource Server** -- registered on the target user pool with a permanent identifier, display name, and the configured OAuth scopes
- **AWS Tags** -- resource metadata tags applied automatically for tracking and governance

## Before You Deploy

### Planton Setup

- **AWS Provider Connection** -- an active connection in the Connect module with credentials for the target AWS account.
- **Planton Runner** -- required when using Runner-based credential delivery.

### AWS Account

- **A Cognito User Pool** -- the directory whose access tokens will carry these scopes. Reference an AwsCognitoUserPool Cloud Resource or provide the pool ID directly.

## Deploy

### Console

Open the deployment store, find **AWS Cognito Resource Server**, and click **Deploy**. Start from the **API Scopes** preset to pre-populate a working configuration.

### CLI

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsCognitoResourceServer
metadata:
  name: orders-api-scopes
  org: acme-corp
  env: prod
spec:
  region: us-west-2
  userPoolId:
    valueFrom:
      kind: AwsCognitoUserPool
      name: app-auth
      fieldPath: status.outputs.user_pool_id
  identifier: https://api.example.com
  name: Orders API
  scopes:
    - scopeName: orders.read
      scopeDescription: Read order history
    - scopeName: orders.write
      scopeDescription: Create and update orders
```

```shell
planton apply -f cognito-resource-server.yaml
```

### InfraChart

Wire the pool with ValueFromRef, then have app clients request the emitted `scope_identifiers` in `allowedOauthScopes`.

## After Deployment

Stack outputs include the resource server identifier, the full `scope_identifiers` list (the exact strings clients request), and the resolved user pool ID. Add those scope strings to an AwsCognitoUserPoolClient's allowed OAuth scopes.

## Related Resources

- **AwsCognitoUserPool** -- the user directory this resource server attaches to
- **AwsCognitoUserPoolClient** -- app clients that request these scopes
- **AwsCognitoIdentityProvider** -- federated sign-in providers on the same pool
