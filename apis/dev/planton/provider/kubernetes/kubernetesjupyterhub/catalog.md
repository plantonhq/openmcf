# JupyterHub

Declares one multi-user JupyterHub install -- the official Zero to JupyterHub Helm chart (`4.4.0` = JupyterHub 5.5.0) rendered SECURED-BY-DEFAULT: the chart's own sign-in default accepts ANY username with NO password, and it never ships -- an empty `authentication` block means a module-generated shared password (the `<name>-auth` Secret), or declare your identity provider (GitHub, Google, or any OIDC issuer -- a composed KubernetesKeycloak realm's endpoints slot straight in). Every teammate who signs in gets their own JupyterLab server pod with a persistent per-user home volume, spawned on demand and culled when idle. Two truths shape the design: chart resource names are FIXED bare strings (`hub`, `proxy`, `proxy-public`) so exactly ONE JupyterHub fits per namespace, and user homes are DATA -- `claim-<username>` PVCs created at runtime that deliberately survive uninstall. Uses a Kubernetes Provider Connection for cluster access.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Kubernetes Namespace** -- created only when `createNamespace` is `true`; otherwise deploys into an existing namespace
- **The Helm release** -- the official `jupyterhub/jupyterhub` chart at the pinned `chartVersion`, rendering:
  - The hub (`hub`, port 8081) -- authentication, spawning, the admin panel and REST API; SINGLE-REPLICA BY DESIGN (no HA upstream; a Recreate strategy pinned), its state on a sqlite PVC or your declared PostgreSQL/MySQL
  - The configurable-http-proxy (`proxy`) and its front-door Service (`proxy-public`) -- EVERY byte of user traffic flows through it
  - Per-user JupyterLab server pods (`jupyter-<username>`) and their home PVCs (`claim-<username>`) -- created AT RUNTIME as users sign in, not by the deploy
  - The packing user-scheduler, warm placeholder pods, the image pre-puller hooks and the idle-culler service, per your scheduling/culling declarations
- **The authentication Secret** -- with shared-password sign-in and no BYO Secret, the module generates the password into `<name>-auth` (readable by cluster admins; sign-in tokens ride environment indirection, never rendered values)
- **NetworkPolicies** -- the chart's hub/proxy/single-user policies render by default (`network_policy_enabled` unset = on)
- **Kubernetes Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with kubeconfig credentials for the target Kubernetes cluster. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline kubeconfig authentication.

### Kubernetes Cluster

- **A namespace with no other JupyterHub in it** -- chart resource names are fixed bare strings; a second install in the same namespace collides. One hub, one namespace.
- **A hub database that ALREADY EXISTS, when declared** -- the postgres/mysql arms connect to a database that must pre-exist (the hub creates tables, never the database). On a composed KubernetesPostgres, declare `jupyterhub` at bootstrap (`initdb`); the spec's foreign-key defaults then wire the host to its read-write Service and the password to its application-user Secret.
- **Credential Secrets in THIS namespace** -- the database password Secret and any OAuth client-secret Secret are read from the workload's own namespace. Co-locate the hub with its database, or replicate the credential Secret.
- **RWX storage for the shared-volume home mode** -- one PVC mounted by every user's server pod is a ReadWriteMany concern on multi-node clusters. The default dynamic per-user mode has no such requirement.
- **A plan for user homes** -- `claim-<username>` PVCs survive this resource's destroy, but DELETING THE NAMESPACE DELETES EVERY USER'S HOME. Back them up with your storage story.

## Deploy

### Console

Open the deployment store, find **JupyterHub**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Team Notebooks** preset for the fastest real multi-user platform (generated shared password, per-user 10Gi homes, hour-idle culling), or **Production OIDC** for identity-provider sign-in on a composed PostgreSQL in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesJupyterHub
metadata:
  name: team-notebooks
  org: acme-corp
  env: dev
