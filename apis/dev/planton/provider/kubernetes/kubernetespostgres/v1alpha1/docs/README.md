# KubernetesPostgres: Research and Design

## Introduction

KubernetesPostgres declares one production-grade PostgreSQL cluster
reconciled by CloudNativePG — the CNCF PostgreSQL operator installed by
KubernetesCloudNativePgOperator (operator chart 0.29.0 = operator
1.30.0; the Barman Cloud plugin chart 0.7.0 = plugin v0.13.0; both pins
verified against the chart-repo index). The spec renders a
`postgresql.cnpg.io/v1` Cluster custom resource and, when backups are
declared, the companion `barmancloud.cnpg.io/v1` ObjectStore and
`postgresql.cnpg.io/v1` ScheduledBackup resources — one Planton resource
for the whole database story.

## Upstream Architecture

CloudNativePG manages PostgreSQL instances DIRECTLY — no StatefulSet in
between. Each instance is a Pod (plus its PVCs) the operator creates,
watches, and replaces itself; the **instance manager** runs as PID 1 of
every instance pod, supervising the `postgres` process, reporting
readiness, applying configuration, and executing the operator's
promotion and demotion decisions. That split — operator as control
plane, instance manager as its in-pod agent — is what makes safe
automated failover and rolling-update choreography possible.

Everything the operator creates derives from the Cluster's
`metadata.name`:

- **Instance pods** `<name>-1`, `<name>-2`, ... — one primary, the rest
  streaming replicas.
- **Traffic Services** `<name>-rw` (the current primary), `<name>-ro`
  (replicas only), `<name>-r` (any ready instance). Applications
  connect through the Services, never a pod: after a failover the `-rw`
  Service re-points automatically.
- **Credential Secrets** `<name>-app` (the application owner; also
  carries ready-made `uri`/`jdbc-uri` connection strings) and
  `<name>-superuser` (only while superuser access is enabled).
- **The ServiceAccount**, named after the cluster — the identity the
  instance pods run as, and the seam keyless backup identity annotates.

## Plugin-Based Backups vs the Deprecated In-Tree Path

CloudNativePG historically embedded Barman object-store support in the
Cluster resource (`barmanObjectStore`). That in-tree path is deprecated
upstream and scheduled for removal; its replacement is CNPG-I — the
operator's plugin interface — with the Barman Cloud plugin as the
object-store implementation. The spec models ONLY the plugin path:

- The backup block renders an **ObjectStore** resource (named after the
  cluster) holding destination, backend credentials, and WAL/base-backup
  tuning, plus the Cluster's `plugins` entry designating the plugin as
  WAL archiver — archiving starts as soon as the cluster is healthy.
- Each declared schedule renders a **ScheduledBackup**
  (`<cluster>-<schedule>`), explicitly `method: plugin` — never the
  deprecated in-tree method. Schedules use SIX-field cron (seconds
  first); at least one schedule is what makes point-in-time recovery
  real, since WAL alone cannot be replayed without a base backup to
  start from.
- A recovery bootstrap renders a SECOND ObjectStore
  (`<name>-recovery-source`), read through a synthetic
  `externalClusters` entry named `origin` that carries the source's
  `serverName` as a plugin parameter (the ObjectStore CRD forbids
  `serverName` inline). Recovery never applies a retention policy — the
  plugin must not prune the source cluster's archive. The spec warns
  against pointing backups and recovery at the same destination path:
  the new cluster would overwrite the archive it restored from.

The prerequisite falls out of the design: the operator must be installed
with `barman_cloud_plugin.enabled`, or backup blocks cannot function.

## Credential Materialization

Nothing sensitive ever appears inline in a rendered custom resource.
Every declared credential materializes as a deterministic Kubernetes
Secret, and the operator/plugin see only secretKeyRef pointers:

- `<name>-app-provided` / `<name>-superuser-provided` /
  `<name>-role-<role>` — `kubernetes.io/basic-auth` pairs CloudNativePG
  WATCHES: rotating the value rotates the database password. When
  initdb declares an owner password, the operator adopts the provided
  secret instead of generating `<name>-app` — the outputs point at the
  EFFECTIVE secret either way.
- `<name>-ext-<external-cluster>` — the password for an external server
  (the operator builds a passfile from it).
- `<name>-backup-creds` / `<name>-recovery-creds` (+ `-endpoint-ca` for
  self-signed S3-compatible endpoints, + `-region` because the
  ObjectStore CRD models the S3 region as a secret reference) — the
  object-store credentials. Keyless arms render the backend's
  ambient-identity flag (`inheritFromIAMRole` / `gkeEnvironment` /
  `inheritFromAzureAD`) and need no credential Secret at all.

Names are deterministic — never engine-generated suffixes — so both
engines agree byte-for-byte and import recipes can derive them blind.

## Bootstrap: Three Ways to Be Born

