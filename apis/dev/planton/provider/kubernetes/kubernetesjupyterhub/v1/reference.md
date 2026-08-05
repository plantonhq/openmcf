# KubernetesJupyterHub

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1`

**KubernetesJupyterHubSpec** deploys JupyterHub — the multi-user
notebook platform — from the official `jupyterhub` Helm chart
(https://hub.jupyter.org/helm-chart/, the Zero to JupyterHub
distribution).

WHAT GETS INSTALLED: the hub (JupyterHub + KubeSpawner + the
configured authenticator), the configurable-http-proxy pod routing
every user request, the user-scheduler (packs user pods onto busy
nodes so autoscalers can reclaim the rest), optional placeholder
pods that pre-warm capacity, image-puller DaemonSets that pre-pull
the notebook image onto every node, and the idle-culler that stops
abandoned notebook servers.

HOW IT WORKS: every user gets their OWN notebook server — the hub
spawns a dedicated pod (and, by default, a dedicated persistent
home volume) per user at login, and the proxy routes each user to
their server. These per-user pods and volumes are created by the
hub AT RUNTIME — they are not part of this resource's deploy and
they survive it: uninstalling the release deletes the hub, proxy
and scheduling machinery, but user home PVCs (named `claim-<user>`)
remain until deleted explicitly. Treat user homes as data.

NAMES ARE CHART-FIXED: resources render with bare names (hub,
proxy, proxy-public, user-scheduler…), so one JupyterHub per
namespace — deploy each instance into its own namespace (the
default posture here).

SECURED BY DEFAULT: the chart's own default authenticator accepts
ANY username with NO password. This kind never ships that — when
`authentication` is empty, the shared sign-in password is
module-generated into `<name>-auth` and wired to the hub through an
environment variable, never through rendered Helm values. The
chart's three internal auth materials (proxy token, cookie secret,
crypt keys) are chart-managed and stable across upgrades.

STATE lives in the hub database: sqlite on a small PVC (the
default — fine for most installs, the hub is single-replica by
design) or an external PostgreSQL (a KubernetesPostgres composes
naturally; the password rides a mounted Secret, never rendered
values).

EXPOSURE: the proxy-public Service stays ClusterIP here (the chart
default is LoadBalancer — deliberately overridden); expose it via
first-class kinds (KubernetesService, Gateway API routes,
KubernetesIngress) over the exported service handle. The chart's
ingress/httpRoute blocks and the autohttps/ACME machinery are never
modeled — TLS terminates at the composed exposure layer.

The typed fields below cover the chart's meaningful configuration
surface; `helm_values` remains as the escape hatch (merged last,
Helm `-f` semantics, identical on both engines) for the long tail —
per-pod securityContext overrides, network-policy egress shaping,
hub.extraConfig snippets, services/loadRoles, LDAP auth — a safety
valve, never the primary interface.

## Example

```yaml
# Full-surface shape: the external-PostgreSQL hub database, OIDC sign-in
# (Keycloak-shaped), sized user servers with a spawn-menu, the packing
# scheduler with warm placeholders, tuned culling and an escape-hatch
# entry — the offline plan/preview proof for the widest typed rendering.
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesJupyterHub
metadata:
  name: jupyterhub-full
