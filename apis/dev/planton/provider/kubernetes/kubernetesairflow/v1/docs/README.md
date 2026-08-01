# Kubernetes Airflow — design notes

## Grain

One resource = one Airflow installation from the official `airflow`
chart (https://airflow.apache.org; chart 1.22.0 = Airflow 3.2.2 at the
pin — the chart pin governs). At the chart's default naming scheme
(`useStandardNaming: false`) the release name IS the fullname, so
every child is `<name>-<suffix>` and the exported outputs are
deterministic. Names are capped at 40 characters — the longest
always-rendered child suffix is `-run-airflow-migrations` (23) and
Kubernetes caps names at 63; both engines fail loudly instead of
letting child creation fail midway. This kind models the Airflow 3
line only (CEL-enforced ≥ 3.0.0): the Airflow 2 webserver/flower-era
chart surfaces are deliberately unmodeled.

## The connection-Secret contract

The chart consumes database and broker connections as Secrets carrying
a full URI under the `connection` key; its alternative (split
connection values) renders the password INTO a Helm-created Secret —
the forbidden path. The module therefore owns every connection Secret:
it reads the referenced credential Secret AT APPLY TIME (Terraform: a
data source deferred behind a module-created Secret, so offline plans
never dial a cluster; Pulumi: a DryRun-gated read) and composes
`postgresql://user:pass@host:port/db?sslmode=…` (metadata),
`db+postgresql://…` (the Celery result backend — SQLAlchemy's scheme
prefix, chart-verified) and `redis://…` (broker). With PgBouncer on,
the URIs are rewritten exactly as the chart's own templates would —
host to `<name>-pgbouncer`, database to the ini aliases
(`<name>-metadata`/`<name>-result-backend`) — and the module composes
`pgbouncer.ini` + `users.txt` into the pooler's config Secret
(`pgbouncer.configSecretName`), because the chart's own rendering path
embeds the raw password in values.

## Render-time random credentials, replaced

The chart generates the Redis password, Fernet key, API/webserver
session keys and JWT secret with template `randAlphaNum` — a NEW value
on every upgrade render (upstream's own comment calls the pre-install
hook that papers over this a hack). The module generates each once
(random provider, generation-shape arguments ignored after creation)
into stable Secrets and passes only their NAMES. The Fernet key is 32
random BYTES in URL-safe base64 (the exact shape Fernet requires); its
Secret carries a companion `fernet-key-std-b64` key — the
standard-base64 form the random_bytes importer takes — so blind
round-trips derive the import value mechanically.

## The admin bootstrap path

The chart's `createUserJob` renders the admin password as a LITERAL
POD ARGUMENT (`-p admin` at defaults). The module overrides the Job's
args to read `$(ADMIN_PASSWORD)` from a job-scoped env var backed by
the admin Secret (`<name>-admin-auth`, module-generated unless
bring-your-own) — Kubernetes expands `$(VAR)` in args natively, so the
credential never appears in a rendered pod spec. The chart's default
`config` pins the FAB auth manager, so password login and
`airflow users create` are the chart's own posture at this pin.

## Executor and broker pairing

`executor` mirrors the chart: a plain string, comma-separated
multi-executor allowed, custom class paths allowed (schema-regex
mirrored as CEL). This kind DEFAULTS to `KubernetesExecutor` — the
zero-dependency Kubernetes-native path — where the chart defaults to
CeleryExecutor; the divergence is deliberate, documented on the field,
and validation-enforced both ways (a Celery family member requires
`broker`; a broker without an explicit Celery executor is rejected,
because the unset-executor default is broker-less). The bundled Redis
arm keeps upstream's own `redis:7.2-bookworm` pin (the last
BSD-licensed line — noted on the field); the composed-Valkey arm is
the licensing-clean production path.

## Log surfaces

`logging.persistence` is the shared-volume path (ReadWriteMany needed
past one replica). `logging.elasticsearch`/`opensearch` configure the
UI's log READ path only (secret-composed like every connection);
shipping logs INTO the backend is the log pipeline's job — a
KubernetesOtelCollector composes naturally. Airflow's remote WRITE
path (S3/GCS remote_logging) is an airflow.cfg surface — it rides
`helm_values` under `config.logging`, taught rather than typed at this
pin.

## Deliberate exclusions

Kerberos, the chart's own OTel collector sidecar (contrib 0.70.0 — a
convenience sidecar; the catalog's OTel kinds are the real path),
Flower (default off upstream), cleanup/databaseCleanup CronJobs,
`multiNamespaceMode`, per-component env/volumes and the
KEDA-bypasses-PgBouncer posture (`workers.keda.usePgbouncer=false` —
its `kedaConnection` key is chart-rendered from split values, which
the secret-name path deliberately avoids) all ride `helm_values`.
Airflow 3's dag-bundles (GitDagBundle) need a connection in the
Airflow DB — a day-2 surface recorded as a follow-up question;
git-sync is the modeled delivery path.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
