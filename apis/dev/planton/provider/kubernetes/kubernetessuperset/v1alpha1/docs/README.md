# KubernetesSuperset — design notes

## Distribution

The official ASF chart served from https://apache.github.io/superset,
chart 0.22.4 = Superset 6.1.0 (Apache-2.0 end to end). The default
image `apachesuperset.docker.scarf.sh/apache/superset` is the ASF's
scarf.sh download gateway in front of the official image — override
the repository for air-gapped mirrors. The module pins
`fullnameOverride` to the resource name; the longest chart child is
`<name>-celerybeat` (11-char suffix) → names ≤ 52 chars, fail-loud in
both engines.

## The credential architecture

The chart's runtime-credential contract is ONE environment Secret
every component consumes (`envFrom`), and the rendered
`superset_config.py` builds every connection FROM ENVIRONMENT. The
module exploits that shape end to end:

- `secretEnv.create=false` turns the chart's own env Secret (which
  renders credentials from values) OFF; the module composes
  `<name>-env` with the non-secret facts (DB_HOST/PORT/USER/NAME,
  REDIS_*, the admin identity) plus the module-GENERATED material
  (SUPERSET_SECRET_KEY, ADMIN_PASSWORD, the websocket JWT_SECRET).
- REFERENCED credentials — the metadata-database password (DB_PASS)
  and the cache password (REDIS_PASSWORD) — arrive as `extraEnvRaw`
  secretKeyRef entries pointing at the composed resources' own
  Secrets: explicit env beats envFrom (the chart's own
  bring-your-own mechanism), nothing is copied, and no apply-time
  secret reads exist in either engine.
- Two chart config blocks render `cache.password` LITERALLY from
  values (the results backend and the async-queries backends) — this
  kind never sets that value; on authed stores a module-owned
  configOverrides snippet redefines those blocks reading environment.

## The admin bootstrap

`init.createAdmin` stays FALSE (the chart renders the admin password
literally into its config-Secret init script, and its config template
hard-fails on an empty password). The module overrides `init.command`
instead: the chart's own rendered init script (schema migration +
role init) runs first, then an idempotent create-admin step reads
ADMIN_USERNAME/ADMIN_EMAIL/ADMIN_PASSWORD from environment.

## SECRET_KEY truths (app source at the pin)

`SECRET_KEY = os.environ.get("SUPERSET_SECRET_KEY") or CHANGE_ME` —
and the server REFUSES to start on the insecure default. The key also
encrypts datasource credentials stored in the metadata database:
rotation without `superset re-encrypt-secrets` (old key in
PREVIOUS_SECRET_KEY) orphans every stored connection. The module
generates it once and shape-ignores the random.

## Component truths

- The worker Deployment only renders when the cache is declared — a
  worker without a broker crash-loops (CEL-enforced upstream of the
  chart).
- Flower ships with NO authentication of its own — off by default and
  taught on the field.
- The websocket server is a COMMUNITY image
  (`oneacrefund/superset-websocket`, chart default tag `latest`) —
  pin deliberately; it reads JWT_SECRET from environment natively,
  and the module points Superset's async-queries JWT at the same
  variable.
- The MCP server requires the `fastmcp` OPTIONAL extra, which the
  official image does NOT include — custom image or bootstrap pip
  install, or the pods crash-loop.
- The init Job is a post-install/post-upgrade Helm hook (the chart's
  own day-2 re-migration path); schema migration runs inside the
  release wait against the composed database.

## Drivers

The official image is the driver-less "lean" build stage (the
`[postgres]` extras ride only the dev/ci variants — verified live:
the server exits at boot with "No module named 'psycopg2'"). The
module's default bootstrap script installs the exact psycopg2 pin the
app's [postgres] extra declares, so the metadata database works out
of the box; a custom `bootstrap_script` replaces that default and
must keep a psycopg2 install. Installs must target the app's venv:
the image's plain `pip` belongs to the system interpreter and its
installs stay invisible to the app (verified live) — use `uv pip
install --python /app/.venv/bin/python <driver>`. Other datasources
(Trino, Elasticsearch…) add their drivers the same way (internet at
container start, re-runs per restart) or via a custom image — the
production posture.