spec:
  namespace:
    value: jupyterhub-full
  createNamespace: true
  chartVersion: 4.4.0
  hub:
    database:
      postgres:
        host:
          value: hub-pg-rw
        port: 5432
        databaseName: jupyterhub
        username: jupyterhub
        passwordSecret:
          secretName:
            value: hub-pg-app
          secretKey: password
    concurrentSpawnLimit: 32
    activeServerLimit: 200
    allowNamedServers: true
    namedServerLimitPerUser: 2
    shutdownOnLogout: true
    resources:
      requests:
        cpu: 200m
        memory: 512Mi
      limits:
        memory: 1Gi
  authentication:
    oidc:
      clientId: jupyterhub
      clientSecretSecret:
        secretName: kc-oauth
        secretKey: client-secret
      oauthCallbackUrl: https://hub.example.com/hub/oauth_callback
      authorizeUrl: https://kc.example.com/realms/eng/protocol/openid-connect/auth
      tokenUrl: https://kc.example.com/realms/eng/protocol/openid-connect/token
      userdataUrl: https://kc.example.com/realms/eng/protocol/openid-connect/userinfo
      scopes:
        - openid
        - email
      usernameClaim: email
      loginService: Keycloak
    adminUsers:
      - ada
  proxy:
    serviceType: ClusterIP
    resources:
      requests:
        cpu: 100m
        memory: 128Mi
  singleUser:
    image:
      repository: quay.io/jupyter/scipy-notebook
      tag: "2026-07-28"
    memoryGuarantee: 2G
    memoryLimit: 4G
    cpuGuarantee: "0.5"
    cpuLimit: "2"
    defaultUrl: /lab
    startTimeoutSeconds: 600
    storage:
      dynamic:
        capacity: 20Gi
    profiles:
      - displayName: Small
        description: Default 2G workspace
        default: true
        memoryGuarantee: 2G
        memoryLimit: 2G
      - displayName: Large
        description: Big-memory ETL
        memoryLimit: 8G
        cpuLimit: "4"
  scheduling:
    userSchedulerEnabled: true
    userPlaceholderReplicas: 2
  culling:
    enabled: true
    timeoutSeconds: 1800
    everySeconds: 300
  prePuller:
    hookEnabled: true
    continuousEnabled: true
  helmValues: |
    debug:
      enabled: true
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.namespace` | `string \| valueFrom` | yes |  | KubernetesNamespace (`spec.name`) |
| `spec.createNamespace` | `bool` |  |  |  |
| `spec.chartVersion` | `string` |  | `4.4.0` |  |
| `spec.hub` | `KubernetesJupyterHubHub` |  |  |  |
| `spec.hub.database` | `KubernetesJupyterHubDatabase` |  |  |  |
| `spec.hub.database.sqlitePvc` | `KubernetesJupyterHubSqlitePvc` |  |  |  |
| `spec.hub.database.sqlitePvc.storageSize` | `string` |  | `1Gi` |  |
| `spec.hub.database.sqlitePvc.storageClass` | `string` |  |  |  |
| `spec.hub.database.postgres` | `KubernetesJupyterHubPostgres` |  |  |  |
| `spec.hub.database.postgres.host` | `string \| valueFrom` | yes |  | KubernetesPostgres (`status.outputs.rw_service`) |
| `spec.hub.database.postgres.port` | `int32` |  | `5432` |  |
| `spec.hub.database.postgres.databaseName` | `string` |  | `jupyterhub` |  |
| `spec.hub.database.postgres.username` | `string` |  | `jupyterhub` |  |
| `spec.hub.database.postgres.passwordSecret` | `KubernetesJupyterHubPasswordSecret` | yes |  |  |
| `spec.hub.database.postgres.passwordSecret.secretName` | `string \| valueFrom` | yes |  | KubernetesPostgres (`status.outputs.password_secret.name`) |
| `spec.hub.database.postgres.passwordSecret.secretKey` | `string` |  | `password` |  |
| `spec.hub.database.mysql` | `KubernetesJupyterHubMysql` |  |  |  |
| `spec.hub.database.mysql.host` | `string \| valueFrom` | yes |  | KubernetesMysql (`status.outputs.primary_service`) |
| `spec.hub.database.mysql.port` | `int32` |  | `3306` |  |
| `spec.hub.database.mysql.databaseName` | `string` |  | `jupyterhub` |  |
| `spec.hub.database.mysql.username` | `string` |  | `jupyterhub` |  |
| `spec.hub.database.mysql.passwordSecret` | `KubernetesJupyterHubMysqlPasswordSecret` | yes |  |  |
| `spec.hub.database.mysql.passwordSecret.secretName` | `string \| valueFrom` | yes |  | KubernetesMysql (`status.outputs.root_password_secret.name`) |
| `spec.hub.database.mysql.passwordSecret.secretKey` | `string` |  | `password` |  |
| `spec.hub.concurrentSpawnLimit` | `int32` |  | `64` |  |
| `spec.hub.activeServerLimit` | `int32` |  |  |  |
| `spec.hub.allowNamedServers` | `bool` |  |  |  |
| `spec.hub.namedServerLimitPerUser` | `int32` |  |  |  |
| `spec.hub.shutdownOnLogout` | `bool` |  |  |  |
| `spec.hub.resources` | `ContainerResources` |  |  |  |
| `spec.hub.resources.limits` | `CpuMemory` |  |  |  |
| `spec.hub.resources.limits.cpu` | `string` |  |  |  |
| `spec.hub.resources.limits.memory` | `string` |  |  |  |
| `spec.hub.resources.requests` | `CpuMemory` |  |  |  |
| `spec.hub.resources.requests.cpu` | `string` |  |  |  |
| `spec.hub.resources.requests.memory` | `string` |  |  |  |
| `spec.authentication` | `KubernetesJupyterHubAuth` |  |  |  |
| `spec.authentication.sharedPassword` | `KubernetesJupyterHubDummyAuth` |  |  |  |
| `spec.authentication.sharedPassword.passwordSecret` | `KubernetesJupyterHubExistingSecretRef` |  |  |  |
| `spec.authentication.sharedPassword.passwordSecret.secretName` | `string` | yes |  |  |
| `spec.authentication.sharedPassword.passwordSecret.secretKey` | `string` |  | `password` |  |
| `spec.authentication.native` | `KubernetesJupyterHubNativeAuth` |  |  |  |
| `spec.authentication.native.openSignup` | `bool` |  |  |  |
| `spec.authentication.native.minimumPasswordLength` | `int32` |  | `8` |  |
| `spec.authentication.github` | `KubernetesJupyterHubGithubAuth` |  |  |  |
| `spec.authentication.github.clientId` | `string` | yes |  |  |
| `spec.authentication.github.clientSecretSecret` | `KubernetesJupyterHubExistingSecretRef` | yes |  |  |
| `spec.authentication.github.clientSecretSecret.secretName` | `string` | yes |  |  |
| `spec.authentication.github.clientSecretSecret.secretKey` | `string` |  | `password` |  |
| `spec.authentication.github.oauthCallbackUrl` | `string` | yes |  |  |
| `spec.authentication.github.allowedOrganizations` | `[]string` |  |  |  |
| `spec.authentication.google` | `KubernetesJupyterHubGoogleAuth` |  |  |  |
| `spec.authentication.google.clientId` | `string` | yes |  |  |
| `spec.authentication.google.clientSecretSecret` | `KubernetesJupyterHubExistingSecretRef` | yes |  |  |
| `spec.authentication.google.clientSecretSecret.secretName` | `string` | yes |  |  |
| `spec.authentication.google.clientSecretSecret.secretKey` | `string` |  | `password` |  |
| `spec.authentication.google.oauthCallbackUrl` | `string` | yes |  |  |
| `spec.authentication.google.hostedDomains` | `[]string` |  |  |  |
| `spec.authentication.oidc` | `KubernetesJupyterHubOidcAuth` |  |  |  |
| `spec.authentication.oidc.clientId` | `string` | yes |  |  |
| `spec.authentication.oidc.clientSecretSecret` | `KubernetesJupyterHubExistingSecretRef` | yes |  |  |
| `spec.authentication.oidc.clientSecretSecret.secretName` | `string` | yes |  |  |
| `spec.authentication.oidc.clientSecretSecret.secretKey` | `string` |  | `password` |  |
| `spec.authentication.oidc.oauthCallbackUrl` | `string` | yes |  |  |
| `spec.authentication.oidc.authorizeUrl` | `string` | yes |  |  |
| `spec.authentication.oidc.tokenUrl` | `string` | yes |  |  |
| `spec.authentication.oidc.userdataUrl` | `string` | yes |  |  |
| `spec.authentication.oidc.scopes` | `[]string` |  |  |  |
| `spec.authentication.oidc.usernameClaim` | `string` |  | `preferred_username` |  |
| `spec.authentication.oidc.loginService` | `string` |  | `OIDC` |  |
| `spec.authentication.adminUsers` | `[]string` |  |  |  |
| `spec.authentication.allowedUsers` | `[]string` |  |  |  |
| `spec.proxy` | `KubernetesJupyterHubProxy` |  |  |  |
| `spec.proxy.serviceType` | `string` |  | `ClusterIP` |  |
| `spec.proxy.serviceAnnotations` | `map<string, string>` |  |  |  |
| `spec.proxy.resources` | `ContainerResources` |  |  |  |
| `spec.proxy.resources.limits` | `CpuMemory` |  |  |  |
| `spec.proxy.resources.limits.cpu` | `string` |  |  |  |
| `spec.proxy.resources.limits.memory` | `string` |  |  |  |
| `spec.proxy.resources.requests` | `CpuMemory` |  |  |  |
| `spec.proxy.resources.requests.cpu` | `string` |  |  |  |
| `spec.proxy.resources.requests.memory` | `string` |  |  |  |
| `spec.singleUser` | `KubernetesJupyterHubSingleUser` |  |  |  |
| `spec.singleUser.image` | `KubernetesJupyterHubImage` |  |  |  |
| `spec.singleUser.image.repository` | `string` | yes |  |  |
| `spec.singleUser.image.tag` | `string` | yes |  |  |
| `spec.singleUser.memoryGuarantee` | `string` |  | `1G` |  |
| `spec.singleUser.memoryLimit` | `string` |  |  |  |
| `spec.singleUser.cpuGuarantee` | `string` |  |  |  |
| `spec.singleUser.cpuLimit` | `string` |  |  |  |
| `spec.singleUser.storage` | `KubernetesJupyterHubUserStorage` |  |  |  |
| `spec.singleUser.storage.dynamic` | `KubernetesJupyterHubDynamicStorage` |  |  |  |
| `spec.singleUser.storage.dynamic.capacity` | `string` |  | `10Gi` |  |
| `spec.singleUser.storage.dynamic.storageClass` | `string` |  |  |  |
| `spec.singleUser.storage.static` | `KubernetesJupyterHubStaticStorage` |  |  |  |
| `spec.singleUser.storage.static.pvcName` | `string` | yes |  |  |
| `spec.singleUser.storage.static.subPath` | `string` |  | `{username}` |  |
| `spec.singleUser.storage.none` | `KubernetesJupyterHubNoStorage` |  |  |  |
| `spec.singleUser.defaultUrl` | `string` |  |  |  |
| `spec.singleUser.startTimeoutSeconds` | `int32` |  | `300` |  |
| `spec.singleUser.extraEnv` | `map<string, string>` |  |  |  |
| `spec.singleUser.profiles` | `[]KubernetesJupyterHubProfile` |  |  |  |
| `spec.singleUser.profiles[].displayName` | `string` | yes |  |  |
| `spec.singleUser.profiles[].description` | `string` |  |  |  |
| `spec.singleUser.profiles[].default` | `bool` |  |  |  |
| `spec.singleUser.profiles[].image` | `KubernetesJupyterHubImage` |  |  |  |
| `spec.singleUser.profiles[].image.repository` | `string` | yes |  |  |
| `spec.singleUser.profiles[].image.tag` | `string` | yes |  |  |
| `spec.singleUser.profiles[].memoryGuarantee` | `string` |  |  |  |
| `spec.singleUser.profiles[].memoryLimit` | `string` |  |  |  |
| `spec.singleUser.profiles[].cpuGuarantee` | `string` |  |  |  |
| `spec.singleUser.profiles[].cpuLimit` | `string` |  |  |  |
| `spec.scheduling` | `KubernetesJupyterHubScheduling` |  |  |  |
| `spec.scheduling.userSchedulerEnabled` | `bool` |  | `true` |  |
| `spec.scheduling.userPlaceholderReplicas` | `int32` |  | `0` |  |
| `spec.scheduling.coreNodeSelector` | `map<string, string>` |  |  |  |
| `spec.scheduling.userNodeSelector` | `map<string, string>` |  |  |  |
| `spec.culling` | `KubernetesJupyterHubCulling` |  |  |  |
| `spec.culling.enabled` | `bool` |  | `true` |  |
| `spec.culling.timeoutSeconds` | `int32` |  | `3600` |  |
| `spec.culling.everySeconds` | `int32` |  | `600` |  |
| `spec.culling.maxAgeSeconds` | `int32` |  |  |  |
| `spec.culling.cullUsers` | `bool` |  |  |  |
| `spec.prePuller` | `KubernetesJupyterHubPrePuller` |  |  |  |
| `spec.prePuller.hookEnabled` | `bool` |  | `true` |  |
| `spec.prePuller.continuousEnabled` | `bool` |  | `true` |  |
| `spec.networkPolicyEnabled` | `bool` |  | `true` |  |
| `spec.helmValues` | `string` |  |  |  |

## Field Details

### spec.namespace

`string | valueFrom` · required

Namespace to install into. Accepts a literal namespace name or a
reference to a KubernetesNamespace resource. One JupyterHub per
namespace — the chart's resource names are fixed (hub, proxy,
proxy-public…), so a second release in the same namespace
collides.

- references: KubernetesNamespace (`spec.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesNamespace, name: <that resource's name>, fieldPath: spec.name}} -- a bare string does not parse

