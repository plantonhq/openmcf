# HA external database preset

The production posture: two Grafana replicas behind one Service with
ALL state — dashboards, users, sessions, preferences — in an external
Postgres. Any pod can die, be drained or be upgraded and users notice
nothing, because no pod holds anything. This is the only honest way
to scale Grafana: the embedded SQLite database cannot be shared, and
the spec refuses `replicas > 1` without the `database` block for
exactly that reason.

The database is a prerequisite, not a side effect: it must exist
before the install (Grafana creates its tables, never the database),
and its password rides an existing Secret through environment
expansion — it never lands in Grafana's rendered configuration. The
host is shown as a literal; a `value_from` reference to a
KubernetesPostgres resource resolves its read-write endpoint and
orders this Grafana after its database. Admin credentials come from a
team-owned Secret here so rotation follows the team's own machinery;
drop the block to let the chart generate them instead.

The ServiceMonitor toggle closes the loop — Grafana's own /metrics
flow into the same stack its dashboards read from, so the pane of
glass is itself on the pane of glass.

Change first: the four placeholders (database host, credentials
Secret, stack name, root URL), then session affinity expectations —
with two replicas, in-flight sessions ride the database, so no
sticky-session configuration is needed at the ingress layer.

See [03-ha-external-db.yaml](./03-ha-external-db.yaml) for the
manifest.
