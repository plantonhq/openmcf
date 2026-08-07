# KubernetesNeo4j: Research and Design

## Introduction

KubernetesNeo4j deploys one Neo4j graph database server from the
official `neo4j` Helm chart (https://helm.neo4j.com/neo4j, pinned
2026.6.0 — chart versions track Neo4j calendar releases) as a single
Helm release named after `metadata.name`. The typed spec renders into
chart values; a declared admin password materializes as a
module-owned Kubernetes Secret the chart consumes; and the chart's
LoadBalancer default is deliberately overridden to ClusterIP.

## The Deployment Landscape

Graph databases earn their place in the stack wherever the data is
relationship-shaped: knowledge graphs that a GraphRAG retrieval layer
walks at query time, agent-memory stores that accumulate entities and
their connections across sessions, dependency and identity graphs.
Neo4j is the standard engine for those workloads — the spec models it
as such — and the deployment questions this component answers are the
Kubernetes ones: grain, credentials, exposure, storage, and TLS.

### The grain: one release, one server

The chart's own grain is a StatefulSet of exactly one pod, and the
component preserves it rather than inventing a replicas knob:

- **Community edition** (the default) is single-instance by license —
  the spec REJECTS `cluster_name` on community with a CEL rule rather
  than letting a doomed install proceed.
- **Enterprise clustering** is built by installing MULTIPLE
  KubernetesNeo4j resources that share the same `cluster_name`. The
  chart's `neo4j.name` value is what members share (it controls the
  pod selector labels, anti-affinity matching, and the shared
  `<neo4j.name>-lb-neo4j` Service name); the module renders it as
  `cluster_name` when set, else `metadata.name`. Each member is its
  own first-class resource — individually addressable, sized, and
  upgraded — and each requires `accept_license_agreement: true` (the
  chart's own shape is the STRING "yes", which the module renders
  only on the affirmative) plus a valid Neo4j license.

### Licensing posture

The Helm chart repository is Apache-2.0 (its own LICENSE). The
community engine is GPLv3 — this component REFERENCES the official
image, it never distributes the engine. Enterprise is a commercial
license, gated by the spec's `accept_license_agreement` +
`edition: enterprise` CEL rule.

## The Auth Secret Contract

The chart's credential contract is specific, and the module honors it
exactly:

- The chart's `neo4j.passwordFromSecret` names a Secret carrying ONE
  key, `NEO4J_AUTH`, whose value is `neo4j/<password>` (the `neo4j`
  admin user).
- The chart LOOKS THE SECRET UP AT TEMPLATE TIME — its helper fails
  the install when the Secret is missing or lacks the key. The Secret
  must therefore exist BEFORE the Helm release; both modules wire the
  explicit dependency (the auth Secret is created first, always).
- With the `password` arm declared, the module materializes
  `<metadata.name>-auth` with exactly that contract — marked sensitive
  in state on both engines — and because `passwordFromSecret` is set,
  the chart renders no auth Secret of its own: the module-owned Secret
  is the ONLY place the credential lands, and it never transits
  rendered chart values.
- The `existing_secret` arm points the chart at a Secret the user owns
  (same contract, same exists-before-install requirement); absent auth
  lets the chart generate a random password and log it once at first
  startup — which is why `auth_secret_name` exports empty in that
  case.

## The ClusterIP Override

The chart ships `services.neo4j.spec.type: LoadBalancer` (verified in
the pinned chart values) — on every install, that would provision a
cloud load balancer or hang Pending on clusters without one. The
module pins the type to ClusterIP unless `spec.service.type` says
otherwise: exposure composes from first-class kinds (KubernetesIngress
for HTTP/Browser, TCP routes or an explicit LoadBalancer for bolt)
over the exported service handle. In-cluster clients never touch this
Service at all — they use the chart's always-created default Service,
named after the release (the chart's fullname helper renders the
release name when no name overrides are set), which is what makes the
exported `service_name`, `bolt_endpoint` (7687) and `http_endpoint`
(7474) deterministic.

## TLS: the Key-Name Bridge

