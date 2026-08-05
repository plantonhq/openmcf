# KubernetesArgocd

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1`

**KubernetesArgocdSpec** deploys Argo CD — the declarative GitOps
continuous-delivery controller — from the official `argo-cd` Helm
chart (https://argoproj.github.io/argo-helm).

WHAT GETS INSTALLED: the application controller (reconciles
Applications against Git), the API/UI server, the repo server
(clones and renders manifests), the ApplicationSet controller
(templated app generation), the notifications controller (on by
default), Dex (SSO broker, on by default), and a Redis used as a
throwaway cache — bundled single-pod by default, switchable to the
redis-ha subchart or an external endpoint via `redis`.

FIRST LOGIN: on first start Argo CD generates an admin password into
the fixed-name Secret `argocd-initial-admin-secret` (key `password`)
in the install namespace — exported as an output handle. Because the
name is fixed by the application, ONE generated-password Argo CD per
namespace; a second instance needs its own namespace. Rotate by
changing the password through the CLI/API and deleting that Secret.

DECLARING APPLICATIONS: this kind installs the CONTROL PLANE only.
Applications, AppProjects and ApplicationSets are Kubernetes custom
resources declared like any other manifest (KubernetesManifest, a
chart, or Argo CD's own UI/CLI) once the control plane runs.

REPOSITORY CREDENTIALS: private-repo credentials are Kubernetes
Secrets labeled `argocd.argoproj.io/secret-type: repository` (or
`repo-creds` for credential templates) that Argo CD discovers
natively — compose them from KubernetesSecret or
KubernetesExternalSecret; nothing in this spec transports Git
credentials. `repositories` below registers PUBLIC (anonymous)
repositories only.

EXPOSURE: services stay ClusterIP; expose the server via first-class
kinds (KubernetesIngress, Gateway API kinds) over the exported
service handle. When exposure terminates TLS in front of Argo CD,
set `server.insecure` so the server speaks plain HTTP behind the
proxy, and set `domain` so SSO redirects and CLI hints carry the
public name.

The typed fields below cover the chart's meaningful configuration
surface; `helm_values` remains as the escape hatch for chart values
beyond them (merged last, Helm `-f` semantics, identical on both
engines) — extensions, per-component overrides, notification
templates — a safety valve, never the primary interface. Never put
secret material in `helm_values`; every credential path in this spec
rides existing Secrets.

## Example

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesArgocd
metadata:
  name: test-argocd
spec:
  namespace:
    value: gitops
  createNamespace: true
  chartVersion: 10.2.1
  adminEnabled: true
  domain: argocd.example.com
  controller:
    replicas: 1
    resources:
      requests:
        cpu: 250m
        memory: 512Mi
      limits:
        cpu: "1"
        memory: 1Gi
  server:
    insecure: true
    resources:
      requests:
        cpu: 100m
        memory: 128Mi
    autoscaling:
      enabled: true
      minReplicas: 2
      maxReplicas: 5
  repoServer:
    autoscaling:
      enabled: true
      minReplicas: 2
      maxReplicas: 4
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
  applicationSet:
    replicas: 2
  notifications:
    enabled: true
  dex:
    enabled: true
  commitServer:
    enabled: true
  redis:
    ha:
      replicas: 3
  sso:
    oidc:
      name: Okta
      issuer: https://example.okta.com
      clientId: argocd
      clientSecretSecret:
        name: argocd-oidc
        key: clientSecret
    dexConfig: |
      connectors:
        - type: github
          id: github
          name: GitHub
          config:
            clientID: abc123
            clientSecret: $argocd-github-sso:clientSecret
            orgs:
              - name: example-org
  rbac:
    policyDefault: role:readonly
    policyCsv: |
      g, example-org:platform, role:admin
    scopes: "[groups]"
  execEnabled: true
  reconciliationTimeout: 180s
  repositories:
    - name: platform
      url: https://github.com/example-org/platform
    - name: charts
      url: https://charts.example.com
      type: helm
  crds:
    install: true
    keep: true
  serviceMonitorsEnabled: true
  image:
    repository: my.registry.com/argoproj/argocd
    tag: v3.4.5
    pullSecretName: mirror-pull
  scheduling:
    nodeSelector:
      role: platform
    tolerations:
      - key: platform
        operator: Equal
        value: "true"
        effect: NoSchedule
    priorityClassName: platform-critical
  helmValues: |
    notifications:
      argocdUrl: https://argocd.example.com
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.namespace` | `string \| valueFrom` | yes |  | KubernetesNamespace (`spec.name`) |
| `spec.createNamespace` | `bool` |  |  |  |
| `spec.chartVersion` | `string` |  | `10.2.1` |  |
| `spec.adminEnabled` | `bool` |  | `true` |  |
| `spec.domain` | `string` |  |  |  |
| `spec.controller` | `KubernetesArgocdComponent` |  |  |  |
| `spec.controller.replicas` | `int32` |  | `1` |  |
| `spec.controller.resources` | `ContainerResources` |  |  |  |
| `spec.controller.resources.limits` | `CpuMemory` |  |  |  |
| `spec.controller.resources.limits.cpu` | `string` |  |  |  |
| `spec.controller.resources.limits.memory` | `string` |  |  |  |
| `spec.controller.resources.requests` | `CpuMemory` |  |  |  |
| `spec.controller.resources.requests.cpu` | `string` |  |  |  |
| `spec.controller.resources.requests.memory` | `string` |  |  |  |
| `spec.server` | `KubernetesArgocdServer` |  |  |  |
| `spec.server.replicas` | `int32` |  | `1` |  |
| `spec.server.resources` | `ContainerResources` |  |  |  |
| `spec.server.resources.limits` | `CpuMemory` |  |  |  |
| `spec.server.resources.limits.cpu` | `string` |  |  |  |
| `spec.server.resources.limits.memory` | `string` |  |  |  |
| `spec.server.resources.requests` | `CpuMemory` |  |  |  |
| `spec.server.resources.requests.cpu` | `string` |  |  |  |
| `spec.server.resources.requests.memory` | `string` |  |  |  |
| `spec.server.autoscaling` | `KubernetesArgocdAutoscaling` |  |  |  |
| `spec.server.autoscaling.enabled` | `bool` |  |  |  |
| `spec.server.autoscaling.minReplicas` | `int32` |  | `1` |  |
| `spec.server.autoscaling.maxReplicas` | `int32` |  | `5` |  |
| `spec.server.insecure` | `bool` |  |  |  |
| `spec.repoServer` | `KubernetesArgocdScalableComponent` |  |  |  |
| `spec.repoServer.replicas` | `int32` |  | `1` |  |
| `spec.repoServer.resources` | `ContainerResources` |  |  |  |
| `spec.repoServer.resources.limits` | `CpuMemory` |  |  |  |
| `spec.repoServer.resources.limits.cpu` | `string` |  |  |  |
| `spec.repoServer.resources.limits.memory` | `string` |  |  |  |
| `spec.repoServer.resources.requests` | `CpuMemory` |  |  |  |
| `spec.repoServer.resources.requests.cpu` | `string` |  |  |  |
| `spec.repoServer.resources.requests.memory` | `string` |  |  |  |
| `spec.repoServer.autoscaling` | `KubernetesArgocdAutoscaling` |  |  |  |
| `spec.repoServer.autoscaling.enabled` | `bool` |  |  |  |
| `spec.repoServer.autoscaling.minReplicas` | `int32` |  | `1` |  |
| `spec.repoServer.autoscaling.maxReplicas` | `int32` |  | `5` |  |
| `spec.applicationSet` | `KubernetesArgocdComponent` |  |  |  |
| `spec.applicationSet.replicas` | `int32` |  | `1` |  |
| `spec.applicationSet.resources` | `ContainerResources` |  |  |  |
| `spec.applicationSet.resources.limits` | `CpuMemory` |  |  |  |
| `spec.applicationSet.resources.limits.cpu` | `string` |  |  |  |
| `spec.applicationSet.resources.limits.memory` | `string` |  |  |  |
| `spec.applicationSet.resources.requests` | `CpuMemory` |  |  |  |
| `spec.applicationSet.resources.requests.cpu` | `string` |  |  |  |
| `spec.applicationSet.resources.requests.memory` | `string` |  |  |  |
| `spec.notifications` | `KubernetesArgocdToggleableComponent` |  |  |  |
| `spec.notifications.enabled` | `bool` |  |  |  |
| `spec.notifications.resources` | `ContainerResources` |  |  |  |
| `spec.notifications.resources.limits` | `CpuMemory` |  |  |  |
| `spec.notifications.resources.limits.cpu` | `string` |  |  |  |
| `spec.notifications.resources.limits.memory` | `string` |  |  |  |
| `spec.notifications.resources.requests` | `CpuMemory` |  |  |  |
| `spec.notifications.resources.requests.cpu` | `string` |  |  |  |
| `spec.notifications.resources.requests.memory` | `string` |  |  |  |
| `spec.dex` | `KubernetesArgocdToggleableComponent` |  |  |  |
| `spec.dex.enabled` | `bool` |  |  |  |
| `spec.dex.resources` | `ContainerResources` |  |  |  |
| `spec.dex.resources.limits` | `CpuMemory` |  |  |  |
| `spec.dex.resources.limits.cpu` | `string` |  |  |  |
| `spec.dex.resources.limits.memory` | `string` |  |  |  |
| `spec.dex.resources.requests` | `CpuMemory` |  |  |  |
| `spec.dex.resources.requests.cpu` | `string` |  |  |  |
| `spec.dex.resources.requests.memory` | `string` |  |  |  |
| `spec.commitServer` | `KubernetesArgocdToggleableComponent` |  |  |  |
| `spec.commitServer.enabled` | `bool` |  |  |  |
| `spec.commitServer.resources` | `ContainerResources` |  |  |  |
| `spec.commitServer.resources.limits` | `CpuMemory` |  |  |  |
| `spec.commitServer.resources.limits.cpu` | `string` |  |  |  |
| `spec.commitServer.resources.limits.memory` | `string` |  |  |  |
| `spec.commitServer.resources.requests` | `CpuMemory` |  |  |  |
| `spec.commitServer.resources.requests.cpu` | `string` |  |  |  |
| `spec.commitServer.resources.requests.memory` | `string` |  |  |  |
| `spec.redis` | `KubernetesArgocdRedis` |  |  |  |
| `spec.redis.bundled` | `KubernetesArgocdRedisBundled` |  |  |  |
| `spec.redis.bundled.resources` | `ContainerResources` |  |  |  |
| `spec.redis.bundled.resources.limits` | `CpuMemory` |  |  |  |
| `spec.redis.bundled.resources.limits.cpu` | `string` |  |  |  |
| `spec.redis.bundled.resources.limits.memory` | `string` |  |  |  |
| `spec.redis.bundled.resources.requests` | `CpuMemory` |  |  |  |
| `spec.redis.bundled.resources.requests.cpu` | `string` |  |  |  |
| `spec.redis.bundled.resources.requests.memory` | `string` |  |  |  |
| `spec.redis.ha` | `KubernetesArgocdRedisHa` |  |  |  |
| `spec.redis.ha.replicas` | `int32` |  | `3` |  |
| `spec.redis.external` | `KubernetesArgocdRedisExternal` |  |  |  |
| `spec.redis.external.host` | `string \| valueFrom` | yes |  | KubernetesValkey (`status.outputs.service`) |
| `spec.redis.external.port` | `int32` |  | `6379` |  |
| `spec.redis.external.credentialsSecretName` | `string` |  |  |  |
| `spec.sso` | `KubernetesArgocdSso` |  |  |  |
| `spec.sso.oidc` | `KubernetesArgocdOidc` |  |  |  |
| `spec.sso.oidc.name` | `string` | yes |  |  |
| `spec.sso.oidc.issuer` | `string` | yes |  |  |
| `spec.sso.oidc.clientId` | `string` | yes |  |  |
| `spec.sso.oidc.clientSecretSecret` | `KubernetesArgocdSecretKeyRef` |  |  |  |
| `spec.sso.oidc.clientSecretSecret.name` | `string` | yes |  |  |
| `spec.sso.oidc.clientSecretSecret.key` | `string` | yes |  |  |
| `spec.sso.dexConfig` | `string` |  |  |  |
| `spec.rbac` | `KubernetesArgocdRbac` |  |  |  |
| `spec.rbac.policyDefault` | `string` |  |  |  |
| `spec.rbac.policyCsv` | `string` |  |  |  |
| `spec.rbac.scopes` | `string` |  | `[groups]` |  |
| `spec.execEnabled` | `bool` |  |  |  |
| `spec.reconciliationTimeout` | `string` |  |  |  |
| `spec.repositories` | `[]KubernetesArgocdRepository` |  |  |  |
| `spec.repositories[].name` | `string` | yes |  |  |
| `spec.repositories[].url` | `string` | yes |  |  |
| `spec.repositories[].type` | `string` |  | `git` |  |
| `spec.crds` | `KubernetesArgocdCrds` |  |  |  |
| `spec.crds.install` | `bool` |  | `true` |  |
| `spec.crds.keep` | `bool` |  | `true` |  |
| `spec.serviceMonitorsEnabled` | `bool` |  |  |  |
| `spec.image` | `KubernetesArgocdImage` |  |  |  |
| `spec.image.repository` | `string` |  |  |  |
| `spec.image.tag` | `string` |  |  |  |
| `spec.image.pullSecretName` | `string` |  |  |  |
| `spec.scheduling` | `KubernetesArgocdScheduling` |  |  |  |
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

