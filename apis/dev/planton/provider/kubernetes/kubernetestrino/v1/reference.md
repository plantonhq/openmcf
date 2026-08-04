# KubernetesTrino

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1`

**KubernetesTrinoSpec** deploys Trino — the distributed SQL query
engine that lets you query data where it lives (data lakes, object
stores, relational databases) and JOIN across sources in one query
(https://trino.io).

WHAT GETS INSTALLED: the official trino Helm chart (trinodb org)
renders a coordinator Deployment (the brain: parses, plans and
schedules queries, serves the REST API and Web UI on port 8080)
and a worker Deployment (the muscle: executes query splits). Every
configuration surface — config.properties, catalogs, access
control — renders into ConfigMaps with checksum annotations, so
config changes roll the pods automatically.

SECURED BY DEFAULT: upstream ships NO authentication — anyone who
can reach the Service can query every catalog. This kind enables
PASSWORD (file) authentication by default with a module-generated
admin user (`<name>-auth` Secret), and configures the
internal-communication shared secret Trino requires once
authentication is on. Password auth normally demands HTTPS; the
module sets `http-server.authentication.allow-insecure-over-http`
because in-cluster traffic rides the ClusterIP Service and TLS
terminates at composed exposure kinds (or the `https` arm here).

CREDENTIALS ARE SECRET-NATIVE: catalog passwords and the internal
shared secret never appear in rendered ConfigMaps — properties
reference environment variables (`${ENV:VAR}`, Trino's own secrets
mechanism) and the variables arrive via Secret references.

QUERY DATA IMMEDIATELY: the chart's default `tpch` and `tpcds`
sample catalogs (in-image data generators) stay available until
disabled — `SELECT count(*) FROM tpch.tiny.nation` works on a
fresh install. Declare postgres/mysql catalogs to query real
databases; a KubernetesPostgres composes naturally.

EXPOSURE: the Service stays ClusterIP; expose it via first-class
kinds over the exported coordinator handle.

## Example

```yaml
# Full-surface shape: composed postgres + mysql catalogs beside the
# in-image samples, a custom memory catalog, fault-tolerant execution
# spooling to an S3-compatible store with env-substituted keys, worker
# graceful shutdown + HPA, access-control/resource-group/session
# documents, JMX metrics with ServiceMonitors, a NetworkPolicy and the
# env escape hatches — the offline plan/preview proof for the widest
# typed rendering.
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesTrino
metadata:
  name: trino-full
