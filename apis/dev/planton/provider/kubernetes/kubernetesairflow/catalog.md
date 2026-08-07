# Airflow

Declares one Apache Airflow 3 install -- the official Helm chart (`1.22.0` = Airflow 3.2.2) rendered SECRET-NATIVELY: the metadata database is a REQUIRED external PostgreSQL or MySQL (the chart's bundled database subchart is non-production by upstream's own definition and always disabled by the module), and every credential rides Secret references -- the SQLAlchemy connection URI, the Celery broker URL, and the security keys all compose into module-owned Secrets at apply time, never into rendered Helm values. One decision shapes everything else: the EXECUTOR. Empty = `KubernetesExecutor` (every task runs as its own ephemeral pod -- no broker to operate, no idle worker fleet; note the chart's own default is CeleryExecutor, deliberately flipped here), or declare a Celery-family executor for a standing worker fleet on a message broker -- low task latency, KEDA scaling on real queue depth, and the broker becomes required (the spec holds the pairing in both directions). Uses a Kubernetes Provider Connection for cluster access.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Kubernetes Namespace** -- created only when `createNamespace` is `true`; otherwise deploys into an existing namespace
- **The Helm release** -- the official `apache-airflow/airflow` chart at the pinned `chartVersion`, rendering:
  - The API server (`<name>-api-server`, port 8080) -- UI + REST API; safe to scale, state lives in the database
  - The scheduler -- the heart; more than one replica = HA via database row locks (PostgreSQL recommended)
  - The DAG processor -- always standalone on Airflow 3
  - The triggerer -- deferred/async operators; enabled by default, with a 100Gi PVC unless sized down
  - With a Celery executor: the worker StatefulSet (per-worker 100Gi PVCs unless persistence is disabled or sized down), optionally KEDA-scaled, plus the chosen broker (the chart's bundled Redis or your composed/external one)
  - The database migration Job and (by default) the admin-user bootstrap Job
- **Module-owned Secrets** -- the metadata connection URI (`key: connection`), the Celery broker URL (`<name>-broker-url`), the admin credential (`<name>-admin-auth`), and STABLE generated security keys (Fernet, API secret, JWT) for any not brought-your-own -- the chart's render-time randoms (regenerated on every upgrade, logging out every UI session) never ship
- **StatsD exporter sidecar** -- on by default (scrape port 9102); no ServiceMonitor ships -- compose scraping via your Prometheus stack's additional scrape configs
- **Kubernetes Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with kubeconfig credentials for the target Kubernetes cluster. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline kubeconfig authentication.

### Kubernetes Cluster

- **A metadata database that ALREADY EXISTS** -- the migration Job creates tables, never the database. On a composed KubernetesPostgres, declare the database (`airflow`) and its owning user at bootstrap (`initdb`); the spec's foreign-key defaults then wire the host to its read-write Service and the password to its application-user Secret.
- **Credential Secrets in THIS namespace** -- the module reads the database (and broker) password Secrets at apply time to compose Airflow's connection Secrets, and a Secret can only be read from the workload's own namespace. Co-locate Airflow with its database, or replicate the credential Secret.
- **A name within budget** -- keep `metadata.name` at 40 characters or fewer: the module derives Job names 23 characters long against the Kubernetes 63-character cap. Both engines fail loudly over the budget.
- **The KEDA operator, when autoscaling workers** -- a KubernetesKeda composes naturally. Without it, the ScaledObject the chart renders has nothing to reconcile it.
- **RWX storage for shared volumes beyond one replica** -- the DAGs shared volume and the task-logs volume are ReadWriteMany concerns once more than one pod mounts them.

## Deploy

### Console

Open the deployment store, find **Airflow**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Dev Kubernetes Executor** preset for the smallest useful Airflow 3 (no broker, no workers), or **Production Celery** for the full queue-driven fleet in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesAirflow
metadata:
  name: daily-pipelines
  org: acme-corp
  env: prod
spec:
  namespace:
    value: "data-platform"
  create_namespace: true
  database:
    postgres:
      host:
        value: "airflow-pg-rw.data-platform.svc"
      password_secret:
        secret_name:
          value: "airflow-pg-app"
  dags:
    git_sync:
      repo: https://github.com/acme-corp/pipelines.git
      ref: main
```

```shell
planton apply -f daily-pipelines.yaml
```

This declares the untouched-defaults install: Airflow 3.2.2 (the chart pin's release), the KubernetesExecutor (no broker, no workers), an existing PostgreSQL reached by hostname with its password read from a co-located Secret, and DAGs git-synced from the repository's default branch on every component -- a merged PR is a deployed pipeline. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire Airflow to a database managed by another Cloud Resource:

```yaml
spec:
  namespace:
    value: "data-platform"
  create_namespace: true
  database:
    postgres:
      host:
        valueFrom:
          kind: KubernetesPostgres
          name: airflow-db
          fieldPath: status.outputs.rw_service
      password_secret:
        secret_name:
          valueFrom:
            kind: KubernetesPostgres
            name: airflow-db
            fieldPath: status.outputs.password_secret.name
```

The InfraPipeline deploys the PostgreSQL cluster first, then declares Airflow against it -- the SAME parent for both fields: the host resolves the cluster's read-write Service (always the current primary) and the password its application-user Secret.

## Key Configuration

These are the most important decisions when configuring Airflow. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**The executor decision, read this first** -- `executor` is a free string, deliberately not an enum: bare class names (`KubernetesExecutor`, `CeleryExecutor`, `LocalExecutor`), comma-separated multi-executor lists (Airflow 3 per-DAG routing), and dotted custom classes are all legal. Empty = `KubernetesExecutor` -- every task its own ephemeral pod, no broker, no idle workers. The spec holds the executor-broker pairing in BOTH directions: a Celery-family executor REQUIRES a `broker`; a `broker` without one is REJECTED.

**The broker, secret-native in every arm** -- `bundled_redis` (the chart's single-instance StatefulSet; upstream pins `redis:7.2-bookworm` -- the last BSD-licensed Redis line; `bundled_redis: {}` with every dial blank is a legal chart-defaults state) | `valkey` (the licensing-clean production path -- a KubernetesValkey composes; an empty password is an UNAUTHENTICATED broker, legal for dev only) | `existing_broker_url_secret` (the escape arm: a Secret carrying the full URL under the `connection` key -- SQS, external TLS Redis). Whatever the arm, the module composes `<name>-broker-url` at apply time -- the chart's render-time random password never ships.

**DAG delivery has THREE postures** -- `git_sync` (recommended: a sidecar on every component keeps DAGs current; push to deploy), `persistence` (CI publishes into a shared volume; RWX beyond one replica), or NO `dags` block at all -- DAGs baked into a custom image, paired with `images.airflow_repository`/`airflow_tag`. Git credentials pair by clone-URL form: HTTPS with `credentials_secret`, SSH with `ssh_key_secret` (pin `known_hosts` with SSH). And the `ref` field writes BOTH git-sync env generations -- including the empty string, which neutralizes the chart's own `v2-2-stable` branch default (empty = HEAD, verified live).

**The three 100Gi defaults** -- the triggerer's PVC, each Celery worker's PVC, and the shared task-logs volume all default to 100Gi in the chart -- far more than most installs need (`1Gi` is plenty on dev clusters). Size them down deliberately; a ten-worker KEDA burst against the worker default would claim a terabyte.

**KEDA scales on REAL queue depth** -- it polls the metadata database for queued and running tasks; `min_replicas` 0 scales the fleet to ZERO between runs (the point of queue-driven scaling), and the plain replica count is ignored while KEDA owns it. `min_replicas` must not exceed `max_replicas` -- held at authoring time.

**PgBouncer is PostgreSQL-only** -- Airflow opens MANY short-lived connections (every task heartbeat is one); production installs on PostgreSQL should enable it (pool defaults 10/5/100). The spec rejects it against a MySQL database, and the module composes the PgBouncer config Secret from the declared credentials -- the chart's password-in-values path is never used.

**Security keys are STABLE by default** -- empty fields = module-generated Fernet key, API secret key, and JWT secret with stable values (the chart would regenerate them on every upgrade render, logging out every UI session). Bring your own by naming Secrets with the chart's exact expected keys. THE DESTRUCTIVE TRUTH: losing the Fernet key orphans every connection password and variable Airflow has stored -- back up `fernet_key_secret_name` (a stack output), share it across DR replicas, rotate only via Airflow's documented procedure.

**Task logs live somewhere, decide where** -- empty `logging` = pod-local (lost on rotation, dev only). `persistence` puts them on the shared volume; `remote_read` points the UI's READ path at an Elasticsearch or OpenSearch your tasks ALREADY ship logs to -- shipping is your log pipeline's job (an OTel collector composes), never implied by this field.

**The escape hatch, merged LAST** -- `helm_values` is raw chart values with Helm `-f` semantics for what the typed fields don't model (Kerberos, per-component env/volumes, cleanup CronJobs, network policies). NEVER secret material -- credentials ride Secrets through the typed fields; NEVER re-enable the postgresql subchart.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **KubernetesPostgres** | `database.postgres.host` | `status.outputs.rw_service` |
| **KubernetesPostgres** | `database.postgres.password_secret.secret_name` | `status.outputs.password_secret.name` |
| **KubernetesValkey** | `broker.valkey.host` | its service output (the FK default) |
| **KubernetesKeda** (runtime prerequisite) | -- | the KEDA operator must run on the cluster when `workers.keda.enabled` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `namespace` | Namespace Airflow runs in | Co-locating exposure kinds and diagnostics |
| `api_server_service` | The API server Service (`<name>-api-server`) -- the UI + REST API front door | The handle exposure kinds route to |
| `api_server_endpoint` | In-cluster endpoint (`http://<name>-api-server.<namespace>.svc.cluster.local:8080`) | What in-cluster clients and REST callers use |
| `admin_password_secret` | Secret + key holding the bootstrap admin password (empty when `admin_user.create` is false) | Signing in; wiring automation that drives the REST API |
| `metadata_connection_secret_name` | Module-owned Secret holding the database connection URI (key `connection`) | Data-lineage sidecars reaching the same database |
| `broker_url_secret_name` | Module-owned Secret holding the Celery broker URL (empty without a Celery executor) | External Celery tooling (e.g. Flower) |
| `fernet_key_secret_name` | Secret holding the Fernet key -- BACK IT UP | DR replicas sharing the key; the backup itself |
| `port_forward_command` | kubectl port-forward one-liner for the Airflow UI | Reaching the UI from a workstation |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Dev on the KubernetesExecutor** -- The smallest useful Airflow 3: API server, scheduler, DAG processor and triggerer against a composed KubernetesPostgres in the same namespace, every task its own ephemeral pod -- no broker, no workers. Start from the **Dev Kubernetes Executor** preset.

**Production Celery with KEDA** -- The full queue-driven fleet: Celery workers KEDA-scaled on real queue depth (and back to zero between runs), DAGs git-synced on every component, PgBouncer pooling the many short-lived connections, two API servers and two schedulers for availability, task logs on a shared volume. Start from the **Production Celery** preset.

**Air-gapped / DAGs baked into the image** -- No `dags` block, `images.airflow_repository`/`airflow_tag` pointing at your custom image (set `airflow_version` when it's built from a different release -- the chart gates version-specific rendering on it), and the remaining image fields at your private mirror.

## Works With

- [**Kubernetes Postgres**](/cloud-catalog/kubernetes-postgres) -- the recommended metadata database; the `database.postgres` foreign-key defaults point at it
- [**Kubernetes Valkey**](/cloud-catalog/kubernetes-valkey) -- the licensing-clean Celery broker; the `broker.valkey` foreign-key defaults point at it
- [**Kubernetes Namespace**](/cloud-catalog/kubernetes-namespace) -- provides the namespace for the deployment
- [**Kubernetes Keda**](/cloud-catalog/kubernetes-keda) -- the runtime prerequisite for queue-driven worker autoscaling
- [**Kubernetes Open Search**](/cloud-catalog/kubernetes-open-search) -- a remote-read backend for task logs; the `logging.opensearch` host default composes its client Service
