# KubernetesOpenFga

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1`

**KubernetesOpenFgaSpec** installs OpenFGA — the CNCF authorization
engine implementing Google-Zanzibar-style relationship-based access
control — from the official `openfga` chart
(https://openfga.github.io/helm-charts, chart 0.3.x = OpenFGA 1.18+).

This component deploys the ENGINE (the OpenFGA server + its
datastore wiring). Authorization DATA — stores, authorization
models, relationship tuples — is managed through OpenFGA's own API
against the exported endpoints (the platform's OpenFGA provider
kinds OpenFgaStore / OpenFgaAuthorizationModel /
OpenFgaRelationshipTuple compose naturally: deploy the server here,
point an OpenFGA provider credential at `api_http_endpoint`, and
declare the authorization data as first-class resources).

PLAYGROUND, DELIBERATELY ABSENT: the chart ships its demo
playground ENABLED by default; this module always disables it
(verified at OpenFGA v1.18.1: upstream turned the playground off by
default for security, the server REFUSES TO START when it is
combined with any authentication method, and at this version it
binds pod-local only — the chart's playground Service port cannot
reach it anyway). Evaluate models with the `fga` CLI or the VS Code
extension against the API instead.

SCHEMA MIGRATIONS run as an init container in every server pod
(`openfga migrate` — idempotent, embedded migrations, no network
fetch). The chart's default hook-Job migration mode is not used:
a post-install hook Job deadlocks engines that wait on rollout
readiness, and its post-delete hook would dial the database during
uninstall.

## Example

```yaml
# Full-surface development manifest — exercises every module-rendered
# arm so the offline plan/preview proofs cover what the kind-cluster
# lanes exclude (the OIDC-free authn Secret materialization, tuned
# pools/limits, tracing, HPA, scheduling, ServiceMonitor, and the
# helm_values escape hatch with its fullnameOverride re-pin).
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesOpenFga
metadata:
  name: openfga-dev
