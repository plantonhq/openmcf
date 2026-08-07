---
title: "Minimal preset"
description: "The smallest declarable schema registry: one replica against a dev Kafka cluster's plaintext listener, everything else on defaults. The registry creates the `_schemas` topic on first start and serves..."
type: "preset"
rank: "01"
presetSlug: "01-minimal"
componentSlug: "karapace-schema-registry"
componentTitle: "Karapace Schema Registry"
provider: "kubernetes"
icon: "package"
order: 1
---

# Minimal preset

The smallest declarable schema registry: one replica against a dev
Kafka cluster's plaintext listener, everything else on defaults. The
registry creates the `_schemas` topic on first start and serves the
standard Schema Registry REST API at the exported endpoint — existing
Confluent SR client libraries, Connect converters and consoles work
against it unchanged.

The trade-offs are the dev trade-offs: the schemas topic is created
at replication factor 1 (a single broker restart makes schemas
briefly unavailable; losing the broker's disk loses them), the API is
open to anyone who can reach the Service, and the connection is
plaintext. The production-ha preset addresses all three; note that
the schemas-topic replication factor applies AT CREATION — graduating
later means reassigning the topic with Kafka tooling, so prefer
starting production registries from the production preset.

See [01-minimal.yaml](./01-minimal.yaml) for the manifest.
