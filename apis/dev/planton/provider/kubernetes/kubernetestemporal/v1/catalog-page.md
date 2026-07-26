# Temporal

The durable workflow engine. Business logic that survives crashes,
deploys and weeks of waiting — sagas, human-in-the-loop approvals,
AI-agent pipelines and scheduled orchestration — written as ordinary
code in your language's SDK, executed exactly-once by the engine, with
a Web UI that shows every workflow's full history.

## Highlights

- **Bring your own database, composed by reference** — PostgreSQL (an
  in-cluster `KubernetesPostgres` pairs naturally), MySQL 8, or
  external Cassandra; credentials ride existing Secrets and never
  render into values.
- **The whole engine, typed** — per-service sizing for
  frontend/history/matching/worker, declarative Temporal namespaces
  with retention, dynamic-config limits applied without restarts, and
  history/visibility archival to S3, GCS or a mounted filesystem.
- **Schema handled** — the chart's schema Jobs prepare both databases
  before the server starts; skip them only if you run your own schema
  pipeline.
- **Honest immutability** — `num_history_shards` is taught as the
  one-way door it is, stated explicitly instead of defaulted silently.
- **Clean lifecycle** — no CRDs, no bundled databases or monitoring
  stacks; destroy leaves nothing behind but the database you own.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