### spec.createNamespace

`bool`

When true, the namespace is created (with the standard Planton
governance labels) before installing and deleted with the
resource. When false, the namespace must already exist. KNOW
THIS: deleting the namespace also deletes every user's home PVC
living in it — back user homes up before tearing an instance
down.

### spec.chartVersion

`string` · optional (explicit presence)

Helm chart version to install (e.g. "4.4.0" — chart 4.4.0 ships
JupyterHub 5.5.0). Versions must exist as SERVED charts in the
repository index (https://hub.jupyter.org/helm-chart/). The chart
requires Kubernetes 1.28+.

- default: `4.4.0`

### spec.hub

`KubernetesJupyterHubHub`

The hub — JupyterHub itself: its database, spawn throttles and
server-lifecycle switches.

- rule: named_server_limit_per_user only applies when allow_named_servers is true — enable named servers or remove the limit.

### spec.hub.database

`KubernetesJupyterHubDatabase`

The hub database — users, running-server records, tokens.
Empty = sqlite on a 1Gi PVC (the chart default; right for most
installs).

### spec.hub.database.sqlitePvc

`KubernetesJupyterHubSqlitePvc`

sqlite on a dedicated PVC (the chart default) — zero external
dependencies, right for most installs since the hub is
single-replica anyway. The PVC (named `hub-db-dir`) carries
ALL hub state.

### spec.hub.database.sqlitePvc.storageSize

`string` · optional (explicit presence)

PVC size. Empty = "1Gi" (the chart default — hub state is
small).

- default: `1Gi`
- rule: {"string":{"pattern":"^[0-9]+(\\.[0-9]+)?(Ei|Pi|Ti|Gi|Mi|Ki|E|P|T|G|M|k|m)?$"}}

### spec.hub.database.sqlitePvc.storageClass

`string`

Storage class for the hub-db PVC. Empty = the cluster default.

### spec.hub.database.postgres

`KubernetesJupyterHubPostgres`

External PostgreSQL — a KubernetesPostgres composes naturally.
Choose this to survive PVC loss, to snapshot hub state with
your database fleet, or on clusters without dynamic volumes.

### spec.hub.database.postgres.host

`string | valueFrom` · required

Database server host. Defaults compose a KubernetesPostgres
resource's read-write Service; any reachable PostgreSQL works as
a literal value.

- references: KubernetesPostgres (`status.outputs.rw_service`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesPostgres, name: <that resource's name>, fieldPath: status.outputs.rw_service}} -- a bare string does not parse

### spec.hub.database.postgres.port

`int32` · optional (explicit presence)

Database server port. Empty = 5432.

- default: `5432`
- rule: {"int32":{"lte":65535,"gte":1}}

### spec.hub.database.postgres.databaseName

`string` · optional (explicit presence)

Database name holding hub state. Must EXIST before install (the
hub creates tables, never the database) — on a
KubernetesPostgres, declare it at bootstrap (`initdb.database`).
Empty = "jupyterhub".

- default: `jupyterhub`
- rule: {"string":{"pattern":"^[a-zA-Z_][a-zA-Z0-9_$]*$"}}

### spec.hub.database.postgres.username

`string` · optional (explicit presence)

Database user with full rights inside `database_name` (ownership
is simplest — declare the same user as the database owner at
bootstrap). Empty = "jupyterhub".

- default: `jupyterhub`

### spec.hub.database.postgres.passwordSecret

`KubernetesJupyterHubPasswordSecret` · required

The Secret holding the user's password. The hub mounts it
through the chart's existing-secret seam — it never lands in
rendered values. Defaults compose a KubernetesPostgres
resource's application-user Secret. Same-namespace constraint (a
Kubernetes rule): the Secret must live in the namespace
JupyterHub installs into.

- rule: {"required":true}

### spec.hub.database.postgres.passwordSecret.secretName

`string | valueFrom` · required

Name of the Secret. Defaults compose a KubernetesPostgres
resource's application-user Secret (`<cluster>-app`).
Same-namespace constraint applies.

- references: KubernetesPostgres (`status.outputs.password_secret.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesPostgres, name: <that resource's name>, fieldPath: status.outputs.password_secret.name}} -- a bare string does not parse

### spec.hub.database.postgres.passwordSecret.secretKey

`string` · optional (explicit presence)

Key inside the Secret holding the password. Empty = "password"
(the KubernetesPostgres application-Secret convention).

- default: `password`

### spec.hub.database.mysql

`KubernetesJupyterHubMysql`

External MySQL 8+ — a KubernetesMysql composes naturally.

### spec.hub.database.mysql.host

`string | valueFrom` · required

Database server host. Defaults compose a KubernetesMysql
resource's client Service.

- references: KubernetesMysql (`status.outputs.primary_service`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesMysql, name: <that resource's name>, fieldPath: status.outputs.primary_service}} -- a bare string does not parse

### spec.hub.database.mysql.port

`int32` · optional (explicit presence)

Database server port. Empty = 3306.

- default: `3306`
- rule: {"int32":{"lte":65535,"gte":1}}

### spec.hub.database.mysql.databaseName

`string` · optional (explicit presence)

Database name holding hub state. Must EXIST before install.
Empty = "jupyterhub".

- default: `jupyterhub`
- rule: {"string":{"pattern":"^[a-zA-Z_][a-zA-Z0-9_$]*$"}}

### spec.hub.database.mysql.username

`string` · optional (explicit presence)

Database user with full rights inside `database_name`. Empty =
"jupyterhub".

- default: `jupyterhub`

### spec.hub.database.mysql.passwordSecret

`KubernetesJupyterHubMysqlPasswordSecret` · required

The Secret holding the user's password. Defaults compose a
KubernetesMysql resource's root credential Secret; for a
dedicated hub user (recommended), point at that user's Secret.
Same-namespace constraint applies.

- rule: {"required":true}

### spec.hub.database.mysql.passwordSecret.secretName

`string | valueFrom` · required

Name of the Secret. Defaults compose a KubernetesMysql
resource's root credential Secret. Same-namespace constraint
applies.

- references: KubernetesMysql (`status.outputs.root_password_secret.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesMysql, name: <that resource's name>, fieldPath: status.outputs.root_password_secret.name}} -- a bare string does not parse

### spec.hub.database.mysql.passwordSecret.secretKey

`string` · optional (explicit presence)

Key inside the Secret holding the password. Empty = "password".

- default: `password`

### spec.hub.concurrentSpawnLimit

`int32` · optional (explicit presence)

How many users may be SPAWNING simultaneously (not how many may
be running) — the login-storm throttle. Empty = 64 (the chart
default).

- default: `64`
- rule: {"int32":{"gte":1}}

### spec.hub.activeServerLimit

`int32` · optional (explicit presence)

Hard cap on RUNNING notebook servers across all users — the
cluster-capacity guard. Users beyond it see "server limit
reached" at spawn. Empty = unlimited (the chart default).

- rule: {"int32":{"gte":1}}

### spec.hub.allowNamedServers

`bool`

Let each user run multiple NAMED servers (e.g. one per project)
alongside their default server. Off by default — each named
server is a full pod + volume.

### spec.hub.namedServerLimitPerUser

`int32` · optional (explicit presence)

Cap on named servers per user. Only meaningful with
`allow_named_servers`. Empty = unlimited.

- rule: {"int32":{"gte":1}}

### spec.hub.shutdownOnLogout

`bool`

Stop the user's notebook server when they log out. Off by
default (matching the chart) — the culler usually handles
abandonment better than logout does.

### spec.hub.resources

`ContainerResources`

CPU/memory for the hub container. Empty = the chart's defaults
(no requests). The hub is light — 512Mi/0.5cpu carries hundreds
of users.

### spec.hub.resources.limits

`CpuMemory`

The resource limits for the container.
Specify the maximum amount of CPU and memory that the container can use.

### spec.hub.resources.limits.cpu

`string`

### spec.hub.resources.limits.memory

`string`

### spec.hub.resources.requests

`CpuMemory`

The resource requests for the container.
Specify the minimum amount of CPU and memory that the container is guaranteed.

### spec.hub.resources.requests.cpu

`string`

### spec.hub.resources.requests.memory

`string`

### spec.authentication

`KubernetesJupyterHubAuth`

WHO can sign in and HOW. Empty = username/shared-password sign-in
with a module-generated password (`<name>-auth`, key `password`)
— the chart's own default (any username, NO password) never
ships.

### spec.authentication.sharedPassword

`KubernetesJupyterHubDummyAuth`

Shared-password sign-in (upstream's DummyAuthenticator +
password): ANY username, ONE shared password — teams, labs,
evaluations. The username becomes the user's identity (home
volume, admin rights), so users must pick consistent names.

### spec.authentication.sharedPassword.passwordSecret

`KubernetesJupyterHubExistingSecretRef`

Existing Secret holding the shared password. Empty = the module
GENERATES it into `<name>-auth` (key `password`) — stable across
upgrades. Either way the password reaches the hub as an env var;
it never renders into values.

### spec.authentication.sharedPassword.passwordSecret.secretName

`string` · required

Name of the Secret. Same-namespace constraint applies.

- rule: {"required":true}

### spec.authentication.sharedPassword.passwordSecret.secretKey

`string` · optional (explicit presence)

Key inside the Secret. Empty = "password".

- default: `password`

### spec.authentication.native

`KubernetesJupyterHubNativeAuth`

Self-service accounts (NativeAuthenticator, ships in the hub
image): users sign UP with a username/password; admins approve
them in the hub UI before first login (unless open_signup).
No external identity provider needed.

### spec.authentication.native.openSignup

`bool`

Let new signups log in immediately WITHOUT admin approval. Off
by default — with it on, anyone reaching the login page can
create a working account, so pair it with `allowed_users` or a
private network.

### spec.authentication.native.minimumPasswordLength

`int32` · optional (explicit presence)

Minimum signup password length. Empty = 8.

- default: `8`
- rule: {"int32":{"gte":1}}

### spec.authentication.github

`KubernetesJupyterHubGithubAuth`

Sign in with GitHub (OAuthenticator). Gate by organization
membership via `allowed_organizations`.

### spec.authentication.github.clientId

`string` · required

The OAuth app's client ID (public, renders into values).

- rule: {"required":true}

### spec.authentication.github.clientSecretSecret

`KubernetesJupyterHubExistingSecretRef` · required

The Secret holding the OAuth app's client secret. Reaches the
hub as an env var (the authenticator's native environment
fallback) — never rendered into values. Same-namespace
constraint applies.

- rule: {"required":true}

### spec.authentication.github.clientSecretSecret.secretName

`string` · required

Name of the Secret. Same-namespace constraint applies.

- rule: {"required":true}

### spec.authentication.github.clientSecretSecret.secretKey

`string` · optional (explicit presence)

Key inside the Secret. Empty = "password".

- default: `password`

### spec.authentication.github.oauthCallbackUrl

`string` · required

The OAuth callback URL as registered with the provider —
`<your-jupyterhub-url>/hub/oauth_callback`. Must match the
registration exactly (scheme and host included).

- rule: {"required":true}

### spec.authentication.github.allowedOrganizations

`[]string`

GitHub organizations whose members may sign in (org name, or
`org:team` to gate by team). Empty = any GitHub identity —
effectively public; set this in production.

### spec.authentication.google

`KubernetesJupyterHubGoogleAuth`

Sign in with Google (OAuthenticator). Gate by workspace domain
via `hosted_domains`.

### spec.authentication.google.clientId

`string` · required

The OAuth client ID (public, renders into values).

- rule: {"required":true}

### spec.authentication.google.clientSecretSecret

`KubernetesJupyterHubExistingSecretRef` · required

The Secret holding the OAuth client secret (env-var indirection,
never rendered). Same-namespace constraint applies.

- rule: {"required":true}

### spec.authentication.google.clientSecretSecret.secretName

`string` · required

Name of the Secret. Same-namespace constraint applies.

- rule: {"required":true}

### spec.authentication.google.clientSecretSecret.secretKey

`string` · optional (explicit presence)

Key inside the Secret. Empty = "password".

- default: `password`

### spec.authentication.google.oauthCallbackUrl

`string` · required

The OAuth callback URL as registered —
`<your-jupyterhub-url>/hub/oauth_callback`.

- rule: {"required":true}

### spec.authentication.google.hostedDomains

`[]string`

Google Workspace domains whose users may sign in (e.g.
"example.com"). Empty = any Google account — set this in
production.

### spec.authentication.oidc

`KubernetesJupyterHubOidcAuth`

Sign in with any OIDC/OAuth2 provider (GenericOAuthenticator)
— Keycloak composes naturally (a KubernetesKeycloak realm's
endpoints slot straight in), as do Okta, Auth0, Dex.

### spec.authentication.oidc.clientId

`string` · required

The OAuth client ID (public, renders into values).

- rule: {"required":true}

### spec.authentication.oidc.clientSecretSecret

`KubernetesJupyterHubExistingSecretRef` · required

The Secret holding the OAuth client secret (env-var indirection,
never rendered). Same-namespace constraint applies.

- rule: {"required":true}

### spec.authentication.oidc.clientSecretSecret.secretName

`string` · required

Name of the Secret. Same-namespace constraint applies.

- rule: {"required":true}

### spec.authentication.oidc.clientSecretSecret.secretKey

`string` · optional (explicit presence)

Key inside the Secret. Empty = "password".

- default: `password`

### spec.authentication.oidc.oauthCallbackUrl

`string` · required

The OAuth callback URL as registered with the provider —
`<your-jupyterhub-url>/hub/oauth_callback`.

- rule: {"required":true}

### spec.authentication.oidc.authorizeUrl

`string` · required

The provider's authorization endpoint (e.g.
`https://keycloak.example.com/realms/eng/protocol/openid-connect/auth`).

