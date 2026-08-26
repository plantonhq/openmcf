# Machine-to-Machine (Client Credentials)

**Rank**: 2

## When to Use

- Service-to-service authentication with no user present
- Backend jobs, daemons, and internal APIs authorizing against each other

## What It Provides

- The `client_credentials` grant with a generated client secret (confidential client; ForceNew)
- Custom scopes from an `AwsCognitoResourceServer` -- the ONLY scopes this grant can request (built-in OIDC scopes describe a user, and there is none)
- Short-lived 30-minute access tokens

## What You Might Add

- More scopes as the resource server's vocabulary grows -- reference its `scope_identifiers` output for the exact strings
- A second client per calling service, so each service holds its own credential and can be revoked independently
