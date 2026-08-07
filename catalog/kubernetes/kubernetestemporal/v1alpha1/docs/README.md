# Kubernetes Temporal — design notes

## Grain

One resource = one Temporal cluster from the official `temporal` chart
(go.temporal.io/helm-charts; the chart pin governs the app version).
The release is named after `metadata.name` and `fullnameOverride` pins
every child name (`<name>-frontend`, `<name>-history`, `<name>-web`,
...), so the exported outputs are deterministic. Names are capped at
45 characters — the chart's componentname helper truncates the
FULLNAME (not the component) to fit 63, and the longest component
suffix is `internal-frontend`; both engines fail loudly instead of
letting the chart truncate silently.

## Bring your own database — the oneof

The chart's current line ships NO bundled subcharts — no Cassandra,
Elasticsearch, Prometheus or Grafana to disable; declaring their
legacy keys makes the chart itself fail rendering. The spec models
that truth as a required `database` oneof: `postgres` (FK-composes a
`KubernetesPostgres` — host defaults to its read-write Service,
credential to its application Secret), `mysql` (FK-composes a
`KubernetesMysql`), or external `cassandra` (contact points you
operate; the catalog has no Cassandra kind). Cassandra lost visibility
support upstream in v1.21, so that arm REQUIRES a SQL `visibility`
block — a CEL rule enforces it at validation time, not at install
time.

## The credential path

Every store renders `existingSecret` + `secretKey`: the chart wires a
secretKeyRef into the server and schema-Job pods and STRIPS the
Helm-side keys before writing the server config. Because
existingSecret is always set, the chart's own per-store password
Secret (which would embed an inline password) is never created —
nothing credential-bearing appears in values, manifests, or
plan/preview output.

## Schema Jobs rendered explicitly

`createDatabase`/`manageSchema` are ALWAYS rendered on every store —
the chart's getStore helper silently defaults BOTH to true when unset,
and an unintended create-database attempt fails against
least-privilege users. The spec's `create_databases` /
`skip_schema_setup` map onto them one to one. `connectAddr` carries
`host:port` — the schema-Job env template parses host and port out of
that exact form.

## Dynamic config and archival

The typed dynamic-config subset covers the history/blob size and count
limits (each rendered as one global `{value, constraints: {}}` entry —
the server's own key format); everything else rides `helm_values`
under `server.dynamicConfig`. Archival enables the capability at the
cluster level (s3store/gstorage/filestore providers) and sets
namespace-default URIs; cloud credentials are ambient (IRSA/workload
identity) — nothing credential-bearing renders.

## Cross-engine parity

Both engines render byte-identical chart values from one shared
resolution order: typed values first, the `helm_values` escape hatch
merged over them with Helm `-f` semantics, and `fullnameOverride`
re-pinned LAST — the one deliberate exception to the escape hatch's
last-word contract, because every exported output derives from the
fullname. The 1.29-image compatibility shims are rendered OFF
explicitly (the chart defaults them on; the pinned line runs 1.31+).

## Deliberate exclusions

mTLS certificate mounts, JWT authorization, multi-cluster replication
(`clusterMetadata`), per-service scheduling and extra dynamic-config
keys — reachable through `helm_values`, never the primary interface.
Elasticsearch-backed advanced visibility is not modeled at all: the
SQL visibility store carries the search surface.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
