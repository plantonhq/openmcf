# KubernetesSolr: Research and Design

## Introduction

KubernetesSolr declares one Apache SolrCloud cluster as a
`solr.apache.org/v1beta1` `SolrCloud` custom resource reconciled by
the Apache Solr Operator (KubernetesSolrOperator, the registry
prerequisite — chart 0.9.1 / operator v0.9.1). Everything else — the
node StatefulSet, the provided ZooKeeper ensemble, Services, PVCs, the
basic-auth bootstrap Secrets, Ingress exposure — is the operator's to
create from the SolrCloud spec; the modules render the CR and export
the operator's deterministic names.

## The Deployment Landscape

SolrCloud pairs a Solr node fleet with a ZooKeeper quorum holding
cluster state. Operating it by hand means ordering restarts so no
shard loses all its replicas, moving replicas before scale-downs
delete their pods, and keeping the ZooKeeper wiring consistent — the
expertise the Apache Solr project encodes in its own operator. This
kind is deliberately thin: one custom resource, rendered identically
on both engines.

## Upstream Architecture

The operator reconciles the declared resource into:

- **The node StatefulSet** (`<name>-solrcloud`) — replicas, image
  (`solr` at `spec.version`), JVM (`java_mem`, `solr_opts`,
  `gc_tune`), storage, and scheduling.
- **The common Service** (`<name>-solrcloud-common`) — fronting all
  nodes on port 80 (443 with TLS); the exported `internal_endpoint`
  derives from it. A headless Service serves direct pod discovery.
- **The ZooKeeper wiring** — see below.
- **The security bootstrap Secrets** — see below.

### ZooKeeper: provided or external

Every SolrCloud needs an ensemble; the spec's `zookeeper` oneof
models the two arms:

- **Provided** (the default — an empty block renders NOTHING and the
  operator itself defaults to a provided 3-node ensemble): the
  zookeeper-operator bundled with the KubernetesSolrOperator chart
  provisions a per-cluster ensemble, reachable at the operator-named
  client service `<name>-solrcloud-zookeeper-client:2181`. Replicas
  (default 3 — the quorum minimum for production; 1 works for
  development), persistence (the zookeeper-operator's own default is a
  20Gi PVC), resources, and chroot are tunable.
- **External**: a connection string to an ensemble outside the
  operator's management, plus a chroot.

The exported `zookeeper_connection_string` reports whichever the
cluster actually uses, chroot included when it diverges from "/".

### Storage: ephemeral by default, honestly

The operator default is emptyDir — data is LOST when a pod leaves its
node (eviction, drain, node failure). That default is preserved, not
hidden: `storage.persistent` is the arm for real workloads (PVC per
node, size required, `reclaim_policy` defaulting to Retain so data
outlives the resource), and `storage.ephemeral` exists for explicit
throwaway intent with an optional size cap.

### The security bootstrap: two Secrets, one contract

`security.authentication_type: basic` (the one operator-managed type)
bootstraps security.json with THREE users carrying random passwords —
verified in the pinned operator source:

- **`<name>-solrcloud-security-bootstrap`** holds the `admin` and
  `solr` user passwords plus the bootstrapped security.json. The
  operator applies that security.json ONCE (through the setup-zk init
  container) and never updates it.
- **`<name>-solrcloud-basic-auth`** (kubernetes.io/basic-auth, fields
  `username`/`password`) holds the `k8s-oper` user — the credentials
  the OPERATOR itself uses for API requests against secured pods. This
  is the exported `basic_auth_secret_name`.

THE ROTATION CONTRACT (upstream): if a password is later rotated
through Solr's security API, the corresponding Secret must be updated
too — otherwise the operator keeps using the stale credential and
locks itself out. Bringing your own `basic_auth_secret` replaces the
generated one (the output is then exported empty, because the operator
only generates credentials when the user did not bring any);
`bootstrap_security_json` replaces the bootstrapped security.json
entirely — applied once, never updated.

`probes_require_auth` extends authentication to the probe endpoints
(operator default: probes stay open).

### TLS: keystore-based, with a client-auth caveat

