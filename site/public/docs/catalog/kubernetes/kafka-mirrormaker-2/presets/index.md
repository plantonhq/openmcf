---
title: "Presets"
description: "Ready-to-deploy configuration presets for Kafka MirrorMaker 2"
type: "preset-list"
componentSlug: "kafka-mirrormaker-2"
componentTitle: "Kafka MirrorMaker 2"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-migrate-from-msk"
    rank: "01"
    title: "Migrate from MSK preset"
    excerpt: "The MSK exit ramp: one mirror from an Amazon MSK cluster (SCRAM listener, TLS trust from a CA-bundle Secret, credentials from a Secret) into a Strimzi-managed target, with IdentityReplicationPolicy..."
  - slug: "02-migrate-from-confluent-cloud"
    rank: "02"
    title: "Migrate from Confluent Cloud preset"
    excerpt: "The Confluent exit ramp: one mirror from a Confluent Cloud cluster into a Strimzi-managed target. Confluent's client contract is SASL PLAIN — the API key rides as the username and the API secret..."
  - slug: "03-active-passive-dr"
    rank: "03"
    title: "Active-passive DR preset"
    excerpt: "Standing disaster-recovery replication between two Strimzi-managed clusters: everything on the active cluster mirrors continuously into a passive standby, with consumer-group checkpoints flowing so..."
---

# Kafka MirrorMaker 2 Presets

Ready-to-deploy configuration presets for Kafka MirrorMaker 2. Each preset is a complete manifest you can copy, customize, and deploy.
