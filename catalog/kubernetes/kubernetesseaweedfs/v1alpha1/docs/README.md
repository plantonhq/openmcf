# KubernetesSeaweedFs: Research and Design

## Introduction

KubernetesSeaweedFs deploys a SeaweedFS object store from the official
`seaweedfs` Helm chart (https://seaweedfs.github.io/seaweedfs/helm,
pinned 4.40.0 — chart versions track SeaweedFS releases; 4.40.0 ships
appVersion 4.40) as a single Helm release named after `metadata.name`.
The typed spec renders into chart values; the S3 gateway is on by
default with auth on; every data volume is deliberately mapped to a
PersistentVolumeClaim instead of the chart's hostPath default; and
everything stays ClusterIP — exposure composes from first-class kinds.

SeaweedFS is Apache-2.0 end to end — the engine, the image
(`chrislusf/seaweedfs`), and the chart — so there is no licensing gate
anywhere in the spec.

## SeaweedFS Architecture

### Masters, volumes, filer

SeaweedFS separates cluster metadata from object data, and the
separation is the topology:

- **Master servers** coordinate the cluster: they track which logical
  volumes exist, which volume servers host them, and they assign file
  ids on write. Masters form a Raft quorum — one master is a working
  dev store, three (an odd count; an even count wastes a node without
  adding failure tolerance) is the HA shape. Master state is small;
  the 5Gi PVC default is generous.
- **Volume servers** store the actual object bytes. This is the only
  tier that grows with stored data, and the tier `replication`
  placement codes spread copies across. Each volume server announces
  its capacity to the masters and serves reads/writes directly.
