# KubernetesSuperset

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1alpha1`

**KubernetesSupersetSpec** deploys Apache Superset — the
business-intelligence platform where analysts explore data, build
charts and assemble dashboards over any SQL database
(https://superset.apache.org).

WHAT GETS INSTALLED: the official ASF superset Helm chart renders
the web application (Flask/gunicorn on port 8088), a Celery worker
fleet for async queries/thumbnails/reports, optional beat
(schedules), flower (Celery monitoring), websocket and MCP
deployments, and an init Job that migrates the metadata schema and
bootstraps the admin user.

TWO EXTERNAL BACKENDS, COMPOSITION-FIRST:

  - The METADATA DATABASE (dashboards, charts, users, saved
    queries) is REQUIRED and always external — a
    KubernetesPostgres composes naturally. The chart's bundled
    PostgreSQL subchart rides a frozen legacy image line and
    never ships from this kind.
  - The CACHE/BROKER (query cache, Celery queues, async results)
    is a Redis-protocol store — a KubernetesValkey composes
    naturally. Absent = cache disabled: a web-only Superset where
    every query runs synchronously and workers, beat, flower,
    websockets and MCP must stay off. The bundled Redis subchart
    (same frozen legacy line) never ships.

SECURED BY DEFAULT: the session-signing SECRET_KEY is
module-generated (Superset refuses to start on its insecure
default) and STABLE — rotating it logs out every session and
orphans the encrypted datasource credentials stored in the
metadata database, so it is generated once and kept. The
bootstrap admin password is module-generated into
`<name>-admin-auth` — the chart's documented admin/admin default
never ships.

CREDENTIALS ARE SECRET-NATIVE: the chart consumes ALL runtime
credentials through one environment Secret; this kind turns the
chart's own Secret OFF and composes `<name>-env` AT APPLY TIME
from the referenced credential Secrets — nothing
credential-bearing appears in rendered values or manifests.

EXPOSURE: the Service stays ClusterIP; expose it via first-class
kinds over the exported service handle.

## Example

```yaml
# Full-surface shape: the composed external metadata database with SSL,
# an authed external cache with Celery workers + beat + flower, the
# websocket and MCP arms, feature flags, config overrides reading
# env-sourced secrets, a bootstrap driver install and the env escape
# hatches — the offline plan/preview proof for the widest typed
# rendering.
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesSuperset
metadata:
  name: superset-full
