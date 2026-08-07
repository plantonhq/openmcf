# KubernetesMlflow

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1alpha1`

**KubernetesMlflowSpec** deploys the MLflow tracking server and
model registry — where ML experiments, runs, metrics, models and
their artifacts live (https://mlflow.org).

WHAT GETS INSTALLED: the module renders its own typed manifests
(Deployment, Service, Secrets, and optionally a PVC and a
garbage-collection CronJob) around the OFFICIAL MLflow image
(`ghcr.io/mlflow/mlflow`). MLflow publishes no Helm chart — this
module IS the distribution, shaped after upstream's own deployment
reference.

TWO KINDS OF STATE, TWO BACKENDS:

  - The BACKEND STORE holds experiments, runs, metrics, params and
    the model registry: PostgreSQL (a KubernetesPostgres composes
    naturally — the production path) or sqlite on a PVC (the
    zero-dependency default; single replica only).
  - The ARTIFACT STORE holds the big binary outputs (models,
    datasets, plots): an S3-compatible store (a KubernetesSeaweedFs
    composes naturally), AWS S3, GCS, Azure Blob, or a local PVC
    (the default). The server PROXIES artifact traffic — clients
    talk only to MLflow and never need store credentials.

SECURED BY DEFAULT: upstream's server is OPEN unless auth is
configured, and upstream's own auth example ships
admin/password1234. This kind enables basic authentication by
default with a module-generated admin password (`<name>-admin-auth`)
— the upstream defaults never ship. Auth state (users,
permissions) follows the backend store: the same PostgreSQL
database, or a sqlite file beside the tracking data on the PVC.

CREDENTIALS ARE SECRET-NATIVE: the backend-store connection URI
(which embeds the database password) is composed AT APPLY TIME
into a module-owned Secret from the referenced credential Secret
and reaches the server as an environment variable — nothing
credential-bearing appears in any rendered manifest. Artifact-store
credentials ride environment variables from their Secrets the same
way.

EXPOSURE: the Service stays ClusterIP; expose it via first-class
kinds over the exported service handle.

## Example

```yaml
# Full-surface shape: the external-PostgreSQL backend store, the
# S3-compatible artifact store (SeaweedFS-shaped), secured basic auth
# with a private-by-default permission posture, garbage collection,
# metrics with a ServiceMonitor, scheduling and the env/args escape
# hatches — the offline plan/preview proof for the widest typed
# rendering.
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesMlflow
metadata:
  name: mlflow-full
