# KubernetesFlinkDeployment: Research and Design

## Introduction

KubernetesFlinkDeployment declares a Flink cluster — the Kubernetes
CR (`flinkdeployments.flink.apache.org/v1beta1`) that the Apache
Flink Kubernetes Operator reconciles into a JobManager, its
TaskManagers, and (in application mode) the job they run. The
prerequisite is a KubernetesFlinkOperator whose watch scope covers
this namespace.

## Two Shapes, One Field Apart

- **Application cluster** (`job` set): the cluster exists to run that
  one job and follows its lifecycle — the recommended production
  shape, isolation per pipeline. The production pattern bakes the job
  jar into a custom image (`local:///…` jar URIs point inside the
  image; remote `https://`/`s3://` schemes download at submission).
- **Session cluster** (`job` absent): an empty Flink runtime that
  accepts jobs submitted separately — FlinkSessionJob CRs or the REST
  API — for many short-lived jobs sharing warm capacity. In native
  mode a session cluster holds zero TaskManagers until a job asks.

## The State Truths: The Operator's Own Validator, Made Visible

The operator's validator enforces three rules; this spec mirrors them
as authoring-time validation so they fail with a message before the
apply instead of at admission (or, webhook-less, at reconcile):

1. **Non-stateless upgrades need `state.checkpoints_dir`.** Any
   `job.upgrade_mode` other than `stateless` requires it — the
   checkpoints are where the state to carry across upgrades lives.
2. **Savepoint mode additionally needs `state.savepoints_dir`.** The
   operator triggers a savepoint there on every upgrade (and
   `savepoint_trigger_nonce` writes manual savepoints there too).
3. **JobManager standbys need HA.** `job_manager.replicas` beyond 1
   are warm standbys coordinated through HA metadata —
   `state.high_availability` (with its `storage_dir`) is required,
   and HA is also what the operator's rollback feature rides. It
   pairs naturally with `last-state` upgrades.

The upgrade modes, honestly:

- **`stateless`** (default): restart from clean state — correct ONLY
  for pipelines with no state worth carrying (enrichment, routing,
  filtering). On a stateful job it silently discards state on every
  upgrade.
- **`last-state`**: restore the newest checkpoint — fast.
- **`savepoint`**: drain through a savepoint — the operationally
  safest.

Every directory must live on storage every pod can reach — object
storage, not a local path. The `state` fields are Flink config sugar
(`state.checkpoints.dir`, `state.savepoints.dir`,
`high-availability.storageDir`) with the operator's requirements made
visible; the module renders `high-availability.type: kubernetes`, the
current key form at the pinned operator.

## The S3 Seam: Secret-Native by Design

`flink_configuration` renders into a ConfigMap in clear text —
credentials must never enter it. `state.s3` is the seam that keeps
them out:

- The ENDPOINT and `path_style_access` render into Flink config
  (`s3.endpoint`, `s3.path.style.access` — path-style true is the
  default, correct for in-cluster object stores; explicit false is
  the AWS-S3-itself posture).
- The CREDENTIALS ride pod environment from Secret references:
  `AWS_ACCESS_KEY_ID` / `AWS_SECRET_ACCESS_KEY` as secretKeyRefs in
  the CR's pod template, read at pod runtime. Secret reads cannot
  cross namespaces — co-locate the deployment with its credentials or
  replicate the Secret.
- The foreign-key defaults compose a KubernetesSeaweedFs: its S3
  endpoint and its generated credentials Secret
  (`admin_access_key_id` / `admin_secret_access_key` for the admin
  identity). Deploy the store first and declare the bucket on it —
  Flink does not create buckets.
- Omit `state.s3` entirely where pods reach object storage ambiently
  (IRSA/workload identity on a cloud bucket).

### The builtin_plugin_jar truth

The official Flink images ship the S3 filesystem plugin DISABLED,
bundled under `/opt/flink/opt`. `builtin_plugin_jar` names the exact
jar to activate (rendered as the `ENABLE_BUILT_IN_PLUGINS` pod env
var) — e.g. `flink-s3-fs-hadoop-2.1.0.jar`, where the version in the
name must match the image's Flink patch version. Without it, every
`s3://` path fails at runtime with "unsupported filesystem scheme".
Custom images that bake the plugin into `/opt/flink/plugins` leave it
empty.

## The Sizing Arithmetic

Total task slots = TaskManager count ×
`taskmanager.numberOfTaskSlots` (a `flink_configuration` key, default
1), and a job needs `parallelism` slots to run. An under-slotted
cluster holds the job in a scheduling wait, not an error — the
failure mode to know before it looks like a hang.

Mode decides who counts TaskManagers:

