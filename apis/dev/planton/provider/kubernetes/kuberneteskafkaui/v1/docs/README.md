# KubernetesKafkaUi: Research and Design

## Introduction

KubernetesKafkaUi deploys kafbat UI — the Apache-2.0, community-run
web console for Kafka (the actively maintained continuation of the
provectus kafka-ui project) — from its official Helm chart. The spec
types the console's meaningful surface: multi-cluster wiring with
sibling foreign keys (Kafka bootstrap, schema registry, Connect REST
endpoints), per-cluster TLS/SASL connections, the observe-only
switch, the single login-form account, sizing and exposure knobs, and
the `helm_values` escape hatch.

## The Chart, and Where Versions Come From

The module installs the SERVED chart — `kafka-ui` from
https://ui.charts.kafbat.io — pinned at 1.6.4, which ships app
version v1.5.0. Chart versions are picked from the served repository
index, not the source tree's Chart.yaml (the tree moves ahead of
what the repository actually serves). Chart and app versions move
separately; the chart pin governs. The container image comes from the
chart's default (`ghcr.io/kafbat/kafka-ui` at the chart's
appVersion); only the REGISTRY seam is typed (`image_registry`, for
air-gapped mirrors) — pinning a divergent app image against a pinned
chart is a drift generator, so it deliberately is not.

The release name and chart fullname are pinned to `metadata.name`
(`fullnameOverride`) — the catalog's Helm-kind convention: outputs
stay deterministic, several consoles coexist in one cluster, and
exposure kinds address the Service by the resource name alone.

## One Console, Many Clusters

