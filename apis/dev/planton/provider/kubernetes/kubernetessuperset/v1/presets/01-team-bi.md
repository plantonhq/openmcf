# Team BI

The smallest honest Superset: the web application against a composed
PostgreSQL metadata database (the ONE required input — dashboards,
charts, users and the encrypted datasource credentials live there).
Secured by default: the session-signing key is module-generated
(Superset refuses to start on its insecure default) and the bootstrap
admin signs in with the generated password from the
`superset-admin-auth` Secret — the chart's documented admin/admin
default never ships.

Declare the `superset` database at the composed PostgreSQL's bootstrap
(`initdb.database`). Without a cache this is the web-only shape —
every query runs synchronously; add `spec.cache` (a KubernetesValkey
composes naturally) for async queries, thumbnails and reports.