spec:
  namespace:
    value: trino-full
  createNamespace: true
  nodeEnvironment: analytics_prod
  logLevel: INFO
  maxQueryMemory: 8GB
  coordinator:
    jvm:
      maxHeapPercent: 70
    maxQueryMemoryPerNode: 2GB
    resources:
      requests:
        cpu: 500m
        memory: 2Gi
      limits:
        memory: 4Gi
  workers:
    jvm:
      maxHeapPercent: 75
    maxQueryMemoryPerNode: 2GB
    resources:
      requests:
        cpu: "1"
        memory: 4Gi
      limits:
        memory: 6Gi
    gracefulShutdown:
      enabled: true
      gracePeriodSeconds: 180
    hpa:
      maxReplicas: 6
      targetCpuUtilizationPercent: 60
      targetMemoryUtilizationPercent: 0
  auth:
    enabled: true
    adminUsername: analyst
  catalogs:
    sampleCatalogsEnabled: true
    postgres:
      - name: warehouse
        host:
          value: warehouse-pg-rw
        port: 5432
        database: warehouse
        username: app
        passwordSecret:
          secretName:
            value: warehouse-pg-app
          secretKey: password
        additionalProperties:
          - postgresql.array-mapping=AS_ARRAY
    mysql:
      - name: orders
        host:
          value: orders-mysql
        username: trino_ro
        passwordSecret:
          secretName:
            value: orders-mysql-trino
          secretKey: password
    custom:
      scratch: |
        connector.name=memory
        memory.max-data-per-node=256MB
  faultTolerantExecution:
    retryPolicy: TASK
    exchangeManager:
      baseDirectories:
        - s3://trino-exchange
      additionalProperties:
        - exchange.s3.region=us-east-1
        - exchange.s3.endpoint=http://exchange-s3.trino-full.svc.cluster.local:8333
        - exchange.s3.aws-access-key=${ENV:EXCHANGE_ACCESS_KEY}
        - exchange.s3.aws-secret-key=${ENV:EXCHANGE_SECRET_KEY}
  accessControlRules: |
    {"catalogs":[{"user":"analyst","catalog":".*","allow":"all"},{"group":"readers","catalog":"warehouse","allow":"read-only"}]}
  resourceGroupsConfig: |
    {"rootGroups":[{"name":"global","softMemoryLimit":"80%","hardConcurrencyLimit":100,"maxQueued":200}],"selectors":[{"group":"global"}]}
  sessionPropertiesConfig: |
    [{"group":"global.*","sessionProperties":{"query_max_execution_time":"4h"}}]
  additionalConfigProperties:
    - http-server.process-forwarded=true
  extraEnvFromSecret:
    EXCHANGE_ACCESS_KEY:
      secretName: exchange-s3-secret
      secretKey: admin_access_key_id
    EXCHANGE_SECRET_KEY:
      secretName: exchange-s3-secret
      secretKey: admin_secret_access_key
  extraEnv:
    JAVA_TOOL_OPTIONS: -XX:+UseStringDeduplication
  metrics:
    enabled: true
    serviceMonitorEnabled: true
  networkPolicyEnabled: true
  service:
    type: ClusterIP
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.namespace` | `string \| valueFrom` | yes |  | KubernetesNamespace (`spec.name`) |
| `spec.createNamespace` | `bool` |  |  |  |
| `spec.image` | `KubernetesTrinoImage` |  |  |  |
| `spec.image.registry` | `string` |  |  |  |
| `spec.image.repository` | `string` |  | `trinodb/trino` |  |
| `spec.image.tag` | `string` |  | `480` |  |
| `spec.imagePullSecrets` | `[]string` |  |  |  |
| `spec.nodeEnvironment` | `string` |  | `production` |  |
| `spec.logLevel` | `string` |  | `INFO` |  |
| `spec.coordinator` | `KubernetesTrinoCoordinator` |  |  |  |
| `spec.coordinator.jvm` | `KubernetesTrinoJvm` |  |  |  |
| `spec.coordinator.jvm.maxHeapSize` | `string` |  |  |  |
| `spec.coordinator.jvm.maxHeapPercent` | `int32` |  |  |  |
| `spec.coordinator.maxQueryMemoryPerNode` | `string` |  | `1GB` |  |
| `spec.coordinator.heapHeadroomPerNode` | `string` |  |  |  |
| `spec.coordinator.includeInScheduling` | `bool` |  |  |  |
| `spec.coordinator.resources` | `ContainerResources` |  |  |  |
| `spec.coordinator.resources.limits` | `CpuMemory` |  |  |  |
| `spec.coordinator.resources.limits.cpu` | `string` |  |  |  |
| `spec.coordinator.resources.limits.memory` | `string` |  |  |  |
| `spec.coordinator.resources.requests` | `CpuMemory` |  |  |  |
| `spec.coordinator.resources.requests.cpu` | `string` |  |  |  |
| `spec.coordinator.resources.requests.memory` | `string` |  |  |  |
| `spec.coordinator.scheduling` | `KubernetesTrinoScheduling` |  |  |  |
| `spec.coordinator.scheduling.nodeSelector` | `map<string, string>` |  |  |  |
| `spec.coordinator.scheduling.tolerations` | `[]WorkloadToleration` |  |  |  |
| `spec.coordinator.scheduling.tolerations[].key` | `string` |  |  |  |
| `spec.coordinator.scheduling.tolerations[].operator` | `string` |  |  |  |
| `spec.coordinator.scheduling.tolerations[].value` | `string` |  |  |  |
| `spec.coordinator.scheduling.tolerations[].effect` | `string` |  |  |  |
| `spec.coordinator.scheduling.tolerations[].tolerationSeconds` | `int64` |  |  |  |
| `spec.workers` | `KubernetesTrinoWorkers` |  |  |  |
| `spec.workers.replicas` | `int32` |  | `2` |  |
| `spec.workers.jvm` | `KubernetesTrinoJvm` |  |  |  |
| `spec.workers.jvm.maxHeapSize` | `string` |  |  |  |
| `spec.workers.jvm.maxHeapPercent` | `int32` |  |  |  |
| `spec.workers.maxQueryMemoryPerNode` | `string` |  | `1GB` |  |
| `spec.workers.heapHeadroomPerNode` | `string` |  |  |  |
| `spec.workers.resources` | `ContainerResources` |  |  |  |
| `spec.workers.resources.limits` | `CpuMemory` |  |  |  |
| `spec.workers.resources.limits.cpu` | `string` |  |  |  |
| `spec.workers.resources.limits.memory` | `string` |  |  |  |
| `spec.workers.resources.requests` | `CpuMemory` |  |  |  |
| `spec.workers.resources.requests.cpu` | `string` |  |  |  |
| `spec.workers.resources.requests.memory` | `string` |  |  |  |
| `spec.workers.scheduling` | `KubernetesTrinoScheduling` |  |  |  |
| `spec.workers.scheduling.nodeSelector` | `map<string, string>` |  |  |  |
| `spec.workers.scheduling.tolerations` | `[]WorkloadToleration` |  |  |  |
| `spec.workers.scheduling.tolerations[].key` | `string` |  |  |  |
| `spec.workers.scheduling.tolerations[].operator` | `string` |  |  |  |
| `spec.workers.scheduling.tolerations[].value` | `string` |  |  |  |
| `spec.workers.scheduling.tolerations[].effect` | `string` |  |  |  |
| `spec.workers.scheduling.tolerations[].tolerationSeconds` | `int64` |  |  |  |
| `spec.workers.gracefulShutdown` | `KubernetesTrinoGracefulShutdown` |  |  |  |
| `spec.workers.gracefulShutdown.enabled` | `bool` |  |  |  |
| `spec.workers.gracefulShutdown.gracePeriodSeconds` | `int32` |  | `120` |  |
| `spec.workers.hpa` | `KubernetesTrinoWorkerHpa` |  |  |  |
| `spec.workers.hpa.maxReplicas` | `int32` |  |  |  |
| `spec.workers.hpa.targetCpuUtilizationPercent` | `int32` |  | `50` |  |
| `spec.workers.hpa.targetMemoryUtilizationPercent` | `int32` |  | `80` |  |
| `spec.workers.keda` | `KubernetesTrinoWorkerKeda` |  |  |  |
| `spec.workers.keda.minReplicas` | `int32` |  |  |  |
| `spec.workers.keda.maxReplicas` | `int32` |  |  |  |
| `spec.workers.keda.pollingIntervalSeconds` | `int32` |  |  |  |
| `spec.workers.keda.cooldownPeriodSeconds` | `int32` |  |  |  |
| `spec.workers.keda.triggers` | `string` | yes |  |  |
| `spec.maxQueryMemory` | `string` |  | `4GB` |  |
| `spec.auth` | `KubernetesTrinoAuth` |  |  |  |
| `spec.auth.enabled` | `bool` |  | `true` |  |
| `spec.auth.adminUsername` | `string` |  | `trino` |  |
| `spec.auth.existingPasswordDbSecret` | `KubernetesTrinoExistingSecret` |  |  |  |
| `spec.auth.existingPasswordDbSecret.secretName` | `string` | yes |  |  |
| `spec.auth.groupsSecret` | `KubernetesTrinoExistingSecret` |  |  |  |
| `spec.auth.groupsSecret.secretName` | `string` | yes |  |  |
| `spec.https` | `KubernetesTrinoHttps` |  |  |  |
| `spec.https.enabled` | `bool` |  |  |  |
| `spec.https.port` | `int32` |  | `8443` |  |
| `spec.https.keystoreSecret` | `KubernetesTrinoKeystoreSecret` |  |  |  |
| `spec.https.keystoreSecret.secretName` | `string` | yes |  |  |
| `spec.https.keystoreSecret.secretKey` | `string` |  | `keystore.jks` |  |
| `spec.catalogs` | `KubernetesTrinoCatalogs` |  |  |  |
| `spec.catalogs.sampleCatalogsEnabled` | `bool` |  | `true` |  |
| `spec.catalogs.postgres` | `[]KubernetesTrinoPostgresCatalog` |  |  |  |
| `spec.catalogs.postgres[].name` | `string` | yes |  |  |
| `spec.catalogs.postgres[].host` | `string \| valueFrom` | yes |  | KubernetesPostgres (`status.outputs.rw_service`) |
| `spec.catalogs.postgres[].port` | `int32` |  | `5432` |  |
| `spec.catalogs.postgres[].database` | `string` | yes |  |  |
| `spec.catalogs.postgres[].username` | `string` |  | `app` |  |
| `spec.catalogs.postgres[].passwordSecret` | `KubernetesTrinoPostgresPasswordSecret` | yes |  |  |
| `spec.catalogs.postgres[].passwordSecret.secretName` | `string \| valueFrom` | yes |  | KubernetesPostgres (`status.outputs.password_secret.name`) |
| `spec.catalogs.postgres[].passwordSecret.secretKey` | `string` |  | `password` |  |
| `spec.catalogs.postgres[].additionalProperties` | `[]string` |  |  |  |
| `spec.catalogs.mysql` | `[]KubernetesTrinoMysqlCatalog` |  |  |  |
| `spec.catalogs.mysql[].name` | `string` | yes |  |  |
| `spec.catalogs.mysql[].host` | `string \| valueFrom` | yes |  | KubernetesMysql (`status.outputs.primary_service`) |
| `spec.catalogs.mysql[].port` | `int32` |  | `3306` |  |
| `spec.catalogs.mysql[].username` | `string` |  | `root` |  |
| `spec.catalogs.mysql[].passwordSecret` | `KubernetesTrinoMysqlPasswordSecret` | yes |  |  |
| `spec.catalogs.mysql[].passwordSecret.secretName` | `string \| valueFrom` | yes |  | KubernetesMysql (`status.outputs.root_password_secret.name`) |
| `spec.catalogs.mysql[].passwordSecret.secretKey` | `string` |  | `password` |  |
| `spec.catalogs.mysql[].additionalProperties` | `[]string` |  |  |  |
| `spec.catalogs.custom` | `map<string, string>` |  |  |  |
| `spec.faultTolerantExecution` | `KubernetesTrinoFaultTolerantExecution` |  |  |  |
| `spec.faultTolerantExecution.retryPolicy` | `string` | yes |  |  |
| `spec.faultTolerantExecution.exchangeManager` | `KubernetesTrinoExchangeManager` | yes |  |  |
| `spec.faultTolerantExecution.exchangeManager.baseDirectories` | `[]string` | yes |  |  |
| `spec.faultTolerantExecution.exchangeManager.additionalProperties` | `[]string` |  |  |  |
| `spec.accessControlRules` | `string` |  |  |  |
| `spec.resourceGroupsConfig` | `string` |  |  |  |
| `spec.sessionPropertiesConfig` | `string` |  |  |  |
| `spec.eventListenerProperties` | `[]string` |  |  |  |
| `spec.additionalConfigProperties` | `[]string` |  |  |  |
| `spec.extraEnv` | `map<string, string>` |  |  |  |
| `spec.extraEnvFromSecret` | `map<string, KubernetesTrinoSecretKeyRef>` |  |  |  |
| `spec.extraEnvFromSecret.*.secretName` | `string` | yes |  |  |
| `spec.extraEnvFromSecret.*.secretKey` | `string` | yes |  |  |
| `spec.metrics` | `KubernetesTrinoMetrics` |  |  |  |
| `spec.metrics.enabled` | `bool` |  |  |  |
| `spec.metrics.serviceMonitorEnabled` | `bool` |  |  |  |
| `spec.metrics.exporterImage` | `string` |  |  |  |
| `spec.networkPolicyEnabled` | `bool` |  |  |  |
| `spec.service` | `KubernetesTrinoService` |  |  |  |
| `spec.service.type` | `string` |  | `ClusterIP` |  |
| `spec.service.annotations` | `map<string, string>` |  |  |  |
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
THIS: catalog credential Secrets are read by the pods at
runtime, and a Secret can only be referenced from the workload's
OWN namespace — co-locate Trino with the database Secrets its
catalogs reference.

### spec.image

`KubernetesTrinoImage`

Container image. Empty = the official `trinodb/trino` image at
the Trino release this kind is built against. The chart takes
the SPLIT form (registry + repository) — a mirror override sets
`registry` (e.g. "my-registry.example.com") while `repository`
keeps the path.

### spec.image.registry

`string`

Image registry host. Empty = Docker Hub.

### spec.image.repository

`string` · optional (explicit presence)

Image repository (path without the registry host). Empty =
"trinodb/trino".

- default: `trinodb/trino`

### spec.image.tag

`string` · optional (explicit presence)

Image tag. Empty = "480" (the Trino release this kind is built
against).

- default: `480`

### spec.imagePullSecrets

`[]string`

Names of image-pull Secrets in the same namespace, for pulling
from private mirrors.

### spec.nodeEnvironment

`string` · optional (explicit presence)

Trino node environment name — shows in the Web UI and groups
nodes into one cluster identity. Empty = "production". Must be
lowercase alphanumeric/underscore (Trino rejects anything else
at startup).

- default: `production`
- rule: {"string":{"pattern":"^[a-z0-9][a-z0-9_]*$"}}

### spec.logLevel

`string` · optional (explicit presence)

Log level for the `io.trino` logger. Empty = "INFO".

- default: `INFO`
- rule: {"string":{"in":["DEBUG","INFO","WARN","ERROR"]}}

### spec.coordinator

`KubernetesTrinoCoordinator`

The coordinator: query planning, scheduling, the REST API and
Web UI.

### spec.coordinator.jvm

`KubernetesTrinoJvm`

JVM heap sizing for the coordinator container.

### spec.coordinator.jvm.maxHeapSize

`string`

Fixed max heap (e.g. "8G" — the chart default when nothing is
set). Ignored when `max_heap_percent` is set.

### spec.coordinator.jvm.maxHeapPercent

`int32` · optional (explicit presence)

Max heap as a percentage of the container memory limit
(`-XX:MaxRAMPercentage`). Requires a memory LIMIT on the
container. When set, the module clears the fixed `-Xmx`.

- rule: {"int32":{"lte":95,"gte":10}}

### spec.coordinator.maxQueryMemoryPerNode

`string` · optional (explicit presence)

Per-node memory ceiling for a single query on the coordinator
(`query.max-memory-per-node`). Empty = "1GB".

- default: `1GB`
- rule: {"string":{"pattern":"^[0-9]+(\\.[0-9]+)?(B|kB|MB|GB|TB)$"}}

### spec.coordinator.heapHeadroomPerNode

`string`

Heap memory reserved for non-query work
(`memory.heap-headroom-per-node`). Empty = the Trino default
(30% of heap).

### spec.coordinator.includeInScheduling

`bool`

Also schedule query work ON the coordinator
(`node-scheduler.include-coordinator`). With `workers.replicas`
0 this gives a true single-node Trino — fine for dev/small
loads; on real clusters keep it false so coordinator resources
stay dedicated to planning and scheduling.

### spec.coordinator.resources

`ContainerResources`

CPU/memory for the coordinator container. Empty = no requests
(fine for dev; size for production). KNOW THIS: the chart's JVM
default is an 8G max heap — when setting container limits,
either size them above the heap or switch to
`jvm.max_heap_percent` so the heap follows the limit.

### spec.coordinator.resources.limits

`CpuMemory`

The resource limits for the container.
Specify the maximum amount of CPU and memory that the container can use.

### spec.coordinator.resources.limits.cpu

`string`

### spec.coordinator.resources.limits.memory

`string`

### spec.coordinator.resources.requests

`CpuMemory`

The resource requests for the container.
Specify the minimum amount of CPU and memory that the container is guaranteed.

### spec.coordinator.resources.requests.cpu

`string`

### spec.coordinator.resources.requests.memory

`string`

### spec.coordinator.scheduling

`KubernetesTrinoScheduling`

Pod scheduling for the coordinator.

### spec.coordinator.scheduling.nodeSelector

`map<string, string>`

Node selector.

### spec.coordinator.scheduling.tolerations

`[]WorkloadToleration`

Tolerations.

### spec.coordinator.scheduling.tolerations[].key

`string`

Taint key to tolerate. Empty key with operator "Exists" tolerates every taint.

### spec.coordinator.scheduling.tolerations[].operator

`string`

How key/value match: "Equal" (default — value must match too) or "Exists"
(key presence alone matches).

- rule: Toleration operator must be either "Equal" or "Exists"

### spec.coordinator.scheduling.tolerations[].value

`string`

Taint value to match when operator is "Equal".

### spec.coordinator.scheduling.tolerations[].effect

`string`

Which taint effect is tolerated: "NoSchedule", "PreferNoSchedule", or
"NoExecute". Empty tolerates all effects for the key.

- rule: Toleration effect must be one of "NoSchedule", "PreferNoSchedule", or "NoExecute"

### spec.coordinator.scheduling.tolerations[].tolerationSeconds

`int64` · optional (explicit presence)

For "NoExecute" taints only: how many seconds already-running pods stay bound
after the taint appears. Unset means tolerate forever.

### spec.workers

`KubernetesTrinoWorkers`

The worker fleet: query execution. Empty = 2 workers with chart
defaults.

### spec.workers.replicas

`int32` · optional (explicit presence)

Worker count. Empty = 2. Set 0 for the single-node shape
(pair with `coordinator.include_in_scheduling`) — otherwise
queries have no execution capacity. Ignored when an autoscaling
arm is set.

- default: `2`
- rule: {"int32":{"gte":0}}

### spec.workers.jvm

`KubernetesTrinoJvm`

JVM heap sizing for worker containers.

### spec.workers.jvm.maxHeapSize

`string`

Fixed max heap (e.g. "8G" — the chart default when nothing is
set). Ignored when `max_heap_percent` is set.

### spec.workers.jvm.maxHeapPercent

`int32` · optional (explicit presence)

Max heap as a percentage of the container memory limit
(`-XX:MaxRAMPercentage`). Requires a memory LIMIT on the
container. When set, the module clears the fixed `-Xmx`.

- rule: {"int32":{"lte":95,"gte":10}}

### spec.workers.maxQueryMemoryPerNode

`string` · optional (explicit presence)

Per-node memory ceiling for a single query on each worker
(`query.max-memory-per-node`). Empty = "1GB".

- default: `1GB`
- rule: {"string":{"pattern":"^[0-9]+(\\.[0-9]+)?(B|kB|MB|GB|TB)$"}}

### spec.workers.heapHeadroomPerNode

`string`

Heap memory reserved for non-query work. Empty = the Trino
default.

### spec.workers.resources

`ContainerResources`

CPU/memory per worker container. Empty = no requests. The same
heap-vs-limit rule as the coordinator applies.

### spec.workers.resources.limits

`CpuMemory`

The resource limits for the container.
Specify the maximum amount of CPU and memory that the container can use.

### spec.workers.resources.limits.cpu

`string`

### spec.workers.resources.limits.memory

`string`

### spec.workers.resources.requests

`CpuMemory`

The resource requests for the container.
Specify the minimum amount of CPU and memory that the container is guaranteed.

### spec.workers.resources.requests.cpu

`string`

### spec.workers.resources.requests.memory

`string`

### spec.workers.scheduling

`KubernetesTrinoScheduling`

Pod scheduling for workers.

### spec.workers.scheduling.nodeSelector

`map<string, string>`

Node selector.

### spec.workers.scheduling.tolerations

`[]WorkloadToleration`

Tolerations.

### spec.workers.scheduling.tolerations[].key

`string`

Taint key to tolerate. Empty key with operator "Exists" tolerates every taint.

### spec.workers.scheduling.tolerations[].operator

`string`

How key/value match: "Equal" (default — value must match too) or "Exists"
(key presence alone matches).

- rule: Toleration operator must be either "Equal" or "Exists"

### spec.workers.scheduling.tolerations[].value

`string`

Taint value to match when operator is "Equal".

### spec.workers.scheduling.tolerations[].effect

`string`

Which taint effect is tolerated: "NoSchedule", "PreferNoSchedule", or
"NoExecute". Empty tolerates all effects for the key.

- rule: Toleration effect must be one of "NoSchedule", "PreferNoSchedule", or "NoExecute"

### spec.workers.scheduling.tolerations[].tolerationSeconds

`int64` · optional (explicit presence)

For "NoExecute" taints only: how many seconds already-running pods stay bound
after the taint appears. Unset means tolerate forever.

### spec.workers.gracefulShutdown

`KubernetesTrinoGracefulShutdown`

Drain workers before termination (`shutdown.grace-period`) so
running queries finish instead of failing. Strongly recommended
with autoscaling.

### spec.workers.gracefulShutdown.enabled

`bool`

Enable coordinated drain: a terminating worker stops accepting
new tasks and finishes running ones within the grace period.

### spec.workers.gracefulShutdown.gracePeriodSeconds

`int32` · optional (explicit presence)

Drain window in seconds (`shutdown.grace-period`). Empty = 120.
The chart requires the pod termination budget to be at least
TWICE this value — the module sets
`terminationGracePeriodSeconds` to 2× automatically.

- default: `120`
- rule: {"int32":{"gte":1}}

### spec.workers.hpa

`KubernetesTrinoWorkerHpa`

Horizontal Pod Autoscaler on CPU/memory utilization (needs a
metrics server and worker resource REQUESTS to compute
utilization against).

### spec.workers.hpa.maxReplicas

`int32`

Upper replica bound. Required.

- rule: {"int32":{"gte":1}}

### spec.workers.hpa.targetCpuUtilizationPercent

`int32` · optional (explicit presence)

Target average CPU utilization percentage (of worker CPU
requests). Empty = 50 (the chart default). Set 0 to disable the
CPU metric.

- default: `50`
- rule: {"int32":{"lte":100,"gte":0}}

### spec.workers.hpa.targetMemoryUtilizationPercent

`int32` · optional (explicit presence)

Target average memory utilization percentage (of worker memory
requests). Empty = 80 (the chart default). Set 0 to disable the
memory metric.

- default: `80`
- rule: {"int32":{"lte":100,"gte":0}}

### spec.workers.keda

`KubernetesTrinoWorkerKeda`

KEDA event-driven autoscaling (scale on Prometheus query
metrics, scale to zero). Requires the KEDA operator on the
cluster — a KubernetesKeda composes naturally.

### spec.workers.keda.minReplicas

`int32` · optional (explicit presence)

Lower replica bound. Empty = 0 (scale to zero between queries —
pair with graceful shutdown).

- rule: {"int32":{"gte":0}}

### spec.workers.keda.maxReplicas

`int32`

Upper replica bound. Required.

- rule: {"int32":{"gte":1}}

### spec.workers.keda.pollingIntervalSeconds

`int32` · optional (explicit presence)

Seconds between trigger evaluations. Empty = 30.

- rule: {"int32":{"gte":1}}

### spec.workers.keda.cooldownPeriodSeconds

`int32` · optional (explicit presence)

Seconds to wait after the last active trigger before scaling
down. Empty = 300.

- rule: {"int32":{"gte":0}}

### spec.workers.keda.triggers

`string` · required

KEDA trigger list as raw YAML (the `triggers:` array content —
e.g. a prometheus trigger on Trino's own
`trino_execution_ClusterSizeMonitor_RequiredWorkers` metric).
Required — KEDA without triggers scales nothing.

- rule: {"required":true}

### spec.maxQueryMemory

`string` · optional (explicit presence)

Cluster-wide query memory ceiling (`query.max-memory` — the
total distributed memory one query may consume across all
nodes). Empty = "4GB". Size together with the per-node limits
under coordinator/workers.

- default: `4GB`
- rule: {"string":{"pattern":"^[0-9]+(\\.[0-9]+)?(B|kB|MB|GB|TB)$"}}

### spec.auth

`KubernetesTrinoAuth`

Authentication. Empty = ENABLED with a module-generated admin
user — the open, anyone-can-query server never ships.

### spec.auth.enabled

`bool` · optional (explicit presence)

Require authentication. Empty = true — the open server never
ships by default. Disabling means anyone who can reach the
Service can query every catalog and impersonate any user.

- default: `true`

### spec.auth.adminUsername

`string` · optional (explicit presence)

The bootstrap admin username. Empty = "trino". The module
generates this user's password into the `<name>-auth` Secret
(key `password`) and writes the matching bcrypt entry into the
server's password file — exported for clients and verifiers.
Ignored when `existing_password_db_secret` is set.

- default: `trino`
- rule: {"string":{"pattern":"^[a-z0-9][a-z0-9._-]*$"}}

### spec.auth.existingPasswordDbSecret

`KubernetesTrinoExistingSecret`

Bring-your-own password file: an existing Secret whose
`password.db` key holds the full htpasswd-format file
(bcrypt entries, one `user:hash` per line). Replaces the
module-generated admin file entirely — manage every user
yourself.

### spec.auth.existingPasswordDbSecret.secretName

`string` · required

Name of the Secret. Same-namespace constraint applies.

- rule: {"required":true}

### spec.auth.groupsSecret

`KubernetesTrinoExistingSecret`

Group mappings: an existing Secret whose `group.db` key holds
the group file (`group_name:user_1,user_2` per line) — pairs
with `access_control_rules` group selectors.

### spec.auth.groupsSecret.secretName

`string` · required

Name of the Secret. Same-namespace constraint applies.

- rule: {"required":true}

### spec.https

`KubernetesTrinoHttps`

HTTPS directly on the Trino pods (a JKS keystore from an
existing Secret). Most deployments leave this empty and
terminate TLS at composed exposure kinds instead.

### spec.https.enabled

`bool`

Serve HTTPS from the Trino containers.

### spec.https.port

`int32` · optional (explicit presence)

HTTPS port. Empty = 8443.

- default: `8443`
- rule: {"int32":{"lte":65535,"gte":1}}

### spec.https.keystoreSecret

`KubernetesTrinoKeystoreSecret`

Existing Secret holding the JKS keystore. Mounted onto the pods
and wired as `http-server.https.keystore.path`.

### spec.https.keystoreSecret.secretName

`string` · required

Name of the Secret. Same-namespace constraint applies.

- rule: {"required":true}

### spec.https.keystoreSecret.secretKey

`string` · optional (explicit presence)

Key inside the Secret holding the keystore. Empty =
"keystore.jks".

- default: `keystore.jks`

### spec.catalogs

`KubernetesTrinoCatalogs`

The catalogs — each one connects Trino to a data source and
becomes a queryable prefix (`SELECT ... FROM <catalog>.<schema>.<table>`).
Empty = the chart's tpch/tpcds sample catalogs only.

- rule: Catalog names must be unique across postgres, mysql and custom catalogs — each becomes one <name>.properties file.
- rule: Catalog name 'system' is reserved (Trino's internal catalog), and 'tpch'/'tpcds' collide with the sample catalogs while sample_catalogs_enabled is on — pick another name or disable the samples.

### spec.catalogs.sampleCatalogsEnabled

`bool` · optional (explicit presence)

Keep the chart's built-in `tpch` and `tpcds` sample catalogs
(in-image synthetic data generators — zero dependencies, great
for trying queries). Empty = true; disable for a
production-only catalog list.

- default: `true`

### spec.catalogs.postgres

`[]KubernetesTrinoPostgresCatalog`

PostgreSQL catalogs — each connects one PostgreSQL database. A
KubernetesPostgres composes naturally (host and credential FK
onto its outputs).

### spec.catalogs.postgres[].name

`string` · required

Catalog name — the query prefix
(`SELECT ... FROM <name>.<schema>.<table>`). Lowercase
alphanumeric/underscore.

- rule: {"required":true,"string":{"pattern":"^[a-z][a-z0-9_]*$"}}

### spec.catalogs.postgres[].host

`string | valueFrom` · required

Database server host. Defaults compose a KubernetesPostgres
resource's read-write Service.

- references: KubernetesPostgres (`status.outputs.rw_service`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesPostgres, name: <that resource's name>, fieldPath: status.outputs.rw_service}} -- a bare string does not parse

### spec.catalogs.postgres[].port

`int32` · optional (explicit presence)

Database server port. Empty = 5432.

- default: `5432`
- rule: {"int32":{"lte":65535,"gte":1}}

### spec.catalogs.postgres[].database

`string` · required

Database to connect to (Trino sees its schemas/tables).

- rule: {"required":true}

### spec.catalogs.postgres[].username

`string` · optional (explicit presence)

Database user Trino connects as — its grants bound what queries
through this catalog can read/write. Empty = "app" (the
KubernetesPostgres application-user convention).

- default: `app`

### spec.catalogs.postgres[].passwordSecret

`KubernetesTrinoPostgresPasswordSecret` · required

The Secret holding the user's password (delivered as an
environment variable the catalog references — never rendered).
Defaults compose a KubernetesPostgres resource's
application-user Secret. Same-namespace constraint applies.

- rule: {"required":true}

### spec.catalogs.postgres[].passwordSecret.secretName

`string | valueFrom` · required

Name of the Secret. Defaults compose a KubernetesPostgres
resource's application-user Secret (`<cluster>-app`).
Same-namespace constraint applies.

- references: KubernetesPostgres (`status.outputs.password_secret.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesPostgres, name: <that resource's name>, fieldPath: status.outputs.password_secret.name}} -- a bare string does not parse

