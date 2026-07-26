# KubernetesOpenSearch: Research and Design

## Introduction

KubernetesOpenSearch declares one OpenSearch cluster as an
`opensearch.opster.io/v1` `OpenSearchCluster` custom resource
reconciled by the OpenSearch Kubernetes Operator
(KubernetesOpenSearchOperator, the registry prerequisite — pinned at
the chart 2.8.0 / operator 2.8.0 stable pairing). One resource carries
the whole cluster story: node pools, security posture, Dashboards,
monitoring, keystore entries, snapshot repositories, and projected
volumes.

## The Deployment Landscape

OpenSearch is the Apache-2.0 fork of Elasticsearch: a drop-in
replacement for the Elasticsearch APIs at the 7.10 fork line, with its
own 2.x/3.x feature line since. The fork exists because of licensing —
the Apache-2.0 posture is what lets the engine, its clients, and its
tooling stay open end to end — and the migration story follows from
it: existing Elasticsearch clients and integrations speak to
OpenSearch unchanged at the API-compatibility line, while new work
tracks the OpenSearch feature line.

Running it on Kubernetes without an operator means hand-rolling
cluster-manager quorums, per-node TLS, security-plugin initialization,
and drain-ordered rolling upgrades. The operator encodes all of that,
so this kind is deliberately thin: it renders ONE custom resource and
exports the operator's deterministic names.

## Upstream Architecture

The operator reconciles the declared resource into:

- **Node StatefulSets per pool** (`<cluster>-<pool>`) — each pool with
  its own roles, sizing, JVM, storage shape, scheduling, and PDB.
  Roles are underscore forms (`cluster_manager`, `data`, `ingest`,
  `ml`, `remote_cluster_client`, `search`); the CRD deliberately
  leaves the list open for new upstream roles, so a typo (the dashed
  "cluster-manager") travels all the way to the pod and fails only at
  node startup.
- **The main Service** — the operator names it after the CR's
  serviceName, which the modules pin to the cluster name, making the
  exported `service_name` deterministic. A `<name>-discovery` Service
  serves internal discovery.
- **A temporary bootstrap pod** that forms the initial quorum (tunable
  via `bootstrap`, rarely needed).
- **TLS Secrets and the security bootstrap** — see below.
- **The Dashboards Deployment and Service** (`<name>-dashboards`, port
  5601 — hardcoded upstream, not configurable through the CRD) when
  `dashboards.enabled`.

### The security truth

Three facts, all verified in the pinned operator source:

1. **The HTTP layer is https in EVERY posture.** The operator's own
   cluster URL builder returns `https://` unconditionally and the node
   readiness probe curls https. With `spec.security` absent, the TLS
   reconciler generates nothing and the opensearchproject image's DEMO
   security configuration serves the HTTP layer over TLS instead.
2. **The bootstrapped admin credentials are the demo credentials.**
   The operator creates the `<name>-admin-password` Secret
   unconditionally, but at this release it does not generate a random
   password — without a custom security config it uses the image's
   well-known demo admin credentials. The Secret is an honest handle,
   not a security guarantee.
3. **A custom `security.config` replaces the bootstrap entirely.** All
   three secrets are then typically required: the security-plugin YAML
   files, an admin client certificate for securityadmin.sh, and
   admin credentials for the operator's own API access (drain
   coordination, health checks — the bootstrapped credentials do not
   exist in that case). The module exports
   `admin_credentials_secret_name` as EMPTY then, because the
   user-provided secret is authoritative.

Production guidance follows: custom security config, or immediate
rotation of the admin password through the security API.

### TLS posture

`security.transport_tls` / `security.http_tls` empty-with-generate
(the default) has the operator issue a CA and per-layer certificates.
Two component-level decisions:

- **Per-node transport certificates default true** — the stronger
  posture; the operator's OWN default is a single shared certificate.
  The modules always render the field explicitly so the component
  default governs.
- **Provided certificates are the cert-manager seam** — `secret`
  references a KubernetesCertificate's output Secret; `nodes_dn` is
  REQUIRED with provided transport certificates (the security plugin
  rejects inter-node connections from unlisted DNs), and `admin_dn`
  names the DNs allowed to run securityadmin.sh.

### Dashboards: a section, not a component