spec:
  namespace:
    value: openfga-dev
  createNamespace: true
  chartVersion: 0.3.10
  datastore:
    postgres:
      host:
        value: openfga-dev-pg-rw
      port: 5432
      database: openfga
      username: openfga
      passwordSecret:
        secretName:
          value: openfga-dev-pg-app
        secretKey: password
      sslMode: disable
    migrationTimeout: 5m
    maxOpenConns: 40
    maxIdleConns: 10
    connMaxIdleTime: 5m
    connMaxLifetime: 1h
  authn:
    preshared:
      keys:
        - dev-only-proof-key
  metrics:
    enabled: true
    serviceMonitorEnabled: true
    enableRpcHistograms: true
  tracing:
    enabled: true
    otlpEndpoint: otel-collector.observability:4317
    sampleRatio: "0.5"
  log:
    level: debug
    format: json
  tuning:
    maxTuplesPerWrite: 50
    maxTypesPerAuthorizationModel: 50
    maxChecksPerBatchCheck: 25
    listObjectsDeadline: 5s
    listObjectsMaxResults: 500
    listUsersDeadline: 5s
    listUsersMaxResults: 500
    requestTimeout: 5s
    checkQueryCache:
      enabled: true
      limit: 5000
      ttl: 10s
  resources:
    requests:
      cpu: 100m
      memory: 128Mi
    limits:
      cpu: "1"
      memory: 512Mi
  hpa:
    enabled: true
    minReplicas: 2
    maxReplicas: 6
    targetCpuUtilizationPercent: 75
    targetMemoryUtilizationPercent: 80
  scheduling:
    nodeSelector:
      workload: authz
    tolerations:
      - key: authz
        operator: Equal
        value: "true"
        effect: NoSchedule
  serviceAccountAnnotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::111122223333:role/openfga-dev
  helmValues: |
    extraEnvVars:
      - name: OPENFGA_LOG_TIMESTAMP_FORMAT
        value: ISO8601
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.namespace` | `string \| valueFrom` | yes |  | KubernetesNamespace (`spec.name`) |
| `spec.createNamespace` | `bool` |  |  |  |
| `spec.chartVersion` | `string` |  | `0.3.10` |  |
| `spec.replicas` | `int32` |  | `3` |  |
| `spec.datastore` | `KubernetesOpenFgaDatastore` | yes |  |  |
| `spec.datastore.postgres` | `KubernetesOpenFgaPostgres` |  |  |  |
| `spec.datastore.postgres.host` | `string \| valueFrom` | yes |  | KubernetesPostgres (`status.outputs.rw_service`) |
| `spec.datastore.postgres.port` | `int32` |  | `5432` |  |
| `spec.datastore.postgres.database` | `string` | yes |  |  |
| `spec.datastore.postgres.username` | `string` | yes |  |  |
| `spec.datastore.postgres.passwordSecret` | `KubernetesOpenFgaPasswordSecret` | yes |  |  |
| `spec.datastore.postgres.passwordSecret.secretName` | `string \| valueFrom` | yes |  | KubernetesPostgres (`status.outputs.password_secret.name`) |
| `spec.datastore.postgres.passwordSecret.secretKey` | `string` |  | `password` |  |
| `spec.datastore.postgres.sslMode` | `string` |  | `disable` |  |
| `spec.datastore.mysql` | `KubernetesOpenFgaMysql` |  |  |  |
| `spec.datastore.mysql.host` | `string \| valueFrom` | yes |  | KubernetesMysql (`status.outputs.primary_service`) |
| `spec.datastore.mysql.port` | `int32` |  | `3306` |  |
| `spec.datastore.mysql.database` | `string` | yes |  |  |
| `spec.datastore.mysql.username` | `string` | yes |  |  |
| `spec.datastore.mysql.passwordSecret` | `KubernetesOpenFgaPasswordSecret` | yes |  |  |
| `spec.datastore.mysql.passwordSecret.secretName` | `string \| valueFrom` | yes |  | KubernetesPostgres (`status.outputs.password_secret.name`) |
| `spec.datastore.mysql.passwordSecret.secretKey` | `string` |  | `password` |  |
| `spec.datastore.memory` | `KubernetesOpenFgaMemory` |  |  |  |
| `spec.datastore.migrationTimeout` | `string` |  | `3m` |  |
| `spec.datastore.maxOpenConns` | `int32` |  |  |  |
| `spec.datastore.maxIdleConns` | `int32` |  |  |  |
| `spec.datastore.connMaxIdleTime` | `string` |  |  |  |
| `spec.datastore.connMaxLifetime` | `string` |  |  |  |
| `spec.authn` | `KubernetesOpenFgaAuthn` |  |  |  |
| `spec.authn.preshared` | `KubernetesOpenFgaPresharedAuthn` |  |  |  |
| `spec.authn.preshared.keys` | `[]string` (sensitive) |  |  |  |
| `spec.authn.preshared.existingKeysSecretName` | `string` |  |  |  |
| `spec.authn.oidc` | `KubernetesOpenFgaOidcAuthn` |  |  |  |
| `spec.authn.oidc.issuer` | `string` | yes |  |  |
| `spec.authn.oidc.audience` | `string` | yes |  |  |
| `spec.metrics` | `KubernetesOpenFgaMetrics` |  |  |  |
| `spec.metrics.enabled` | `bool` |  | `true` |  |
| `spec.metrics.serviceMonitorEnabled` | `bool` |  |  |  |
| `spec.metrics.enableRpcHistograms` | `bool` |  |  |  |
| `spec.tracing` | `KubernetesOpenFgaTracing` |  |  |  |
| `spec.tracing.enabled` | `bool` |  |  |  |
| `spec.tracing.otlpEndpoint` | `string` |  |  |  |
| `spec.tracing.sampleRatio` | `string` |  |  |  |
| `spec.log` | `KubernetesOpenFgaLog` |  |  |  |
| `spec.log.level` | `string` |  | `info` |  |
| `spec.log.format` | `string` |  | `json` |  |
| `spec.tuning` | `KubernetesOpenFgaTuning` |  |  |  |
| `spec.tuning.maxTuplesPerWrite` | `int32` |  |  |  |
| `spec.tuning.maxTypesPerAuthorizationModel` | `int32` |  |  |  |
| `spec.tuning.maxChecksPerBatchCheck` | `int32` |  |  |  |
| `spec.tuning.listObjectsDeadline` | `string` |  |  |  |
| `spec.tuning.listObjectsMaxResults` | `int32` |  |  |  |
| `spec.tuning.listUsersDeadline` | `string` |  |  |  |
| `spec.tuning.listUsersMaxResults` | `int32` |  |  |  |
| `spec.tuning.requestTimeout` | `string` |  |  |  |
| `spec.tuning.checkQueryCache` | `KubernetesOpenFgaCheckQueryCache` |  |  |  |
| `spec.tuning.checkQueryCache.enabled` | `bool` |  |  |  |
| `spec.tuning.checkQueryCache.limit` | `int32` |  |  |  |
| `spec.tuning.checkQueryCache.ttl` | `string` |  |  |  |
| `spec.tuning.experimentals` | `[]string` |  |  |  |
| `spec.resources` | `ContainerResources` |  |  |  |
| `spec.resources.limits` | `CpuMemory` |  |  |  |
| `spec.resources.limits.cpu` | `string` |  |  |  |
| `spec.resources.limits.memory` | `string` |  |  |  |
| `spec.resources.requests` | `CpuMemory` |  |  |  |
| `spec.resources.requests.cpu` | `string` |  |  |  |
| `spec.resources.requests.memory` | `string` |  |  |  |
| `spec.hpa` | `KubernetesOpenFgaHpa` |  |  |  |
| `spec.hpa.enabled` | `bool` |  |  |  |
| `spec.hpa.minReplicas` | `int32` |  | `1` |  |
| `spec.hpa.maxReplicas` | `int32` |  | `10` |  |
| `spec.hpa.targetCpuUtilizationPercent` | `int32` |  | `80` |  |
| `spec.hpa.targetMemoryUtilizationPercent` | `int32` |  | `80` |  |
| `spec.scheduling` | `KubernetesOpenFgaScheduling` |  |  |  |
| `spec.scheduling.nodeSelector` | `map<string, string>` |  |  |  |
| `spec.scheduling.tolerations` | `[]WorkloadToleration` |  |  |  |
| `spec.scheduling.tolerations[].key` | `string` |  |  |  |
| `spec.scheduling.tolerations[].operator` | `string` |  |  |  |
| `spec.scheduling.tolerations[].value` | `string` |  |  |  |
| `spec.scheduling.tolerations[].effect` | `string` |  |  |  |
| `spec.scheduling.tolerations[].tolerationSeconds` | `int64` |  |  |  |
| `spec.serviceAccountAnnotations` | `map<string, string>` |  |  |  |
| `spec.helmValues` | `string` |  |  |  |

## Field Details

### spec.namespace

`string | valueFrom` · required

Namespace to install into. Accepts a literal namespace name or a
reference to a KubernetesNamespace resource. NOTE the datastore
password rides a secretKeyRef, which can only read Secrets in
this same namespace — co-locate OpenFGA with its database or
replicate the credential Secret.

- references: KubernetesNamespace (`spec.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesNamespace, name: <that resource's name>, fieldPath: spec.name}} -- a bare string does not parse

