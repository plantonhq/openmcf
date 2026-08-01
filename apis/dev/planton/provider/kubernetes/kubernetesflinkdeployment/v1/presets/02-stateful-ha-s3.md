# Stateful HA-on-S3 preset

The full durable-state posture, every requirement the operator
enforces made explicit:

**Savepoint upgrades.** `upgradeMode: savepoint` drains the job
through a savepoint on every spec change — the operationally safest
path — and REQUIRES both `state.checkpointsDir` and
`state.savepointsDir` (the operator rejects the deployment without
them; this spec validates the same rule at authoring time).

**JobManager standbys behind Kubernetes HA.** `replicas: 2` runs a
warm standby coordinated through HA metadata in
`highAvailability.storageDir` — required for standbys, and what lets
a lost JobManager RECOVER the job from its newest checkpoint instead
of restarting it from scratch.

**The S3 seam, Secret-native.** All three directories land on a
composed `KubernetesSeaweedFs` (referenced by name — the FK defaults
wire its S3 endpoint and its generated credentials Secret; the pods
read the keys at RUNTIME through secretKeyRefs, so nothing
credential-bearing ever renders into config). Deploy the store first,
in this same namespace — secretKeyRefs cannot cross namespaces — and
declare the bucket (`flink-state`) on the store: Flink does not
create buckets.

**The plugin truth.** The official Flink images ship the S3
filesystem plugin DISABLED; `builtinPluginJar` names the exact jar
under `/opt/flink/opt` in YOUR image (its version must match the
image's Flink patch version) — without it every `s3://` path fails at
runtime with "unsupported filesystem scheme". Custom images that bake
the plugin into `/opt/flink/plugins` leave it empty.

Change first: the store reference (`objects`) and the plugin jar's
patch version to match your image.

See [02-stateful-ha-s3.yaml](./02-stateful-ha-s3.yaml) for the
manifest.