spec:
  namespace:
    value: superset-full
  createNamespace: true
  metadataDatabase:
    host:
      value: superset-pg-rw
    port: 5432
    databaseName: superset
    username: superset
    passwordSecret:
      secretName:
        value: superset-pg-app
      secretKey: password
    ssl:
      enabled: true
      mode: require
  cache:
    host:
      value: superset-valkey
    port: 6379
    passwordSecret:
      secretName:
        value: superset-valkey-password
      secretKey: password
    cacheDb: 1
    celeryDb: 0
  web:
    hpa:
      minReplicas: 2
      maxReplicas: 6
      targetCpuUtilizationPercent: 70
    resources:
      requests:
        cpu: 500m
        memory: 1Gi
      limits:
        memory: 3Gi
  worker:
    replicas: 2
    resources:
      requests:
        cpu: 500m
        memory: 1Gi
      limits:
        memory: 3Gi
  beat:
    enabled: true
  flower:
    enabled: true
  websockets:
    enabled: true
    image:
      tag: "0.1.4"
    replicas: 2
  mcp:
    enabled: true
  init:
    admin:
      username: bi-admin
      email: bi-admin@example.com
    loadExamples: false
  featureFlags:
    DASHBOARD_RBAC: true
    ALERT_REPORTS: true
    GLOBAL_ASYNC_QUERIES: true
  configOverrides:
    mapbox: |
      MAPBOX_API_KEY = env('MAPBOX_API_KEY', '')
  extraEnvFromSecret:
    MAPBOX_API_KEY:
      secretName: superset-mapbox
      secretKey: api_key
  extraEnv:
    GUNICORN_TIMEOUT: "300"
  bootstrapScript: |
    #!/bin/bash
    if [ ! -f ~/bootstrap ]; then pip install trino elasticsearch-dbapi; echo "done" > ~/bootstrap; fi
  service:
    type: ClusterIP
  scheduling:
    nodeSelector:
      workload: analytics
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.namespace` | `string \| valueFrom` | yes |  | KubernetesNamespace (`spec.name`) |
| `spec.createNamespace` | `bool` |  |  |  |
| `spec.image` | `KubernetesSupersetImage` |  |  |  |
| `spec.image.repository` | `string` |  | `apachesuperset.docker.scarf.sh/apache/superset` |  |
| `spec.image.tag` | `string` |  | `6.1.0` |  |
| `spec.imagePullSecrets` | `[]string` |  |  |  |
| `spec.metadataDatabase` | `KubernetesSupersetMetadataDatabase` | yes |  |  |
| `spec.metadataDatabase.host` | `string \| valueFrom` | yes |  | KubernetesPostgres (`status.outputs.rw_service`) |
| `spec.metadataDatabase.port` | `int32` |  | `5432` |  |
| `spec.metadataDatabase.databaseName` | `string` |  | `superset` |  |
| `spec.metadataDatabase.username` | `string` |  | `superset` |  |
| `spec.metadataDatabase.passwordSecret` | `KubernetesSupersetPostgresPasswordSecret` | yes |  |  |
| `spec.metadataDatabase.passwordSecret.secretName` | `string \| valueFrom` | yes |  | KubernetesPostgres (`status.outputs.password_secret.name`) |
| `spec.metadataDatabase.passwordSecret.secretKey` | `string` |  | `password` |  |
| `spec.metadataDatabase.ssl` | `KubernetesSupersetDatabaseSsl` |  |  |  |
| `spec.metadataDatabase.ssl.enabled` | `bool` |  |  |  |
| `spec.metadataDatabase.ssl.mode` | `string` |  | `require` |  |
| `spec.cache` | `KubernetesSupersetCache` |  |  |  |
| `spec.cache.host` | `string \| valueFrom` | yes |  | KubernetesValkey (`status.outputs.service`) |
| `spec.cache.port` | `int32` |  | `6379` |  |
| `spec.cache.username` | `string` |  |  |  |
| `spec.cache.passwordSecret` | `KubernetesSupersetCachePasswordSecret` |  |  |  |
| `spec.cache.passwordSecret.secretName` | `string \| valueFrom` | yes |  | KubernetesValkey (`status.outputs.password_secret.name`) |
| `spec.cache.passwordSecret.secretKey` | `string` |  | `password` |  |
| `spec.cache.cacheDb` | `int32` |  |  |  |
| `spec.cache.celeryDb` | `int32` |  |  |  |
| `spec.secretKeySecret` | `KubernetesSupersetSecretKeyRef` |  |  |  |
| `spec.secretKeySecret.secretName` | `string` | yes |  |  |
| `spec.secretKeySecret.secretKey` | `string` | yes |  |  |
| `spec.web` | `KubernetesSupersetWeb` |  |  |  |
| `spec.web.replicas` | `int32` |  | `1` |  |
| `spec.web.hpa` | `KubernetesSupersetHpa` |  |  |  |
| `spec.web.hpa.minReplicas` | `int32` |  |  |  |
| `spec.web.hpa.maxReplicas` | `int32` |  |  |  |
| `spec.web.hpa.targetCpuUtilizationPercent` | `int32` |  |  |  |
| `spec.web.resources` | `ContainerResources` |  |  |  |
| `spec.web.resources.limits` | `CpuMemory` |  |  |  |
| `spec.web.resources.limits.cpu` | `string` |  |  |  |
| `spec.web.resources.limits.memory` | `string` |  |  |  |
| `spec.web.resources.requests` | `CpuMemory` |  |  |  |
| `spec.web.resources.requests.cpu` | `string` |  |  |  |
| `spec.web.resources.requests.memory` | `string` |  |  |  |
| `spec.worker` | `KubernetesSupersetWorker` |  |  |  |
| `spec.worker.enabled` | `bool` |  | `true` |  |
| `spec.worker.replicas` | `int32` |  | `1` |  |
| `spec.worker.hpa` | `KubernetesSupersetHpa` |  |  |  |
| `spec.worker.hpa.minReplicas` | `int32` |  |  |  |
| `spec.worker.hpa.maxReplicas` | `int32` |  |  |  |
| `spec.worker.hpa.targetCpuUtilizationPercent` | `int32` |  |  |  |
| `spec.worker.resources` | `ContainerResources` |  |  |  |
| `spec.worker.resources.limits` | `CpuMemory` |  |  |  |
| `spec.worker.resources.limits.cpu` | `string` |  |  |  |
| `spec.worker.resources.limits.memory` | `string` |  |  |  |
| `spec.worker.resources.requests` | `CpuMemory` |  |  |  |
| `spec.worker.resources.requests.cpu` | `string` |  |  |  |
| `spec.worker.resources.requests.memory` | `string` |  |  |  |
| `spec.beat` | `KubernetesSupersetBeat` |  |  |  |
| `spec.beat.enabled` | `bool` |  |  |  |
| `spec.flower` | `KubernetesSupersetFlower` |  |  |  |
| `spec.flower.enabled` | `bool` |  |  |  |
| `spec.websockets` | `KubernetesSupersetWebsockets` |  |  |  |
| `spec.websockets.enabled` | `bool` |  |  |  |
| `spec.websockets.image` | `KubernetesSupersetWsImage` |  |  |  |
| `spec.websockets.image.repository` | `string` |  | `oneacrefund/superset-websocket` |  |
| `spec.websockets.image.tag` | `string` |  | `latest` |  |
| `spec.websockets.replicas` | `int32` |  | `1` |  |
| `spec.mcp` | `KubernetesSupersetMcp` |  |  |  |
| `spec.mcp.enabled` | `bool` |  |  |  |
| `spec.init` | `KubernetesSupersetInit` |  |  |  |
| `spec.init.admin` | `KubernetesSupersetAdminUser` |  |  |  |
| `spec.init.admin.username` | `string` |  | `admin` |  |
| `spec.init.admin.email` | `string` |  | `admin@superset.local` |  |
| `spec.init.admin.passwordSecret` | `KubernetesSupersetSecretKeyRef` |  |  |  |
| `spec.init.admin.passwordSecret.secretName` | `string` | yes |  |  |
| `spec.init.admin.passwordSecret.secretKey` | `string` | yes |  |  |
| `spec.init.loadExamples` | `bool` |  |  |  |
| `spec.featureFlags` | `map<string, bool>` |  |  |  |
| `spec.configOverrides` | `map<string, string>` |  |  |  |
| `spec.extraEnv` | `map<string, string>` |  |  |  |
| `spec.extraEnvFromSecret` | `map<string, KubernetesSupersetSecretKeyRef>` |  |  |  |
| `spec.extraEnvFromSecret.*.secretName` | `string` | yes |  |  |
| `spec.extraEnvFromSecret.*.secretKey` | `string` | yes |  |  |
| `spec.bootstrapScript` | `string` |  |  |  |
| `spec.service` | `KubernetesSupersetService` |  |  |  |
| `spec.service.type` | `string` |  | `ClusterIP` |  |
| `spec.service.annotations` | `map<string, string>` |  |  |  |
| `spec.scheduling` | `KubernetesSupersetScheduling` |  |  |  |
| `spec.scheduling.nodeSelector` | `map<string, string>` |  |  |  |
| `spec.scheduling.tolerations` | `[]WorkloadToleration` |  |  |  |
| `spec.scheduling.tolerations[].key` | `string` |  |  |  |
| `spec.scheduling.tolerations[].operator` | `string` |  |  |  |
| `spec.scheduling.tolerations[].value` | `string` |  |  |  |
| `spec.scheduling.tolerations[].effect` | `string` |  |  |  |
| `spec.scheduling.tolerations[].tolerationSeconds` | `int64` |  |  |  |
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
governance labels) before installing and deleted with the
resource. When false, the namespace must already exist. KNOW
THIS: the module reads the database (and cache) credential
Secrets to compose the environment Secret, and a Secret can
only be read from the workload's OWN namespace — co-locate
Superset with its database (the default composition) or
replicate the credential Secrets into this namespace.

### spec.image

`KubernetesSupersetImage`

Container image. Empty = the official image the chart pins
(`apachesuperset.docker.scarf.sh/apache/superset` — the ASF's
scarf.sh download gateway in front of the official image;
override the repository for air-gapped mirrors or a custom
image). KNOW THIS (verified live and in the image's own build
file): the published image is the driver-less "lean" build stage
— even the PostgreSQL metadata-database driver rides only the
dev/ci variants. The module's default bootstrap script installs
the exact psycopg2 pin at container start; for production
(air-gap, no pip-at-boot), bake a custom image with the drivers
and set bootstrap_script to a no-op.

### spec.image.repository

`string` · optional (explicit presence)

Image repository including any registry host. Empty =
"apachesuperset.docker.scarf.sh/apache/superset" (the chart
default — the ASF's scarf gateway to the official image).

- default: `apachesuperset.docker.scarf.sh/apache/superset`

### spec.image.tag

`string` · optional (explicit presence)

Image tag. Empty = "6.1.0" (the Superset release this kind is
built against).

- default: `6.1.0`

### spec.imagePullSecrets

`[]string`

Names of image-pull Secrets in the same namespace, for pulling
from private mirrors.

### spec.metadataDatabase

`KubernetesSupersetMetadataDatabase` · required

The metadata database — dashboards, charts, users, saved
queries, and the ENCRYPTED datasource credentials. REQUIRED and
always external; a KubernetesPostgres composes naturally.

- rule: {"required":true}

### spec.metadataDatabase.host

`string | valueFrom` · required

Database server host. Defaults compose a KubernetesPostgres
resource's read-write Service.

- references: KubernetesPostgres (`status.outputs.rw_service`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesPostgres, name: <that resource's name>, fieldPath: status.outputs.rw_service}} -- a bare string does not parse

### spec.metadataDatabase.port

`int32` · optional (explicit presence)

Database server port. Empty = 5432.

- default: `5432`
- rule: {"int32":{"lte":65535,"gte":1}}

### spec.metadataDatabase.databaseName

`string` · optional (explicit presence)

Database name holding Superset's metadata. Must EXIST before
install — on a KubernetesPostgres, declare it at bootstrap
(`initdb.database`). Empty = "superset".

- default: `superset`
- rule: {"string":{"pattern":"^[a-zA-Z_][a-zA-Z0-9_$]*$"}}

### spec.metadataDatabase.username

`string` · optional (explicit presence)

Database user with full rights inside `database_name` (the init
Job creates and migrates the schema). Empty = "superset".

- default: `superset`

### spec.metadataDatabase.passwordSecret

`KubernetesSupersetPostgresPasswordSecret` · required

The Secret holding the user's password (composed into the
environment Secret at apply time — never rendered). Defaults
compose a KubernetesPostgres resource's application-user
Secret. Same-namespace constraint applies.

- rule: {"required":true}

### spec.metadataDatabase.passwordSecret.secretName

`string | valueFrom` · required

Name of the Secret. Defaults compose a KubernetesPostgres
resource's application-user Secret (`<cluster>-app`).
Same-namespace constraint applies.

- references: KubernetesPostgres (`status.outputs.password_secret.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesPostgres, name: <that resource's name>, fieldPath: status.outputs.password_secret.name}} -- a bare string does not parse

