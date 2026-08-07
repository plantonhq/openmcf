# KubernetesRayCluster: Research and Design

## Introduction

KubernetesRayCluster declares a Ray cluster — the Kubernetes CR
(`rayclusters.ray.io/v1`) that the KubeRay operator reconciles into a
head pod, worker groups, and the Services that Ray clients, jobs, and
dashboards connect to. The prerequisite is a KubernetesKubeRayOperator
whose watch scope covers this namespace (cluster-wide with the
operator's defaults).

A Ray cluster is the shared distributed-compute service ML and data
teams connect to: notebooks and applications dial the CLIENT port
(10001), jobs submit through the DASHBOARD port (8265, the Job
Submission API), and Ray's own processes coordinate through the GCS
port (6379) on the head.

## The State Truth: GCS Lives in Head-Pod Memory

The head pod's GCS (Global Control Store) holds the cluster's control
state IN MEMORY. Without GCS fault tolerance, deleting or losing the
HEAD pod loses every job, actor, and worker registration — the
operator rebuilds an empty cluster. Worker pods are always
expendable; the head is not.

`gcs_fault_tolerance` renders the CR's `gcsFaultToleranceOptions`:
the control state moves to an external Redis-protocol store, and a
replaced head RECOVERS the cluster instead of rebuilding an empty
one. The composition story:

- **`redis_address`** accepts a literal `host:port` or a
  KubernetesValkey reference (its write Service endpoint — always the
  primary). A CEL rule makes the endpoint mandatory whenever the arm
  is enabled.
- **`redis_password_secret`** rides a secretKeyRef — on a
  KubernetesValkey, its module-materialized auth Secret, where the
  key is the ACL username (`default` unless you declared users).
  Secret reads cannot cross namespaces: co-locate the cluster with
  its Valkey or replicate the Secret.
- **`external_storage_namespace`** namespaces the keys so several Ray
  clusters share one store. Empty is safe — the operator derives one
  from the cluster's UID; set it explicitly only when the state must
  survive delete-and-recreate of the CR itself.

## Version/Image Lockstep

`ray_version` must match the Ray inside the image: the operator reads
`rayVersion` to shape its commands but runs the image as given, so a
mismatch fails at runtime, not at apply. The spec derives the default
image from the version (`rayproject/ray:<ray_version>`), which keeps
the lockstep true by construction — a custom image (dependencies
baked in, CUDA variants like `rayproject/ray:2.52.0-gpu`) overrides
it deliberately and owns keeping the versions identical. 2.52.0 is
the version KubeRay v1.6.x ships samples for.

## Authentication: Secure by Default, Rendered Always

Ray ≥ 2.52 / KubeRay ≥ 1.6 support token authentication: the operator
provisions a bearer-token Secret and the dashboard, job submission,
and client surfaces all require it. Two defaults collide here:

- The OPERATOR's nil-`authOptions` default is DISABLED — the legacy
  open cluster, where anyone who can reach port 8265 runs arbitrary
  code (the Job Submission API is an arbitrary-code-execution surface
  by design).
- The CATALOG's default is TOKEN.

The modules therefore render `authOptions` ALWAYS, never leaving it
to the CR default: an absent block would silently deploy the open
cluster. Empty auth or empty mode means token; only an explicit
`disabled` opts out — for fenced labs. In token mode the credential
handle is exported as `auth_token_secret` (key `auth_token`, the
operator's token-key contract); `existing_token_secret_name` is the
bring-your-own arm, telling the operator to skip generating a Secret
and read yours instead.

## The Head: Unloaded by Default

The modules OWN two head `rayStartParams` entries and merge them last
so user entries cannot override them:

- **`dashboard-host: "0.0.0.0"`** — the dashboard binds localhost
  otherwise and the head Service cannot answer (upstream's own
  complete sample sets exactly this).
- **`num-cpus: "0"`** unless `schedule_tasks_on_head` — the head
  advertises zero CPUs, so the Ray scheduler keeps application work
  off the control plane. A task-loaded head starves the GCS and
  destabilizes the whole cluster; upstream guidance starts production
  heads at 4 CPU / 8Gi and scales with cluster size. Single-node labs
  flip the field: with the production default and no worker groups
  there would be zero capacity to run anything.

Head resources are REQUIRED by validation — an unsized head is the
classic Ray outage. Pin the head to stable capacity: without GCS
fault tolerance, losing its node loses the cluster's state.

## Worker Groups, Autoscaling, and Accelerators

Worker groups are homogeneous sets (same resources, same Ray start
parameters); a CPU group and a GPU group is the classic two-group
shape, and group names must be unique (the operator derives pod names
per group). Each group's resources are required — Ray schedules
against declared capacity, and unsized workers make its scheduler
blind.

`autoscaling` injects the operator's in-tree autoscaler sidecar into
the head pod. It scales each group between `min_replicas` and
`max_replicas` on the Ray scheduler's demand — queued tasks and
actors, application-aware, unlike HPA. `replicas` becomes the initial
size only; a CEL rule orders `min_replicas <= replicas <=
max_replicas`. The operator's defaults stay authoritative when the
knobs are unset (60s idle timeout, Default upscaling).

Accelerators ride `extra_resource_limits`, which land in the
container LIMITS ONLY: Kubernetes rejects requests-without-limits for
extended resources, and Ray discovers accelerators from the container
limits. Pair the group's node selector and tolerations so the pods
land on accelerator nodes; with autoscaling, a GPU group can start at
zero and materialize only when a task asks for GPUs.

`suspend` is the idle-GPU cost lever: the operator deletes head and
worker PODS but keeps the declaration and (with GCS fault tolerance)
the external state; un-suspend to resume.

## Design Decisions

- **The CR applies without a cluster connection at plan time**
  (Terraform: the alekc/kubectl provider's `kubectl_manifest`, not
  `kubernetes_manifest`) — a RayCluster can be PLANNED before the
  operator's CRDs exist, which is what lets an infra chart deploy the
  operator and its clusters in one run.
- **No readiness wait, deliberately**: cluster readiness depends on
  the operator (multi-GB Ray image pulls, autoscaler sidecar
  injection, GCS startup) — never on applying the resource. The same
  never-block-on-a-controller posture as the sibling operator-CR
  modules.
- **Background deletion, explicitly**: the OPERATOR owns the
  RayCluster's cascade — its finalizer tears down the pods, the head
  Service, and the generated token Secret. Foreground propagation
  would block the delete on children the operator keeps reconciling.
- **Values render only when declared** so the operator's defaulting
  stays authoritative — with two exceptions: `authOptions` (see
  above) and the module-owned head `rayStartParams` entries. The v1
  CRD defaults `replicas`/`minReplicas`/`maxReplicas` itself
  (verified against the pinned schema); `rayStartParams` renders
  always (`{}` when empty) because the operator's own Go type
  serializes it unconditionally — matching that keeps server-side
  apply diffs quiet.
- **The 40-character name budget**: the operator derives
  `<name>-head-svc` (9-char suffix) and per-group worker pod names
  (`<name>-<group>-worker-<random>`) against the 63-character
  Kubernetes cap; both modules fail loudly past the budget, and short
  worker group names matter for the same reason.

## Naming Contracts and Endpoints

| What | Value | Notes |
|---|---|---|
| CR | `rayclusters.ray.io/v1` | Reconciled by the KubeRay operator (KubernetesKubeRayOperator) |
| Head Service | `<name>-head-svc` | The operator's naming contract; exported as `head_service` |
| Client endpoint | `<head_service>.<namespace>.svc.cluster.local:10001` | `ray.init("ray://…")` from notebooks and applications |
| Dashboard/Job API | `…:8265` | The Job Submission API and web dashboard; authenticated in token mode |
| GCS endpoint | `…:6379` | Ray's coordination port; what `ray start --address` joins |
| Token Secret key | `auth_token` | The operator's token-key contract; exported as `auth_token_secret` |
| Default image | `rayproject/ray:<ray_version>` | Docker Hub; version-locked by construction |

## IaC Twins

Pulumi (`module/raycluster_cr.go` + `module/locals.go`) and Terraform
(`raycluster_cr.tf` + `locals.tf`) render the same CR body (field
names are the CRD's own JSON keys, verified against the pinned ray.io
v1 schema) and the same outputs: the always-rendered auth block, the
module-owned head start params, the limits-only accelerator
placement, and the background-deletion posture. Keep them in
lockstep.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