Solr's TLS model is JVM keystores: `pkcs12_secret` + 
`keystore_password_secret` (JVMs refuse password-less PKCS#12), an
optional separate truststore, and `client_auth` (None/Want/Need).
`Need` demands a client certificate from every caller — INCLUDING the
operator, whose probes and reconciliation calls fail unless it carries
its own mTLS identity (the KubernetesSolrOperator `mtls` block). A
KubernetesCertificate with a pkcs12 keystore output is the natural
producer of the keystore Secret.

### Backup repositories: registration here, backups elsewhere

`backup_repositories` registers named targets on the cluster; backup
operations themselves are operational verbs (SolrBackup resources or
Solr API calls). Three backends, exactly one per repository:

- **s3** — region (required even for S3-compatible endpoints; any
  value works there), bucket, optional base location and endpoint
  (MinIO and friends). Credentials are declared Secret references, or
  an EMPTY credentials block for the keyless path: the nodes' ambient
  identity (an IRSA-bound ServiceAccount on EKS).
- **gcs** — bucket plus a service-account key JSON from a Secret, or
  keyless via GKE Workload Identity.
- **volume** — an existing PVC that MUST be multi-writer
  (ReadWriteMany or NFS), mounted to every node; `directory` defaults
  to the cluster name.

Solr modules a declared repository needs load automatically;
`solr_modules` adds the rest (`analytics`, `ltr`, ...).

### Managed updates and scaling

The operator's Managed update strategy (the default) restarts pods
shard-aware and in parallel, bounded by `max_pods_unavailable`
(default 25%) and `max_shard_replicas_unavailable` (default 1) —
`StatefulSet` (plain ordinal order) and `Manual` (you delete pods) are
the opt-outs, and `restart_schedule` adds cron-scheduled restarts.
Scale operations move data by default: `vacate_pods_on_scale_down` and
`populate_pods_on_scale_up` (Solr 9.3+) both default true — replicas
are moved, not dropped.

### Exposure: the operator's own, or composed

The `external` block models the operator's built-in exposure — Ingress
(per-node paths through an ingress controller) or ExternalDNS
(per-node DNS records) under a `domain_name`, with
`use_external_address` making Solr advertise its EXTERNAL addresses
(what CloudSolrClient outside the cluster needs), and
`hide_common`/`hide_nodes` trimming the surface. For simple HTTP
access, composing a KubernetesIngress or Gateway API route over the
common service handle is equally valid and keeps exposure a
first-class graph node. In-cluster clients never need either: the
common Service is the entry point.

## Design Decisions

- **One custom resource, nothing else.** The modules render the
  SolrCloud CR and export the operator's deterministic names — the
  StatefulSet, Services, PVCs and Secrets are all operator-created.
- **Typed CR rendering on both engines.** The Pulumi module renders
  the CR with the typed crd2pulumi SDK (drift against the pinned CRD
  fails at compile time); the Terraform module applies through
  `kubectl_manifest` (alekc/kubectl), which needs no cluster
  connection at plan time — an infra chart can plan the operator and
  its clusters in one run. Unset optionals are omitted so the
  apiserver applies the CRD's own defaults; presence discipline is
  byte-for-byte identical across engines.
- **No await machinery, deliberately.** Cluster readiness depends on
  the operator (image pulls, ZooKeeper quorum, node startup) that is
  not part of applying the resource — the never-block-on-a-controller
  posture of every operator-CR kind in the catalog.
- **The CRD's own value casing is honored at the seam** — the spec's
  lowercase `basic` renders as the CRD's capitalized `Basic`; the
  int-or-string update budgets render integer strings as YAML numbers
  ("2" becomes 2) and percentages as strings ("25%").
- **The common service port is deliberately not modeled** — the
  operator's own default (80, or 443 with TLS) governs, and the
  exported endpoint carries the effective scheme and port. `pod_port`
  (default 8983) is the node listener.
- **Backup repository names follow the CRD's own pattern** —
  alphanumeric with interior dashes/underscores, max 100 characters
  (broader than a DNS label; uppercase and underscores are accepted).

## Version Pins and Naming Contracts

| What | Value | Notes |
|---|---|---|
| CR API | `solr.apache.org/v1beta1` | Reconciled by the Solr operator (pre-1.0: minor versions can change the CRD API) |
| Solr image | `solr` (Docker Hub official) at `spec.version` | Repository overridable for air-gap; the tag always comes from `version` |
| StatefulSet | `<name>-solrcloud` | |
| Common Service | `<name>-solrcloud-common` (80, or 443 with TLS) | Exported as `common_service_name` |
| Basic-auth Secret | `<name>-solrcloud-basic-auth` (the `k8s-oper` user) | Exported as `basic_auth_secret_name`; empty when security is off or user-provided |
| Bootstrap Secret | `<name>-solrcloud-security-bootstrap` (`admin`, `solr`, security.json) | Applied once, never updated by the operator |
| Provided ZooKeeper | `<name>-solrcloud-zookeeper-client:2181` | Exported in `zookeeper_connection_string` (plus chroot) |
| Pod port | 8983 (spec and operator default) | |

## IaC Twins

Pulumi (`module/solr_cloud.go`, typed crd2pulumi SDK) and Terraform
(`kubectl_manifest` + null-prune locals) render identical CR bodies —
same arms rendered and omitted, the same int-or-string coercions — and
derive the same output names blind from `metadata.name`. Keep
`locals.go`/`solr_cloud.go` and `locals.tf` in lockstep.
