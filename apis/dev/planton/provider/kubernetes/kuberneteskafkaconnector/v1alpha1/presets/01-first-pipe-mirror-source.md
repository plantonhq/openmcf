# First-pipe preset (stock-image MirrorSource)

The smallest real pipe that actually runs on the stock image. The
Strimzi Connect image carries ONLY the three MirrorMaker 2 connector
classes (MirrorSource, MirrorCheckpoint, MirrorHeartbeat — verified
live against the workers' own plugin listing; Kafka's FileStream
example connectors are NOT on the distribution's classpath). So the
zero-machinery first declaration is a MirrorSourceConnector mirroring
a topic from the cluster back to itself under an alias: real records
cross a real connector, and nothing has to be installed first.

The two contracts this preset teaches transfer to every real
connector: the namespace MUST be the Connect cluster's own (a
connector elsewhere is silently never reconciled), and
`connect_cluster` is the binding rendered as the
`strimzi.io/cluster` label. `auto_restart` is on even here — the
habit worth forming before the pipes matter.

Nothing else transfers: a self-mirror is a smoke test, not an
integration pattern. Graduate to the Debezium CDC preset for a
production-shaped source — its class arrives through the Connect
cluster's `image`, `plugins`, or `build` arms.

See [01-first-pipe-mirror-source.yaml](./01-first-pipe-mirror-source.yaml)
for the manifest.
