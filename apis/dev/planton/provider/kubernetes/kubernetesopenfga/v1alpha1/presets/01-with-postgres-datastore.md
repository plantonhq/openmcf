# PostgreSQL datastore preset

The production OpenFGA shape: three stateless server replicas sharing
a PostgreSQL datastore, schema migrations running as an init container
in every pod (idempotent — `openfga migrate` gates each rollout on the
database being reachable and current), and the API guarded by
pre-shared keys.

The datastore references a `KubernetesPostgres` resource by name: the
host resolves to its read-write Service (always the current primary)
and the password rides a secretKeyRef against the operator-maintained
`<cluster>-app` credential Secret — nothing credential-bearing ever
renders into values or manifests. Declare the `openfga` database at
the Postgres resource's bootstrap (`initdb`), and co-locate OpenFGA
with the database namespace (a secretKeyRef reads only its own
namespace's Secrets).

Authentication here points at a Secret YOU maintain
(`openfga-api-keys`, comma-separated keys under the data key `keys` —
the chart's contract). Alternatively declare `authn.preshared.keys`
inline and the module materializes them into a managed Secret. Without
any `authn` block the API is OPEN to everything that can reach the
Service — fine on a lab cluster, never in production.

Authorization DATA is not deployment config: create stores, models and
tuples through the API (the `fga` CLI, or the platform's OpenFgaStore /
OpenFgaAuthorizationModel / OpenFgaRelationshipTuple resources against
the exported endpoint).

Change first: the two `my-postgres` references, then `replicas` and
`tuning` to your check volume.

See [01-with-postgres-datastore.yaml](./01-with-postgres-datastore.yaml)
for the manifest.