### spec.catalogs.postgres[].passwordSecret.secretKey

`string` · optional (explicit presence)

Key inside the Secret holding the password. Empty = "password"
(the KubernetesPostgres application-Secret convention).

- default: `password`

### spec.catalogs.postgres[].additionalProperties

`[]string`

Additional catalog properties appended to this catalog's file
(connector tuning — e.g. `postgresql.array-mapping`). Never put
secret literals here.

### spec.catalogs.mysql

`[]KubernetesTrinoMysqlCatalog`

MySQL catalogs — each connects one MySQL server. A
KubernetesMysql composes naturally.

### spec.catalogs.mysql[].name

`string` · required

Catalog name — the query prefix. Lowercase
alphanumeric/underscore.

- rule: {"required":true,"string":{"pattern":"^[a-z][a-z0-9_]*$"}}

### spec.catalogs.mysql[].host

`string | valueFrom` · required

Database server host. Defaults compose a KubernetesMysql
resource's client Service.

- references: KubernetesMysql (`status.outputs.primary_service`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesMysql, name: <that resource's name>, fieldPath: status.outputs.primary_service}} -- a bare string does not parse

### spec.catalogs.mysql[].port

`int32` · optional (explicit presence)

Database server port. Empty = 3306.

- default: `3306`
- rule: {"int32":{"lte":65535,"gte":1}}

### spec.catalogs.mysql[].username

`string` · optional (explicit presence)

Database user Trino connects as. Empty = "root".

- default: `root`

### spec.catalogs.mysql[].passwordSecret

`KubernetesTrinoMysqlPasswordSecret` · required

The Secret holding the user's password. Defaults compose a
KubernetesMysql resource's root credential Secret; for a
dedicated read-only Trino user (recommended), point at that
user's Secret. Same-namespace constraint applies.

- rule: {"required":true}

### spec.catalogs.mysql[].passwordSecret.secretName

`string | valueFrom` · required

Name of the Secret. Defaults compose a KubernetesMysql
resource's root credential Secret. Same-namespace constraint
applies.

- references: KubernetesMysql (`status.outputs.root_password_secret.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesMysql, name: <that resource's name>, fieldPath: status.outputs.root_password_secret.name}} -- a bare string does not parse

### spec.catalogs.mysql[].passwordSecret.secretKey

`string` · optional (explicit presence)

Key inside the Secret holding the password. Empty = "password".

- default: `password`

### spec.catalogs.mysql[].additionalProperties

`[]string`

Additional catalog properties appended to this catalog's file.
Never put secret literals here.

### spec.catalogs.custom

`map<string, string>`

Any other connector: catalog name → raw catalog properties
(`connector.name=...` plus its configuration). Reference
environment variables for secret values (`${ENV:VAR}` paired
with `extra_env_from_secret`) — properties render into a
ConfigMap.

### spec.faultTolerantExecution

`KubernetesTrinoFaultTolerantExecution`

Fault-tolerant execution: queries survive worker failures by
spooling exchange data to durable storage and retrying failed
tasks. Off by default (the standard pipeline execution).

### spec.faultTolerantExecution.retryPolicy

`string` · required

Retry granularity: "TASK" (retry failed tasks — the batch/ETL
mode, pairs with worker autoscaling and spot capacity) or
"QUERY" (retry whole failed queries — the interactive mode).

- rule: {"required":true,"string":{"in":["QUERY","TASK"]}}

### spec.faultTolerantExecution.exchangeManager

`KubernetesTrinoExchangeManager` · required

Where exchange data spools. Required — fault tolerance without
durable spooling loses the data it needs to retry.

- rule: {"required":true}

### spec.faultTolerantExecution.exchangeManager.baseDirectories

`[]string` · required

Spooling destinations: local paths (single-node/dev) or object
store URIs (`s3://bucket` — production; a KubernetesSeaweedFs
bucket works via its S3 endpoint). At least one.

- rule: {"repeated":{"minItems":"1"}}

### spec.faultTolerantExecution.exchangeManager.additionalProperties

`[]string`

Additional exchange-manager properties (S3 endpoint, region —
e.g. `exchange.s3.endpoint=...`). Reference environment
variables for keys (`exchange.s3.aws-access-key=${ENV:VAR}`
with `extra_env_from_secret`) — never literals.

### spec.accessControlRules

`string`

System access control rules (the `file` provider) as the JSON
rules document — who may access which catalogs/schemas/tables.
Empty = no system access control (authenticated users see
everything). See
https://trino.io/docs/current/security/file-system-access-control.html
for the document shape.

### spec.resourceGroupsConfig

`string`

Resource-group configuration (the `file` manager) as the JSON
document — concurrency/memory/queue limits per user group.
Empty = no resource groups.

### spec.sessionPropertiesConfig

`string`

Session-property overrides per resource group / user, as the
JSON document for the `file` session property manager. Empty =
none.

### spec.eventListenerProperties

`[]string`

Event-listener properties (`event-listener.name=...` plus its
configuration) — query audit/logging integrations. Empty =
none. Reference environment variables for any secret values
(`${ENV:VAR}` with `extra_env_from_secret`), never literals.

### spec.additionalConfigProperties

`[]string`

Additional raw `config.properties` lines applied to BOTH
coordinator and workers — the escape hatch for tuning
properties the typed fields do not model. These render into a
ConfigMap: reference environment variables for secret values
(`${ENV:VAR}` with `extra_env_from_secret`), never literals.

### spec.extraEnv

`map<string, string>`

Extra environment variables injected into every Trino container
(plain values). For secret values use `extra_env_from_secret`.

### spec.extraEnvFromSecret

`map<string, KubernetesTrinoSecretKeyRef>`

Extra environment variables sourced from existing Secrets —
name → Secret reference. Pair with `${ENV:VAR}` references in
catalog properties or config properties to keep credentials out
of rendered ConfigMaps (Trino's own secrets mechanism).

### spec.extraEnvFromSecret.*.secretName

`string` · required

Name of the Secret. Same-namespace constraint applies.

- rule: {"required":true}

### spec.extraEnvFromSecret.*.secretKey

`string` · required

Key inside the Secret.

- rule: {"required":true}

### spec.metrics

`KubernetesTrinoMetrics`

JMX metrics and Prometheus scraping.

- rule: ServiceMonitors scrape the JMX exporter — enable metrics.enabled too, or remove service_monitor_enabled.

### spec.metrics.enabled

`bool`

Enable JMX metrics and the bundled Prometheus JMX exporter
sidecar (HTTP metrics on every pod).

### spec.metrics.serviceMonitorEnabled

`bool`

Render ServiceMonitors (coordinator + worker) for operator-based
scraping. Requires the Prometheus operator CRDs on the cluster —
a KubernetesKubePrometheusStack composes naturally.

### spec.metrics.exporterImage

`string`

JMX-exporter sidecar image override. Empty = the chart default
(`bitnamilegacy/jmx-exporter` — a FROZEN legacy mirror that
receives no updates since the Bitnami retirement; override for
production mirrors).

### spec.networkPolicyEnabled

`bool`

Restrict pod traffic with a NetworkPolicy (Trino pods may talk
to each other; add extra ingress rules for clients/scrapers).
Requires a CNI that enforces NetworkPolicy.

### spec.service

`KubernetesTrinoService`

The coordinator Service — this kind keeps it ClusterIP (compose
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

### spec.helmValues

`string`

Raw Helm values (YAML) merged LAST over everything this spec
renders — the escape hatch for chart surfaces the typed fields
do not model (probes, lifecycle hooks, sidecars, configMounts).
The module re-pins security-critical values after the merge
(authentication wiring, the shared secret) — those cannot be
silently disabled from here.

## Validation Rules

- `https.keystore_required`: HTTPS on the Trino pods needs the keystore: set https.keystore_secret (a Secret holding the JKS keystore) when https.enabled is true.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesTrino, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.namespace` | `string` | Namespace Trino runs in. |
| `status.outputs.coordinator_service` | `string` | Name of the coordinator Service (`<name>-trino-coordinator` at default naming) — the handle exposure kinds route to. |
| `status.outputs.coordinator_endpoint` | `string` | In-cluster coordinator endpoint, `http://<coordinator_service>.<namespace>.svc.cluster.local:8080` — the URL SQL clients and BI tools connect to (with auth on, pair it with the admin credential below). |
| `status.outputs.admin_username` | `string` | The bootstrap admin username (`auth.admin_username`, default "trino"). Empty when auth is disabled or a bring-your-own password file is declared. |
| `status.outputs.admin_password_secret` | `KubernetesSecretKey` | The admin credential: the Secret and key holding the bootstrap admin password (module-generated `<name>-auth`, key `password`). Empty when auth is disabled or bring-your-own. |
| `status.outputs.admin_password_secret.name` | `string` | The name of the Kubernetes Secret. |
| `status.outputs.admin_password_secret.key` | `string` | The key within the Kubernetes Secret. |
| `status.outputs.worker_service` | `string` | Name of the worker Service (`<name>-trino-worker` at default naming) — internal; exported for network-policy composition. |
| `status.outputs.port_forward_command` | `string` | Port-forward command for reaching the Trino Web UI from a workstation when no exposure is composed. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespace` | KubernetesNamespace | `spec.name` |
| `spec.catalogs.postgres[].host` | KubernetesPostgres | `status.outputs.rw_service` |
| `spec.catalogs.postgres[].passwordSecret.secretName` | KubernetesPostgres | `status.outputs.password_secret.name` |
| `spec.catalogs.mysql[].host` | KubernetesMysql | `status.outputs.primary_service` |
| `spec.catalogs.mysql[].passwordSecret.secretName` | KubernetesMysql | `status.outputs.root_password_secret.name` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