The chart's `ssl` scopes (`bolt`, `https`) mount a private key and
certificate from Secrets, reading `private.key` and `public.crt` (its
subPath defaults). cert-manager-issued Secrets — including those from
a KubernetesCertificate — carry `tls.key`/`tls.crt` instead. The
module renders both `privateKey.secretName` and
`publicCertificate.secretName` against the ONE declared scope Secret
and does NOT silently rewrite key names: bridging is explicit — mirror
the certificate data into a Secret with the chart's expected keys, or
produce a Secret carrying those keys directly. Empty `ssl` means
plaintext in-cluster, with TLS composed at the exposure layer instead.

## Sizing and Storage

- **The chart enforces a resource floor**: installs below 500m CPU /
  2Gi memory are REJECTED. The module never defaults below it — empty
  `resources` renders nothing and the chart's own defaults (1000m/2Gi)
  apply. The chart's primary resources shape is flat {cpu, memory}
  applied to both requests and limits; declared limits render into the
  full-format limits sub-map.
- **The chart requires a data-volume mode** (its templates fail
  without one), so the module ALWAYS renders the data volume:
  `dynamic` with the declared StorageClass, else `defaultStorageClass`
  — size resolved to the spec default (10Gi) either way.
- **Memory tuning is typed**: the `memory` block renders the three
  neo4j.conf keys (`server.memory.heap.initial_size`,
  `server.memory.heap.max_size`, `server.memory.pagecache.size`),
  winning over duplicates in the free-form `config` map — the typed
  block is the declared interface for those keys. Empty = Neo4j
  auto-computes from container memory (the chart default).

## Design Decisions

- **The install is blocking.** The Helm release waits for the server
  to become Ready (atomic, 600s timeout): a database that never starts
  — bad image, unschedulable pod, unbindable volume, a JVM that OOMs
  on boot — fails THIS apply, not the first driver connection. Neo4j
  recovers/upgrades store files on startup, so the budget is generous.
- **`neo4j.name` always renders** — the chart REQUIRES it (nothing
  defaults it to the release name), so the module renders
  `cluster_name`-or-`metadata.name` unconditionally.
- **Scheduling splits where the chart splits it**: `nodeSelector` at
  the chart's top level; tolerations, `podAntiAffinity` (chart default
  true — meaningful for Enterprise members, harmless standalone), and
  `priorityClassName` under `podSpec`.
- **The image override resolves the repository default**: the chart
  fails when any separated image field is set while `repository` is
  empty, so the module fills `neo4j` whenever the block renders
  anything.
- **Chart-default-matching values render only on divergence** — with
  the services.neo4j type as the one deliberate always-rendered
  override.

## Version Pins and Naming Contracts

| What | Value | Notes |
|---|---|---|
| Chart | `neo4j` at https://helm.neo4j.com/neo4j | Pinned 2026.6.0; calendar-versioned with Neo4j releases |
| Auth Secret | `<name>-auth`, key `NEO4J_AUTH` = `neo4j/<password>` | Module-materialized (password arm); must exist before the release |
| Default Service | `<name>` (the release name) | Exported as `service_name`; bolt 7687, http 7474 |
| Exposure Service | `<neo4j.name>-lb-neo4j` | Chart default LoadBalancer — pinned ClusterIP by the module |
| ssl Secret keys | `private.key` / `public.crt` | cert-manager Secrets carry `tls.key`/`tls.crt` — bridge explicitly |
| Resource floor | 500m CPU / 2Gi memory | Chart-enforced; chart defaults 1000m/2Gi |
| Data volume | 10Gi (spec default), always rendered | `dynamic` or `defaultStorageClass` mode |
| Editions | community (GPLv3 engine) / enterprise (commercial) | Chart Apache-2.0; enterprise requires license acceptance |

## IaC Twins

Pulumi (`module/values.go` + `module/secrets.go`) and Terraform
(`locals.tf` + `kubernetes_secret_v1.auth` + `helm_release`) render
identical chart values and the same auth Secret (same name, key, and
contents, sensitive in state), with the same
Secret-before-release ordering. Keep the typed-value rendering in
lockstep.
