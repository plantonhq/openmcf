# Trino

The distributed SQL query engine: query your data where it lives —
PostgreSQL, MySQL, object stores, data lakes — and JOIN across sources
in one query. A coordinator plans, a worker fleet executes, and every
data source is just another catalog prefix in your SQL.

## Highlights

- **Queryable in one apply** — the in-image `tpch`/`tpcds` sample
  catalogs answer real SQL on a fresh install; composed
  PostgreSQL/MySQL catalogs turn your databases into federated,
  JOIN-able sources.
- **Secured by default** — upstream's no-authentication posture never
  ships: PASSWORD auth with a generated admin credential and the
  required internal shared secret are on from the first apply.
- **Secret-native catalogs** — catalog passwords reach Trino as
  environment variables referenced through `${ENV:...}` (Trino's own
  secrets mechanism); nothing credential-bearing lands in rendered
  configuration.
- **Elastic workers** — fixed count, HPA on utilization, or KEDA on
  Prometheus metrics with scale-to-zero; graceful shutdown drains
  running queries before any worker terminates.
- **Fault-tolerant execution** — task-level retries with spooled
  exchanges over an S3-compatible bucket: the ETL posture that
  survives worker loss.

Governance rides along: file-based access-control rules, resource
groups (per-group concurrency/memory/queue budgets), session-property
policies and event listeners are all first-class fields.