Helm chart version to install (e.g. "10.2.1" — chart 10.2.1 ships
Argo CD v3.4.5). Versions must exist as SERVED charts in the
repository index (https://argoproj.github.io/argo-helm).

- default: `10.2.1`

### spec.adminEnabled

`bool` · optional (explicit presence)

Enable the local `admin` user. Default true — with it comes the
generated `argocd-initial-admin-secret` (see the FIRST LOGIN note
on the spec). Disable only AFTER SSO works, or the UI locks
everyone out.

- default: `true`

### spec.domain

`string`

The public domain users reach Argo CD at (e.g.
"argocd.example.com" — name only, no scheme). SSO redirect URLs,
CLI login hints and notification links embed it — set it whenever
exposure is composed in front of this Argo CD. Empty = the chart's
placeholder ("argocd.example.com").

### spec.controller

`KubernetesArgocdComponent`

The application controller — the reconciliation engine. Scale
replicas only when sharding across many target clusters (each
replica owns a shard; the chart wires the shard count
automatically).

### spec.controller.replicas

`int32` · optional (explicit presence)

Number of replicas. Empty = 1.

- default: `1`
- rule: {"int32":{"gte":1}}

### spec.controller.resources

`ContainerResources`

CPU and memory for the component's container. Empty = no
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

### spec.server

`KubernetesArgocdServer`

The API/UI server.

### spec.server.replicas

`int32` · optional (explicit presence)

Number of replicas. Empty = 1. Ignored when `autoscaling` is
enabled.

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

### spec.server.autoscaling

`KubernetesArgocdAutoscaling`

Horizontal pod autoscaling for the server (requires
metrics-server).

- rule: autoscaling max_replicas must be greater than or equal to min_replicas

### spec.server.autoscaling.enabled

`bool`

Enable the HPA.

### spec.server.autoscaling.minReplicas

`int32` · optional (explicit presence)

Minimum replicas. Empty = 1 (the chart default); the HA recipe uses
2.

- default: `1`
- rule: {"int32":{"gte":1}}

### spec.server.autoscaling.maxReplicas

`int32` · optional (explicit presence)

Maximum replicas. Empty = 5 (the chart default).

- default: `5`
- rule: {"int32":{"gte":1}}

### spec.server.insecure

`bool`

Serve plain HTTP instead of the server's self-signed TLS. Set it
when a composed ingress/gateway terminates TLS in front of Argo CD
(the standard exposure pattern) — leaving the default self-signed
listener behind a TLS-terminating proxy double-encrypts and breaks
gRPC CLI traffic.

### spec.repoServer

`KubernetesArgocdScalableComponent`

The repo server — clones repositories and renders manifests. The
component to scale first when sync latency grows with many apps.

### spec.repoServer.replicas

`int32` · optional (explicit presence)

Number of replicas. Empty = 1. Ignored when `autoscaling` is
enabled (the HPA owns the count).

- default: `1`
- rule: {"int32":{"gte":1}}

### spec.repoServer.resources

`ContainerResources`

CPU and memory for the component's container. Empty = no
requests/limits (the chart default). Declare requests when
autoscaling — a CPU-utilization HPA without a CPU request never
scales.

### spec.repoServer.resources.limits

`CpuMemory`

The resource limits for the container.
Specify the maximum amount of CPU and memory that the container can use.

### spec.repoServer.resources.limits.cpu

`string`

### spec.repoServer.resources.limits.memory

`string`

### spec.repoServer.resources.requests

`CpuMemory`

The resource requests for the container.
Specify the minimum amount of CPU and memory that the container is guaranteed.

### spec.repoServer.resources.requests.cpu

`string`

### spec.repoServer.resources.requests.memory

`string`

### spec.repoServer.autoscaling

`KubernetesArgocdAutoscaling`

Horizontal pod autoscaling (requires metrics-server on the
cluster).

- rule: autoscaling max_replicas must be greater than or equal to min_replicas

### spec.repoServer.autoscaling.enabled

`bool`

Enable the HPA.

### spec.repoServer.autoscaling.minReplicas

`int32` · optional (explicit presence)

Minimum replicas. Empty = 1 (the chart default); the HA recipe uses
2.

- default: `1`
- rule: {"int32":{"gte":1}}

### spec.repoServer.autoscaling.maxReplicas

`int32` · optional (explicit presence)

Maximum replicas. Empty = 5 (the chart default).

- default: `5`
- rule: {"int32":{"gte":1}}

### spec.applicationSet

`KubernetesArgocdComponent`

The ApplicationSet controller (always installed by the chart).

### spec.applicationSet.replicas

`int32` · optional (explicit presence)

Number of replicas. Empty = 1.

- default: `1`
- rule: {"int32":{"gte":1}}

### spec.applicationSet.resources

`ContainerResources`

CPU and memory for the component's container. Empty = no
requests/limits (the chart default).

### spec.applicationSet.resources.limits

`CpuMemory`

The resource limits for the container.
Specify the maximum amount of CPU and memory that the container can use.

### spec.applicationSet.resources.limits.cpu

`string`

### spec.applicationSet.resources.limits.memory

`string`

### spec.applicationSet.resources.requests

`CpuMemory`

The resource requests for the container.
Specify the minimum amount of CPU and memory that the container is guaranteed.

### spec.applicationSet.resources.requests.cpu

`string`

### spec.applicationSet.resources.requests.memory

`string`

### spec.notifications

`KubernetesArgocdToggleableComponent`

The notifications controller — delivers sync/health events to
chat/webhook/email services. Chart default: enabled. Notifier
credentials ride Secrets referenced from the notifications
configuration (`$<secret-name>:<key>` indirection) — configure
notifiers via `helm_values` (`notifications.notifiers`, templates,
triggers).

### spec.notifications.enabled

`bool` · optional (explicit presence)

Enable the component. Empty = the chart default (notifications and
dex: enabled; commit server: disabled).

### spec.notifications.resources

`ContainerResources`

CPU and memory for the component's container. Empty = no
requests/limits (the chart default).

### spec.notifications.resources.limits

`CpuMemory`

The resource limits for the container.
Specify the maximum amount of CPU and memory that the container can use.

### spec.notifications.resources.limits.cpu

`string`

### spec.notifications.resources.limits.memory

`string`

### spec.notifications.resources.requests

`CpuMemory`

The resource requests for the container.
Specify the minimum amount of CPU and memory that the container is guaranteed.

### spec.notifications.resources.requests.cpu

`string`

### spec.notifications.resources.requests.memory

`string`

### spec.dex

`KubernetesArgocdToggleableComponent`

Dex — the bundled SSO broker for connector-based logins (GitHub,
LDAP, SAML, ...). Chart default: enabled. Only does something once
`sso.dex_config` declares connectors; disable when using direct
OIDC (`sso.oidc`) or admin-only access.

### spec.dex.enabled

`bool` · optional (explicit presence)

Enable the component. Empty = the chart default (notifications and
dex: enabled; commit server: disabled).

### spec.dex.resources

`ContainerResources`

CPU and memory for the component's container. Empty = no
requests/limits (the chart default).

### spec.dex.resources.limits

`CpuMemory`

The resource limits for the container.
Specify the maximum amount of CPU and memory that the container can use.

### spec.dex.resources.limits.cpu

`string`

### spec.dex.resources.limits.memory

`string`

### spec.dex.resources.requests

`CpuMemory`

The resource requests for the container.
Specify the minimum amount of CPU and memory that the container is guaranteed.

### spec.dex.resources.requests.cpu

`string`

### spec.dex.resources.requests.memory

`string`

### spec.commitServer

`KubernetesArgocdToggleableComponent`

The commit server — the manifest-hydration (rendered-manifests)
feature. Chart default: disabled.

### spec.commitServer.enabled

`bool` · optional (explicit presence)

Enable the component. Empty = the chart default (notifications and
dex: enabled; commit server: disabled).

### spec.commitServer.resources

`ContainerResources`

CPU and memory for the component's container. Empty = no
requests/limits (the chart default).

### spec.commitServer.resources.limits

`CpuMemory`

The resource limits for the container.
Specify the maximum amount of CPU and memory that the container can use.

### spec.commitServer.resources.limits.cpu

`string`

### spec.commitServer.resources.limits.memory

`string`

### spec.commitServer.resources.requests

`CpuMemory`

The resource requests for the container.
Specify the minimum amount of CPU and memory that the container is guaranteed.

### spec.commitServer.resources.requests.cpu

`string`

### spec.commitServer.resources.requests.memory

`string`

### spec.redis

`KubernetesArgocdRedis`

The Redis cache arm. Empty = the chart's bundled single-pod Redis —
right for most installs (the cache is disposable; losing it costs a
re-sync, not state).

### spec.redis.bundled

`KubernetesArgocdRedisBundled`

The chart's bundled single-pod Redis (the default when the whole
block is empty).

