---
title: "Presets"
description: "Ready-to-deploy configuration presets for Karapace Schema Registry"
type: "preset-list"
componentSlug: "karapace-schema-registry"
componentTitle: "Karapace Schema Registry"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-minimal"
    rank: "01"
    title: "Minimal preset"
    excerpt: "The smallest declarable schema registry: one replica against a dev Kafka cluster's plaintext listener, everything else on defaults. The registry creates the `_schemas` topic on first start and serves..."
  - slug: "02-production-ha"
    rank: "02"
    title: "Production HA preset"
    excerpt: "The production posture: two replicas with automatic leader election (followers forward writes to the leader at its pod-IP identity — no external coordination), a SASL_SSL connection to the Kafka..."
  - slug: "03-with-rest-proxy-and-tls"
    rank: "03"
    title: "REST proxy and TLS preset"
    excerpt: "The hardened single-replica shape: the registry API served over HTTPS from a cert-manager-issued Secret, HTTP Basic authentication from a hot-reloaded authfile Secret, and the Kafka REST-proxy role..."
---

# Karapace Schema Registry Presets

Ready-to-deploy configuration presets for Karapace Schema Registry. Each preset is a complete manifest you can copy, customize, and deploy.
