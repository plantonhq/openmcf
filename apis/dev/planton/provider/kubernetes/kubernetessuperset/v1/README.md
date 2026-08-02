# Kubernetes Superset

Deploys Apache Superset — the business-intelligence platform — from
the official ASF Helm chart: the web application, a Celery worker
fleet for async queries/thumbnails/reports, optional beat (schedules),
flower (Celery monitoring), websocket and MCP servers, and an init Job
that migrates the metadata schema and bootstraps the admin.

## Two backends, always composed

The METADATA DATABASE (dashboards, charts, users, and the ENCRYPTED
datasource credentials) is the one required input — a
KubernetesPostgres composes naturally; declare the `superset` database
at its bootstrap. The CACHE/BROKER is a Redis-protocol store — a
KubernetesValkey composes naturally; without it you get a web-only
Superset (synchronous queries; workers, beat, flower, websockets and
MCP stay off — enforced). The chart's bundled postgresql/redis
subcharts ride frozen legacy image lines and never ship from here.

## Secured by default

The session-signing SECRET_KEY is module-generated (Superset refuses
to start on its insecure default) and deliberately STABLE — it also
encrypts stored datasource credentials, so rotation goes through
Superset's own `re-encrypt-secrets` procedure, never a casual change.
The bootstrap admin password is module-generated into
`<name>-admin-auth` (the chart's documented admin/admin default never
ships) and reaches the init step as an environment variable — never a
rendered literal.

## The credential contract

The chart consumes ALL runtime credentials through one environment
Secret. This kind turns the chart's copy OFF and composes `<name>-env`
itself: non-secret connection facts plus generated material live
there; the database and cache passwords arrive in the pods as
environment references to the composed resources' OWN Secrets —
nothing is copied, nothing renders.

## Drivers

The official image is the driver-less "lean" build stage — verified
live, it ships NO database driver at all, not even the PostgreSQL
driver its own metadata database needs (that rides only the dev/ci
image variants). The module's default bootstrap script installs the
exact psycopg2 pin at container start, so the metadata database works
out of the box; a custom `bootstrap_script` REPLACES that default and
must keep a psycopg2 install (or the pods crash-loop). Installs must
target the app's venv — the image's plain `pip` is the system
interpreter's and its installs stay invisible to the app (verified
live); use `uv pip install --python /app/.venv/bin/python <driver>`.
Extra data sources (Trino, Elasticsearch…) add their python drivers
the same way (needs internet, re-runs each restart) — for production,
bake a custom image and set `bootstrap_script` to a no-op. A composed
KubernetesTrino makes every federated catalog chartable from one
Superset connection once the `trino` driver is present.

## Exposure

The Service stays ClusterIP (port 8088); compose exposure kinds over
the exported `service` handle.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
