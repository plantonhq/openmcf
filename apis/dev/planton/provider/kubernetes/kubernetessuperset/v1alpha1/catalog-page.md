# Superset

The open-source business-intelligence platform: explore any SQL
database, build charts in a no-code builder or SQL Lab, and assemble
them into shareable dashboards with row-level security and
dashboard-level RBAC.

## Highlights

- **Composition-first** — the metadata database composes a
  KubernetesPostgres and the cache/broker composes a KubernetesValkey;
  both credentials ride environment references to the composed
  resources' own Secrets. The chart's bundled database/redis
  subcharts (frozen legacy images upstream) never ship.
- **Secured by default** — the session-signing key is generated
  (Superset refuses its insecure default) and kept stable (it
  encrypts stored datasource credentials); the admin password is
  generated and delivered through environment, never rendered.
- **The full Celery stack, typed** — workers for async SQL Lab and
  thumbnails, beat for scheduled alerts & reports, flower for queue
  monitoring (off by default: it ships unauthenticated), websockets
  for live async-query push, and the MCP server for AI-agent access.
- **BI over federation** — point a datasource at a composed
  KubernetesTrino and every federated catalog becomes chartable from
  one connection.

Feature flags, config-override python snippets and env escape hatches
keep the long tail of `superset_config.py` reachable without leaving
the typed spec.