spec:
  namespace:
    value: mlflow-full
  createNamespace: true
  server:
    replicas: 2
    workers: 8
    resources:
      requests:
        cpu: 250m
        memory: 512Mi
      limits:
        memory: 2Gi
  backendStore:
    postgres:
      host:
        value: mlflow-pg-rw
      port: 5432
      databaseName: mlflow
      username: mlflow
      passwordSecret:
        secretName:
          value: mlflow-pg-app
        secretKey: password
  artifactStore:
    s3Compatible:
      endpoint:
        value: http://artifacts-s3.mlflow-full.svc.cluster.local:8333
      bucket: mlflow-artifacts
      prefix: runs
      credentialsSecret:
        secretName:
          value: artifacts-s3-secret
        accessKeyIdKey: admin_access_key_id
        secretAccessKeyKey: admin_secret_access_key
  auth:
    enabled: true
    adminUsername: mlops
    defaultPermission: NO_PERMISSIONS
  gc:
    enabled: true
    schedule: "30 2 * * 0"
    olderThan: 72h
  service:
    type: ClusterIP
  metrics:
    enabled: true
    serviceMonitorEnabled: true
  scheduling:
    nodeSelector:
      workload: ml
  extraEnv:
    MLFLOW_HTTP_REQUEST_TIMEOUT: "120"
  extraArgs:
    - --gunicorn-opts=--timeout 120
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.namespace` | `string \| valueFrom` | yes |  | KubernetesNamespace (`spec.name`) |
| `spec.createNamespace` | `bool` |  |  |  |
| `spec.server` | `KubernetesMlflowServer` |  |  |  |
| `spec.server.replicas` | `int32` |  | `1` |  |
| `spec.server.image` | `KubernetesMlflowImage` |  |  |  |
| `spec.server.image.repository` | `string` |  | `ghcr.io/mlflow/mlflow` |  |
| `spec.server.image.tag` | `string` |  | `v3.15.0-full` |  |
| `spec.server.workers` | `int32` |  |  |  |
| `spec.server.resources` | `ContainerResources` |  |  |  |
| `spec.server.resources.limits` | `CpuMemory` |  |  |  |
| `spec.server.resources.limits.cpu` | `string` |  |  |  |
| `spec.server.resources.limits.memory` | `string` |  |  |  |
| `spec.server.resources.requests` | `CpuMemory` |  |  |  |
| `spec.server.resources.requests.cpu` | `string` |  |  |  |
| `spec.server.resources.requests.memory` | `string` |  |  |  |
| `spec.backendStore` | `KubernetesMlflowBackendStore` |  |  |  |
| `spec.backendStore.sqlitePvc` | `KubernetesMlflowSqlitePvc` |  |  |  |
| `spec.backendStore.sqlitePvc.storageSize` | `string` |  | `5Gi` |  |
| `spec.backendStore.sqlitePvc.storageClass` | `string` |  |  |  |
| `spec.backendStore.postgres` | `KubernetesMlflowPostgres` |  |  |  |
| `spec.backendStore.postgres.host` | `string \| valueFrom` | yes |  | KubernetesPostgres (`status.outputs.rw_service`) |
| `spec.backendStore.postgres.port` | `int32` |  | `5432` |  |
| `spec.backendStore.postgres.databaseName` | `string` |  | `mlflow` |  |
| `spec.backendStore.postgres.username` | `string` |  | `mlflow` |  |
| `spec.backendStore.postgres.passwordSecret` | `KubernetesMlflowPasswordSecret` | yes |  |  |
| `spec.backendStore.postgres.passwordSecret.secretName` | `string \| valueFrom` | yes |  | KubernetesPostgres (`status.outputs.password_secret.name`) |
| `spec.backendStore.postgres.passwordSecret.secretKey` | `string` |  | `password` |  |
| `spec.backendStore.mysql` | `KubernetesMlflowMysql` |  |  |  |
| `spec.backendStore.mysql.host` | `string \| valueFrom` | yes |  | KubernetesMysql (`status.outputs.primary_service`) |
| `spec.backendStore.mysql.port` | `int32` |  | `3306` |  |
| `spec.backendStore.mysql.databaseName` | `string` |  | `mlflow` |  |
| `spec.backendStore.mysql.username` | `string` |  | `mlflow` |  |
| `spec.backendStore.mysql.passwordSecret` | `KubernetesMlflowMysqlPasswordSecret` | yes |  |  |
| `spec.backendStore.mysql.passwordSecret.secretName` | `string \| valueFrom` | yes |  | KubernetesMysql (`status.outputs.root_password_secret.name`) |
| `spec.backendStore.mysql.passwordSecret.secretKey` | `string` |  | `password` |  |
| `spec.artifactStore` | `KubernetesMlflowArtifactStore` |  |  |  |
| `spec.artifactStore.pvc` | `KubernetesMlflowPvcArtifacts` |  |  |  |
| `spec.artifactStore.pvc.storageSize` | `string` |  | `10Gi` |  |
| `spec.artifactStore.pvc.storageClass` | `string` |  |  |  |
| `spec.artifactStore.s3Compatible` | `KubernetesMlflowS3CompatibleArtifacts` |  |  |  |
| `spec.artifactStore.s3Compatible.endpoint` | `string \| valueFrom` | yes |  | KubernetesSeaweedFs (`status.outputs.s3_endpoint`) |
| `spec.artifactStore.s3Compatible.bucket` | `string` | yes |  |  |
| `spec.artifactStore.s3Compatible.prefix` | `string` |  |  |  |
| `spec.artifactStore.s3Compatible.credentialsSecret` | `KubernetesMlflowS3CredentialsSecret` | yes |  |  |
| `spec.artifactStore.s3Compatible.credentialsSecret.secretName` | `string \| valueFrom` | yes |  | KubernetesSeaweedFs (`status.outputs.s3_credentials_secret_name`) |
| `spec.artifactStore.s3Compatible.credentialsSecret.accessKeyIdKey` | `string` |  | `admin_access_key_id` |  |
| `spec.artifactStore.s3Compatible.credentialsSecret.secretAccessKeyKey` | `string` |  | `admin_secret_access_key` |  |
| `spec.artifactStore.awsS3` | `KubernetesMlflowAwsS3Artifacts` |  |  |  |
| `spec.artifactStore.awsS3.bucket` | `string` | yes |  |  |
| `spec.artifactStore.awsS3.prefix` | `string` |  |  |  |
| `spec.artifactStore.awsS3.region` | `string` | yes |  |  |
| `spec.artifactStore.awsS3.credentialsSecret` | `KubernetesMlflowAwsCredentialsSecret` |  |  |  |
| `spec.artifactStore.awsS3.credentialsSecret.secretName` | `string` | yes |  |  |
| `spec.artifactStore.awsS3.credentialsSecret.accessKeyIdKey` | `string` |  | `access_key_id` |  |
| `spec.artifactStore.awsS3.credentialsSecret.secretAccessKeyKey` | `string` |  | `secret_access_key` |  |
| `spec.artifactStore.gcs` | `KubernetesMlflowGcsArtifacts` |  |  |  |
| `spec.artifactStore.gcs.bucket` | `string` | yes |  |  |
| `spec.artifactStore.gcs.prefix` | `string` |  |  |  |
| `spec.artifactStore.gcs.credentialsSecret` | `KubernetesMlflowGcpCredentialsSecret` |  |  |  |
| `spec.artifactStore.gcs.credentialsSecret.secretName` | `string` | yes |  |  |
| `spec.artifactStore.gcs.credentialsSecret.secretKey` | `string` |  | `credentials.json` |  |
| `spec.artifactStore.azureBlob` | `KubernetesMlflowAzureBlobArtifacts` |  |  |  |
| `spec.artifactStore.azureBlob.storageAccount` | `string` | yes |  |  |
| `spec.artifactStore.azureBlob.container` | `string` | yes |  |  |
| `spec.artifactStore.azureBlob.prefix` | `string` |  |  |  |
| `spec.artifactStore.azureBlob.credentialsSecret` | `KubernetesMlflowAzureCredentialsSecret` | yes |  |  |
| `spec.artifactStore.azureBlob.credentialsSecret.secretName` | `string` | yes |  |  |
| `spec.artifactStore.azureBlob.credentialsSecret.secretKey` | `string` |  | `access_key` |  |
| `spec.auth` | `KubernetesMlflowAuth` |  |  |  |
| `spec.auth.enabled` | `bool` |  | `true` |  |
| `spec.auth.adminUsername` | `string` |  | `admin` |  |
| `spec.auth.adminPasswordSecret` | `KubernetesMlflowSecretKeyRef` |  |  |  |
| `spec.auth.adminPasswordSecret.secretName` | `string` | yes |  |  |
| `spec.auth.adminPasswordSecret.secretKey` | `string` |  | `password` |  |
| `spec.auth.defaultPermission` | `string` |  | `READ` |  |
| `spec.gc` | `KubernetesMlflowGc` |  |  |  |
| `spec.gc.enabled` | `bool` |  |  |  |
| `spec.gc.schedule` | `string` |  | `0 3 * * *` |  |
| `spec.gc.olderThan` | `string` |  | `30d` |  |
| `spec.service` | `KubernetesMlflowService` |  |  |  |
| `spec.service.type` | `string` |  | `ClusterIP` |  |
| `spec.service.annotations` | `map<string, string>` |  |  |  |
| `spec.metrics` | `KubernetesMlflowMetrics` |  |  |  |
| `spec.metrics.enabled` | `bool` |  |  |  |
| `spec.metrics.serviceMonitorEnabled` | `bool` |  |  |  |
| `spec.scheduling` | `KubernetesMlflowScheduling` |  |  |  |
| `spec.scheduling.nodeSelector` | `map<string, string>` |  |  |  |
| `spec.scheduling.tolerations` | `[]WorkloadToleration` |  |  |  |
| `spec.scheduling.tolerations[].key` | `string` |  |  |  |
| `spec.scheduling.tolerations[].operator` | `string` |  |  |  |
| `spec.scheduling.tolerations[].value` | `string` |  |  |  |
| `spec.scheduling.tolerations[].effect` | `string` |  |  |  |
| `spec.scheduling.tolerations[].tolerationSeconds` | `int64` |  |  |  |
| `spec.extraEnv` | `map<string, string>` |  |  |  |
| `spec.extraEnvFromSecret` | `map<string, KubernetesMlflowSecretKeyRef>` |  |  |  |
| `spec.extraEnvFromSecret.*.secretName` | `string` | yes |  |  |
| `spec.extraEnvFromSecret.*.secretKey` | `string` |  | `password` |  |
| `spec.extraArgs` | `[]string` |  |  |  |

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
governance labels) before installing and deleted with the
resource. When false, the namespace must already exist. KNOW
THIS: the module reads the database credential Secret to compose
the connection URI, and a Secret can only be read from the
workload's OWN namespace — co-locate MLflow with its database
(the default composition) or replicate the credential Secret
into this namespace.

### spec.server

`KubernetesMlflowServer`

The tracking server container: replicas, image and sizing.

### spec.server.replicas

`int32` · optional (explicit presence)

Server replicas. Empty = 1. More than one requires the postgres
backend AND an object artifact store (enforced) — the server
itself is stateless once both states live externally.

- default: `1`
- rule: {"int32":{"gte":1}}

### spec.server.image

`KubernetesMlflowImage`

Container image. Empty = the official image
(`ghcr.io/mlflow/mlflow`, `-full` variant at this kind's pinned
MLflow release) — override for private mirrors or a custom
image. The variant matters (verified live: the pod crash-loops
otherwise): the bare `vX.Y.Z` image is core MLflow only — no
database drivers (psycopg2/PyMySQL), no object-store clients
(boto3/GCS/Azure), and no basic-auth dependency (Flask-WTF) —
so the postgres/mysql backends, every remote artifact store,
and the secured-by-default sign-in all require the `-full`
variant (published upstream since v3.9.0). A mirror or custom
image must carry the same dependency set.

### spec.server.image.repository

`string` · optional (explicit presence)

Image repository including any registry host. Empty =
"ghcr.io/mlflow/mlflow".

- default: `ghcr.io/mlflow/mlflow`

### spec.server.image.tag

`string` · optional (explicit presence)

Image tag. Empty = "v3.15.0-full" (the MLflow release this kind
is built against, in the `-full` variant that ships the database
drivers, object-store clients and auth dependency — the bare
variant boots none of those arms; verified live at the pin).

- default: `v3.15.0-full`

### spec.server.workers

`int32` · optional (explicit presence)

Number of uvicorn worker processes. Empty = the server default
(4). More workers = more concurrent API requests per pod.

- rule: {"int32":{"lte":64,"gte":1}}

### spec.server.resources

`ContainerResources`

CPU/memory for the server container. Empty = no requests (fine
for dev; size for production — tracking servers are mostly
memory-bound on large experiment queries).

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

### spec.backendStore

`KubernetesMlflowBackendStore`

The backend store — experiments, runs, metrics, the model
registry, and (with auth on) users and permissions. Empty =
sqlite on a PVC (single replica only — fine for a team server;
declare postgres for production).

### spec.backendStore.sqlitePvc

`KubernetesMlflowSqlitePvc`

sqlite on a dedicated PVC — zero external dependencies,
single replica only.

### spec.backendStore.sqlitePvc.storageSize

`string` · optional (explicit presence)

PVC size. Empty = "5Gi". Metric time-series for large
experiment fleets grow — move to postgres before this fills.

- default: `5Gi`
- rule: {"string":{"pattern":"^[0-9]+(\\.[0-9]+)?(Ei|Pi|Ti|Gi|Mi|Ki|E|P|T|G|M|k|m)?$"}}

### spec.backendStore.sqlitePvc.storageClass

`string`

Storage class. Empty = the cluster default.

### spec.backendStore.postgres

`KubernetesMlflowPostgres`

PostgreSQL — the production path. A KubernetesPostgres
composes naturally; any reachable PostgreSQL (RDS, Cloud SQL)
works with literal values.

### spec.backendStore.postgres.host

`string | valueFrom` · required

Database server host. Defaults compose a KubernetesPostgres
resource's read-write Service.

- references: KubernetesPostgres (`status.outputs.rw_service`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesPostgres, name: <that resource's name>, fieldPath: status.outputs.rw_service}} -- a bare string does not parse

### spec.backendStore.postgres.port

`int32` · optional (explicit presence)

Database server port. Empty = 5432.

- default: `5432`
- rule: {"int32":{"lte":65535,"gte":1}}

### spec.backendStore.postgres.databaseName

`string` · optional (explicit presence)

Database name holding MLflow's tracking state (and, with auth
on, its user/permission tables). Must EXIST before install —
on a KubernetesPostgres, declare it at bootstrap
(`initdb.database`). Empty = "mlflow".

- default: `mlflow`
- rule: {"string":{"pattern":"^[a-zA-Z_][a-zA-Z0-9_$]*$"}}

### spec.backendStore.postgres.username

`string` · optional (explicit presence)

Database user with full rights inside `database_name` (the
server creates and migrates its tables at startup; ownership is
simplest). Empty = "mlflow".

- default: `mlflow`

### spec.backendStore.postgres.passwordSecret

`KubernetesMlflowPasswordSecret` · required

The Secret holding the user's password (composed into the URI
Secret at apply time — never rendered). Defaults compose a
KubernetesPostgres resource's application-user Secret.
Same-namespace constraint applies.

- rule: {"required":true}

### spec.backendStore.postgres.passwordSecret.secretName

`string | valueFrom` · required

Name of the Secret. Defaults compose a KubernetesPostgres
resource's application-user Secret (`<cluster>-app`).
Same-namespace constraint applies.

- references: KubernetesPostgres (`status.outputs.password_secret.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesPostgres, name: <that resource's name>, fieldPath: status.outputs.password_secret.name}} -- a bare string does not parse

