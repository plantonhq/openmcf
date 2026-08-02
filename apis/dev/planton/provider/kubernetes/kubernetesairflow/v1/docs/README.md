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
embeds the raw password in values. The pooler's metrics-exporter
sidecar (unconditional at this pin) gets the same treatment: the
chart's own stats Secret composes the exporter DSN from split values
(defaults `postgres:postgres` — a user the auth_file never carries, so
the sidecar crash-loops; verified live), so the module composes
`<name>-pgbouncer-stats` (key `connection`) with the metadata user —
who is also the ini's `stats_users` grant — through the chart's
`statsSecretName` seam.

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

## The bootstrap Jobs run hook-less

The chart defaults its migration and create-user Jobs to post-install
Helm HOOKS — and Helm runs post-install hooks only AFTER the release
wait completes, while every component's `wait-for-airflow-migrations`
init container blocks on the migrations that hook would apply. Under
any wait-style install (both engines wait; Argo CD's sync waits the
same way) that is a deadlock by construction: verified live — no Job
existed while every init container crash-looped on "unapplied
migrations". The module renders `useHelmHooks: false` on BOTH Jobs, so
they apply WITH the release as ordinary resources and the init
containers converge on them. The chart's own
`ttlSecondsAfterFinished: 300` default self-deletes the finished Jobs,
which keeps day-2 applies clean (each upgrade recreates them fresh —
no immutable-Job patch conflicts).

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

## `load_examples` rides the env var, not airflow.cfg

The official image bakes `AIRFLOW__CORE__LOAD_EXAMPLES=False` as a
container env, and Airflow's precedence puts env above airflow.cfg —
a cfg-only `config.core.load_examples: "True"` is silently defeated
(verified live: examples never parsed). The module therefore sets
`AIRFLOW__CORE__LOAD_EXAMPLES=True` through the chart's top-level
`env` list when the field is true — the chart's own docs prescribe
the same route.

## git-sync `ref` must neutralize the chart's legacy `branch`

The chart renders BOTH env generations unconditionally —
`GITSYNC_REF` from `ref` (v4) and `GIT_SYNC_BRANCH` from `branch`
(legacy) — and git-sync v4 translates the deprecated `--branch` OVER
`--ref`. A ref-only render therefore silently syncs the chart's
default branch (`v2-2-stable`, verified live) while the values looked
correct. The module always writes BOTH keys to `dags.git_sync.ref`'s
value — including the empty string, which clears the chart defaults
so Empty = HEAD holds.

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