- rule: {"required":true}

### spec.authentication.oidc.tokenUrl

`string` · required

The provider's token endpoint.

- rule: {"required":true}

### spec.authentication.oidc.userdataUrl

`string` · required

The provider's userinfo endpoint (where the authenticator reads
the user's identity claims).

- rule: {"required":true}

### spec.authentication.oidc.scopes

`[]string`

OAuth scopes to request. Empty = ["openid", "email", "profile"].

### spec.authentication.oidc.usernameClaim

`string` · optional (explicit presence)

The claim in the userinfo response that becomes the JupyterHub
username. Empty = "preferred_username" (the OIDC standard
claim).

- default: `preferred_username`

### spec.authentication.oidc.loginService

`string` · optional (explicit presence)

The provider name shown on the login button ("Sign in with
<login_service>"). Empty = "OIDC".

- default: `OIDC`

### spec.authentication.adminUsers

`[]string`

Usernames granted JupyterHub admin rights (manage users, access
servers, use the admin panel). For OAuth methods this is the
provider-side username (e.g. the GitHub login).

### spec.authentication.allowedUsers

`[]string`

Explicit allow-list of usernames permitted to sign in. Empty =
any authenticated identity may sign in (for shared_password that
means ANY username knowing the shared password — set this list
to pin the roster).

### spec.proxy

`KubernetesJupyterHubProxy`

The proxy tier — how traffic reaches JupyterHub and what the CHP
pod gets to work with.

### spec.proxy.serviceType

`string` · optional (explicit presence)

The proxy-public Service — the instance's front door. This kind
defaults it to ClusterIP (the chart's own default is
LoadBalancer — deliberately overridden): expose it by composing
first-class kinds over the exported service handle, or set
"LoadBalancer"/"NodePort" here for direct exposure.

- default: `ClusterIP`
- rule: {"string":{"in":["ClusterIP","LoadBalancer","NodePort"]}}

### spec.proxy.serviceAnnotations

`map<string, string>`

Annotations for the proxy-public Service — the cloud
load-balancer configuration surface (NLB type, internal LB,
static IPs) when `service_type` is LoadBalancer.

### spec.proxy.resources

`ContainerResources`

CPU/memory for the CHP container. Empty = the chart's defaults
(no requests).

### spec.proxy.resources.limits

`CpuMemory`

The resource limits for the container.
Specify the maximum amount of CPU and memory that the container can use.

### spec.proxy.resources.limits.cpu

`string`

### spec.proxy.resources.limits.memory

`string`

### spec.proxy.resources.requests

`CpuMemory`

The resource requests for the container.
Specify the minimum amount of CPU and memory that the container is guaranteed.

### spec.proxy.resources.requests.cpu

`string`

### spec.proxy.resources.requests.memory

`string`

### spec.singleUser

`KubernetesJupyterHubSingleUser`

The per-user notebook servers — image, sizing, home storage and
the profile menu. This is what your USERS live in; size it for
them, not for the hub.

### spec.singleUser.image

`KubernetesJupyterHubImage`

The notebook image every user server runs. Empty = the chart's
sample image (quay.io/jupyterhub/k8s-singleuser-sample — minimal
JupyterLab; fine for evaluation, replace it with a
jupyter/docker-stacks image or your own for real work). The
image must carry `jupyterhub-singleuser` matching the hub's
major version.

### spec.singleUser.image.repository

`string` · required

Image repository including any registry host (e.g.
"quay.io/jupyter/scipy-notebook").

- rule: {"required":true}

### spec.singleUser.image.tag

`string` · required

Image tag (e.g. "2026-07-28").

- rule: {"required":true}

### spec.singleUser.memoryGuarantee

`string` · optional (explicit presence)

Per-user memory GUARANTEE (the pod's request — capacity planning
multiplies this by concurrent users). Empty = "1G" (the chart
default). Kubernetes quantity or decimal bytes.

- default: `1G`

### spec.singleUser.memoryLimit

`string`

Per-user memory LIMIT (the OOM ceiling — a runaway notebook dies
at it instead of taking the node down). Empty = unlimited (the
chart default; set one in production).

### spec.singleUser.cpuGuarantee

`string`

Per-user CPU guarantee (e.g. "0.5"). Empty = none (the chart
default).

### spec.singleUser.cpuLimit

`string`

Per-user CPU limit. Empty = unlimited.

### spec.singleUser.storage

`KubernetesJupyterHubUserStorage`

Per-user home storage. Empty = a dynamically-provisioned 10Gi
ReadWriteOnce PVC per user, named `claim-<username>`, mounted at
/home/jovyan — KEPT when the user's server stops AND when this
resource is destroyed (user homes are data; delete `claim-*`
PVCs explicitly, or the namespace, to reclaim them).

### spec.singleUser.storage.dynamic

`KubernetesJupyterHubDynamicStorage`

A dynamically-provisioned PVC per user (the chart default) —
each user's home survives server stops and hub upgrades.

### spec.singleUser.storage.dynamic.capacity

`string` · optional (explicit presence)

Per-user volume size. Empty = "10Gi" (the chart default).

- default: `10Gi`
- rule: {"string":{"pattern":"^[0-9]+(\\.[0-9]+)?(Ei|Pi|Ti|Gi|Mi|Ki|E|P|T|G|M|k|m)?$"}}

### spec.singleUser.storage.dynamic.storageClass

`string`

Storage class for user volumes. Empty = the cluster default.

### spec.singleUser.storage.static

`KubernetesJupyterHubStaticStorage`

One existing shared volume for all users, each user in a
subPath. The PVC (and its storage class) must support the
concurrent access your users create — ReadWriteMany for
multi-node clusters.

### spec.singleUser.storage.static.pvcName

`string` · required

Name of the existing PVC holding all user homes.

- rule: {"required":true}

### spec.singleUser.storage.static.subPath

`string` · optional (explicit presence)

Per-user subPath template within the volume. Empty =
"{username}".

- default: `{username}`

### spec.singleUser.storage.none

`KubernetesJupyterHubNoStorage`

No persistent home — user workspaces are EPHEMERAL and vanish
when the server stops (culling included). Only for stateless
teaching setups where all work lands in external systems.

### spec.singleUser.defaultUrl

`string`

The URL users land on after their server starts. Empty = the
image's default (JupyterLab on current images). Set "/lab" to
force JupyterLab or "/tree" for the classic notebook UI.

### spec.singleUser.startTimeoutSeconds

`int32` · optional (explicit presence)

Seconds the hub waits for a user server to become ready before
failing the spawn. Empty = 300 (the chart default) — raise it
for very large images on slow-pulling nodes (or rely on the
pre-puller, which is the real fix).

- default: `300`
- rule: {"int32":{"gte":30}}

### spec.singleUser.extraEnv

`map<string, string>`

Extra environment variables injected into every user server
(plain values — for secrets, mount them via `helm_values`
singleuser.extraFiles or bake them into your image's tooling).

### spec.singleUser.profiles

`[]KubernetesJupyterHubProfile`

The server-options menu shown at spawn time. Empty = no menu;
every user gets the default image/sizing above. Each profile can
override the image and sizing — "Small CPU", "GPU workstation",
"Big-memory ETL".

### spec.singleUser.profiles[].displayName

`string` · required

The name users see in the server-options menu.

- rule: {"required":true}

### spec.singleUser.profiles[].description

`string`

One-line description under the name.

### spec.singleUser.profiles[].default

`bool`

Pre-select this profile in the menu. At most one profile should
set it.

### spec.singleUser.profiles[].image

`KubernetesJupyterHubImage`

Image override for this profile. Empty = the spec-level
single_user image.

### spec.singleUser.profiles[].image.repository

`string` · required

Image repository including any registry host (e.g.
"quay.io/jupyter/scipy-notebook").

- rule: {"required":true}

### spec.singleUser.profiles[].image.tag

`string` · required

Image tag (e.g. "2026-07-28").

- rule: {"required":true}

### spec.singleUser.profiles[].memoryGuarantee

`string`

Memory guarantee override (e.g. "4G").

### spec.singleUser.profiles[].memoryLimit

`string`

Memory limit override.

### spec.singleUser.profiles[].cpuGuarantee

`string`

CPU guarantee override.

### spec.singleUser.profiles[].cpuLimit

`string`

CPU limit override.

### spec.scheduling

`KubernetesJupyterHubScheduling`

User-pod scheduling machinery: the packing scheduler, capacity
placeholders and pod priorities. Empty = the chart defaults
(user-scheduler on, no placeholders).

### spec.scheduling.userSchedulerEnabled

`bool` · optional (explicit presence)

Run the chart's user-scheduler (a kube-scheduler configured to
PACK user pods onto the busiest fitting node instead of
spreading them) — what lets cluster autoscalers actually scale
user capacity DOWN. Empty = true (the chart default).

- default: `true`

### spec.scheduling.userPlaceholderReplicas

`int32` · optional (explicit presence)

Number of user-placeholder pods — dummy pods at lower priority
than real users that keep spare capacity warm; a real user
evicts one instantly instead of waiting for a node to boot.
Sized like one default user each. Setting this above 0 enables
the chart's pod-priority machinery. Empty = 0.

- default: `0`
- rule: {"int32":{"gte":0}}

### spec.scheduling.coreNodeSelector

`map<string, string>`

Node selector applied to the hub, proxy and other CORE pods (the
chart's per-component nodeSelector surface; user pods schedule
via the user-scheduler and `helm_values`
singleuser.nodeSelector).

### spec.scheduling.userNodeSelector

`map<string, string>`

Node selector applied to USER pods.

### spec.culling

`KubernetesJupyterHubCulling`

Idle-server culling. Empty = enabled, stop notebook servers idle
for an hour (the chart default) — the single most important knob
for not paying for abandoned notebooks.

### spec.culling.enabled

`bool` · optional (explicit presence)

Cull idle notebook servers. Empty = true (the chart default) —
disable only when notebooks must never be stopped automatically.

- default: `true`

### spec.culling.timeoutSeconds

`int32` · optional (explicit presence)

Seconds of inactivity before a server is culled. Empty = 3600
(the chart default).

- default: `3600`
- rule: {"int32":{"gte":60}}

### spec.culling.everySeconds

`int32` · optional (explicit presence)

Seconds between cull sweeps. Empty = 600 (the chart default).

- default: `600`
- rule: {"int32":{"gte":30}}

### spec.culling.maxAgeSeconds

`int32` · optional (explicit presence)

Cull servers older than this many seconds REGARDLESS of
activity (a hard session lifetime). Empty/0 = no age cap (the
chart default).

- rule: {"int32":{"gte":0}}

### spec.culling.cullUsers

`bool`

Also delete the culled USER record (not just their server) —
for authenticators minting throwaway identities. Off by default
(matching the chart); never enable with real user rosters.

### spec.prePuller

`KubernetesJupyterHubPrePuller`

Notebook-image pre-pulling. Empty = the chart defaults (both
pullers on). The hook puller runs BEFORE install/upgrade
completes — first installs wait for the notebook image to land on
every node, which is exactly what makes first spawns fast.

### spec.prePuller.hookEnabled

`bool` · optional (explicit presence)

Pre-pull the notebook image on every node BEFORE
install/upgrade completes (a hook DaemonSet + an awaiter Job the
release waits on). Makes first spawns fast at the cost of slower
installs — the install is pulling gigabytes to every node.
Empty = true (the chart default).

- default: `true`

### spec.prePuller.continuousEnabled

`bool` · optional (explicit presence)

Keep a DaemonSet running that pre-pulls the notebook image onto
NEW nodes as they join — spawns on fresh autoscaled nodes skip
the pull. Empty = true (the chart default).

- default: `true`

### spec.networkPolicyEnabled

`bool` · optional (explicit presence)

Render the chart's NetworkPolicies for hub, proxy and user pods
(user pods get a deny-by-default posture that still allows DNS,
the hub and the internet; the cloud metadata server is blocked).
Empty = true (the chart default). Requires a CNI that enforces
NetworkPolicy for any actual effect.

- default: `true`

### spec.helmValues

`string`

Additional Helm values merged LAST (Helm `-f` semantics,
identical on both engines) — the escape hatch for chart values
the typed fields do not model: hub.extraConfig snippets,
hub.services/loadRoles, per-component securityContext, LDAP
auth, network-policy egress rules, singleuser lifecycle hooks.
YAML document as a string. Never put secret material here —
everything in Helm values lands base64-readable inside the
chart-owned hub Secret; credentials belong in the typed
secret-reference fields, which ride environment indirection
instead. Never set proxy.https/ingress/httpRoute — exposure
composes from first-class kinds.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesJupyterHub, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.namespace` | `string` | Namespace JupyterHub runs in. |
| `status.outputs.proxy_public_service` | `string` | Name of the proxy-public Service — the instance's front door (all user traffic, login included, enters here). The handle exposure kinds route to. Chart-fixed name: "proxy-public". |
| `status.outputs.endpoint` | `string` | In-cluster endpoint of the front door, `http://proxy-public.<namespace>.svc.cluster.local:80`. |
| `status.outputs.hub_service` | `string` | Name of the hub's own Service ("hub", port 8081) — the hub REST API handle for in-cluster automation (`/hub/api`); user traffic goes through proxy-public instead. |
| `status.outputs.shared_password_secret` | `KubernetesSecretKey` | The shared sign-in password: the Secret and key holding it when the shared-password method is active (module-generated `<name>-auth` unless an existing Secret was declared). Empty for OAuth/OIDC/native sign-in — those identities live with the provider. |
| `status.outputs.shared_password_secret.name` | `string` | The name of the Kubernetes Secret. |
| `status.outputs.shared_password_secret.key` | `string` | The key within the Kubernetes Secret. |
| `status.outputs.port_forward_command` | `string` | Port-forward command for reaching JupyterHub from a workstation when no exposure is composed (`kubectl port-forward svc/proxy-public -n <namespace> 8080:80`). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespace` | KubernetesNamespace | `spec.name` |
| `spec.hub.database.postgres.host` | KubernetesPostgres | `status.outputs.rw_service` |
| `spec.hub.database.postgres.passwordSecret.secretName` | KubernetesPostgres | `status.outputs.password_secret.name` |
| `spec.hub.database.mysql.host` | KubernetesMysql | `status.outputs.primary_service` |
| `spec.hub.database.mysql.passwordSecret.secretName` | KubernetesMysql | `status.outputs.root_password_secret.name` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
