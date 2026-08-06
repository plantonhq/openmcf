# Kubernetes Mlflow

Deploys MLflow — the experiment-tracking server and model registry —
from the OFFICIAL MLflow image with module-rendered manifests. MLflow
publishes no Helm chart, so this module IS the distribution: a
Deployment, a Service, the credential Secrets, the volumes the chosen
arms need, and optionally a garbage-collection CronJob and a
ServiceMonitor.

## When NOT to Use This

Databricks-managed MLflow already gives you a tracking server — this
kind is for the self-hosted story: your experiments, your models, your
cluster. And if all you need is metrics dashboards, that is the
observability tier, not MLflow.

## Two kinds of state

The BACKEND STORE holds experiments, runs, metrics and the model
registry: sqlite on a PVC by default (zero dependencies, single
replica), or a composed KubernetesPostgres for production. The
ARTIFACT STORE holds the big binaries (models, datasets): a PVC by
default, or S3-compatible object storage — a KubernetesSeaweedFs
composes naturally — or AWS S3 (keyless via ambient pod identity or
declared keys), GCS, Azure Blob.

## The credential path

The database connection URI embeds the password, so the module composes
it AT DEPLOY TIME from the referenced credential Secret into a
module-owned Secret (`<name>-backend-uri`), and the server reads it as
an environment variable — nothing credential-bearing appears in any
rendered manifest. Artifact-store keys ride environment variables from
their Secrets the same way. And because the server PROXIES all artifact
traffic (`--serve-artifacts`), your notebooks and pipelines need only
the tracking URI and their MLflow login — storage credentials never
leave the server pod.

## Secured by default

Upstream's server is OPEN unless auth is configured, and upstream's own
auth example ships admin/password1234. Neither ever ships from here:
basic authentication is ON by default with a module-generated admin
password (`<name>-admin-auth`, the exported credential handle). Auth
state (users, permissions) lives in the backend store — as durable as
the tracking data itself. `default_permission: NO_PERMISSIONS` makes
every experiment private to its creator until shared.

## Scaling honesty

One replica on the default arms — sqlite is a single-writer file and a
ReadWriteOnce artifact volume binds one pod (enforced at validation).
With postgres + an object store the server is stateless: scale
`server.replicas` for API load.

## Exposure

The Service stays ClusterIP; compose exposure from first-class kinds
over the exported `service` handle. Set MLFLOW_TRACKING_URI in training
jobs to the exported `tracking_endpoint`.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
