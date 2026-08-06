# KubernetesArgoWorkflows

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1alpha1`

**KubernetesArgoWorkflowsSpec** deploys Argo Workflows — the
Kubernetes-native workflow engine for DAG/step pipelines (CI jobs,
data and ML pipelines, batch orchestration) — from the official
`argo-workflows` Helm chart (https://argoproj.github.io/argo-helm).

WHAT GETS INSTALLED: the workflow controller (turns Workflow custom
resources into pods), the Argo server (UI + REST API, on by
default), and the workflow runner ServiceAccount + RBAC that
workflow pods execute under.

RUNNING WORKFLOWS: this kind installs the ENGINE only. Workflows,
WorkflowTemplates and CronWorkflows are Kubernetes custom resources
declared like any other manifest (KubernetesManifest, a chart, the
Argo CLI/UI) once the engine runs. Submit against the runner service
account (`workflow_service_account`, default "argo-workflow") —
pods run with ITS permissions, never the controller's.

CRD INSTALL REACHES THE INTERNET by default: the chart applies the
full-schema CRDs through a hook Job that downloads them from the
chart's GitHub release at install time (they are too large to
template inline). Air-gapped clusters set
`crds.full_schema = false` (chart-templated minified CRDs, no
download, weaker server-side validation) or mirror via
`crds.base_url`.

ARTIFACTS AND ARCHIVE are the two durability seams: declare
`artifact_repository` so steps can pass files (S3-compatible — a
KubernetesSeaweedFs pairs naturally — GCS, or Azure Blob), and
`archive` (Postgres/MySQL — a KubernetesPostgres pairs naturally) so
workflow history survives the Workflow CRs' garbage collection.
Without an archive, history lives only in etcd until
`retention_policy` (or TTL) prunes it.

EXPOSURE: the server Service stays ClusterIP; expose via first-class
kinds (KubernetesIngress, Gateway API kinds) over the exported
service handle.

The typed fields below cover the chart's meaningful configuration
surface; `helm_values` remains as the escape hatch for chart values
beyond them (merged last, Helm `-f` semantics, identical on both
engines) — workflow defaults documents, executor tuning, extra env —
a safety valve, never the primary interface. Never put secret
material in `helm_values`; every credential path in this spec rides
existing Secrets.

## Example

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesArgoWorkflows
metadata:
  name: test-workflows
spec:
  namespace:
    value: pipelines
  createNamespace: true
  chartVersion: 1.0.23
  controller:
    replicas: 2
    resources:
      requests:
        cpu: 250m
        memory: 512Mi
      limits:
        cpu: "1"
        memory: 1Gi
    workflowNamespaces:
      - team-a
      - team-b
    instanceId: platform
    parallelism: 50
    namespaceParallelism: 10
  server:
    enabled: true
    replicas: 2
    resources:
      requests:
        cpu: 100m
        memory: 128Mi
    authModes:
      - client
    baseHref: https://workflows.example.com/
  workflowServiceAccount: pipeline-runner
  artifactRepository:
    archiveLogs: true
    s3:
      bucket: workflow-artifacts
      endpoint:
        value: objects-s3.storage.svc.cluster.local:8333
      insecure: true
      credentialsSecret:
        secretName:
          value: objects-s3-auth
        accessKeyIdKey: accesskey
        secretAccessKeyKey: secretkey
  archive:
    engine: postgres
    host:
      value: workflows-db-rw.data.svc.cluster.local
    database: argo_archive
    credentialsSecret:
      name:
        value: workflows-db-app
      usernameKey: username
      passwordKey: password
    sslMode: require
  retentionPolicy:
    completed: 10
    failed: 3
    errored: 3
  crds:
    install: true
    keep: true
    fullSchema: true
    baseUrl: https://mirror.internal/argo-crds
  serviceMonitorEnabled: true
  image:
    registry: my.registry.com
    tag: v4.0.8
    pullSecretName: mirror-pull
  scheduling:
    nodeSelector:
      role: pipelines
    tolerations:
      - key: pipelines
        operator: Equal
        value: "true"
        effect: NoSchedule
    priorityClassName: batch-normal
  helmValues: |
    controller:
      workflowDefaults:
        spec:
          ttlStrategy:
            secondsAfterCompletion: 86400
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.namespace` | `string \| valueFrom` | yes |  | KubernetesNamespace (`spec.name`) |
| `spec.createNamespace` | `bool` |  |  |  |
| `spec.chartVersion` | `string` |  | `1.0.23` |  |
| `spec.controller` | `KubernetesArgoWorkflowsController` |  |  |  |
| `spec.controller.replicas` | `int32` |  | `1` |  |
| `spec.controller.resources` | `ContainerResources` |  |  |  |
| `spec.controller.resources.limits` | `CpuMemory` |  |  |  |
| `spec.controller.resources.limits.cpu` | `string` |  |  |  |
| `spec.controller.resources.limits.memory` | `string` |  |  |  |
| `spec.controller.resources.requests` | `CpuMemory` |  |  |  |
| `spec.controller.resources.requests.cpu` | `string` |  |  |  |
| `spec.controller.resources.requests.memory` | `string` |  |  |  |
| `spec.controller.workflowNamespaces` | `[]string` |  |  |  |
| `spec.controller.instanceId` | `string` |  |  |  |
| `spec.controller.parallelism` | `int32` |  |  |  |
| `spec.controller.namespaceParallelism` | `int32` |  |  |  |
| `spec.server` | `KubernetesArgoWorkflowsServer` |  |  |  |
| `spec.server.enabled` | `bool` |  | `true` |  |
| `spec.server.replicas` | `int32` |  | `1` |  |
| `spec.server.resources` | `ContainerResources` |  |  |  |
| `spec.server.resources.limits` | `CpuMemory` |  |  |  |
| `spec.server.resources.limits.cpu` | `string` |  |  |  |
| `spec.server.resources.limits.memory` | `string` |  |  |  |
| `spec.server.resources.requests` | `CpuMemory` |  |  |  |
| `spec.server.resources.requests.cpu` | `string` |  |  |  |
| `spec.server.resources.requests.memory` | `string` |  |  |  |
| `spec.server.authModes` | `[]string` |  |  |  |
| `spec.server.secure` | `bool` |  |  |  |
| `spec.server.baseHref` | `string` |  |  |  |
| `spec.workflowServiceAccount` | `string` |  | `argo-workflow` |  |
| `spec.artifactRepository` | `KubernetesArgoWorkflowsArtifactRepository` |  |  |  |
| `spec.artifactRepository.archiveLogs` | `bool` |  |  |  |
| `spec.artifactRepository.s3` | `KubernetesArgoWorkflowsArtifactS3` |  |  |  |
| `spec.artifactRepository.s3.bucket` | `string` | yes |  |  |
| `spec.artifactRepository.s3.endpoint` | `string \| valueFrom` |  |  | KubernetesSeaweedFs (`status.outputs.s3_endpoint`) |
| `spec.artifactRepository.s3.region` | `string` |  |  |  |
| `spec.artifactRepository.s3.insecure` | `bool` |  |  |  |
| `spec.artifactRepository.s3.useAmbientCredentials` | `bool` |  |  |  |
| `spec.artifactRepository.s3.credentialsSecret` | `KubernetesArgoWorkflowsS3CredentialsSecret` |  |  |  |
| `spec.artifactRepository.s3.credentialsSecret.secretName` | `string \| valueFrom` | yes |  | KubernetesSeaweedFs (`status.outputs.s3_credentials_secret_name`) |
| `spec.artifactRepository.s3.credentialsSecret.accessKeyIdKey` | `string` |  | `admin_access_key_id` |  |
| `spec.artifactRepository.s3.credentialsSecret.secretAccessKeyKey` | `string` |  | `admin_secret_access_key` |  |
| `spec.artifactRepository.gcs` | `KubernetesArgoWorkflowsArtifactGcs` |  |  |  |
| `spec.artifactRepository.gcs.bucket` | `string` | yes |  |  |
| `spec.artifactRepository.gcs.credentialsSecretName` | `string` |  |  |  |
| `spec.artifactRepository.azure` | `KubernetesArgoWorkflowsArtifactAzure` |  |  |  |
| `spec.artifactRepository.azure.endpoint` | `string` | yes |  |  |
| `spec.artifactRepository.azure.container` | `string` | yes |  |  |
| `spec.artifactRepository.azure.credentialsSecretName` | `string` |  |  |  |
| `spec.archive` | `KubernetesArgoWorkflowsArchive` |  |  |  |
| `spec.archive.engine` | `enum` | yes |  |  |
| `spec.archive.host` | `string \| valueFrom` | yes |  | KubernetesPostgres (`status.outputs.rw_service`) |
| `spec.archive.port` | `int32` |  |  |  |
| `spec.archive.database` | `string` | yes |  |  |
| `spec.archive.credentialsSecret` | `KubernetesArgoWorkflowsArchiveCredentials` | yes |  |  |
| `spec.archive.credentialsSecret.name` | `string \| valueFrom` | yes |  | KubernetesPostgres (`status.outputs.password_secret.name`) |
| `spec.archive.credentialsSecret.usernameKey` | `string` |  | `username` |  |
| `spec.archive.credentialsSecret.passwordKey` | `string` |  | `password` |  |
| `spec.archive.sslMode` | `string` |  |  |  |
| `spec.retentionPolicy` | `KubernetesArgoWorkflowsRetentionPolicy` |  |  |  |
| `spec.retentionPolicy.completed` | `int32` |  |  |  |
| `spec.retentionPolicy.failed` | `int32` |  |  |  |
| `spec.retentionPolicy.errored` | `int32` |  |  |  |
| `spec.crds` | `KubernetesArgoWorkflowsCrds` |  |  |  |
| `spec.crds.install` | `bool` |  | `true` |  |
| `spec.crds.keep` | `bool` |  | `true` |  |
| `spec.crds.fullSchema` | `bool` |  | `true` |  |
| `spec.crds.baseUrl` | `string` |  |  |  |
| `spec.serviceMonitorEnabled` | `bool` |  |  |  |
| `spec.image` | `KubernetesArgoWorkflowsImage` |  |  |  |
| `spec.image.registry` | `string` |  |  |  |
| `spec.image.tag` | `string` |  |  |  |
| `spec.image.pullSecretName` | `string` |  |  |  |
| `spec.scheduling` | `KubernetesArgoWorkflowsScheduling` |  |  |  |
| `spec.scheduling.nodeSelector` | `map<string, string>` |  |  |  |
| `spec.scheduling.tolerations` | `[]WorkloadToleration` |  |  |  |
| `spec.scheduling.tolerations[].key` | `string` |  |  |  |
| `spec.scheduling.tolerations[].operator` | `string` |  |  |  |
| `spec.scheduling.tolerations[].value` | `string` |  |  |  |
| `spec.scheduling.tolerations[].effect` | `string` |  |  |  |
| `spec.scheduling.tolerations[].tolerationSeconds` | `int64` |  |  |  |
| `spec.scheduling.priorityClassName` | `string` |  |  |  |
| `spec.helmValues` | `string` |  |  |  |

