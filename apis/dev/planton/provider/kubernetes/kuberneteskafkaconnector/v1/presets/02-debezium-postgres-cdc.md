# Debezium Postgres CDC preset

A production-shaped change-data-capture source: the Debezium Postgres
connector streams row-level changes from a Postgres database into
Kafka topics (`orders.public.orders`, ...), using Postgres's built-in
`pgoutput` logical-decoding plugin so nothing has to be installed in
the database server.

Two things make this preset the reference shape:

- **The credential is a config-provider reference.** The database
  password is `${secrets:kafka/orders-db-credentials:password}` — the
  workers resolve it from the Kubernetes Secret at connector start.
  It never lands in this resource, in IaC state, or in `kubectl get`
  output. The prerequisite lives on the CONNECT cluster: enable the
  KubernetesSecretConfigProvider through its worker `config`
  (`config.providers` entries), and grant the Connect service account
  read access to the Secret.
- **The class must be on the workers.** Pair this with a
  KubernetesKafkaConnect cluster whose `image`, `plugins` or `build`
  arm carries the Debezium Postgres plugin (that kind's
  debezium-prebuilt-image and operator-built-image presets are the
  companions).

`tasks_max: 1` is not a limitation to tune away — a Postgres
replication slot is a single stream, so the connector runs one task
regardless. `auto_restart` with a raised cap keeps a flaky database
connection from parking the pipe FAILED overnight.

See [02-debezium-postgres-cdc.yaml](./02-debezium-postgres-cdc.yaml)
for the manifest.