### spec.redis.bundled.resources

`ContainerResources`

CPU and memory for the Redis container. Empty = no requests/limits
(the chart default).

### spec.redis.bundled.resources.limits

`CpuMemory`

The resource limits for the container.
Specify the maximum amount of CPU and memory that the container can use.

### spec.redis.bundled.resources.limits.cpu

`string`

### spec.redis.bundled.resources.limits.memory

`string`

### spec.redis.bundled.resources.requests

`CpuMemory`

The resource requests for the container.
Specify the minimum amount of CPU and memory that the container is guaranteed.

### spec.redis.bundled.resources.requests.cpu

`string`

### spec.redis.bundled.resources.requests.memory

`string`

### spec.redis.ha

`KubernetesArgocdRedisHa`

The redis-ha subchart: a 3-replica Sentinel cluster behind
HAProxy. KNOW THIS: its pods carry a REQUIRED anti-affinity —
the cluster needs at least 3 schedulable worker nodes or the
pods stay Pending.

### spec.redis.ha.replicas

`int32` · optional (explicit presence)

Number of Redis replicas. Empty = 3 (the subchart default — also
the Sentinel quorum floor; do not go below it).

- default: `3`
- rule: {"int32":{"gte":3}}

### spec.redis.external

`KubernetesArgocdRedisExternal`

