# Migrate from Confluent Cloud preset

The Confluent exit ramp: one mirror from a Confluent Cloud cluster
into a Strimzi-managed target. Confluent's client contract is SASL
PLAIN — the API key rides as the username and the API secret comes
from a Kubernetes Secret (`confluent-api-secret`, key `password`),
so no credential ever lands in the manifest. As in every migration
preset, IdentityReplicationPolicy on BOTH connectors keeps original
topic names, and the checkpoint connector's group-offset sync makes
consumer cutover a bootstrap repoint.

Notes on the connection shape:

- **The credential pair is the composition seam** — declare the API
  secret as a KubernetesSecret (or ExternalSecret) and reference it
  here; rotating the key is a Secret update, not a manifest change.
- **TLS trust**: this preset declares only the SASL block, mirroring
  the kind's full-surface development manifest. Confluent endpoints
  terminate TLS with public-CA certificates; when your posture
  requires explicit trust declaration, add a `tls` block naming a
  CA-bundle Secret, exactly as the MSK preset does.
- **Schemas do not mirror.** MirrorMaker 2 moves topics, records and
  offsets — subjects in Confluent Schema Registry move separately
  (register them into a KubernetesKarapace sibling before cutover;
  Karapace speaks the same SR API).

See
[02-migrate-from-confluent-cloud.yaml](./02-migrate-from-confluent-cloud.yaml)
for the manifest.