- **`native`** (default, operator-recommended): Flink asks Kubernetes
  for TaskManagers on demand from the job's parallelism —
  `task_manager.replicas` is ignored there.
- **`standalone`**: fixed resources, no Flink↔Kubernetes API
  traffic — `task_manager.replicas` is the worker count.

Both tiers default to 1 CPU / 2Gi (the upstream example sizing), and
Flink derives its JVM, managed, and network memory from the container
memory — resize instead of tuning JVM flags.

## Version/Image Lockstep

`flink_version` is the CR's own enum (`v2_1` = Flink 2.1; `v1_20` is
the last 1.x LTS line). The default image derives from it
(`flink:<major>.<minor>`), keeping the lockstep true by construction;
a custom image must carry exactly that Flink version — the operator
shapes its submission protocol from the declared version, and a
mismatch fails at runtime, not at apply. The operator refuses
changing `flink_version` after a last-state suspend: restore first,
then upgrade.

## Config Ownership

`flink_configuration` is Flink's own key space. Three ownership
rules:

- **The operator forbids keys it owns** (`kubernetes.cluster-id`,
  `kubernetes.namespace`, the HA cluster id) — they derive from this
  resource; setting them is rejected at admission.
- **The module-rendered state keys merge LAST** over user entries —
  the typed `state` fields are the truth, and a colliding user entry
  loses deliberately.
- **Credentials never enter it** (see the S3 seam above).

## Design Decisions

- **The CR applies without a cluster connection at plan time**
  (Terraform: the alekc/kubectl provider's `kubectl_manifest`) — a
  FlinkDeployment can be PLANNED before the operator's CRDs exist,
  which is what lets an infra chart deploy the operator and its
  clusters in one run.
- **No readiness wait, deliberately**: readiness depends on the
  operator (image pulls, job submission, TaskManager registration) —
  never on applying the resource; the same never-block-on-a-controller
  posture as the sibling operator-CR modules.
- **Background deletion, explicitly**: the OPERATOR owns the
  FlinkDeployment's cascade — its finalizer tears down the
  JobManager, TaskManagers, and Services; foreground propagation
  would block the delete on children the operator keeps reconciling.
- **Both resource blocks render ALWAYS**: the operator's validator
  requires resource memory on both tiers, so JobManager and
  TaskManager blocks render even when unset, at the spec's defaults.
  The CR's Resource takes CPU as a NUMBER (cores) and memory as a
  string — the modules convert Kubernetes quantities (`1000m` → 1.0),
  preferring requests for CPU and limits for memory (the ceiling
  Flink sizes its JVM from).
- **Everything else renders only when declared**, so the operator's
  defaulting stays authoritative — `mode` only when standalone, `job`
  sub-fields only on divergence from CR defaults, the pod template
  only when it would carry something (scheduling, image-pull Secrets,
  or the S3 credential env; `flink-main-container` is the operator's
  merge contract for the main container).
- **Declarative operations as fields**: `job.state: suspended` is the
  declarative pause (a stateful upgrade_mode retains state for the
  resume), `restart_nonce` forces a restart without a spec change,
  `savepoint_trigger_nonce` triggers a savepoint now,
  `initial_savepoint_path` bootstraps the FIRST start when migrating
  a job in, and `allow_non_restored_state` permits restoring past a
  removed operator instead of failing.

## Naming Contracts and Endpoints

| What | Value | Notes |
|---|---|---|
| CR | `flinkdeployments.flink.apache.org/v1beta1` | Reconciled by the Flink Kubernetes Operator (KubernetesFlinkOperator) |
| REST Service | `<name>-rest` | The operator's naming contract; the Flink REST API and web UI on 8081 — exported as `rest_service` |
| REST endpoint | `<rest_service>.<namespace>.svc.cluster.local:8081` | Where session-mode jobs submit and job status reads from |
| TaskManager pods | `<name>-taskmanager-N-M` | Behind the 45-char name budget (63-char Kubernetes cap) |
| Default image | `flink:<major>.<minor>` (from `flink_version`) | Docker Hub; version-locked by construction |
| Service account | `flink` (default) | The KubernetesFlinkOperator chart's job service account |

## IaC Twins

Pulumi (`module/flinkdeployment_cr.go` + `module/locals.go`) and
Terraform (`flinkdeployment_cr.tf` + `locals.tf`) render the same CR
body (field names are the CR's own JSON keys, verified against the
pinned operator's API classes) and the same outputs: the same
state-key merge order, the same Secret-native S3 env, the same
always-rendered resource blocks with the same quantity-to-cores
conversion, and the same background-deletion posture. Keep them in
lockstep.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