### spec.createNamespace

`bool`

When true, the namespace is created (with the standard Planton
governance labels) before installing and deleted with the
resource. When false, the namespace must already exist.

### spec.chartVersion

`string` · optional (explicit presence)

Helm chart version to install (e.g. "0.3.10" = OpenFGA v1.18.1).
Versions must exist in the SERVED index at
https://openfga.github.io/helm-charts.

- default: `0.3.10`

### spec.replicas

`int32` · optional (explicit presence)

Server replicas. OpenFGA is stateless (all state lives in the
datastore), so replicas scale reads and checks linearly. Ignored
when hpa is enabled (the autoscaler owns the count), and forced
to 1 on the memory datastore (each replica would hold its own
divergent world).

- default: `3`
- rule: {"int32":{"lte":50,"gte":1}}

### spec.datastore

`KubernetesOpenFgaDatastore` · required

The datastore backing all authorization data.

- rule: {"required":true}
- rule: Pick a datastore engine: postgres (recommended), mysql, or memory (evaluation only).

### spec.datastore.postgres

`KubernetesOpenFgaPostgres`

PostgreSQL — the recommended production engine. Defaults
compose a KubernetesPostgres resource.

### spec.datastore.postgres.host

`string | valueFrom` · required

PostgreSQL host — a Service name (same namespace) or a full FQDN
(cross-namespace or external). Accepts a literal or a reference
to a KubernetesPostgres resource (its read-write Service —
always the current primary). The port is declared separately.

