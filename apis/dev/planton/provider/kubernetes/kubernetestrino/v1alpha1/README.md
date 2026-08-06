# Kubernetes Trino

Deploys Trino — the distributed SQL query engine — from the official
trinodb Helm chart. One coordinator plans and schedules; a worker
fleet executes; every configured data source becomes a CATALOG you
query (and JOIN across) with plain SQL:
`SELECT ... FROM warehouse.public.orders o JOIN lake.events.clicks c ON ...`.

## Secured by default

Upstream Trino ships with NO authentication — anyone who can reach the
Service can query every catalog as any user. That never ships from
here: an empty `auth` block means PASSWORD (file) authentication with
a module-generated admin credential (`<name>-auth` Secret, exported as
the credential handle), and the module configures the
internal-communication shared secret Trino requires once
authentication is on. Trino enforces password authentication only on
secure requests (verified in the server source at the pin — the
`allow-insecure-over-http` flag would bypass the password file
entirely over HTTP, not extend it), so the module sets
`http-server.process-forwarded=true`: terminate TLS at composed
exposure kinds (their `X-Forwarded-Proto: https` engages the password
file) or use the `https` keystore arm; plain-HTTP data-plane requests
are refused outright.

## Catalogs are the product

The chart's in-image `tpch`/`tpcds` sample catalogs stay available
until disabled — a fresh install answers
`SELECT count(*) FROM tpch.tiny.nation` immediately. Declared
`catalogs.postgres` / `catalogs.mysql` entries compose the catalog's
database kinds (host and credential FK onto their outputs); any other
connector rides `catalogs.custom` as raw properties. KNOW THIS: every
catalog renders into a ConfigMap — credentials must ride `${ENV:VAR}`
references (Trino's own secrets mechanism); the typed arms wire this
automatically, and `extra_env_from_secret` pairs with custom catalogs.

## Sizing and scaling

Prefer `jvm.max_heap_percent` with container memory limits (the heap
follows the limit); the chart's fixed default is an 8G max heap —
size limits accordingly or switch to percent. Workers scale by count,
by HPA (CPU/memory utilization), or by KEDA (Prometheus triggers,
scale to zero) — enable `graceful_shutdown` so scaling in never kills
running queries. `coordinator.include_in_scheduling` with
`workers.replicas: 0` gives a true single-node Trino for dev.

## Fault-tolerant execution

`fault_tolerant_execution` turns on task/query retries with exchange
data spooled to durable storage (an S3-compatible bucket composes
naturally) — the batch/ETL posture that survives worker loss and
pairs with spot capacity.

## Exposure

The coordinator Service stays ClusterIP; compose exposure kinds over
the exported `coordinator_service` handle. BI tools (a
KubernetesSuperset composes naturally) connect to the exported
`coordinator_endpoint` with the admin credential.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
