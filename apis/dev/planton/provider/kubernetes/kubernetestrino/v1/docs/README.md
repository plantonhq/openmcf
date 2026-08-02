# KubernetesTrino — design notes

## Distribution

The official trinodb Helm chart (https://trinodb.github.io/charts),
chart 1.42.2 = Trino 480 (Apache-2.0). The module pins
`fullnameOverride` to the resource name, so child names are
deterministic: `<name>-coordinator`, `<name>-worker`,
`<name>-catalog`, `<name>-schemas-volume-coordinator`.

## The credential architecture

Every properties surface in this chart — `config.properties`, every
catalog, access control, resource groups, the exchange manager —
renders into ConfigMaps. Nothing credential-bearing may therefore ride
values. Trino's own secrets mechanism (`${ENV:VAR}` substitution,
supported in ALL properties files including catalogs) is the delivery
path:

- The internal-communication shared secret (REQUIRED once
  authentication is on) lives in the module-owned `<name>-internal`
  Secret and reaches `config.properties` as
  `${ENV:TRINO_INTERNAL_SHARED_SECRET}`.
- Each typed catalog's `connection-password` renders as
  `${ENV:TRINO_CATALOG_<NAME>_PASSWORD}`, delivered by a secretKeyRef
  env var pointing at the referenced credential Secret.
- The admin credential is one random: the `<name>-auth` Secret's
  `password` key holds the plaintext and `password.db` holds the
  bcrypt htpasswd line the chart mounts through
  `auth.passwordAuthSecret` — a verified pairing by construction.

## Authentication truths (server source at the pin)

- PASSWORD (file) auth requires HTTPS for clients unless
  `http-server.authentication.allow-insecure-over-http=true` — the
  module sets it when the `https` arm is off because in-cluster
  traffic rides ClusterIP and TLS terminates at composed exposure.
- The chart auto-renders `password-authenticator.properties` when the
  authentication type contains PASSWORD; the password file mounts at
  `<config-path>/auth/password/password.db`.

## Sizing truths

- The chart's JVM default is a FIXED 8G max heap. `max_heap_percent`
  (rendered as `-XX:MaxRAMPercentage`) requires the fixed `-Xmx` to be
  UNSET — the module null-deletes it when percent is chosen.
- `query.max-memory` (cluster-wide) and `query.max-memory-per-node`
  take SI data sizes (`4GB`, not `4Gi`) — the spec patterns enforce
  the Trino grammar.
- Graceful shutdown requires the pod termination budget to be at least
  twice the drain window — the module sets exactly 2×.

## Name budgets (chart templates at the pin)

`<name>-schemas-volume-coordinator` (27-char suffix) renders
unconditionally → names ≤ 36 chars;
`<name>-resource-groups-volume-coordinator` (36-char suffix) renders
when resource groups are declared → names ≤ 27 chars then. Both
checked fail-loud in both engines.

## Pins

Chart 1.42.2 = Trino 480; image `trinodb/trino:480` (split
registry/repository form, tag always pinned explicitly). The JMX
exporter sidecar default is `bitnamilegacy/jmx-exporter:1.4.0` — a
frozen legacy mirror; override `metrics.exporter_image` on production
mirrors.
