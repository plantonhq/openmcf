# Argo Workflows

Deploy [Argo Workflows](https://argoproj.github.io/argo-helm) — the Kubernetes-native workflow engine for DAG/step pipelines (CI jobs, data and ML pipelines, batch orchestration) — from the official `argo-workflows` Helm chart. Every workflow step runs as a pod; the engine turns declarative pipeline documents into running containers with retries, artifacts, and history.

The grain is deliberate: **this kind installs the ENGINE only.** Workflows, WorkflowTemplates and CronWorkflows are Kubernetes custom resources declared like any other manifest (a KubernetesManifest, a chart, the Argo CLI/UI) once the engine runs. Submit against the runner ServiceAccount (`workflow_service_account`, default `argo-workflow`) — workflow pods run with ITS permissions, never the controller's.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Helm release** (official `argo-workflows` chart, default pin `1.0.23` — ships Argo Workflows v4.0.8, named `metadata.name`) — the **workflow controller** (turns Workflow CRs into pods), the **Argo server** (UI + REST API, on by default, ClusterIP), and the **runner ServiceAccount + Role/RoleBinding** workflow pods execute under — created in the install namespace and in every `controller.workflow_namespaces` entry
- **The Argo Workflows CRDs** — installed by default through the chart's hook Job, which **downloads the full-schema CRD YAMLs from the chart's GitHub release at install time** (they exceed inline-template limits). This is the one place the install reaches the internet; see the air-gap note below
- **Kubernetes Namespace** — created only when `create_namespace` is true; otherwise the namespace must already exist

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** — an active connection in the Connect module with credentials for the target cluster.

### Cluster Side

- **An object-storage bucket** — only if you declare `artifact_repository`. S3 (or an in-cluster **Kubernetes SeaweedFS** S3 endpoint), GCS, or Azure Blob. The bucket must exist; Argo Workflows never creates it.
- **A relational database** — only if you declare `archive`. Postgres or MySQL, with the database already created (the controller creates its tables, never the database) and a credentials Secret in the install namespace. An in-cluster **Kubernetes Postgres** pairs naturally.
- **Internet reachability (or a mirror)** — the default CRD install downloads from GitHub. Air-gapped clusters set `crds.full_schema = false` (chart-templated minified CRDs, no download, weaker server-side validation) or point `crds.base_url` at an internal mirror.
- **kube-prometheus-stack** — only if you set `service_monitor_enabled` (the monitoring.coreos.com CRDs must exist first).

## Deploy

### Console

Open the deployment store, find **Argo Workflows**, and click **Deploy**. The creation wizard walks you through namespace placement, the chart pin, the runner identity (the ServiceAccount and the namespaces it lands in), the controller (replicas, instance ID, parallelism caps), the server and its access posture (auth modes, TLS, base URL), the artifact repository, the workflow archive, CR retention, the CRD lifecycle, metrics, the air-gap image path, placement, and the Helm-values escape hatch. Start from the **Dev Pipelines** preset in the [Presets](#presets) tab.

### CLI

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesArgoWorkflows
metadata:
  name: dev-workflows
  org: acme-corp
  env: dev
spec:
  namespace:
    value: dev-workflows
  create_namespace: true
```

```shell
planton apply -f argo-workflows.yaml
```

This near-empty spec is a complete engine: the controller, the server in `client` auth mode (callers present a Kubernetes token and act with its permissions), and the `argo-workflow` runner ServiceAccount. Submit a Workflow CR and it runs; open the UI over the exported port-forward command. No artifact repository and no archive yet — steps cannot pass files, and history lives only as Workflow CRs until something prunes them.

### InfraChart

Compose Argo Workflows behind its namespace with a reference, and the InfraPipeline orders the deploys:

```yaml
spec:
  namespace:
    valueFrom:
      kind: KubernetesNamespace
      name: pipelines-namespace
      fieldPath: spec.name
  create_namespace: false
```

## Key Configuration

**The runner identity is the security boundary.** Workflow pods execute as `workflow_service_account` (default `argo-workflow`) with that account's permissions — never the controller's. Grant IRSA/Workload-Identity annotations or extra RBAC to THIS account when workflows need cloud or cluster access; never widen the controller's. `controller.workflow_namespaces` places the runner ServiceAccount + Role/RoleBinding in additional namespaces so workflows can run there — a workflow submitted in a namespace without that identity fails at pod creation. KNOW THIS: the controller WATCHES cluster-wide either way; the list places RBAC, it does not scope the watch.

**Two durability seams — both empty by default, both deliberate.** Declare `artifact_repository` so steps can pass files (and, with `archive_logs`, so the UI shows step logs after pods are gone): exactly one backend — S3-compatible, GCS, or Azure Blob. Declare `archive` (Postgres/MySQL) so run history survives the Workflow CRs' garbage collection; without it, history lives only in etcd until `retention_policy` or TTLs prune it. With the archive keeping history, `retention_policy` caps the CRs kept per outcome (completed/failed/errored) — explicit `0` means "keep none of that outcome," which is honest and preserved; empty means keep everything.

**Object storage is keyless first.** The S3 arm enforces exactly one credential path: `use_ambient_credentials` (IRSA / workload identity on the RUNNER service account) XOR a `credentials_secret_name` holding `accesskey`/`secretkey` — a KubernetesSeaweedFs credentials Secret already has this shape, and its S3 endpoint composes by reference into `endpoint`. GCS and Azure are keyless by absence: leave the credentials Secret name empty for GKE workload identity / Azure workload identity on the runner account.

**Server auth modes are the front-door posture.** Empty = `client`: callers present a Kubernetes bearer token and act with ITS permissions — auditable, least-privilege, the right default. `server` runs every request with the server's own ServiceAccount — no login, anyone who reaches the endpoint has the server's full power; acceptable only behind trusted network boundaries. `sso` opens OIDC login (configure the chart's `server.sso` block via `helm_values`; its client ID/secret ride existing Secrets by the chart's own contract). `secure` flips the server to HTTPS with its self-signed certificate — probes, clients, and the exported endpoint output all follow the scheme. The server Service stays ClusterIP; expose via KubernetesIngress or Gateway API kinds over the exported handle, and set `base_href` when a public URL fronts it.

**THE CRD DELETION TRAP:** `crds.keep` defaults to true for a reason — removing the CRDs deletes EVERY Workflow, WorkflowTemplate and CronWorkflow in the cluster with them. Turn it off only on throwaway clusters. And `crds.install = false` is only for clusters where another release already owns the CRDs.

**Several engines, one cluster — the instance ID is the mechanism.** Set `controller.instance_id` and this controller reconciles ONLY Workflow CRs labeled with its instance ID, ignoring the rest — team A's engine, team B's engine and a platform engine coexist without stealing each other's work. Controller `parallelism` (cluster-wide) and `namespace_parallelism` cap how many workflows run at once; extra controller `replicas` are hot standbys behind leader election, not workload sharing.

**The image override maps onto the chart's registry/repository split.** `image.registry` replaces the registry part (default `quay.io`) for all three components — workflow-controller, argocli, argoexec — while the repository paths stay upstream; a mirror that re-paths repositories overrides those via `helm_values`. `image.tag` pins all three; `pull_secret_name` names an existing image-pull Secret.

**`helm_values` merges last** — the escape hatch for chart surface beyond the typed fields (workflowDefaults documents, executor resources, the `server.sso` block, extra env, per-component priority classes). Anything here silently overrides the typed fields on every deploy; never put secret material in it — every credential path in this spec rides existing Secrets.

## Outputs and Dependencies

### What This Component Consumes

| Field | References | Purpose |
|-------|-----------|---------|
| `spec.namespace` | KubernetesNamespace (`spec.name`) | Where the engine runs |
| `spec.artifact_repository.s3.endpoint` | KubernetesSeaweedFs (`status.outputs.s3_endpoint`) | In-cluster S3-compatible artifact storage |
| `spec.artifact_repository.s3.credentials_secret_name` | Existing Secret (`accesskey`/`secretkey`) | Declared S3 keys — only when ambient identity is unavailable |
| `spec.artifact_repository.gcs.credentials_secret_name` | Existing Secret (`serviceAccountKey`) | GCS key — only when GKE workload identity is unavailable |
| `spec.artifact_repository.azure.credentials_secret_name` | Existing Secret (`account-access-key`) | Azure account key — only when workload identity is unavailable |
| `spec.archive.host` | KubernetesPostgres (`status.outputs.rw_service`) | The archive database host |
| `spec.archive.credentials_secret` | Existing Secret (name + username/password keys) | Archive database credentials |
| `spec.image.pull_secret_name` | Existing Secret | Air-gap/private-mirror image pulls |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `namespace` | Namespace the engine runs in | Application deployment manifests |
| `release_name` | Helm release name (= metadata.name); child names derive from it | Operational tooling |
| `server_service` | The Argo server Service (`<name>-server`, port 2746). Empty when the server is disabled | Ingress/Gateway exposure |
| `server_kube_endpoint` | In-cluster server URL — plain HTTP by default, HTTPS when `server.secure` is set. Empty when the server is disabled | API clients, CLI configuration |
| `workflow_service_account` | The runner ServiceAccount name — annotate THIS account for IRSA/workload identity | Cloud-identity wiring for workflow pods |
| `port_forward_command` | Copy-paste `kubectl port-forward` for the server UI | Local development access |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

| Preset | Posture | What it carries |
|--------|---------|-----------------|
| **Dev Pipelines** | dev loop / first pipelines | The smallest useful engine: controller, UI in `client` auth mode, and the runner identity — no artifact repository, no archive; declare both before pipelines produce anything worth keeping |
| **Durable Pipelines** | production pipelines | Both durability seams filled: an in-cluster S3-compatible artifact store (composes from a KubernetesSeaweedFs) with archived logs, a Postgres archive (composes from a KubernetesPostgres), a working-set retention policy, controller resources, and ServiceMonitors |
| **Multi Team** | several engines in one cluster | The instance-ID mechanism: this controller reconciles only its own labeled Workflows, the runner identity lands in the team's app namespaces, parallelism caps keep noisy pipelines contained, and a team-scoped runner name keeps identity grants unmistakable |

## Works With

- **Kubernetes SeaweedFS** — in-cluster S3-compatible artifact storage: its `s3_endpoint` output composes into the s3 arm by reference, and its credentials Secret already carries the `accesskey`/`secretkey` shape.
- **Kubernetes Postgres** — the workflow archive: its `rw_service` output composes into `archive.host`, and its app-credentials Secret is the one the archive reads.
- **Kubernetes kube-prometheus-stack** — scrapes the controller's ServiceMonitors when `service_monitor_enabled` is set.
- **Kubernetes Manifest** — the vehicle for Workflow, WorkflowTemplate and CronWorkflow CRs once the engine runs.
- **Kubernetes Namespace** — referenced placement; the InfraPipeline orders namespace-first.
- **Kubernetes Ingress / Gateway API kinds** — HTTP exposure over the exported server Service handle.
