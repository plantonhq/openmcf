# KubernetesMlflow — research notes

Design substance for this kind comes from the MLflow application
source at the pin (mlflow/mlflow v3.15.0, Apache-2.0) and its in-repo
deployment reference. The distribution ruling: MLflow publishes NO
served Helm chart — the repository carries a chart under `charts/`,
but its CI workflow only lints, templates and kind-tests it (no
release/publish step exists; version 0.1.0, appVersion lagging), so
there is nothing served to pin. The mature community chart was
disqualified on image provenance (it runs a personal third-party
rebuild of MLflow, not the official image) and on dead bundled
dependencies. This module therefore renders MODULE-OWNED typed
manifests around the OFFICIAL image `ghcr.io/mlflow/mlflow:v3.15.0-full`
(registry-verified), using the unpublished in-repo chart as the
canonical pod-shape reference: `mlflow server` command args,
MLFLOW_BACKEND_STORE_URI via env (the secretKeyRef seam is upstream's
own design), `/health` probes, `mlflow gc` as a CronJob, and the
Recreate-when-RWO strategy rule.

## Image-variant truth (verified live — the pod crash-loops otherwise)

The bare `ghcr.io/mlflow/mlflow:vX.Y.Z` image is `pip install mlflow`
on python-slim and NOTHING else (its Dockerfile at the pin). Upstream's
own docker README states the consequence plainly: "Most integrations
(backend store databases, artifact stores, etc.) will not work without
additional packages." Concretely missing: the database drivers
(psycopg2-binary/PyMySQL — a postgres backend dies at boot with
`No module named 'psycopg2'`), the object-store clients (boto3/GCS/
Azure — every remote artifact destination), and Flask-WTF (the `auth`
extra — basic auth itself). The `-full` variant
(`pip install mlflow[extras,azure,db,gateway,genai,auth]`, published
since v3.9.0) is the only official image that serves this kind's
modeled surface, so it is the module default. A private mirror or
custom image must carry the same dependency set.

## Server contract at the pin

- `mlflow server --host 0.0.0.0 --port 5000 --workers N` (uvicorn
  workers; upstream default 4).
- `--artifacts-destination <uri> --serve-artifacts`: the server
  proxies artifact upload/download through its own API — clients never
  carry store credentials. Direct-access artifact roots
  (`--default-artifact-root` with client-side credentials) are
  deliberately unmodeled: the proxied posture is both simpler and
  safer, and multi-replica correctness is identical.
- `/health` answers unauthenticated even with auth on — the probe
  contract.
- `--expose-prometheus <path>` enables the /metrics endpoint (the
  metrics arm; the ServiceMonitor is module-rendered and requires the
  Prometheus operator CRDs).

## Auth truth (server source, `mlflow/server/auth/`)

`mlflow server --app-name basic-auth` + an ini file via
MLFLOW_AUTH_CONFIG_PATH. The ini carries `admin_username`,
`admin_password` (UPSTREAM DEFAULT admin/password1234 — never ships;
the module generates or composes the real value into the Secret-mounted
ini), `default_permission` (READ/EDIT/MANAGE/NO_PERMISSIONS) and
`database_uri` for the auth tables. The module points the auth database
at the backend store (same PostgreSQL database on the database arms —
upstream-supported; a sqlite file beside the tracking data otherwise),
so auth state is exactly as durable as tracking state. Passwords store
PBKDF2-hashed.

The auth app additionally REFUSES to start without
`MLFLOW_FLASK_SERVER_SECRET_KEY` (CSRF protection — `create_app`
raises; verified live at the pin), and upstream requires the key be
CONSISTENT across servers. The module generates it into the
`<name>-auth-config` Secret (key `flask_secret_key`) and env-wires it,
so every replica shares one value by construction.

## Storage arms

- Backend store: sqlite on the `<name>-data` PVC
  (`sqlite:////mlflow/data/mlflow.db`) XOR postgres
  (KubernetesPostgres FK; the module composes
  `postgresql://user:pass@host:port/db` at apply time from the
  referenced credential Secret — the connection-secret class) XOR
  mysql (`mysql+pymysql://`, plan-proven).
- Artifact store: `<name>-artifacts` PVC XOR S3-compatible
  (KubernetesSeaweedFs FK: endpoint + `<name>-s3-secret` credentials;
  MLFLOW_S3_ENDPOINT_URL + AWS_* env — boto3's automatic addressing
  path-styles custom endpoints, which is why no path-style knob exists
  here or upstream) XOR AWS S3 (keyless ambient identity or declared
  keys) XOR GCS (keyless or a mounted service-account key) XOR Azure
  Blob (storage-account key).

## Replica truths (CEL-enforced)

More than one replica requires the postgres backend (sqlite is a
single-writer file — and auth state follows the backend store) AND an
object artifact store (the PVC arm is ReadWriteOnce). The Deployment
strategy follows the same volume truth: Recreate whenever a PVC mounts,
RollingUpdate otherwise (upstream's own rule).

## Garbage collection

`mlflow gc --older-than <window> --backend-store-uri $(URI)` as a
CronJob: soft-deleted runs/experiments become unrecoverable once
collected; the retention window is the undo period. The URI rides env
expansion from the Secret-sourced variable, so it never appears in the
rendered spec. On the sqlite arm the job mounts the data PVC (RWO —
meaningful mostly on single-node clusters; postgres is the real gc
story).

## Deliberate exclusions

- The AI-gateway/serving surfaces (`mlflow gateway`, model serving):
  different products with different lifecycles — candidates for their
  own kinds on demand evidence.
- TLS termination at the server (`--uvicorn-opts` certfiles): TLS
  composes at the exposure layer.
- Workspaces/multi-tenancy flags beyond `default_permission`: young
  surface at the pin; `extra_args`/`extra_env` cover early adopters.