An external Redis-compatible endpoint (a managed cache or a
KubernetesValkey resource).

### spec.redis.external.host

`string | valueFrom` · required

Redis host (name or FQDN, no port). Accepts a literal host or a
reference to a KubernetesValkey resource (Valkey speaks the Redis
protocol).

- references: KubernetesValkey (`status.outputs.service`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesValkey, name: <that resource's name>, fieldPath: status.outputs.service}} -- a bare string does not parse

### spec.redis.external.port

`int32` · optional (explicit presence)

Redis port. Empty = 6379.

- default: `6379`
- rule: {"int32":{"lte":65535,"gte":1}}

### spec.redis.external.credentialsSecretName

`string`

Existing Secret holding the Redis credentials — key
`redis-password` required by the chart's contract; add
`redis-username` when the user is not `default`. Empty =
unauthenticated Redis.

### spec.sso

`KubernetesArgocdSso`

Single sign-on configuration (direct OIDC or Dex connectors).

### spec.sso.oidc

`KubernetesArgocdOidc`

Direct OIDC against the identity provider (no Dex in the path).

### spec.sso.oidc.name

`string` · required

Display name on the login button (e.g. "Okta", "Azure AD").

- rule: {"required":true}

### spec.sso.oidc.issuer

`string` · required

The OIDC issuer URL (e.g.
"https://login.microsoftonline.com/<tenant>/v2.0").

- rule: {"required":true}

### spec.sso.oidc.clientId

`string` · required

The OAuth client ID registered for Argo CD.

- rule: {"required":true}

### spec.sso.oidc.clientSecretSecret

`KubernetesArgocdSecretKeyRef`

The OAuth client secret, read at runtime from an existing Secret in
Argo CD's namespace. KNOW THIS: the referenced Secret must carry
the label `app.kubernetes.io/part-of: argocd` or Argo CD refuses
to read it. Empty = a public client (PKCE).

### spec.sso.oidc.clientSecretSecret.name

`string` · required

Secret name. The Secret must live in Argo CD's namespace.

- rule: {"required":true}

### spec.sso.oidc.clientSecretSecret.key

`string` · required

Key within the Secret.

- rule: {"required":true}

### spec.sso.dexConfig

`string`

Dex connector configuration as YAML (the `connectors:` document
from Dex's documentation). Requires the dex component enabled (the
default). SECRET INDIRECTION: never inline client secrets — write
`$<secret-name>:<key>` and Argo CD resolves it at runtime from a
Secret in its namespace labeled `app.kubernetes.io/part-of:
argocd` (compose that Secret from KubernetesSecret or
KubernetesExternalSecret).

### spec.rbac

`KubernetesArgocdRbac`

Argo CD's own RBAC (who may do what INSIDE Argo CD — distinct from
Kubernetes RBAC).

