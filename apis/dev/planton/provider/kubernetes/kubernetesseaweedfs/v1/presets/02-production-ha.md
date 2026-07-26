# Production HA preset

A highly available store: a 3-master Raft quorum, 3 volume servers on
explicit fast storage with the `"001"` replication code (one extra
copy of every object on another server — the store survives losing a
volume server), a dedicated 2-replica S3 gateway Deployment scaled
independently of filer metadata, resource requests and limits on
every tier, and the admin console enabled behind module-generated
credentials with persisted state.

Know where the availability boundary sits: the filer is ONE pod, by
design — its embedded leveldb metadata store is per-pod, so raising
filer replicas without first wiring a shared external store
(Postgres/MySQL through `WEED_*` environment variables and
`helm_values`) would give each filer its own divergent namespace.
Replication also is not backup: `"001"` protects against server loss,
not against deletion or corruption written through the API. And the
declared topology must fit the code — dropping to fewer than 2 volume
servers under `"001"` makes writes fail.

Change first: both `storage_class` placeholders (a literal class
name, or a valueFrom reference to a KubernetesStorageClass), then
`volume.data_volume.size` for the data you actually expect — it is
per volume-server pod, and with `"001"` every object is stored twice.
Revisit the gateway replica count once real S3 traffic exists; the
dedicated Deployment is the knob for API throughput.

See [02-production-ha.yaml](./02-production-ha.yaml) for the
manifest.
