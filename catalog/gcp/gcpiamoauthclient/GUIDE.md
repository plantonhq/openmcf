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

## Only confidential clients exist (GCP's rule, not ours)

`clientType` is immutable, and only `CONFIDENTIAL_CLIENT` (server-side
apps that can keep a secret) can actually be created: GCP's API lists
`PUBLIC_CLIENT` in its enum but rejects creating one with "Client type
is not supported" — verified at the raw API, no field combination
unlocks it. SPAs and mobile apps that need a public OAuth client have
no workforce-federation path today; when GCP ships support, the spec's
validation re-admits the value. A wrong choice would mean
destroy-and-recreate plus re-registering every consumer against the new
client ID — which today cannot arise, since there is exactly one legal
choice.

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