### spec.backendStore.postgres.passwordSecret.secretKey

`string` · optional (explicit presence)

Key inside the Secret holding the password. Empty = "password"
(the KubernetesPostgres application-Secret convention).

- default: `password`

### spec.backendStore.mysql

`KubernetesMlflowMysql`

MySQL 8+ — a KubernetesMysql composes naturally.

### spec.backendStore.mysql.host

`string | valueFrom` · required

Database server host. Defaults compose a KubernetesMysql
resource's client Service.

- references: KubernetesMysql (`status.outputs.primary_service`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesMysql, name: <that resource's name>, fieldPath: status.outputs.primary_service}} -- a bare string does not parse

### spec.backendStore.mysql.port

`int32` · optional (explicit presence)

Database server port. Empty = 3306.

- default: `3306`
- rule: {"int32":{"lte":65535,"gte":1}}

### spec.backendStore.mysql.databaseName

`string` · optional (explicit presence)

Database name holding MLflow's tracking state. Must EXIST before
install. Empty = "mlflow".

- default: `mlflow`
- rule: {"string":{"pattern":"^[a-zA-Z_][a-zA-Z0-9_$]*$"}}

### spec.backendStore.mysql.username

`string` · optional (explicit presence)

Database user with full rights inside `database_name`. Empty =
"mlflow".

- default: `mlflow`

### spec.backendStore.mysql.passwordSecret

`KubernetesMlflowMysqlPasswordSecret` · required

The Secret holding the user's password. Defaults compose a
KubernetesMysql resource's root credential Secret; for a
dedicated MLflow user (recommended), point at that user's
Secret. Same-namespace constraint applies.

- rule: {"required":true}

### spec.backendStore.mysql.passwordSecret.secretName

`string | valueFrom` · required

Name of the Secret. Defaults compose a KubernetesMysql
resource's root credential Secret. Same-namespace constraint
applies.

- references: KubernetesMysql (`status.outputs.root_password_secret.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesMysql, name: <that resource's name>, fieldPath: status.outputs.root_password_secret.name}} -- a bare string does not parse

### spec.backendStore.mysql.passwordSecret.secretKey

`string` · optional (explicit presence)

Key inside the Secret holding the password. Empty = "password".

- default: `password`

### spec.artifactStore

`KubernetesMlflowArtifactStore`

The artifact store — where models, datasets and other run
artifacts physically live. The server proxies all artifact
traffic, so clients never need these credentials. Empty = a
local PVC (artifacts live with the server — fine to start;
declare an object store for production durability and size).

### spec.artifactStore.pvc

`KubernetesMlflowPvcArtifacts`

Artifacts on a local PVC, served through the tracking server —
zero external dependencies, single replica only.

### spec.artifactStore.pvc.storageSize

`string` · optional (explicit presence)

PVC size. Empty = "10Gi". Models are big — size for your
artifact volume or use an object store.

- default: `10Gi`
- rule: {"string":{"pattern":"^[0-9]+(\\.[0-9]+)?(Ei|Pi|Ti|Gi|Mi|Ki|E|P|T|G|M|k|m)?$"}}

### spec.artifactStore.pvc.storageClass

`string`

Storage class. Empty = the cluster default.

### spec.artifactStore.s3Compatible

`KubernetesMlflowS3CompatibleArtifacts`

Any S3-compatible object store — a KubernetesSeaweedFs
composes naturally (endpoint and credentials FK onto its
outputs); MinIO-compatible endpoints work with literal
values.

### spec.artifactStore.s3Compatible.endpoint

`string | valueFrom` · required

The S3 API endpoint (e.g.
`http://seaweedfs-s3.<ns>.svc.cluster.local:8333`). Defaults
compose a KubernetesSeaweedFs resource's S3 endpoint.

- references: KubernetesSeaweedFs (`status.outputs.s3_endpoint`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesSeaweedFs, name: <that resource's name>, fieldPath: status.outputs.s3_endpoint}} -- a bare string does not parse

### spec.artifactStore.s3Compatible.bucket

`string` · required

Bucket holding the artifacts (must exist — on a
KubernetesSeaweedFs, declare it under `s3.buckets`). Rendered as
`s3://<bucket>[/<prefix>]`.

- rule: {"required":true}

### spec.artifactStore.s3Compatible.prefix

`string`

Key prefix within the bucket. Empty = the bucket root.

### spec.artifactStore.s3Compatible.credentialsSecret

`KubernetesMlflowS3CredentialsSecret` · required

The Secret holding the access credentials. Defaults compose a
KubernetesSeaweedFs resource's S3 credentials Secret
(`<name>-s3-secret`). Same-namespace constraint applies.
Addressing style is automatic: MLflow's S3 client (boto3) uses
path-style for custom endpoints — SeaweedFS and MinIO-compatible
stores work without configuration.

- rule: {"required":true}

### spec.artifactStore.s3Compatible.credentialsSecret.secretName

`string | valueFrom` · required

Name of the Secret. Defaults compose a KubernetesSeaweedFs
resource's S3 credentials Secret. Same-namespace constraint
applies.

- references: KubernetesSeaweedFs (`status.outputs.s3_credentials_secret_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesSeaweedFs, name: <that resource's name>, fieldPath: status.outputs.s3_credentials_secret_name}} -- a bare string does not parse

### spec.artifactStore.s3Compatible.credentialsSecret.accessKeyIdKey

`string` · optional (explicit presence)

Key NAME within the Secret (a reference, not secret material)
holding the access key ID. Empty = "admin_access_key_id"
(the KubernetesSeaweedFs convention).

- default: `admin_access_key_id`

### spec.artifactStore.s3Compatible.credentialsSecret.secretAccessKeyKey

`string` · optional (explicit presence)

Key holding the secret access key. Empty =
"admin_secret_access_key".

- default: `admin_secret_access_key`

### spec.artifactStore.awsS3

`KubernetesMlflowAwsS3Artifacts`

AWS S3. Keyless (IRSA/pod identity — annotate the
ServiceAccount via `helm-less` extra config) or declared keys.

### spec.artifactStore.awsS3.bucket

`string` · required

Bucket name.

- rule: {"required":true}

### spec.artifactStore.awsS3.prefix

`string`

Key prefix within the bucket. Empty = the bucket root.

### spec.artifactStore.awsS3.region

`string` · required

AWS region of the bucket (e.g. "us-west-2").

- rule: {"required":true}

### spec.artifactStore.awsS3.credentialsSecret

`KubernetesMlflowAwsCredentialsSecret`

The Secret holding AWS access keys (keys
`access_key_id`/`secret_access_key` by default). Empty = KEYLESS
— the pod's ambient identity (IRSA / EKS Pod Identity) must
grant bucket access.

### spec.artifactStore.awsS3.credentialsSecret.secretName

`string` · required

Name of the Secret. Same-namespace constraint applies.

- rule: {"required":true}

### spec.artifactStore.awsS3.credentialsSecret.accessKeyIdKey

`string` · optional (explicit presence)

Key NAME within the Secret (a reference, not secret material)
holding the access key ID. Empty = "access_key_id".

- default: `access_key_id`

### spec.artifactStore.awsS3.credentialsSecret.secretAccessKeyKey

`string` · optional (explicit presence)

Key holding the secret access key. Empty = "secret_access_key".

- default: `secret_access_key`

### spec.artifactStore.gcs

`KubernetesMlflowGcsArtifacts`

Google Cloud Storage — keyless via Workload Identity, or a
service-account key Secret.

### spec.artifactStore.gcs.bucket

`string` · required

Bucket name (rendered as `gs://<bucket>[/<prefix>]`).

- rule: {"required":true}

### spec.artifactStore.gcs.prefix

`string`

Key prefix within the bucket. Empty = the bucket root.

### spec.artifactStore.gcs.credentialsSecret

`KubernetesMlflowGcpCredentialsSecret`

The Secret holding a service-account JSON key (key
`credentials.json` by default; mounted and wired via
GOOGLE_APPLICATION_CREDENTIALS). Empty = KEYLESS — the pod's
ambient identity (GKE Workload Identity) must grant bucket
access.

### spec.artifactStore.gcs.credentialsSecret.secretName

`string` · required

Name of the Secret. Same-namespace constraint applies.

- rule: {"required":true}

### spec.artifactStore.gcs.credentialsSecret.secretKey

`string` · optional (explicit presence)

Key holding the service-account JSON. Empty =
"credentials.json".

- default: `credentials.json`

### spec.artifactStore.azureBlob

`KubernetesMlflowAzureBlobArtifacts`

Azure Blob Storage with a storage-account access key.

### spec.artifactStore.azureBlob.storageAccount

`string` · required

Storage account name.

- rule: {"required":true}

### spec.artifactStore.azureBlob.container

`string` · required

Blob container name.

- rule: {"required":true}

### spec.artifactStore.azureBlob.prefix

`string`

Key prefix within the container. Empty = the container root.

### spec.artifactStore.azureBlob.credentialsSecret

`KubernetesMlflowAzureCredentialsSecret` · required

The Secret holding the storage-account access key (key
`access_key` by default; wired via AZURE_STORAGE_ACCESS_KEY).

- rule: {"required":true}

### spec.artifactStore.azureBlob.credentialsSecret.secretName

`string` · required

Name of the Secret. Same-namespace constraint applies.

- rule: {"required":true}

### spec.artifactStore.azureBlob.credentialsSecret.secretKey

`string` · optional (explicit presence)

Key holding the access key. Empty = "access_key".

- default: `access_key`

### spec.auth

`KubernetesMlflowAuth`

Basic authentication. Empty = ENABLED with a module-generated
admin password — never the open server, never upstream's
admin/password1234 default.

### spec.auth.enabled

`bool` · optional (explicit presence)

Require authentication. Empty = true — the open server never
ships by default. Disabling leaves EVERY experiment and model
writable by anyone who can reach the Service.

- default: `true`

### spec.auth.adminUsername

`string` · optional (explicit presence)

The bootstrap admin username. Empty = "admin".

- default: `admin`

### spec.auth.adminPasswordSecret

`KubernetesMlflowSecretKeyRef`

Existing Secret holding the admin password. Empty = the module
GENERATES it into `<name>-admin-auth` (key `password`) —
upstream's admin/password1234 default never ships. Either way
the password reaches the server through its auth configuration
Secret, never a rendered manifest.

### spec.auth.adminPasswordSecret.secretName

`string` · required

Name of the Secret. Same-namespace constraint applies.

- rule: {"required":true}

### spec.auth.adminPasswordSecret.secretKey

`string` · optional (explicit presence)

Key inside the Secret. Empty = "password".

- default: `password`

### spec.auth.defaultPermission

`string` · optional (explicit presence)

Permission newly-created users get on resources they did not
create. Empty = "READ" (upstream default). NO_PERMISSIONS makes
every experiment private to its creator until shared.

- default: `READ`
- rule: {"string":{"in":["READ","EDIT","MANAGE","NO_PERMISSIONS"]}}

### spec.gc

`KubernetesMlflowGc`

Periodic hard-deletion of soft-deleted runs and experiments
(`mlflow gc`). Off by default — deleted runs stay restorable
until you enable it.

### spec.gc.enabled

`bool`

Run `mlflow gc` on a schedule (a CronJob). Deleted runs become
UNRECOVERABLE once collected.

### spec.gc.schedule

`string` · optional (explicit presence)

Cron schedule. Empty = "0 3 * * *" (daily at 03:00).

- default: `0 3 * * *`
- rule: {"string":{"pattern":"^(@(annually|yearly|monthly|weekly|daily|hourly))|((((\\d+,)+\\d+|(\\d+([/\\-])\\d+)|\\d+|\\*(/\\d+)?) ?){5})$"}}

### spec.gc.olderThan

`string` · optional (explicit presence)

Only collect runs deleted longer ago than this (e.g. "30d",
"72h") — the undo window. Empty = "30d".

- default: `30d`
- rule: {"string":{"pattern":"^[0-9]+(d|h|m|s)$"}}

### spec.service

`KubernetesMlflowService`

The MLflow Service — this kind keeps it ClusterIP (compose
exposure kinds over the exported handle); the annotations
surface carries cloud LB configuration when you flip the type.

### spec.service.type

`string` · optional (explicit presence)

Service type. Empty = "ClusterIP" — compose exposure kinds over
the exported handle instead of exposing directly.

- default: `ClusterIP`
- rule: {"string":{"in":["ClusterIP","LoadBalancer","NodePort"]}}

### spec.service.annotations

`map<string, string>`

Service annotations — the cloud load-balancer configuration
surface when `type` is LoadBalancer.

### spec.metrics

`KubernetesMlflowMetrics`

Prometheus metrics: expose them on the server and optionally
render a ServiceMonitor for operator-based scraping (requires
the Prometheus operator CRDs — a KubernetesKubePrometheusStack
composes naturally).

- rule: A ServiceMonitor scrapes the metrics endpoint — enable metrics.enabled too, or remove service_monitor_enabled.

### spec.metrics.enabled

`bool`

Expose Prometheus metrics on the server (`/metrics`).

### spec.metrics.serviceMonitorEnabled

`bool`

Render a ServiceMonitor for operator-based scraping. Requires
the Prometheus operator CRDs on the cluster (a
KubernetesKubePrometheusStack composes naturally) — enabling it
without them fails the deploy.

### spec.scheduling

`KubernetesMlflowScheduling`

Pod scheduling for the server (and GC job) pods.

### spec.scheduling.nodeSelector

`map<string, string>`

Node selector for the server (and GC) pods.

### spec.scheduling.tolerations

`[]WorkloadToleration`

Tolerations for the server (and GC) pods.

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

### spec.extraEnv

`map<string, string>`

Extra environment variables injected into the server container
(plain values — e.g. MLFLOW_* tuning flags). For secret values
use `extra_env_from_secret`.

### spec.extraEnvFromSecret

`map<string, KubernetesMlflowSecretKeyRef>`

Extra environment variables sourced from existing Secrets —
name → Secret reference. The escape hatch for auth/storage
integrations beyond the typed arms.

### spec.extraEnvFromSecret.*.secretName

`string` · required

Name of the Secret. Same-namespace constraint applies.

- rule: {"required":true}

### spec.extraEnvFromSecret.*.secretKey

`string` · optional (explicit presence)

Key inside the Secret. Empty = "password".

- default: `password`

### spec.extraArgs

`[]string`

Extra command-line arguments appended to `mlflow server` — the
escape hatch for server flags the typed fields do not model
(e.g. --gunicorn-opts). Never put secret material here;
arguments are visible in the rendered pod spec.

## Validation Rules

- `server.replicas_need_postgres`: Running more than one MLflow replica requires the postgres backend store — sqlite is a single-writer file on a ReadWriteOnce volume and a second replica would corrupt it (this also covers auth state, which follows the backend store).
- `server.replicas_need_object_artifacts`: Running more than one MLflow replica requires an object artifact store (s3_compatible, aws_s3, gcs or azure_blob) — the PVC artifact arm is ReadWriteOnce and mounts on a single pod.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesMlflow, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.namespace` | `string` | Namespace MLflow runs in. |
| `status.outputs.service` | `string` | Name of the MLflow Service (`<name>`) — the handle exposure kinds route to. |
| `status.outputs.tracking_endpoint` | `string` | In-cluster tracking endpoint, `http://<name>.<namespace>.svc.cluster.local:5000` — set it as MLFLOW_TRACKING_URI in training jobs and pipelines (with auth on, pair it with MLFLOW_TRACKING_USERNAME/PASSWORD). |
| `status.outputs.admin_password_secret` | `KubernetesSecretKey` | The admin credential: the Secret and key holding the bootstrap admin password (module-generated `<name>-admin-auth` unless an existing Secret was declared). Username lives in the spec (`auth.admin_username`, default "admin"). Empty when auth is disabled. |
| `status.outputs.admin_password_secret.name` | `string` | The name of the Kubernetes Secret. |
| `status.outputs.admin_password_secret.key` | `string` | The key within the Kubernetes Secret. |
| `status.outputs.backend_store_uri_secret_name` | `string` | Name of the module-owned Secret holding the backend-store connection URI (key `uri`) — empty on the sqlite arm (no external database to share). |
| `status.outputs.port_forward_command` | `string` | Port-forward command for reaching the MLflow UI from a workstation when no exposure is composed (`kubectl port-forward svc/<name> -n <namespace> 5000:5000`). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespace` | KubernetesNamespace | `spec.name` |
| `spec.backendStore.postgres.host` | KubernetesPostgres | `status.outputs.rw_service` |
| `spec.backendStore.postgres.passwordSecret.secretName` | KubernetesPostgres | `status.outputs.password_secret.name` |
| `spec.backendStore.mysql.host` | KubernetesMysql | `status.outputs.primary_service` |
| `spec.backendStore.mysql.passwordSecret.secretName` | KubernetesMysql | `status.outputs.root_password_secret.name` |
| `spec.artifactStore.s3Compatible.endpoint` | KubernetesSeaweedFs | `status.outputs.s3_endpoint` |
| `spec.artifactStore.s3Compatible.credentialsSecret.secretName` | KubernetesSeaweedFs | `status.outputs.s3_credentials_secret_name` |

## See Also

- [Overview](../README.md)
