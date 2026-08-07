# Minimal stock Connect preset

The smallest declarable Connect cluster: one worker on the stock
image against a dev Kafka cluster, no TLS, no authentication. The
stock image carries only the MirrorMaker 2 connector classes
(MirrorSource/MirrorCheckpoint/MirrorHeartbeat — Kafka's FileStream
examples are not on the distribution's classpath) — still enough to
prove the pipe machinery end-to-end with a KubernetesKafkaConnector
(a MirrorSource self-mirror) before any plugin story exists.

The `"-1"` storage replication-factor entries are deliberate:
Connect's internal-topic default of 3 cannot be satisfied on a
single-broker cluster, and the workers wedge creating their topics.
Graduate to `"3"` (and a TLS+SCRAM connection, and a real plugin arm)
with the sibling presets before anything depends on the pipes.

The group identity (group.id and the three storage topics) defaults
from `metadata.name` — two Connect clusters sharing a Kafka cluster
must keep distinct names.

See [01-minimal-stock-connect.yaml](./01-minimal-stock-connect.yaml)
for the manifest.