kafbat UI is multi-cluster by design: the app config carries a
`kafka.clusters` list, and each entry is independent — its own
bootstrap, security posture, schema registry, and Connect clusters.
The spec mirrors that shape as `clusters[]`, with unique names (the
console's display and API identifiers) and three sibling seams per
entry:

- `bootstrap_servers` ← KubernetesKafka
  (`internal_bootstrap_endpoint`)
- `schema_registry.url` ← KubernetesKarapace (`endpoint`) — Karapace
  serves the Confluent SR API the console speaks
- `kafka_connect[].address` ← KubernetesKafkaConnect
  (`rest_api_endpoint`)

`read_only` renders the app's per-cluster observe-only switch: every
mutating action (topic create/delete, message produce, config edits)
disappears from the UI and API for that cluster. It is an app-side
switch, not a Kafka ACL — the right posture for production clusters
on a shared console, and a complement to (not a substitute for)
least-privilege SASL credentials.

## The Single-User LOGIN_FORM Fact

The typed `auth` models exactly one account, and that is an upstream
fact, not a modeling shortcut: with `auth.type: LOGIN_FORM` the app's
form login authenticates against Spring Boot's DEFAULT security user
— verified in the app source: `io.kafbat.ui.config.auth.
BasicAuthSecurityConfig` enables form login and registers NO user
store of its own, so the only credential is
`spring.security.user.name` / `spring.security.user.password`. A
repeated `users` list in the spec would validate and then silently
authenticate only one of its entries — the reason `auth.user` is a
single object (CEL admits only `login_form`; the message says where
multi-user lives).

Multi-user, OAuth2/OIDC, LDAP and console RBAC are real upstream
features — configured through the app's Spring configuration, which
rides `helm_values` (`yamlApplicationConfig`'s auth/rbac trees).
Deliberately untyped: that surface is large, moves with upstream, and
composes badly with a portable spec.

## The Placeholder / secretMappings Mechanism

The chart writes the app configuration verbatim into a ConfigMap
(`config.yml`). Anything rendered there is world-readable to anyone
who can read ConfigMaps, lands in Helm release history, and in both
engines' state files — so the module NEVER renders a credential value
into the application config. Instead:

1. Every password position in the rendered config carries a
   `${ENV_VAR}` placeholder — the SASL JAAS line, schema-registry and
   Connect basic-auth passwords, the Spring login password.
2. The chart's `envs.secretMappings` wires each env var to a
   Kubernetes Secret key (rendered as `valueFrom.secretKeyRef` in the
   Deployment).
3. kafbat UI is a Spring Boot app: Spring's property resolution
   expands `${ENV_VAR}` in the mounted config against the container's
   environment at startup — the credential exists only inside the
   running container.

Two credential kinds feed the mappings: REFERENCED credentials
(cluster SASL, registry, Connect `password_secret` entries) map
straight to their source Secrets — the module copies nothing; the one
LITERAL credential (the console login password) is materialized into
the module-owned `<name>-secrets` Secret (key
`console-user-password`). Env-var names are deterministic and
index-based (`KAFKA_CLUSTER_<i>_PASSWORD`,
`KAFKA_CLUSTER_<i>_SCHEMA_REGISTRY_PASSWORD`,
`KAFKA_CLUSTER_<i>_CONNECT_<j>_PASSWORD`, `KAFKA_UI_USER_PASSWORD`)
so both engines emit identical placeholders.

## The PEM Truststore Mechanism

Kafka clients accept `ssl.truststore.type=PEM` with
`ssl.truststore.location` pointing at a plain PEM certificate file
(KIP-651) — no JKS/PKCS12 conversion, no truststore password. The
module exploits that: each TLS cluster's CA Secret mounts as-is at an
index-named path (`/etc/kafkaui/cluster-<i>-ca`), and the rendered
client properties point the truststore at the mounted key — so a
Strimzi cluster-CA Secret (`ca.crt`) works directly by reference.
The module also owns the derived `security.protocol`
(SSL / SASL_PLAINTEXT / SASL_SSL from the tls/sasl combination) and
the JAAS line (ScramLoginModule for SCRAM mechanisms,
PlainLoginModule for PLAIN); the spec forbids credentials in
`properties` — module-owned security entries win.

## Design Decisions

- **Wait-for-ready, atomic installs.** Unlike the catalog's
  operator-CR kinds (which never block on a controller), a Helm
  release of a stateless console SHOULD gate on readiness: a console
  that never starts — bad image, unschedulable pod, unresolvable
  cluster config — fails the deploy, not the first browser hit. Both
  engines set wait + atomic with a 600s timeout.
- **helm_values merges LAST, Helm `-f` semantics** — maps deep-merge
  with the override winning, lists replace; identical on both engines
  (the Terraform provider does this natively; the Pulumi module
  implements the same merge). The full-surface development manifest
  proves an override winning over a typed field.
- **The console Secret exists only when auth does** — no auth, no
  Secret, no mappings; the app renders `auth.type: DISABLED`
  explicitly.
- **Volumes/volumeMounts pass through the chart verbatim** — the
  chart forwards them to the Deployment, which is what makes the
  PEM-mount mechanism chart-native rather than a post-render hack.

## Deliberately Unmodeled

All reachable through `helm_values` (the chart/app surface), none
through typed fields:

- **OAuth2/OIDC, LDAP, multi-user, console RBAC** — the app's Spring
  security configuration; large, upstream-moving, and already
  documented upstream in chart-values terms.
- **Ingress** — the chart ships an ingress block; this catalog
  composes exposure from first-class kinds against the exported
  service handles instead, with credentials in place first.
- **Per-cluster mTLS client certificates (keystores)** — the typed
  TLS surface is CA TRUST (the overwhelmingly common posture for a
  console); client-certificate keystores need JKS material and
  passwords, and a console holding produce-capable mTLS identities is
  an anti-pattern — SASL credentials scope better.
- **Probes, security contexts, extra volumes, annotations** — chart
  passthroughs with sane defaults.

## Version Pins and Naming Contracts

| What | Value | Notes |
|---|---|---|
| Chart | `kafka-ui` 1.6.4 at https://ui.charts.kafbat.io | Served-index pin; ships appVersion v1.5.0 |
| Image | chart default (`ghcr.io/kafbat/kafka-ui`) | Only the registry seam is typed (`image_registry`) |
| Release / fullname | `metadata.name` | `fullnameOverride`; the Service IS the resource name |
| Console Secret | `<name>-secrets` (key `console-user-password`) | Only when `auth` is declared |
| CA mounts | `/etc/kafkaui/cluster-<i>-ca` | Index-named, one per TLS cluster |
| Placeholder env vars | `KAFKA_CLUSTER_<i>_...`, `KAFKA_UI_USER_PASSWORD` | Deterministic across engines |

## IaC Twins

Pulumi (`module/helm_release.go` + `secret.go`) and Terraform
(`main.tf`: `kubernetes_secret_v1.console` +
`helm_release.kafka_ui`, values rendered in `locals.tf`) emit the
same chart values documents: the same `yamlApplicationConfig`
(cluster entries, placeholders, auth tree), the same secretMappings,
the same TLS volumes/mounts, the same wait/atomic posture, with
`helm_values` merged last on both. Keep them in lockstep.
