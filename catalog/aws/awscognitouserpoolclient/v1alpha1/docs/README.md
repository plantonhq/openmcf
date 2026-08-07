# AWS Cognito User Pool Client -- Architecture and Design

## Overview

An app client is the OAuth 2.0 / OIDC contract between one application and a Cognito user pool. The pool is the directory; the client is how a specific application is allowed to use it: which grants it may run, where it may redirect, how long its tokens live, and whether it holds a secret.

It is a first-class resource because it has everything the honest-decompose test looks for: a pool serves many applications (many-per-parent), the client ID is referenced by downstream systems (JWT authorizer audiences, ALB authentication actions, application configs), and clients rotate on their own schedule without touching the pool or its users.

## Public vs Confidential -- Decided at Creation

`generate_secret` is ForceNew. A **public** client (SPA, mobile app) cannot protect a secret, so it authenticates with PKCE; a **confidential** client (server-side app, machine-to-machine service) holds the AWS-minted secret and presents it at the token endpoint. Changing your mind means replacing the client -- applications get a new client ID.

## Grant Types

- **Authorization Code** (`code`) -- the recommended grant for anything with a user in front of it; pair with PKCE for public clients.
- **Client Credentials** (`client_credentials`) -- machine-to-machine. It requires a secret, can ONLY request custom scopes minted by a resource server (built-in OIDC scopes describe a user, and there is none), and AWS rejects mixing it with user-facing grants on one client -- create a separate M2M client.
- **Implicit** (`implicit`) -- legacy; tokens leak into browser history. Modeled for completeness, avoid for new applications.

## Token Model

Cognito issues three JWTs: the ID token (identity claims), the access token (scopes and group membership -- what APIs authorize on), and the refresh token. Each lifetime pairs a value with a unit; AWS bounds the results (access/ID: 5 minutes to 24 hours; refresh: 60 minutes to 10 years) and the spec enforces those bounds at validation time in whatever unit is chosen.

**Refresh-token rotation** (`refresh_token_rotation`) makes each refresh a one-time-use event: the exchange issues a new refresh token and retires the presented one, with an optional 0-60s grace window for clients that lose the response. It shrinks the blast radius of a stolen refresh token from "the full validity period" to "one exchange".

## Federated Sign-In Ordering

`supported_identity_providers` lists what the client offers at sign-in: the literal `COGNITO` for the pool's own directory plus identity-provider names. The entries are references: pointing them at `AwsCognitoIdentityProvider` resources gives the deployment graph the correct ordering guarantee -- the provider exists before any client lists it -- while literal names remain available for providers managed outside the graph. When the field is omitted, AWS enables all of the pool's providers for the client.

## Security Posture

- `prevent_user_existence_errors: ENABLED` returns identical errors for wrong-password and no-such-user, closing the user-enumeration probe.
- `enable_token_revocation` (AWS default true) lets sign-out actually revoke tokens.
- `enable_propagate_additional_user_context_data` forwards the END USER's IP and user-agent to threat protection when a backend authenticates on the user's behalf (otherwise Cognito only sees the server's IP); it requires a client secret because only confidential clients can assert it truthfully.

## Composition

The `client_id` output is the join key of the serverless auth story:

- **API Gateway JWT authorizer**: `audiences` entries reference `status.outputs.client_id`; the `issuer` comes from the pool.
- **ALB `authenticate_cognito`**: `user_pool_client_id` references the same output.
- **Application config**: SDK sign-in embeds the client ID (and, for confidential clients, the `client_secret` output -- sensitive).