### spec.rbac.policyDefault

`string`

Role every authenticated user falls back to when no policy
matches (e.g. "role:readonly"). Empty = no fallback — users log in
and see nothing until a policy grants them access (the safe
default).

### spec.rbac.policyCsv

`string`

RBAC policy in Argo CD's CSV form — `p` rows grant, `g` rows bind
groups/users to roles (e.g.
"g, my-org:platform-team, role:admin").

### spec.rbac.scopes

`string` · optional (explicit presence)

OIDC token claims examined for group bindings. Empty = "[groups]".

- default: `[groups]`

### spec.execEnabled

`bool`

Allow `argocd exec` / the UI terminal into managed pods. Default
false — an interactive shell through the API server is a serious
capability; enable deliberately and pair with an RBAC `exec` rule.

### spec.reconciliationTimeout

`string`

How often Argo CD polls repositories for new commits, as a Go
duration (e.g. "120s", "3m"). Empty = the chart default (120s with
up to 60s jitter). Lower = faster syncs, more repo-server and Git
load; webhook-driven setups can raise it.

- rule: {"string":{"pattern":"^$|^\\d+(s|m|h)$"}}

### spec.repositories

`[]KubernetesArgocdRepository`

PUBLIC (anonymous) repositories registered declaratively. For
credentialed repositories see the REPOSITORY CREDENTIALS note on
the spec — credentials are composed as labeled Secrets, never
transported through this spec.

