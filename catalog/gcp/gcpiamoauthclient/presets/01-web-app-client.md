# Web App Client

The standard server-side shape: a confidential client with one managed
credential whose secret GCP generates — the app authenticates users
through the code flow and keeps sessions alive with refresh tokens.

## What it configures

- `clientType: CONFIDENTIAL_CLIENT` — the app runs server-side and can
  keep a secret. Immutable.
- Both grant types: `AUTHORIZATION_CODE_GRANT` for the login flow,
  `REFRESH_TOKEN_GRANT` for long-lived sessions.
- `cloud-platform` and `openid` scopes — API access plus the identity
  claim.
- One credential named `primary`; its server-generated secret is the
  `client_secret` output. Wire it to consumers via ValueFromRef (e.g.
  into a GcpSecretManagerSecret) — never copy-paste it.

## Adjust before deploying

- **allowedRedirectUris** — replace with the app's real callback URL, or
  better, a ValueFromRef to the serving resource's URL output so the
  registration follows the deployment.
- **allowedScopes** — trim `cloud-platform` if the app only needs
  identity; a broad scope turns any token leak into broad exposure.
- Removing the `primary` credential later takes two applies: set
  `disabled: true` first, then delete the entry — GCP refuses to delete
  an enabled credential.

## When to choose something else

If the application runs in a browser or on a device (no place to hide a
secret), note that GCP currently rejects public-client creation
("Client type is not supported") — SPAs go through a confidential
backend-for-frontend instead; public clients cannot
carry credentials, by design.
