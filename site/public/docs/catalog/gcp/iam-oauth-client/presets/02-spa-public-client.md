---
title: "SPA Public Client"
description: "The browser-app shape: a public client using the authorization code flow with PKCE. No credentials exist because none can — a secret shipped in a browser bundle is public by definition, and GCP..."
type: "preset"
rank: "02"
presetSlug: "02-spa-public-client"
componentSlug: "iam-oauth-client"
componentTitle: "IAM OAuth Client"
provider: "gcp"
icon: "package"
order: 2
---

# SPA Public Client

The browser-app shape: a public client using the authorization code flow
with PKCE. No credentials exist because none can — a secret shipped in a
browser bundle is public by definition, and GCP refuses to attach
credentials to a `PUBLIC_CLIENT`.

## What it configures

- `clientType: PUBLIC_CLIENT` — the app cannot keep a secret. Immutable.
- `AUTHORIZATION_CODE_GRANT` only — the code-plus-PKCE flow SPAs use;
  refresh tokens for public clients are a deliberate omission here.
- `openid` scope — the identity claim, nothing more.
- No `credentials` — the `client_secret` output stays empty.

## Adjust before deploying

- **allowedRedirectUris** — replace with the SPA's real callback route,
  or a ValueFromRef to the hosting resource's URL output.
- **allowedScopes** — add scopes only for APIs the SPA calls directly;
  most SPAs should call their own backend instead.

## When to choose something else

If a server-side component handles the OAuth exchange, use the **Web App
Client** preset — a confidential client with a managed credential keeps
tokens off the browser entirely.