- **The filer** provides the file/bucket namespace on top of the raw
  volume layer — a path-to-file-id mapping held in a metadata store —
  and hosts the S3 API. Buckets are directories under the filer's
  `/buckets` folder (the chart's `WEED_FILER_BUCKETS_FOLDER` default).

Each tier scales independently through its own `replicas` field: more
masters for quorum resilience, more volume servers for capacity and
replica placement, a dedicated S3 gateway Deployment for API
throughput. The chart wires the tiers together itself (its global
`WEED_CLUSTER_SW_MASTER` / `WEED_CLUSTER_SW_FILER` env defaults
resolve the in-release addresses); nothing in the spec needs to.

### The needle/volume storage model

SeaweedFS descends from the Haystack design: instead of one file per
object on the underlying filesystem, objects ("needles") are appended
into large pre-allocated volume files, and each needle is addressed by
a compact file id — volume id, needle key, cookie. The master assigns
file ids; reading an object is one lookup in the volume server's
needle index plus one disk seek into the volume file. The index is
small enough to keep in memory (the default `index_mode`), which is
what makes reads O(1) disk operations even for billions of small
files — the classic failure mode of one-file-per-object stores.

Two spec fields tune this model directly:

- `master.volume_size_limit_mb` (default 1000, the chart's default)
  caps each volume file's size before the master assigns a new one —
  raise it for large-object stores.
- `volume.max_volumes` caps volume files per server; 0 (the default)
  auto-sizes from free disk, which is the right answer almost always.

### Why it is light where Ceph-class systems are heavy

A Ceph-class system carries cluster machinery sized for
exabyte-scale, multi-protocol storage: monitor quorums, manager
daemons, per-disk OSD daemons, placement-group rebalancing. SeaweedFS
is a single Go binary per role with a deliberately small metadata
plane — masters track volumes (thousands), not objects (billions), so
the coordination burden stays constant as the store grows. The
practical consequence for this catalog: a working S3 store fits in
three small pods on default PVCs, and the HA shape is still just
StatefulSets with placement codes — no CRD forest, no dedicated
storage operators, no minimum node count beyond what the replication
code itself demands.

## The Deployment Landscape on Kubernetes

Three ways to run SeaweedFS on Kubernetes exist in the wild:

- **The official Helm chart** (`seaweedfs` at
  https://seaweedfs.github.io/seaweedfs/helm, maintained in the
  upstream repository under `k8s/charts/seaweedfs`) — templates every
  tier, the S3 gateway in both shapes, the credential Secret, the
  bucket-creation hook, the admin console, ServiceMonitors, and the
  optional arms (SFTP, maintenance worker, COSI, all-in-one). This is
  the canonical deployment path and what this component installs.
- **Operators** — third-party and not part of the upstream project's
  release train; none carries the chart's surface or its maintenance
  cadence.
- **Raw manifests** — re-derive by hand exactly what the chart
  already templates (per-tier StatefulSets, peer discovery env,
  security wiring), with none of the upgrade story.

The chart is the upstream-maintained artifact with the full feature
surface, so the component builds on it and keeps the release the unit
of ownership: one release, named `metadata.name`, per resource.

## The S3 Surface

### On by default

This kind exists to serve S3, and the spec encodes that:
`s3.enabled` and `s3.enable_auth` are optional bools DEFAULTING TO
TRUE — an empty spec gets an authenticated S3 endpoint. Setting
`enabled: false` explicitly produces a pure filer/POSIX deployment.

The chart's own `s3.enabled`/`filer.s3.enabled` defaults are false;
the component's modules flip the right one based on the declared
shape:

- **Embedded (default)**: the gateway serves from the filer pods
  (`filer.s3.enabled`) — one fewer Deployment, metadata and API share
  fate and sizing.
- **Dedicated (`s3.dedicated`)**: the gateway runs as its own
  Deployment (`s3.enabled`) with its own replicas and resources —
  scale the API independently of metadata.

Both shapes expose the same `<name>-s3` Service on port 8333, so the
exported `s3_endpoint` is stable across the shape change. The chart's
s3-secret and bucket hook read `enableAuth` and
`existingConfigSecret` from BOTH the `s3.*` and `filer.s3.*` paths,
so the modules render the common keys on both — only the enabled
flags differ.

### The credential model

With auth on, the chart materializes the `<name>-s3-secret` Secret
carrying two credential pairs: an admin identity
(`admin_access_key_id` / `admin_secret_access_key`) and a read-only
identity (`read_access_key_id` / `read_secret_access_key`). The
Secret is generated once and is stable across upgrades and kept on
uninstall — credentials survive the release lifecycle, and the stack
outputs export the Secret name, never the material.

To own the identity list yourself, `s3.existing_config_secret` names
a Secret carrying the chart's contract — one key,
`seaweedfs_s3_config`, holding the inline JSON identities document
(every user, key pair, and permission). When set, the chart generates
nothing; the output points at your Secret instead.

`s3.enable_auth: false` serves an OPEN in-cluster S3 endpoint — dev
only, and the credentials output exports empty as the honest signal.

### Bucket lifecycle

Buckets declared in `s3.buckets` render into the chart's
`createBuckets` list, consumed by its install/upgrade hook. Each
entry is typed:

- `name` — S3 naming rules, validated in the spec (3–63 chars,
  lowercase, starts/ends alphanumeric).
- `anonymous_read` — unauthenticated reads on this bucket (public
  objects behind your ingress); auth still applies to writes.
- `ttl` — SeaweedFS TTL syntax, 1–255 followed by `m|h|d|w|M|y`
  (e.g. `"7d"`); empty = objects live forever.
- `object_lock` — S3 Object Lock, IRREVERSIBLE and forces versioning
  (WORM compliance workloads).
- `versioning` — S3 versioning (renders as the chart's `Enabled`).

The hook creates; it never destroys. Removing a bucket entry from the
spec does NOT delete the bucket or its data — bucket deletion is a
data operation performed against the store, never an IaC side effect.
This is deliberate: an accidental manifest edit must not be able to
destroy objects.

### Addressing

Empty `s3.domain_name` means path-style requests only
(`http://<endpoint>/<bucket>/<key>`) — the in-cluster norm, and what
every SDK supports with a "force path style" flag. Declaring
`domain_name` enables virtual-hosted-style addressing
(`{bucket}.{domain_name}`) for clients that require it — meaningful
once exposure with a wildcard host is composed on top.

## Replication Placement Codes

`replication` is the SeaweedFS XYZ placement code: X replicas in
other data centers, Y in other racks of the same data center, Z in
other servers of the same rack. `"001"` keeps one extra copy on
another server (needs at least 2 volume servers); `"010"` needs rack
topology; `"100"` needs multiple data centers. Empty means no
replication (`"000"`) — the single-node default.

Setting the code flips the chart's
`global.seaweedfs.enableReplication` and overrides the master's
`defaultReplication` and the filer's `defaultReplicaPlacement`
together — one field, consistent placement everywhere. The
non-negotiable constraint: the declared volume topology must be able
to place every copy, or writes fail. Rack and data-center identity
per volume server (`volume.rack`, `volume.dataCenter`, or the chart's
named `volumes` groups for zone-aware fleets) rides `helm_values`.

## Storage and Index Tuning

### PVCs, not hostPath — a deliberate posture

The chart's out-of-the-box storage is hostPath under `/ssd` and
`/storage` — bare-metal grain, where nodes have known disks and pods
pin to them. That default breaks on every managed cloud and kind
cluster. The modules therefore map every data volume to a
`persistentVolumeClaim` (master 5Gi, volume 30Gi, filer 10Gi, admin
10Gi when declared) and every logs volume to `emptyDir`. StorageClass
renders only when declared — absent means the cluster's default class
(the chart renders `storageClassName` empty, which Kubernetes treats
as nil). Bare-metal fleets can restore hostPath through the escape
hatch.

The volume tier's `dataDirs` is a LIST in the chart; the typed
surface models the canonical single-PVC entry (named `data`), sized
per pod. Exotic layouts — multiple data dirs, hostPath fleets — ride
`helm_values`, with the caveat that lists REPLACE on merge: an
override provides the whole list. The chart's resize hook (enabled by
default) patches PVCs automatically when a declared
`persistentVolumeClaim` size changes.

### Volume-server knobs

- `min_free_space_percent` (default 1, the upstream default) marks
  all volumes read-only when free disk drops below the threshold —
  the guard against writing a volume file past the disk.
- `filer.encrypt_volume_data` encrypts object data on the volume
  servers as it is written through the filer — ciphertext lands on
  disk; keys stay in filer metadata.

### Needle-index modes

`volume.index_mode` selects the memory/performance balance for the
needle index: `memory` (the upstream default — fastest, index
rebuilt from disk on start), `leveldb`, `leveldbMedium`, or
`leveldbLarge` (least memory, for very large stores). Empty defers to
the upstream default. The chart also supports splitting the index
onto its own persistent volume (`volume.idx`) for fleets where index
rebuild time matters — that rides `helm_values`; note the chart
REJECTS emptyDir for `idx` (an ephemeral index next to persistent
data forces a full rebuild every restart).

## Filer Metadata Stores

The filer's namespace lives in a metadata store, and the chart's
default is embedded leveldb on the filer's data volume
(`WEED_LEVELDB2_ENABLED: "true"` in the chart's filer env) — zero
extra infrastructure, and the reason the filer gets a PVC. The
embedded store is PER-POD: run 1 filer unless a shared external store
is wired.

External stores are configured exactly the way upstream configures
everything — `WEED_*` environment variables. The typed
`filer.extra_environment_vars` map is the declared surface for
plaintext options (e.g. `WEED_FILER_OPTIONS_RECURSIVE_DELETE:
"true"`, or the `WEED_MYSQL_*` / `WEED_POSTGRES_*` connection
settings); credential-bearing variables belong in a Kubernetes Secret
wired through `filer.secretExtraEnvironmentVars` in `helm_values`
(secretKeyRef entries — references, never material). With a shared
Postgres/MySQL store in place, filer replicas scale meaningfully.

## The Admin Console

`admin.enabled` installs the SeaweedFS management console — cluster
state, volumes, buckets, maintenance — on port 23646 (the exported
`admin_endpoint`). The security invariant: the console is NEVER
installed open. The chart requires `userKey`/`pwKey` alongside
`existingSecret`, and the modules always point it at a credentials
Secret — the module-materialized `<name>-admin-auth` (user `admin`,
random password; keys `user`/`password`) or the Secret referenced by
`existing_auth_secret`. There is no configuration that yields an
unauthenticated console.

Console state (configuration, maintenance bookkeeping) is in-memory
by default — lost on restart, fine for inspection-only use. Declaring
`admin.data_volume` persists it on a PVC (default 10Gi). The
maintenance WORKER — the pods that execute vacuum, balance and
erasure-coding jobs the console coordinates — is a separate chart arm
that rides `helm_values`.

## Monitoring

`service_monitor_enabled` maps to the chart's ONE
`global.seaweedfs.monitoring.enabled` flag, which gates
ServiceMonitors on every enabled tier (each tier's metrics port
defaults on in the chart). It requires the Prometheus Operator CRDs —
on a cluster without them, the install fails on unknown kinds, which
is why the flag defaults false.

## What the Typed Spec Covers — the 90/10 Rationale

The typed fields are the 90%: the decisions every deployment makes
(topology, storage, S3 posture, buckets, replication, the console,
monitoring, the image path) — each validated, defaulted, documented,
and rendered identically by both engines.

| Concern | Typed surface |
|---|---|
| Topology | `master.replicas`, `volume.replicas`, `filer.replicas`, `s3.dedicated.replicas` |
| Storage | per-tier `data_volume` (size + StorageClass), always PVC |
| S3 posture | `s3.enabled`, `s3.enable_auth`, `s3.existing_config_secret`, `s3.domain_name`, `s3.dedicated` |
| Buckets | `s3.buckets` (name, `anonymous_read`, `ttl`, `object_lock`, `versioning`) |
| Durability | `replication` (XYZ code), `filer.encrypt_volume_data` |
| Tuning | `master.volume_size_limit_mb`, `volume.max_volumes`, `volume.index_mode`, `volume.min_free_space_percent`, per-tier `resources` |
| Filer options | `filer.extra_environment_vars` (`WEED_*` keys) |
| Console | `admin.enabled`, `admin.existing_auth_secret`, `admin.data_volume` |
| Observability | `service_monitor_enabled` |
| Supply chain | `chart_version`, `image` (registry/repository/tag) |

The 10% rides `helm_values`, merged LAST with Helm `-f` semantics
(maps deep-merge with the later document winning, lists replace) —
identical on both engines:

- Per-tier scheduling: nodeSelector, tolerations, affinity, topology
  spread constraints, priority classes
- Probes, sidecars, init containers, extra volumes/mounts
- The SFTP arm (`sftp.*`)
- The maintenance worker (`worker.*`)
- mTLS: `global.seaweedfs.enableSecurity` plus the `certificates`
  section — requires cert-manager
- External filer stores: `filer.secretExtraEnvironmentVars`
  (credentials as secretKeyRef entries)
- Topology-aware volume groups (the chart's named `volumes` map) and
  hostPath/multi-dir storage layouts
- The COSI driver (`cosi.*`)
- The all-in-one dev mode (`allInOne.*`)

The escape hatch is a safety valve, never the primary interface — and
never the place for secrets: credential material belongs in
Kubernetes Secrets (`s3.existing_config_secret`,
`admin.existing_auth_secret`, `secretExtraEnvironmentVars`
references).

## Version Pins and Naming Contracts

| What | Value | Notes |
|---|---|---|
| Chart | `seaweedfs` at https://seaweedfs.github.io/seaweedfs/helm | Pinned 4.40.0 (appVersion 4.40); versioned with SeaweedFS releases |
| Release / fullname | `metadata.name` | `fullnameOverride` pins it — children derive deterministically |
| Children | `<name>-master`, `<name>-volume`, `<name>-filer`, `<name>-s3`, `<name>-admin` | The chart's componentName grammar |
| S3 Service | `<name>-s3`, port 8333 | Same Service for embedded and dedicated shapes |
| S3 credentials | `<name>-s3-secret` — `admin_access_key_id` / `admin_secret_access_key` / `read_access_key_id` / `read_secret_access_key` | Chart-generated once; stable across upgrades, kept on uninstall |
| S3 config contract | key `seaweedfs_s3_config` (inline JSON identities) | For `existing_config_secret` |
| Filer Service | `<name>-filer`, port 8888 | File namespace HTTP API |
| Master Service | `<name>-master`, port 9333 | Cluster coordination |
| Admin console | `<name>-admin`, port 23646 | Auth Secret `<name>-admin-auth`, keys `user`/`password` |
| Buckets folder | `/buckets` on the filer | Chart env default |
| Data volumes | PVC: master 5Gi, volume 30Gi, filer 10Gi, admin 10Gi | Component defaults; chart default is hostPath |
| Replication | XYZ code, empty = `"000"` | Flips `enableReplication`, overrides master + filer placement |

## IaC Twins

Pulumi (`module/values.go`) and Terraform (`locals.tf` +
`helm_release`) render identical chart values from the typed spec —
same `fullnameOverride`, same per-tier PVC/emptyDir mapping, same
dual-path S3 rendering, same bucket list shape — and both merge
`helm_values` last with the same semantics (Terraform passes
`values = [yamlencode(typed), helm_values]` and the provider merges
in order). Both materialize the same `<name>-admin-auth` Secret when
the console is enabled without an existing Secret. Keep the
typed-value rendering in lockstep.