Bootstrap is immutable — it describes how the cluster came to exist, and
the operator ignores changes after the first reconcile. Exactly one
method (spec-enforced oneof):

- **initdb** — a fresh empty database; the standard path. Carries the
  application database/owner, optional declared owner password,
  data-page checksums (cannot be enabled later), encoding/locale,
  post-init SQL, and optionally **import**: a logical
  pg_dump/pg_restore from a declared external cluster — the
  cross-version, cross-architecture migration path (works from RDS,
  Cloud SQL, anything reachable). `microservice` imports one database
  into the new application database; `monolith` recreates databases and
  roles one-to-one.
- **recovery** — restore from an object-store backup: disaster recovery,
  cloning, and point-in-time recovery (at most one target selector:
  time, LSN, named restore point, or first-consistent-state).
- **pg_basebackup** — physical streaming from a live server declared as
  an external cluster: the binary-identical migration path (same major
  version, same architecture).

## Version Pins

| What | Pin | Ships |
|---|---|---|
| Operator chart (`cloudnative-pg`) | 0.29.0 | operator 1.30.0 |
| Plugin chart (`plugin-barman-cloud`) | 0.7.0 | plugin v0.13.0 |

Chart and app versions move separately; the chart pins govern and are
verified against the chart-repo index. The PostgreSQL image itself is
per-database (`image_name`; empty rides the operator's default for its
release).

## Deliberate Exclusions

The typed spec covers the surface a production database actually
exercises. The following CloudNativePG capabilities are deliberately NOT
modeled — each is reachable today by declaring the raw custom resources
through KubernetesManifest, until demand admits a typed field:

- **Tablespaces** — a specialist layout concern (separate volumes per
  tablespace); the separate WAL volume covers the common I/O-isolation
  need.
- **Replica / distributed clusters** — cross-cluster topologies
  (a cluster replicating from another cluster) multiply the spec's
  surface for a posture few installations run.
- **The Pooler (PgBouncer)** — connection pooling is its own resource
  with its own sizing story; folding it in would blur the
  one-resource-one-cluster contract.
- **LDAP authentication** — enterprise directory integration; pg_hba
  rules cover the common authentication surface.
- **Image catalogs** — the indirection (ClusterImageCatalog) solves
  fleet-wide image governance, which is a platform concern, not a
  per-database one; `image_name` pins directly.
- **Extension image-volumes** — mounting extension images requires
  operator feature gates; not a stable default posture.
- **The pod monitor knob** — deprecated upstream in favor of creating
  the PodMonitor separately; the exporter itself (port 9187) is always
  on and the `monitoring` block types its TLS/query knobs.
- **Managed services tuning** — reshaping the operator-created Services
  is an edge posture; the three-service contract is the value.
- **Ephemeral / projected volumes** — niche pod-spec plumbing.

## Exposure Is Composed

No ingress block exists, deliberately. The cluster is in-cluster
plumbing reachable at `kube_endpoint`; external exposure is a separate
first-class kind (a KubernetesService of type LoadBalancer, or a TCP
route on a Gateway) targeting the exported service names. The
certificate seam cooperates: `server_alt_dns_names` puts external
hostnames on the operator-generated certificate, or
`certificates.server_tls_secret` wires a cert-manager-issued
certificate (a KubernetesCertificate reference) with its
`server_ca_secret`.

## Render Semantics

Both engines render the same resources byte-for-byte: the optional
namespace, the credential Secrets, the ObjectStore(s), the Cluster, and
the ScheduledBackups. Unset optionals are omitted entirely so the
apiserver applies the CRD's own defaults; chart-default-matching values
(anti-affinity `preferred`, `enable_pdb` true, role ensure `present`,
connection limit −1) render only when they diverge. The Pulumi module
uses typed CRD bindings (field drift against the pinned CRD fails at
compile time); the Terraform module applies the same bodies through
server-side apply with no cluster connection needed at plan time.

## Outputs

`namespace`, `cluster_name`, the three service names (`rw_service`,
`ro_service`, `r_service`), `kube_endpoint`, `port_forward_command`,
`username_secret` / `password_secret` (pointing at the EFFECTIVE
application secret — operator-generated or module-provided), and
`superuser_secret_name` (empty unless superuser access is enabled).

## E2E

The behavioral facts are properties of the platform, not of any one test
run:

- The readiness proof is the Cluster reporting healthy instances and the
  `-rw` Service routing to the primary.
- The backup proof is plugin-era end to end: ObjectStore + Cluster
  plugin wiring + an immediate ScheduledBackup against in-cluster MinIO
  (the S3-compatible arm — declared keys, endpoint URL, no cloud
  account), waiting for the Backup to reach Completed.
- The failover proof deletes the primary pod and watches the `-rw`
  Service re-point to a promoted replica.
- Declared credentials (owner password, role passwords, superuser, MinIO
  keys) exercise every secret-materialization path.
