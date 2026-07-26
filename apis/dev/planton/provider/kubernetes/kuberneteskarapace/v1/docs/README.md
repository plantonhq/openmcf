# KubernetesKarapace: Research and Design

## Introduction

KubernetesKarapace declares a Karapace schema registry — Aiven's
Apache-2.0 implementation of the Schema Registry API — as
module-owned Kubernetes manifests: a registry Deployment and Service,
an optional REST-proxy Deployment and Service, and an optional
module-materialized SASL Secret. There is no Helm chart and no
operator upstream; the module IS the deployment story, with every
meaningful configuration surface typed on the spec.

## The Engine Decision: Why Karapace

The schema-registry role in an open catalog has one incumbent
(Confluent Schema Registry) and one serious open alternative. The
evidence pointed at Karapace:

- **License.** Confluent SR is under the Confluent Community License
  (not open source); Karapace is Apache-2.0 — the same license as
  this catalog.
- **Drop-in compatibility.** Upstream states Karapace is a drop-in
  replacement for both the Schema Registry and the Kafka REST proxy
  APIs, compatible with Schema Registry 6.1.1 at the API level —
  existing SR client libraries, Connect converters, and consoles work
  unchanged. (Normalization support has gaps versus Confluent —
  currently Protobuf-only upstream.)
- **The same architecture.** Schemas live in a compacted Kafka topic
  (`_schemas` by convention) on the connected cluster — exactly how
  Confluent SR stores them. No database appears in the stack, and
  operational intuition transfers.
- **Leader/replica built in.** Replicas coordinate leadership through
  a Kafka consumer group; followers forward writes to the leader.
  Availability comes from `replicas: 2`, not from external
  coordination.

The pinned image is upstream's own `ghcr.io/aiven-open/karapace` at
the 6.2.1 release; `spec.image` overrides for mirrors and pins.

## Module-Owned Manifests: How Configuration Reaches the Engine

Karapace reads configuration from environment variables with the
`KARAPACE_` prefix (the pydantic settings mechanism — config key X
becomes `KARAPACE_<uppercased X>`). The container entrypoints come
from upstream's own deployment reference (the production image
declares no entrypoint of its own): `python3 -m karapace` for the
registry role, `python3 -m karapace.kafka_rest_apis` for the REST
proxy, with `KARAPACE_KARAPACE_REGISTRY` / `KARAPACE_KARAPACE_REST`
selecting the served surface. File-based configuration (TLS material,
the authfile) mounts from Secrets at fixed paths under
`/etc/karapace/`, with env vars pointing inside.

Probes hit `/_health` — the engine ships that path in its
skip-auth list and the upstream image's own HEALTHCHECK curls it, so
liveness and readiness keep working when HTTP authentication is
enabled. Probe scheme follows `server_tls` (HTTPS probes skip
certificate verification — what makes cert-manager-issued
certificates with Service SANs work).

## The Pod-IP Advertised Hostname

The subtlest design point in the module. The leader coordinator
publishes `advertised_protocol://advertised_hostname:port` through
the consumer group as the master URL, and followers forward writes to
it. Upstream's compose reference gives every instance its OWN
identity — one advertised hostname per container. Two wrong answers
and the right one, on Kubernetes:

- **A shared Service name** would make followers forward writes to
  themselves (the Service load-balances across all replicas,
  including the follower doing the forwarding).
- **The pod's bare name** does not resolve in cluster DNS for
  Deployment pods (no headless-Service/StatefulSet contract here).
- **The POD IP** — injected via the downward API
  (`status.podIP`) — is each pod's directly-reachable self-identity:
  the Kubernetes twin of compose's container hostname.

The engine's own fallback makes explicitness mandatory: unset,
`advertised_hostname` falls back to `host`, which the module sets to
`0.0.0.0` for serving — an unroutable advertisement.

### The server_tls × replicas caveat

The pod-IP design has one consequence, carried as a caveat ON THE
SPEC: with `server_tls`, followers forward writes to the leader at
`https://<pod-ip>:<port>`, and a certificate issued for a DNS name
does not cover pod IPs. `server_tls` therefore pairs with
`replicas: 1`; multi-replica TLS postures terminate TLS at an
Ingress/Gateway in front of plain-HTTP replicas. The advertised
protocol follows the serving scheme (an https server advertising http
would break forwarding) — the module owns that coupling.

## The Kafka Connection

`security_protocol` is typed with the four Kafka client postures and
CEL-enforced pairings: SSL forms require `tls` (at minimum the CA to
trust), SASL forms require `sasl` — and, the rule that prevents the
silent failure mode: `sasl` requires the protocol EXPLICITLY set to
SASL_PLAINTEXT or SASL_SSL, because the protocol defaults to
PLAINTEXT when unset and the engine would silently ignore the
declared credentials.

TLS material mounts from Secrets (CA at `/etc/karapace/kafka-ca`,
mutual-TLS client identity at `/etc/karapace/kafka-cert`), with the
KubernetesKafka / KubernetesKafkaUser composition defaults on the
Secret names and key-name fallbacks matching Strimzi's layouts
(`ca.crt`, `user.crt`/`user.key`, `password`).

**The SASL password never rides the pod spec**: a referenced
`password_secret` wires straight into a secretKeyRef; a literal
`password` is materialized into the module-owned `<name>-sasl` Secret
first (exactly one of the two, CEL-enforced). Pod specs are readable
by anyone with get-pod RBAC; Secret values have their own ACL.

## Registry Behavior

