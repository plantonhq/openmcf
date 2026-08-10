# GcpIamOauthClient Guide

The judgment this guide protects: an OAuth client is an identity contract
between an application and Google Cloud. Its ID, type, and location are
immutable, its secrets are server-generated, and its deletion is only
soft — every casual change here is either a forced re-registration or a
30-day wait.

## Rotation is two applies by API design

GCP refuses to delete an ENABLED credential — a hard API rule, not a
module choice. The rotation story: add a second credential (`credentials`
gains an entry), cut consumers over to the new secret, then set the old
credential's `disabled: true` in one apply and remove its entry in the
next. Deleting the entry while it is still enabled fails the apply. The
`client_secret` output always carries the FIRST credential's secret, so
after removing the old entry the remaining credential is again the first
and downstream ValueFromRef consumers pick up the new secret without a
manifest change.

## Public or confidential — decide once

`clientType` is immutable. `CONFIDENTIAL_CLIENT` is for server-side apps
that can keep a secret; `PUBLIC_CLIENT` is for SPAs and mobile apps that
cannot — credentials cannot even be attached to one, which is the point:
a secret shipped inside a browser bundle is public by definition.
Choosing wrong means destroy-and-recreate plus re-registering every
consumer against the new client ID.

## Redirect URIs should be references, not strings

A literal redirect URI drifts the day the app's address changes — the
OAuth flow breaks with an error users see and operators cannot grep for.
Each `allowedRedirectUris` entry accepts a ValueFromRef, so point it at
the serving resource's URL output (e.g. a GcpCloudRun `url`) and the
registration follows the deployment automatically.

## Scope minimization

`allowedScopes` is an allow-list the client can never exceed. Grant what
the application actually calls — `openid`/`email`/`profile` for identity,
a service scope only if the app touches that API. A client scoped to
`cloud-platform` "for convenience" turns any token leak into full-cloud
exposure.

## Names outlive deletion

Deleted clients are soft-deleted for ~30 days and the `oauthClientId`
stays reserved for the window. Never plan to delete-and-recreate under
the same ID, and keep throwaway experiments on throwaway IDs so the good
names stay available.