OpenSearch Dashboards (the Kibana-role console) is a section of the
same custom resource. Its version defaults to the cluster's version
(module-owned defaulting — Dashboards refuses mismatched clusters, and
the CRD's field is required), TLS is a knob, `base_path` serves
path-rewriting proxies, and `opensearch_credentials_secret` is needed
only with a custom security config (the operator wires its
bootstrapped credentials otherwise).

### Snapshot repositories and the keystore

Snapshot repositories register on the cluster with `type` and
`settings` passing straight through to the snapshot API. The
credential story is deliberate:

- **Declared keys go in the KEYSTORE** — `keystore` entries load
  existing Secrets into the OpenSearch keystore on every node before
  startup, with `key_mappings` renaming Secret keys into keystore keys
  (e.g. a Secret's `accessKey` becomes `s3.client.default.access_key`
  — the operator's own documented recipe). Settings never carry
  credentials.
- **Keyless paths skip the keystore** — instance or workload identity
  on the nodes (IRSA on EKS, Workload Identity on GKE) authenticates
  the repository plugins without any declared secret.
- **The matching repository plugin must be on the nodes** —
  `repository-s3`, `repository-gcs`, ... via `plugins_list` (an
  install-at-startup download; air-gapped clusters bake a custom
  image).

## Design Decisions

- **One custom resource, nothing else.** Node StatefulSets, Services,
  TLS Secrets, the admin bootstrap and Dashboards are all
  operator-created; the modules render the CR and export the
  operator's deterministic names.
- **No ingress resources.** Exposure composes from first-class kinds
  (KubernetesIngress, Gateway API kinds) over the exported service
  handles; the Dashboards service type is the one quick-LoadBalancer
  knob, and `service_annotations` is the cloud-controller injection
  surface.
- **Typed CR rendering on both engines.** The Pulumi module renders
  the CR with the typed crd2pulumi SDK (field/structure drift against
  the pinned CRD fails at compile time); the Terraform module applies
  through `kubectl_manifest` (alekc/kubectl), which needs no cluster
  connection at plan time — an infra chart can plan the operator and
  its clusters in one run, before the CRDs exist. Unset optionals are
  omitted entirely so the apiserver applies the CRD's own defaults —
  presence discipline is kept byte-for-byte identical across engines.
- **No await machinery, deliberately.** Cluster readiness depends on
  the operator (image pulls, quorum bootstrap, security
  initialization) that is not part of applying the resource — the
  never-block-on-a-controller posture of every operator-CR kind in the
  catalog.
- **`set_vm_max_map_count` defaults true** — the CRD/operator default
  is FALSE, but most distributions ship a kernel default below what
  OpenSearch requires and pods crash-loop without the privileged init
  container; the component default is rendered explicitly. Disable
  only on pre-tuned nodes or where privileged init containers are
  forbidden.
- **The image override folds into the CRD's single image string**
  (`repo:tag`), and a private image's `pull_secret_name` joins
  `image_pull_secrets` (deduplicated) — a private image travels with
  its own pull secret.
- **`https` is exported honestly** — the endpoint scheme never
  pretends a plaintext path exists.

## Version Pins and Naming Contracts

| What | Value | Notes |
|---|---|---|
| CR API | `opensearch.opster.io/v1` | The stable operator's API group (the 3.x operator line migrates to `opensearch.org`) |
| Operator | 2.8.0 (via KubernetesOpenSearchOperator) | Supports OpenSearch 2.19.x through the 3.x line |
| Node image | `opensearchproject/opensearch` at `spec.version` | Overridable for air-gap |
| Main Service | `<name>` (serviceName pinned to the cluster name) | Exported as `service_name` |
| Admin Secret | `<name>-admin-password` (fields `username`/`password`) | Demo credentials without a custom security config; exported empty with one |
| Dashboards Service | `<name>-dashboards`, port 5601 | Port hardcoded upstream |
| StatefulSets | `<cluster>-<pool>` | Pool names: DNS-safe, max 20 chars |
| HTTP port | 9200 (spec and CRD default) | Changing it changes every advertised endpoint |

## IaC Twins

Pulumi (`module/cluster.go`, typed crd2pulumi SDK) and Terraform
(`kubectl_manifest` + null-prune locals) render identical CR bodies —
same keys rendered and omitted, the same explicit component defaults
(per-node transport certs, vm.max_map_count) — and derive the same
output names. Keep `locals.go`/`cluster.go` and `locals.tf` in
lockstep.
