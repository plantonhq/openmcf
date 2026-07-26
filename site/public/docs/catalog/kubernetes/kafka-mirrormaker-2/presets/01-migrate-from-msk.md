---
title: "Migrate from MSK preset"
description: "The MSK exit ramp: one mirror from an Amazon MSK cluster (SCRAM listener, TLS trust from a CA-bundle Secret, credentials from a Secret) into a Strimzi-managed target, with IdentityReplicationPolicy..."
type: "preset"
rank: "01"
presetSlug: "01-migrate-from-msk"
componentSlug: "kafka-mirrormaker-2"
componentTitle: "Kafka MirrorMaker 2"
provider: "kubernetes"
icon: "package"
order: 1
---

# Migrate from MSK preset

The MSK exit ramp: one mirror from an Amazon MSK cluster (SCRAM
listener, TLS trust from a CA-bundle Secret, credentials from a
Secret) into a Strimzi-managed target, with
IdentityReplicationPolicy on BOTH connectors so topics keep their
original names — the posture consumers expect after a migration.

What each piece is doing:

- **IdentityReplicationPolicy on both connectors** is the load-bearing
  recipe. The default policy prefixes mirrored topics with the source
  alias (`prod-msk.orders`); identity keeps `orders` as `orders`. It
  must be set on the checkpoint connector too — offset translation
  maps onto TARGET topic names, and a mismatch translates onto names
  that do not exist.
- **`sync.group.offsets.enabled: "true"`** writes translated consumer
  offsets into the target's group state, so at cutover consumers just
  repoint their bootstrap and resume — no reprocessing.
- **`tasks_max: 8`** sizes the source connector against the source's
  partition volume; `replicas: 3` spreads the tasks.
- **`auto_restart`** keeps a multi-day migration alive through
  transient MSK or network wobbles.

Cutover is still a step you perform: watch mirroring lag (metrics are
on), stop producers, let the mirror drain, repoint clients at the
target, then retire this resource.

See [01-migrate-from-msk.yaml](./01-migrate-from-msk.yaml) for the
manifest.