- **`topic_name`** (default `_schemas`) — the Confluent convention
  existing tooling expects; the registry creates the topic on first
  start.
- **`replication_factor`** — applies AT CREATION only. The upstream
  default of 1 is fine for dev and a data-loss risk in production
  (the schemas topic is the registry's entire state); set 3 on
  multi-broker clusters. Changing it later means reassigning the
  existing topic with Kafka tooling, not editing the field — the spec
  says so explicitly.
- **`compatibility`** (default BACKWARD) — the default mode for new
  subjects; the _TRANSITIVE variants check against ALL prior
  versions. Per-subject overrides ride the standard SR config API at
  runtime.
- **`group_id`** (default `metadata.name`) — the leader-election
  consumer group. Unique per installation sharing a Kafka cluster:
  two installations sharing a group id fight over leadership and
  corrupt each other's view of the schemas topic.
- **`master_election_strategy`** (`lowest`/`highest`) — member
  ordering for election; an upstream knob typed for completeness.
- **`log_level`** — the upstream default is DEBUG, too noisy for
  production; the spec defaults to INFO.

## HTTP-Layer Authentication

`basic` XOR `oidc` (CEL-enforced):

- **`basic`** mounts a Karapace authfile (a JSON users/permissions
  document) from a Secret; the engine hot-reloads it on change, so
  rotating credentials is a Secret update with no restart.
- **`oidc`** validates JWT bearer tokens against the IdP's JWKS
  endpoint. The spec requires an `https://` JWKS URL — a plain-HTTP
  JWKS source would let an in-path attacker substitute signing keys
  (the upstream engine refuses it outside dev overrides). Optional
  issuer/audience claim checks.

Both leave `/_health` unauthenticated (the engine's skip-auth
list) — the probes keep working.

## The REST-Proxy Role

`rest_proxy.enabled` deploys the SAME image as a second Deployment
(`<name>-rest`) with the role flags flipped: Kafka's REST proxy API
(produce/consume/admin over HTTP), wired to the registry Service for
schema lookups (scheme follows the registry's `server_tls` posture)
and to the same Kafka cluster with identical connection env and TLS
mounts. Independently sized (replicas/port/resources); the exported
`rest_proxy_endpoint` is empty when the role is off — an honest
signal for composition. The REST proxy always serves plain HTTP
(`server_tls` covers the registry API only).

## Design Decisions

- **Native manifests, not a packaged chart** — upstream ships no
  chart or operator, so the module renders Deployments/Services
  directly on both engines (Pulumi SDK resources, Terraform
  `kubernetes_*_v1` resources) — typed, diffable, no third-party
  packaging dependency.
- **Per-role selector identity** — both Deployments run the same
  image in the same namespace; a role-specific `app` label
  (`<name>` / `<name>-rest`) keeps each Service from selecting the
  other role's pods.
- **Registry-scoped scheduling** — `node_selector`/`tolerations`
  apply to the registry pods per the spec contract; the REST-proxy
  role carries only replicas/port/resources.
- **Probes tolerate topic replay** — the engine replays the schemas
  topic at startup; initial delay and a generous failure threshold
  keep liveness from restarting a warming registry.

## Deliberately Unmodeled

Karapace exposes more knobs than a deployment spec should carry.
Because the module owns the manifests (no helm_values escape hatch
exists on this kind), exclusions are real exclusions — the raw
alternative is deploying the engine yourself via
KubernetesDeployment/KubernetesManifest:

- **statsd / sentry / telemetry knobs** — upstream's metrics
  emission, Sentry error reporting and OTel wiring are
  operational-tooling choices, not deployment posture; the
  Prometheus-native path for this catalog is cluster-level scraping.
- **Name-strategy fine-tuning** — subject-name strategy behavior
  belongs to producers/clients and per-subject SR API configuration,
  not the registry deployment.
- **`master_eligibility` per pod** — upstream can pin
  producer-only replicas; meaningless under one spec that sizes a
  homogeneous replica set.
- **Schema-reader/cache tuning** — internal performance knobs with
  sane engine defaults; typing them would freeze implementation
  detail into the API.

## Version Pins and Naming Contracts

| What | Value | Notes |
|---|---|---|
| Image | `ghcr.io/aiven-open/karapace:6.2.1` | Upstream's own image; `spec.image` overrides |
| API compatibility | Schema Registry 6.1.1 (API level) | Upstream's stated compat target |
| Registry Deployment/Service | `<name>` | Endpoint exported as `endpoint` |
| REST-proxy Deployment/Service | `<name>-rest` | `rest_proxy_endpoint`; empty when off |
| SASL Secret (literal passwords) | `<name>-sasl` (key `password`) | Module-materialized |
| Schemas topic | `_schemas` (default) | Exported as `schemas_topic` |
| Leader-election group | `metadata.name` (default) | Unique per installation per Kafka cluster |
| Health endpoint | `/_health` | Unauthenticated (skip-auth list); probed by both roles |
| Ports | 8081 registry / 8082 REST proxy (defaults) | |

## IaC Twins

Pulumi (`module/deployment.go`, `service.go`, `secret.go`) and
Terraform (`main.tf`: `kubernetes_deployment_v1.registry`,
`kubernetes_service_v1.registry`,
`kubernetes_deployment_v1.rest_proxy`,
`kubernetes_service_v1.rest_proxy`,
`kubernetes_secret_v1.sasl_password`) render the same env-var sets,
mount paths, probes, and naming contracts. Mount paths and the env
mechanism must stay byte-identical across engines.