### spec.metadataDatabase.passwordSecret.secretKey

`string` · optional (explicit presence)

Key inside the Secret holding the password. Empty = "password"
(the KubernetesPostgres application-Secret convention).

- default: `password`

### spec.metadataDatabase.ssl

`KubernetesSupersetDatabaseSsl`

Require SSL on the database connection (`sslmode`).

### spec.metadataDatabase.ssl.enabled

`bool`

Require SSL on the metadata-database connection.

### spec.metadataDatabase.ssl.mode

`string` · optional (explicit presence)

The `sslmode` to request. Empty = "require".

- default: `require`
- rule: {"string":{"in":["require","verify-ca","verify-full"]}}

### spec.cache

`KubernetesSupersetCache`

The Redis-protocol cache/broker — query cache, Celery queues,
async results. A KubernetesValkey composes naturally. Absent =
cache disabled (web-only Superset: synchronous queries, no
workers/beat/flower/websockets/MCP).

### spec.cache.host

`string | valueFrom` · required

Cache host. Defaults compose a KubernetesValkey resource's
client Service; any reachable Redis-protocol host works as a
literal value.

- references: KubernetesValkey (`status.outputs.service`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesValkey, name: <that resource's name>, fieldPath: status.outputs.service}} -- a bare string does not parse

### spec.cache.port

`int32` · optional (explicit presence)

Cache port. Empty = 6379.

- default: `6379`
- rule: {"int32":{"lte":65535,"gte":1}}

### spec.cache.username

`string`

ACL username, when the store uses named users (a
KubernetesValkey ACL user's name). Empty = password-only AUTH
(or no auth when `password_secret` is also empty).

### spec.cache.passwordSecret

`KubernetesSupersetCachePasswordSecret`

The Secret holding the cache password (composed into the
environment Secret at apply time; auth-dependent config blocks
read it from environment). Empty = an auth-less store. Defaults
compose a KubernetesValkey resource's password Secret when
referenced. Same-namespace constraint applies.

### spec.cache.passwordSecret.secretName

`string | valueFrom` · required

Name of the Secret. Defaults compose a KubernetesValkey
resource's password Secret. Same-namespace constraint applies.

- references: KubernetesValkey (`status.outputs.password_secret.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesValkey, name: <that resource's name>, fieldPath: status.outputs.password_secret.name}} -- a bare string does not parse

### spec.cache.passwordSecret.secretKey

`string` · optional (explicit presence)

Key inside the Secret holding the password. For a KubernetesValkey
auth Secret this is the ACL username ("default" unless you declared
users) — that Secret keys entries BY USERNAME, one key per user.
Empty = "password" (the generic existing-Secret convention; wrong
for Valkey auth Secrets, where an unset key deploys cleanly and
then fails authentication at runtime).

- default: `password`

### spec.cache.cacheDb

`int32` · optional (explicit presence)

Redis database number for the query cache. Empty = 1 (the
chart default).

- rule: {"int32":{"lte":15,"gte":0}}

### spec.cache.celeryDb

`int32` · optional (explicit presence)

Redis database number for Celery queues/results. Empty = 0
(the chart default).

- rule: {"int32":{"lte":15,"gte":0}}

### spec.secretKeySecret

`KubernetesSupersetSecretKeyRef`

The session-signing key (Flask SECRET_KEY — also encrypts
stored datasource credentials). Empty = module-generated into
`<name>-secret-key`, STABLE across applies. KNOW THIS: changing
the key logs out every session AND orphans stored datasource
credentials — rotate only via Superset's own
`superset re-encrypt-secrets` procedure with the old key in
PREVIOUS_SECRET_KEY.

### spec.secretKeySecret.secretName

`string` · required

Name of the Secret. Same-namespace constraint applies.

- rule: {"required":true}

### spec.secretKeySecret.secretKey

`string` · required

Key inside the Secret.

- rule: {"required":true}

### spec.web

`KubernetesSupersetWeb`

The web application (`supersetNode`) — the UI and REST API.

### spec.web.replicas

`int32` · optional (explicit presence)

Web replicas. Empty = 1. The web tier is stateless (state lives
in the metadata database and cache) — scale freely. Ignored
when `hpa` is set.

- default: `1`
- rule: {"int32":{"gte":1}}

### spec.web.hpa

`KubernetesSupersetHpa`

Horizontal Pod Autoscaler for the web tier (needs a metrics
server and resource REQUESTS).

### spec.web.hpa.minReplicas

`int32` · optional (explicit presence)

Lower replica bound. Empty = 1.

- rule: {"int32":{"gte":1}}

### spec.web.hpa.maxReplicas

`int32`

Upper replica bound. Required.

- rule: {"int32":{"gte":1}}

### spec.web.hpa.targetCpuUtilizationPercent

`int32` · optional (explicit presence)

Target average CPU utilization percentage (of the container's
CPU requests). Empty = 80 (the chart default).

- rule: {"int32":{"lte":100,"gte":1}}

### spec.web.resources

`ContainerResources`

CPU/memory for the web container. Empty = no requests (fine
for dev; size for production — chart rendering and SQL Lab are
memory-hungry).

### spec.web.resources.limits

`CpuMemory`

The resource limits for the container.
Specify the maximum amount of CPU and memory that the container can use.

### spec.web.resources.limits.cpu

`string`

### spec.web.resources.limits.memory

`string`

### spec.web.resources.requests

`CpuMemory`

The resource requests for the container.
Specify the minimum amount of CPU and memory that the container is guaranteed.

### spec.web.resources.requests.cpu

`string`

### spec.web.resources.requests.memory

`string`

### spec.worker

`KubernetesSupersetWorker`

The Celery worker fleet — async SQL Lab queries, thumbnails,
alerts and reports execution. Requires `cache`. Empty (with
cache set) = 1 worker with chart defaults.

### spec.worker.enabled

`bool` · optional (explicit presence)

Run the worker fleet. Empty = true when `cache` is set (async
queries, thumbnails, reports execute here).

- default: `true`

### spec.worker.replicas

`int32` · optional (explicit presence)

Worker replicas. Empty = 1.

- default: `1`
- rule: {"int32":{"gte":1}}

### spec.worker.hpa

`KubernetesSupersetHpa`

Horizontal Pod Autoscaler for workers.

### spec.worker.hpa.minReplicas

`int32` · optional (explicit presence)

Lower replica bound. Empty = 1.

- rule: {"int32":{"gte":1}}

### spec.worker.hpa.maxReplicas

`int32`

Upper replica bound. Required.

- rule: {"int32":{"gte":1}}

### spec.worker.hpa.targetCpuUtilizationPercent

`int32` · optional (explicit presence)

Target average CPU utilization percentage (of the container's
CPU requests). Empty = 80 (the chart default).

- rule: {"int32":{"lte":100,"gte":1}}

### spec.worker.resources

`ContainerResources`

CPU/memory per worker container.

### spec.worker.resources.limits

`CpuMemory`

The resource limits for the container.
Specify the maximum amount of CPU and memory that the container can use.

### spec.worker.resources.limits.cpu

`string`

### spec.worker.resources.limits.memory

`string`

### spec.worker.resources.requests

`CpuMemory`

The resource requests for the container.
Specify the minimum amount of CPU and memory that the container is guaranteed.

### spec.worker.resources.requests.cpu

`string`

### spec.worker.resources.requests.memory

`string`

### spec.beat

`KubernetesSupersetBeat`

Celery beat — the scheduler that FIRES alerts, reports and
cache-warmup on their schedules. Off by default; required for
alerts & reports to actually run. Requires `cache` and the
worker.

### spec.beat.enabled

`bool`

Run the beat scheduler (exactly one pod — beat is a singleton
by design). Required for alerts, reports and scheduled cache
warmup to fire.

### spec.flower

`KubernetesSupersetFlower`

Flower — the Celery monitoring UI (port 5555). Off by default.
KNOW THIS: flower ships with NO authentication of its own —
anyone who can reach its Service sees task payloads (including
executed SQL). Keep it off or fence it with a NetworkPolicy.

### spec.flower.enabled

`bool`

Run flower. Unauthenticated — see the spec-level warning.

### spec.websockets

`KubernetesSupersetWebsockets`

The websocket server for GLOBAL_ASYNC_QUERIES in `ws` transport
mode (live query-result push instead of polling). Off by
default. KNOW THIS: the chart's websocket image is a
COMMUNITY-maintained build (`oneacrefund/superset-websocket`) —
pin the tag deliberately; the JWT shared between Superset and
the websocket server is module-generated.

### spec.websockets.enabled

`bool`

Run the websocket server. Enabling it switches Superset's
async-query transport to `ws` automatically (the chart wires
the in-cluster URL); the shared JWT is module-generated into
`<name>-ws-jwt`.

### spec.websockets.image

`KubernetesSupersetWsImage`

Websocket server image. Empty = the chart default
(`oneacrefund/superset-websocket:latest` — a COMMUNITY build
with an UNPINNED tag; pin a digest-verified tag for
production).

### spec.websockets.image.repository

`string` · optional (explicit presence)

Image repository. Empty = "oneacrefund/superset-websocket"
(community-maintained — see the arm's warning).

- default: `oneacrefund/superset-websocket`

### spec.websockets.image.tag

`string` · optional (explicit presence)

Image tag. Empty = "latest" (the chart default — UNPINNED; pin
deliberately for production).

- default: `latest`

### spec.websockets.replicas

`int32` · optional (explicit presence)

Websocket server replicas. Empty = 1.

- default: `1`
- rule: {"int32":{"gte":1}}

### spec.mcp

`KubernetesSupersetMcp`

The Superset MCP server — exposes dashboards/datasets to AI
agents over the Model Context Protocol. Off by default. KNOW
THIS: requires the `fastmcp` python extra, which the official
image does NOT include — provide a custom image (or a
bootstrap_script pip install) or the pods crash-loop.

### spec.mcp.enabled

`bool`

Run the MCP server (port 5008). Requires an image carrying the
`fastmcp` extra — see the spec-level warning.

### spec.init

`KubernetesSupersetInit`

Init behavior: schema migration always runs; admin bootstrap
and example dashboards are configured here.

### spec.init.admin

`KubernetesSupersetAdminUser`

The bootstrap admin user. Empty = created with username "admin"
and a module-generated password.

### spec.init.admin.username

`string` · optional (explicit presence)

Admin username. Empty = "admin".

- default: `admin`

### spec.init.admin.email

`string` · optional (explicit presence)

Admin email (Superset requires one). Empty =
"admin@superset.local".

- default: `admin@superset.local`

### spec.init.admin.passwordSecret

`KubernetesSupersetSecretKeyRef`

Existing Secret holding the admin password. Empty = the module
GENERATES it into `<name>-admin-auth` (key `password`) — the
chart's documented admin/admin default never ships. Either way
the password reaches the init step as an environment variable,
never a rendered literal.

### spec.init.admin.passwordSecret.secretName

`string` · required

Name of the Secret. Same-namespace constraint applies.

- rule: {"required":true}

### spec.init.admin.passwordSecret.secretKey

`string` · required

Key inside the Secret.

- rule: {"required":true}

### spec.init.loadExamples

`bool`

Load Superset's example dashboards and datasets at init
(downloads example data — needs internet from the init pod;
adds minutes to the first install).

### spec.featureFlags

`map<string, bool>`

Superset feature flags
(https://superset.apache.org/docs/configuration/configuring-superset#feature-flags)
— e.g. DASHBOARD_RBAC, ALERT_REPORTS, GLOBAL_ASYNC_QUERIES.

### spec.configOverrides

`map<string, string>`

Raw python snippets appended to the generated
`superset_config.py` (the chart's configOverrides) — name →
snippet. The escape hatch for any config.py setting the typed
fields do not model (OAuth providers, row limits, timeouts).
Read secret values from environment variables
(`os.environ.get(...)` with `extra_env_from_secret`) — the
rendered config lives in a Secret, but literals here also sit
in Helm values.

### spec.extraEnv

`map<string, string>`

Extra environment variables for every Superset container (plain
values — e.g. GUNICORN_TIMEOUT, SERVER_WORKER_AMOUNT). For
secret values use `extra_env_from_secret`.

### spec.extraEnvFromSecret

`map<string, KubernetesSupersetSecretKeyRef>`

Extra environment variables sourced from existing Secrets —
name → Secret reference. The escape hatch for OAuth client
secrets, API keys (MAPBOX_API_KEY) and datasource credentials
consumed by config_overrides snippets.

### spec.extraEnvFromSecret.*.secretName

`string` · required

Name of the Secret. Same-namespace constraint applies.

- rule: {"required":true}

### spec.extraEnvFromSecret.*.secretKey

`string` · required

Key inside the Secret.

- rule: {"required":true}

### spec.bootstrapScript

`string`

Bootstrap script run at every container start (before the
server). The chart's mechanism for `pip install`-ing database
drivers — the official image is the driver-less "lean" build
stage that ships NO metadata-database driver at all (verified
live: the server exits at boot with "No module named 'psycopg2'"
without one). Empty = the module's own default, which installs
the exact psycopg2 pin the app's [postgres] extra declares.
Declaring a script REPLACES that default — include a psycopg2
install (or bake it into a custom image) or the pods crash-loop.
KNOW THIS (verified live): installs must target the app's venv —
the image's plain `pip` is the SYSTEM interpreter's and its
installs stay invisible to the app; use `uv pip install
--python /app/.venv/bin/python <driver>` (uv is the image's own
tool). pip at container start needs internet from the pod and
re-runs on every restart — for production, bake a custom image
instead and set this to a no-op.

### spec.service

`KubernetesSupersetService`

The Superset Service — this kind keeps it ClusterIP (compose
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

### spec.scheduling

`KubernetesSupersetScheduling`

Pod scheduling applied to all Superset pods.

### spec.scheduling.nodeSelector

`map<string, string>`

Node selector.

### spec.scheduling.tolerations

`[]WorkloadToleration`

Tolerations.

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

### spec.helmValues

`string`

Raw Helm values (YAML) merged LAST over everything this spec
renders — the escape hatch for chart surfaces the typed fields
do not model (probes, lifecycle hooks, extra containers,
per-component overrides). The module re-pins security-critical
values after the merge (the bundled postgresql/redis subcharts
stay OFF, the chart's own env Secret stays OFF) — those cannot
be silently re-enabled from here.

## Validation Rules

- `worker.requires_cache`: Celery workers need the cache/broker: set spec.cache (a composed KubernetesValkey or any Redis-protocol store) or remove the worker.
- `beat.requires_worker_and_cache`: Celery beat schedules tasks that WORKERS execute over the cache/broker — enable the worker and set spec.cache, or remove beat.
- `flower.requires_cache`: Flower monitors Celery over the broker — set spec.cache or remove flower.
- `websockets.requires_cache`: The websocket server streams async query events through Redis — set spec.cache or remove websockets.
- `mcp.requires_cache`: The MCP server requires the cache/broker (its pods wait for Redis at startup) — set spec.cache or remove mcp.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesSuperset, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.namespace` | `string` | Namespace Superset runs in. |
| `status.outputs.service` | `string` | Name of the Superset web Service — the handle exposure kinds route to (port 8088). |
| `status.outputs.endpoint` | `string` | In-cluster endpoint, `http://<service>.<namespace>.svc.cluster.local:8088` — the URL browsers and API clients reach Superset at (behind composed exposure). |
| `status.outputs.admin_username` | `string` | The bootstrap admin username (`init.admin.username`, default "admin"). |
| `status.outputs.admin_password_secret` | `KubernetesSecretKey` | The admin credential: the Secret and key holding the bootstrap admin password (module-generated `<name>-admin-auth` unless an existing Secret was declared). |
| `status.outputs.admin_password_secret.name` | `string` | The name of the Kubernetes Secret. |
| `status.outputs.admin_password_secret.key` | `string` | The key within the Kubernetes Secret. |
| `status.outputs.env_secret_name` | `string` | Name of the module-owned environment Secret (`<name>-env`) — the chart's runtime-credential contract, composed at apply time; exported for audit and advanced composition. |
| `status.outputs.secret_key_secret_name` | `string` | Name of the Secret holding the session-signing SECRET_KEY (module-generated `<name>-secret-key` unless an existing Secret was declared) — needed for `superset re-encrypt-secrets` rotation procedures. |
| `status.outputs.port_forward_command` | `string` | Port-forward command for reaching the Superset UI from a workstation when no exposure is composed. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespace` | KubernetesNamespace | `spec.name` |
| `spec.metadataDatabase.host` | KubernetesPostgres | `status.outputs.rw_service` |
| `spec.metadataDatabase.passwordSecret.secretName` | KubernetesPostgres | `status.outputs.password_secret.name` |
| `spec.cache.host` | KubernetesValkey | `status.outputs.service` |
| `spec.cache.passwordSecret.secretName` | KubernetesValkey | `status.outputs.password_secret.name` |

## See Also

- [Overview](../README.md)
