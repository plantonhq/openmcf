---
title: "Active-passive DR preset"
description: "Standing disaster-recovery replication between two Strimzi-managed clusters: everything on the active cluster mirrors continuously into a passive standby, with consumer-group checkpoints flowing so..."
type: "preset"
rank: "03"
presetSlug: "03-active-passive-dr"
componentSlug: "kafka-mirrormaker-2"
componentTitle: "Kafka MirrorMaker 2"
provider: "kubernetes"
icon: "package"
order: 3
---

# Active-passive DR preset

Standing disaster-recovery replication between two Strimzi-managed
clusters: everything on the active cluster mirrors continuously into
a passive standby, with consumer-group checkpoints flowing so
failed-over consumers resume near where the primary's consumers left
off.

The deliberate contrast with the migration presets: this preset keeps
the DEFAULT replication policy, so mirrored topics arrive as
`primary.<topic>` on the DR cluster. For a standing mirror that is a
feature, not a wart — provenance stays visible, a future fail-back
mirror cannot loop records into the topics it reads, and DR-side
tooling can tell mirrored data from anything written locally.
Failing-over applications consume the prefixed names (or you rename
at fail-over time); if your DR doctrine demands identical names,
apply the migration presets' IdentityReplicationPolicy recipe — on
both connectors — and accept the loop-prevention trade-off.

Operational notes:

- **Checkpoints ARE the DR story.** Records without translated
  offsets mean every consumer restarts from earliest or latest at
  fail-over; `sync.group.offsets.enabled: "true"` keeps the DR
  cluster's group state warm. The 300-second refresh intervals bound
  how stale topic and group discovery can get.
- **Both connections are TLS + SCRAM** with distinct reader/writer
  principals — least privilege per side, credentials in Secrets.
- **Fail-over is a decision, not an automation.** This resource moves
  data; promoting the DR cluster, repointing clients, and standing up
  the reverse mirror are runbook steps.

See [03-active-passive-dr.yaml](./03-active-passive-dr.yaml) for the
manifest.