## Field Details

### spec.namespace

`string | valueFrom` · required

Namespace to install into. Accepts a literal namespace name or a
reference to a KubernetesNamespace resource.

- references: KubernetesNamespace (`spec.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesNamespace, name: <that resource's name>, fieldPath: spec.name}} -- a bare string does not parse

### spec.createNamespace

`bool`

When true, the namespace is created (with the standard Planton
governance labels) before installing and deleted with the resource.
When false, the namespace must already exist.

### spec.chartVersion

`string` · optional (explicit presence)

Helm chart version to install (e.g. "1.0.23" — chart 1.0.23 ships
Argo Workflows v4.0.8). Versions must exist as SERVED charts in
the repository index (https://argoproj.github.io/argo-helm).

- default: `1.0.23`

### spec.controller

`KubernetesArgoWorkflowsController`

The workflow controller.

### spec.controller.replicas

`int32` · optional (explicit presence)

Number of controller replicas. Empty = 1. Extra replicas are HOT
STANDBYS behind leader election — they take over on failure, they
do not share the workload.

- default: `1`
- rule: {"int32":{"gte":1}}

### spec.controller.resources

`ContainerResources`

CPU and memory for the controller container. Empty = no
requests/limits (the chart default).

### spec.controller.resources.limits

`CpuMemory`

The resource limits for the container.
Specify the maximum amount of CPU and memory that the container can use.

### spec.controller.resources.limits.cpu

`string`

### spec.controller.resources.limits.memory

`string`

### spec.controller.resources.requests

`CpuMemory`

The resource requests for the container.
Specify the minimum amount of CPU and memory that the container is guaranteed.

### spec.controller.resources.requests.cpu

`string`

### spec.controller.resources.requests.memory

`string`

### spec.controller.workflowNamespaces

`[]string`

Additional namespaces workflows will run in. The chart creates the
runner ServiceAccount + Role/RoleBinding in each listed namespace
AND in the install namespace (always included) — a workflow
submitted in a namespace without that identity fails at pod
creation. Empty = the install namespace only. KNOW THIS: the
controller WATCHES cluster-wide either way; this list places the
runner RBAC, it does not scope the watch.

### spec.controller.instanceId

`string`

Instance ID label this controller claims. Set it when running
SEVERAL Argo Workflows installs in one cluster: each controller
only reconciles Workflow CRs labeled with its instance ID, and
ignores the rest. Empty = unlabeled single-instance behavior.

### spec.controller.parallelism

`int32` · optional (explicit presence)

Maximum number of workflows running at once, cluster-wide
(additional ones queue as Pending). Empty = unlimited.

- rule: {"int32":{"gte":1}}

### spec.controller.namespaceParallelism

`int32` · optional (explicit presence)

Maximum number of workflows running at once PER NAMESPACE. Empty =
unlimited.

- rule: {"int32":{"gte":1}}

### spec.server

`KubernetesArgoWorkflowsServer`

The Argo server (UI + REST API).

### spec.server.enabled

`bool` · optional (explicit presence)

Run the server. Empty = true (the chart default). Disable for
controller-only installs where workflows are submitted purely as
CRs.

- default: `true`

### spec.server.replicas

`int32` · optional (explicit presence)

Number of server replicas. Empty = 1.

- default: `1`
- rule: {"int32":{"gte":1}}

### spec.server.resources

`ContainerResources`

CPU and memory for the server container. Empty = no requests/limits
(the chart default).

### spec.server.resources.limits

`CpuMemory`

The resource limits for the container.
Specify the maximum amount of CPU and memory that the container can use.

### spec.server.resources.limits.cpu

`string`

### spec.server.resources.limits.memory

`string`

### spec.server.resources.requests

`CpuMemory`

The resource requests for the container.
Specify the minimum amount of CPU and memory that the container is guaranteed.

### spec.server.resources.requests.cpu

`string`

### spec.server.resources.requests.memory

`string`

### spec.server.authModes

`[]string`

Authentication modes the server accepts. Empty = ["client"] (the
application default): callers present a Kubernetes bearer token
and act with ITS permissions. "server": the server performs every
request with its own ServiceAccount — no login, anyone who can
reach the endpoint has the server's full power; acceptable only
behind trusted network boundaries. "sso": OIDC login (configure
the chart's `server.sso` block via `helm_values`; its client
ID/secret ride existing Secrets by the chart's own contract —
declare "sso" here to open the mode).

- rule: {"repeated":{"unique":true,"items":{"string":{"in":["server","client","sso"]}}}}

### spec.server.secure

`bool`

Serve HTTPS with the server's self-signed certificate instead of
plain HTTP. The chart default is false (plain HTTP behind
composed exposure). KNOW THIS: flipping it changes the scheme
probes and clients use — the exported endpoint output follows it.

### spec.server.baseHref

`string`

The public base URL when composing exposure in front of the server
(e.g. "https://workflows.example.com/"). SSO redirects and UI
links embed it. Empty = relative links.

### spec.workflowServiceAccount

`string` · optional (explicit presence)

Name of the ServiceAccount workflow pods run as (created by the
chart, with Role/RoleBinding in every watched namespace). Empty =
"argo-workflow". Grant IRSA/Workload-Identity annotations or extra
RBAC to THIS account when workflows need cloud or cluster access —
never widen the controller's account.

- default: `argo-workflow`

### spec.artifactRepository

`KubernetesArgoWorkflowsArtifactRepository`

Where workflow steps read and write artifacts. Empty = no artifact
repository — steps can still run, but input/output artifacts and
archived logs have nowhere to live.

### spec.artifactRepository.archiveLogs

`bool`

Also archive each step's main-container logs into the repository
(the UI can then show logs after pods are gone).

### spec.artifactRepository.s3

`KubernetesArgoWorkflowsArtifactS3`

An S3-compatible object store — AWS S3 or an in-cluster
KubernetesSeaweedFs.

- rule: declare exactly one credential path — credentials_secret for declared keys, or use_ambient_credentials for keyless pod identity

### spec.artifactRepository.s3.bucket

`string` · required

Bucket name. The bucket must exist — Argo Workflows never creates
it.

- rule: {"required":true}

### spec.artifactRepository.s3.endpoint

`string | valueFrom`

S3 endpoint as host[:port] — no scheme (e.g. "s3.amazonaws.com",
or an in-cluster SeaweedFS S3 endpoint). Accepts a literal
endpoint or a reference to a KubernetesSeaweedFs resource (its S3
endpoint). Empty = "s3.amazonaws.com".

- references: KubernetesSeaweedFs (`status.outputs.s3_endpoint`)
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesSeaweedFs, name: <that resource's name>, fieldPath: status.outputs.s3_endpoint}} -- a bare string does not parse

### spec.artifactRepository.s3.region

`string`

AWS region (real S3 only; in-cluster stores ignore it).

### spec.artifactRepository.s3.insecure

`bool`

Speak plain HTTP to the endpoint (in-cluster stores without TLS).
Never set it towards real S3.

### spec.artifactRepository.s3.useAmbientCredentials

`bool`

Keyless access: sign requests with the pod's ambient cloud
identity (IRSA / workload identity on the RUNNER service account)
instead of declared keys. Mutually exclusive with
`credentials_secret` (enforced below).

### spec.artifactRepository.s3.credentialsSecret

`KubernetesArgoWorkflowsS3CredentialsSecret`

The existing Secret holding the access keys. Defaults compose a
KubernetesSeaweedFs resource's generated credentials Secret
(`<name>-s3-secret`) as-is — its name by reference, its admin key
pair through the key-name defaults below. Required unless
`use_ambient_credentials`.

### spec.artifactRepository.s3.credentialsSecret.secretName

`string | valueFrom` · required

Name of the Secret. Accepts a literal name or a reference to a
KubernetesSeaweedFs resource (its S3 credentials Secret). The
Secret must live in the install namespace (the controller reads
it there).

- references: KubernetesSeaweedFs (`status.outputs.s3_credentials_secret_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesSeaweedFs, name: <that resource's name>, fieldPath: status.outputs.s3_credentials_secret_name}} -- a bare string does not parse

### spec.artifactRepository.s3.credentialsSecret.accessKeyIdKey

`string` · optional (explicit presence)

Key NAME within the Secret (a reference, not secret material)
holding the access key ID. Empty = "admin_access_key_id"
(the KubernetesSeaweedFs convention — the generated `-s3-secret`
composes with zero key configuration). For a Secret shaped to the
argo-workflows chart's documented example (keys
`accesskey`/`secretkey`), set both key names explicitly.