### spec.repositories[].name

`string` · required

Registration name (also the key of the rendered Secret; keep it
DNS-label-ish).

- rule: {"required":true}

### spec.repositories[].url

`string` · required

Repository URL (e.g. "https://github.com/org/repo" or an
"https://" Helm repository URL).

- rule: {"required":true}

### spec.repositories[].type

`string` · optional (explicit presence)

Repository type. Empty = "git"; set "helm" for Helm repositories.

- default: `git`
- rule: {"string":{"in":["git","helm"]}}

### spec.crds

`KubernetesArgocdCrds`

CRD lifecycle. The chart templates the Application/AppProject/
ApplicationSet CRDs and keeps them on uninstall by default —
deleting CRDs cascades to every Application in the cluster.

### spec.crds.install

`bool` · optional (explicit presence)

Install (and upgrade) the Application/AppProject/ApplicationSet
CRDs with the release. KNOW THIS: the chart TEMPLATES these CRDs,
so they carry the installing release's Helm ownership metadata —
and CRDs kept by a destroyed release still carry it. Disable
whenever the cluster already has the CRDs (from another live
release OR kept behind a destroyed one): a second release cannot
adopt them and its install fails on ownership validation.

- default: `true`

### spec.crds.keep

`bool` · optional (explicit presence)

