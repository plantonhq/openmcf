# AWS Cognito Resource Server -- Architecture and Design

## Overview

A resource server is how a Cognito user pool models an API it protects. It has two jobs: naming the API (the `identifier`, conventionally the API's audience URI) and defining the scope vocabulary access tokens may carry for it. Every scope it defines becomes requestable by app clients as `{identifier}/{scope_name}`.

## Why It Is First-Class

The honest-decompose test: a pool protects many APIs (many-per-parent), each resource server's scope identifiers are referenced by app clients' `allowed_oauth_scopes` (referenced-by-others), and an API's scope vocabulary evolves with the API rather than the pool (independent lifecycle).

## The Machine-to-Machine Connection

Custom scopes exist chiefly for the `client_credentials` grant. A machine-to-machine client has no user, so the built-in OIDC scopes (`openid`, `email`, `profile`) are meaningless for it -- AWS only lets it request CUSTOM scopes. Without a resource server, an M2M client literally has nothing to ask for. The composed shape:

1. The resource server mints `https://api.example.com/read`.
2. An `AwsCognitoUserPoolClient` with `allowed_oauth_flows: [client_credentials]` lists that scope.
3. The client exchanges its secret for an access token whose `scope` claim carries it.
4. The API validates the token (issuer = the pool's `issuer` output) and authorizes on the scope.

## Identity and Lifecycle

`identifier` is ForceNew -- it is the resource server's identity within the pool and the prefix baked into every issued token. Scopes update in place: adding one is immediate; removing one stops future tokens from carrying it, while already-issued tokens keep it until they expire (plan API-side authorization changes accordingly).

## Naming Convention

Use the API's HTTPS audience URI as the identifier (`https://api.example.com`). It is unique, self-describing, and matches what JWT middleware ecosystems expect an audience/scope prefix to look like. Any 1-256 character string works, but URIs keep multi-API pools legible.