- default: `admin_access_key_id`

### spec.artifactRepository.s3.credentialsSecret.secretAccessKeyKey

`string` · optional (explicit presence)

Key holding the secret access key. Empty =
"admin_secret_access_key".

- default: `admin_secret_access_key`

### spec.artifactRepository.gcs

`KubernetesArgoWorkflowsArtifactGcs`

Google Cloud Storage.

### spec.artifactRepository.gcs.bucket

`string` · required

Bucket name.

- rule: {"required":true}

### spec.artifactRepository.gcs.credentialsSecretName

`string`

Existing Secret holding a service-account JSON key (key
`serviceAccountKey`). Empty = keyless: GKE Workload Identity on
the runner service account.

### spec.artifactRepository.azure

`KubernetesArgoWorkflowsArtifactAzure`

Azure Blob Storage.

### spec.artifactRepository.azure.endpoint

`string` · required

Storage-account blob endpoint (e.g.
"https://mystorageaccount.blob.core.windows.net").

- rule: {"required":true}

### spec.artifactRepository.azure.container

`string` · required

Container name.

- rule: {"required":true}

### spec.artifactRepository.azure.credentialsSecretName

`string`

Existing Secret holding the storage-account access key (key
`account-access-key`). Empty = keyless: Azure Workload Identity /
managed identity on the runner service account.

