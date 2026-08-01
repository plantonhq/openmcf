# Production — OIDC sign-in on PostgreSQL

The durable, org-scale shape: sign-in delegates to your identity
provider over OIDC (a KubernetesKeycloak realm's endpoints slot
straight in — Okta, Auth0 and Dex work identically), hub state lives in
a composed KubernetesPostgres (survives volume loss, snapshots with
your database fleet), users pick their machine size from a spawn menu,
and two warm placeholder pods keep "start my server" instant while the
cluster autoscaler catches up.

The database references do all the wiring: the host resolves to the
Postgres read-write Service and the credential to the
operator-maintained application Secret, which the hub reads through a
mounted Secret — nothing password-shaped is ever written into this
manifest or into rendered chart values. What the preset expects from
the database side: the KubernetesPostgres declares `jupyterhub` as its
bootstrap database (owner `jupyterhub`).

Two things to prepare before deploying: the `kc-oauth` Secret carrying
your OAuth client secret (key `client-secret`), and the OAuth client
registration whose callback URL is exactly
`<your-hub-url>/hub/oauth_callback`.

Change first: `oauth_callback_url` and the three provider endpoints to
your real identity provider, and the profile sizings to your real
machine tiers. User homes are per-user PVCs named `claim-<username>` —
they survive this resource's destroy; plan their backup with your
storage story.

See [02-production-oidc.yaml](./02-production-oidc.yaml) for the
manifest.