spec:
  namespace:
    value: "notebooks"
  create_namespace: true
  authentication:
    shared_password: {}
    admin_users:
      - ada
  single_user:
    image:
      repository: quay.io/jupyter/scipy-notebook
      tag: "2026-07-28"
    memory_limit: 2G
```

```shell
planton apply -f team-notebooks.yaml
```

This declares the secured team default: JupyterHub 5.5.0 (the chart pin's release), shared-password sign-in with the password module-generated into `team-notebooks-auth`, one admin, a real scientific-Python image in place of the chart's evaluation sample, a 2G ceiling per user, and everything else by absence -- sqlite hub state on a 1Gi PVC, dynamic 10Gi per-user homes, hour-idle culling. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the hub's database to a PostgreSQL managed by another Cloud Resource:

```yaml
spec:
  namespace:
    value: "jupyterhub"
  create_namespace: true
  hub:
    database:
      postgres:
        host:
          valueFrom:
            kind: KubernetesPostgres
            name: hub-db
            fieldPath: status.outputs.rw_service
        password_secret:
          secret_name:
            valueFrom:
              kind: KubernetesPostgres
              name: hub-db
              fieldPath: status.outputs.password_secret.name
```

The InfraPipeline deploys the PostgreSQL cluster first, then declares JupyterHub against it -- the SAME parent for both fields: the host resolves the cluster's read-write Service (always the current primary) and the password its application-user Secret. Declare the `jupyterhub` database at the Postgres kind's bootstrap (`initdb`).

## Key Configuration

These are the most important decisions when configuring JupyterHub. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Sign-in, read this first** -- the chart's own default accepts ANY username with NO password; it NEVER ships. An empty `authentication` block is the secured default: shared-password sign-in with the password module-generated into `<name>-auth`. Five explicit arms: `shared_password` (optionally BYO via `password_secret`) | `native` (per-user accounts inside the hub; `open_signup` means anyone reaching the login page can create a WORKING account) | `github` (client ID + client-secret Secret + callback URL; set `allowed_organizations` in production or ANY GitHub identity signs in) | `google` (set `hosted_domains` or any Google account signs in) | `oidc` (six endpoints/fields; Keycloak, Okta, Auth0, Dex all compose). `admin_users` and `allowed_users` sit OUTSIDE the arm choice and apply to every method -- an empty allowed-users roster admits any authenticated identity.

**Hub state, sqlite is right for most installs** -- the hub is SINGLE-REPLICA BY DESIGN, so no `database` block (sqlite on a 1Gi PVC) carries hundreds of users. Declare `postgres` or `mysql` when hub state must survive volume loss or snapshot with your database fleet -- the database itself must PRE-EXIST (declare it at the Postgres kind's `initdb`). Named servers pair with their limit: `named_server_limit_per_user` requires `allow_named_servers` -- the spec holds the pair.

**Per-user sizing has two dials that differ** -- `memory_guarantee` (default 1G) is what the scheduler RESERVES; `memory_limit` (empty = UNLIMITED) is where the kernel kills the kernel. Set a limit in production -- one runaway `pd.read_csv` on an unlimited pod evicts its node neighbors. The chart's sample image is evaluation-grade: point `single_user.image` at a real notebook image (`quay.io/jupyter/scipy-notebook` or your own) for real work.

**Home storage has THREE postures** -- no `storage` block = a dynamic per-user PVC (10Gi default) -- the right default; `static` mounts ONE shared PVC with per-user sub-paths (RWX on multi-node clusters, quota is shared); `none` is the VANISHING WORKSPACE -- everything a user creates disappears when their server stops, and the idle culler stops servers by default. Ephemeral is legal and taught, never accidental.

**The spawn menu** -- `profiles` puts a machine-size menu at login ("Standard 2G", "Big-memory ETL", "GPU workstation"), each row overriding image or sizing with blank overrides falling back to the `single_user` baseline. At most one row is `default: true`. No profiles = no menu -- every user gets the baseline directly.

**Capacity that scales DOWN** -- `user_scheduler_enabled` (default on) packs user pods onto the fewest nodes so autoscalers can reclaim the rest; `user_placeholder_replicas` keeps N warm decoys a real user evicts INSTANTLY instead of waiting for a node boot. Culling is the cost model: idle servers stop after `timeout_seconds` (default 3600) -- homes survive culls; never enable `cull_users` with real user rosters (it deletes hub-database user records).

**The proxy fronts everything, privately by default** -- the chart's own `proxy-public` default is `LoadBalancer` (a public IP on install); this kind DELIBERATELY inverts it to ClusterIP. Reach the hub by port-forward, or compose exposure over the exported Service handle. HTTPS is deliberately unmodeled here -- TLS terminates at your composed ingress/gateway kind.

**The escape hatch, merged LAST** -- `helm_values` is raw chart values with Helm `-f` semantics for what the typed fields don't model (per-pod tolerations, custom hub config snippets, registry credentials). NEVER secret material -- the chart-owned hub Secret embeds the ENTIRE rendered values document, readable by anyone who can read Secrets; NEVER `proxy.https` or ingress keys -- exposure composes from first-class kinds.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **KubernetesPostgres** | `hub.database.postgres.host` | `status.outputs.rw_service` |
| **KubernetesPostgres** | `hub.database.postgres.password_secret.secret_name` | `status.outputs.password_secret.name` |
| **KubernetesMysql** | `hub.database.mysql.host` | its primary Service output (the FK default) |
| **KubernetesKeycloak** (composes naturally) | `authentication.oidc.*` endpoints | the realm's OIDC endpoints |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `namespace` | Namespace JupyterHub runs in | Co-locating the database, exposure kinds and diagnostics |
| `proxy_public_service` | The front-door Service (`proxy-public`) -- every user request enters here | The handle exposure kinds route to |
| `endpoint` | In-cluster endpoint of the front door | What in-cluster clients use |
| `hub_service` | The hub's internal Service (`hub`, port 8081) -- the REST API | Automation driving the hub API (spawning, user admin) |
| `shared_password_secret` | Secret + key holding the generated shared password (EMPTY with GitHub/Google/OIDC/native sign-in) | Distributing the team password; rotating it |
| `port_forward_command` | kubectl port-forward one-liner for the hub UI | Reaching the UI from a workstation on the ClusterIP default |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Team notebooks in minutes** -- The secured default: shared-password sign-in (generated into `<name>-auth`), a real scientific-Python image, per-user 10Gi homes that survive restarts, idle servers culled after an hour. Start from the **Team Notebooks** preset.

**Production on your identity provider** -- OIDC sign-in against a composed KubernetesKeycloak realm (Okta/Auth0/Dex identical), hub state in a composed KubernetesPostgres, a spawn menu of machine sizes, warm placeholder pods keeping "start my server" instant. Start from the **Production OIDC** preset.

**Classroom / ephemeral workshops** -- `storage.none` (no home PVCs to clean up), `native` accounts with `open_signup` inside a private network, aggressive culling (`timeout_seconds: 1800`, `max_age_seconds: 28800`) -- the vanishing workspace as a feature, chosen deliberately.

## Works With

- [**Kubernetes Postgres**](/cloud-catalog/kubernetes-postgres) -- the durable hub-state database; the `hub.database.postgres` foreign-key defaults point at it
- [**Kubernetes Mysql**](/cloud-catalog/kubernetes-mysql) -- the MySQL alternative for hub state
- [**Kubernetes Keycloak**](/cloud-catalog/kubernetes-keycloak) -- the natural OIDC issuer; a realm's endpoints slot into `authentication.oidc`
- [**Kubernetes Namespace**](/cloud-catalog/kubernetes-namespace) -- provides the namespace for the deployment
- [**Kubernetes Http Endpoint**](/cloud-catalog/kubernetes-http-endpoint) -- exposure over the exported `proxy_public_service` handle, where TLS terminates