### spec.archive

`KubernetesArgoWorkflowsArchive`

The workflow archive — completed workflows written to a relational
database so history survives CR garbage collection.

### spec.archive.engine

`enum` · required

Database engine.

- rule: {"required":true,"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `kubernetes_argo_workflows_archive_engine_unspecified` -- Unspecified. Declare the engine explicitly.
- `postgres` -- PostgreSQL — pairs naturally with a KubernetesPostgres resource.
- `mysql` -- MySQL — pairs naturally with a KubernetesMysql resource.

### spec.archive.host

`string | valueFrom` · required

Database host — name or FQDN, no port (e.g.
"workflows-db-rw.data.svc.cluster.local"). Accepts a literal host
or a reference to a KubernetesPostgres resource (its read-write
Service; same-namespace referencing may use the short name,
cross-namespace needs the FQDN).

- references: KubernetesPostgres (`status.outputs.rw_service`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesPostgres, name: <that resource's name>, fieldPath: status.outputs.rw_service}} -- a bare string does not parse

### spec.archive.port

`int32` · optional (explicit presence)

Database port. Empty = the engine's default (5432 postgres, 3306
mysql).

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.archive.database

`string` · required

Database name. The database must exist — the controller creates
its tables, never the database.

- rule: {"required":true}

### spec.archive.credentialsSecret

`KubernetesArgoWorkflowsArchiveCredentials` · required

Existing Secret holding the database credentials.

- rule: {"required":true}

### spec.archive.credentialsSecret.name

`string | valueFrom` · required

Secret name. Accepts a literal name or a reference to a
KubernetesPostgres resource (its operator-maintained application
Secret — whose `username`/`password` keys are exactly the key-name
defaults below, so it composes untouched). The Secret must live in
the install namespace (the controller reads it there).

- references: KubernetesPostgres (`status.outputs.password_secret.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesPostgres, name: <that resource's name>, fieldPath: status.outputs.password_secret.name}} -- a bare string does not parse

### spec.archive.credentialsSecret.usernameKey

`string` · optional (explicit presence)

Key holding the username. Empty = "username".

- default: `username`

### spec.archive.credentialsSecret.passwordKey

`string` · optional (explicit presence)

Key holding the password. Empty = "password".

- default: `password`

### spec.archive.sslMode

`string`

Require TLS towards the database ("require", "verify-ca",
"verify-full"). Empty = plain ("disable") — the in-cluster
default.

- rule: {"string":{"in":["","disable","require","verify-ca","verify-full"]}}

### spec.retentionPolicy

`KubernetesArgoWorkflowsRetentionPolicy`

How many COMPLETED Workflow CRs of each outcome the controller
keeps in the cluster (older ones are deleted; the `archive` keeps
their history). Empty = keep everything until TTLs or humans clean
up.

### spec.retentionPolicy.completed

`int32` · optional (explicit presence)

Completed-successfully workflows to keep.

- rule: {"int32":{"gte":0}}

### spec.retentionPolicy.failed

`int32` · optional (explicit presence)

Failed workflows to keep.

- rule: {"int32":{"gte":0}}

### spec.retentionPolicy.errored

`int32` · optional (explicit presence)

Errored workflows to keep.

- rule: {"int32":{"gte":0}}

### spec.crds

`KubernetesArgoWorkflowsCrds`

CRD lifecycle (see the CRD note on the spec).

### spec.crds.install

`bool` · optional (explicit presence)

Install (and upgrade) the Argo Workflows CRDs with the release.
Disable only when another release in the cluster already owns
them.

- default: `true`

### spec.crds.keep

`bool` · optional (explicit presence)

Keep the CRDs when this resource is destroyed. Default true —
REMOVING the CRDs DELETES EVERY Workflow, WorkflowTemplate and
CronWorkflow in the cluster with them. Turn off only on throwaway
clusters.

- default: `true`

### spec.crds.fullSchema

`bool` · optional (explicit presence)

Install the FULL-SCHEMA CRDs via the chart's hook Job, which
DOWNLOADS them from the chart's GitHub release at install time
(they exceed inline-template limits). Default true. Set false on
air-gapped clusters: the chart falls back to templated MINIFIED
CRDs (no download, no full server-side validation), or mirror the
files and point `base_url` at the mirror.

- default: `true`

### spec.crds.baseUrl

`string`

Base URL the full-schema hook Job downloads CRD YAMLs from. Empty
= the chart's GitHub release for the pinned version. Set it to an
internal mirror on restricted networks.

### spec.serviceMonitorEnabled

`bool`

Create ServiceMonitors for the controller's /metrics (requires the
Prometheus Operator CRDs — deploy KubernetesKubePrometheusStack
first). Chart default: false.

### spec.image

`KubernetesArgoWorkflowsImage`

Override the Argo Workflows images (air-gap path). One override
serves controller, server and executor — the chart composes
{registry}/{repository}:{tag} per component.

### spec.image.registry

`string`

Registry to pull all Argo Workflows images from (e.g.
"my.registry.com"). Empty = "quay.io". The per-component
repository paths (argoproj/workflow-controller, argoproj/argocli,
argoproj/argoexec) stay upstream; override those via `helm_values`
when a mirror re-paths them.

### spec.image.tag

`string`

Tag for all Argo Workflows images. Empty = the chart's appVersion
for the pinned chart_version.

### spec.image.pullSecretName

`string`

Name of an existing image-pull Secret in the namespace, for
private mirrors (applied to all components).

### spec.scheduling

`KubernetesArgoWorkflowsScheduling`

Scheduling for the controller and server pods.

### spec.scheduling.nodeSelector

`map<string, string>`

Node selector for the controller and server pods (workflow pods
schedule per their own Workflow specs).

### spec.scheduling.tolerations

`[]WorkloadToleration`

Tolerations for the controller and server pods.

### spec.scheduling.tolerations[].key

`string`

Taint key to tolerate. Empty key with operator "Exists" tolerates every taint.

### spec.scheduling.tolerations[].operator

`string`

How key/value match: "Equal" (default — value must match too) or "Exists"
(key presence alone matches).

- rule: Toleration operator must be either "Equal" or "Exists"

### spec.scheduling.tolerations[].value

`string`

Taint value to match when operator is "Equal".

### spec.scheduling.tolerations[].effect

`string`

Which taint effect is tolerated: "NoSchedule", "PreferNoSchedule", or
"NoExecute". Empty tolerates all effects for the key.

- rule: Toleration effect must be one of "NoSchedule", "PreferNoSchedule", or "NoExecute"

### spec.scheduling.tolerations[].tolerationSeconds

`int64` · optional (explicit presence)

For "NoExecute" taints only: how many seconds already-running pods stay bound
after the taint appears. Unset means tolerate forever.

### spec.scheduling.priorityClassName

`string`

Priority class name for the controller and server pods.

### spec.helmValues

`string`

Escape hatch: additional chart values as a YAML document, merged
LAST over everything the typed fields render (Helm `-f` semantics,
identical on both engines). For the chart surface beyond the typed
fields (workflowDefaults, executor resources, sso, extra env,
priority-class per component, ...) — never the substitute for
them. Do not put secrets here; credential material belongs in
referenced Secrets.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesArgoWorkflows, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.namespace` | `string` | Namespace Argo Workflows is installed in. |
| `status.outputs.release_name` | `string` | Helm release name (equals metadata.name). |
| `status.outputs.server_service` | `string` | Name of the Argo server Service (`<name>-server`) — the backend handle exposure kinds (KubernetesIngress, KubernetesHttpRoute) reference. Empty when the server is disabled. |
| `status.outputs.server_kube_endpoint` | `string` | In-cluster endpoint of the Argo server (e.g. "http://main-server.pipelines.svc.cluster.local:2746"). Plain HTTP by default; HTTPS when `server.secure` is set. Empty when the server is disabled. |
| `status.outputs.workflow_service_account` | `string` | Name of the ServiceAccount workflow pods run as — annotate THIS account for IRSA/Workload-Identity when workflows need cloud access. |
| `status.outputs.port_forward_command` | `string` | Command to port-forward the UI to a workstation (http://localhost:2746 unless `server.secure`). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespace` | KubernetesNamespace | `spec.name` |
| `spec.artifactRepository.s3.endpoint` | KubernetesSeaweedFs | `status.outputs.s3_endpoint` |
| `spec.artifactRepository.s3.credentialsSecret.secretName` | KubernetesSeaweedFs | `status.outputs.s3_credentials_secret_name` |
| `spec.archive.host` | KubernetesPostgres | `status.outputs.rw_service` |
| `spec.archive.credentialsSecret.name` | KubernetesPostgres | `status.outputs.password_secret.name` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