Keep the CRDs when this resource is destroyed (the chart stamps
`helm.sh/resource-policy: keep`). Default true — REMOVING the CRDs
DELETES EVERY Application, AppProject and ApplicationSet in the
cluster with them. Turn off only on throwaway clusters.

- default: `true`

### spec.serviceMonitorsEnabled

`bool`

Create ServiceMonitors for every component's /metrics (requires the
Prometheus Operator CRDs — deploy KubernetesKubePrometheusStack
first). Chart default: false.

### spec.image

`KubernetesArgocdImage`

Override the Argo CD image (air-gap path). One image serves all
Argo CD components; Dex and Redis images stay chart-defaults
(override via `helm_values` when mirroring those too).

### spec.image.repository

`string`

Image repository including registry, e.g.
"my.registry.com/argoproj/argocd". Empty = "quay.io/argoproj/argocd".

### spec.image.tag

`string`

Image tag. Empty = the chart's appVersion for the pinned
chart_version.

### spec.image.pullSecretName

`string`

Name of an existing image-pull Secret in the namespace, for
private mirrors (applied to all components).

### spec.scheduling

`KubernetesArgocdScheduling`

Scheduling applied to ALL components (the chart's global block).

### spec.scheduling.nodeSelector

`map<string, string>`

Node selector for all Argo CD pods.

### spec.scheduling.tolerations

`[]WorkloadToleration`

Tolerations for all Argo CD pods.

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

Priority class name for all Argo CD pods.

### spec.helmValues

`string`

Escape hatch: additional chart values as a YAML document, merged
LAST over everything the typed fields render (Helm `-f` semantics,
identical on both engines). For the chart surface beyond the typed
fields (per-component env/volumes, notifications notifiers and
templates, server extensions, network policies, ...) — never the
substitute for them. Do not put secrets here; credential material
belongs in referenced Secrets.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesArgocd, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.namespace` | `string` | Namespace Argo CD is installed in. |
| `status.outputs.release_name` | `string` | Helm release name (equals metadata.name). |
| `status.outputs.server_service` | `string` | Name of the API/UI server Service (`<name>-server`) — the backend handle exposure kinds (KubernetesIngress, KubernetesHttpRoute) reference. |
| `status.outputs.server_kube_endpoint` | `string` | In-cluster endpoint of the server (e.g. "https://main-argocd-server.gitops.svc.cluster.local"). HTTPS by default; plain HTTP when `server.insecure` is set. |
| `status.outputs.initial_admin_secret_name` | `string` | Name of the Secret carrying the generated initial admin password (key `password`) — Argo CD creates it at first start. Empty when the admin user is disabled. The name is fixed by the application: "argocd-initial-admin-secret". |
| `status.outputs.port_forward_command` | `string` | Command to port-forward the UI to a workstation (https://localhost:8080 unless `server.insecure`). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespace` | KubernetesNamespace | `spec.name` |
| `spec.redis.external.host` | KubernetesValkey | `status.outputs.service` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