- references: KubernetesPostgres (`status.outputs.rw_service`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesPostgres, name: <that resource's name>, fieldPath: status.outputs.rw_service}} -- a bare string does not parse

### spec.datastore.postgres.port

`int32` · optional (explicit presence)

PostgreSQL port. Empty = 5432.

- default: `5432`
- rule: {"int32":{"lte":65535,"gt":0}}

### spec.datastore.postgres.database

`string` · required

Database to hold the OpenFGA schema (on a KubernetesPostgres:
declare it at bootstrap via initdb, e.g. "openfga").

- rule: {"required":true}

### spec.datastore.postgres.username

`string` · required

Database username (on a KubernetesPostgres this is the bootstrap
owner role — ownership covers everything the migrations need).

- rule: {"required":true}

### spec.datastore.postgres.passwordSecret

`KubernetesOpenFgaPasswordSecret` · required

The user's password, read from an existing Secret.

- rule: {"required":true}

### spec.datastore.postgres.passwordSecret.secretName

`string | valueFrom` · required

Secret name. Accepts a literal or a reference to a
KubernetesPostgres resource (its `<cluster>-app` credential
Secret, maintained by the operator across failovers). KNOW THIS
(a Kubernetes constraint, not a chart one): a secretKeyRef can
only read Secrets in the workload's OWN namespace — co-locate
OpenFGA with its database or replicate the Secret.

- references: KubernetesPostgres (`status.outputs.password_secret.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesPostgres, name: <that resource's name>, fieldPath: status.outputs.password_secret.name}} -- a bare string does not parse

### spec.datastore.postgres.passwordSecret.secretKey

`string` · optional (explicit presence)

Key within the Secret holding the password. Empty = "password"
(the key a KubernetesPostgres application Secret uses).

- default: `password`

### spec.datastore.postgres.sslMode

`string` · optional (explicit presence)

Postgres sslmode for the connection (disable, require,
verify-ca, verify-full). Empty = "disable" — the in-cluster
plaintext default; set verify-full for external/managed
databases.

- default: `disable`
- rule: sslmode must be one of: disable, require, verify-ca, verify-full.

### spec.datastore.mysql

`KubernetesOpenFgaMysql`

MySQL 8. Defaults compose a KubernetesMysql resource.

### spec.datastore.mysql.host

`string | valueFrom` · required

MySQL host — a Service name (same namespace) or a full FQDN.
Accepts a literal or a reference to a KubernetesMysql resource
(its primary Service — writes always land on the current
primary). The port is declared separately.

- references: KubernetesMysql (`status.outputs.primary_service`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesMysql, name: <that resource's name>, fieldPath: status.outputs.primary_service}} -- a bare string does not parse

### spec.datastore.mysql.port

`int32` · optional (explicit presence)

MySQL port. Empty = 3306.

- default: `3306`
- rule: {"int32":{"lte":65535,"gt":0}}

### spec.datastore.mysql.database

`string` · required

Database to hold the OpenFGA schema.

- rule: {"required":true}

### spec.datastore.mysql.username

`string` · required

Database username.

- rule: {"required":true}

### spec.datastore.mysql.passwordSecret

`KubernetesOpenFgaPasswordSecret` · required

The user's password, read from an existing Secret.

- rule: {"required":true}

### spec.datastore.mysql.passwordSecret.secretName

`string | valueFrom` · required

Secret name. Accepts a literal or a reference to a
KubernetesPostgres resource (its `<cluster>-app` credential
Secret, maintained by the operator across failovers). KNOW THIS
(a Kubernetes constraint, not a chart one): a secretKeyRef can
only read Secrets in the workload's OWN namespace — co-locate
OpenFGA with its database or replicate the Secret.

- references: KubernetesPostgres (`status.outputs.password_secret.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesPostgres, name: <that resource's name>, fieldPath: status.outputs.password_secret.name}} -- a bare string does not parse

### spec.datastore.mysql.passwordSecret.secretKey

`string` · optional (explicit presence)

Key within the Secret holding the password. Empty = "password"
(the key a KubernetesPostgres application Secret uses).

- default: `password`

### spec.datastore.memory

`KubernetesOpenFgaMemory`

In-memory — NEVER for real use: data is lost on every restart
and replicas are forced to 1 (each pod would hold its own
divergent authorization world). Exists as the zero-dependency
evaluation arm.

### spec.datastore.migrationTimeout

`string` · optional (explicit presence)

How long the schema-migration init container retries an
unreachable database before failing (Go duration; the server's
own default is "1m"). Widen when the database is provisioned
concurrently with OpenFGA (a composed install's normal shape).

- default: `3m`
- rule: {"string":{"pattern":"^\\d+(\\.\\d+)?(ms|s|m|h)$"}}

### spec.datastore.maxOpenConns

`int32` · optional (explicit presence)

Connection pool: maximum open connections per replica.
Empty = the server default (30).

- rule: {"int32":{"gt":0}}

### spec.datastore.maxIdleConns

`int32` · optional (explicit presence)

Connection pool: maximum idle connections per replica.
Empty = the server default (10).

- rule: {"int32":{"gt":0}}

### spec.datastore.connMaxIdleTime

`string` · optional (explicit presence)

Maximum idle time for a pooled connection (Go duration).

- rule: {"string":{"pattern":"^\\d+(\\.\\d+)?(ms|s|m|h)$"}}

### spec.datastore.connMaxLifetime

`string` · optional (explicit presence)

Maximum lifetime of a pooled connection (Go duration).

- rule: {"string":{"pattern":"^\\d+(\\.\\d+)?(ms|s|m|h)$"}}

### spec.authn

`KubernetesOpenFgaAuthn`

API authentication. Unset = NO authentication — anyone who can
reach the Service can read and write every store. Fine on a lab
cluster, never in production.

### spec.authn.preshared

`KubernetesOpenFgaPresharedAuthn`

Static pre-shared API keys sent as `Authorization: Bearer <key>`.

- rule: Provide the API keys exactly one way: declare them in `keys` (materialized into a managed Secret) or point `existing_keys_secret_name` at a Secret you maintain.

### spec.authn.preshared.keys

`[]string` · sensitive

Declare keys here and the module materializes them into a
module-owned Secret (`<name>-authn-keys`) before the release;
only the Secret NAME ever renders. Mutually exclusive with
existing_keys_secret_name.

### spec.authn.preshared.existingKeysSecretName

`string`

OR reference an existing Secret carrying the comma-separated key
list under the data key `keys` (the chart's contract).

### spec.authn.oidc

`KubernetesOpenFgaOidcAuthn`

OIDC bearer tokens validated against an issuer.

### spec.authn.oidc.issuer

`string` · required

Issuer URL (the server fetches its JWKS for token validation).

- rule: {"required":true,"string":{"uri":true}}

### spec.authn.oidc.audience

`string` · required

Required `aud` claim value.

- rule: {"required":true}

### spec.metrics

`KubernetesOpenFgaMetrics`

Prometheus metrics (chart default on, port 2112) and the
optional ServiceMonitor.

### spec.metrics.enabled

`bool` · optional (explicit presence)

Serve /metrics on port 2112. Empty = true (the chart default).

- default: `true`

### spec.metrics.serviceMonitorEnabled

`bool`

Also create a ServiceMonitor (requires the Prometheus Operator
CRDs on the cluster; the install FAILS without them).

### spec.metrics.enableRpcHistograms

`bool`

Export per-RPC latency histograms (higher cardinality, better
latency visibility).

### spec.tracing

`KubernetesOpenFgaTracing`

OpenTelemetry trace export.

- rule: Tracing needs somewhere to send spans — set otlp_endpoint (host:port of an OTLP gRPC collector) when tracing is enabled.

### spec.tracing.enabled

`bool`

Enable OTLP trace export.

### spec.tracing.otlpEndpoint

`string`

OTLP gRPC endpoint (host:port, e.g. a KubernetesOtelCollector or
KubernetesSignoz collector Service). Required when enabled.

### spec.tracing.sampleRatio

`string` · optional (explicit presence)

Sampling ratio 0.0–1.0. Empty = the server default (0.2).

- rule: {"string":{"pattern":"^(0(\\.\\d+)?|1(\\.0+)?)$"}}

### spec.log

`KubernetesOpenFgaLog`

Server log settings.

### spec.log.level

`string` · optional (explicit presence)

Log level: debug, info (default), warn, error, panic, fatal.
The OpenFGA server itself also accepts "none", but the Helm chart's
closed values schema rejects it, so it is not offered here.

- default: `info`
- rule: Log level must be one of: debug, info, warn, error, panic, fatal.

### spec.log.format

`string` · optional (explicit presence)

Log format: json (default) or text.

- default: `json`
- rule: Log format must be either "json" or "text".

### spec.tuning

`KubernetesOpenFgaTuning`

Query limits and performance tuning. Every field maps to a
server flag; empty fields keep the server's own defaults.

### spec.tuning.maxTuplesPerWrite

`int32` · optional (explicit presence)

Maximum tuples per single Write call. Server default: 100.

- rule: {"int32":{"gt":0}}

### spec.tuning.maxTypesPerAuthorizationModel

`int32` · optional (explicit presence)

Maximum type definitions per authorization model. Server
default: 100.

- rule: {"int32":{"gt":0}}

### spec.tuning.maxChecksPerBatchCheck

`int32` · optional (explicit presence)

Maximum checks in one BatchCheck call. Server default: 50.

- rule: {"int32":{"gt":0}}

### spec.tuning.listObjectsDeadline

`string` · optional (explicit presence)

Deadline for ListObjects queries (Go duration). Server
default: "3s".

- rule: {"string":{"pattern":"^\\d+(\\.\\d+)?(ms|s|m|h)$"}}

### spec.tuning.listObjectsMaxResults

`int32` · optional (explicit presence)

Maximum results from ListObjects. Server default: 1000.

- rule: {"int32":{"gte":0}}

### spec.tuning.listUsersDeadline

`string` · optional (explicit presence)

Deadline for ListUsers queries (Go duration). Server
default: "3s".

- rule: {"string":{"pattern":"^\\d+(\\.\\d+)?(ms|s|m|h)$"}}

### spec.tuning.listUsersMaxResults

`int32` · optional (explicit presence)

Maximum results from ListUsers. Server default: 1000.

- rule: {"int32":{"gte":0}}

### spec.tuning.requestTimeout

`string` · optional (explicit presence)

Global request timeout (Go duration). Server default: "3s".

- rule: {"string":{"pattern":"^\\d+(\\.\\d+)?(ms|s|m|h)$"}}

### spec.tuning.checkQueryCache

`KubernetesOpenFgaCheckQueryCache`

Cache Check results for hot authorization paths.

### spec.tuning.checkQueryCache.enabled

`bool`

Enable the cache. Trades authorization-change propagation delay
(up to ttl) for large Check throughput gains.

### spec.tuning.checkQueryCache.limit

`int32` · optional (explicit presence)

Maximum cached entries. Empty = the server default (10000).

- rule: {"int32":{"gt":0}}

### spec.tuning.checkQueryCache.ttl

`string` · optional (explicit presence)

Entry time-to-live (Go duration). Empty = the server
default ("10s").

- rule: {"string":{"pattern":"^\\d+(\\.\\d+)?(ms|s|m|h)$"}}

### spec.tuning.experimentals

`[]string`

Experimental server features. The chart's CLOSED values schema
accepts exactly these four (the server knows more, but anything
else fails the install with an opaque schema error — validated
here instead): enable-check-optimizations, enable-access-control,
enable-list-objects-optimizations, pipeline_list_objects.
KNOW THIS (server contract at v1.18.1): the list REPLACES the
server's own default experimental set (which ships
"pipeline_list_objects" enabled) — declaring any value here drops
the defaults unless you re-list them.

- rule: The chart's closed values schema accepts only: enable-check-optimizations, enable-access-control, enable-list-objects-optimizations, pipeline_list_objects.

### spec.resources

`ContainerResources`

CPU and memory for the server container. The chart ships NO
defaults at all; these are modest laboratory defaults — size
real installs to check volume.

### spec.resources.limits

`CpuMemory`

The resource limits for the container.
Specify the maximum amount of CPU and memory that the container can use.

### spec.resources.limits.cpu

`string`

### spec.resources.limits.memory

`string`

### spec.resources.requests

`CpuMemory`

The resource requests for the container.
Specify the minimum amount of CPU and memory that the container is guaranteed.

### spec.resources.requests.cpu

`string`

### spec.resources.requests.memory

`string`

### spec.hpa

`KubernetesOpenFgaHpa`

Horizontal pod autoscaling (CPU/memory targets). Not available on
the memory datastore: the chart pins memory-engine replicas to 1
only while autoscaling is OFF — an HPA would scale pods that each
hold their own divergent authorization world.

- rule: The autoscaler's minimum replica count cannot exceed its maximum.

### spec.hpa.enabled

`bool`

Enable the HorizontalPodAutoscaler (replicas then belongs to the
autoscaler).

### spec.hpa.minReplicas

`int32` · optional (explicit presence)

Minimum replicas. Empty = 1.

- default: `1`
- rule: {"int32":{"gte":1}}

### spec.hpa.maxReplicas

`int32` · optional (explicit presence)

Maximum replicas. Empty = 10 (the platform default; the chart's
own fallback when nothing renders is 100).

- default: `10`
- rule: {"int32":{"gte":1}}

### spec.hpa.targetCpuUtilizationPercent

`int32` · optional (explicit presence)

Target average CPU utilization percent. Empty = 80.

- default: `80`
- rule: {"int32":{"lte":100,"gte":1}}

### spec.hpa.targetMemoryUtilizationPercent

`int32` · optional (explicit presence)

Target average memory utilization percent. Empty = 80.

- default: `80`
- rule: {"int32":{"lte":100,"gte":1}}

### spec.scheduling

`KubernetesOpenFgaScheduling`

Pod scheduling constraints.

### spec.scheduling.nodeSelector

`map<string, string>`

Schedule onto nodes carrying these labels.

### spec.scheduling.tolerations

`[]WorkloadToleration`

Tolerations for tainted nodes.

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

### spec.serviceAccountAnnotations

`map<string, string>`

Annotations for the server ServiceAccount (the cloud
workload-identity seam).

### spec.helmValues

`string`

Advanced escape hatch: raw Helm values merged LAST (Helm `-f`
semantics) over everything this spec renders — later keys win.
KNOW THIS (verified at chart 0.3.10): the chart ships a CLOSED
values schema (`additionalProperties: false`) — a key the chart
does not define fails the install outright, so this hatch can
only override EXISTING chart values (extraEnvVars for the ~50
server flags without values paths, TLS file wiring, sidecars),
never invent new ones. The module re-pins `fullnameOverride`
after the merge. YAML document as a string.

## Validation Rules

- `spec.hpa.memory_engine`: Autoscaling on the memory datastore scales pods that each hold their own divergent in-memory authorization world — use postgres or mysql to scale out.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesOpenFga, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.namespace` | `string` | Namespace the server runs in. |
| `status.outputs.service` | `string` | The OpenFGA Service name. |
| `status.outputs.api_http_endpoint` | `string` | In-cluster HTTP API endpoint (e.g. "http://fga.authz.svc.cluster.local:8080") — the REST surface SDKs and the platform's OpenFGA provider connect to. |
| `status.outputs.api_grpc_endpoint` | `string` | In-cluster gRPC API endpoint host:port (e.g. "fga.authz.svc.cluster.local:8081") — plaintext gRPC. |
| `status.outputs.authn_keys_secret_name` | `string` | Name of the module-owned Secret holding declared pre-shared API keys (`<name>-authn-keys`, data key `keys`); empty when authn is unset or rides an existing Secret. |
| `status.outputs.port_forward_command` | `string` | Copy-paste command for reaching the HTTP API from a workstation. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespace` | KubernetesNamespace | `spec.name` |
| `spec.datastore.postgres.host` | KubernetesPostgres | `status.outputs.rw_service` |
| `spec.datastore.postgres.passwordSecret.secretName` | KubernetesPostgres | `status.outputs.password_secret.name` |
| `spec.datastore.mysql.host` | KubernetesMysql | `status.outputs.primary_service` |
| `spec.datastore.mysql.passwordSecret.secretName` | KubernetesPostgres | `status.outputs.password_secret.name` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
